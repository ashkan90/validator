package validator

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

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
			//r.Validate.Form.(url.Values)[mConfirmationKey][0]
			mConfirmationVal := mGetNthWithKey(r.Validate.Form, mConfirmationKey, 0).(string)

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
