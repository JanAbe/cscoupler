package server

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/rs/cors"

	pg "github.com/janabe/cscoupler/database/postgres"
	d "github.com/janabe/cscoupler/domain"
	"github.com/janabe/cscoupler/handlers"
	ser "github.com/janabe/cscoupler/services"
	"github.com/janabe/cscoupler/util"
)

// Server struct, conveying the application
type Server struct {
	db       *sql.DB
	handlers []handlers.Handler

	userService           ser.UserService
	companyService        ser.CompanyService
	studentService        ser.StudentService
	projectService        ser.ProjectService
	inviteLinkService     ser.InviteLinkService
	representativeService ser.RepresentativeService

	userRepo           d.UserRepository
	studentRepo        d.StudentRepository
	companyRepo        d.CompanyRepository
	projectRepo        d.ProjectRepository
	inviteLinkRepo     d.InviteLinkRepository
	representativeRepo d.RepresentativeRepository
}

// NewServer creates a new server which can be run
// to start the app
func NewServer(db *sql.DB) *Server {
	server := Server{db: db}
	server.initRepos()
	server.initServices()
	server.initHandlers()
	return &server
}

// Run runs the server
func (s *Server) Run() {
	fmt.Println("Running server, listening on port 3000...")
	// log.Fatal(http.ListenAndServeTLS(":3000", "./server/cert.pem", "./server/key.pem", nil))
	mux := http.DefaultServeMux
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8080"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
	})
	h := c.Handler(mux)
	log.Fatal(http.ListenAndServe(":3000", h))
}

func (s *Server) initRepos() {
	s.userRepo = pg.UserRepo{DB: s.db}
	s.inviteLinkRepo = pg.InviteLinkRepo{DB: s.db}
	s.studentRepo = pg.StudentRepo{DB: s.db, UserRepo: s.userRepo.(pg.UserRepo)}
	s.representativeRepo = pg.RepresentativeRepo{DB: s.db, UserRepo: s.userRepo.(pg.UserRepo)}
	s.companyRepo = pg.CompanyRepo{DB: s.db, ReprRepo: s.representativeRepo.(pg.RepresentativeRepo)}
	s.projectRepo = pg.ProjectRepo{DB: s.db}
}

func (s *Server) initServices() {
	s.userService = ser.UserService{UserRepo: s.userRepo}
	s.companyService = ser.CompanyService{CompanyRepo: s.companyRepo}
	s.inviteLinkService = ser.InviteLinkService{InviteLinkRepo: s.inviteLinkRepo}
	s.studentService = ser.StudentService{StudentRepo: s.studentRepo}
	s.representativeService = ser.RepresentativeService{
		RepresentativeRepo: s.representativeRepo,
		CompanyService:     s.companyService,
		UserService:        s.userService,
	}

	s.companyService.ReprService = &s.representativeService
	s.projectService = ser.ProjectService{ProjectRepo: s.projectRepo}
}

func (s *Server) initHandlers() {
	authHandler := handlers.AuthHandler{
		JWTKey:      util.GetJWTSecret("./.secret.json"),
		UserService: s.userService,
	}

	studentHandler := handlers.StudentHandler{
		StudentService: s.studentService,
		AuthHandler:    authHandler,
		Path:           "/students/",
	}

	companyHandler := handlers.CompanyHandler{
		CompanyService: s.companyService,
		AuthHandler:    authHandler,
		Path:           "/companies/",
	}

	representativeHandler := handlers.RepresentativeHandler{
		RepresentativeService: s.representativeService,
		InviteLinkService:     s.inviteLinkService,
		AuthHandler:           authHandler,
		Path:                  "/representatives/",
	}

	projectHandler := handlers.ProjectHandler{
		ProjectService: s.projectService,
		AuthHandler:    authHandler,
		Path:           "/projects/",
	}

	authHandler.Register()
	studentHandler.Register()
	companyHandler.Register()
	representativeHandler.Register()
	projectHandler.Register()
}
