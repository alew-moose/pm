package pkg

import (
	"fmt"

	"github.com/alew-moose/pm/internal/version"
)

type PackageVersionConstraint struct {
	Name              PackageName               `json:"name" yaml:"name"`
	VersionConstraint version.VersionConstraint `json:"ver" yaml:"ver"`
}

func (pvs PackageVersionConstraint) String() string {
	return fmt.Sprintf("%s(ver %s)", pvs.Name, pvs.VersionConstraint)
}

func (pvs *PackageVersionConstraint) Match(pv PackageVersion) bool {
	return pvs.Name == pv.Name && pvs.VersionConstraint.Match(pv.Version)
}

// TODO: remove?
func (pvs PackageVersionConstraint) Validate() error {
	if err := pvs.Name.Validate(); err != nil {
		return err
	}
	// TODO: uncomment + implement?
	// if err := pvs.VersionConstraint.Validate(); err != nil {
	// 	return fmt.Errorf("invalid version spec: %s", err)
	// }
	return nil
}
