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

func NewCategoryRepository(db *gorm.DB) *MovieRepository {
	return &MovieRepository{db: db}
}

func (c *MovieRepository) Migrate() {
	c.db.AutoMigrate(&Movie{})
}
func (c *MovieRepository) Create(mv *Movie) (*Movie, error) {
	zap.L().Debug("movie.repo.create", zap.Reflect("movie", mv.Name))
	if err := c.db.Create(&mv).Error; err != nil {
		zap.L().Error("movie.repo.create failed to create movie", zap.Error(err))
		return nil, err
	}
	return mv, nil
}

// Update take movie id. It can be uses without id with Save. But it creates different movie id's. Here it can be uses with Save() without Where() conditions,
// because id is implemented movie. Here I would like to differentiate usage practices.
func (c *MovieRepository) Update(mv *Movie, id int) (*Movie, error) {
	zap.L().Debug("movie.repo.update", zap.Reflect("movie", mv))
	mv.ID = id
	if err := c.db.Where("id=?", id).Updates(&mv).Error; err != nil {
		zap.L().Error("movie.repo.update failed to update movie", zap.Error(err))
		return nil, err
	}

	return mv, nil

}
func (c *MovieRepository) Upload(movies *MovieList) (int, string, error) {
	var count int64
	err := c.db.Create(&movies).Count(&count).Error
	return int(count), os.Getenv("CREATE_FILE"), err
}

func (c *MovieRepository) GetAll() (*MovieList, error) {
	var movies MovieList
	err := c.db.Find(&movies).Error
	if err != nil {
		return nil, err
	}
	return &movies, nil
}

func (c *MovieRepository) GetAllMoviesWithPagination(pag Pagination) ([]Movie, int, error) {
	var movies []Movie
	var count int64
	err := c.db.Offset(int(pag.Page) - 1*int(pag.PageSize)).Limit(int(pag.PageSize)).Find(&movies).Count(&count).Error
	if err != nil {
		return nil, -1, err
	}
	return movies, int(count), nil
}
