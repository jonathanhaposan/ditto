package main

import (
	"encoding/json"
	"github.com/google/go-cmp/cmp"
	"reflect"
)

func IsEqualJson(s1, s2 []byte) bool {
	var o1 interface{}
	var o2 interface{}

	err := json.Unmarshal(s1, &o1)
	if err != nil {
		return false
	}

	err = json.Unmarshal(s2, &o2)
	if err != nil {
		return false
	}

	v1 := reflect.ValueOf(o1)
	v2 := reflect.ValueOf(o2)
	return compare(v1, v2)
}

func compare(v1, v2 reflect.Value) bool {
	if v1.Kind() == reflect.String {
		return cmp.Equal(v1.Interface(), v2.Interface())
	} else if v1.Kind() == reflect.Map {
		match := true
		for _, key := range v1.MapKeys() {
			valV1 := v1.MapIndex(key)
			valV2 := v2.MapIndex(key)
			if !valV2.IsValid()|| valV2.IsNil(){
				continue
			}
			o1 := reflect.ValueOf(valV1.Interface())
			o2 := reflect.ValueOf(valV2.Interface())
			if !compare(o1, o2) {
				match = false
				break
			}
		}
		return match
	} else if v1.Kind() == reflect.Slice {
		var ix int
		match := true
		for ix < v1.Len() || ix < v2.Len() {
			valV1 := v1.Index(ix)
			valV2 := v2.Index(ix)
			o1 := reflect.ValueOf(valV1.Interface())
			o2 := reflect.ValueOf(valV2.Interface())
			if !compare(o1, o2) {
				match = false
			}
			ix++
		}
		return match
	}
	return false
}
