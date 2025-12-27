package jl

import (
	"encoding/json"
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

type processCtx struct {
	o           *Options
	decodeCount int
}

// processRecursive recursively processes JSON and expands any JSON strings found in values
// into nested structures, returning a map[string]any or []any
func processRecursive(c *processCtx, raw *json.RawMessage) (any, error) {
	if len(*raw) == 0 {
		return nil, nil
	}

	var firstChar = (*raw)[0]
	switch firstChar {
	case charObject:
		return processObject(c, raw)
	case charArray:
		return processArray(c, raw)
	default:
		return processValue(c, raw)
	}
}

func processObject(c *processCtx, raw *json.RawMessage) (any, error) {
	var j JsonObject
	err := json.Unmarshal(*raw, &j)
	if err != nil {
		return nil, err
	}

	result := make(map[string]any)
	sortedKeys := sortedKeys(j)

	for _, k := range sortedKeys {
		raw := j[k]
		value, err := processRecursive(c, &raw)
		if err != nil {
			return nil, err
		}
		result[k] = value
	}

	return result, nil
}

func sortedKeys(obj JsonObject) []string {
	var keys []string
	for k := range obj {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	return keys
}

func processArray(c *processCtx, raw *json.RawMessage) (any, error) {
	var j JsonArray
	err := json.Unmarshal(*raw, &j)
	if err != nil {
		return nil, err
	}

	result := make([]any, 0, len(j))
	for i := range j {
		value, err := processRecursive(c, &j[i])
		if err != nil {
			return nil, err
		}
		result = append(result, value)
	}

	return result, nil
}

func processValue(c *processCtx, raw *json.RawMessage) (any, error) {
	var value any
	err := json.Unmarshal(*raw, &value)
	if err != nil {
		return nil, err
	}

	switch v := value.(type) {
	case string:
		return processStringValue(c, v)
	default:
		return v, nil
	}
}

func processStringValue(c *processCtx, v string) (any, error) {
	if c.decodeCount >= maxDecodeCount {
		return v, nil
	}

	var bv []byte

	if c.o.SplitLF && strings.Contains(v, "\n") {
		var err error
		elements := strings.Split(v, "\n")
		bv, err = json.Marshal(elements)
		if err != nil {
			return nil, err
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
			return nil, err
		}
	} else {
		bv = []byte(v)
	}

	if j := wouldBeJSON(&bv); j != nil {
		c.decodeCount++
		return processRecursive(c, j)
	}

	return v, nil
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
