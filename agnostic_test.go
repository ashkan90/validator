package validator

import "testing"

func TestAgnostic(t *testing.T) {
	data := map[string][]string{
		"name":    {"Emirhan"},
		"surname": {"Ataman"},
	}

	rules := map[string][]string{
		"name": {"confirmation:test"},
	}

	var b = &Builder{}
	b.Load(data, rules)
	b.Required("name")
	b.Confirmation("name")
}
