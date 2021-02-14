package parser

import (
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
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

func testErrorWithConvertSection(t *testing.T, content string, convertFunc func(s *GdResource) (interface{}, error), error bool) {
	tscnFile, err := Parse(strings.NewReader(content))
	assert.NoError(t, err)

	section := tscnFile.Sections[0]
	_, err = convertFunc(section)
	if error {
		assert.Error(t, err)
	} else {
		assert.NoError(t, err)
	}
}

func TestConvertSectionToExtResource(t *testing.T) {
	testConvertFunc := func(s *GdResource) (interface{}, error) {
		return convertSectionToExtResource(s)
	}
	assumeError := func(t *testing.T, content string) {
		testErrorWithConvertSection(t, content, testConvertFunc, true)
	}
	assumeNoError := func(t *testing.T, content string) {
		testErrorWithConvertSection(t, content, testConvertFunc, false)
	}

	assumeError(t, `[gd_scene] [sub_resource]`)
	assumeError(t, `[gd_scene] [ext_resource]`)
	assumeError(t, `[gd_scene] [ext_resource path="res://Test.tscn"]`)
	assumeError(t, `[gd_scene] [ext_resource path="res://Test.tscn" type="PackedScene"]`)
	assumeNoError(t, `[gd_scene] [ext_resource path="res://Test.tscn" type="PackedScene" id=1]`)
}

func TestConvertSectionToSubResource(t *testing.T) {
	testConvertFunc := func(s *GdResource) (interface{}, error) {
		return convertSectionToSubResource(s)
	}
	assumeError := func(t *testing.T, content string) {
		testErrorWithConvertSection(t, content, testConvertFunc, true)
	}
	assumeNoError := func(t *testing.T, content string) {
		testErrorWithConvertSection(t, content, testConvertFunc, false)
	}

	assumeError(t, `[gd_scene] [ext_resource]`)
	assumeError(t, `[gd_scene] [sub_resource]`)
	assumeError(t, `[gd_scene] [sub_resource type="TileSet"]`)
	assumeNoError(t, `[gd_scene] [sub_resource type="TileSet" id=2]`)
}

func TestConvertSectionToEditable(t *testing.T) {
	testConvertFunc := func(s *GdResource) (interface{}, error) {
		return convertSectionToEditable(s)
	}
	assumeError := func(t *testing.T, content string) {
		testErrorWithConvertSection(t, content, testConvertFunc, true)
	}
	assumeNoError := func(t *testing.T, content string) {
		testErrorWithConvertSection(t, content, testConvertFunc, false)
	}

	assumeError(t, `[gd_scene] [ext_resource]`)
	assumeError(t, `[gd_scene] [editable]`)
	assumeNoError(t, `[gd_scene] [editable path="TestNode"]`)
}

func TestConvertSectionToConnection(t *testing.T) {
	testConvertFunc := func(s *GdResource) (interface{}, error) {
		return convertSectionToConnection(s)
	}
	assumeError := func(t *testing.T, content string) {
		testErrorWithConvertSection(t, content, testConvertFunc, true)
	}
	assumeNoError := func(t *testing.T, content string) {
		testErrorWithConvertSection(t, content, testConvertFunc, false)
	}

	assumeError(t, `[gd_scene] [ext_resource]`)
	assumeError(t, `[gd_scene] [connection]`)
	assumeError(t, `[gd_scene] [connection from="."]`)
	assumeError(t, `[gd_scene] [connection from="." to="."]`)
	assumeError(t, `[gd_scene] [connection from="." to="." signal="connect"]`)
	assumeNoError(t, `[gd_scene] [connection from="." to="." signal="connect" method="OnSignalConnect"]`)
	assumeNoError(t, `[gd_scene] [connection from="." to="." signal="connect" method="OnSignalConnect" flags=7]`)
	assumeNoError(t, `[gd_scene] [connection from="." to="." signal="connect" method="OnSignalConnect" flags=7 binds=["Test", Vector3(1, 2, 3)]]`)
}

func TestConvertSectionToUnattachedNode(t *testing.T) {
	testConvertFunc := func(s *GdResource) (interface{}, error) {
		return convertSectionToUnattachedNode(s)
	}
	assumeError := func(t *testing.T, content string) {
		testErrorWithConvertSection(t, content, testConvertFunc, true)
	}
	assumeNoError := func(t *testing.T, content string) {
		testErrorWithConvertSection(t, content, testConvertFunc, false)
	}

	assumeError(t, `[gd_scene] [ext_resource]`)
	assumeError(t, `[gd_scene] [node]`)
	assumeError(t, `[gd_scene] [node name="Test" instance=ExtResource(1,3)]`)
	assumeNoError(t, `[gd_scene] [node name="Test"]`)
	assumeNoError(t, `[gd_scene] [node name="Test" instance=ExtResource(1)]`)
}

func TestBuildNodeTreeWithInvalidTree(t *testing.T) {
	content := `[gd_scene]
[node name="Test" type="Node2D"]
[node name="Test3" parent="NonExistentParent" type="Node2D"]`
	tscnFile, err := Parse(strings.NewReader(content))
	assert.NoError(t, err)

	_, err = buildNodeTree(tscnFile)
	assert.Error(t, err)
}

func TestBuildNodeTreeWithInvalidParentParameter(t *testing.T) {
	content := `[gd_scene]
[node name="Test" type="Node2D"]
[node name="Test2" parent=1234 type="Node2D"]`
	tscnFile, err := Parse(strings.NewReader(content))
	assert.NoError(t, err)

	_, err = buildNodeTree(tscnFile)
	assert.Error(t, err)
}

func TestBuildNodeTreeWithInvalidChildNode(t *testing.T) {
	content := `[gd_scene]
[node name="Root Node" type="Node2D"]
[node parent="." type="Node2D"]`
	tscnFile, err := Parse(strings.NewReader(content))
	assert.NoError(t, err)

	_, err = buildNodeTree(tscnFile)
	assert.Error(t, err)
}

// keep integration tests at the bottom please
func TestIntegrationConvertToGodotSceneFixtures(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	files, err := filepath.Glob(filepath.Join(cwd, "../../test/fixtures", "*"))
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

		if tscnFile.Key == "gd_scene" {
			_, err = tscnFile.ConvertToGodotScene()
			assert.NoError(t, errors.Wrapf(err, "error with fixture: '%s'", file))
		}

		err = f.Close()
		if err != nil {
			panic(err)
		}
	}
}
