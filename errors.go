package validator

import (
	"regexp"
	"strings"
)

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
