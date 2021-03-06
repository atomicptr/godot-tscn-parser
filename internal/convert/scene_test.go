package convert

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/atomicptr/godot-tscn-parser/internal/parser"
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

	tscnFile, err := parser.Parse(strings.NewReader(content))
	assert.NoError(t, err)

	scene, err := ToGodotScene(tscnFile)
	assert.NoError(t, err)

	assert.Len(t, scene.ExtResources, 3)
	assert.Len(t, scene.SubResources, 1)

	node, err := scene.GetNode("Hazards")
	assert.NoError(t, err)

	assert.Equal(t, "Hazards", node.Name)
}

func TestConvertToGodotSceneWithResource(t *testing.T) {
	content := `[gd_resource]`
	tscnFile, err := parser.Parse(strings.NewReader(content))
	assert.NoError(t, err)

	_, err = ToGodotScene(tscnFile)
	assert.Error(t, err)
}

func TestConvertToGodotSceneWithInvalidExtResource(t *testing.T) {
	content := `[gd_scene]
[ext_resource]`
	tscnFile, err := parser.Parse(strings.NewReader(content))
	assert.NoError(t, err)

	_, err = ToGodotScene(tscnFile)
	assert.Error(t, err)
}

func TestConvertToGodotSceneWithInvalidSubResource(t *testing.T) {
	content := `[gd_scene]
[sub_resource]`
	tscnFile, err := parser.Parse(strings.NewReader(content))
	assert.NoError(t, err)

	_, err = ToGodotScene(tscnFile)
	assert.Error(t, err)
}

func TestConvertToGodotSceneWithInvalidEditable(t *testing.T) {
	content := `[gd_scene]
[editable]`
	tscnFile, err := parser.Parse(strings.NewReader(content))
	assert.NoError(t, err)

	_, err = ToGodotScene(tscnFile)
	assert.Error(t, err)
}

func TestConvertToGodotSceneWithInvalidConnection(t *testing.T) {
	content := `[gd_scene]
[connection]`
	tscnFile, err := parser.Parse(strings.NewReader(content))
	assert.NoError(t, err)

	_, err = ToGodotScene(tscnFile)
	assert.Error(t, err)
}

func TestConvertToGodotSceneWithInvalidResourceType(t *testing.T) {
	content := `[gd_scene]
[this_does_not_exist]`
	tscnFile, err := parser.Parse(strings.NewReader(content))
	assert.NoError(t, err)

	_, err = ToGodotScene(tscnFile)
	assert.Error(t, err)
}

func TestConvertToGodotSceneWithInvalidNodeTree(t *testing.T) {
	content := `[gd_scene]
[node parent="." type="Node2D"]`
	tscnFile, err := parser.Parse(strings.NewReader(content))
	assert.NoError(t, err)

	_, err = ToGodotScene(tscnFile)
	assert.Error(t, err)
}

func testErrorWithConvertSection(
	t *testing.T,
	content string,
	convertFunc func(s *parser.GdResource) (interface{}, error),
	isError bool,
) {
	tscnFile, err := parser.Parse(strings.NewReader(content))
	assert.NoError(t, err)

	section := tscnFile.Sections[0]
	_, err = convertFunc(section)
	if isError {
		assert.Error(t, err)
	} else {
		assert.NoError(t, err)
	}
}

func TestConvertSectionToExtResource(t *testing.T) {
	testConvertFunc := func(s *parser.GdResource) (interface{}, error) {
		return convertSectionToExtResource(s)
	}

	table := []struct {
		convertFunc func(s *parser.GdResource) (interface{}, error)
		isError     bool
		content     string
	}{
		{testConvertFunc, true, `[gd_scene] [sub_resource]`},
		{testConvertFunc, true, `[gd_scene] [ext_resource]`},
		{testConvertFunc, true, `[gd_scene] [ext_resource path="res://Test.tscn"]`},
		{testConvertFunc, true, `[gd_scene] [ext_resource path="res://Test.tscn" type="PackedScene"]`},
		{testConvertFunc, false, `[gd_scene] [ext_resource path="res://Test.tscn" type="PackedScene" id=1]`},
	}

	for _, tc := range table {
		testErrorWithConvertSection(t, tc.content, tc.convertFunc, tc.isError)
	}
}

func TestConvertSectionToSubResource(t *testing.T) {
	testConvertFunc := func(s *parser.GdResource) (interface{}, error) {
		return convertSectionToSubResource(s)
	}

	table := []struct {
		convertFunc func(s *parser.GdResource) (interface{}, error)
		isError     bool
		content     string
	}{
		{testConvertFunc, true, `[gd_scene] [ext_resource]`},
		{testConvertFunc, true, `[gd_scene] [sub_resource]`},
		{testConvertFunc, true, `[gd_scene] [sub_resource type="TileSet"]`},
		{testConvertFunc, false, `[gd_scene] [sub_resource type="TileSet" id=2]`},
	}

	for _, tc := range table {
		testErrorWithConvertSection(t, tc.content, tc.convertFunc, tc.isError)
	}
}

func TestConvertSectionToEditable(t *testing.T) {
	testConvertFunc := func(s *parser.GdResource) (interface{}, error) {
		return convertSectionToEditable(s)
	}

	table := []struct {
		convertFunc func(s *parser.GdResource) (interface{}, error)
		isError     bool
		content     string
	}{
		{testConvertFunc, true, `[gd_scene] [ext_resource]`},
		{testConvertFunc, true, `[gd_scene] [editable]`},
		{testConvertFunc, false, `[gd_scene] [editable path="TestNode"]`},
	}

	for _, tc := range table {
		testErrorWithConvertSection(t, tc.content, tc.convertFunc, tc.isError)
	}
}

func TestConvertSectionToConnection(t *testing.T) {
	testConvertFunc := func(s *parser.GdResource) (interface{}, error) {
		return convertSectionToConnection(s)
	}

	table := []struct {
		convertFunc func(s *parser.GdResource) (interface{}, error)
		isError     bool
		content     string
	}{
		{testConvertFunc, true, `[gd_scene] [ext_resource]`},
		{testConvertFunc, true, `[gd_scene] [connection]`},
		{testConvertFunc, true, `[gd_scene] [connection from="."]`},
		{testConvertFunc, true, `[gd_scene] [connection from="." to="."]`},
		{testConvertFunc, true, `[gd_scene] [connection from="." to="." signal="connect"]`},
		{testConvertFunc, false, `[gd_scene] [connection from="." to="." signal="connect" method="OnSignalConnect"]`},
		{testConvertFunc, false, `[gd_scene] [connection from="." to="." signal="connect" method="OnSignalConnect" flags=7]`},
		{testConvertFunc, false, `[gd_scene]
[connection from="." to="." signal="connect" method="OnSignalConnect" flags=7 binds=["Test", Vector3(1, 2, 3)]]`},
	}

	for _, tc := range table {
		testErrorWithConvertSection(t, tc.content, tc.convertFunc, tc.isError)
	}
}

func TestConvertSectionToUnattachedNode(t *testing.T) {
	testConvertFunc := func(s *parser.GdResource) (interface{}, error) {
		return convertSectionToUnattachedNode(s)
	}

	table := []struct {
		convertFunc func(s *parser.GdResource) (interface{}, error)
		isError     bool
		content     string
	}{
		{testConvertFunc, true, `[gd_scene] [ext_resource]`},
		{testConvertFunc, true, `[gd_scene] [node]`},
		{testConvertFunc, true, `[gd_scene] [node name="Test" instance=ExtResource(1,3)]`},
		{testConvertFunc, false, `[gd_scene] [node name="Test"]`},
		{testConvertFunc, false, `[gd_scene] [node name="Test" instance=ExtResource(1)]`},
	}

	for _, tc := range table {
		testErrorWithConvertSection(t, tc.content, tc.convertFunc, tc.isError)
	}
}

func TestBuildNodeTreeWithInvalidTree(t *testing.T) {
	content := `[gd_scene]
[node name="Test" type="Node2D"]
[node name="Test3" parent="NonExistentParent" type="Node2D"]`
	tscnFile, err := parser.Parse(strings.NewReader(content))
	assert.NoError(t, err)

	tree, err := buildNodeTree(tscnFile)
	assert.NoError(t, err)

	unassignableNodes, err := tree.GetNode(internalNodeUnassignableNodes)
	assert.NoError(t, err)

	assert.NotEmpty(t, unassignableNodes.Children, "tree should contain nodes that we can't assign")
}

func TestBuildNodeTreeWithInvalidParentParameter(t *testing.T) {
	content := `[gd_scene]
[node name="Test" type="Node2D"]
[node name="Test2" parent=1234 type="Node2D"]`
	tscnFile, err := parser.Parse(strings.NewReader(content))
	assert.NoError(t, err)

	_, err = buildNodeTree(tscnFile)
	assert.Error(t, err)
}

func TestBuildNodeTreeWithInvalidChildNode(t *testing.T) {
	content := `[gd_scene]
[node name="Root Node" type="Node2D"]
[node parent="." type="Node2D"]`
	tscnFile, err := parser.Parse(strings.NewReader(content))
	assert.NoError(t, err)

	_, err = buildNodeTree(tscnFile)
	assert.Error(t, err)
}

func TestRegressionConvertToGdSceneWithEditableNodeWithMissingChildren(t *testing.T) {
	content := `[gd_scene]
[ext_resource path="res://TestNode.tscn" type="PackedScene" id=3]
[node name="Root" type="Node2D"]
[node name="EditableNode" parent="." instance=ExtResource(3)]
[node name="ChildNodeWeAreOverwriting" parent="EditableNode/A/B/C/D"]
position = Vector2(13, 37)
[editable path="EditableNode"]`
	tscnFile, err := parser.Parse(strings.NewReader(content))
	assert.NoError(t, err)

	scene, err := ToGodotScene(tscnFile)
	assert.NoError(t, err)

	node, err := scene.GetNode("EditableNode/A/B/C/D/ChildNodeWeAreOverwriting")
	assert.NoError(t, err)

	if err != nil {
		return
	}

	assert.Equal(t, "ChildNodeWeAreOverwriting", node.Name)
	assert.Len(t, node.Fields, 1)
}

func TestRegressionConvertToGdSceneWithEditableNodeWithMissingChildrenAndASlightlyComplexerSetup(t *testing.T) {
	content := `[gd_scene]
[ext_resource path="res://TestNode.tscn" type="PackedScene" id=3]
[node name="Root" type="Node2D"]
[node name="EditableNode" parent="." instance=ExtResource(3)]
[node name="CouldThisWork" parent="EditableNode/A/B"]
position = Vector2(42, 0)
[node name="ChildNodeWeAreOverwriting" parent="EditableNode/A/B/CouldThisWork/D"]
position = Vector2(13, 37)
[editable path="EditableNode"]`
	tscnFile, err := parser.Parse(strings.NewReader(content))
	assert.NoError(t, err)

	scene, err := ToGodotScene(tscnFile)
	assert.NoError(t, err)

	node, err := scene.GetNode("EditableNode/A/B/CouldThisWork/D/ChildNodeWeAreOverwriting")
	assert.NoError(t, err)

	if err != nil {
		return
	}

	assert.Equal(t, "ChildNodeWeAreOverwriting", node.Name)
	assert.Len(t, node.Fields, 1)
}
