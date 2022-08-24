package model

import (
	"github.com/google/uuid"
	. "github.com/kadirgonen/movie-api/api/model"
)

func ResponseToUser(u *SignUp) *User {
	return &User{
		Id:        uuid.New().String(),
		FirstName: *u.Firstname,
		LastName:  *u.Lastname,
		Password:  *u.Password,
		Email:     *u.Email,
		IsAdmin:   false,
	}
}
