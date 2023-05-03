package jl

import (
	"encoding/json"
	"os"
)

type Flatter struct {
	pathKeys []PathKey
	value    any
}

type PathKey struct {
	keyType KeyType
	key     string
}

type KeyType int

const (
	keyTypeObject KeyType = iota
	keyTypeArray
)

type Options struct {
	Prettify bool
	ShowErr  bool
}

// Process tries to convert "JSON within JSON" line to JUST Nested JSON line.
// If there would be an error, return original JSON straightfarward.
func Process(o *Options, origJson []byte) []byte {
	var src json.RawMessage
	err := json.Unmarshal(origJson, &src)
	if err != nil {
		if o.ShowErr {
			os.Stdout.Write([]byte(err.Error()))
		}
		return origJson
	}
	pathKeys := []PathKey{}
	flatters := []Flatter{}
	err2 := untangle(&src, &pathKeys, &flatters)
	if err2 != nil {
		if o.ShowErr {
			os.Stdout.Write([]byte(err2.Error()))
		}
		return origJson
	}

	result, err3 := stitch(o, &flatters)
	if err3 != nil {
		if o.ShowErr {
			os.Stdout.Write([]byte(err3.Error()))
		}
		return origJson
	}

	return result
}
