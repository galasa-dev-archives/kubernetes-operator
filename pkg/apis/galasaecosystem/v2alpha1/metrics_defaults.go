package v2alpha1

func DefaultMetrics() ComponentSpec {
	return ComponentSpec{
		Image:           METRICSIMAGE,
		Replicas:        &SINGLEREPLICA,
		ImagePullPolicy: IMAGEPULLPOLICY,
	}
}

func SetMetricsDefaults(metrics ComponentSpec) ComponentSpec {
	if metrics.Image == "" {
		metrics.Image = METRICSIMAGE
	}
	if metrics.ImagePullPolicy == "" {
		metrics.ImagePullPolicy = IMAGEPULLPOLICY
	}
	if metrics.Replicas == nil {
		metrics.Replicas = &SINGLEREPLICA
	}

	return metrics
}
