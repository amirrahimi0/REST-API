package main

import (
	"fmt"
	handlers "golang_project/handler"
	"time"
)

// @title Golang Project API
// @version 1.0
// @description This is a sample server for a Golang project.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:9000
// @BasePath /

func main() {
	fmt.Println("Making Rest API.....")

	endTime := time.Now()
	fmt.Println("Current Time:", endTime.Format(time.RFC1123))

	handlers.HandleRequest()
}
