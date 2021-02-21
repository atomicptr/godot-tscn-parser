package validate

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/atomicptr/godot-tscn-parser/internal/parser"
)

func TestValidatorSceneIsInSupportedFormatNoFormat(t *testing.T) {
	tscn, err := parser.Parse(strings.NewReader(`[gd_scene]`))
	assert.NoError(t, err)
	err = validatorSceneIsInSupportedFormat(tscn)
	assert.NoError(t, err)
}

func TestValidatorSceneIsInSupportedFormatValidFormat(t *testing.T) {
	tscn, err := parser.Parse(strings.NewReader(`[gd_scene format=2]`))
	assert.NoError(t, err)
	err = validatorSceneIsInSupportedFormat(tscn)
	assert.NoError(t, err)
}

func TestValidatorSceneIsInSupportedFormatInvalidFormat1(t *testing.T) {
	tscn, err := parser.Parse(strings.NewReader(`[gd_scene format=1]`))
	assert.NoError(t, err)
	err = validatorSceneIsInSupportedFormat(tscn)
	assert.Error(t, err)
}

func TestValidatorSceneIsInSupportedFormatInvalidFormat3(t *testing.T) {
	tscn, err := parser.Parse(strings.NewReader(`[gd_scene format=3]`))
	assert.NoError(t, err)
	err = validatorSceneIsInSupportedFormat(tscn)
	assert.Error(t, err)
}
