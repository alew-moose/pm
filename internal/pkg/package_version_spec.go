package pkg

import (
	"fmt"

	"github.com/alew-moose/pm/internal/version"
)

type PackageVersionSpec struct {
	Name        PackageName         `json:"name" yaml:"name"`
	VersionSpec version.VersionSpec `json:"ver" yaml:"ver"`
}

func (pvs PackageVersionSpec) String() string {
	return fmt.Sprintf("%s(ver %s)", pvs.Name, pvs.VersionSpec)
}

func (pvs *PackageVersionSpec) Match(pv PackageVersion) bool {
	return pvs.Name == pv.Name && pvs.VersionSpec.Match(pv.Version)
}

func (pvs PackageVersionSpec) Validate() error {
	if err := pvs.Name.Validate(); err != nil {
		return err
	}
	if err := pvs.VersionSpec.Validate(); err != nil {
		return fmt.Errorf("invalid version spec: %s", err)
	}
	return nil
}
