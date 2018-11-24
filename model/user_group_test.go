package model

import "testing"

func TestUserGroupIsValid(t *testing.T) {

	userGroup1 := &UserGroup{
		UserID:  1234,
		GroupID: 4567,
	}

	userGroup2 := &UserGroup{
		UserID:  0,
		GroupID: 5678,
	}

	userGroup3 := &UserGroup{
		UserID:  0,
		GroupID: 0,
	}

	testCases := []struct {
		name      string
		userGroup *UserGroup
		expected  bool
	}{
		{"User group with not null userID and groupID should be true", userGroup1, true},
		{"User group with null userID and not null groupID should be false", userGroup2, false},
		{"User group with null userID and groupID should be false", userGroup3, false},
	}

	for _, test := range testCases {
		if test.userGroup.IsValid() != test.expected {
			t.Error(test.name)
		}
	}
}
