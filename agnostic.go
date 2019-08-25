package validator

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type Builder struct {
	FormData interface{}
	Rules    interface{}
}

func (b *Builder) Load(f interface{}, r interface{}) {
	b.FormData = f
	b.Rules = r
}

func (b Builder) Required(key string) {
	k, v := t(b.FormData, key)
	if inLen(v) == 0 {
		fmt.Printf("%s is required.", k)
	}
}
func (b Builder) Confirmation(key string) {}
func (b Builder) Max(key string) {
	k, v := t(b.FormData, key)
	fmt.Println(len(valueToStrSlice(v)))
	if len(valueToStrSlice(v)) > 0 {
		oV := valueToStrSlice(v)[0] // overload v variable

		_, rV := t(b.Rules, valueToString(k)) // for equal, rV = max:10
		nrV := findFor("max", rV)

		if len(nrV) > 0 {
			cnVal := transformFormValue(nrV[0])
			iV, _ := strconv.Atoi(cnVal)

			if len(oV) > iV {
				fmt.Printf("%s field cannot be greater than %d", valueToString(k), iV)
			}
		}
	}
}
func (b Builder) Min(key string) {
	k, v := t(b.FormData, key)

	if len(valueToStrSlice(v)) > 0 {
		oV := valueToStrSlice(v)[0] // overload v variable

		_, rV := t(b.Rules, valueToString(k))
		nrV := findFor("min", rV) // for equal, rV = min:10

		if len(nrV) > 0 {
			cnVal := transformFormValue(nrV[0])
			iV, _ := strconv.Atoi(cnVal)

			if len(oV) < iV {
				fmt.Printf("%s field cannot be lower than %d", valueToString(k), iV)
			}
		}
	}

}
func (b Builder) Equal(key string) {
	k, v := t(b.FormData, key)
	_, rV := t(b.Rules, valueToString(k)) // for equal, rV = equal:"Emirhan"
	nrV := findFor("equal", rV)
	cnVal := transformFormValue(nrV[0])

	if valueToString(v) != cnVal {
		fmt.Printf("%s field is not equal with '%s'", valueToString(k), cnVal)
	}
}

func findFor(k string, i interface{}) []string {
	var found []string
	if isValue(i) {
		v := i.(reflect.Value)
		if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
			vals := valueToStrSlice(i)
			for _, value := range vals {
				if strings.Contains(value, k) {
					found = append(found, value)
				}
			}
		} else {
			fmt.Printf("unexpected reflect type. given: %v", v)
		}
	}

	return found
}

// i'ye gelecek değerin her zaman reflect.Value tipinde
// olduğu düşünülerek işlem yapılacak.
func valueToString(i interface{}) string {
	v := i.(reflect.Value)
	h := ""

	//
	//if v.Index(0).Kind() != reflect.String {
	//	fmt.Println(v.Index(0).Elem().Index(0).String())
	//}
	if v.Len() > 0 {
		h = v.Index(0).Elem().String()
	}

	return h
}
func valueToStrSlice(i interface{}) []string {
	var t []string
	v := i.(reflect.Value)
	if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
		for n := 0; n < v.Index(0).Elem().Len(); n++ {
			t = append(t, v.Index(0).Elem().Index(n).String())
		}
	}

	return t
}

func transformFormValue(val string) string {
	var v string
	if strings.Contains(val, ":") {
		v = strings.SplitAfter(val, ":")[1]
	}

	return v
}

func inLen(i interface{}) int {
	switch i.(type) {
	case reflect.Value:
		return i.(reflect.Value).Len()
	case []interface{}:
		return len(i.([]interface{}))
	}

	return -1
}

func isValue(i interface{}) bool {
	switch i.(type) {
	case reflect.Value:
		return true
	}

	return false
}
