package godot

// Value is a wrapper for a regular TSCN value which also contains meta data
type Value struct {
	Value interface{}
	MetaData
}

// Type represents more or less something akin to a struct in TSCN files, for instance a Vector
type Type struct {
	Identifier string
	Parameters []interface{}
	MetaData
}

// KeyValuePair a pair which contains a value with its associated key
type KeyValuePair struct {
	Key   string
	Value interface{}
	MetaData
}
