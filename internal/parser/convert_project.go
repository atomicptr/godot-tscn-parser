package parser

import "github.com/atomicptr/godot-tscn-parser/pkg/godot"

// ConvertToGodotScene tries to convert a TscnFile structure to an actual Godot Scene with a node tree
func (tscn *TscnFile) ConvertToGodotProject() (*godot.Project, error) {
	project := &godot.Project{
		Android:     make(map[string]interface{}),
		Audio:       make(map[string]interface{}),
		Application: make(map[string]interface{}),
		Autoload:    make(map[string]interface{}),
		Compression: make(map[string]interface{}),
		Display:     make(map[string]interface{}),
		Editor:      make(map[string]interface{}),
		Filesystem:  make(map[string]interface{}),
		GUI:         make(map[string]interface{}),
		Input:       make(map[string]interface{}),
		LayerNames:  make(map[string]interface{}),
		Locale:      make(map[string]interface{}),
		Logging:     make(map[string]interface{}),
		Memory:      make(map[string]interface{}),
		Network:     make(map[string]interface{}),
		Node:        make(map[string]interface{}),
		Rendering:   make(map[string]interface{}),
		World:       make(map[string]interface{}),
		Rest:        make(map[string]map[string]interface{}),
		Fields:      make(map[string]interface{}),
	}

	insertFieldEntriesFromSection(&GdResource{Fields: tscn.Fields}, project.Fields)

	sectionMap := map[string]map[string]interface{}{
		"android":     project.Android,
		"audio":       project.Audio,
		"application": project.Application,
		"autoload":    project.Autoload,
		"compression": project.Compression,
		"display":     project.Display,
		"editor":      project.Editor,
		"filesystem":  project.Filesystem,
		"gui":         project.GUI,
		"input":       project.Input,
		"layernames":  project.LayerNames,
		"locale":      project.Locale,
		"logging":     project.Logging,
		"memory":      project.Memory,
		"network":     project.Network,
		"node":        project.Node,
		"rendering":   project.Rendering,
		"world":       project.World,
	}

	for _, section := range tscn.Sections {
		m, ok := sectionMap[section.ResourceType]
		if ok {
			insertFieldEntriesFromSection(section, m)
			continue
		}
		project.Rest[section.ResourceType] = make(map[string]interface{})
		insertFieldEntriesFromSection(section, project.Rest[section.ResourceType])
	}

	return project, nil
}
