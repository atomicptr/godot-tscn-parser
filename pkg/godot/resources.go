package godot

type ExtResource struct {
	Path string
	Type string
	Id   int
}

type SubResource struct {
	Type   string
	Id     int
	Fields map[string]Value
}
