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
			in:     `{"a":"{\"c\":\"{\\\"d\\\":34}\",\"b\":12}"}`,
			expect: `{"a":{"b":12,"c":{"d":34}}}`,
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
