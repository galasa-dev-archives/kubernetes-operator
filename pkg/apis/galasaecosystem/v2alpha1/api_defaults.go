package v2alpha1

func DefaultApi() ComponentSpec {
	return ComponentSpec{
		Image:            APIIMAGE,
		Replicas:         &SINGLEREPLICA,
		ImagePullPolicy:  IMAGEPULLPOLICY,
		Storage:          APISTORAGE,
		StorageClassName: STORAGECLASSNAME,
	}
}

func SetApiDefaults(api ComponentSpec) ComponentSpec {
	if api.Image == "" {
		api.Image = APIIMAGE
	}
	if api.ImagePullPolicy == "" {
		api.ImagePullPolicy = IMAGEPULLPOLICY
	}
	if api.Storage == "" {
		api.Storage = APISTORAGE
	}
	if api.StorageClassName == "" {
		api.StorageClassName = STORAGECLASSNAME
	}
	if api.Replicas == nil {
		api.Replicas = &SINGLEREPLICA
	}

	return api
}
