package pagination

import (
	"reflect"
	"testing"

	. "github.com/kadirgonen/movie-api/api/model"
)

func TestNewPage(t *testing.T) {
	type args struct {
		p Pagination
	}
	tests := []struct {
		name string
		args args
		want *Pagination
	}{
		{name: "PaginationCalculateWithTrue", args: args{p: Pagination{Page: int64(1), PageCount: int64(15), PageSize: int64(3), TotalCount: int64(10)}}, want: &Pagination{Page: int64(1), PageCount: int64(4), PageSize: int64(3), TotalCount: int64(10)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewPage(tt.args.p); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPage() = %v, want %v", got, tt.want)
			}
		})
	}
}
