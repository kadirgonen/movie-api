package repository

import (
	"os"

	. "github.com/kadirgonen/movie-api/api/model"
	. "github.com/kadirgonen/movie-api/app/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type MovieRepository struct {
	db *gorm.DB
}

func NewMovieRepository(db *gorm.DB) *MovieRepository {
	return &MovieRepository{db: db}
}

func (m *MovieRepository) Migrate() {
	m.db.AutoMigrate(&Movie{})
}
func (m *MovieRepository) Create(mv *Movie) (*Movie, error) {
	zap.L().Debug("movie.repo.create", zap.Reflect("movie", mv.Name))
	if err := m.db.Create(&mv).Error; err != nil {
		zap.L().Error("movie.repo.create failed to create movie", zap.Error(err))
		return nil, err
	}
	return mv, nil
}

// Update take movie id. It can be uses without id with Save. But it creates different movie id's. Here it can be uses with Save() without Where() conditions,
// because id is implemented movie. Here I would like to differentiate usage practices.
func (m *MovieRepository) Update(mv *Movie, id int) (*Movie, error) {
	zap.L().Debug("movie.repo.update", zap.Reflect("movie", mv))
	mv.ID = id
	if err := m.db.Where("id=?", id).Updates(&mv).Error; err != nil {
		zap.L().Error("movie.repo.update failed to update movie", zap.Error(err))
		return nil, err
	}

	return mv, nil

}

// Delete helps user to delete movie. Here is hard-delete. If you would like to change it soft-delete, just remove "Unscoped()".
func (m *MovieRepository) Delete(id int) (bool, error) {
	zap.L().Debug("movie.repo.delete", zap.Reflect("query", id))
	if err := m.db.Unscoped().Where("id=?", id).Delete(&Movie{}).Error; err != nil {
		zap.L().Error("movie.repo.delete failed to delete movie", zap.Error(err))
		return false, err
	}
	return true, nil
}

// CheckMovie helps user to check product is exist or not
func (m *MovieRepository) CheckMovie(id int) (bool, error) {
	var mv *Movie
	var exists bool = false
	zap.L().Debug("movie.repo.checkproduct")
	r := m.db.Where("id=?", id).Limit(1).Find(&mv)
	if r.RowsAffected > 0 {
		exists = true
	}
	return exists, nil
}
func (m *MovieRepository) Upload(movies *MovieList) (int, string, error) {
	var count int64
	err := m.db.Create(&movies).Count(&count).Error
	return int(count), os.Getenv("CREATE_FILE"), err
}

func (m *MovieRepository) GetAll() (*MovieList, error) {
	var movies MovieList
	err := m.db.Find(&movies).Error
	if err != nil {
		return nil, err
	}
	return &movies, nil
}

func (m *MovieRepository) GetAllMoviesWithPagination(pag Pagination) ([]Movie, int, error) {
	var movies []Movie
	var count int64
	err := m.db.Offset(int(pag.Page) - 1*int(pag.PageSize)).Limit(int(pag.PageSize)).Find(&movies).Count(&count).Error
	if err != nil {
		return nil, -1, err
	}
	return movies, int(count), nil
}
