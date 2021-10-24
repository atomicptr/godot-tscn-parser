package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseFail(t *testing.T) {
	_, err := Parse(strings.NewReader("This is not a proper TSCN file"))
	assert.Error(t, err)
}

func TestParseFileDescriptorWithAttributes(t *testing.T) {
	content := "[gd_scene load_steps=0 format=2]"

	scene, err := Parse(strings.NewReader(content))
	assert.NoError(t, err)

	assert.Equal(t, "gd_scene", scene.Key)
	assert.Equal(t, 2, len(scene.Attributes))
	assert.Equal(t, 0, len(scene.Sections))
	assert.NotNil(t, scene.Pos)

	for _, attribute := range scene.Attributes {
		assert.NotNil(t, attribute.Pos)

		if !assertField(t, attribute, "load_steps", int64(0)) {
			continue
		}
		if !assertField(t, attribute, "format", int64(2)) {
			continue
		}

		assert.Fail(t, fmt.Sprintf("Unknown attribute found '%s", attribute.Key))
	}
}

func TestParseFileDescriptorFields(t *testing.T) {
	content := `; a comment just in case
[gd_scene load_steps=0 format=2]
; another comment right here
int_field = 10
negative_int_field = -10
string_field = "Test"
reference_field = ExtResource( 1337 )
reference_field_multi_args = Vector2( 12.37, 13.37 )
float_field = 13.37
negative_float_field = -69.0 ; nice
bool_field = true
array_field = [ 13.37, 42.0, 12.12 ]
map_field = {
    "string_value": "value",
    "float_value": 13.37
}`
	scene, err := Parse(strings.NewReader(content))
	assert.NoError(t, err)

	for _, field := range scene.Fields {
		if !assertField(t, field, "int_field", int64(10)) {
			continue
		}
		if !assertField(t, field, "negative_int_field", int64(-10)) {
			continue
		}
		if !assertField(t, field, "string_field", "Test") {
			continue
		}
		if !assertField(t, field, "reference_field", "ExtResource", int64(1337)) {
			continue
		}
		if !assertField(t, field, "reference_field_multi_args", "Vector2", 12.37, 13.37) {
			continue
		}
		if !assertField(t, field, "float_field", 13.37) {
			continue
		}
		if !assertField(t, field, "negative_float_field", -69.0) {
			continue
		}
		if !assertField(t, field, "bool_field", true) {
			continue
		}
		if !assertField(t, field, "array_field", 13.37, 42.0, 12.12) {
			continue
		}
		if !assertField(t, field, "map_field",
			keyValuePair{Key: "string_value", Value: "value"},
			keyValuePair{Key: "float_value", Value: 13.37},
		) {
			continue
		}

		assert.Fail(t, fmt.Sprintf("Unknown field found '%s", field.Key))
	}
}

func TestParseFieldWithoutFileDescriptor(t *testing.T) {
	content := "string_value = \"Yes this works!\""
	scene, err := Parse(strings.NewReader(content))
	assert.NoError(t, err)

	assert.Equal(t, 1, len(scene.Fields))
	assert.Equal(t, "Yes this works!", *scene.Fields[0].Value.String)
}

func TestParseNodeWithGroups(t *testing.T) {
	content := `[gd_scene]
[node name="Root" type="Spatial"]

[node name="Spatial" type="Spatial" parent="." groups=[
"test1",
"test2",
"test3",
]]`
	scene, err := Parse(strings.NewReader(content))
	assert.NoError(t, err)

	spatial := scene.Sections[1]

	groupsFound := false
	for _, attrib := range spatial.Attributes {
		if attrib.Key != "groups" {
			continue
		}

		groupsFound = true

		assert.NotNil(t, attrib.Value.Array)
		groups := attrib.Value.Array
		assert.Len(t, groups, 3)
		assert.Equal(t, "test1", *groups[0].String)
		assert.Equal(t, "test2", *groups[1].String)
		assert.Equal(t, "test3", *groups[2].String)
	}

	if !groupsFound {
		assert.Fail(t, "Could not find attribute groups")
	}
}

func TestParseWithEmptyArray(t *testing.T) {
	scene, err := Parse(strings.NewReader("array = []"))
	assert.NoError(t, err)
	field := scene.Fields[0]

	fieldAsArray, ok := field.Value.Raw().([]*GdValue)
	assert.True(t, ok)
	assert.Empty(t, fieldAsArray)
	assert.Equal(t, "[]", field.Value.ToString())
}

func TestParseWithEmptyMap(t *testing.T) {
	scene, err := Parse(strings.NewReader("map = {}"))
	assert.NoError(t, err)
	field := scene.Fields[0]

	fieldAsArray, ok := field.Value.Raw().([]*GdMapField)
	assert.True(t, ok)
	assert.Empty(t, fieldAsArray)
	assert.Equal(t, "{}", field.Value.ToString())
}

func TestFieldDescriptorWithoutArguments(t *testing.T) {
	content := `[gd_scene]`
	scene, err := Parse(strings.NewReader(content))
	assert.NoError(t, err)
	assert.Equal(t, 0, len(scene.Attributes))
}

func TestSectionAttributes(t *testing.T) {
	content := `[gd_scene load_steps=2 format=2]
[ext_resource path="res://CombatSystem/Background/steppes.png" type="Texture" id=1]
[ext_resource path="res://CombatSystem/UserInterface/UILayer.gd" type="Script" id=2]
[node name="CombatDemo" type="Node2D"]
script = ExtResource( 1 )`

	scene, err := Parse(strings.NewReader(content))
	assert.NoError(t, err)

	assert.Equal(t, 3, len(scene.Sections))

	table := []struct {
		Key    string
		Values map[string]interface{}
	}{
		{
			Key: ResourceTypeExtResource,
			Values: map[string]interface{}{
				"path": "res://CombatSystem/Background/steppes.png",
				"type": "Texture",
				"id":   int64(1),
			},
		},
		{
			Key: ResourceTypeExtResource,
			Values: map[string]interface{}{
				"path": "res://CombatSystem/UserInterface/UILayer.gd",
				"type": "Script",
				"id":   int64(2),
			},
		},
		{
			Key: ResourceTypeNode,
			Values: map[string]interface{}{
				"name": "CombatDemo",
				"type": "Node2D",
			},
		},
	}

	for index, expected := range table {
		actual := scene.Sections[index]

		assert.Equal(t, expected.Key, actual.ResourceType)

		for _, attribute := range actual.Attributes {
			assert.Equal(t, expected.Values[attribute.Key], attribute.Value.Raw())
		}
	}
}

func TestReferenceTypeWithKeyValuePairs(t *testing.T) {
	content := `[gd_scene format=2]
reference_type = Object(InputEventKey,"resource_local_to_scene":false,"resource_name":"","device":0,"alt":false)`
	_, err := Parse(strings.NewReader(content))
	assert.NoError(t, err)
}

func TestScientificENotation(t *testing.T) {
	content := `[gd_scene]
scientific_e_float = 1e-05`

	scene, err := Parse(strings.NewReader(content))
	require.NoError(t, err)
	assertField(t, scene.Fields[0], "scientific_e_float", 1e-05)
}

// keep regression tests at the bottom please (above integration tests though)
func TestRegressionFieldNamesStartingWithNumbers(t *testing.T) {
	content := `[gd_scene format=2]
[sub_resource type="TileSet" id=25]
0/name = "TileSet1.svg 0"
0/texture = ExtResource( 2 )
0/tex_offset = Vector2( 0, 0 )
0/modulate = Color( 1, 1, 1, 1 )
0/region = Rect2( 0, 0, 896, 512 )`
	_, err := Parse(strings.NewReader(content))
	assert.NoError(t, err)
}

func TestRegressionEmptyComment(t *testing.T) {
	content := `[gd_scene format=2]
; This is a comment
;
; Notice the comment above? Yeah no comment at all
field_name = "value"`
	_, err := Parse(strings.NewReader(content))
	assert.NoError(t, err)
}

// keep integration tests at the bottom please
func TestIntegrationParseFixtures(t *testing.T) {
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

		_, err = Parse(f)
		assert.NoError(t, errors.Wrapf(err, "error with fixture: '%s'", file))

		err = f.Close()
		if err != nil {
			panic(err)
		}
	}
}
