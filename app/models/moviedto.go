package model

import (
	api "github.com/kadirgonen/movie-api/api/model/movie"
)

func ResponseToMovie(m api.Movie) *Movie {
	return &Movie{
		ID:          *m.ID,
		Name:        *m.Name,
		Description: *m.Description,
		Type:        *m.Type,
	}
}

func MovieToResponse(m *Movie) *api.Movie {
	return &api.Movie{

		Name:        &m.Name,
		Description: &m.Description,
		Type:        &m.Type,
	}
}
