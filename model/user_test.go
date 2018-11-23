package model

import "testing"

func TestUserIsValid(t *testing.T) {

	user1 := &User{
		Name: "Harry",
	}

	user2 := &User{
		Name: "1234",
	}

	user3 := &User{
		Name: "harry1234",
	}

	user4 := &User{
		Name: "&Harry",
	}

	testCases := []struct {
		name     string
		user    *User
		expected bool
	}{
		{"User name with only alphabates should be ture", user1, true},
		{"User name with only numeric should be true", user2, true},
		{"User name with only alphanumeric should be true", user3, true},
		{"User name with  special character should be flase", user4, false},
	}

	for _, test := range testCases {
		if test.user.IsValid() != test.expected {
			t.Error(test.name)
		}
	}
}
