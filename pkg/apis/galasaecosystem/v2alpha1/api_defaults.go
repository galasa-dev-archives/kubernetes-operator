/*
 * Copyright contributors to the Galasa Project
 */
package v2alpha1

func DefaultApi(c *GalasaEcosystemSpec) ComponentSpec {
	return ComponentSpec{
		Image:            APIIMAGE,
		Replicas:         &SINGLEREPLICA,
		ImagePullPolicy:  c.ImagePullPolicy,
		Storage:          APISTORAGE,
		StorageClassName: c.StorageClassName,
		NodeSelector:     c.NodeSelector,
	}
}

func SetApiDefaults(api ComponentSpec, c *GalasaEcosystemSpec) ComponentSpec {
	if api.Image == "" {
		api.Image = APIIMAGE
	}
	if api.ImagePullPolicy == "" {
		api.ImagePullPolicy = c.ImagePullPolicy
	}
	if api.Storage == "" {
		api.Storage = APISTORAGE
	}
	if api.StorageClassName == "" {
		api.StorageClassName = c.StorageClassName
	}
	if api.Replicas == nil {
		api.Replicas = &SINGLEREPLICA
	}
	if api.NodeSelector == nil {
		api.NodeSelector = c.NodeSelector
	}

	return api
}
