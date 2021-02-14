package parser

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

// keep integration tests at the bottom please
func TestIntegrationConvertFixtures(t *testing.T) {
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

		tscnFile, err := Parse(f)
		if err != nil {
			continue
		}

		if tscnFile.Key == TscnTypeGodotScene {
			_, err = tscnFile.ConvertToGodotScene()
			assert.NoError(t, errors.Wrapf(err, "error with fixture: '%s'", file))
		}

		if strings.HasSuffix(file, ".godot") {
			_, err = tscnFile.ConvertToGodotProject()
			assert.NoError(t, errors.Wrapf(err, "error with fixture: '%s", file))
		}

		err = f.Close()
		if err != nil {
			panic(err)
		}
	}
}
