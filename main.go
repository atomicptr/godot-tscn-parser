// TODO: this file will be deleted at some point
package main

import (
	"fmt"
	"github.com/atomicptr/godot-tscn-parser/pkg/tscn"
)

func main() {
	scene, err := tscn.LoadFileAndParse("examples/TestFile.tscn")
	if err != nil {
		panic(err)
	}

	for _, section := range scene.Sections {
		fmt.Println("\nSection:", section.ResourceType, section.Pos)

		if len(section.Attributes) > 0 {
			fmt.Println("Attributes:")

			for _, attr := range section.Attributes {
				fmt.Printf("\t%s = %s [%s]\n", attr.Key, attr.Value.ToString(), attr.Pos)
			}
		}

		if len(section.Fields) > 0 {
			fmt.Println("Fields:")

			for _, field := range section.Fields {
				fmt.Printf("\t%s = %s [%s]\n", field.Key, field.Value.ToString(), field.Pos)
			}
		}
	}
}
