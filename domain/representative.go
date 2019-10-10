package domain

import (
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

// RepresentativeRepository interface
type RepresentativeRepository interface {
	Create(representative Representative) error
	FindByID(id string) (Representative, error)
}

// Representative struct conveying a
// representative/employee of a company
// that is looking in name of the company
// for students
type Representative struct {
	ID        string
	Position  string
	User      User
	CompanyID string
}

// NewRepresentative creates a new representative based on the provided input
func NewRepresentative(pos, companyID string, user User) (Representative, error) {
	if len(strings.TrimSpace(pos)) == 0 {
		return Representative{}, errors.New("provided position can't be empty")
	}

	return Representative{
		ID:        uuid.New().String(),
		Position:  strings.ToLower(pos),
		User:      user,
		CompanyID: companyID,
	}, nil
}

// GenerateInviteLink generates an invite-link to to sent to a non-user
// in order to become a representative of this representatives company
// urlTemplate should contain <[companyID]> and <[inviteID]>.
// e.g.: /signup/representatives/invite/<[companyID]>/<[inviteID]>
func (r Representative) GenerateInviteLink(inviteLinkID, urlTemplate string) (InviteLink, error) {
	if len(strings.TrimSpace(urlTemplate)) == 0 {
		return InviteLink{}, errors.New("provided url can't be empty")
	}

	// todo: this also needs to be updated as the inviteLink-id also needs to be there
	url := regexp.MustCompile(`<\[companyID\]>`).ReplaceAllString(urlTemplate, r.CompanyID)
	url = regexp.MustCompile(`<\[inviteID\]>`).ReplaceAllString(url, inviteLinkID)

	return NewInviteLink(
		inviteLinkID,
		url,
		r.CompanyID,
		time.Hour*24,
	), nil
}

// InviteLink struct conveying an invitelink
// that gets sent to bind a new user to a company
type InviteLink struct {
	ID           string
	URL          string
	CreationDate time.Time
	ExpiryDate   time.Time

	// id of an entity you want the new user to be bound with.
	// e.g. companyID if you want a user to be bound to that company
	// or universityID if you want a user to be bound to that university, etc.
	EntityID string
}

// InviteLinkRepository interface
type InviteLinkRepository interface {
	Create(inviteLink InviteLink) error
	FindByID(id string) (InviteLink, error)
}

// NewInviteLink creates a new InviteLink to be sent
// to a non-user to create a account and be bound to
// the provided entity. InviteLinks are valid for the
// amount of time specified by the validFor parameter
// e.g. 24 hours -> time.Hour * 24
func NewInviteLink(id, url, entityID string, validFor time.Duration) InviteLink {
	// todo: add check to see if link is valid?
	// add check to see if entityID is empty?
	return InviteLink{
		ID:           id,
		URL:          url,
		CreationDate: time.Now(),
		ExpiryDate:   time.Now().Add(validFor),
		EntityID:     entityID,
	}
}
