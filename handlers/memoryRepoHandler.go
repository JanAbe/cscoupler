package handlers

import (
	"errors"

	"github.com/google/uuid"
	"github.com/janabe/cscoupler/domain"
)

type memoryStudentRepo struct {
	db map[string]domain.User
}

func (m memoryStudentRepo) NextID() string {
	return uuid.New().String()
}

func (m memoryStudentRepo) Create(user *domain.User) error {
	m.db[user.ID] = *user
	return nil
}

func (m memoryStudentRepo) FindByID(id string) (*domain.User, error) {
	if user, ok := m.db[id]; ok {
		return &user, nil
	}

	return &domain.User{}, errors.New("No user with id " + id)
}

func (m memoryStudentRepo) FindByEmail(email string) (*domain.User, error) {
	for _, user := range m.db {
		if user.Email == email {
			return &user, nil
		}
	}

	return &domain.User{}, errors.New("No user with email " + email)
}
