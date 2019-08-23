package validator

import (
	"fmt"
	"github.com/ashkan90/golang-in_array"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
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

func (r *Rules) ruleFinder(ruleKey string) (bool, string, string) {
	var found bool
	var foundKey string
	var foundValue string
	for key, value := range r.Plain {
		for _, v := range value {

			if strings.Contains(v, ruleKey) {
				found = true
				foundValue = v
				foundKey = key
			}
		}
	}

	return found, foundKey, foundValue
}
func (r *Rules) ruleFinderWithFormKey(key string) ([]string, []string) {
	keys, values := t(r.Validate.Form, key)

	return keys.([]string), values.([]string)
}
func (r *Rules) ruleFinderAsSlice(ruleKey string) ([]string, []string) {
	var foundKey []string
	var foundValue []string
	for key, value := range r.Plain {
		for _, v := range value {

			if strings.Contains(v, ruleKey) {
				foundValue = append(foundValue, v)
				foundKey = append(foundKey, key)
			}
		}
	}

	return foundKey, foundValue
}
func (r *Rules) ruleComparerProcessor(ruleKey string) {
	ok, k, v := r.ruleFinder(ruleKey)     // k = rule daki key değeri örn: 'name', 'surname'
	values := r.Validate.getFormValues(k) // k'daki değerler dizgesi.

	if strings.Contains(v, token) {
		v = strings.SplitAfter(v, token)[1] // "equal", "23", [1] == 23
	}
	if ok {
		comp := &comparor{
			ruleField: k,
			ruleKey:   ruleKey,
			err:       []string{k, v},
		}
		if ruleKey == R_ARR {
			comp.comp1 = v
			comp.comp2 = values
			r.ruleComparor(*comp)
		} else {
			for _, value := range values {
				comp.comp1 = v
				comp.comp2 = value
				r.ruleComparor(*comp)
			}
		}
	}
}
func (r *Rules) ruleComparor(c comparor) {

	switch c.ruleKey {
	case R_EQ:
		// give error if form value and equal value doesn't match
		if c.comp1 != c.comp2 {
			r.Validate.addError(c.ruleField, c.ruleKey, c.err)
		}
	case R_CONF:
		// if confirmation request's input name don't end with 'confirmation_suffix'
		// we don't need to check.
		if strings.Contains(c.comp1.(string), confirmation_suffix) {
			// incoming confirmation value. eg. value of password.
			inConfirmation := c.comp2.(string)
			// must confirmation with incoming confirmation.
			// this key should contain '_confirmation' suffix.
			// we'll search for confirmation value by 'c.comp1'(password_confirmation)
			mConfirmationKey := c.comp1.(string)
			mConfirmationVal := r.Validate.Form.(url.Values)[mConfirmationKey][0]

			if inConfirmation != mConfirmationVal {
				r.Validate.addError(c.ruleField, c.ruleKey, c.err)
			}
		}
	case R_MAX_LEN:
		// max-len returning value as string so we need to assert type
		if isDigit(c.comp1.(string)) {
			comp1, _ := strconv.Atoi(c.comp1.(string)) // value of rule eg. '10', '255'
			comp2 := c.comp2.(string)                  // value of form input eg. 'emirhan'

			// comp1 is value of max-len which is already type of int.
			// and comp2 is form value type of string, we should be sure
			// comp2's length is greater than comp1
			// eg. comp1 = 50, comp2 = "Hello world" // len = 11
			// Q: Is 11 greater than 50 ?
			// A: No.
			if comp1 < len(comp2) {
				r.Validate.addError(c.ruleField, c.ruleKey, c.err)
			}
		}
	case R_MIN_LEN:
		comp1, _ := strconv.Atoi(c.comp1.(string))
		comp2 := c.comp2.(string)
		if comp1 > len(comp2) {
			r.Validate.addError(c.ruleField, c.ruleKey, c.err)
		}
	case R_ARR:
		// comp1 always be 1 in this rule
		// because we can determine any variable to its type
		// in this case, if length is greater than 1 then do nothing
		// but if incoming form value's length lower than 1 then
		// we need to push an error for field.
		comp1 := 1
		comp2 := c.comp2.([]string)
		if len(comp2) <= comp1 {
			r.Validate.addError(c.ruleField, c.ruleKey, c.err)
		}
	case R_STR:
		// Basically, we're getting type of interface which is c.comp2
		// in this case, we're checking for type of string inside of c.comp2
		// so we can easily use reflect package because of it's easy
		// if type is not String then push an error.
		typ := reflect.TypeOf(c.comp2).Kind()
		if typ != reflect.String {
			r.Validate.addError(c.ruleField, c.ruleKey, c.err)
		}
	case R_INT:
		// Same as checking string. There's checking integer type.
		typ := reflect.TypeOf(c.comp2).Kind()
		if typ != reflect.Int {
			r.Validate.addError(c.ruleField, c.ruleKey, c.err)
		}
	case R_DG:
		typ := reflect.TypeOf(c.comp2).Kind()
		if typ == reflect.String && !isDigit(c.comp2.(string)) {
			r.Validate.addError(c.ruleField, c.ruleKey, c.err)
		} else if typ == reflect.Int {
			r.Validate.addError(c.ruleField, c.ruleKey, c.err)
		}
	case R_BTW:
		comparator, _ := strconv.Atoi(c.comp2.(string))
		values := strings.Split(c.comp1.(string), sub_token)
		to, _ := strconv.Atoi(values[0])
		from, _ := strconv.Atoi(values[1])

		if !(comparator > to && comparator < from) {
			r.Validate.addError(c.ruleField, c.ruleKey, c.err)
		}
	case R_DT:
		// default, date rule uses "Y-m-D".
		// default, date location uses "Europe/Berlin".
		loc, _ := time.LoadLocation("Europe/Berlin")
		const ymd = "2006-01-02"

		_, err := time.ParseInLocation(ymd, c.comp2.(string), loc)
		if err != nil {
			var _err []string
			_err = append(_err, err.Error()+" "+c.ruleField)
			r.Validate.addError(c.ruleField, c.ruleKey, _err)
		}
	case R_FILE:
		fmt.Println(c)
	}
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

func (r *Rules) Confirmation() {
	if exists.In_array(R_CONF, r.Plain, false) {
		r.ruleComparerProcessor(R_CONF)
	}
}
func (r *Rules) Required() {
	if exists.In_array(R_REQ, r.Plain, true) {
		keys, _ := r.ruleFinderAsSlice(R_REQ)
		form := r.Validate.Form.(url.Values)
		for _, value := range keys {
			if !r.Validate.hasErrorKeyAlready(R_REQ, value) && len(form[value]) == 0 {
				r.Validate.addError(value, R_REQ, []string{value})
			}
		}
	}

	//if exists.In_array(R_REQ, r.Plain, true) {
	//	keys, _ := r.ruleFinderAsSlice(R_REQ)
	//	form := &url.Values{}
	//	inititator.Convert(r.Validate.Form, form)
	//	for _, value := range keys {
	//		if !r.Validate.hasErrorKeyAlready(R_REQ, value) && len(form.Get(value)) == 0 {
	//			r.Validate.addError(value, R_REQ, []string{value})
	//		}
	//	}
	//}
}
func (r *Rules) Equal() {
	if exists.In_array(R_EQ, r.Plain, false) {
		r.ruleComparerProcessor(R_EQ)
	}
}
func (r *Rules) Date() {
	if exists.In_array(R_DT, r.Plain, false) {
		r.ruleComparerProcessor(R_DT)
	}
}
func (r *Rules) Max() {
	if exists.In_array(R_MAX_LEN, r.Plain, false) {
		r.ruleComparerProcessor(R_MAX_LEN)
	}
}
func (r *Rules) Min() {
	if exists.In_array(R_MIN_LEN, r.Plain, false) {
		r.ruleComparerProcessor(R_MIN_LEN)
	}

}
func (r *Rules) Array() {
	if exists.In_array(R_ARR, r.Plain, true) {
		r.ruleComparerProcessor(R_ARR)
	}
}
func (r *Rules) File() {
	if exists.In_array(R_FILE, r.Plain, true) {
		el := reflect.ValueOf(r.Validate.Form).Elem()
		val := el.Field(1)
		iter := val.MapRange()
		if val.Len() == 0 {
			keys, _ := r.ruleFinderAsSlice(R_FILE)
			for _, value := range keys {
				r.Validate.addError(value, R_FILE, []string{value})
			}
		} else {
			for iter.Next() {
				fmt.Println(iter.Value().Len() == 0, iter.Value().Len())
				if iter.Value().Len() == 0 {
					r.Validate.addError(iter.Key().String(), R_FILE, []string{iter.Key().String()})
				}
			}
		}

		//for iter.Next() {
		//	fmt.Println(iter.Value().Len(), iter.Value())
		//	if iter.Value().Len() == 1 {
		//		val := iter.Value().Index(0).Elem()
		//		fSize := val.FieldByName("Size")
		//		fmt.Println()
		//	}
		//	//fmt.Println(iter.Value())
		//}
		//
		//for iter.Next() {
		//	fmt.Println(iter.Value())
		//	//if iter.Value().Kind() == reflect.Ptr {
		//	//
		//	//}
		//}
	}
}
func (r *Rules) FileTypes() {}
func (r *Rules) MaxSize()   {}
func (r *Rules) MinSize()   {}
func (r *Rules) Between() {
	if exists.In_array(R_BTW, r.Plain, false) {
		r.ruleComparerProcessor(R_BTW)
	}
}
func (r *Rules) String() {
	if exists.In_array(R_STR, r.Plain, true) {
		r.ruleComparerProcessor(R_STR)
	}
}
func (r *Rules) Integer() {
	if exists.In_array(R_INT, r.Plain, true) {
		r.ruleComparerProcessor(R_INT)
	}
}
func (r *Rules) Digit() {
	if exists.In_array(R_DG, r.Plain, true) {
		r.ruleComparerProcessor(R_DG)
	}
}
func (r *Rules) Email() {}

type Initializer struct {
	rules *Rules
}
type ITransform interface{ Convert(interface{}) }

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
func (i Initializer) Convert(d interface{}, c interface{}) {
	i.reSerialize(d, c)
}
func (i Initializer) reSerialize(v interface{}, con interface{}) {
	iter := reflect.ValueOf(v).MapRange()
	for iter.Next() {
		reflect.ValueOf(con).Set(iter.Value().Elem())
	}
	i.rules.Validate.Form = con
}

type Validate struct {
	Form            interface{}
	Errors          map[string][]string
	FieldNames      []string
	Internalization Internalization
}

func (v Validate) hasErrorAlready(errField string) (bool, []string) {
	val, ok := v.Errors[errField]
	return ok && len(val) > 0, val
}

func (v Validate) hasErrorKeyAlready(errKey, fieldKey string) bool {
	for key, value := range v.Errors {
		for _, v := range value {
			return key == fieldKey && strings.Contains(v, errKey)
		}
	}
	return false
}

func (v Validate) hasNamedErrorAlready(errField, errKey string) bool {
	val, ok := v.Errors[errField]
	if ok {
		for _, value := range val {
			if strings.Contains(value, errKey) {
				return true
			}
		}
	}

	return false
}

func (v *Validate) addError(errField string, errKey string, keys []string) {
	messages := v.Internalization.Messages()
	errMessage := messages[errKey]

	errMessage = errored(errMessage, keys)
	v.Errors[errField] = append(v.Errors[errField], errMessage)
}

func (v Validate) GetErrors() map[string][]string {
	return v.Errors
}

func (v Validate) GetError(key string) []string {
	_, val := v.hasErrorAlready(key)
	return val
}

func (v Validate) getFormValues(key string) []string {
	var found []string
	for k, value := range v.Form.(url.Values) {
		if k == key {
			found = value
		}
	}

	return found
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

//Slice
func sArrayContainsThenRetrieve(d []interface{}, c string) (bool, interface{}) {

	for _, value := range d {
		switch value.(type) {
		case []string:
			iter := value.([]string)
			for _, v := range iter {
				if strings.Contains(v, c) {
					return true, v
				}
			}
		}
	}

	return false, ""
}

// Map
func mGetValuesAsSlice(d interface{}) []interface{} {
	var values []interface{}
	switch d.(type) {
	case map[string]interface{}:
		iter := d.(map[string]interface{})
		for _, value := range iter {
			values = append(values, value)
		}
	case map[string][]string:
		iter := d.(map[string][]string)
		for _, value := range iter {
			values = append(values, value)
		}
	case map[string][]int:
		iter := d.(map[string][]int)
		for _, value := range iter {
			values = append(values, value)
		}
	}

	return values
}
func mSearchAndRetrieve(d interface{}, k string) interface{} {

	var fVal interface{}

	switch d.(type) {
	case map[interface{}]interface{}:
		iter := d.(map[string]interface{})
		for key, value := range iter {
			if key == k {
				fVal = value
			}
		}
	case map[interface{}][]interface{}:
		iter := d.(map[string][]interface{})
		for key, value := range iter {
			if key == k {
				fVal = value
			}
		}
	}

	return fVal
}
func t(d interface{}, k string) (interface{}, interface{}) {
	//var foundKeys interface{}
	//var foundValues interface{}

	foundKeys := reflect.New(reflect.TypeOf([]interface{}{})).Elem()
	foundValues := reflect.New(reflect.TypeOf([]interface{}{})).Elem()

	iter := reflect.ValueOf(d).MapRange()
	for iter.Next() {
		if iter.Key().String() == k {
			foundKeys = reflect.Append(foundKeys, iter.Key())
			if iter.Value().Len() == 1 {
				foundValues = reflect.Append(foundValues, iter.Value().Index(0))
			}
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
func errored(v string, k []string) string {
	var prepared = v
	// Let's image that we've a string like this:
	// :::> {0} is required.
	// we need to group tokens familiar with '{0}' this
	re := regexp.MustCompile(`({\d})`)

	mustReplace := re.FindAllString(v, -1)
	// grouped tokens can be replaced with
	// token values which are real value
	// for 'v'.
	mustReplaceWith := k
	mustReplaceWithLen := len(mustReplaceWith)

	for key, value := range mustReplace {
		// we've to check lengths of slices
		// because if both of them's lengths are not equal,
		// errors can be generated at runtime.
		if key >= 0 && mustReplaceWithLen > key {
			prepared = strings.Replace(prepared, value, mustReplaceWith[key], -1)
		}
	}

	return prepared
}
