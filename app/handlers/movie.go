package handler

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-openapi/strfmt"
	. "github.com/kadirgonen/movie-api/api/model"
	. "github.com/kadirgonen/movie-api/api/model/movie"
	model "github.com/kadirgonen/movie-api/app/models"
	"github.com/kadirgonen/movie-api/app/pkg/config"
	. "github.com/kadirgonen/movie-api/app/pkg/errors"
	. "github.com/kadirgonen/movie-api/app/pkg/middleware"
	. "github.com/kadirgonen/movie-api/app/pkg/pagination"
	. "github.com/kadirgonen/movie-api/app/service"
	"go.uber.org/zap"
)

type MovieHandler struct {
	movieService *MovieService
	cfg          *config.Config
}

func NewMovieHandler(r *gin.RouterGroup, m *MovieService, cfg *config.Config) {
	h := &MovieHandler{movieService: m, cfg: cfg}
	m.Migrate()
	r.POST("/create", AuthorizationMiddleware(h.cfg), h.create)
	r.PUT("/:id", AuthorizationMiddleware(h.cfg), h.update)
	r.DELETE("/:id", AuthorizationMiddleware(h.cfg), h.delete)
	r.POST("/upload", AuthorizationMiddleware(h.cfg), h.upload)
	r.GET("/", PaginationMiddleware(h.cfg), h.getmovies)

}

func (m *MovieHandler) Migrate() {
	m.movieService.Migrate()
}

// create helps to create movie
func (mv *MovieHandler) create(c *gin.Context) {
	var req Movie
	if err := c.Bind(&req); err != nil {
		zap.L().Error("movie.handler.create", zap.Error(err))
		c.JSON(ErrorResponse(NewRestError(http.StatusBadRequest, os.Getenv("CHECK_YOUR_REQUEST"), nil)))
		return
	}

	if err := req.Validate(strfmt.NewFormats()); err != nil {
		zap.L().Error("movie.handler.validate", zap.Error(err))
		c.JSON(ErrorResponse(err))
		return
	}
	m, err := mv.movieService.Create(model.ResponseToMovie(req))
	if err != nil {
		zap.L().Error("movie.handler.create", zap.Error(err))
		c.JSON(ErrorResponse(err))
		return
	}
	c.JSON(http.StatusCreated, APIResponse{Code: http.StatusCreated, Message: os.Getenv("CREATE_MOVIE_SUCCESS"), Details: model.MovieToResponse(m)})
	return

}

// update helps user to update movie by id
func (mv *MovieHandler) update(c *gin.Context) {
	id := c.Param("id")
	a, err := strconv.Atoi(id)
	var req Movie
	if err := c.Bind(&req); err != nil {
		zap.L().Error("movie.handler.update", zap.Error(err))
		c.JSON(ErrorResponse(NewRestError(http.StatusBadRequest, os.Getenv("CHECK_YOUR_REQUEST"), err)))
		return
	}
	if err := req.Validate(strfmt.NewFormats()); err != nil {
		zap.L().Error("movie.handler.validate", zap.Error(err))
		c.JSON(ErrorResponse(err))
		return
	}
	m, err := mv.movieService.Update(model.ResponseToMovie(req), a)
	if err != nil {
		zap.L().Error("movie.handler.update", zap.Error(err))
		c.JSON(ErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, APIResponse{Code: http.StatusOK, Message: os.Getenv("UPDATE_MOVIE_SUCCESS"), Details: model.MovieToResponse(m)})
	return
}

// delete helps user to delete movie by id
func (mv *MovieHandler) delete(c *gin.Context) {
	id := c.Param("id")
	a, err := strconv.Atoi(id)
	res, err := mv.movieService.Delete(a)
	if err != nil {
		zap.L().Error("movie.handler.delete", zap.Error(err))
		c.JSON(ErrorResponse(err))
		return
	}
	if res {
		c.JSON(http.StatusOK, APIResponse{Code: http.StatusOK, Message: os.Getenv("DELETE_MOVIE_SUCCESS")})
		return
	} else {
		c.JSON(ErrorResponse(err))
		return
	}

}

// upload helps to user create bulk movie. If category has implemented before db, upload checks is without any fault.
func (mv *MovieHandler) upload(c *gin.Context) {
	file, handler, err := c.Request.FormFile("myFile")
	if err != nil {
		zap.L().Error("movie.handler.upload", zap.Error(err))
		c.JSON(ErrorResponse(NewRestError(http.StatusInternalServerError, os.Getenv("READ_FILE_FAULT"), nil)))
		return
	}
	defer file.Close()
	zap.L().Debug("movie.handler.upload:", zap.String("filename", fmt.Sprintf("Uploaded File: %+v\n", handler.Filename)))
	zap.L().Debug("movie.handler.upload:", zap.String("filesize", fmt.Sprintf("File Size: %+v\n", handler.Size)))
	zap.L().Debug("movie.handler.upload:", zap.String("MIME", fmt.Sprintf("MIME Header: %+v\n", handler.Header)))

	count, str, err := mv.movieService.Upload(&file)

	if err != nil {
		zap.L().Error("movie.handler.upload", zap.Error(err))
		c.JSON(ErrorResponse(NewRestError(http.StatusInternalServerError, os.Getenv("UPLOAD_FILE_FAULT"), err.Error())))
		return
	}
	c.JSON(http.StatusOK, APIResponse{Code: http.StatusOK, Message: fmt.Sprintf("%v Record", count), Details: str})

}

func (mv *MovieHandler) getmovies(c *gin.Context) {
	val, res := c.Get("Pagination")
	if res == false {
		zap.L().Error("movie.handler.getmovies", zap.Bool("value: ", res))
		c.JSON(ErrorResponse(NewRestError(http.StatusInternalServerError, os.Getenv("NO_CONTEXT"), nil)))
		return
	}
	pag := val.(Pagination)
	p, err := mv.movieService.GetAllMoviesWithPagination(pag)
	if err != nil {
		zap.L().Error("category.handler.getmovies", zap.Error(err))
		c.JSON(ErrorResponse(NewRestError(http.StatusInternalServerError, os.Getenv("PAGINATION_FAULT"), nil)))
		return
	}
	c.JSON(http.StatusOK, APIResponse{Code: http.StatusOK, Message: os.Getenv("MOVIES_GET_SUCCESS"), Details: NewPage(*p)})
	return
}
