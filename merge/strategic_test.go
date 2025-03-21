package merge_test

import (
	"testing"

	"github.com/kanopy-platform/go-library/merge"
	"github.com/stretchr/testify/assert"
)

type mergeItem struct {
	String string            `json:"string,omitempty"`
	Slice  []string          `json:"slice,omitempty" patchStrategy:"merge"`
	Map    map[string]string `json:"map,omitempty"`
}

// the upstream package has more robust testing
// these are here primarily for demonstration purposes/documentation
func TestStrategic(t *testing.T) {
	t.Parallel()

	tests := []struct {
		desc     string
		original *mergeItem
		modified *mergeItem
		extra    []*mergeItem
		want     *mergeItem
	}{
		{
			desc:     "test empty modified object",
			original: &mergeItem{String: "a", Slice: []string{"a"}},
			modified: &mergeItem{},
			want:     &mergeItem{String: "a", Slice: []string{"a"}},
		},
		// the logic of merging slices isn't very straight forward and more info can be found here:
		// https://github.com/kubernetes/design-proposals-archive/blob/b9d542b4fec7264362f3c31ce5f63738a32d115f/cli/preserve-order-in-strategic-merge-patch.md#when-setelementorder-is-not-present-and-patching-a-list
		{
			desc:     "test merging of slices",
			original: &mergeItem{Slice: []string{"b", "c", "a"}},
			modified: &mergeItem{Slice: []string{"a", "b", "d"}},
			want:     &mergeItem{Slice: []string{"c", "a", "b", "d"}},
		},
		{
			desc:     "test merging of maps",
			original: &mergeItem{Map: map[string]string{"o": "o", "test": "o"}},
			modified: &mergeItem{Map: map[string]string{"m": "m", "test": "m"}},
			want:     &mergeItem{Map: map[string]string{"o": "o", "m": "m", "test": "m"}},
		},
		{
			desc:     "test merging multiple objects",
			original: &mergeItem{Map: map[string]string{"o": "o", "test": "o"}},
			modified: &mergeItem{Map: map[string]string{"m": "m", "test": "m"}},
			extra: []*mergeItem{
				{Map: map[string]string{"e1": "e1", "test": "e1"}},
				{Map: map[string]string{"e2": "e2", "test": "e2"}},
			},
			want: &mergeItem{Map: map[string]string{"o": "o", "m": "m", "e1": "e1", "e2": "e2", "test": "e2"}},
		},
	}

	for _, test := range tests {
		o, err := merge.Strategic(test.original, test.modified, test.extra...)
		assert.NoError(t, err)
		assert.Equal(t, test.want, o, test.desc)
	}
}
