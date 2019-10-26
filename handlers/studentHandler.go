package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/janabe/cscoupler/util"

	"github.com/google/uuid"
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
	ID         string   `json:"id"`
	University string   `json:"university"`
	Skills     []string `json:"skills"`
	Experience []string `json:"experience"`
	Status     string   `json:"status"`
	Resume     string   `json:"resume"`
	UserData   UserData `json:"user"`
}

// SignupStudent signs up a new student
func (s StudentHandler) SignupStudent() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			return
		}

		resumePath, err := processResume(r)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var data StudentData

		// check if json is invalid
		err = json.Unmarshal([]byte(r.FormValue("studentData")), &data)
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
			uuid.New().String(),
			data.University,
			data.Skills,
			data.Experience,
			user,
			domain.Available,
			resumePath,
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

// EditStudent edits a student account
func (s StudentHandler) EditStudent() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			return
		}

		studentID := strings.TrimPrefix(r.URL.Path, s.Path+"edit/")
		student, err := s.StudentService.FindByID(studentID)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		resumePath, err := processResume(r)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var updatedData StudentData

		// check if json is invalid
		err = json.Unmarshal([]byte(r.FormValue("studentData")), &updatedData)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		updatedUser, err := domain.NewUser(
			updatedData.UserData.Email,
			updatedData.UserData.Password,
			updatedData.UserData.Firstname,
			updatedData.UserData.Lastname,
			domain.StudentRole,
		)

		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		updatedStudent, err := domain.NewStudent(
			studentID,
			updatedData.University,
			updatedData.Skills,
			updatedData.Experience,
			updatedUser,
			domain.Available,
			resumePath,
		)

		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = s.StudentService.Edit(updatedStudent)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest) // what header to return
			return
		}

		// remove the old resume file of the student
		err = os.Remove(student.Resume)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(updatedStudent.ID)
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

// FetchAllStudents fetches all the students
func (s StudentHandler) FetchAllStudents() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" && r.Method != "OPTIONS" {
			return
		}

		students, err := s.StudentService.FindAll()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Println(err)
			return
		}

		var studentsData []StudentData
		for _, s := range students {
			studentsData = append(studentsData, ToStudentData(s))
		}

		json.NewEncoder(w).Encode(studentsData)
	})
}

// Register registers all student related handlers
func (s StudentHandler) Register() {
	http.Handle(s.Path, LoggingHandler(os.Stdout, s.AuthHandler.Validate("", s.FetchStudentByID())))
	http.Handle(s.Path+"edit/", LoggingHandler(os.Stdout, s.AuthHandler.Validate(domain.StudentRole, s.EditStudent())))
	http.Handle("/signup/student", LoggingHandler(os.Stdout, s.SignupStudent()))
	http.Handle(s.Path+"all/", LoggingHandler(os.Stdout, s.AuthHandler.Validate("", s.FetchAllStudents())))
}

// Helper func to extract the uploaded resume file
// and store it on the server. It returns the absolute path
// to this file
func processResume(r *http.Request) (string, error) {
	// Create copy of the sent file
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("resume")
	if err != nil {
		return "", err
	}

	isPdf := util.HasCorrectContentType(file, "application/pdf")
	if !isPdf {
		return "", errors.New("incorrect content type")
	}

	defer file.Close()

	resumePath, err := filepath.Abs("./resumes/" + uuid.New().String() + "-" + handler.Filename)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	dest, err := os.OpenFile(resumePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer dest.Close()
	io.Copy(dest, file)

	return resumePath, nil
}
