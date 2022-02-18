/*
 * Copyright contributors to the Galasa Project
 */
package v2alpha1

func DefaultEngineController(c *GalasaEcosystemSpec) ComponentSpec {
	return ComponentSpec{
		Image:           CONTROLLERIMAGE,
		Replicas:        &SINGLEREPLICA,
		ImagePullPolicy: c.ImagePullPolicy,
		NodeSelector:    c.NodeSelector,
	}
}

func SetEngineControllerDefaults(enginecontroller ComponentSpec, c *GalasaEcosystemSpec) ComponentSpec {
	if enginecontroller.Image == "" {
		enginecontroller.Image = CONTROLLERIMAGE
	}
	if enginecontroller.ImagePullPolicy == "" {
		enginecontroller.ImagePullPolicy = c.ImagePullPolicy
	}
	if enginecontroller.Replicas == nil {
		enginecontroller.Replicas = &SINGLEREPLICA
	}
	if enginecontroller.NodeSelector == nil {
		enginecontroller.NodeSelector = c.NodeSelector
	}

	return enginecontroller
}
