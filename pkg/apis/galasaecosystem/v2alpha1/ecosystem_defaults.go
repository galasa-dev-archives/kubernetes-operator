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
	if len(c.ComponentsSpec) == 0 {
		c.ComponentsSpec = map[string]ComponentSpec{}
	}
	if cps, exists := c.ComponentsSpec["cps"]; !exists {
		c.ComponentsSpec["cps"] = DefaultCps()
	} else {
		c.ComponentsSpec["cps"] = SetCpsDefaults(cps)
	}
	if ras, exists := c.ComponentsSpec["ras"]; !exists {
		c.ComponentsSpec["ras"] = DefaultRas()
	} else {
		c.ComponentsSpec["ras"] = SetRasDefaults(ras)
	}
	if api, exists := c.ComponentsSpec["api"]; !exists {
		c.ComponentsSpec["api"] = DefaultApi()
	} else {
		c.ComponentsSpec["api"] = SetApiDefaults(api)
	}
	if enginecontroller, exists := c.ComponentsSpec["enginecontroller"]; !exists {
		c.ComponentsSpec["enginecontroller"] = DefaultEngineController()
	} else {
		c.ComponentsSpec["enginecontroller"] = SetEngineControllerDefaults(enginecontroller)
	}
	if metrics, exists := c.ComponentsSpec["metrics"]; !exists {
		c.ComponentsSpec["metrics"] = DefaultMetrics()
	} else {
		c.ComponentsSpec["metrics"] = SetMetricsDefaults(metrics)
	}
	if resmon, exists := c.ComponentsSpec["resmon"]; !exists {
		c.ComponentsSpec["resmon"] = DefaultResmon()
	} else {
		c.ComponentsSpec["resmon"] = SetResmonDefaults(resmon)
	}
}
