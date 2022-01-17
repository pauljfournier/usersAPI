package user

import (
	"testing"
)

type userEscapeTest struct {
	user                User
	userEscapedExpected User
}

func TestEscape(t *testing.T) {
	userModelTests := []userEscapeTest{
		//test without special characters
		{
			user: User{
				FirstName: "FirstName",
				LastName:  "LastName",
				Nickname:  "Nickname",
				Password:  "Password",
				Email:     "Email@email.com",
				Country:   "Country",
			},
			userEscapedExpected: User{
				FirstName: "FirstName",
				LastName:  "LastName",
				Nickname:  "Nickname",
				Password:  "Password",
				Email:     "Email@email.com",
				Country:   "Country",
			},
		},
		//test with special characters
		{
			user: User{
				FirstName: "FirstName£",
				LastName:  "Last§Name",
				Nickname:  "]Nickname",
				Password:  "Pass word",
				Email:     "Eµmail@email.com",
				Country:   "Country@",
			},
			userEscapedExpected: User{
				FirstName: "FirstName%C2%A3",
				LastName:  "Last%C2%A7Name",
				Nickname:  "%5DNickname",
				Password:  "Pass%20word",
				Email:     "E%C2%B5mail@email.com",
				Country:   "Country%40",
			},
		},
	}

	for _, item := range userModelTests {
		item.user.Escape()
		if !item.user.IsSoftEqual(&item.userEscapedExpected) {
			t.Errorf("User.Escape output %v but expected %v", item.user, item.userEscapedExpected)
		}
	}
}

type userValidatorTest struct {
	user        User
	expectedErr bool
}

func TestValidator(t *testing.T) {
	userValidatorTests := []userValidatorTest{
		//test normal behavior
		{
			user: User{
				FirstName: "FirstName",
				LastName:  "LastName",
				Nickname:  "Nickname",
				Password:  "Password",
				Email:     "Email@email.com",
				Country:   "Country",
			},
			expectedErr: false,
		},
		//test one missing
		{
			user: User{
				LastName: "LastName",
				Nickname: "Nickname",
				Password: "Password",
				Email:    "Email@email.com",
				Country:  "Country",
			},
			expectedErr: true,
		},
		{
			user: User{
				FirstName: "FirstName",
				Nickname:  "Nickname",
				Password:  "Password",
				Email:     "Email@email.com",
				Country:   "Country",
			},
			expectedErr: true,
		},
		{
			user: User{
				FirstName: "FirstName",
				LastName:  "LastName",
				Password:  "Password",
				Email:     "Email@email.com",
				Country:   "Country",
			},
			expectedErr: true,
		},
		{
			user: User{
				FirstName: "FirstName",
				LastName:  "LastName",
				Nickname:  "Nickname",
				Email:     "Email@email.com",
				Country:   "Country",
			},
			expectedErr: true,
		},
		{
			user: User{
				FirstName: "FirstName",
				LastName:  "LastName",
				Nickname:  "Nickname",
				Password:  "Password",
				Country:   "Country",
			},
			expectedErr: true,
		},
		{
			user: User{
				FirstName: "FirstName",
				LastName:  "LastName",
				Nickname:  "Nickname",
				Password:  "Password",
				Email:     "Email@email.com",
			},
			expectedErr: true,
		},
		//test multiple missing
		{
			user: User{
				FirstName: "FirstName",
				Email:     "Email@email.com",
				Country:   "Country",
			},
			expectedErr: true,
		},
		//test invalid email
		{
			user: User{
				FirstName: "FirstName",
				LastName:  "LastName",
				Nickname:  "Nickname",
				Password:  "Password",
				Email:     "emailemailcom",
				Country:   "Country",
			},
			expectedErr: true,
		},
	}

	for _, item := range userValidatorTests {
		resultErr := item.user.Validate()
		if !item.expectedErr && resultErr != nil {
			t.Errorf("user.Validate for %v output err %v not expected", item.user, resultErr.Error())
		}
		if item.expectedErr && resultErr == nil {
			t.Errorf("user.Validate for %v output err expected but not found", item.user)
		}
	}
}
