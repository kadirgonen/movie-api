package handler

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	. "github.com/kadirgonen/movie-api/api/model"
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

func NewMovieHandler(r *gin.RouterGroup, c *MovieService, cfg *config.Config) {
	h := &MovieHandler{movieService: c, cfg: cfg}
	c.Migrate()
	r.POST("/upload", AuthorizationMiddleware(h.cfg), h.upload)
	r.GET("/", PaginationMiddleware(h.cfg), h.getmovies)

}

func (c *MovieHandler) Migrate() {
	c.movieService.Migrate()
}

// upload helps to user create bulk movie. If category has implemented before db, upload checks is without any fault.
func (ct *MovieHandler) upload(c *gin.Context) {
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

	count, str, err := ct.movieService.Upload(&file)

	if err != nil {
		zap.L().Error("movie.handler.upload", zap.Error(err))
		c.JSON(ErrorResponse(NewRestError(http.StatusInternalServerError, os.Getenv("UPLOAD_FILE_FAULT"), err.Error())))
		return
	}
	c.JSON(http.StatusOK, APIResponse{Code: http.StatusOK, Message: fmt.Sprintf("%v Record", count), Details: str})

}

func (ct *MovieHandler) getmovies(c *gin.Context) {
	val, res := c.Get("Pagination")
	if res == false {
		zap.L().Error("movie.handler.getmovies", zap.Bool("value: ", res))
		c.JSON(ErrorResponse(NewRestError(http.StatusInternalServerError, os.Getenv("NO_CONTEXT"), nil)))
		return
	}
	pag := val.(Pagination)
	p, err := ct.movieService.GetAllMoviesWithPagination(pag)
	if err != nil {
		zap.L().Error("category.handler.getmovies", zap.Error(err))
		c.JSON(ErrorResponse(NewRestError(http.StatusInternalServerError, os.Getenv("PAGINATION_FAULT"), nil)))
		return
	}
	c.JSON(http.StatusOK, APIResponse{Code: http.StatusOK, Message: os.Getenv("MOVIES_GET_SUCCESS"), Details: NewPage(*p)})
	return
}
