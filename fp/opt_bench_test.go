package fp_test

import (
	"encoding/json"
	"testing"

	"github.com/kanopy-platform/go-library/fp"
)

var benchSink any

type PointerStruct struct {
	Name     string  `json:"name"`
	Address  *string `json:"address,omitempty"`
	Priority *int    `json:"priority,omitempty"`
}

type OptStruct struct {
	Name     string         `json:"name"`
	Address  fp.Opt[string] `json:"address,omitzero"`
	Priority fp.Opt[int]    `json:"priority,omitzero"`
}

func ptrStr(s string) *string { return &s }
func ptrInt(n int) *int       { return &n }

// Struct-level benchmarks (marshal full struct through encoding/json)

func BenchmarkMarshal_Pointer_AllSet(b *testing.B) {
	v := PointerStruct{Name: "Alice", Address: ptrStr("123 Main St"), Priority: ptrInt(1)}
	for b.Loop() {
		benchSink, _ = json.Marshal(v)
	}
}

func BenchmarkMarshal_Opt_AllSet(b *testing.B) {
	v := OptStruct{Name: "Alice", Address: fp.Some("123 Main St"), Priority: fp.Some(1)}
	for b.Loop() {
		benchSink, _ = json.Marshal(v)
	}
}

func BenchmarkMarshal_Pointer_NoneSet(b *testing.B) {
	v := PointerStruct{Name: "Alice"}
	for b.Loop() {
		benchSink, _ = json.Marshal(v)
	}
}

func BenchmarkMarshal_Opt_NoneSet(b *testing.B) {
	v := OptStruct{Name: "Alice"}
	for b.Loop() {
		benchSink, _ = json.Marshal(v)
	}
}

func BenchmarkUnmarshal_Pointer(b *testing.B) {
	data := []byte(`{"name":"Alice","address":"123 Main St","priority":1}`)
	for b.Loop() {
		var v PointerStruct
		_ = json.Unmarshal(data, &v)
	}
}

func BenchmarkUnmarshal_Opt(b *testing.B) {
	data := []byte(`{"name":"Alice","address":"123 Main St","priority":1}`)
	for b.Loop() {
		var v OptStruct
		_ = json.Unmarshal(data, &v)
	}
}

// Direct MarshalJSON benchmarks (isolate the Opt overhead from struct encoding)

func BenchmarkMarshalJSON_Pointer_String(b *testing.B) {
	v := ptrStr("hello world")
	for b.Loop() {
		benchSink, _ = json.Marshal(v)
	}
}

func BenchmarkMarshalJSON_Opt_String(b *testing.B) {
	v := fp.Some("hello world")
	for b.Loop() {
		benchSink, _ = json.Marshal(v)
	}
}

func BenchmarkMarshalJSON_Pointer_Int(b *testing.B) {
	v := ptrInt(42)
	for b.Loop() {
		benchSink, _ = json.Marshal(v)
	}
}

func BenchmarkMarshalJSON_Opt_Int(b *testing.B) {
	v := fp.Some(42)
	for b.Loop() {
		benchSink, _ = json.Marshal(v)
	}
}

func BenchmarkMarshalJSON_Pointer_Nil(b *testing.B) {
	v := (*string)(nil)
	for b.Loop() {
		benchSink, _ = json.Marshal(v)
	}
}

func BenchmarkMarshalJSON_Opt_None(b *testing.B) {
	v := fp.None[string]()
	for b.Loop() {
		benchSink, _ = json.Marshal(v)
	}
}
