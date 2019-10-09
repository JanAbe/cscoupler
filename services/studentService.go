package services

import (
	"fmt"

	"github.com/janabe/cscoupler/domain"
)

// StudentService struct, containing all features
// the app supports regaring students
type StudentService struct {
	StudentRepo domain.StudentRepository
	UserService UserService
}

// Register registers a new Student
func (s StudentService) Register(student domain.Student) error {
	err := s.UserService.Register(student.User)
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = s.StudentRepo.Create(student)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// FindByID finds a student based on an identifier
func (s StudentService) FindByID(id string) (domain.Student, error) {
	student, err := s.StudentRepo.FindByID(id)
	if err != nil {
		fmt.Println(err)
		return domain.Student{}, err
	}

	return student, nil
}

// FindAll finds all students present
func (s StudentService) FindAll() ([]domain.Student, error) {
	students, err := s.StudentRepo.FindAll()
	if err != nil {
		fmt.Println(err)
		return []domain.Student{}, err
	}

	return students, nil
}
