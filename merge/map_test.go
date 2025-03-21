package merge

import (
	"reflect"
	"testing"
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
		got := MapStringString(test.m, test.args...)

		if !reflect.DeepEqual(test.want, got) {
			t.Errorf("%#v is not equal to %#v", test.want, got)
		}
	}
}
