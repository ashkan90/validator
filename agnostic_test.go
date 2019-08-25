package validator

import "testing"

func TestAgnostic(t *testing.T) {
	data := map[string][]string{
		"name":    {"Emirhanlaaaaaanasdsadsadqwdwq"},
		"surname": {"Ataman"},
	}

	rules := map[string][]string{
		"name": {"equal:test", "max:15"},
	}

	var b = &Builder{}
	b.Load(data, rules)
	//b.Required("name")
	//b.Equal("name")
	//b.Max("name")
	b.Min("name")
}
