package handlers

import (
	"fmt"
	"golang_project/auth"
	"golang_project/crud"
	"golang_project/filters"
	"log"
	"net/http"

	_ "golang_project/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

// MainPage handles the main page request
// @Summary Main Page
// @Description This is the main page.
// @Tags main
// @Produce plain
// @Success 200 {string} string "Welcome to the main page!"
// @Router / [get]
func MainPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the main page!\n")
	fmt.Fprintf(w, "Please visit /books to see the list of books\n")
	fmt.Fprintf(w, "Please visit /users to see the list of users")
	fmt.Fprintf(w, "Please visit /login to login")
	fmt.Fprintf(w, "Please visit /login/bookkeepers to login as a bookkeeper")
	fmt.Fprintf(w, "Please visit /admin to see the admin page")
	fmt.Fprintf(w, "Please visit /user to see the user page")
	fmt.Fprintf(w, "Please visit /books/create to create a book")
	fmt.Fprintf(w, "Please visit /books/update to update a book")
	fmt.Fprintf(w, "Please visit /books/delete to delete a book")
	fmt.Fprintf(w, "Please visit /users/update to update a user")
	fmt.Fprintf(w, "Please visit /users/delete to delete a user")
	fmt.Fprintf(w, "Please visit /bookkeepers/update to update a bookkeeper")
	fmt.Fprintf(w, "Please visit /bookkeepers/delete to delete a bookkeeper")
	fmt.Fprintf(w, "Please visit /bookkeepers/create to create a bookkeeper")
	fmt.Fprintf(w, "Please visit /bookkeepers/read to read a bookkeeper")
	fmt.Fprintf(w, "Please visit /secret to see the secret page")
}

// HandleRequest sets up the routes and starts the server
func HandleRequest() {
	http.HandleFunc("/", MainPage)
	http.HandleFunc("/login", auth.LoginUser)
	http.HandleFunc("/login/bookkeepers", auth.LoginBookkeeper)
	http.HandleFunc("/books", crud.HandleBooks)
	http.HandleFunc("/books/read", crud.ReadBook)
	http.HandleFunc("/books/filter/genre", filters.FilterBooksByGenre)
	http.HandleFunc("/books/filter/author", filters.FilterBooksByAuthor)
	http.HandleFunc("/books/filter/year", filters.FilterBooksByPublishedYear)
	http.HandleFunc("/books/filter/advanced", filters.AdvancedFilterBooks)
	http.HandleFunc("/books/search/title", filters.SearchBooksByTitle)
	http.HandleFunc("/users/create", crud.CreateUser)
	http.HandleFunc("/users/read", crud.ReadUser)

	http.Handle("/admin", auth.BookkeeperMiddleware(http.HandlerFunc(auth.AdminHandler)))
	http.Handle("/user", auth.AuthMiddleware(http.HandlerFunc(auth.UserHandler)))

	// Protected routes for bookkeepers
	http.Handle("/books/create", auth.BookkeeperMiddleware(http.HandlerFunc(crud.CreateBook)))
	http.Handle("/books/update", auth.BookkeeperMiddleware(http.HandlerFunc(crud.UpdateBook)))
	http.Handle("/books/delete", auth.BookkeeperMiddleware(http.HandlerFunc(crud.DeleteBook)))
	http.Handle("/users/update", auth.BookkeeperMiddleware(http.HandlerFunc(crud.UpdateUser)))
	http.Handle("/users/delete", auth.BookkeeperMiddleware(http.HandlerFunc(crud.DeleteUser)))
	http.Handle("/bookkeepers/update", auth.BookkeeperMiddleware(http.HandlerFunc(crud.UpdateBookkeeper)))
	http.Handle("/bookkeepers/delete", auth.BookkeeperMiddleware(http.HandlerFunc(crud.DeleteBookkeeper)))
	http.Handle("/bookkeepers/create", auth.BookkeeperMiddleware(http.HandlerFunc(crud.CreateBookkeeper)))
	http.Handle("/bookkeepers/read", auth.BookkeeperMiddleware(http.HandlerFunc(crud.ReadBookkeeper)))
	http.Handle("/secret", auth.BookkeeperMiddleware(http.HandlerFunc(SecretPage)))

	// Swagger endpoint
	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	log.Fatal(http.ListenAndServe(":9000", nil))
}

// SecretPage handles the secret page request
// @Summary Secret Page
// @Description This is the secret page.
// @Tags secret
// @Produce plain
// @Success 200 {string} string "Welcome to the secret page!"
// @Router /secret [get]
func SecretPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the secret page!\n")
	fmt.Fprintf(w, "You have successfully logged in!")
}
