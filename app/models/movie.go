package model

import "gorm.io/gorm"

type Movie struct {
	gorm.Model
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
}
type MovieList []Movie
