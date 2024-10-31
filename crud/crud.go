package crud

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"golang_project/models"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

// CheckErr checks for an error and panics if it exists
func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

// HandleBooks handles the request to list all books
// @Summary List all books
// @Description Get a list of all books
// @Tags books
// @Produce json
// @Success 200 {array} models.Book
// @Router /books [get]
func HandleBooks(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "/Users/amir/Documents/newtestgo/test.db")
	if err != nil {
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM books")
	if err != nil {
		http.Error(w, "Error querying database", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var books []models.Book
	for rows.Next() {
		var book models.Book
		err = rows.Scan(&book.ID, &book.Title, &book.Author, &book.ISBN, &book.PublishedYear, &book.Genre)
		if err != nil {
			http.Error(w, "Error scanning book", http.StatusInternalServerError)
			return
		}
		books = append(books, book)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(books)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
	}
}

// CreateBook handles the request to create a new book
// @Summary Create a new book
// @Description Create a new book with the provided details
// @Tags books
// @Accept json
// @Produce json
// @Param book body models.Book true "Book"
// @Success 201 {string} string "Book created successfully"
// @Router /books/create [post]
func CreateBook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var book models.Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db, err := sql.Open("sqlite3", "/Users/amir/Documents/newtestgo/test.db")
	CheckErr(err)
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO books(Title, Author, ISBN, PublishedYear, Genre) VALUES(?, ?, ?, ?, ?)")
	CheckErr(err)
	_, err = stmt.Exec(book.Title, book.Author, book.ISBN, book.PublishedYear, book.Genre)
	CheckErr(err)

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Book created successfully"))
}

// ReadBook handles the request to read a book by ID
// @Summary Read a book by ID
// @Description Get the details of a book by its ID
// @Tags books
// @Produce json
// @Param id query string true "Book ID"
// @Success 200 {object} models.Book
// @Router /books/read [get]
func ReadBook(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing book ID", http.StatusBadRequest)
		return
	}

	db, err := sql.Open("sqlite3", "/Users/amir/Documents/newtestgo/test.db")
	CheckErr(err)
	defer db.Close()

	row := db.QueryRow("SELECT * FROM books WHERE ID = ?", id)
	var book models.Book
	err = row.Scan(&book.ID, &book.Title, &book.Author, &book.ISBN, &book.PublishedYear, &book.Genre)
	if err != nil {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Book found"))
}

// UpdateBook handles the request to update a book
// @Summary Update a book
// @Description Update the details of an existing book
// @Tags books
// @Accept json
// @Produce json
// @Param book body models.Book true "Book"
// @Success 200 {string} string "Book updated successfully"
// @Router /books/update [put]
func UpdateBook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var book models.Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db, err := sql.Open("sqlite3", "/Users/amir/Documents/newtestgo/test.db")
	CheckErr(err)
	defer db.Close()

	stmt, err := db.Prepare("UPDATE books SET Title = ?, Author = ?, ISBN = ?, PublishedYear = ?, Genre = ? WHERE ID = ?")
	CheckErr(err)
	_, err = stmt.Exec(book.Title, book.Author, book.ISBN, book.PublishedYear, book.Genre, book.ID)
	CheckErr(err)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Book updated successfully"))
}

// DeleteBook handles the request to delete a book
// @Summary Delete a book
// @Description Delete a book by its ID
// @Tags books
// @Param id query string true "Book ID"
// @Success 200 {string} string "Book deleted successfully"
// @Router /books/delete [delete]
func DeleteBook(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing book ID", http.StatusBadRequest)
		return
	}

	db, err := sql.Open("sqlite3", "/Users/amir/Documents/newtestgo/test.db")
	CheckErr(err)
	defer db.Close()

	stmt, err := db.Prepare("DELETE FROM books WHERE ID = ?")
	CheckErr(err)
	_, err = stmt.Exec(id)
	CheckErr(err)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Book deleted successfully"))
}

// CreateUser handles the request to create a new user
// @Summary Create a new user
// @Description Create a new user with the provided details
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.User true "User"
// @Success 201 {string} string "User created successfully"
// @Router /users/create [post]
func CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Hash the user's password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	db, err := sql.Open("sqlite3", "/Users/amir/Documents/newtestgo/test.db")
	CheckErr(err)
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO Users(name, email, membershipdate, is_active, password, role) VALUES(?, ?, ?, ?, ?, ?)")
	CheckErr(err)
	_, err = stmt.Exec(user.Name, user.Email, user.MembershipDate, user.IsActive, hashedPassword, user.Role)
	CheckErr(err)

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User created successfully"))
}

// ReadUser handles the request to read a user by ID
// @Summary Read a user by ID
// @Description Get the details of a user by their ID
// @Tags users
// @Produce json
// @Param id query string true "User ID"
// @Success 200 {object} models.User
// @Router /users/read [get]
func ReadUser(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing user ID", http.StatusBadRequest)
		return
	}

	db, err := sql.Open("sqlite3", "/Users/amir/Documents/newtestgo/test.db")
	CheckErr(err)
	defer db.Close()

	row := db.QueryRow("SELECT ID, name, email, membershipdate, is_active FROM Users WHERE ID = ? AND role='user'", id)
	var user models.User
	err = row.Scan(&user.ID, &user.Name, &user.Email, &user.MembershipDate, &user.IsActive)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error scanning user", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User found"))
}

// UpdateUser handles the request to update a user
// @Summary Update a user
// @Description Update the details of an existing user
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.User true "User"
// @Success 200 {string} string "User updated successfully"
// @Router /users/update [put]
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db, err := sql.Open("sqlite3", "/Users/amir/Documents/newtestgo/test.db")
	CheckErr(err)
	defer db.Close()

	stmt, err := db.Prepare("UPDATE Users SET name = ?, email = ?, membershipdate = ?, is_active = ? WHERE id = ?")
	CheckErr(err)
	_, err = stmt.Exec(user.Name, user.Email, user.MembershipDate, user.IsActive, user.ID)
	CheckErr(err)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User updated successfully"))
}

// DeleteUser handles the request to delete a user
// @Summary Delete a user
// @Description Delete a user by their ID
// @Tags users
// @Param id query string true "User ID"
// @Success 200 {string} string "User deleted successfully"
// @Router /users/delete [delete]
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing user ID", http.StatusBadRequest)
		return
	}

	db, err := sql.Open("sqlite3", "/Users/amir/Documents/newtestgo/test.db")
	CheckErr(err)
	defer db.Close()

	stmt, err := db.Prepare("DELETE FROM Users WHERE id = ?")
	CheckErr(err)
	_, err = stmt.Exec(id)
	CheckErr(err)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User deleted successfully"))
}

// CreateBookkeeper handles the request to create a new bookkeeper
// @Summary Create a new bookkeeper
// @Description Create a new bookkeeper with the provided details
// @Tags bookkeepers
// @Accept json
// @Produce json
// @Param bookkeeper body models.User true "Bookkeeper"
// @Success 201 {string} string "Bookkeeper created successfully"
// @Router /bookkeepers/create [post]
func CreateBookkeeper(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var bookkeeper models.User
	err := json.NewDecoder(r.Body).Decode(&bookkeeper)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db, err := sql.Open("sqlite3", "/Users/amir/Documents/newtestgo/test.db")
	CheckErr(err)
	defer db.Close()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(bookkeeper.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	stmt, err := db.Prepare("INSERT INTO Users(name, email, is_active, password, role, membershipdate) VALUES(?, ?, ?, ?, ?, ?)")
	CheckErr(err)
	_, err = stmt.Exec(bookkeeper.Name, bookkeeper.Email, bookkeeper.IsActive, hashedPassword, bookkeeper.Role, bookkeeper.MembershipDate)
	CheckErr(err)

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Bookkeeper created successfully"))
}

// ReadBookkeeper handles the request to read a bookkeeper by ID
// @Summary Read a bookkeeper by ID
// @Description Get the details of a bookkeeper by their ID
// @Tags bookkeepers
// @Produce json
// @Param id query string true "Bookkeeper ID"
// @Success 200 {object} models.User
// @Router /bookkeepers/read [get]
func ReadBookkeeper(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing bookkeeper ID", http.StatusBadRequest)
		return
	}

	db, err := sql.Open("sqlite3", "/Users/amir/Documents/newtestgo/test.db")
	CheckErr(err)
	defer db.Close()

	row := db.QueryRow("SELECT ID, name, email, is_active FROM Users WHERE ID = ? AND role='admin'", id)
	var bookkeeper models.User
	err = row.Scan(&bookkeeper.ID, &bookkeeper.Name, &bookkeeper.Email, &bookkeeper.IsActive)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Bookkeeper not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error scanning bookkeeper", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bookkeeper)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Bookkeeper found"))
}

// UpdateBookkeeper handles the request to update a bookkeeper
// @Summary Update a bookkeeper
// @Description Update the details of an existing bookkeeper
// @Tags bookkeepers
// @Accept json
// @Produce json
// @Param bookkeeper body models.User true "Bookkeeper"
// @Success 200 {string} string "Bookkeeper updated successfully"
// @Router /bookkeepers/update [put]
func UpdateBookkeeper(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var bookkeeper models.User
	err := json.NewDecoder(r.Body).Decode(&bookkeeper)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db, err := sql.Open("sqlite3", "/Users/amir/Documents/newtestgo/test.db")
	CheckErr(err)
	defer db.Close()

	stmt, err := db.Prepare("UPDATE Users SET name = ?, email = ?, is_active = ? WHERE ID = ? AND role='admin'")
	CheckErr(err)
	_, err = stmt.Exec(bookkeeper.Name, bookkeeper.Email, bookkeeper.IsActive, bookkeeper.ID)
	CheckErr(err)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Bookkeeper updated successfully"))
}

// DeleteBookkeeper handles the request to delete a bookkeeper
// @Summary Delete a bookkeeper
// @Description Delete a bookkeeper by their ID
// @Tags bookkeepers
// @Param id query string true "Bookkeeper ID"
// @Success 200 {string} string "Bookkeeper deleted successfully"
// @Router /bookkeepers/delete [delete]
func DeleteBookkeeper(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing bookkeeper ID", http.StatusBadRequest)
		return
	}

	db, err := sql.Open("sqlite3", "/Users/amir/Documents/newtestgo/test.db")
	CheckErr(err)
	defer db.Close()

	stmt, err := db.Prepare("DELETE FROM Users WHERE ID = ? AND role='admin'")
	CheckErr(err)
	_, err = stmt.Exec(id)
	CheckErr(err)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Bookkeeper deleted successfully"))
}
