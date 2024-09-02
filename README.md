# Golang Project: REST-API Book Management System

This project is a RESTful API for managing books and users, built with Go. It includes features such as CRUD operations for books and users, authentication, and advanced filtering capabilities.

## Features

- CRUD operations for books and users
- Authentication for users and bookkeepers
- Advanced filtering and searching for books
- Swagger documentation for API endpoints
- Dockerized application for easy deployment

## Tech Stack

- Go 1.23.0
- SQLite3 for database
- JWT for authentication
- Swagger for API documentation
- Docker for containerization

## Project Structure
```
golang_project/
├── auth/
│ └── auth.go
├── crud/
│ ├── crud.go
│ └── crud_test.go
├── docs/
│ ├── docs.go
│ ├── swagger.json
│ └── swagger.yaml
├── filters/
│ ├── filters.go
│ └── filters_test.go
├── models/
│ └── models.go
├── go.mod
├── go.sum
├── main.go
└── README.md
```

## Setup and Installation

### Prerequisites

- Go 1.23.0 or later
- Docker and Docker Compose

### Running the Application

1. Clone the repository:
   ```
   git clone https://github.com/yourusername/golang_project.git
   cd golang_project
   ```

2. Build and run the Docker containers:
   ```
   docker-compose up --build
   ```

3. The application will be available at `http://localhost:9000`
4. Swagger UI will be available at `http://localhost:8080`

## API Endpoints

### Main
- `GET /`: Main page with links to other endpoints

### Authentication
- `POST /login`: User login
- `POST /login/bookkeepers`: Bookkeeper login

### Books
- `GET /books`: List all books
- `GET /books/read`: Read a specific book
- `POST /books/create`: Create a new book (Bookkeeper only)
- `PUT /books/update`: Update a book (Bookkeeper only)
- `DELETE /books/delete`: Delete a book (Bookkeeper only)

### Book Filtering
- `GET /books/filter/genre`: Filter books by genre
- `GET /books/filter/author`: Filter books by author
- `GET /books/filter/year`: Filter books by published year
- `POST /books/filter/advanced`: Advanced filtering for books
- `GET /books/search/title`: Search books by title

### Users
- `POST /users/create`: Create a new user
- `GET /users/read`: Read a specific user
- `PUT /users/update`: Update a user (Bookkeeper only)
- `DELETE /users/delete`: Delete a user (Bookkeeper only)

### Bookkeepers
- `POST /bookkeepers/create`: Create a new bookkeeper (Bookkeeper only)
- `GET /bookkeepers/read`: Read a specific bookkeeper (Bookkeeper only)
- `PUT /bookkeepers/update`: Update a bookkeeper (Bookkeeper only)
- `DELETE /bookkeepers/delete`: Delete a bookkeeper (Bookkeeper only)

### Protected Pages
- `GET /admin`: Admin page (Bookkeeper only)
- `GET /user`: User page (Authenticated users)
- `GET /secret`: Secret page (Bookkeeper only)

### Documentation
- `GET /swagger/`: Swagger UI for API documentation

For detailed request/response schemas and examples, please refer to the Swagger UI available at `http://localhost:8080` when running the application.

Note: Endpoints marked with "Bookkeeper only" require authentication as a bookkeeper to access.

For detailed API documentation, please refer to the Swagger UI at `http://localhost:8080`.

## Authentication

The application uses JWT for authentication. To access protected endpoints, you need to include the JWT token in the `Authorization` header of your requests.

## Testing

To run the tests:

`go test ./...`

