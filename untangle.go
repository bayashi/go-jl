package jl

import (
	"encoding/json"
	"fmt"
	"sort"
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

// untangle converts JSON to the `Flatter` structure
func untangle(raw *json.RawMessage, pks *[]PathKey, flatters *[]Flatter, decodeCount int) error {
	var firstChar = (*raw)[0]
	switch firstChar {
	case charObject:
		if err := untangleObject(raw, pks, flatters, decodeCount); err != nil {
			return err
		}
	case charArray:
		if err := untangleArray(raw, pks, flatters, decodeCount); err != nil {
			return err
		}
	default:
		if err := untangleValue(raw, pks, flatters, decodeCount); err != nil {
			return err
		}
	}

	return nil
}

func untangleObject(raw *json.RawMessage, pks *[]PathKey, flatters *[]Flatter, decodeCount int) error {
	var j JsonObject
	err := json.Unmarshal(*raw, &j)
	if err != nil {
		return err
	}
	sortedKeys := sortedKeys(j)
	if len(sortedKeys) == 0 {
		*flatters = append(*flatters, Flatter{pathKeys: *pks, value: map[string]any{}})
	}
	current := make([]PathKey, len(*pks))
	copy(current, *pks)
	for _, k := range sortedKeys {
		*pks = append(current, PathKey{keyType: keyTypeObject, key: k})
		h := j[k]
		untangle(&h, pks, flatters, decodeCount)
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

func untangleArray(raw *json.RawMessage, pks *[]PathKey, flatters *[]Flatter, decodeCount int) error {
	var j JsonArray
	err := json.Unmarshal(*raw, &j)
	if err != nil {
		return err
	}
	if len(j) == 0 {
		*flatters = append(*flatters, Flatter{pathKeys: *pks, value: []any{}})
	}
	current := make([]PathKey, len(*pks))
	copy(current, *pks)
	for i := range j {
		*pks = append(current, PathKey{keyType: keyTypeArray, key: fmt.Sprint(i)})
		untangle(&j[i], pks, flatters, decodeCount)
	}

	return nil
}

func untangleValue(raw *json.RawMessage, pks *[]PathKey, flatters *[]Flatter, decodeCount int) error {
	var value any
	err := json.Unmarshal(*raw, &value)
	if err != nil {
		return err
	}
	switch v := value.(type) {
	case string:
		untangleStringValue(pks, flatters, v, decodeCount)
	default:
		*flatters = append(*flatters, Flatter{pathKeys: *pks, value: v})
	}

	return nil
}

func untangleStringValue(pks *[]PathKey, flatters *[]Flatter, v string, decodeCount int) {
	if decodeCount >= maxDecodeCount {
		*flatters = append(*flatters, Flatter{pathKeys: *pks, value: v})
		return
	}

	bv := []byte(v)
	if j := wouldBeJSON(&bv); j != nil {
		decodeCount++
		untangle(j, pks, flatters, decodeCount)
	} else {
		*flatters = append(*flatters, Flatter{pathKeys: *pks, value: v})
	}
}

func wouldBeJSON(src *[]byte) *json.RawMessage {
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
