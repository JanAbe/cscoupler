package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// StudentHandler struct containing all student
// related handler funcs
// type StudentHandler struct {
// 	StudentService services.StudentService
// 	Path           string
// }

// FetchStudentByID fetches a student by ID
// func (s StudentHandler) FetchStudentByID() http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		if r.Method != "GET" {
// 			return
// 		}
// 		id := strings.TrimPrefix(r.URL.Path, s.Path)
// 		s.StudentService.FindByID(id)
// 	})
// }

// FetchStudentByID ...
var FetchStudentByID = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/students/")
	student, err := studentService.FindByID(id)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	studentData := StudentData{
		University: student.University,
		Skills:     student.Skills,
		Experience: student.Experience,
		UserData: UserData{
			Email:     student.User.Email,
			Firstname: student.User.Firstname,
			Lastname:  student.User.Lastname,
		},
	}

	json.NewEncoder(w).Encode(studentData)
})
