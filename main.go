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
	companyRepo := memory.CompanyRepo{DB: make(map[string]domain.Company)}
	representativeRepo := memory.RepresentativeRepo{DB: make(map[string]domain.Representative)}
	inviteLinkRepo := memory.InviteLinkRepo{DB: make(map[string]domain.InviteLink)}

	userService := services.UserService{
		UserRepo: userRepo,
	}

	studentService := services.StudentService{
		StudentRepo: studentRepo,
		UserService: userService,
	}

	companyService := services.CompanyService{
		CompanyRepo: companyRepo,
	}

	representativeService := services.RepresentativeService{
		RepresentativeRepo: representativeRepo,
		CompanyService:     companyService,
		UserService:        userService,
	}

	inviteLinkService := services.InviteLinkService{
		InviteLinkRepo: inviteLinkRepo,
	}

	companyService.RepresentativeService = &representativeService

	authHandler := handlers.AuthHandler{
		JWTKey:      util.GetJWTSecret("./.secret.json"),
		UserService: userService,
	}

	studentHandler := handlers.StudentHandler{
		StudentService: studentService,
		AuthHandler:    authHandler,
		Path:           "/students/",
	}

	companyHandler := handlers.CompanyHandler{
		CompanyService: companyService,
		AuthHandler:    authHandler,
		Path:           "/companies/",
	}

	representativeHandler := handlers.RepresentativeHandler{
		RepresentativeService: representativeService,
		InviteLinkService:     inviteLinkService,
		AuthHandler:           authHandler,
		Path:                  "/representatives/",
	}

	authHandler.RegisterHandlers()
	studentHandler.RegisterHandlers()
	companyHandler.RegisterHandlers()
	representativeHandler.RegisterHandlers()

	log.Fatal(http.ListenAndServe(":3000", nil))
}
