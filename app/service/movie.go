package service

import (
	"mime/multipart"
	"net/http"
	"os"

	. "github.com/kadirgonen/movie-api/app/models"
	. "github.com/kadirgonen/movie-api/app/repo"
	"go.uber.org/zap"

	. "github.com/kadirgonen/movie-api/api/model"
	. "github.com/kadirgonen/movie-api/app/pkg/errors"
	. "github.com/kadirgonen/movie-api/app/pkg/helper"
)

type MovieService struct {
	MovieRepo *MovieRepository
}

func NewMovieService(c *MovieRepository) *MovieService {
	return &MovieService{MovieRepo: c}
}

func (m *MovieService) Migrate() {
	m.MovieRepo.Migrate()
}
func (m *MovieService) Create(mv *Movie) (*Movie, error) {
	return m.MovieRepo.Create(mv)
}
func (m *MovieService) Update(mv *Movie, id int) (*Movie, error) {
	res, err := m.MovieRepo.CheckMovie(id)
	if err != nil {
		return nil, NewRestError(http.StatusBadRequest, os.Getenv("UPDATE_CHECK_MOVIE_ISSUE"), nil)
	}
	if res {
		return m.MovieRepo.Update(mv, id)
	} else {
		return nil, NewRestError(http.StatusBadRequest, os.Getenv("NO_MOVIE"), nil)
	}

}

func (m *MovieService) Delete(id int) (bool, error) {
	res, err := m.MovieRepo.CheckMovie(id)
	if err != nil {
		return false, NewRestError(http.StatusBadRequest, os.Getenv("DELETE_CHECK_PRODUCT_ISSUE"), nil)
	}
	if res {
		return m.MovieRepo.Delete(id)
	} else {
		return false, NewRestError(http.StatusBadRequest, os.Getenv("NO_PRODUCT"), nil)
	}

}

// Upload helps to read file and compare after that create on db
func (m *MovieService) Upload(file *multipart.File) (int, string, error) {
	var count int
	var str string

	movielist, err := ReadCSV(file)
	if err != nil {
		zap.L().Error("movie.service.readcsv", zap.Error(err))
		return 0, os.Getenv("ERROR"), err
	}
	moviesOnDb, err := m.GetAll()
	if err != nil {
		return 0, os.Getenv("ERROR"), err
	}

	if len(*moviesOnDb) > 0 {
		compared := CompareMovies(moviesOnDb, &movielist)
		if len(compared) > 0 {
			count, str, err = m.MovieRepo.Upload(&compared)
			if err != nil {
				return count, os.Getenv("ERROR"), err
			}
			return count, str, nil

		}
		return 0, os.Getenv("SAME_MOVİE"), nil
	}
	count, str, err = m.MovieRepo.Upload(&movielist)
	if err != nil {
		return count, os.Getenv("ERROR"), err
	}

	return count, str, nil

}

// GetAll returns all categories
func (m *MovieService) GetAll() (*MovieList, error) {
	return m.MovieRepo.GetAll()
}

func (m *MovieService) GetAllMoviesWithPagination(pag Pagination) (*Pagination, error) {
	movies, count, err := m.MovieRepo.GetAllMoviesWithPagination(pag)
	if err != nil {
		return nil, err
	}
	pag.Items = movies
	pag.TotalCount = int64(count)
	return &pag, nil
}
