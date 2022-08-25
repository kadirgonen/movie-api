package pagination

import (
	. "github.com/kadirgonen/movie-api/api/model"
)

var (
	// DefaultPageSize specifies the default page size
	DefaultPageSize = 100
	// MaxPageSize specifies the maximum page size
	MaxPageSize = 10000
	// PageVar specifies the query parameter name for page number

)

func NewPage(p Pagination) *Pagination {
	if p.PageSize <= 0 {
		p.PageSize = int64(DefaultPageSize)
	}
	if int(p.PageSize) > MaxPageSize {
		p.PageSize = int64(MaxPageSize)
	}
	p.PageCount = int64(-1)
	if p.TotalCount >= 0 {
		p.PageCount = (p.TotalCount + p.PageSize - 1) / p.PageSize
	}

	return &Pagination{
		Page:       p.Page,
		PageSize:   p.PageSize,
		TotalCount: p.TotalCount,
		PageCount:  p.PageCount,
		Items:      p.Items,
	}
}
