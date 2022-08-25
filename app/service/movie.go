package service

import (
	"mime/multipart"
	"os"

	. "github.com/kadirgonen/movie-api/app/models"
	. "github.com/kadirgonen/movie-api/app/repo"
	"go.uber.org/zap"

	. "github.com/kadirgonen/movie-api/api/model"
	. "github.com/kadirgonen/movie-api/app/pkg/helper"
)

type MovieService struct {
	MovieRepo *MovieRepository
}

func NewMovieService(c *MovieRepository) *MovieService {
	return &MovieService{MovieRepo: c}
}

func (c *MovieService) Migrate() {
	c.MovieRepo.Migrate()
}

// Upload helps to read file and compare after that create on db
func (c *MovieService) Upload(file *multipart.File) (int, string, error) {
	var count int
	var str string

	movielist, err := ReadCSV(file)
	if err != nil {
		zap.L().Error("movie.service.readcsv", zap.Error(err))
		return 0, os.Getenv("ERROR"), err
	}
	moviesOnDb, err := c.GetAll()
	if err != nil {
		return 0, os.Getenv("ERROR"), err
	}

	if len(*moviesOnDb) > 0 {
		compared := CompareMovies(moviesOnDb, &movielist)
		if len(compared) > 0 {
			count, str, err = c.MovieRepo.Upload(&compared)
			if err != nil {
				return count, os.Getenv("ERROR"), err
			}
			return count, str, nil

		}
		return 0, os.Getenv("SAME_MOVİE"), nil
	}
	count, str, err = c.MovieRepo.Upload(&movielist)
	if err != nil {
		return count, os.Getenv("ERROR"), err
	}

	return count, str, nil

}

// GetAll returns all categories
func (c *MovieService) GetAll() (*MovieList, error) {
	return c.MovieRepo.GetAll()
}

func (c *MovieService) GetAllMoviesWithPagination(pag Pagination) (*Pagination, error) {
	movies, count, err := c.MovieRepo.GetAllMoviesWithPagination(pag)
	if err != nil {
		return nil, err
	}
	pag.Items = movies
	pag.TotalCount = int64(count)
	return &pag, nil
}
