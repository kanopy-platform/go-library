package merge_test

import (
	"testing"

	"github.com/kanopy-platform/go-library/merge"
	"github.com/stretchr/testify/assert"
)

func TestMergePatch(t *testing.T) {
	t.Parallel()

	tests := []struct {
		desc     string
		original string
		patch    string
		want     string
	}{
		{desc: "empty"},
		{
			desc:     "nested map",
			original: `{"map":{"one":{"i":1,"s":"one"}}}`,
			patch:    `{"map":{"one":{"i":2}}}`,
			want:     `{"map":{"one":{"i":2,"s":"one"}}}`,
		},
		{
			desc:     "null field",
			original: `{"cpu":"100m","memory":"1Gi"}`,
			patch:    `{"cpu":null}`,
			want:     `{"cpu":null,"memory":"1Gi"}`,
		},
	}

	for _, test := range tests {
		got, err := merge.PatchJSON([]byte(test.original), []byte(test.patch))
		assert.NoError(t, err, test.desc)
		assert.Equal(t, test.want, string(got), test.desc)
	}
}
