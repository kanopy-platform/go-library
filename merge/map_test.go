package merge

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMergeMapStringString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		m    map[string]string
		args []map[string]string
		want map[string]string
	}{
		// test nil maps
		{
			m:    nil,
			args: []map[string]string{nil, nil},
			want: map[string]string{},
		},
		// test empty maps
		{
			m:    map[string]string{},
			args: []map[string]string{{}, {}},
			want: map[string]string{},
		},
		// test multiple maps
		{
			m: map[string]string{"a": "a"},
			args: []map[string]string{
				{"b": "b"},
				{"c": "c"},
			},
			want: map[string]string{"a": "a", "b": "b", "c": "c"},
		},
		// test overrides
		{
			m: map[string]string{"key": "a"},
			args: []map[string]string{
				{"key": "b"},
				{"key": "c"},
			},
			want: map[string]string{"key": "c"},
		},
	}

	for _, test := range tests {
		got := Maps(test.m, test.args...)
		assert.Equal(t, test.want, got)
	}
}

func TestTypes(t *testing.T) {
	got := Maps(map[int]bool{0: true}, map[int]bool{1: false, 2: true})
	assert.Equal(t, map[int]bool{0: true, 1: false, 2: true}, got)
}
