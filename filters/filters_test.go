package filters

import (
	"bytes"
	"encoding/json"
	"golang_project/auth"
	"golang_project/models"
	"net/http"
	"net/http/httptest"
	"testing"
)

func loginAsBookkeeper(t *testing.T) *http.Cookie {
	credentials := auth.Credentials{
		Username: "amir@gmail.com",
		Password: "1234",
	}
	jsonPayload, err := json.Marshal(credentials)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/login/bookkeepers", bytes.NewBuffer(jsonPayload))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(auth.LoginBookkeeper)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Fatalf("handler returned wrong status code: got %v want %v. Response body: %v",
			status, http.StatusOK, rr.Body.String())
	}

	cookies := rr.Result().Cookies()
	for _, cookie := range cookies {
		if cookie.Name == "token" {
			return cookie
		}
	}

	t.Fatalf("Token cookie not found. Cookies: %v", cookies)
	return nil
}

func TestFilterBooksByGenre(t *testing.T) {
	req, err := http.NewRequest("GET", "/books/filter/genre?genre=Test Genre", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(FilterBooksByGenre)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var books []models.Book
	err = json.Unmarshal(rr.Body.Bytes(), &books)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	for _, book := range books {
		if book.Genre != "Test Genre" {
			t.Errorf("Expected all books to have genre 'Test Genre', but got %s", book.Genre)
		}
	}
}

func TestFilterBooksByAuthor(t *testing.T) {
	req, err := http.NewRequest("GET", "/books/filter/author?author=Test Author", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(FilterBooksByAuthor)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var books []models.Book
	err = json.Unmarshal(rr.Body.Bytes(), &books)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	for _, book := range books {
		if book.Author != "Test Author" {
			t.Errorf("Expected all books to have author 'Test Author', but got %s", book.Author)
		}
	}
}

func TestFilterBooksByPublishedYear(t *testing.T) {
	req, err := http.NewRequest("GET", "/books/filter/year?published_year=2023", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(FilterBooksByPublishedYear)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var books []models.Book
	err = json.Unmarshal(rr.Body.Bytes(), &books)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	for _, book := range books {
		if book.PublishedYear != 2023 {
			t.Errorf("Expected all books to have published year 2023, but got %d", book.PublishedYear)
		}
	}
}

func TestSearchBooksByTitle(t *testing.T) {
	req, err := http.NewRequest("GET", "/books/search/title?title=Test Book", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(SearchBooksByTitle)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var books []models.Book
	err = json.Unmarshal(rr.Body.Bytes(), &books)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	for _, book := range books {
		if book.Title != "Test Book" {
			t.Errorf("Expected all books to have title 'Test Book', but got %s", book.Title)
		}
	}
}

func TestAdvancedFilterBooks(t *testing.T) {
	filter := models.Filter{
		Genre:        "Test Genre",
		Author:       "Test Author",
		PublishedYear: "2023",
		Title:        "Test Book",
		SortOrder:    "asc",
	}
	jsonPayload, err := json.Marshal(filter)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/books/filter/advanced", bytes.NewBuffer(jsonPayload))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(AdvancedFilterBooks)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var books []models.Book
	err = json.Unmarshal(rr.Body.Bytes(), &books)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	for _, book := range books {
		if book.Genre != "Test Genre" || book.Author != "Test Author" || book.PublishedYear != 2023 || book.Title != "Test Book" {
			t.Errorf("Book does not match filter criteria: %+v", book)
		}
	}

	// Check if books are sorted in ascending order by PublishedYear
	for i := 1; i < len(books); i++ {
		if books[i].PublishedYear < books[i-1].PublishedYear {
			t.Errorf("Books are not sorted in ascending order by PublishedYear")
			break
		}
	}
}
