package repository

import (
	. "github.com/kadirgonen/movie-api/app/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

type UserRepositoryInterface interface {
	Save(user *User) (*User, error)
	Login(email string) (*User, error)
	CheckUser(user *User) (bool, error)
	Migrate()
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (u *UserRepository) Migrate() {
	u.db.AutoMigrate(&User{})
}

func (u *UserRepository) Save(user *User) (*User, error) {
	zap.L().Debug("user.repo.save", zap.Reflect("user", user))
	if err := u.db.Create(user).Error; err != nil {
		zap.L().Error("user.repo.Save failed to create user", zap.Error(err))
		return nil, err
	}
	return user, nil
}

func (u *UserRepository) Login(email string) (*User, error) {
	var user User
	zap.L().Debug("user.repo.login", zap.String("email", email))
	if err := u.db.Where("email = ?", email).Find(&user).Error; err != nil {
		zap.L().Error("user.repo.Login failed to login user", zap.Error(err))
		return nil, err
	}
	return &user, nil
}

func (u *UserRepository) CheckUser(user *User) (bool, error) {
	var us *User
	var exists bool = false
	zap.L().Debug("user.repo.checkuser", zap.Reflect("user", user))
	r := u.db.Where("email=?", user.Email).Limit(1).Find(&us)
	if r.RowsAffected > 0 {
		exists = true
	}
	return exists, nil
}
