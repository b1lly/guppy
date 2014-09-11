package main

import (
	"fmt"
	"github.com/b1lly/guppy"
)

type Registry struct {
	Packages map[string][]*guppy.Package
}

// NewRegistry returns an empty registry that you can use to manage your packages
func NewRegistry() *Registry {
	return &Registry{make(map[string][]*guppy.Package)}
}

// Load will fetch all of the packages from the database and update the local
// cache to be the same
func (r *Registry) Load() error {
	pkgs, err := allPackages()
	if err != nil {
		return err
	}

	// Empty our registry and update it with the results from the database
	r.Packages = make(map[string][]*guppy.Package)
	for _, pkg := range pkgs {
		r.Packages[pkg.Name] = append(r.Packages[pkg.Name], pkg)
	}

	return nil
}

// Add is used to add a package to the registry, and also saves it to the database
func (r *Registry) Add(pkg *guppy.Package) error {
	existingPkg := r.PackageByNameAndVersion(pkg.Name, pkg.Version)
	if existingPkg != nil {
		return guppy.PackageError{fmt.Sprintf("Package %s already exists", pkg.Name)}
	}

	err := r.SavePkg(pkg)
	if err != nil {
		fmt.Println(err)
		return guppy.PackageError{"There was a problem saving package to registry, please try again."}
	}

	r.Packages[pkg.Name] = append(r.Packages[pkg.Name], pkg)
	return nil
}

func (r *Registry) Remove(pkg *guppy.Package) {

}

// SavePkg will save the package to the database and update the package id based on
// the MySQL auto_increment value returned after insert
func (r *Registry) SavePkg(pkg *guppy.Package) error {
	err := insertPackage(pkg)
	if err != nil {
		return guppy.PackageError{err.Error()}
	}

	return nil
}

// PackagesByName will return all the packages by package name from the cache
func (r *Registry) PackagesByName(pkgName string) ([]*guppy.Package, bool) {
	pkgs, ok := r.Packages[pkgName]
	return pkgs, ok
}

// PackageByNameAndVersion will look into the cache for a particular package name/version
func (r *Registry) PackageByNameAndVersion(name string, version *guppy.Version) *guppy.Package {
	var (
		pkgs []*guppy.Package
		ok   bool
	)

	if pkgs, ok = r.PackagesByName(name); !ok {
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