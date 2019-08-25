package validator

import (
	"fmt"
	"testing"
)

func TestSlice(t *testing.T) {
	var dump []string
	dump = []string{
		"v1", "v2", "v3",
	}

	iter := NewIterator(dump)

	data := iter.Slice(func(i, v interface{}) interface{} {
		return v
	})

	fmt.Println(data)

}

func TestMap(t *testing.T) {
	dump := map[string][]string{
		"t1": {"v1", "v2", "v3"},
		"t2": {"v1", "v2", "v3"},
		"t3": {"v1", "v2", "v3"},
		"t4": {"v1", "v2", "v3"},
	}

	//fmt.Println(mGetNthWithKey(dump, "t1", 0))
	//
	//var dump []string
	//dump = []string{
	//	"vv", "v33", "v555",
	//}

	iter := NewIterator(dump)
	data := iter.Map(func(k, v interface{}) interface{} {
		//nIter := NewIterator(v)
		//nIter.Slice(func(sV interface{}) interface{} {
		//	fmt.Println(sV)
		//	return sV
		//})
		return v
	})

	fmt.Println(data)

}
