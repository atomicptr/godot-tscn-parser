package validate

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/atomicptr/godot-tscn-parser/internal/parser"
)

func TestValidatorFirstNodeHasNoParentWithBadTestCase(t *testing.T) {
	tscn, err := parser.Parse(strings.NewReader(`[gd_scene format=2] [node name="Test" type="Node2D" parent="."]`))
	assert.NoError(t, err)
	assert.Error(t, validatorFirstNodeHasNoParent(tscn))
}

func TestValidatorFirstNodeHasNoParentWithGoodTestCase(t *testing.T) {
	tscn, err := parser.Parse(strings.NewReader(`[gd_scene format=2] [node name="Test" type="Node2D"]`))
	assert.NoError(t, err)
	assert.NoError(t, validatorFirstNodeHasNoParent(tscn))
}

func TestValidatorOnlyOneRootNode(t *testing.T) {
	tscn, err := parser.Parse(strings.NewReader(`[gd_scene]
[node name="Test" type="Node2D"]
[node name="Test2" type="Node2D"]`))
	assert.NoError(t, err)
	assert.Error(t, validatorOnlyOneRootNode(tscn))
}
