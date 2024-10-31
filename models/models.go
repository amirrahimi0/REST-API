package models

type Book struct {
	ID            int    `json:"id"`
	Title         string `json:"title"`
	Author        string `json:"author"`
	ISBN          string `json:"isbn"`
	PublishedYear int    `json:"published_year"`
	Genre         string `json:"genre"`
}

type User struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	MembershipDate string `json:"membership_date"`
	IsActive       bool   `json:"is_active"`
	Password       string `json:"password"`
	Role           string `json:"role"`
}

type Filter struct {
	Genre         string `json:"genre"`
	Author        string `json:"author"`
	PublishedYear string `json:"published_year"`
	Title         string `json:"title"`
	SortOrder     string `json:"sort_order"`
}
