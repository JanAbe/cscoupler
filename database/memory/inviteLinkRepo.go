package memory

import (
	"errors"

	"github.com/janabe/cscoupler/domain"
)

// InviteLinkRepo ...
type InviteLinkRepo struct {
	DB map[string]domain.InviteLink
}

// Create ...
func (i InviteLinkRepo) Create(inviteLink domain.InviteLink) error {
	i.DB[inviteLink.ID] = inviteLink
	return nil
}

// FindByID ...
func (i InviteLinkRepo) FindByID(id string) (domain.InviteLink, error) {
	if inviteLink, ok := i.DB[id]; ok {
		return inviteLink, nil
	}

	return domain.InviteLink{}, errors.New("no invitelink with id: " + id)
}

// Update ...
func (i InviteLinkRepo) Update(inviteLink domain.InviteLink) error {
	i.DB[inviteLink.ID] = inviteLink
	return nil
}
