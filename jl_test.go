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
		// not JSON
		{in: "", expect: ""},
		{in: "not json", expect: "not json"},
		{in: "{]", expect: "{]"},

		// Blank
		{in: "{}", expect: "{}"},
		{in: "[]", expect: "[]"},
		{in: "{[]}", expect: "{[]}"},
		{in: "[{},{}]", expect: "[{},{}]"},

		// Blank element
		{in: `{"a":{}}`, expect: `{"a":{}}`},
		{in: `{"a":[]}`, expect: `{"a":[]}`},
		{in: `["a",{}]`, expect: `["a",{}]`},

		// null element
		{in: `{"a":null}`, expect: `{"a":null}`},
		{in: `["a":null]`, expect: `["a":null]`},

		// JSON within JSON map
		{
			in:     `{"a":"{\"c\":\"{\\\"d\\\":34}\",\"b\":1.2}"}`,
			expect: `{"a":{"b":1.2,"c":{"d":34}}}`,
		},
		{
			in:     `{"a":"{\"b\":1.2,\"c\":\"{\\\"d\\\":\\\"{\\\\\\\"e\\\\\\\":\\\\\\\"aiko\\\\\\\"}\\\"}\"}"}`,
			expect: `{"a":{"b":1.2,"c":{"d":{"e":"aiko"}}}}`,
		},
		{
			in:     `{"a":"{\"c\":\"{\\\"d\\\":12,\\\"e\\\":34}\",\"b\":\"{\\\"d\\\":12,\\\"e\\\":34}\"}"}`,
			expect: `{"a":{"b":{"d":12,"e":34},"c":{"d":12,"e":34}}}`,
		},

		// JSON within JSON array
		{
			in:     `[["a","{\"b\":1.2}"]]`,
			expect: `[["a",{"b":1.2}]]`,
		},
		{
			in:     `[{"a":"{\"b\":\"[\\\"aiko\\\",\\\"eiko\\\"]\"}"}]`,
			expect: `[{"a":{"b":["aiko","eiko"]}}]`,
		},

		// JSON within JSON complecated
		{
			in:     `{"a":"{\"h\":\"{\\\"i\\\":12,\\\"j\\\":34,\\\"k\\\":\\\"{\\\\\\\"l\\\\\\\":\\\\\\\"[null,0,0.1,3,\\\\\\\\\\\\\\\"5\\\\\\\\\\\\\\\"]\\\\\\\"}\\\"}\"}","b":"[1,2,3,4,\"{\\\"g\\\":\\\"good\\\",\\\"c\\\":null,\\\"d\\\":1.2,\\\"e\\\":12}\"]"}`,
			expect: `{"a":{"h":{"i":12,"j":34,"k":{"l":[null,0,0.1,3,"5"]}}},"b":[1,2,3,4,{"c":null,"d":1.2,"e":12,"g":"good"}]}`,
		},
		{
			in:     `{"message": "{\"level\":\"info\",\"ts\":1557004280.5372975,\"caller\":\"zap/server_interceptors.go:40\",\"msg\":\"finished unary call with code OK\",\"grpc.start_time\":\"2019-05-04T21:11:20Z\",\"system\":\"grpc\",\"span.kind\":\"server\",\"grpc.service\":\"FooService\",\"grpc.method\":\"GetBar\",\"grpc.code\":\"OK\",\"grpc.time_ms\":248.45199584960938}\n","namespace": "foo-service","podName": "foo-86495899d8-m2vfl","containerName": "foo-service"}`,
			expect: `{"containerName":"foo-service","message":{"caller":"zap/server_interceptors.go:40","grpc.code":"OK","grpc.method":"GetBar","grpc.service":"FooService","grpc.start_time":"2019-05-04T21:11:20Z","grpc.time_ms":248.45199584960938,"level":"info","msg":"finished unary call with code OK","span.kind":"server","system":"grpc","ts":1557004280.5372975},"namespace":"foo-service","podName":"foo-86495899d8-m2vfl"}`,
		},
	}

	for _, tt := range tts {
		got := Process(&Options{}, []byte(tt.in))
		actually.Got(string(got)).Expect(tt.expect).X().Same(t)
	}
}

func prettified() string {
	return `{
 "a": {
  "b": 12
 }
}`
}

func TestProcessOptions(t *testing.T) {
	tts := []struct {
		in      string
		expect  string
		options *Options
	}{
		{
			in:      `{"a":"{\"b\":12}"}`,
			expect:  prettified(),
			options: &Options{Prettify: true},
		},
		{
			in:      `{"a":"key\t{\"b\":12}"}`,
			expect:  `{"a":{"key":{"b":12}}}`,
			options: &Options{SplitTab: true},
		},
		{
			in:      `{"a":"key\t{\"b\":12}\tfoo"}`,
			expect:  `{"a":["key",{"b":12},"foo"]}`,
			options: &Options{SplitTab: true},
		},
	}

	for _, tt := range tts {
		got := Process(tt.options, []byte(tt.in))
		actually.Got(string(got)).Expect(tt.expect).X().Same(t)
	}
}

// Example test for error output
func ExampleProcess() {
	Process(&Options{ShowErr: true}, []byte("{"))
	// Output:
	// unexpected end of JSON input
}
