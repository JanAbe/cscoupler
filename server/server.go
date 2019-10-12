package server

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

// Server struct, conveying the application
type Server struct {
	handlers []handlers.Handler
}

// NewServer creates a new server which can be run
// to start the app
func NewServer() *Server {
	server := Server{}
	server.init()
	server.registerHandlers()

	return &server
}

// Run runs the server
func (s *Server) Run() {
	fmt.Println("Running server, listening on port 3000...")
	log.Fatal(http.ListenAndServeTLS(":3000", "./server/cert.pem", "./server/key.pem", nil))
}

// registerHandlers registers all handlers of Server s
func (s *Server) registerHandlers() {
	for _, handler := range s.handlers {
		handler.Register()
	}
}

// init creates all necessary repositories,
// services and handlers.
func (s *Server) init() {
	userRepo := memory.UserRepo{DB: make(map[string]domain.User)}
	studentRepo := memory.StudentRepo{DB: make(map[string]domain.Student)}
	companyRepo := memory.CompanyRepo{DB: make(map[string]domain.Company)}
	representativeRepo := memory.RepresentativeRepo{DB: make(map[string]domain.Representative)}
	inviteLinkRepo := memory.InviteLinkRepo{DB: make(map[string]domain.InviteLink)}

	userService := services.UserService{UserRepo: userRepo}
	companyService := services.CompanyService{CompanyRepo: companyRepo}
	inviteLinkService := services.InviteLinkService{InviteLinkRepo: inviteLinkRepo}
	studentService := services.StudentService{StudentRepo: studentRepo, UserService: userService}

	representativeService := services.RepresentativeService{
		RepresentativeRepo: representativeRepo,
		CompanyService:     companyService,
		UserService:        userService,
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

	s.handlers = append(s.handlers, authHandler)
	s.handlers = append(s.handlers, studentHandler)
	s.handlers = append(s.handlers, companyHandler)
	s.handlers = append(s.handlers, representativeHandler)
}
