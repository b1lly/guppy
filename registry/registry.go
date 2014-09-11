package main

import (
	"fmt"
	"github.com/b1lly/guppy"
)

type Registry struct {
	Packages map[string][]*guppy.Package
}

func NewRegistry() *Registry {
	return &Registry{make(map[string][]*guppy.Package)}
}

func (r *Registry) Add(pkg *guppy.Package) error {
	existingPkg := r.GetByNameAndVersion(pkg.Name, pkg.Version)
	if existingPkg != nil {
		return guppy.PackageError{fmt.Sprintf("Package %s already exists", pkg.Name)}
	}

	if r.SavePkg(pkg) != nil {
		return guppy.PackageError{"There was a problem saving package to registry, please try again."}
	}

	r.Packages[pkg.Name] = append(r.Packages[pkg.Name], pkg)
	return nil
}

func (r *Registry) Remove(pkg *guppy.Package) {

}

func (r *Registry) SavePkg(pkg *guppy.Package) error {
	return nil
}

func (r *Registry) GetAllByName(pkgName string) ([]*guppy.Package, bool) {
	pkgs, ok := r.Packages[pkgName]
	return pkgs, ok
}

func (r *Registry) GetByNameAndVersion(name string, version *guppy.Version) *guppy.Package {
	var (
		pkgs []*guppy.Package
		ok   bool
	)

	if pkgs, ok = r.GetAllByName(name); !ok {
		return nil
	}

	for _, p := range pkgs {
		if p.Version.Major == version.Major &&
			p.Version.Minor == version.Minor {
			return p
		}
	}

	return nil
}