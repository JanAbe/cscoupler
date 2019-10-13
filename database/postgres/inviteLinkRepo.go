package postgres

import (
	"database/sql"

	d "github.com/janabe/cscoupler/domain"
)

// InviteLinkRepo ...
type InviteLinkRepo struct {
	DB *sql.DB
}

// Create ...
func (i InviteLinkRepo) Create(inviteLink d.InviteLink) error {
	return nil
}

// FindByID ...
func (i InviteLinkRepo) FindByID(id string) (d.InviteLink, error) {
	return d.InviteLink{}, nil
}

// Update ...
func (i InviteLinkRepo) Update(inviteLink d.InviteLink) error {
	return nil
}
