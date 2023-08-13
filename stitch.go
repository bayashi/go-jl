package jl

import (
	"encoding/json"
	"strconv"
)

// stitch converts untangled data to JSON string
func stitch(o *Options, flatters *[]Flatter) ([]byte, error) {
	var result any
	for _, flatter := range *flatters {
		result = stitchUp(flatter.pathKeys, flatter.value, result)
	}

	if o.NoPrettify {
		return json.Marshal(result)
	} else {
		return json.MarshalIndent(result, "", " ")
	}
}

func stitchUp(pks []PathKey, value any, result any) any {
	if len(pks) == 0 {
		return value
	}

	switch pks[0].keyType {
	case keyTypeObject:
		return stitchUpObject(pks[0].key, pks[1:], value, result)
	case keyTypeArray:
		i, _ := strconv.Atoi(pks[0].key)
		return stitchUpArray(i, pks[1:], value, result)
	}

	return nil
}

func stitchUpObject(key string, pks []PathKey, value any, result any) any {
	if result == nil {
		result = make(map[string]any)
	}

	re := result.(map[string]any)

	if len(pks) == 0 {
		re[key] = value
	} else {
		re[key] = stitchUp(pks, value, re[key])
	}

	return re
}

func stitchUpArray(key int, pks []PathKey, value any, result any) any {
	if result == nil {
		result = make([]any, 0)
	}

	re := expand(result.([]any), key+1)

	if len(pks) == 0 {
		re[key] = value
	} else {
		re[key] = stitchUp(pks, value, re[key])
	}

	return re
}

func expand(slice []any, size int) []any {
	for i := len(slice); i < size; i++ {
		slice = append(slice, nil)
	}

	return slice
}
