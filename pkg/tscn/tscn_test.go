package tscn

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/atomicptr/godot-tscn-parser/pkg/godot"
)

func TestParseScene(t *testing.T) {
	content := `[gd_scene]
[ext_resource path="res://Test.tscn" type="PackedScene" id=1]
[node name="Root" type="Node2D"]
[node name="Child" type="Node2D" parent="."]
position = Vector2(13, 37)`
	scene, err := ParseScene(strings.NewReader(content))
	assert.NoError(t, err)
	assert.Len(t, scene.ExtResources, 1)
}

func TestParseSceneWithInvalidFormat(t *testing.T) {
	content := `[gd_scene`
	_, err := ParseScene(strings.NewReader(content))
	assert.Error(t, err)
}

func TestParseProject(t *testing.T) {
	content := `
config_version=4
[application]
config/name="Test Game"
run/main_scene="res://World.tscn"
config/icon="res://icon.png"
[display]
window/size/width=320
window/size/height=180
window/size/test_width=1280
window/size/test_height=720
window/stretch/mode="2d"
window/stretch/aspect="keep"
[input]
[rendering]
environment/default_environment="res://default_env.tres"`
	project, err := ParseProject(strings.NewReader(content))
	assert.NoError(t, err)
	assert.Len(t, project.Application, 3)
	assert.Len(t, project.Display, 6)
	assert.Len(t, project.Input, 0)
	assert.Len(t, project.Rendering, 1)

	configName := project.Application["config/name"].(godot.Value)
	assert.Equal(t, "Test Game", configName.Value)
}

func TestParseProjectWithInvalidFormat(t *testing.T) {
	content := `[test`
	_, err := ParseProject(strings.NewReader(content))
	assert.Error(t, err)
}

func TestParseResource(t *testing.T) {
	content := `[gd_resource type="Environment" load_steps=2 format=2]
[sub_resource type="ProceduralSky" id=1]
[resource]
background_mode = 2
background_sky = SubResource( 1 )`
	resource, err := ParseResource(strings.NewReader(content))
	assert.NoError(t, err)
	assert.Len(t, resource.SubResources, 1)
	assert.Len(t, resource.Fields, 2)
}

func TestParseResourceWithInvalidFormat(t *testing.T) {
	content := `[gd_resource`
	_, err := ParseResource(strings.NewReader(content))
	assert.Error(t, err)
}
