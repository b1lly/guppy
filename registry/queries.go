package main

import (
	"fmt"

	"github.com/b1lly/guppy"
)

var pkgsSelectQuery = "SELECT Name, Remote, CommitHash, Version FROM guppy.packages"

func allPackages() ([]*guppy.Package, error) {
	return selectPackages(pkgsSelectQuery)
}

func packagesByName(name string) ([]*guppy.Package, error) {
	return selectPackages(fmt.Sprintf("%v WHERE Name='%v'", pkgsSelectQuery, name))
}

func selectPackages(query string) ([]*guppy.Package, error) {
	rows, err := DB.Query(query)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	var packages []*guppy.Package
	for rows.Next() {
		pkg := guppy.Package{}

		err = rows.Scan(&pkg.Name, &pkg.Remote, &pkg.CommitHash, &pkg.Version)
		packages = append(packages, &pkg)
	}

	return packages, err
}

func insertPackage(pkg *guppy.Package) error {
	query := `
		INSERT INTO guppy.packages (Name, Remote, CommitHash, Version)
		VALUES ('%v', '%v', '%v', '%v')
	`
	res, err := DB.Exec(fmt.Sprintf(query, pkg.Name, pkg.Remote, pkg.CommitHash, pkg.Version.String()))
	if err != nil {
		return err
	}

	pkg.Id, err = res.LastInsertId()
	if err != nil {
		return err
	}

	return nil
}