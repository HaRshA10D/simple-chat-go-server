package model

import "testing"

func TestGroupIsValid(t *testing.T) {

	group1 := &Group{
		Name: "Abhinav",
	}

	group2 := &Group{
		Name: "",
	}

	group3 := &Group{
		Name: "Harry//",
	}

	testCases := []struct {
		name     string
		group    *Group
		expected bool
	}{
		{"Valid group name should return true", group1, true},
		{"Group name with empty name should return false", group2, false},
		{"Group name with special charecters other than _ should return false", group3, false},
	}

	for _, test := range testCases {
		if test.group.IsValid() != test.expected {
			t.Error(test.name)
		}
	}
}
