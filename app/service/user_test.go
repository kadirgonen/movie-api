package service_test

import (
	"errors"
	"testing"

	. "github.com/kadirgonen/movie-api/app/models"
	. "github.com/kadirgonen/movie-api/app/pkg/helper"
	. "github.com/kadirgonen/movie-api/app/repo"
	. "github.com/kadirgonen/movie-api/app/service"
)

func TestUserService_CheckUser(t *testing.T) {
	type fields struct {
		userRepo UserRepositoryInterface
	}
	type args struct {
		user *User
	}
	mockRepo := &mockUserRepository{
		items: []User{
			{Id: "1", FirstName: "İsim", LastName: "Soyisim", Email: "registeredbefore@gmail.com"},
		},
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "userCheck_Success",
			fields: fields{
				userRepo: mockRepo,
			},
			args: args{
				user: &User{
					Id:    "1",
					Email: "registeredbefore@gmail.com",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &UserService{
				UserRepo: tt.fields.userRepo,
			}
			_, err := a.CheckUser(tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthService.Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUserService_Save(t *testing.T) {
	type fields struct {
		userRepo UserRepositoryInterface
	}
	type args struct {
		user *User
	}

	mockRepo := &mockUserRepository{
		items: []User{
			{Id: "1", FirstName: "İsim", LastName: "Soyisim", Email: "registeredbefore@gmail.com"},
		},
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *User
		wantErr bool
	}{
		{
			name: "userSave_Success",
			fields: fields{
				userRepo: mockRepo,
			},
			args: args{
				user: &User{
					Id:        "2",
					FirstName: "İsim1",
					LastName:  "Soyisim1",
					Email:     "unregisteredbefore@gmail.com",
					Password:  "aAbcdef12.",
					IsAdmin:   false,
				},
			},
			want: &User{
				Id:        "2",
				FirstName: "İsim1",
				LastName:  "Soyisim1",
				Email:     "unregisteredbefore@gmail.com",
				Password:  "$2a$14$y6G6s3wIMnIU08rKPNA2s.oFn39T9fGQ7kaTuTF5O2Qb4L5JghfrS",
				IsAdmin:   false,
			},
			wantErr: false,
		},
		{
			name: "userSave_Fail",
			fields: fields{
				userRepo: mockRepo,
			},
			args: args{
				user: &User{
					Id:        "1",
					FirstName: "İsim1",
					LastName:  "Soyisim1",
					Email:     "registeredbefore@gmail.com",
					Password:  "aAbcdef12.",
					IsAdmin:   false,
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &UserService{
				UserRepo: tt.fields.userRepo,
			}
			_, err := a.Save(tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserService.Save() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func TestUserService_Login(t *testing.T) {
	type fields struct {
		UserRepo UserRepositoryInterface
	}
	type args struct {
		email    string
		password string
	}
	pass, _ := HashPassword("12345")

	mockRepo := &mockUserRepository{
		items: []User{
			{Id: "1", FirstName: "İsim", LastName: "Soyisim", Email: "registeredbefore@gmail.com", Password: pass},
		},
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *User
		wantErr bool
	}{
		{name: "UserLoginSuccess", fields: fields{UserRepo: mockRepo}, args: args{email: "registeredbefore@gmail.com", password: "12345"}, wantErr: false},
		{name: "UserLoginFail", fields: fields{UserRepo: mockRepo}, args: args{email: "unregisteredbefore@gmail.com", password: "12345"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserService{
				UserRepo: tt.fields.UserRepo,
			}
			_, err := u.Login(tt.args.email, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserService.Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// if !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("UserService.Login() = %v, want %v", got, tt.want)
			// }
		})
	}
}

type mockUserRepository struct {
	items []User
}

var (
	errCRUD      = errors.New("Mock: Error crud operation")
	userNotFound = errors.New("User not found")
)

func (u *mockUserRepository) CheckUser(user *User) (bool, error) {
	for _, v := range u.items {
		if user.Email == v.Email {
			return false, errors.New("User already exist!")
		}
	}
	return true, nil

}

func (u *mockUserRepository) Save(user *User) (*User, error) {
	for _, item := range u.items {
		if item.FirstName == user.FirstName || item.Email == user.Email {
			return nil, errCRUD
		}
	}
	u.items = append(u.items, *user)
	return user, nil
}

func (u *mockUserRepository) Login(email string) (*User, error) {
	for _, item := range u.items {
		if item.Email == email {
			return &User{Id: "1", Password: item.Password}, nil
		}
	}
	return nil, errors.New("User not exist!")
}

func (u *mockUserRepository) Migrate() {

}
