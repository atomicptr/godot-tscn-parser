package convert

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/atomicptr/godot-tscn-parser/internal/parser"
	"github.com/atomicptr/godot-tscn-parser/pkg/godot"
)

func TestToGodotImport(t *testing.T) {
	content := `[remap]
importer="texture"
type="StreamTexture"
[deps]
source_file="res://icon.png"
[customsection]
customfield=1337`
	tscnFile, err := parser.Parse(strings.NewReader(content))
	assert.NoError(t, err)

	imp, err := ToGodotImport(tscnFile)
	assert.NoError(t, err)

	assert.Len(t, imp.Remap, 2)
	remapType := imp.Remap["type"].(godot.Value)
	assert.Equal(t, "StreamTexture", remapType.Value)

	assert.Len(t, imp.Deps, 1)
	depsSourceFile := imp.Deps["source_file"].(godot.Value)
	assert.Equal(t, "res://icon.png", depsSourceFile.Value)

	customField, ok := imp.Rest["customsection"]["customfield"]
	assert.True(t, ok)
	customFieldValue := customField.(godot.Value)
	assert.Equal(t, int64(1337), customFieldValue.Value)
}
