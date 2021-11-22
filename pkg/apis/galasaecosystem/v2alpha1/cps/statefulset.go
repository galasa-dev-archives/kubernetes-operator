package cps

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func (c *Cps) getStatefulSet() *appsv1.StatefulSet {
	labels := map[string]string{
		"app": c.Name,
	}
	s := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:            c.Name,
			Namespace:       c.Namespace,
			Labels:          labels,
			OwnerReferences: c.Owner,
		},
		Spec: appsv1.StatefulSetSpec{
			ServiceName: c.Name + "-internal-service",
			Replicas:    c.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:   c.Name,
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					NodeSelector: c.NodeSelector,
					Containers: []corev1.Container{
						{
							Name:            "etcd",
							Image:           c.Image,
							ImagePullPolicy: corev1.PullPolicy(c.ImagePullPolicy),
							Ports: []corev1.ContainerPort{
								{
									Name:          "peer",
									ContainerPort: int32(2379),
								},
								{
									Name:          "client",
									ContainerPort: int32(2380),
								},
							},
							LivenessProbe: &corev1.Probe{
								InitialDelaySeconds: 60,
								PeriodSeconds:       60,
								Handler: corev1.Handler{
									TCPSocket: &corev1.TCPSocketAction{
										Port: intstr.FromInt(2379),
									},
								},
							},
							Env: []corev1.EnvVar{
								{
									Name:  "INITIAL_CLUSTER_SIZE",
									Value: "1",
								},
								{
									Name:  "SET_NAME",
									Value: c.Name,
								},
								{
									Name:  "SERVICE",
									Value: c.Name + "-internal-service",
								},
								{
									Name:  "ETCDCTL_API",
									Value: "3",
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "cps-datadir",
									MountPath: "/var/run/etcd",
								},
							},
							Lifecycle: &corev1.Lifecycle{
								PreStop: &corev1.Handler{
									Exec: &corev1.ExecAction{
										Command: []string{
											"/bin/sh",
											"-ec",
											`
											EPS=""
											for i in $(seq 0 $((${INITIAL_CLUSTER_SIZE} - 1))); do
												EPS="${EPS}${EPS:+,}http://${SET_NAME}-${i}.${SERVICE}:2379"
											done
											HOSTNAME=$(hostname)
											member_hash() {
												etcdctl member list | grep http://${HOSTNAME}:2380 | cut -d':' -f1 | cut -d'[' -f1
											}
											SET_ID=${HOSTNAME##*[^0-9]}
											if [ "${SET_ID}" -ge ${INITIAL_CLUSTER_SIZE} ]; then
												echo "Removing ${HOSTNAME} from etcd cluster"
												ETCDCTL_ENDPOINT=${EPS} etcdctl member remove $(member_hash)
												if [ $? -eq 0 ]; then
													# Remove everything otherwise the cluster will no longer scale-up
													rm -rf /var/run/etcd/*
												fi
											fi
											`,
										},
									},
								},
							},
							Command: []string{
								"/bin/sh",
								"-ec",
								`
								HOSTNAME=$(hostname)
								# store member id into PVC for later member replacement
								collect_member() {
			
									#while ! etcdctl member list &>/dev/null; do sleep 1; done
									etcdctl member list | grep http://${HOSTNAME}:2380 | awk -F',' '{print $1}' > /var/run/etcd/member_id
									exit 0
								}
								eps() {
									EPS=""
									for i in $(seq 0 $((${INITIAL_CLUSTER_SIZE} - 1))); do
										EPS="${EPS}${EPS:+,}http://${SET_NAME}-${i}:2379"
									done
									echo ${EPS}
								}
								member_hash() {
									etcdctl member list | grep http://${HOSTNAME}:2380 | awk -F',' '{print $1}'
								}
								# we should wait for other pods to be up before trying to join
								# otherwise we got "no such host" errors when trying to resolve other members
								for i in $(seq 0 $((${INITIAL_CLUSTER_SIZE} - 1))); do
									while true; do
										echo "Waiting for ${SET_NAME}-${i}.${SERVICE} to come up"
										ping -W 1 -c 1 ${SET_NAME}-${i}.${SERVICE} > /dev/null && break
										sleep 1s
									done
								done
								# re-joining after failure?
								if [ -e /var/run/etcd/default.etcd ]; then
									echo "Re-joining etcd member"
									member_id=$(cat /var/run/etcd/member_id)
									# re-join member
									ETCDCTL_ENDPOINT=$(eps) etcdctl member update ${member_id} http://${HOSTNAME}.${SERVICE}:2380 | true
									exec etcd --name ${HOSTNAME} \
										--listen-peer-urls http://0.0.0.0:2380 \
										--listen-client-urls http://0.0.0.0:2379\
										--advertise-client-urls http://${HOSTNAME}.${SERVICE}:2379 \
										--data-dir /var/run/etcd/default.etcd
								fi
								# etcd-SET_ID
								SET_ID=${HOSTNAME##*[^0-9]}
								# adding a new member to existing cluster (assuming all initial pods are available)
								if [ "${SET_ID}" -ge ${INITIAL_CLUSTER_SIZE} ]; then
									export ETCDCTL_ENDPOINT=$(eps)
									# member already added?
									MEMBER_HASH=$(member_hash)
									if [ -n "${MEMBER_HASH}" ]; then
										# the member hash exists but for some reason etcd failed
										# as the datadir has not be created, we can remove the member
										# and retrieve new hash
										etcdctl member remove ${MEMBER_HASH}
									fi
									echo "Adding new member"
									etcdctl member add ${HOSTNAME} http://${HOSTNAME}.${SERVICE}:2380 | grep "^ETCD_" > /var/run/etcd/new_member_envs
									if [ $? -ne 0 ]; then
										echo "Exiting"
										rm -f /var/run/etcd/new_member_envs
										exit 1
									fi
									cat /var/run/etcd/new_member_envs
									source /var/run/etcd/new_member_envs
									echo "Collect 1"
									collect_member &
									exec etcd --name ${HOSTNAME} \
										--listen-peer-urls http://0.0.0.0:2380 \
										--listen-client-urls http://0.0.0.0:2379 \
										--advertise-client-urls http://${HOSTNAME}.${SERVICE}:2379 \
										--data-dir /var/run/etcd/default.etcd \
										--initial-advertise-peer-urls http://${HOSTNAME}.${SERVICE}:2380 \
										--initial-cluster ${ETCD_INITIAL_CLUSTER} \
										--initial-cluster-state ${ETCD_INITIAL_CLUSTER_STATE}
								fi
								PEERS=""
								for i in $(seq 0 $((${INITIAL_CLUSTER_SIZE} - 1))); do
									PEERS="${PEERS}${PEERS:+,}${SET_NAME}-${i}=http://${SET_NAME}-${i}:2380"
								done
								echo "Collect 2"
								collect_member &
								# join member
								exec etcd --name ${HOSTNAME} \
									--initial-advertise-peer-urls http://${HOSTNAME}:2380 \
									--listen-peer-urls http://0.0.0.0:2380 \
									--listen-client-urls http://0.0.0.0:2379 \
									--advertise-client-urls http://${HOSTNAME}:2379 \
									--initial-cluster-token etcd-cluster-1 \
									--initial-cluster ${PEERS} \
									--initial-cluster-state new \
									--data-dir /var/run/etcd/default.etcd
								`,
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "cps-datadir",
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: c.Name + "-pvc",
								},
							},
						},
					},
				},
			},
		},
	}
	s.SetGroupVersionKind(schema.FromAPIVersionAndKind("apps/v1", "StatefulSet"))
	return s
}
