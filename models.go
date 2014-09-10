package guppy

import (
	"strconv"
	"strings"
)

type Package struct {
	Name       string
	Version    Version
	Remote     string
	CommitHash string
}

type Version struct {
	Major int
	Minor int
	Patch int
}

// Set will take a string representation of versions and convert it to
// a Version struct. It fills in empty fields with the value of 0
func (v *Version) Set(version string) *Version {
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
	for i := len(vers); i < 4; i++ {
		vers = append(vers, 0)
	}

	return &Version{vers[0], vers[1], vers[2]}
}
