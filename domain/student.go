package domain

import (
	"errors"
	"strings"
	"time"
)

// Student struct
type Student struct {
	firstname  string
	lastname   string
	university string
	birthday   time.Time
	skills     []string
	experience []string
	user       User
	status     Status
	// iets met cv en afbeelding, hoe sla je dat op en hoe geef je het aan in de struct enz.
}

// StudentRepository interface
type StudentRepository interface {
	NextID() string
	Create(student *Student) error
	FindByID(id string) (*Student, error)
	FindAll() ([]*Student, error) // todo: think of a way to return all students 1 by 1, not all in one go
}

// NewStudent creates a new student based on the provided input args
func NewStudent(fname, lname, uni string,
	dob time.Time,
	skills, exp []string,
	user User, status Status) (Student, error) {

	if len(strings.TrimSpace(fname)) == 0 {
		return Student{}, errors.New("provided firstname can't be empty")
	}

	if len(strings.TrimSpace(lname)) == 0 {
		return Student{}, errors.New("provided lastname can't be empty")
	}

	if len(strings.TrimSpace(uni)) == 0 {
		return Student{}, errors.New("provided univeristy can't empty")
	}

	// Maybe add check to see if user is >=18 or >=16 or something?
	if dob.After(time.Now()) {
		return Student{}, errors.New("provided date of birth is invalid. it can't be after the current date")
	}

	return Student{
		firstname:  fname,
		lastname:   lname,
		university: uni,
		birthday:   dob,
		skills:     skills,
		experience: exp,
		user:       user,
		status:     status,
	}, nil

}
