package convert

import (
	"github.com/atomicptr/godot-tscn-parser/internal/parser"
	"github.com/atomicptr/godot-tscn-parser/pkg/godot"
)

// ToGodotImport tries to convert a TscnFile structure to a godot.Import
func ToGodotImport(tscn *parser.TscnFile) (*godot.Import, error) {
	imp := &godot.Import{
		Remap:  make(map[string]interface{}),
		Deps:   make(map[string]interface{}),
		Params: make(map[string]interface{}),
		Rest:   make(map[string]map[string]interface{}),
		MetaData: godot.MetaData{
			LexerPosition: tscn.Pos,
		},
	}

	insertFieldEntriesFromSection(
		&parser.GdResource{
			ResourceType: tscn.Key,
			Attributes:   tscn.Attributes,
			Fields:       tscn.Fields,
			Pos:          tscn.Pos,
		},
		imp.Remap,
	)

	sectionMap := map[string]map[string]interface{}{
		"deps":   imp.Deps,
		"params": imp.Params,
	}

	for _, section := range tscn.Sections {
		m, ok := sectionMap[section.ResourceType]
		if ok {
			insertFieldEntriesFromSection(section, m)
			continue
		}
		imp.Rest[section.ResourceType] = make(map[string]interface{})
		insertFieldEntriesFromSection(section, imp.Rest[section.ResourceType])
	}

	return imp, nil
}
