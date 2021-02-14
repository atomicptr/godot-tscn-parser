package godot

type Connection struct {
	From   string
	To     string
	Signal string
	Method string
	Flags  int64
	Binds  Value
	MetaData
}
