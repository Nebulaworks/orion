package applicant

import "testing"

func TestIsValidEmail(t *testing.T) {
	cases := map[string]bool{
		"valid@example.com":          true,
		"missingdomain":              false,
		"double@at@symbols.com":      false,
		"notld@booze":                true,
		"contains,commas@commas.com": false,
	}
	for input, expected := range cases {
		got := isValidEmail(input)
		if expected != got {
			t.Logf("error: %v should be %v but got %v", input, expected, got)
			t.Fail()
		}
	}
}

func TestIsValidRole(t *testing.T) {
	cases := map[string]bool{
		"invalid!chars,andsuch":    false,
		"Senior Software Engineer": true,
	}
	for input, expected := range cases {
		got := isValidRole(input)
		if expected != got {
			t.Logf("error: %v should be %v but got %v", input, expected, got)
			t.Fail()
		}
	}
}
