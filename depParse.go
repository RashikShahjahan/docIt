package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func depParse(dirPath string) {
	g := graph.NewGraph()
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() || path == dirPath {
			return nil
		}
		pkgPath, err := filepath.Rel(dirPath, path)
		if err != nil {
			return err
		}
		pkg := graph.NewPackage(pkgPath)
		if err := g.AddPackage(pkg); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		fmt.Printf("Error generating dependency graph: %v\n", err)
	}

	// Perform topological sort on the graph
	order, err := g.TopologicalSort()
	if err != nil {
		fmt.Printf("Error performing topological sort: %v\n", err)
	}

	// Print the packages in topological order
	fmt.Println("Packages in topological order:")
	for _, pkg := range order {
		fmt.Printf("- %s\n", pkg.ImportPath())
	}
}
