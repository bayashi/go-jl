package jl

import (
	"github.com/goccy/go-json"
)

// Options is just an option data for a process
type Options struct {
	NoPrettify bool
	ShowErr    bool
	SplitTab   bool
	SplitLF    bool
	Skip       int
}

// Process tries to convert "JSON within JSON" line to JUST nested JSON line.
// If there would be an error, return original JSON straightforward.
func Process(o *Options, origJson []byte) ([]byte, error) {
	if o.Skip > 0 && len(origJson) < o.Skip {
		return origJson, nil
	}

	var src json.RawMessage
	err := json.Unmarshal(origJson, &src)
	if err != nil {
		return origJson, err
	}

	c := &processCtx{
		o:           o,
		decodeCount: 0,
	}

	result, err2 := processRecursive(c, &src)
	if err2 != nil {
		return origJson, err2
	}

	if o.NoPrettify {
		return json.Marshal(result)
	} else {
		return json.MarshalIndent(result, "", " ")
	}
}
