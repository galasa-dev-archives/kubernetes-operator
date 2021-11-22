package v2alpha1

const (
	IMAGEPULLPOLICY  string = "Always"
	STORAGECLASSNAME string = "default"
	CPSIMAGE         string = "quay.io/coreos/etcd:v3.4.3"
	CPSSTORAGE       string = "1Gi"
	RASIMAGE         string = "couchdb:2.3.1"
	RASSTORAGE       string = "5Gi"
	APISTORAGE       string = "1Gi"
)

var (
	GALASAVERSION   string = "0.18.0"
	APIIMAGE        string = "docker.galasa.dev/galasadev/galasa-boot-embedded-amd64:" + GALASAVERSION
	RESMONIMAGE     string = "docker.galasa.dev/galasadev/galasa-boot-embedded-amd64:" + GALASAVERSION
	METRICSIMAGE    string = "docker.galasa.dev/galasadev/galasa-boot-embedded-amd64:" + GALASAVERSION
	CONTROLLERIMAGE string = "docker.galasa.dev/galasadev/galasa-boot-embedded-amd64:" + GALASAVERSION
	SINGLEREPLICA   int32  = 1
)
