package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/janabe/cscoupler/domain"
	e "github.com/janabe/cscoupler/errors"
	"github.com/janabe/cscoupler/services"
)

// StudentHandler struct containing all
// student related handler funcs
type StudentHandler struct {
	StudentService services.StudentService
	AuthHandler    AuthHandler
	Path           string
}

// StudentData is a struct that corresponds to incoming student data
type StudentData struct {
	University string   `json:"university"`
	Skills     []string `json:"skills"`
	Experience []string `json:"experience"`
	UserData   UserData `json:"user"`
}

// SignupStudent signs up a new student
func (s StudentHandler) SignupStudent() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			return
		}

		var data StudentData

		// check if json is invalid
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		user, err := domain.NewUser(
			data.UserData.Email,
			data.UserData.Password,
			data.UserData.Firstname,
			data.UserData.Lastname,
			domain.StudentRole,
		)

		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		student, err := domain.NewStudent(
			data.University,
			data.Skills,
			data.Experience,
			user,
			domain.Available,
		)

		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = s.StudentService.Register(student)
		if err == e.ErrorEmailAlreadyUsed {
			fmt.Println(err)
			w.WriteHeader(http.StatusConflict)
			return
		}

		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		json.NewEncoder(w).Encode(student.ID)
	})
}

// FetchStudentByID fetches a student based on ID
// path = /students/... where the dots are a student ID
func (s StudentHandler) FetchStudentByID() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			return
		}

		id := strings.TrimPrefix(r.URL.Path, s.Path)
		student, err := s.StudentService.FindByID(id)

		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		studentData := ToStudentData(student)

		json.NewEncoder(w).Encode(studentData)
	})
}

// Register registers all student related handlers
func (s StudentHandler) Register() {
	http.Handle(s.Path, LoggingHandler(os.Stdout, s.AuthHandler.Validate("", s.FetchStudentByID())))
	http.Handle("/signup/student", LoggingHandler(os.Stdout, s.SignupStudent()))
}
