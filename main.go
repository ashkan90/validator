package validator

import (
	"reflect"
	"strings"
)

const (
	R_REQ     = "required"
	R_MAX_LEN = "max-len"
	R_MIN_LEN = "min-len"
	R_EQ      = "equal"
	R_CONF    = "confirmation"
	R_DT      = "date"
	R_ARR     = "array"
	R_FILE    = "file"
	R_MA_SIZE = "max-size"
	R_MI_SIZE = "min-size"
	R_BTW     = "between"
	R_STR     = "string"
	R_INT     = "integer"
	R_DG      = "digit"
	R_EM      = "email"
)

const (
	token     = ":"
	sub_token = "-"
)

const confirmation_suffix = "_confirmation"

var inititator Initializer

type comparor struct {
	ruleField string
	ruleKey   string
	err       []string
	comp1     interface{}
	comp2     interface{}
}

type Rules struct {
	Plain    map[string][]string
	Validate *Validate
}

func (r *Rules) Prepare(rules map[string]string, form interface{}) {
	r.Validate.Form = form

	// initialize map for length of rules
	// if i didn't initialized it then it gives me error about
	// nil map
	fieldRules := make(map[string][]string, len(rules))
	var fieldNames []string

	r.Validate.Errors = make(map[string][]string)
	for key := range rules {
		r.Validate.Errors[key] = []string{}
	}

	for key, value := range rules {
		fieldNames = append(fieldNames, key)
		fieldRules[key] = strings.Split(value, "|")
	}

	r.Plain = fieldRules
	r.Validate.FieldNames = fieldNames

	inititator = Initializer{rules: r}

}

type Initializer struct {
	rules *Rules
}

func (i Initializer) Run() map[string][]string {
	i.rules.Required()     // done
	i.rules.Confirmation() // done
	i.rules.Equal()        // done
	i.rules.Date()         // done
	i.rules.Max()          // done
	i.rules.Min()          // done
	i.rules.Array()        // done
	i.rules.String()       // done
	i.rules.Integer()      // done
	i.rules.Digit()        // done
	i.rules.Between()      // done
	i.rules.File()         // done

	return i.rules.Validate.Errors
}

type Validate struct {
	Form            interface{}
	Errors          map[string][]string
	FieldNames      []string
	Internalization Internalization
}

func (v Validate) getFormValues(key string) []string {
	var found interface{}
	found = mGetNthWithKey(v.Form, key, -1)
	//for k, value := range v.Form.(url.Values) {
	//	if k == key {
	//		found = value
	//	}
	//}

	return found.([]string)
}

type Internalization struct{}

type IMessage interface{ Messages() map[string]string }

func (i Internalization) Messages() map[string]string {
	return map[string]string{
		R_REQ:     "{0} is required.",
		R_MAX_LEN: "{0} field's length cannot be greater than '{1}'",
		R_MIN_LEN: "{0} field's length cannot be lower than '{1}'",
		R_EQ:      "{0} field must be '{1}'",
		R_ARR:     "{0} field must be type of Array/Slice",
		R_STR:     "{0} field must be type of String",
		R_INT:     "{0} field must be type of Integer",
		R_DG:      "{0} field must be digit",
		R_BTW:     "{0} field must between {1}",
		R_CONF:    "{0} field is not confirmed.",
		R_DT:      "{0} field is not valid date format.",
		R_FILE:    "{0} field is not type of File.",
	}
}

func Load(rules map[string]string, form interface{}) Initializer {
	rStruct := new(Rules)
	rStruct.Validate = new(Validate)
	rStruct.Prepare(rules, form)
	return Initializer{
		rules: rStruct,
	}
}

// validator.Load().Prepare(...)
// validator.Run()

func t(d interface{}, k string) (interface{}, interface{}) {
	//var foundKeys interface{}
	//var foundValues interface{}

	foundKeys := reflect.New(reflect.TypeOf([]interface{}{})).Elem()
	foundValues := reflect.New(reflect.TypeOf([]interface{}{})).Elem()

	iter := reflect.ValueOf(d).MapRange()
	for iter.Next() {
		if iter.Key().String() == k {
			foundKeys = reflect.Append(foundKeys, iter.Key())
			foundValues = reflect.Append(foundValues, iter.Value())
		}
		//for i := 0; i < iter.Value().Len(); i++ {
		//	iVal := iter.Value().Index(i)
		//	fmt.Println(iVal)
		//	//if strings.Contains(iVal.String(), s) {
		//
		//	//	//foundKeys = append(foundKeys, iter.Key())
		//	//	//foundValues = append(foundValues, iVal)
		//	//}
		//}
	}

	return foundKeys, foundValues
}
