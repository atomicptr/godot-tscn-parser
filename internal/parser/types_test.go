package parser

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestTscnFileGetAttribute(t *testing.T) {
	content := "[gd_scene value1=1234 value2=\"test\"]"
	scene, err := Parse(strings.NewReader(content))
	assert.NoError(t, err)
	value1, err := scene.GetAttribute("value1")
	assert.NoError(t, err)
	assert.Equal(t, int64(1234), value1.Raw())
	value2, err := scene.GetAttribute("value2")
	assert.NoError(t, err)
	assert.Equal(t, "test", value2.Raw())
	_, err = scene.GetAttribute("value3")
	assert.Error(t, err)
}

func TestGdResourceGetAttribute(t *testing.T) {
	content := `[gd_scene]
[node value1=1234 value2="test"]`
	scene, err := Parse(strings.NewReader(content))
	assert.NoError(t, err)
	node := scene.Sections[0]
	value1, err := node.GetAttribute("value1")
	assert.NoError(t, err)
	assert.Equal(t, int64(1234), value1.Raw())
	value2, err := node.GetAttribute("value2")
	assert.NoError(t, err)
	assert.Equal(t, "test", value2.Raw())
	_, err = node.GetAttribute("value3")
	assert.Error(t, err)
}

func TestGdResourceGetField(t *testing.T) {
	content := `[gd_scene]
[node attr=1234]
value1=1234
value2="test"`
	scene, err := Parse(strings.NewReader(content))
	assert.NoError(t, err)
	node := scene.Sections[0]
	value1, err := node.GetField("value1")
	assert.NoError(t, err)
	assert.Equal(t, int64(1234), value1.Raw())
	value2, err := node.GetField("value2")
	assert.NoError(t, err)
	assert.Equal(t, "test", value2.Raw())
	_, err = node.GetField("value3")
	assert.Error(t, err)
}

func TestGdMapFieldToString(t *testing.T) {
	val := "value"
	kv := GdMapField{
		Key:   "string_field",
		Value: &GdValue{String: &val},
	}

	assert.Equal(t, `"string_field": "value"`, kv.ToString())
}

func TestGdValueRawReturnNil(t *testing.T) {
	val := GdValue{}
	assert.Nil(t, val.Raw())
}

func TestGdValueRawWithMap(t *testing.T) {
	content := `map={
"key1": "value1",
"key2": "value2"
}`
	scene, _ := Parse(strings.NewReader(content))

	m := scene.Fields[0]
	v := m.Value.Raw()

	assert.NotNil(t, v)
	assert.Len(t, v, 2)
}

func TestGdValueRawWithKeyValuePair(t *testing.T) {
	content := `obj_with_kv = Object("key":"value")`
	scene, _ := Parse(strings.NewReader(content))

	obj := scene.Fields[0]
	v := obj.Value.Type.Parameters[0].Raw()

	assert.NotNil(t, v)

	kv := v.(GdMapField)
	assert.Equal(t, "key", kv.Key)
	assert.Equal(t, "value", *kv.Value.String)
}

func TestGdValueRawWithArray(t *testing.T) {
	content := `array = [ 1, 2, 3 ]`
	scene, _ := Parse(strings.NewReader(content))

	arr := scene.Fields[0]
	values := arr.Value.Raw()

	assert.Len(t, values, 3)
}

func TestGdValueRawWithType(t *testing.T) {
	content := `obj = Object(1, 2, 3, 4)`
	scene, _ := Parse(strings.NewReader(content))

	obj := scene.Fields[0]
	gdTypeRaw := obj.Value.Raw()
	gdType := gdTypeRaw.(GdType)

	assert.Equal(t, "Object", gdType.Key)
	assert.Len(t, gdType.Parameters, 4)
}

func TestGdValueRawWithBasicTypes(t *testing.T) {
	content := `obj = Object("string", 42, -13.37, true, null)`
	scene, _ := Parse(strings.NewReader(content))

	obj := scene.Fields[0]
	gdTypeRaw := obj.Value.Raw()

	gdType := gdTypeRaw.(GdType)

	expectedParams := []interface{}{"string", int64(42), -13.37, true, nil}

	for index, param := range gdType.Parameters {
		actual := param.Raw()
		expected := expectedParams[index]
		assert.Equal(t, expected, actual)
	}
}

func TestGdValueToStringWithMap(t *testing.T) {
	val := "value1"
	val2 := "value2"
	m := GdValue{Map: []*GdMapField{
		{Key: "key1", Value: &GdValue{String: &val}},
		{Key: "key2", Value: &GdValue{String: &val2}},
	}}

	assert.Equal(t, `{"key1": "value1", "key2": "value2"}`, m.ToString())
}

func TestGdValueToStringWithKeyValuePair(t *testing.T) {
	val := "value"
	kv := GdMapField{
		Key:   "string_field",
		Value: &GdValue{String: &val},
	}
	value := GdValue{KeyValuePair: &kv}
	assert.Equal(t, `"string_field": "value"`, value.ToString())
}

func TestGdValueToStringWithArray(t *testing.T) {
	s, i, f, b, n := "str", int64(42), 13.37, true, true
	arr := GdValue{Array: []*GdValue{
		{String: &s},
		{Integer: &i},
		{Float: &f},
		{Bool: &b},
		{Null: &n},
	}}

	assert.Equal(t, `["str", 42, 13.370000, true, null]`, arr.ToString())
}

func TestGdValueToStringWithGdType(t *testing.T) {
	i := int64(42)
	v := GdValue{Type: &GdType{
		Key: "Object",
		Parameters: []*GdValue{
			{Integer: &i},
		},
	}}
	assert.Equal(t, `Object (42)`, v.ToString())
}

func TestGdValueToStringWithInvalidType(t *testing.T) {
	v := GdValue{}
	assert.Equal(t, "null", v.ToString())
}
