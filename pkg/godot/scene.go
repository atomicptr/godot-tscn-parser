package godot

// FormatVersion defines the version of the TSCN file format
// Snippet from the C++ code: https://github.com/godotengine/godot/blob/master/scene/resources/resource_format_text.cpp#L39
const FormatVersion = 2

type Scene struct {
	ExtResources map[int64]*ExtResource
	SubResources map[int64]*SubResource
	RootNode     *Node
	MetaData
}
