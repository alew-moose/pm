package pkg

import (
	"fmt"
	"regexp"
)

type PackageName string

var packageNameRe = regexp.MustCompile(`^[\w-]+$`)

func (pn PackageName) Validate() error {
	if !packageNameRe.MatchString(string(pn)) {
		return fmt.Errorf("invalid package name %q", string(pn))
	}
	return nil
}
