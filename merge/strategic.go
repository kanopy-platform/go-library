package merge

import (
	"encoding/json"

	"k8s.io/apimachinery/pkg/util/strategicpatch"
)

// Strategic takes any number of generic objects and merges them using the strategic merge patch strategy.
// https://github.com/kubernetes/community/blob/6690abcd6b833f46550f5eaba2ec17a9e39b35c4/contributors/devel/sig-api-machinery/strategic-merge-patch.md
func Strategic[V any](original, modified *V, extra ...*V) (*V, error) {
	if original == nil {
		return nil, nil
	}

	originalBytes, err := json.Marshal(original)
	if err != nil {
		return nil, err
	}

	out := new(V)
	objects := append([]*V{modified}, extra...)

	for _, o := range objects {
		if o == nil {
			continue
		}

		modifiedBytes, err := json.Marshal(o)
		if err != nil {
			return nil, err
		}

		originalBytes, err = strategicpatch.StrategicMergePatch(originalBytes, modifiedBytes, out)
		if err != nil {
			return nil, err
		}
	}

	if err := json.Unmarshal(originalBytes, out); err != nil {
		return nil, err
	}

	return out, nil
}
