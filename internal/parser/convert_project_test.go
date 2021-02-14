package parser

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/atomicptr/godot-tscn-parser/pkg/godot"
)

func TestConvertToGodotProject(t *testing.T) {
	content := `config_version=4
[application]
config/name="Your first Godot Game"
[customsection]
customfield=1337`
	tscnFile, err := Parse(strings.NewReader(content))
	assert.NoError(t, err)

	project, err := tscnFile.ConvertToGodotProject()
	assert.NoError(t, err)

	configVersion, ok := project.Fields["config_version"]
	assert.True(t, ok)
	configVersionValue := configVersion.(godot.Value)
	assert.Equal(t, int64(4), configVersionValue.Value)

	name, ok := project.Application["config/name"]
	assert.True(t, ok)
	nameValue := name.(godot.Value)
	assert.Equal(t, "Your first Godot Game", nameValue.Value)

	customField, ok := project.Rest["customsection"]["customfield"]
	assert.True(t, ok)
	customFieldValue := customField.(godot.Value)
	assert.Equal(t, int64(1337), customFieldValue.Value)
}
