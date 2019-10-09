package memory

import (
	"errors"

	"github.com/janabe/cscoupler/domain"
)

// StudentRepo ...
type StudentRepo struct {
	DB map[string]domain.Student
}

// Create ...
func (s StudentRepo) Create(student domain.Student) error {
	s.DB[student.ID] = student
	return nil
}

// FindByID ...
func (s StudentRepo) FindByID(id string) (domain.Student, error) {
	if student, ok := s.DB[id]; ok {
		return student, nil
	}

	return domain.Student{}, errors.New("No student with id " + id)
}

// FindAll ...
func (s StudentRepo) FindAll() ([]domain.Student, error) {
	students := []domain.Student{}
	for _, student := range s.DB {
		students = append(students, student)
	}

	return students, nil
}
