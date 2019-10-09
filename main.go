package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/janabe/cscoupler/database/memory"
	"github.com/janabe/cscoupler/domain"
	"github.com/janabe/cscoupler/handlers"
	"github.com/janabe/cscoupler/services"
	"github.com/janabe/cscoupler/util"
)

// HandlerFunc creates a handler from a normal func

func main() {
	fmt.Println("Running server, listening on port 3000...")

	http.Handle("/", http.FileServer(http.Dir("./views")))

	userRepo := memory.UserRepo{DB: make(map[string]domain.User)}
	studentRepo := memory.StudentRepo{DB: make(map[string]domain.Student)}

	userService := services.UserService{UserRepo: userRepo}
	studentService := services.StudentService{StudentRepo: studentRepo, UserService: userService}

	authHandler := handlers.AuthHandler{
		JWTKey:      util.GetJWTSecret("./.secret.json"),
		UserService: userService,
	}

	studentHandler := handlers.StudentHandler{
		StudentService: studentService,
		AuthHandler:    authHandler,
		Path:           "/students/",
	}

	authHandler.RegisterHandlers()
	studentHandler.RegisterHandlers()

	log.Fatal(http.ListenAndServe(":3000", nil))
}
