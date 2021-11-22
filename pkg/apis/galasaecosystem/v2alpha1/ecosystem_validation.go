/*
 * Copyright contributors to the Galasa Project
 */
package v2alpha1

import "fmt"

func Validate(c *GalasaEcosystem) error {
	if c.Spec.Hostname == "" {
		return fmt.Errorf("hostname is a required field")
	}
	return nil
}
