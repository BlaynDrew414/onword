package models

type Book struct {
	Title     string `json:"title,omitempty" validate:"required"`
	Subtitle  string `json:"subtitle,omitempty" validate:"required"`
	Author    string `json:"author,omitempty" validate:"required"`
	BookCover string `json:"bookCover,omitempty"`
}
