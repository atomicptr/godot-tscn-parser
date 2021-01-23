package main

import (
	"fmt"
	"github.com/atomicptr/godot-tscn-parser/pkg/godot"
	"github.com/atomicptr/godot-tscn-parser/pkg/tscn"
	"os"
)

const filename = "./TestFile.tscn"

func main() {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer func() {
		err := f.Close()
		if err != nil {
			panic(err)
		}
	}()

	scene, err := tscn.ParseScene(f)
	if err != nil {
		panic(err)
	}

	fmt.Printf("# Godot Scene %s [%s]\n", filename, scene.LexerPosition)

	if len(scene.ExtResources) > 0 {
		fmt.Println("## ExtResources:")
		for _, res := range scene.ExtResources {
			fmt.Printf("\t%s (id=%d, path='%s') [%s]\n", res.Type, res.Id, res.Path, res.LexerPosition)
		}
	}

	if len(scene.SubResources) > 0 {
		fmt.Println("## SubResources:")
		for _, res := range scene.SubResources {
			fmt.Printf("\t%s (id=%d) [%s]\n", res.Type, res.Id, res.LexerPosition)
		}
	}

	if scene.Node != nil {
		fmt.Println("## Node Tree:")
		printNodes(scene.Node)
	}
}

func printNodes(node *godot.Node) {
	printNodesWithIndent(node, 0)
}

func printNodesWithIndent(node *godot.Node, indent int) {
	nodeType := node.Type
	if nodeType == "" {
		nodeType = node.Instance
	}
	if nodeType == "" {
		nodeType = "<None>"
	}

	printIndent(indent)
	fmt.Printf("%s (%s) [%s]:\n", node.Name, nodeType, node.MetaData.LexerPosition)

	for field, value := range node.Fields {
		printIndent(indent + 1)
		fmt.Printf("%s = %v\n", field, value)
	}

	for _, childNode := range node.Children {
		printIndent(indent + 1)
		fmt.Println("Children:")
		printNodesWithIndent(childNode, indent+2)
	}
}

func printIndent(indent int) {
	for i := 0; i < indent; i++ {
		fmt.Print("\t")
	}
}
