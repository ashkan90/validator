package validator

import (
	exists "github.com/ashkan90/golang-in_array"
	"net/url"
	"strings"
)

const (
	R_REQ 	= "required"
	R_MAX 	= "max"
	R_MIN 	= "min"
	R_EQ  	= "equal"
	R_CONF 	= "confirmation"
	R_DT	= "date"

)

type Rules struct {
	Plain map[string][]string
	Validate *Validate
}

func (r Rules) ruleFinder(ruleKey string) (bool, string, string) {
	var foundKey string
	var foundValue string
	for key, value := range r.Plain {
		for _, v := range value {
			if strings.Contains(v, ruleKey) {
				foundValue = v
				foundKey = key

				return true, foundKey, foundValue
			}
		}
	}

	return false, foundKey, foundValue
}
func (r Rules) ruleComparerProcessor(ruleKey string, ruleError []string) {
	ok, k, v := r.ruleFinder(ruleKey)
	values := r.Validate.getFormValues(k)// k'daki değerler dizgesi.

	v = strings.SplitAfter(v, ":")[1] // "equal", "23", [1] == 23
	if ok {
		for _, value := range values {
			if v != value {
				r.Validate.addError(k, ruleKey, ruleError)
			}
		}
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

}

func (r Rules) Confirmation() {  }
func (r Rules) Required() {
	if exists.In_array(R_REQ, r.Plain) {
		form := r.Validate.Form.(url.Values)
		for _, value := range r.Validate.FieldNames {
			if !r.Validate.hasErrorKeyAlready(R_REQ) && len(form[value]) == 0 {
				r.Validate.addError(value, R_REQ, []string{})
			}
		}
	}
}
func (r Rules) Equal() {
	r.ruleComparerProcessor(R_EQ, []string{})
}
func (r Rules) Date() { }
func (r Rules) Max() {
	r.ruleComparerProcessor(R_MAX, []string{})
}
func (r Rules) Min() {
	r.ruleComparerProcessor(R_MIN, []string{})
}

type Initializer struct {
	rules *Rules
}

func (i Initializer) Run() map[string][]string {
	i.rules.Confirmation()
	i.rules.Required()
	i.rules.Equal()
	i.rules.Date()
	i.rules.Max()
	i.rules.Min()

	return i.rules.Validate.Errors
}


type Validate struct {
	Form interface{}
	Errors map[string][]string
	FieldNames []string
	Internalization Internalization
}

func (v Validate) hasErrorAlready(errField string) bool {
	val, ok := v.Errors[errField]
	return ok && len(val) > 0
}

func (v Validate) hasErrorKeyAlready(errKey string) bool {
	for _, value := range v.Errors {
		for _, v := range value {
			return strings.Contains(v, errKey)
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
	messages 	:= v.Internalization.Messages()
	errMessage 	:= messages[errKey]
	v.Errors[errField] = append(v.Errors[errField], errMessage)
}

func (v Validate) GetErrors() map[string][]string {
	return v.Errors
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


type Internalization struct {}

type IMessage interface { Messages() }

func (i Internalization) Messages() map[string]string {
	return map[string]string{
		R_CONF	: "",
		R_DT	: "",
		R_REQ	: "{0} is required.",
		R_MAX	: "{0} field cannot be greater than '{1}'",
		R_MIN	: "{0} field cannot be lower than '{1}'",
		R_EQ	: "{0} field must be '{1}'",
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

