package crud

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

func TestCreateBook(t *testing.T) {
	book := models.Book{
		Title:         "Test Book",
		Author:        "Test Author",
		ISBN:          "1234567890",
		PublishedYear: 2023,
		Genre:         "Test Genre",
	}
	jsonPayload, err := json.Marshal(book)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/books/create", bytes.NewBuffer(jsonPayload))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Login as bookkeeper
	cookie := loginAsBookkeeper(t)
	req.AddCookie(cookie)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateBook)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	expected := "Book created successfully"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestReadBook(t *testing.T) {
	req, err := http.NewRequest("GET", "/books/read?id=6", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ReadBook)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestUpdateBook(t *testing.T) {
	book := models.Book{
		ID:            1,
		Title:         "Updated Test Book",
		Author:        "Updated Test Author",
		ISBN:          "0987654321",
		PublishedYear: 2023,
		Genre:         "Updated Test Genre",
	}
	jsonPayload, err := json.Marshal(book)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("PUT", "/books/update", bytes.NewBuffer(jsonPayload))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Login as bookkeeper
	cookie := loginAsBookkeeper(t)
	req.AddCookie(cookie)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(UpdateBook)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := "Book updated successfully"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestDeleteBook(t *testing.T) {
	req, err := http.NewRequest("DELETE", "/books/delete?id=1", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Login as bookkeeper
	cookie := loginAsBookkeeper(t)
	req.AddCookie(cookie)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(DeleteBook)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := "Book deleted successfully"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}