/*
 * Copyright contributors to the Galasa Project
 */
package v2alpha1

func DefaultRas(c *GalasaEcosystemSpec) ComponentSpec {
	return ComponentSpec{
		Image:            RASIMAGE,
		Replicas:         &SINGLEREPLICA,
		ImagePullPolicy:  c.ImagePullPolicy,
		Storage:          RASSTORAGE,
		StorageClassName: c.StorageClassName,
		NodeSelector:     c.NodeSelector,
	}
}

func SetRasDefaults(ras ComponentSpec, c *GalasaEcosystemSpec) ComponentSpec {
	if ras.Image == "" {
		ras.Image = RASIMAGE
	}
	if ras.ImagePullPolicy == "" {
		ras.ImagePullPolicy = c.ImagePullPolicy
	}
	if ras.Storage == "" {
		ras.Storage = RASSTORAGE
	}
	if ras.StorageClassName == "" {
		ras.StorageClassName = c.StorageClassName
	}
	if ras.Replicas == nil {
		ras.Replicas = &SINGLEREPLICA
	}
	if ras.NodeSelector == nil {
		ras.NodeSelector = c.NodeSelector
	}

	return ras
}
