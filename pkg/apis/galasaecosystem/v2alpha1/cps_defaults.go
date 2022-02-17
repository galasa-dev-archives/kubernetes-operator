/*
 * Copyright contributors to the Galasa Project
 */
package v2alpha1

func DefaultCps(c *GalasaEcosystemSpec) ComponentSpec {
	return ComponentSpec{
		Image:            CPSIMAGE,
		Replicas:         &SINGLEREPLICA,
		ImagePullPolicy:  c.ImagePullPolicy,
		Storage:          CPSSTORAGE,
		StorageClassName: c.StorageClassName,
		NodeSelector:     c.NodeSelector,
	}
}

func SetCpsDefaults(cps ComponentSpec, c *GalasaEcosystemSpec) ComponentSpec {
	if cps.Image == "" {
		cps.Image = CPSIMAGE
	}
	if cps.ImagePullPolicy == "" {
		cps.ImagePullPolicy = c.ImagePullPolicy
	}
	if cps.Storage == "" {
		cps.Storage = CPSSTORAGE
	}
	if cps.StorageClassName == "" {
		cps.StorageClassName = c.StorageClassName
	}
	if cps.Replicas == nil {
		cps.Replicas = &SINGLEREPLICA
	}
	if cps.NodeSelector == nil {
		cps.NodeSelector = c.NodeSelector
	}

	return cps
}
