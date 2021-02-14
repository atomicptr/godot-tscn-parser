package godot

// Connection connects a Signal From a node To another node and calls its Method.
type Connection struct {
	From   string
	To     string
	Signal string
	Method string
	Flags  int64
	Binds  Value
	MetaData
}
