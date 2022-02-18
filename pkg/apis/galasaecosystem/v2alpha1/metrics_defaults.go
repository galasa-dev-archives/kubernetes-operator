/*
 * Copyright contributors to the Galasa Project
 */
package v2alpha1

func DefaultMetrics(c *GalasaEcosystemSpec) ComponentSpec {
	return ComponentSpec{
		Image:           METRICSIMAGE,
		Replicas:        &SINGLEREPLICA,
		ImagePullPolicy: c.ImagePullPolicy,
		NodeSelector:    c.NodeSelector,
	}
}

func SetMetricsDefaults(metrics ComponentSpec, c *GalasaEcosystemSpec) ComponentSpec {
	if metrics.Image == "" {
		metrics.Image = METRICSIMAGE
	}
	if metrics.ImagePullPolicy == "" {
		metrics.ImagePullPolicy = c.ImagePullPolicy
	}
	if metrics.Replicas == nil {
		metrics.Replicas = &SINGLEREPLICA
	}
	if metrics.NodeSelector == nil {
		metrics.NodeSelector = c.NodeSelector
	}

	return metrics
}
