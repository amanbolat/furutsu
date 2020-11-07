package user

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
)

type User struct {
	Id        string    `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	FullName  string    `json:"full_name"`
	CreatedAt time.Time `sql:"created_at"`
	UpdatedAt time.Time `sql:"updated_at"`
}

func (u User) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Username, validation.Required, validation.Length(1, 32)),
		validation.Field(&u.Password, validation.Required, validation.Length(3, 32)),
		validation.Field(&u.FullName, validation.Required, validation.Length(1, 100)),
	)
}
