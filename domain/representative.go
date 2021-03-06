package domain

import (
	"errors"
	"regexp"
	"strings"
	"time"
)

// RepresentativeRepository interface
type RepresentativeRepository interface {
	Create(representative Representative) error
	FindByID(id string) (Representative, error)
	Update(representative Representative) error
}

// Representative struct conveying a
// representative/employee of a company
// that is looking in name of the company
// for students
type Representative struct {
	ID        string
	JobTitle  string
	User      User
	CompanyID string
}

// NewRepresentative creates a new representative based on the provided input
func NewRepresentative(id, jobTitle, companyID string, user User) (Representative, error) {
	if len(strings.TrimSpace(jobTitle)) == 0 {
		return Representative{}, errors.New("provided jobTitle can't be empty")
	}

	return Representative{
		ID:        id,
		JobTitle:  strings.ToLower(jobTitle),
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
		r.ID,
		false,
		time.Hour*24,
	), nil
}

// CreateProject creates a new project for the company of
// the representative
func (r Representative) CreateProject(projectID, desc, comp, dur string, recs []string) (Project, error) {
	return NewProject(projectID, desc, comp, dur, r.CompanyID, recs)
}
