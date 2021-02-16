package convert

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/atomicptr/godot-tscn-parser/internal/parser"
	"github.com/atomicptr/godot-tscn-parser/pkg/godot"
)

func TestToGodotResource(t *testing.T) {
	content := `[gd_resource type="Environment" load_steps=2 format=2]

[sub_resource type="ProceduralSky" id=1]

[resource]
background_mode = 2
background_sky = SubResource( 1 )`
	tscn, err := parser.Parse(strings.NewReader(content))
	assert.NoError(t, err)

	res, err := ToGodotResource(tscn)
	assert.NoError(t, err)

	assert.Equal(t, "Environment", res.Type)
	assert.Len(t, res.SubResources, 1)
	assert.Equal(t, "ProceduralSky", res.SubResources[1].Type)
	assert.Len(t, res.Fields, 2)

	backgroundMode := res.Fields["background_mode"].(godot.Value)
	assert.Equal(t, int64(2), backgroundMode.Value)
}
