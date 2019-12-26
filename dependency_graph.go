package gogurt

import (
	"fmt"
)

func DependencyGraph(packageMap map[string]Package) {
	fmt.Println("digraph packages {")
	for _, p := range packageMap {
		for _, dependency := range p.Dependencies() {
			fmt.Printf("\t\"%s\" -> \"%s\";\n", p.Name(), dependency.Name())
		}
	}
	fmt.Println("}")
}
