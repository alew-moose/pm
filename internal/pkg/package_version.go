package pkg

import (
	"fmt"
	"regexp"

	"github.com/alew-moose/pm/internal/version"
)

type PackageVersion struct {
	Name    PackageName
	Version version.Version
}

func (pv PackageVersion) String() string {
	return fmt.Sprintf("%s-%s", pv.Name, pv.Version)
}

func (pv PackageVersion) Validate() error {
	if err := pv.Name.Validate(); err != nil {
		return err
	}
	if err := pv.Version.Validate(); err != nil {
		return err
	}
	return nil
}

var packageVersionRe = regexp.MustCompile(`^(.+)-(.+)$`)

func PackageVersionFromString(s string) (PackageVersion, error) {
	var packageVersion PackageVersion
	var err error

	matches := packageVersionRe.FindStringSubmatch(s)
	if len(matches) != 3 {
		return packageVersion, fmt.Errorf("invalid package version %q", s)
	}

	packageVersion.Name = PackageName(matches[1])
	if err := packageVersion.Name.Validate(); err != nil {
		return packageVersion, err
	}

	packageVersion.Version, err = version.VersionFromString(matches[2])
	if err != nil {
		return packageVersion, fmt.Errorf("invalid version: %q", err)
	}

	if err := packageVersion.Validate(); err != nil {
		return packageVersion, fmt.Errorf("invalid package version %q", s)
	}

	return packageVersion, nil
}
