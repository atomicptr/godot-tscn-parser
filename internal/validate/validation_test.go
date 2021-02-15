package validate

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/atomicptr/godot-tscn-parser/internal/parser"
)

func TestIntegrationTscnFileFormat(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	files, err := filepath.Glob(filepath.Join(cwd, "..", "..", "test", "fixtures", "*"))
	if err != nil {
		panic(err)
	}

	assert.NotEmpty(t, files)

	for _, file := range files {
		// ignore the README.md file
		if filepath.Base(file) == "README.md" {
			continue
		}

		f, err := os.Open(filepath.Clean(file))
		if err != nil {
			panic(err)
		}

		tscnFile, err := parser.Parse(f)
		if err != nil {
			panic(err)
		}

		err = TscnFileFormat(tscnFile)
		assert.NoError(t, err)

		err = f.Close()
		if err != nil {
			panic(err)
		}
	}
}
