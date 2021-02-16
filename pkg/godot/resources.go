package godot

// Resource is the model for a .tres file
type Resource struct {
	// Type determines the resource type
	Type string
	// ExtResources is a map of external resources, key is their ID
	ExtResources map[int64]*ExtResource
	// SubResources is a map of internal resources, key is their ID
	SubResources map[int64]*SubResource
	// Fields contains the fields attached to the resource
	Fields map[string]interface{}
	// MetaData contains extra data like the lexer position
	MetaData
}

// ExtResource is a link to resources not contained within the TSCN file itself
// Documentation: https://docs.godotengine.org/en/stable/development/file_formats/tscn.html#external-resources
type ExtResource struct {
	Path string
	Type string
	ID   int64
	MetaData
}

// SubResource a TSCN file can contain meshes, materials and other data which are contained within this type
// Documentation: https://docs.godotengine.org/en/stable/development/file_formats/tscn.html#internal-resources
type SubResource struct {
	Type   string
	ID     int64
	Fields map[string]interface{}
	MetaData
}
