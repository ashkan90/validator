package validator

import (
	"reflect"
	"strconv"
	"strings"
)

//Slice
func sArrayContainsThenRetrieve(d interface{}, c string) (bool, interface{}) {

	//for _, value := range d {
	//	switch value.(type) {
	//	case []string:
	//		iter := value.([]string)
	//		for _, v := range iter {
	//			if strings.Contains(v, c) {
	//				return true, v
	//			}
	//		}
	//	}
	//}
	var found bool
	var fVal interface{}
	iter := NewIterator(d)
	iter.Slice(func(i, v interface{}) interface{} {
		if strings.Contains(v.(string), c) {
			found = true
			fVal = v
		}

		return ""
	})

	return found, fVal
}

// Map
func mSearchAndRetrieve(d interface{}, key interface{}) interface{} {

	var fVal interface{}

	iter := NewIterator(d)
	iter.Map(func(k, v interface{}) interface{} {
		if key == k {
			fVal = v
		}
		return ""
	})
	return fVal
}
func mGetNthWithKey(d interface{}, k string, nth int) interface{} {

	var fVal interface{}
	sliceOfVal := mSearchAndRetrieve(d, k)

	if nth != -1 {
		fVal = index(sliceOfVal, nth)
	} else {
		fVal = sliceOfVal
	}

	return fVal
}

func index(d interface{}, i int) interface{} {
	var found interface{}
	v := reflect.ValueOf(d)
	if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
		found = v.Index(i).Interface()
	}

	return found
}

func lenOnKey(d interface{}, k interface{}) int {
	sliceOfVal := mSearchAndRetrieve(d, k)
	return reflect.ValueOf(sliceOfVal).Len()
}

func isDigit(s string) bool {
	var holding bool
	for _, value := range s {
		v, ok := strconv.Atoi(string(value))
		if ok == nil {
			if v >= 0 && v <= 10 {
				holding = true
			}
		}
	}

	return holding
}
