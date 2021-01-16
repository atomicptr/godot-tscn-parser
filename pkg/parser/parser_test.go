package parser

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestParseFileDescriptorWithAttributes(t *testing.T) {
	content := "[gd_scene load_steps=0 format=2]"

	scene, err := Parse(strings.NewReader(content))
	assert.Nil(t, err)

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
	assert.Nil(t, err)

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
	assert.Nil(t, err)

	assert.Equal(t, 1, len(scene.Fields))
	assert.Equal(t, "Yes this works!", *scene.Fields[0].Value.String)
}

func TestFieldDescriptorWithoutArguments(t *testing.T) {
	content := `[gd_scene]`
	scene, err := Parse(strings.NewReader(content))
	assert.Nil(t, err)
	assert.Equal(t, 0, len(scene.Attributes))
}

func TestSectionAttributes(t *testing.T) {
	content := `[gd_scene load_steps=2 format=2]
[ext_resource path="res://CombatSystem/Background/steppes.png" type="Texture" id=1]
[ext_resource path="res://CombatSystem/UserInterface/UILayer.gd" type="Script" id=2]
[node name="CombatDemo" type="Node2D"]
script = ExtResource( 1 )`

	scene, err := Parse(strings.NewReader(content))
	assert.Nil(t, err)

	assert.Equal(t, 3, len(scene.Sections))

	table := []struct {
		Key    string
		Values map[string]interface{}
	}{
		{
			Key: "ext_resource",
			Values: map[string]interface{}{
				"path": "res://CombatSystem/Background/steppes.png",
				"type": "Texture",
				"id":   int64(1),
			},
		},
		{
			Key: "ext_resource",
			Values: map[string]interface{}{
				"path": "res://CombatSystem/UserInterface/UILayer.gd",
				"type": "Script",
				"id":   int64(2),
			},
		},
		{
			Key: "node",
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
