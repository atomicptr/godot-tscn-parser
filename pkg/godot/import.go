package godot

// Import is the configuration of well Godot imports (*.import file)
type Import struct {
	Remap  map[string]interface{}
	Deps   map[string]interface{}
	Params map[string]interface{}
	Rest   map[string]map[string]interface{}
	MetaData
}
