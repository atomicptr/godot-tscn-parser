package convert

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"github.com/atomicptr/godot-tscn-parser/internal/parser"
	"github.com/atomicptr/godot-tscn-parser/pkg/godot"
)

func TestInsertFieldEntriesFromSection(t *testing.T) {
	content := `[gd_scene] [custom] a=1 b=2 c=3`
	tscn, err := parser.Parse(strings.NewReader(content))
	assert.NoError(t, err)

	m := make(map[string]interface{})
	insertFieldEntriesFromSection(tscn.Sections[0], m)

	assert.Len(t, m, 3)
	b := m["b"].(godot.Value)
	assert.Equal(t, int64(2), b.Value)
}

func TestConvertGdValueForArray(t *testing.T) {
	content := `value = [1, 2, TestType(1337)]`
	tscn, _ := parser.Parse(strings.NewReader(content))
	v, ok := convertGdValue(tscn.Fields[0].Value).(godot.Value)
	assert.True(t, ok)

	value, ok := v.Value.([]interface{})
	assert.True(t, ok)
	assert.Len(t, value, 3)
}

func TestConvertGdValueForMap(t *testing.T) {
	content := `value = {"a": 1, "b": [1, 2], "c": {"d": 5}}`
	tscn, _ := parser.Parse(strings.NewReader(content))
	v, ok := convertGdValue(tscn.Fields[0].Value).(godot.Value)
	assert.True(t, ok)

	value, ok := v.Value.(map[string]interface{})
	assert.True(t, ok)
	assert.Len(t, value, 3)
}

func TestConvertGdValueForType(t *testing.T) {
	content := `value = CustomType(1337, OtherType(12, 3))`
	tscn, _ := parser.Parse(strings.NewReader(content))
	v, ok := convertGdValue(tscn.Fields[0].Value).(godot.Type)
	assert.True(t, ok)

	assert.Equal(t, "CustomType", v.Identifier)
	assert.Len(t, v.Parameters, 2)
}

func TestConvertGdValueForMapField(t *testing.T) {
	content := `value = CustomType("a": 13, "b": 37)`
	tscn, _ := parser.Parse(strings.NewReader(content))
	v, ok := convertGdValue(tscn.Fields[0].Value).(godot.Type)
	assert.True(t, ok)

	aKey := v.Parameters[0].(godot.KeyValuePair).Key
	aValue := v.Parameters[0].(godot.KeyValuePair).Value.(godot.Value).Value.(int64)
	bKey := v.Parameters[1].(godot.KeyValuePair).Key
	bValue := v.Parameters[1].(godot.KeyValuePair).Value.(godot.Value).Value.(int64)

	assert.Equal(t, "a", aKey)
	assert.Equal(t, int64(13), aValue)
	assert.Equal(t, "b", bKey)
	assert.Equal(t, int64(37), bValue)
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

		if strings.HasPrefix(file, ".import") {
			_, err = ToGodotImport(tscnFile)
			assert.NoError(t, errors.Wrapf(err, "error with fixture: '%s", file))
		}

		err = f.Close()
		if err != nil {
			panic(err)
		}
	}
}
