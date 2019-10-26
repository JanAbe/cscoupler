package domain

import (
	"errors"
	"strings"
)

// Student struct
type Student struct {
	ID               string
	University       string
	Skills           []string
	Experiences      []string
	ShortExperiences []string
	Status           Status
	User             User
	Resume           string // path to the resume of the student
}

// StudentRepository interface
type StudentRepository interface {
	Create(student Student) error
	Update(student Student) error
	FindByID(id string) (Student, error)
	FindAll() ([]Student, error) // todo: think of a way to return all students 1 by 1, not all in one go
}

// NewStudent creates a new student based on the provided input args
func NewStudent(id, uni string,
	skills []string,
	experiences []string,
	shortExperiences []string,
	user User,
	status Status,
	resume string) (Student, error) {

	if len(strings.TrimSpace(uni)) == 0 {
		return Student{}, errors.New("provided univeristy can't empty")
	}

	return Student{
		ID:               id,
		University:       strings.ToLower(uni),
		Skills:           skills,
		Experiences:      experiences,
		ShortExperiences: shortExperiences,
		User:             user,
		Status:           status,
		Resume:           resume,
	}, nil

}
