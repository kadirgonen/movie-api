package repository

import (
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	. "github.com/kadirgonen/movie-api/app/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var email = "user@gmail.com"
var email1 = "test@gmail.com"
var user = User{
	FirstName: "DummyName",
	LastName:  "DummyLastName",
	Email:     email,
	Password:  "12345",
}

func NewMock() (DB *gorm.DB, mock sqlmock.Sqlmock) {
	var (
		db *sql.DB
	)
	db, mock, _ = sqlmock.New()
	DB, _ = gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	return DB, mock
}

func TestUserRepository_CheckUser(t *testing.T) {
	db, mock := NewMock()
	repo := &UserRepository{db}
	query := `SELECT * FROM "users" WHERE email=$1 AND "users"."deleted_at" IS NULL LIMIT 1`
	rows := sqlmock.NewRows([]string{"email", "password", "first_name", "last_name"}).AddRow(user.Email, user.Password, user.FirstName, user.LastName)

	t.Run("checkUser_True", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(email).WillReturnRows(rows)
		resp, err := repo.CheckUser(&user)
		assert.Equal(t, resp, true)
		require.NoError(t, err)
	})

	t.Run("checkUser_False", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(email1).WillReturnRows(rows)
		resp, err := repo.CheckUser(&user)
		assert.Equal(t, resp, false)
		require.NoError(t, err)
	})
}
