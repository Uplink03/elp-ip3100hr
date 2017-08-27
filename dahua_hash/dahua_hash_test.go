package dahua_hash

import "testing"

func TestDahuaHash(t *testing.T) {
	for _, test := range testCases {
		hashed := DahuaHash(test.original)
		if hashed != test.hashed {
			t.Errorf("Input: %s. Expected: %s. Got: %s", test.original, test.hashed, hashed)
		}
	}
}

type testCase struct {
	original, hashed string
}

var testCases = []testCase{
	testCase{"", "tlJwpbo6"},
	testCase{"admin", "6QNMIQGe"},
	testCase{"password", "mF95aD4o"},
	testCase{"qwerty12345", "FJHwdAA1"},
}
