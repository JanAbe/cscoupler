package services

import (
	"github.com/google/uuid"
	d "github.com/janabe/cscoupler/domain"
)

// InviteLinkService struct, containing all features
// the app supports regarding invite links
type InviteLinkService struct {
	InviteLinkRepo d.InviteLinkRepository
}

// todo: refactor/ improve code. This is not conform open/closed principal
// if i want to add a feature to create an invitelink for universities,
// to bind a new student to that university, i need to create a new func in
// this file. This can be improved i think.

// CreateRepresentativeInvite creates a new invitelink
// for new representatives.
// Path should be a relative path, like: /signup/representatives/
func (i InviteLinkService) CreateRepresentativeInvite(path string, r d.Representative) (d.InviteLink, error) {
	// todo: replace util.URL to the domain name of the client, otherwise the created invitelink
	// points to the endpoint of the server and not to the front-end
	urlTemplate := "/signup" + path + "invite/<[companyID]>/<[inviteID]>"
	inviteLinkID := uuid.New().String()
	inviteLink, err := r.GenerateInviteLink(inviteLinkID, urlTemplate)
	if err != nil {
		return d.InviteLink{}, err
	}

	err = i.InviteLinkRepo.Create(inviteLink)
	if err != nil {
		return d.InviteLink{}, err
	}

	return inviteLink, nil
}

// FindByID fetches an inviteLink based on id
func (i InviteLinkService) FindByID(id string) (d.InviteLink, error) {
	inviteLink, err := i.InviteLinkRepo.FindByID(id)
	if err != nil {
		return d.InviteLink{}, err
	}

	return inviteLink, nil
}

// FindByCreator fetches all inviteLinks that are created by the provided id
func (i InviteLinkService) FindByCreator(representativeID string) ([]d.InviteLink, error) {
	inviteLinks, err := i.InviteLinkRepo.FindByCreator(representativeID)
	if err != nil {
		return []d.InviteLink{}, err
	}

	return inviteLinks, nil
}

// Update updates the invitelink
func (i InviteLinkService) Update(inviteLink d.InviteLink) error {
	err := i.InviteLinkRepo.Update(inviteLink)
	if err != nil {
		return err
	}

	return nil
}
