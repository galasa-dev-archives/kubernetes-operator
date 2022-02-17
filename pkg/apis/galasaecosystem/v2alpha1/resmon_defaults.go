/*
 * Copyright contributors to the Galasa Project
 */
package v2alpha1

func DefaultResmon(c *GalasaEcosystemSpec) ComponentSpec {
	return ComponentSpec{
		Image:           RESMONIMAGE,
		Replicas:        &SINGLEREPLICA,
		ImagePullPolicy: c.ImagePullPolicy,
		NodeSelector:    c.NodeSelector,
	}
}

func SetResmonDefaults(resmon ComponentSpec, c *GalasaEcosystemSpec) ComponentSpec {
	if resmon.Image == "" {
		resmon.Image = RESMONIMAGE
	}
	if resmon.ImagePullPolicy == "" {
		resmon.ImagePullPolicy = c.ImagePullPolicy
	}
	if resmon.Replicas == nil {
		resmon.Replicas = &SINGLEREPLICA
	}
	if resmon.NodeSelector == nil {
		resmon.NodeSelector = c.NodeSelector
	}

	return resmon
}
