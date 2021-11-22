package v2alpha1

func DefaultResmon() ComponentSpec {
	return ComponentSpec{
		Image:           RESMONIMAGE,
		Replicas:        &SINGLEREPLICA,
		ImagePullPolicy: IMAGEPULLPOLICY,
	}
}

func SetResmonDefaults(resmon ComponentSpec) ComponentSpec {
	if resmon.Image == "" {
		resmon.Image = RESMONIMAGE
	}
	if resmon.ImagePullPolicy == "" {
		resmon.ImagePullPolicy = IMAGEPULLPOLICY
	}
	if resmon.Replicas == nil {
		resmon.Replicas = &SINGLEREPLICA
	}

	return resmon
}
