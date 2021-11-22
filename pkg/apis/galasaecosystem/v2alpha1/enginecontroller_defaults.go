package v2alpha1

func DefaultEngineController() ComponentSpec {
	return ComponentSpec{
		Image:           CONTROLLERIMAGE,
		Replicas:        &SINGLEREPLICA,
		ImagePullPolicy: IMAGEPULLPOLICY,
	}
}

func SetEngineControllerDefaults(engiencontroller ComponentSpec) ComponentSpec {
	if engiencontroller.Image == "" {
		engiencontroller.Image = CONTROLLERIMAGE
	}
	if engiencontroller.ImagePullPolicy == "" {
		engiencontroller.ImagePullPolicy = IMAGEPULLPOLICY
	}
	if engiencontroller.Replicas == nil {
		engiencontroller.Replicas = &SINGLEREPLICA
	}

	return engiencontroller
}
