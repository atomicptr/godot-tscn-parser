package godot

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
