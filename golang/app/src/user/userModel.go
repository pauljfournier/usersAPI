package user

import (
	"errors"
	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"net/url"
	"strings"
	"time"
)

func ErrRequiredValue(field string) error {
	return errors.New(field + " value required")
}

// User represents the schema for the User
type User struct {
	ID        string    `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName string    `bson:"first_name" json:"first_name,omitempty"`
	LastName  string    `bson:"last_name" json:"last_name,omitempty"`
	Nickname  string    `bson:"nickname" json:"nickname,omitempty"`
	Password  string    `bson:"password" json:"password,omitempty"`
	Email     string    `bson:"email" json:"email,omitempty"`
	Country   string    `bson:"country" json:"country,omitempty"`
	CreatedAt time.Time `bson:"created_at" json:"created_at,omitempty"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at,omitempty"`
}

//Escape User for safety
func (u *User) Escape() {
	u.FirstName = strings.ReplaceAll(url.QueryEscape(u.FirstName), "+", "%20")
	u.LastName = strings.ReplaceAll(url.QueryEscape(u.LastName), "+", "%20")
	u.Nickname = strings.ReplaceAll(url.QueryEscape(u.Nickname), "+", "%20")
	u.Password = strings.ReplaceAll(url.QueryEscape(u.Password), "+", "%20")
	u.Email = strings.ReplaceAll(strings.ReplaceAll(url.QueryEscape(u.Email), "+", "%20"), "%40", "@")
	u.Country = strings.ReplaceAll(url.QueryEscape(u.Country), "+", "%20")
}

// Validate validates User email and returns validation errors.
func (u *User) Validate() error {
	if u.FirstName == "" {
		return ErrRequiredValue("first_name")
	}
	if u.LastName == "" {
		return ErrRequiredValue("last_name")
	}
	if u.Nickname == "" {
		return ErrRequiredValue("nickname")
	}
	if u.Password == "" {
		return ErrRequiredValue("password")
	}
	if u.Email == "" {
		return ErrRequiredValue("email")
	}
	if u.Country == "" {
		return ErrRequiredValue("country")
	}

	// Verify that the name value is valid
	return validation.ValidateStruct(u,
		validation.Field(&u.Email, validation.Required, is.Email),
	)
}

func (u *User) IsSoftEqual(u2 *User) bool {
	return u.FirstName == u2.FirstName &&
		u.LastName == u2.LastName &&
		u.Nickname == u2.Nickname &&
		u.Password == u2.Password &&
		u.Email == u2.Email &&
		u.Country == u2.Country
}
