package godot

// FormatVersion defines the version of the TSCN file format
// Snippet from the C++ code: https://github.com/godotengine/godot/blob/master/scene/resources/resource_format_text.cpp#L39
const FormatVersion = 2

type Scene struct {
	// ExtResources is a map of external resources, key is their ID
	ExtResources map[int64]*ExtResource
	// SubResources is a map of internal resources, key is their ID
	SubResources map[int64]*SubResource
	// Node is the root node of the scene tree
	*Node
	// MetaData contains extra data like the lexer position
	MetaData
}
