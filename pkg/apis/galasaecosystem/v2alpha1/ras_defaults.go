package v2alpha1

func DefaultRas() ComponentSpec {
	return ComponentSpec{
		Image:            RASIMAGE,
		Replicas:         &SINGLEREPLICA,
		ImagePullPolicy:  IMAGEPULLPOLICY,
		Storage:          RASSTORAGE,
		StorageClassName: STORAGECLASSNAME,
	}
}

func SetRasDefaults(ras ComponentSpec) ComponentSpec {
	if ras.Image == "" {
		ras.Image = RASIMAGE
	}
	if ras.ImagePullPolicy == "" {
		ras.ImagePullPolicy = IMAGEPULLPOLICY
	}
	if ras.Storage == "" {
		ras.Storage = RASSTORAGE
	}
	if ras.StorageClassName == "" {
		ras.StorageClassName = STORAGECLASSNAME
	}
	if ras.Replicas == nil {
		ras.Replicas = &SINGLEREPLICA
	}

	return ras
}
