package services

import (
	"github.com/janabe/cscoupler/domain"
)

// StudentService struct, containing all features
// the app supports regaring students
type StudentService struct {
	StudentRepo domain.StudentRepository
}

// Register registers a new Student
func (s StudentService) Register(student domain.Student) error {
	err := s.StudentRepo.Create(student)
	if err != nil {
		return err
	}

	return nil
}

// Edit edits a student's information
func (s StudentService) Edit(student domain.Student) error {
	err := s.StudentRepo.Update(student)
	if err != nil {
		return err
	}

	return nil
}

// FindByID finds a student based on an identifier
func (s StudentService) FindByID(id string) (domain.Student, error) {
	student, err := s.StudentRepo.FindByID(id)
	if err != nil {
		return domain.Student{}, err
	}

	return student, nil
}

// FindAll finds all students present
func (s StudentService) FindAll() ([]domain.Student, error) {
	students, err := s.StudentRepo.FindAll()
	if err != nil {
		return []domain.Student{}, err
	}

	return students, nil
}
