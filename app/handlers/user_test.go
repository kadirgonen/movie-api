package handler

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	. "github.com/kadirgonen/movie-api/app/models"
	"github.com/kadirgonen/movie-api/app/pkg/config"
	"github.com/stretchr/testify/assert"
)

var (
	Email     = "deneme@gmail.com"
	Firstname = "DenemeIsim"
	Lastname  = "DenemeSoyIsim"
	Password  = "12345"
)

var TruePayLoad = []byte(
	`{"email":"test@gmail.com","firstname":"dummyname","lastname":"dummylastname","password":"12345"}`,
)

var BadPayLoad = []byte(
	`{"e_mail":"test@gmail.com","first_name":"dummyname","lastname":"dummylastname","password":"12345"}`,
)

func TestUserHandler_signup(t *testing.T) {
	cfg := &config.Config{
		JWTConfig: config.JWTConfig{
			AccessSessionTime:  300,
			RefreshSessionTime: 900,
			SecretKey:          "secretkey",
		},
	}
	t.Run("userSignUp_UserExist", func(t *testing.T) {

		mockUserService := &mockUserService{
			items: []User{
				{Email: "test@gmail.com", FirstName: "dummyname", LastName: "dummylastname", Password: "12345"},
			},
			cfg: cfg,
		}

		handler := &UserHandler{userService: mockUserService, cfg: cfg}
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/signup", nil)
		c.Request.Header.Set("content-Type", "application/json")
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(TruePayLoad))
		handler.signup(c)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("userSignUp_UserSuccess", func(t *testing.T) {

		mockUserService := &mockUserService{
			items: []User{},
			cfg:   cfg,
		}

		handler := &UserHandler{userService: mockUserService, cfg: cfg}
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/signup", nil)
		c.Request.Header.Set("content-Type", "application/json")
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(TruePayLoad))
		handler.signup(c)
		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("userSignUp_BadBody", func(t *testing.T) {

		mockUserService := &mockUserService{
			items: []User{},
			cfg:   cfg,
		}

		handler := &UserHandler{userService: mockUserService, cfg: cfg}
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/signup", nil)
		c.Request.Header.Set("content-Type", "application/json")
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(BadPayLoad))
		handler.signup(c)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

// Duplicate Error Code
var UserExist = fmt.Errorf("23505")

type mockUserService struct {
	items []User
	cfg   *config.Config
}

func (m *mockUserService) Save(user *User) (*User, error) {
	for _, item := range m.items {
		if *&item.Email == *&user.Email {
			return nil, UserExist
		}
	}
	m.items = append(m.items, *user)
	return user, nil
}

func (m *mockUserService) Login(email string, password string) (*User, error) {
	return nil, nil
}

func (m *mockUserService) CheckUser(user *User) (bool, error) {
	return false, nil
}

func (u *mockUserService) Migrate() {

}
