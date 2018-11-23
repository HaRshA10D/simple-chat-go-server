package model

import "testing"

func TestGroupMessageIsValid(t *testing.T) {

	groupMessage1 := &GroupMessage{
		UserID: 1234,
		GroupID: 2345,
		Message: "hello",
	}

	groupMessage2 := &GroupMessage{
		UserID: 1234,
		GroupID: 2345,
		Message: "",
	}

	groupMessage3 := &GroupMessage{
		UserID: 0,
		GroupID: 0,
		Message: "hello",
	}

	groupMessage4 := &GroupMessage{
		UserID: 0,
		GroupID: 0,
		Message: "",
	}

	testCases := []struct {
		name     string
		groupMessage    *GroupMessage
		expected bool
	}{
		{"Group message with not null userID, groupID and message should be true", groupMessage1, true},
		{"Group message with not null userID, groupID and null message should be false", groupMessage2, false},
		{"Group message with null userID, groupID and not null message should be false", groupMessage3, false},
		{"Group message with null userID, groupID and message should be false", groupMessage4, false},
	}

	for _, test := range testCases {
		if test.groupMessage.IsValid() != test.expected {
			t.Error(test.name)
		}
	}
}
