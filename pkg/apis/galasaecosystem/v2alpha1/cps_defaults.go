package v2alpha1

func DefaultCps() ComponentSpec {
	return ComponentSpec{
		Image:            CPSIMAGE,
		Replicas:         &SINGLEREPLICA,
		ImagePullPolicy:  IMAGEPULLPOLICY,
		Storage:          CPSSTORAGE,
		StorageClassName: STORAGECLASSNAME,
	}
}

func SetCpsDefaults(cps ComponentSpec) ComponentSpec {
	if cps.Image == "" {
		cps.Image = CPSIMAGE
	}
	if cps.ImagePullPolicy == "" {
		cps.ImagePullPolicy = IMAGEPULLPOLICY
	}
	if cps.Storage == "" {
		cps.Storage = CPSSTORAGE
	}
	if cps.StorageClassName == "" {
		cps.StorageClassName = STORAGECLASSNAME
	}
	if cps.Replicas == nil {
		cps.Replicas = &SINGLEREPLICA
	}

	return cps
}
