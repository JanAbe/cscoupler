package domain

import "time"

// InviteLink struct conveying an invitelink
// that gets sent to bind a new user to a company
type InviteLink struct {
	ID         string
	URL        string
	CreatedAt  time.Time
	ExpiryDate time.Time
	CreatedBy  string // id of representative that created the invitation

	// keeps track if the link has been used or not
	// This can be used to make sure a link can only gets used once
	Used bool
}

// InviteLinkRepository interface
type InviteLinkRepository interface {
	Create(inviteLink InviteLink) error
	FindByID(id string) (InviteLink, error)
	Update(inviteLink InviteLink) error
	FindByCreator(representativeID string) ([]InviteLink, error)
}

// NewInviteLink creates a new InviteLink to be sent
// to a non-user to create a account and be bound to
// the provided entity. InviteLinks are valid for the
// amount of time specified by the validFor parameter
// e.g. 24 hours -> time.Hour * 24
func NewInviteLink(id, url, createdBy string, used bool, validFor time.Duration) InviteLink {
	// todo: add check to see if link is valid?
	return InviteLink{
		ID:         id,
		URL:        url,
		CreatedAt:  time.Now(),
		ExpiryDate: time.Now().Add(validFor),
		Used:       used,
		CreatedBy:  createdBy,
	}
}

// HasExpired checks if the expiry date of
// an invitelink has been reached
func (i InviteLink) HasExpired() bool {
	if time.Now().After(i.ExpiryDate) {
		return true
	}

	return false
}

// HasBeenUsed checks if the invitelink
// has been used or not
func (i InviteLink) HasBeenUsed() bool {
	return i.Used
}
