package convert

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"github.com/atomicptr/godot-tscn-parser/internal/parser"
)

func TestInsertFieldEntriesFromSection(t *testing.T) {
	t.Fail()
}

func TestConvertGdValue(t *testing.T) {
	t.Fail()
}

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

		tscnFile, err := parser.Parse(f)
		if err != nil {
			continue
		}

		if tscnFile.Key == TscnTypeGodotScene {
			_, err = ToGodotScene(tscnFile)
			assert.NoError(t, errors.Wrapf(err, "error with fixture: '%s'", file))
		}

		if strings.HasSuffix(file, ".godot") {
			_, err = ToGodotProject(tscnFile)
			assert.NoError(t, errors.Wrapf(err, "error with fixture: '%s", file))
		}

		if tscnFile.Key == TscnTypeGodotResource {
			_, err = ToGodotResource(tscnFile)
			assert.NoError(t, errors.Wrapf(err, "error with fixture: '%s'", file))
		}

		err = f.Close()
		if err != nil {
			panic(err)
		}
	}
}
