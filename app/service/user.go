package service

import (
	"net/http"
	"os"

	. "github.com/kadirgonen/movie-api/app/models"
	. "github.com/kadirgonen/movie-api/app/pkg/errors"
	. "github.com/kadirgonen/movie-api/app/pkg/helper"
	. "github.com/kadirgonen/movie-api/app/repo"
)

type UserService struct {
	UserRepo UserRepositoryInterface
}

type UserServiceInterface interface {
	CheckUser(user *User) (bool, error)
	Save(user *User) (*User, error)
	Login(email string, password string) (*User, error)
	Migrate()
}

func NewUserService(u UserRepositoryInterface) *UserService {
	return &UserService{UserRepo: u}
}

// CheckUser helps to identify user is registered before or not
func (u *UserService) CheckUser(user *User) (bool, error) {
	return u.UserRepo.CheckUser(user)

}

// Save helps to create user after check verify e-mail,password
func (u *UserService) Save(user *User) (*User, error) {
	if VerifyEMail(user.Email) {
		ver := VerifyPassword(user.Password)
		if ver == true {
			hash, err := HashPassword(user.Password)
			if err != nil {
				return nil, NewRestError(http.StatusBadRequest, os.Getenv("ISSUE_PASSWORD"), nil)
			}
			user.Password = hash
			return u.UserRepo.Save(user)
		}
		return nil, NewRestError(http.StatusBadRequest, os.Getenv("WEAK_PASSWORD"), nil)
	}
	return nil, NewRestError(http.StatusBadRequest, os.Getenv("INVALID_MAIL"), nil)
}

// Login helps user to enter system after checking e-mail,password
func (u *UserService) Login(email string, password string) (*User, error) {
	user, err := u.UserRepo.Login(email)
	if err != nil {
		return nil, NewRestError(http.StatusBadRequest, os.Getenv("ISSUE_DECODE_PASSWORD"), nil)
	}
	res := CheckPasswordHash(password, user.Password)
	if !res {
		return nil, NewRestError(http.StatusBadRequest, os.Getenv("WRONG_EMAIL_PASSWORD"), nil)
	}
	return user, nil
}

func (u *UserService) Migrate() {
	u.UserRepo.Migrate()
}
