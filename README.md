# Godot TSCN Parser

[![Travis CI](https://api.travis-ci.com/atomicptr/godot-tscn-parser.svg?branch=master)](https://travis-ci.com/atomicptr/godot-tscn-parser)
[![Go Report Card](https://goreportcard.com/badge/github.com/atomicptr/godot-tscn-parser)](https://goreportcard.com/report/github.com/atomicptr/godot-tscn-parser)
[![Coverage Status](https://coveralls.io/repos/github/atomicptr/godot-tscn-parser/badge.svg?branch=master)](https://coveralls.io/github/atomicptr/godot-tscn-parser?branch=master)

Go library for parsing
the [Godot TSCN file format](https://docs.godotengine.org/en/stable/development/file_formats/tscn.html).

Powered by the great [participle](https://github.com/alecthomas/participle) parser library.

## Usage

```go
package main

import (
	"fmt"
	"os"

	"github.com/atomicptr/godot-tscn-parser/pkg/tscn"
)

func main() {
	// open the file
	f, err := os.Open("./path/to/my/scene.tscn")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// parse the scene, this accepts an io.Reader
	scene, err := tscn.ParseScene(f)
	if err != nil {
		panic(err)
	}

	// get the node "Sprite" which is a child of "Player" which is a child of
	// of the root node
	playerSpriteNode, err := scene.GetNode("Player/Sprite")
	if err != nil {
		panic(err)
	}

	// access a field, keep in mind that TSCN files only store non default values
	position := playerSpriteNode.Fields["position"]
	fmt.Printf("Player/Sprite is at position %v\n", position)
}
```

## FAQ

### My TSCN file isn't working, can you fix it?

Please open an issue with your TSCN file.

Or even better, open a pull request which adds the file to test/fixtures.

### Why can't I read a certain value?

Godot does not store default values in TSCN files. Which means you'll only see stuff you've added or changed yourself in
these files.

### What are "Volatile Nodes"?

If you have an instanced scene in your scene, and you mark it as editable (and actually edit something), you might find
nodes with the type "VolatileNode" in your tree. Since Godot doesn't store unchanged things, these nodes are a band aid
for this library to render a proper tree structure. You can basically ignore them.

## License

MIT