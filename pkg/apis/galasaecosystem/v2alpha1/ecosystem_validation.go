package v2alpha1

import "fmt"

func Validate(c *GalasaEcosystem) error {
	if c.Spec.Hostname == "" {
		return fmt.Errorf("Hostname is a required field")
	}
	return nil
}
