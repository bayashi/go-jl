package jl

import (
	"testing"

	"github.com/bayashi/actually"
)

func TestProcess(t *testing.T) {
	tts := []struct {
		in     string
		expect string
	}{
		{in: "", expect: ""},
		{in: "not json", expect: "not json"},
		{in: "{]", expect: "{]"},

		{in: "{}", expect: "{}"},
		{in: "{ }", expect: "{ }"},
		{in: "[]", expect: "[]"},
		{in: "[ ]", expect: "[ ]"},
		{in: "{[]}", expect: "{[]}"},

		{
			in:     `{"a":"{\"c\":\"{\\\"d\\\":34}\",\"b\":1.2}"}`,
			expect: `{"a":{"b":1.2,"c":{"d":34}}}`,
		},
		{
			in:     `{"a":"{\"b\":1.2,\"c\":\"{\\\"d\\\":\\\"{\\\\\\\"e\\\\\\\":\\\\\\\"aiko\\\\\\\"}\\\"}\"}"}`,
			expect: `{"a":{"b":1.2,"c":{"d":{"e":"aiko"}}}}`,
		},
		{
			in: `{"a":"{\"c\":\"{\\\"d\\\":12,\\\"e\\\":34}\",\"b\":\"{\\\"d\\\":12,\\\"e\\\":34}\"}"}`,
			expect: `{"a":{"b":{"d":12,"e":34},"c":{"d":12,"e":34}}}`,
		},

		{
			in:     `[["a","{\"b\":1.2}"]]`,
			expect: `[["a",{"b":1.2}]]`,
		},
		{
			in:     `[{"a":"{\"b\":\"[\\\"aiko\\\",\\\"eiko\\\"]\"}"}]`,
			expect: `[{"a":{"b":["aiko","eiko"]}}]`,
		},
	}

	for _, tt := range tts {
		got := Process(&Options{}, b(tt.in))
		actually.Got(got).Expect(b(tt.expect)).X().Same(t)
	}
}

func b(b string) []byte {
	return []byte(b)
}
