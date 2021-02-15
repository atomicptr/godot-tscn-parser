package convert

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntArrayContains(t *testing.T) {
	table := []struct {
		num      int
		arr      []int
		expected bool
	}{
		{1, []int{1}, true},
		{2, []int{1}, false},
		{1, []int{1, 2, 3, 4}, true},
		{4, []int{1, 2, 3, 4}, true},
		{2, []int{1, 2, 3, 4}, true},
		{2, []int{}, false},
		{2, nil, false},
	}

	for _, tc := range table {
		assert.Equal(t, tc.expected, intArrayContains(tc.num, tc.arr))
	}
}
