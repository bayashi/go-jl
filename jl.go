package jl

import (
	"encoding/json"
)

// Flatter stores path data for each value
type Flatter struct {
	pathKeys []PathKey
	value    any
}

// PathKey represents a path
type PathKey struct {
	keyType KeyType
	key     string
}

// KeyType represents either an object or an array
type KeyType int

const (
	keyTypeObject KeyType = iota
	keyTypeArray
)

// Options is just an option data for a process
type Options struct {
	NoPrettify bool
	ShowErr    bool
	SplitTab   bool
	SplitLF    bool
}

// Process tries to convert "JSON within JSON" line to JUST nested JSON line.
// If there would be an error, return original JSON straightforward.
func Process(o *Options, origJson []byte) ([]byte, error) {
	var src json.RawMessage
	err := json.Unmarshal(origJson, &src)
	if err != nil {
		return origJson, err
	}

	pathKeys := []PathKey{}
	flatters := []Flatter{}

	c := &untangleCtx{
		o:           o,
		raw:         &src,
		pks:         &pathKeys,
		flatters:    &flatters,
		decodeCount: 0,
	}
	err2 := untangle(c)
	if err2 != nil {
		return origJson, err2
	}

	result, err3 := stitch(o, &flatters)
	if err3 != nil {
		return origJson, err3
	}

	return result, nil
}
