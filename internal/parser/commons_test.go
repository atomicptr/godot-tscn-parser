package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type keyValuePair struct {
	Key   string
	Value interface{}
}

func assertField(t *testing.T, field *GdField, key string, values ...interface{}) bool {
	if field.Key != key {
		return false
	}

	switch field.Value.Raw().(type) {
	case GdType:
		typeVal, ok := field.Value.Raw().(GdType)
		assert.True(t, ok)
		assert.Equal(t, values[0], typeVal.Key)
		assert.Equal(t, len(values)-1, len(typeVal.Parameters))
		for index, param := range typeVal.Parameters {
			assert.Equal(t, values[index+1], param.Raw())
		}
	case []*GdValue:
		typeVal, ok := field.Value.Raw().([]*GdValue)
		assert.True(t, ok)
		for index, value := range typeVal {
			assert.Equal(t, values[index], value.Raw())
		}
	case []*GdMapField:
		typeVal, ok := field.Value.Raw().([]*GdMapField)
		assert.True(t, ok)
		for _, kv := range typeVal {
			for _, kvRaw := range values {
				kv2, ok := kvRaw.(keyValuePair)
				assert.True(t, ok)

				if kv.Key == kv2.Key {
					assert.Equal(t, kv2.Value, kv.Value.Raw())
				}
			}
		}
	default:
		assert.Equal(t, values[0], field.Value.Raw())
	}

	return true
}
