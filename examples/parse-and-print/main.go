package main

import (
	"fmt"
	"github.com/atomicptr/godot-tscn-parser/pkg/tscn"
)

func main() {
	scene, err := tscn.LoadFileAndParse("./TestFile.tscn")
	if err != nil {
		panic(err)
	}

	if scene.Key != "" {
		fmt.Println("File Descriptor Type:", scene.Key, "[", scene.Pos, "]")
	}

	if len(scene.Attributes) > 0 {
		fmt.Println("Attributes:")

		for _, attribute := range scene.Attributes {
			fmt.Printf("\t%s = %s [%s]\n", attribute.Key, attribute.Value.ToString(), attribute.Pos)
		}
	}

	if len(scene.Fields) > 0 {
		fmt.Println("Fields:")

		for _, field := range scene.Fields {
			fmt.Printf("\t%s = %s [%s]\n", field.Key, field.Value.ToString(), field.Pos)
		}
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
			fmt.Println("Attributes:")

			for _, field := range section.Fields {
				fmt.Printf("\t%s = %s [%s]\n", field.Key, field.Value.ToString(), field.Pos)
			}
		}
	}
}
