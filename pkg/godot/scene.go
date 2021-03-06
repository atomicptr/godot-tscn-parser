package godot

// Scene is the data representation of a Godot TSCN Scene, which contains resources and a node tree
type Scene struct {
	// ExtResources is a map of external resources, key is their ID
	ExtResources map[int64]*ExtResource
	// SubResources is a map of internal resources, key is their ID
	SubResources map[int64]*SubResource
	// Node is the root node of the scene tree
	*Node
	// Editables is a list of editable scenes
	Editables []*Editable
	// Connections is a list of event triggers and to which nodes they're connected to
	Connections []*Connection
	// MetaData contains extra data like the lexer position
	MetaData
}
