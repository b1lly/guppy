package main

import (
	"fmt"

	"github.com/b1lly/guppy"
)

var pkgsSelectQuery = "SELECT Name, Remote, CommitHash, Version FROM guppy.packages"

func QuerySelectPackages(query string) ([]*guppy.Package, error) {
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

func AllPackages() ([]*guppy.Package, error) {
	return QuerySelectPackages(pkgsSelectQuery)
}

func PackageByName(name string) ([]*guppy.Package, error) {
	return QuerySelectPackages(fmt.Sprintf("%v WHERE Name='%v'", pkgsSelectQuery, name))
}

