package filters

import (
	"database/sql"
	"encoding/json"
	"golang_project/models"
	"net/http"
	"strings"
)

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

// FilterBooksByGenre filters books by genre
// @Summary Filter Books by Genre
// @Description Filter books by genre
// @Tags books
// @Produce json
// @Param genre query string true "Genre"
// @Success 200 {array} models.Book
// @Router /books/filter/genre [get]
func FilterBooksByGenre(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "/Users/amir/Documents/newtestgo/test.db")
	CheckErr(err)
	defer db.Close()

	genre := r.URL.Query().Get("genre")

	rows, err := db.Query("SELECT ID, Title, Author, ISBN, PublishedYear, Genre FROM books WHERE Genre = ?", genre)
	CheckErr(err)
	defer rows.Close()

	var books []models.Book
	for rows.Next() {
		var book models.Book
		err = rows.Scan(&book.ID, &book.Title, &book.Author, &book.ISBN, &book.PublishedYear, &book.Genre)
		CheckErr(err)
		books = append(books, book)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

// FilterBooksByAuthor filters books by author
// @Summary Filter Books by Author
// @Description Filter books by author
// @Tags books
// @Produce json
// @Param author query string true "Author"
// @Success 200 {array} models.Book
// @Router /books/filter/author [get]
func FilterBooksByAuthor(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "/Users/amir/Documents/newtestgo/test.db")
	CheckErr(err)
	defer db.Close()

	author := r.URL.Query().Get("author")

	rows, err := db.Query("SELECT ID, Title, Author, ISBN, PublishedYear, Genre FROM books WHERE Author = ?", author)
	CheckErr(err)
	defer rows.Close()

	var books []models.Book
	for rows.Next() {
		var book models.Book
		err = rows.Scan(&book.ID, &book.Title, &book.Author, &book.ISBN, &book.PublishedYear, &book.Genre)
		CheckErr(err)
		books = append(books, book)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

// FilterBooksByPublishedYear filters books by published year
// @Summary Filter Books by Published Year
// @Description Filter books by published year
// @Tags books
// @Produce json
// @Param published_year query string true "Published Year"
// @Success 200 {array} models.Book
// @Router /books/filter/year [get]
func FilterBooksByPublishedYear(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "/Users/amir/Documents/newtestgo/test.db")
	CheckErr(err)
	defer db.Close()

	publishedYear := r.URL.Query().Get("published_year")

	rows, err := db.Query("SELECT ID, Title, Author, ISBN, PublishedYear, Genre FROM books WHERE PublishedYear = ?", publishedYear)
	CheckErr(err)
	defer rows.Close()

	var books []models.Book
	for rows.Next() {
		var book models.Book
		err = rows.Scan(&book.ID, &book.Title, &book.Author, &book.ISBN, &book.PublishedYear, &book.Genre)
		CheckErr(err)
		books = append(books, book)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

// SearchBooksByTitle searches books by title
// @Summary Search Books by Title
// @Description Search books by title
// @Tags books
// @Produce json
// @Param title query string true "Title"
// @Success 200 {array} models.Book
// @Router /books/search/title [get]
func SearchBooksByTitle(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "/Users/amir/Documents/newtestgo/test.db")
	CheckErr(err)
	defer db.Close()

	title := r.URL.Query().Get("title")

	rows, err := db.Query("SELECT ID, Title, Author, ISBN, PublishedYear, Genre FROM books WHERE Title = ?", title)
	CheckErr(err)
	defer rows.Close()

	var books []models.Book
	for rows.Next() {
		var book models.Book
		err = rows.Scan(&book.ID, &book.Title, &book.Author, &book.ISBN, &book.PublishedYear, &book.Genre)
		CheckErr(err)
		books = append(books, book)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

// AdvancedFilterBooks filters books based on multiple criteria
// @Summary Advanced Filter Books
// @Description Filter books based on multiple criteria
// @Tags books
// @Produce json
// @Param filter body models.Filter true "Filter"
// @Success 200 {array} models.Book
// @Router /books/filter/advanced [post]
func AdvancedFilterBooks(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "/Users/amir/Documents/newtestgo/test.db")
	CheckErr(err)
	defer db.Close()

	var filter models.Filter
	err = json.NewDecoder(r.Body).Decode(&filter)
	CheckErr(err)

	var filters []string
	var args []interface{}

	if filter.Genre != "" {
		filters = append(filters, "Genre = ?")
		args = append(args, filter.Genre)
	}
	if filter.Author != "" {
		filters = append(filters, "Author = ?")
		args = append(args, filter.Author)
	}
	if filter.PublishedYear != "" {
		filters = append(filters, "PublishedYear = ?")
		args = append(args, filter.PublishedYear)
	}
	if filter.Title != "" {
		filters = append(filters, "Title = ?")
		args = append(args, filter.Title)
	}

	query := "SELECT ID, Title, Author, ISBN, PublishedYear, Genre FROM books"
	if len(filters) > 0 {
		query += " WHERE " + strings.Join(filters, " AND ")
	}

	// Add sorting order based on published year
	if filter.SortOrder == "asc" {
		query += " ORDER BY PublishedYear ASC"
	} else if filter.SortOrder == "desc" {
		query += " ORDER BY PublishedYear DESC"
	}

	rows, err := db.Query(query, args...)
	CheckErr(err)
	defer rows.Close()

	var books []models.Book
	for rows.Next() {
		var book models.Book
		err = rows.Scan(&book.ID, &book.Title, &book.Author, &book.ISBN, &book.PublishedYear, &book.Genre)
		CheckErr(err)
		books = append(books, book)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}
