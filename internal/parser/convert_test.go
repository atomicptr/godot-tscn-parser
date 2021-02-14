package parser

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestConvertToGodotScene(t *testing.T) {
	content := `[gd_scene load_steps=4]
[ext_resource path="res://Player/Player.tscn" type="PackedScene" id=1]
[ext_resource path="res://World/tile_set.svg" type="Texture" id=2]
[ext_resource path="res://World/Hazard.tscn" type="PackedScene" id=3]

[sub_resource type="ConvexPolygonShape2D" id=1]
points = PoolVector2Array( 16, 64, 128, 64, 128, 128, 16, 128 )

[node name="RootNode" type="Node2D"]

[node name="Hazards" type="Area2D" parent="."]

[node name="TrapFloorSpikes" parent="Hazards" instance=ExtResource( 3 )]
position = Vector2( 687.645, -209.178 )

[node name="TrapFloorSpikes2" parent="Hazards" instance=ExtResource( 3 )]
position = Vector2( 811.747, -209.178 )

[editable path="Hazards"]
[connection signal="area_entered" from="Hazards" to="." method="_on_Hazards_area_entered"]`

	tscnFile, err := Parse(strings.NewReader(content))
	assert.NoError(t, err)

	scene, err := tscnFile.ConvertToGodotScene()
	assert.NoError(t, err)

	assert.Len(t, scene.ExtResources, 3)
	assert.Len(t, scene.SubResources, 1)

	node, err := scene.GetNode("Hazards")
	assert.NoError(t, err)

	assert.Equal(t, "Hazards", node.Name)
}

func TestConvertToGodotSceneWithResource(t *testing.T) {
	content := `[gd_resource]`
	tscnFile, err := Parse(strings.NewReader(content))
	assert.NoError(t, err)

	_, err = tscnFile.ConvertToGodotScene()
	assert.Error(t, err)
}

func TestConvertToGodotSceneWithInvalidExtResource(t *testing.T) {
	content := `[gd_scene]
[ext_resource]`
	tscnFile, err := Parse(strings.NewReader(content))
	assert.NoError(t, err)

	_, err = tscnFile.ConvertToGodotScene()
	assert.Error(t, err)
}

func TestConvertToGodotSceneWithInvalidSubResource(t *testing.T) {
	content := `[gd_scene]
[sub_resource]`
	tscnFile, err := Parse(strings.NewReader(content))
	assert.NoError(t, err)

	_, err = tscnFile.ConvertToGodotScene()
	assert.Error(t, err)
}

func TestConvertToGodotSceneWithInvalidEditable(t *testing.T) {
	content := `[gd_scene]
[editable]`
	tscnFile, err := Parse(strings.NewReader(content))
	assert.NoError(t, err)

	_, err = tscnFile.ConvertToGodotScene()
	assert.Error(t, err)
}

func TestConvertToGodotSceneWithInvalidConnection(t *testing.T) {
	content := `[gd_scene]
[connection]`
	tscnFile, err := Parse(strings.NewReader(content))
	assert.NoError(t, err)

	_, err = tscnFile.ConvertToGodotScene()
	assert.Error(t, err)
}

func TestConvertToGodotSceneWithInvalidResourceType(t *testing.T) {
	content := `[gd_scene]
[this_does_not_exist]`
	tscnFile, err := Parse(strings.NewReader(content))
	assert.NoError(t, err)

	_, err = tscnFile.ConvertToGodotScene()
	assert.Error(t, err)
}

func TestConvertToGodotSceneWithInvalidNodeTree(t *testing.T) {
	content := `[gd_scene]
[node parent="." type="Node2D"]`
	tscnFile, err := Parse(strings.NewReader(content))
	assert.NoError(t, err)

	_, err = tscnFile.ConvertToGodotScene()
	assert.Error(t, err)
}
