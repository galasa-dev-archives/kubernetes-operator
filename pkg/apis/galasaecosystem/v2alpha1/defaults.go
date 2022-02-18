/*
 * Copyright contributors to the Galasa Project
 */
package v2alpha1

const (
	IMAGEPULLPOLICY  string = "Always"
	STORAGECLASSNAME string = "default"
	CPSIMAGE         string = "quay.io/coreos/etcd:v3.4.3"
	CPSSTORAGE       string = "1Gi"
	RASIMAGE         string = "couchdb:2.3.1"
	RASSTORAGE       string = "5Gi"
	APISTORAGE       string = "1Gi"
)

var (
	GALASAVERSION   string = "latest"
	APIIMAGE        string = "icr.io/galasadev/galasa-boot-embedded-amd64:" + GALASAVERSION
	RESMONIMAGE     string = "icr.io/galasadev/galasa-boot-embedded-amd64:" + GALASAVERSION
	METRICSIMAGE    string = "icr.io/galasadev/galasa-boot-embedded-amd64:" + GALASAVERSION
	CONTROLLERIMAGE string = "icr.io/galasadev/galasa-boot-embedded-amd64:" + GALASAVERSION
	SIMBANKIMAGE    string = "icr.io/galasadev/galasa-simplatform-amd64:" + GALASAVERSION
	SINGLEREPLICA   int32  = 1
	SIMBANK         bool   = false
	SIMBANKVERSION  string = "0.15.0"
)
