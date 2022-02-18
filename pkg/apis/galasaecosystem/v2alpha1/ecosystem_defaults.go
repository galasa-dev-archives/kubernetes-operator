/*
 * Copyright contributors to the Galasa Project
 */
package v2alpha1

import (
	"context"
)

func (c *GalasaEcosystem) SetDefaults(ctx context.Context) {
	c.Spec.SetDefaults(ctx)
}

func (c *GalasaEcosystemSpec) SetDefaults(ctx context.Context) {
	if c.BusyboxImage == "" {
		c.BusyboxImage = "busybox:latest"
	}
	if c.GalasaVersion == "" {
		c.GalasaVersion = GALASAVERSION
	} else {
		GALASAVERSION = c.GalasaVersion
	}
	if c.SimbankVersion == "" {
		c.SimbankVersion = SIMBANKVERSION
	} else {
		SIMBANKVERSION = c.SimbankVersion
	}
	if len(c.ComponentsSpec) == 0 {
		c.ComponentsSpec = map[string]ComponentSpec{}
	}
	if cps, exists := c.ComponentsSpec["cps"]; !exists {
		c.ComponentsSpec["cps"] = DefaultCps(c)
	} else {
		c.ComponentsSpec["cps"] = SetCpsDefaults(cps, c)
	}
	if ras, exists := c.ComponentsSpec["ras"]; !exists {
		c.ComponentsSpec["ras"] = DefaultRas(c)
	} else {
		c.ComponentsSpec["ras"] = SetRasDefaults(ras, c)
	}
	if api, exists := c.ComponentsSpec["api"]; !exists {
		c.ComponentsSpec["api"] = DefaultApi(c)
	} else {
		c.ComponentsSpec["api"] = SetApiDefaults(api, c)
	}
	if enginecontroller, exists := c.ComponentsSpec["enginecontroller"]; !exists {
		c.ComponentsSpec["enginecontroller"] = DefaultEngineController(c)
	} else {
		c.ComponentsSpec["enginecontroller"] = SetEngineControllerDefaults(enginecontroller, c)
	}
	if metrics, exists := c.ComponentsSpec["metrics"]; !exists {
		c.ComponentsSpec["metrics"] = DefaultMetrics(c)
	} else {
		c.ComponentsSpec["metrics"] = SetMetricsDefaults(metrics, c)
	}
	if resmon, exists := c.ComponentsSpec["resmon"]; !exists {
		c.ComponentsSpec["resmon"] = DefaultResmon(c)
	} else {
		c.ComponentsSpec["resmon"] = SetResmonDefaults(resmon, c)
	}
	if simbank, exists := c.ComponentsSpec["simbankSpec"]; !exists {
		c.ComponentsSpec["simbankSpec"] = DefaultSimbank(c)
	} else {
		c.ComponentsSpec["simbankSpec"] = SetSimbankDefaults(simbank, c)
	}
}
