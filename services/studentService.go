package services

import (
	"fmt"

	"github.com/janabe/cscoupler/domain"
)

// StudentService struct, containing all features
// the app supports regaring students
type StudentService struct {
	studentRepo domain.StudentRepository
}

// Create creates a new Student
func (s StudentService) Create(student domain.Student) {
	err := s.studentRepo.Create(student)
	if err != nil {
		fmt.Println(err)
	}
}

// FindByID finds a student based on an identifier
func (s StudentService) FindByID(id string) domain.Student {
	student, err := s.studentRepo.FindByID(id)
	if err != nil {
		fmt.Println(err)
	}

	return student
}

// FindAll finds all students present
func (s StudentService) FindAll() []domain.Student {
	students, err := s.studentRepo.FindAll()
	if err != nil {
		fmt.Println(err)
	}

	return students
}
