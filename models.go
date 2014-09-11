package guppy

import (
	"fmt"
	"strconv"
	"strings"
)

type PackageError struct {
	Msg string
}

func (e PackageError) Error() string {
	return fmt.Sprintf("%v", e.Msg)
}

type Package struct {
	Id int64
	Name       string
	Version    *Version
	Remote     string
	CommitHash string
}

func NewPackage(name, version, remote, hash string) (*Package, error) {
	if name == "" {
		return nil, PackageError{"no package name provided"}
	}

	if remote == "" {
		return nil, PackageError{"no remote name specified"}
	}

	if hash == "" {
		return nil, PackageError{"no commit hash provided"}
	}

	// TODO(billy) Validate that repository & commit hash exist

	return &Package{0, name, NewVersion(version), remote, hash}, nil
}

type Version struct {
	Major int
	Minor int
	Patch int
}

func (v *Version) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}

func (v *Version) Scan(val interface{}) error {
	*v = *NewVersion(string(val.([]uint8)))
	return nil
}

// NewVersion will take a string representation of versions and convert it to
// a Version struct. It fills in empty fields with the value of '0'
func NewVersion(version string) *Version {
	segments := strings.Split(version, ".")
	if len(segments) == 0 {
		return &Version{0, 0, 0}
	}

	// Convert our string to integers for storage
	var vers []int
	for _, seg := range segments {
		i, err := strconv.Atoi(seg)

		if err != nil {
			vers = append(vers, 0)
			continue
		}
		vers = append(vers, i)
	}

	// Fill in remaining version fields with 0 (if necessary)
	for i := len(vers); i < 3; i++ {
		vers = append(vers, 0)
	}

	return &Version{vers[0], vers[1], vers[2]}
}