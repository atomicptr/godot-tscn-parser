package validate

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/atomicptr/godot-tscn-parser/internal/parser"
)

func TestValidatorExtResourceRequiredAttributes(t *testing.T) {
	tscn, err := parser.Parse(strings.NewReader(`[gd_scene]
[ext_resource path="res://Player.tscn" type="PackedScene" id=5]`))
	assert.NoError(t, err)
	assert.NoError(t, validatorExtResourceRequiredAttributes(tscn))
}

func TestValidatorExtResourceRequiredAttributesNoPath(t *testing.T) {
	tscn, err := parser.Parse(strings.NewReader(`[gd_scene]
[ext_resource type="PackedScene" id=5]`))
	assert.NoError(t, err)
	assert.Error(t, validatorExtResourceRequiredAttributes(tscn))
}

func TestValidatorExtResourceRequiredAttributesNoType(t *testing.T) {
	tscn, err := parser.Parse(strings.NewReader(`[gd_scene]
[ext_resource path="res://Player.tscn" id=5]`))
	assert.NoError(t, err)
	assert.Error(t, validatorExtResourceRequiredAttributes(tscn))
}

func TestValidatorExtResourceRequiredAttributesNoId(t *testing.T) {
	tscn, err := parser.Parse(strings.NewReader(`[gd_scene]
[ext_resource path="res://Player.tscn" type="PackedScene"]`))
	assert.NoError(t, err)
	assert.Error(t, validatorExtResourceRequiredAttributes(tscn))
}
