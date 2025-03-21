package merge

import (
	"encoding/json"

	jsonpatch "github.com/evanphx/json-patch/v5"
	"sigs.k8s.io/yaml"
)

func emptyBytes(b []byte) bool {
	return len(b) == 0 || string(b) == "null"
}

func argBytes(original, patch any) ([]byte, []byte, error) {
	originalBytes, err := json.Marshal(original)
	if err != nil {
		return nil, nil, err
	}

	patchBytes, err := json.Marshal(patch)
	if err != nil {
		return nil, nil, err
	}

	return originalBytes, patchBytes, nil
}

func PatchJSON(original, patch []byte) ([]byte, error) {
	switch {
	case emptyBytes(original):
		return patch, nil
	case emptyBytes(patch):
		return original, nil
	}

	// use MergeMergePatches to preserve `null` fields that can be passed on as a json patch.
	return jsonpatch.MergeMergePatches(original, patch)
}

func PatchYAML(original, patch []byte) ([]byte, error) {
	originalJSON, err := yaml.YAMLToJSON(original)
	if err != nil {
		return nil, err
	}

	patchJSON, err := yaml.YAMLToJSON(patch)
	if err != nil {
		return nil, err
	}

	return PatchJSON(originalJSON, patchJSON)
}

func PatchInterface(original, patch any) error {
	if original == nil || patch == nil {
		return nil
	}

	originalBytes, patchBytes, err := argBytes(original, patch)
	if err != nil {
		return err
	}

	mergedBytes, err := PatchJSON(originalBytes, patchBytes)
	if err != nil {
		return err
	}

	return json.Unmarshal(mergedBytes, original)
}

func Patch[T any](original, patch *T) (*T, error) {
	switch {
	case original == nil:
		return patch, nil
	case patch == nil:
		return original, nil
	}

	originalBytes, patchBytes, err := argBytes(original, patch)
	if err != nil {
		return nil, err
	}

	mergedBytes, err := PatchJSON(originalBytes, patchBytes)
	if err != nil {
		return nil, err
	}

	out := new(T)
	if err := json.Unmarshal(mergedBytes, out); err != nil {
		return nil, err
	}

	return out, nil
}
