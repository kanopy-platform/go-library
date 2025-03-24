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

// PatchJSON takes 2 []byte arguments and merges them as JSON Merge Patch documents.
// This can be used for general merging of JSON documents, but it also preserves `null`
// fields, which is useful if the merged document is used as a patch document or
// gets passed to a downstream encoder.
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

// PatchYAML does the same thing as PatchJSON, but with YAML []byte args.
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

// PatchInterface marshals interface objects to []byte and processes with PatchJSON.
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

// Patch marshals generic objects to []byte and processes with PatchJSON.
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
