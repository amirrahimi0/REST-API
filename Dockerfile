# Use the official Golang image as a parent image
FROM golang:1.23.0-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application's source code
COPY . .

# Build the application
RUN go build -o main .

# Expose port 9000
EXPOSE 9000

# Run the application
CMD ["./golang_project"]