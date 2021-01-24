package godot

// Value is a wrapper for a regular TSCN value which also contains meta data
type Value struct {
	Value interface{}
	MetaData
}

type Type struct {
	Identifier string
	Parameters []interface{}
	MetaData
}

type KeyValuePair struct {
	Key   string
	Value interface{}
	MetaData
}
