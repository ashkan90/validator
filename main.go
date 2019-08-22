package validator

import (
	"github.com/ashkan90/golang-in_array"
	"net/url"
	"regexp"
	"strconv"
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
func (r *Rules) ruleComparerProcessor(ruleKey string) {
	ok, k, v := r.ruleFinder(ruleKey)     // k = rule daki key değeri örn: 'name', 'surname'
	values := r.Validate.getFormValues(k) // k'daki değerler dizgesi.

	if strings.Contains(v, ":") {
		v = strings.SplitAfter(v, ":")[1] // "equal", "23", [1] == 23
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

	if c.ruleKey == R_EQ {
		if c.comp1 != c.comp2 {
			r.Validate.addError(c.ruleField, c.ruleKey, c.err)
		}
	} else if c.ruleKey == R_MAX_LEN {
		if isDigit(c.comp1.(string)) {
			comp1, _ := strconv.Atoi(c.comp1.(string)) // value of rule eg. '10', '255'
			comp2 := c.comp2.(string)                  // value of form input eg. 'emirhan'

			if comp1 < len(comp2) {
				r.Validate.addError(c.ruleField, c.ruleKey, c.err)
			}
		}

	} else if c.ruleKey == R_MIN_LEN {
		comp1, _ := strconv.Atoi(c.comp1.(string))
		comp2 := c.comp2.(string)
		if comp1 > len(comp2) {
			r.Validate.addError(c.ruleField, c.ruleKey, c.err)
		}
	} else if c.ruleKey == R_ARR {
		comp1 := 1
		comp2 := c.comp2.([]string)
		if len(comp2) <= comp1 {
			r.Validate.addError(c.ruleField, c.ruleKey, c.err)
		}
	} else if c.ruleKey == R_STR {
		v, err := strconv.Atoi(c.comp2.(string))
		if err != nil && v == 0 {
			r.Validate.addError(c.ruleField, c.ruleKey, c.err)
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

func (r *Rules) Confirmation() {}
func (r *Rules) Required() {

	//params := &subRulesParameters{
	//	rk: R_REQ,
	//	expr: "==",
	//	compare: 0,
	//	rules: r.Plain,
	//	validate: r.Validate,
	//	priority: true,
	//	deep: true,
	//}
	//
	//length(params)

	if exists.In_array(R_REQ, r.Plain, true) {
		form := r.Validate.Form.(url.Values)
		for _, value := range r.Validate.FieldNames {
			if !r.Validate.hasErrorKeyAlready(R_REQ, value) && len(form[value]) == 0 {
				r.Validate.addError(value, R_REQ, []string{value})
			}
		}
	}
}
func (r *Rules) Equal() {
	if exists.In_array(R_MAX_LEN, r.Plain, false) {
		r.ruleComparerProcessor(R_EQ)
	}
}
func (r *Rules) Date() {}
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
func (r *Rules) File()    {}
func (r *Rules) MaxSize() {}
func (r *Rules) MinSize() {}
func (r *Rules) Between() {}
func (r *Rules) String() {
	if exists.In_array(R_STR, r.Plain, true) {
		r.ruleComparerProcessor(R_STR)
	}
}
func (r *Rules) Integer() {}
func (r *Rules) Digit()   {}
func (r *Rules) Email()   {}

//
//var params = &subRulesParameters{}
//
//
//type subRulesParameters struct {
//	rk			string
//	expr		string
//	compare 	int
//	rules		map[string][]string
//	validate 	*Validate
//	priority 	bool
//	deep 		bool
//}
//
//func length(sp *subRulesParameters) {
//	fmt.Println(sp.rk, sp.rules)
//	if exists.In_array(sp.rk, sp.rules, sp.deep) {
//		form := sp.validate.Form.(url.Values)
//		if sp.deep {
//			for _, value := range sp.validate.FieldNames {
//				switch sp.expr {
//				case "==":
//					if !sp.validate.hasErrorKeyAlready(sp.rk, value) && len(form[value]) == sp.compare {
//						sp.validate.addError(value, sp.rk, []string{value})
//					}
//				case "!=":
//					if !sp.validate.hasErrorKeyAlready(sp.rk, value) && len(form[value]) != sp.compare {
//						sp.validate.addError(value, sp.rk, []string{value})
//					}
//				case ">=":
//					if !sp.validate.hasErrorKeyAlready(sp.rk, value) && len(form[value]) >= sp.compare {
//						sp.validate.addError(value, sp.rk, []string{value})
//					}
//				case "<=":
//					if !sp.validate.hasErrorKeyAlready(sp.rk, value) && len(form[value]) <= sp.compare {
//						sp.validate.addError(value, sp.rk, []string{value})
//					}
//				case ">":
//					fmt.Println(len(form[value]))
//					if !sp.validate.hasErrorKeyAlready(sp.rk, value) && len(form[value]) > sp.compare {
//						sp.validate.addError(value, sp.rk, []string{value})
//					}
//				case "<":
//					if !sp.validate.hasErrorKeyAlready(sp.rk, value) && len(form[value]) < sp.compare {
//						sp.validate.addError(value, sp.rk, []string{value})
//					}
//				}
//
//			}
//		} else {
//			for _, value := range sp.validate.FieldNames {
//				switch sp.expr {
//				case "==":
//					if len(form[value]) == sp.compare {
//						sp.validate.addError(value, sp.rk, []string{value})
//					}
//				case "!=":
//					if len(form[value]) != sp.compare {
//						sp.validate.addError(value, sp.rk, []string{value})
//					}
//				case ">=":
//					if len(form[value]) >= sp.compare {
//						sp.validate.addError(value, sp.rk, []string{value})
//					}
//				case "<=":
//					if len(form[value]) <= sp.compare {
//						sp.validate.addError(value, sp.rk, []string{value})
//					}
//				case ">":
//					if len(form[value]) > sp.compare {
//						sp.validate.addError(value, sp.rk, []string{value})
//					}
//				case "<":
//					if len(form[value]) < sp.compare {
//						sp.validate.addError(value, sp.rk, []string{value})
//					}
//				}
//
//			}
//		}
//	}
//}

type Initializer struct {
	rules *Rules
}

func (i Initializer) Run() map[string][]string {
	i.rules.Confirmation()
	i.rules.Required() // done
	i.rules.Equal()    // done
	i.rules.Date()
	i.rules.Max()    // done
	i.rules.Min()    // done
	i.rules.Array()  // done
	i.rules.String() // done
	//i.rules.File() //

	return i.rules.Validate.Errors
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

type IMessage interface{ Messages() }

func (i Internalization) Messages() map[string]string {
	return map[string]string{
		R_CONF:    "",
		R_DT:      "",
		R_REQ:     "{0} is required.",
		R_MAX_LEN: "{0} field's length cannot be greater than '{1}'",
		R_MIN_LEN: "{0} field's length cannot be lower than '{1}'",
		R_EQ:      "{0} field must be '{1}'",
		R_ARR:     "{0} field must be type of Array/Slice",
		R_STR:     "{0} field must be type of String",
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

func isDigit(s string) bool {
	var holding bool
	for _, value := range s {
		if value >= 0 || value <= 10 {
			holding = true
		}
	}

	return holding
}

func errored(v string, k []string) string {
	var prepared string
	prepared = v
	re := regexp.MustCompile(`({\d})`)

	mustReplace := re.FindAllString(v, -1)
	mustReplaceWith := k
	mustReplaceWithLen := len(mustReplaceWith)

	for key, value := range mustReplace {
		if key >= 0 && mustReplaceWithLen > key {
			prepared = strings.Replace(prepared, value, mustReplaceWith[key], -1)
		}
	}

	return prepared
}
