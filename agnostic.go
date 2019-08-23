package validator

import (
	"fmt"
	"reflect"
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

func (b Builder) Confirmation(key string) {
	k, v := t(b.FormData, key)
	_, rV := t(b.Rules, valueToString(k)) // for confirmation, v = confirmation:"Emirhan"
	cnVal := transformFormValue(valueToString(rV))

	fmt.Println(v, rV, cnVal)

	if valueToString(v) != cnVal {
		fmt.Printf("%s field is not confirmed as '%s'", valueToString(k), cnVal)
	}
}

// i'ye gelecek değerin her zaman reflect.Value tipinde
// olduğu düşünülerek işlem yapılacak.
func valueToString(i interface{}) string {
	v := i.(reflect.Value)
	h := ""
	if v.Len() > 0 {
		h = v.Index(0).Elem().String()
	}

	return h
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
