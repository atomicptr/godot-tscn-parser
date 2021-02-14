package tscn

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestParseSceneWithInvalidScene(t *testing.T) {
	content := `[gd_scene`
	_, err := ParseScene(strings.NewReader(content))
	assert.Error(t, err)
}
