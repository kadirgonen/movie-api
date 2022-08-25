package mw

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	. "github.com/kadirgonen/movie-api/api/model"
	"github.com/kadirgonen/movie-api/app/pkg/config"
)

func PaginationMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		Page := c.Request.URL.Query().Get("page")
		PageSize := c.Request.URL.Query().Get("pagesize")
		intPageID := 0
		intPageSize := 0
		var err error
		pagevalue, _ := strconv.Atoi(Page)
		pagesizevalue, _ := strconv.Atoi(PageSize)
		if pagevalue <= 0 || pagesizevalue <= 0 {
			c.JSON(http.StatusForbidden, APIResponse{Code: http.StatusForbidden, Message: "Wrong Usage!"})
			c.Abort()
			return
		}
		if Page != "" || pagevalue <= 0 {
			intPageID, err = strconv.Atoi(Page)
			if err != nil {
				c.JSON(http.StatusForbidden, APIResponse{Code: http.StatusForbidden, Message: "Page cannot read!"})
				c.Abort()
				return
			}

		}
		if PageSize != "" || pagesizevalue <= 0 {
			intPageSize, err = strconv.Atoi(PageSize)
			if err != nil {
				c.JSON(http.StatusForbidden, APIResponse{Code: http.StatusForbidden, Message: "PageSize cannot read!"})
				c.Abort()
				return
			}
		}

		c.Set("Pagination", Pagination{Page: int64(intPageID), PageSize: int64(intPageSize)})
		c.Next()
		c.Abort()
		return
	}
}
