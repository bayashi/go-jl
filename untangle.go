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

// untangle converts JSON to the `Flatter` structure
func untangle(o *Options, raw *json.RawMessage, pks *[]PathKey, flatters *[]Flatter, decodeCount int) error {
	var firstChar = (*raw)[0]
	switch firstChar {
	case charObject:
		if err := untangleObject(o, raw, pks, flatters, decodeCount); err != nil {
			return err
		}
	case charArray:
		if err := untangleArray(o, raw, pks, flatters, decodeCount); err != nil {
			return err
		}
	default:
		if err := untangleValue(o, raw, pks, flatters, decodeCount); err != nil {
			return err
		}
	}

	return nil
}

func untangleObject(o *Options, raw *json.RawMessage, pks *[]PathKey, flatters *[]Flatter, decodeCount int) error {
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
		untangle(o, &h, pks, flatters, decodeCount)
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

func untangleArray(o *Options, raw *json.RawMessage, pks *[]PathKey, flatters *[]Flatter, decodeCount int) error {
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
		untangle(o, &j[i], pks, flatters, decodeCount)
	}

	return nil
}

func untangleValue(o *Options, raw *json.RawMessage, pks *[]PathKey, flatters *[]Flatter, decodeCount int) error {
	var value any
	err := json.Unmarshal(*raw, &value)
	if err != nil {
		return err
	}
	switch v := value.(type) {
	case string:
		if err := untangleStringValue(o, pks, flatters, v, decodeCount); err != nil {
			return err
		}
	default:
		*flatters = append(*flatters, Flatter{pathKeys: *pks, value: v})
	}

	return nil
}

func untangleStringValue(o *Options, pks *[]PathKey, flatters *[]Flatter, v string, decodeCount int) error {
	if decodeCount >= maxDecodeCount {
		*flatters = append(*flatters, Flatter{pathKeys: *pks, value: v})
		return nil
	}

	var bv []byte

	if o.SplitLF && strings.Contains(v, "\n") {
		var err error
		elements := strings.Split(v, "\n")
		bv, err = json.Marshal(elements)
		if err != nil {
			return err
		}
	} else if o.SplitTab && strings.Contains(v, "\t") {
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
		decodeCount++
		untangle(o, j, pks, flatters, decodeCount)
	} else {
		*flatters = append(*flatters, Flatter{pathKeys: *pks, value: v})
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
