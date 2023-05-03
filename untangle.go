package jl

import (
	"encoding/json"
	"fmt"
	"sort"
)

const (
	charObject = 123 // "{" -> char code 123
	charArray  = 91  // "[" -> char code 91
)

type (
	JsonObject map[string]json.RawMessage
	JsonArray  []json.RawMessage
)

func untangle(raw *json.RawMessage, pks *[]PathKey, flatters *[]Flatter) error {
	var firstChar = (*raw)[0]
	switch firstChar {
	case charObject:
		var j JsonObject
		err := json.Unmarshal(*raw, &j)
		if err != nil {
			return err
		}
		sortedKeys := sortedKeys(j)
		if len(sortedKeys) == 0 {
			return fmt.Errorf("blank object")
		}
		current := *pks
		for _, k := range sortedKeys {
			*pks = append(current, PathKey{keyType: keyTypeObject, key: k})
			h := j[k]
			untangle(&h, pks, flatters)
		}
	case charArray:
		var j JsonArray
		err := json.Unmarshal(*raw, &j)
		if err != nil {
			return err
		}
		if len(j) == 0 {
			return fmt.Errorf("blank array")
		}
		current := *pks
		for i, v := range j {
			*pks = append(current, PathKey{keyType: keyTypeArray, key: fmt.Sprint(i)})
			untangle(&v, pks, flatters)
		}
	default:
		var value any
		err := json.Unmarshal(*raw, &value)
		if err != nil {
			return err
		}
		switch v := value.(type) {
		case string:
			bv := []byte(v)
			if j := wouldBeJSON(&bv); j != nil {
				untangle(j, pks, flatters)
			} else {
				*flatters = append(*flatters, Flatter{pathKeys: *pks, value: v})
			}
		default:
			*flatters = append(*flatters, Flatter{pathKeys: *pks, value: v})
		}
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
