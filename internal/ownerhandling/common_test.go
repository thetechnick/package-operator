package ownerhandling

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemove(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		slice    []string
		index    int
		expected []string
	}{
		{
			name:     "remove from middle",
			slice:    []string{"a", "b", "c", "d"},
			index:    1,
			expected: []string{"a", "d", "c"},
		},
		{
			name:     "remove first element",
			slice:    []string{"a", "b", "c"},
			index:    0,
			expected: []string{"c", "b"},
		},
		{
			name:     "remove last element",
			slice:    []string{"a", "b", "c"},
			index:    2,
			expected: []string{"a", "b"},
		},
		{
			name:     "remove from single element slice",
			slice:    []string{"a"},
			index:    0,
			expected: []string{},
		},
		{
			name:     "remove from two element slice - first",
			slice:    []string{"a", "b"},
			index:    0,
			expected: []string{"b"},
		},
		{
			name:     "remove from two element slice - second",
			slice:    []string{"a", "b"},
			index:    1,
			expected: []string{"a"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := remove(tt.slice, tt.index)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRemoveWithInts(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		slice    []int
		index    int
		expected []int
	}{
		{
			name:     "remove from middle",
			slice:    []int{1, 2, 3, 4},
			index:    1,
			expected: []int{1, 4, 3},
		},
		{
			name:     "remove first element",
			slice:    []int{1, 2, 3},
			index:    0,
			expected: []int{3, 2},
		},
		{
			name:     "remove last element",
			slice:    []int{1, 2, 3},
			index:    2,
			expected: []int{1, 2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := remove(tt.slice, tt.index)
			assert.Equal(t, tt.expected, result)
		})
	}
}
