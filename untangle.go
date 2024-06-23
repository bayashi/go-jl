package jl

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

const (
	charObject = 123 // "{" -> char code 123
	charArray  = 91  // "[" -> char code 91

	maxDecodeCount = 99
)

type (
	JsonObject map[string]json.RawMessage
	JsonArray  []json.RawMessage
)

type untangleCtx struct {
	o           *Options
	raw         *json.RawMessage
	pks         *[]PathKey
	flatters    *[]Flatter
	decodeCount int
}

// untangle converts JSON to the `Flatter` structure
func untangle(c *untangleCtx) error {
	var firstChar = (*c.raw)[0]
	switch firstChar {
	case charObject:
		if err := untangleObject(c); err != nil {
			return err
		}
	case charArray:
		if err := untangleArray(c); err != nil {
			return err
		}
	default:
		if err := untangleValue(c); err != nil {
			return err
		}
	}

	return nil
}

func untangleObject(c *untangleCtx) error {
	var j JsonObject
	err := json.Unmarshal(*c.raw, &j)
	if err != nil {
		return err
	}
	sortedKeys := sortedKeys(j)
	if len(sortedKeys) == 0 {
		*c.flatters = append(*c.flatters, Flatter{pathKeys: *c.pks, value: map[string]any{}})
	}
	current := make([]PathKey, len(*c.pks))
	copy(current, *c.pks)
	for _, k := range sortedKeys {
		*c.pks = append(current, PathKey{keyType: keyTypeObject, key: k})
		h := j[k]
		c.raw = &h
		untangle(c)
	}

	return nil
}

func sortedKeys(obj JsonObject) []string {
	var keys []string
	for k := range obj {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	return keys
}

func untangleArray(c *untangleCtx) error {
	var j JsonArray
	err := json.Unmarshal(*c.raw, &j)
	if err != nil {
		return err
	}
	if len(j) == 0 {
		*c.flatters = append(*c.flatters, Flatter{pathKeys: *c.pks, value: []any{}})
	}
	current := make([]PathKey, len(*c.pks))
	copy(current, *c.pks)
	for i := range j {
		*c.pks = append(current, PathKey{keyType: keyTypeArray, key: fmt.Sprint(i)})
		c.raw = &j[i]
		untangle(c)
	}

	return nil
}

func untangleValue(c *untangleCtx) error {
	var value any
	err := json.Unmarshal(*c.raw, &value)
	if err != nil {
		return err
	}
	switch v := value.(type) {
	case string:
		if err := untangleStringValue(c, v); err != nil {
			return err
		}
	default:
		*c.flatters = append(*c.flatters, Flatter{pathKeys: *c.pks, value: v})
	}

	return nil
}

func untangleStringValue(c *untangleCtx, v string) error {
	if c.decodeCount >= maxDecodeCount {
		*c.flatters = append(*c.flatters, Flatter{pathKeys: *c.pks, value: v})
		return nil
	}

	var bv []byte

	if c.o.SplitLF && strings.Contains(v, "\n") {
		var err error
		elements := strings.Split(v, "\n")
		bv, err = json.Marshal(elements)
		if err != nil {
			return err
		}
	} else if c.o.SplitTab && strings.Contains(v, "\t") {
		var err error
		elements := strings.Split(v, "\t")
		if len(elements) == 2 {
			kv := map[string]string{}
			kv[elements[0]] = elements[1]
			bv, err = json.Marshal(kv)
		} else {
			bv, err = json.Marshal(elements)
		}
		if err != nil {
			return err
		}
	} else {
		bv = []byte(v)
	}

	if j := wouldBeJSON(&bv); j != nil {
		c.decodeCount++
		c.raw = j
		untangle(c)
	} else {
		*c.flatters = append(*c.flatters, Flatter{pathKeys: *c.pks, value: v})
	}

	return nil
}

func wouldBeJSON(src *[]byte) *json.RawMessage {
	if len(*src) == 0 {
		return nil
	}

	if firstChar := (*src)[0]; firstChar != charObject && firstChar != charArray {
		return nil
	}
	var jsonRaw json.RawMessage
	err := json.Unmarshal(*src, &jsonRaw)
	if err != nil {
		return nil // invalid
	}

	return &jsonRaw
}
