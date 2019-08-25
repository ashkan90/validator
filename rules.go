package validator

import (
	"fmt"
	"github.com/ashkan90/golang-in_array"
	"reflect"
)

func (r *Rules) Confirmation() {
	if exists.In_array(R_CONF, r.Plain, false) {
		r.ruleComparerProcessor(R_CONF)
	}
}
func (r *Rules) Required() {
	if exists.In_array(R_REQ, r.Plain, true) {
		keys, _ := r.ruleFinderAsSlice(R_REQ)
		//form := r.Validate.Form.(url.Values)
		form := r.Validate.Form
		for _, value := range keys {
			if !r.Validate.hasErrorKeyAlready(R_REQ, value) && lenOnKey(form, value) == 0 {
				//if !r.Validate.hasErrorKeyAlready(R_REQ, value) && lenOnKey(form, value)len(form[value]) == 0 {
				r.Validate.addError(value, R_REQ, []string{value})
			}
		}
	}

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
