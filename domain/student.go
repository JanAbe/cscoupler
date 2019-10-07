package domain

import (
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
	// iets met cv en afbeelding, hoe sla je dat op en hoe geef je het aan in de struct enz.
}

// StudentRepository interface
type StudentRepository interface {
	NextID() string
	Create(student *Student) error
	FindByID(id string) (*Student, error)
	FindAll() ([]*Student, error) // todo: think of a way to return all students 1 by 1, not all in one go
}
