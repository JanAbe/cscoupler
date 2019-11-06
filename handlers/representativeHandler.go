package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"

	"github.com/janabe/cscoupler/domain"
	e "github.com/janabe/cscoupler/errors"
	"github.com/janabe/cscoupler/services"
)

// RepresentativeHandler struct containing all
// representative related handler funcs
type RepresentativeHandler struct {
	RepresentativeService services.RepresentativeService
	InviteLinkService     services.InviteLinkService
	AuthHandler           AuthHandler
	Path                  string
}

// RepresentativeData is a struct that corresponds to incoming
// representative data
type RepresentativeData struct {
	JobTitle  string   `json:"jobTitle"`
	CompanyID string   `json:"companyID"`
	UserData  UserData `json:"user"`
}

// SignupRepresentative signs up a representative and binds
// it to the companyID present in the invite-link URL
// Format for invite-links: /signup/representatives/invite/[companyID]/[invitelinkID]
func (r RepresentativeHandler) SignupRepresentative() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			return
		}

		ids := strings.TrimPrefix(req.URL.Path, "/signup"+r.Path+"invite/")
		companyID := strings.Split(ids, "/")[0]
		inviteID := strings.Split(ids, "/")[1]

		inviteLink, err := r.InviteLinkService.FindByID(inviteID)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if inviteLink.HasExpired() {
			// what to return if an invitelink has expired
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if inviteLink.HasBeenUsed() {
			// what to return if an invitelink has been used already
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var data RepresentativeData

		// check if json is invalid
		err = json.NewDecoder(req.Body).Decode(&data)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		user, err := domain.NewUser(
			data.UserData.Email,
			data.UserData.Password,
			data.UserData.Firstname,
			data.UserData.Lastname,
			domain.RepresentativeRole,
		)

		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		representative, err := domain.NewRepresentative(
			data.JobTitle,
			companyID,
			user,
		)

		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = r.RepresentativeService.Register(representative)
		if err == e.ErrorEmailAlreadyUsed {
			fmt.Println(err)
			w.WriteHeader(http.StatusConflict)
			return
		}

		if err == e.ErrorEntityNotFound {
			fmt.Println(err)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		inviteLink.Used = true
		err = r.InviteLinkService.Update(inviteLink)
		if err != nil {
			fmt.Println(err)
			// which status to return?
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		json.NewEncoder(w).Encode(representative.ID)
	})
}

// FetchRepresentativeByID fetches a representative by ID
func (r RepresentativeHandler) FetchRepresentativeByID() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Method != "GET" {
			return
		}

		id := strings.TrimPrefix(req.URL.Path, r.Path)
		representative, err := r.RepresentativeService.FindByID(id)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		representativeData := ToRepresentativeData(representative)

		json.NewEncoder(w).Encode(representativeData)
	})
}

// FetchCreatedInvitations fetch all created invitations by the representative
func (r RepresentativeHandler) FetchCreatedInvitations() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Method != "GET" {
			return
		}

		cookie, _ := req.Cookie("token")
		token, _ := r.AuthHandler.GetToken(cookie)
		representativeID := token.Claims.(jwt.MapClaims)["ID"].(string)

		_, err := r.RepresentativeService.FindByID(representativeID)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		invitations, err := r.InviteLinkService.FindByCreator(representativeID)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(invitations)
	})
}

// MakeInviteLink makes an invite link for the representative to sent
// to colleagues.
func (r RepresentativeHandler) MakeInviteLink() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Method != "GET" {
			return
		}

		cookie, _ := req.Cookie("token")
		token, _ := r.AuthHandler.GetToken(cookie)
		representativeID := token.Claims.(jwt.MapClaims)["ID"].(string)

		repr, err := r.RepresentativeService.FindByID(representativeID)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		inviteLink, err := r.InviteLinkService.CreateRepresentativeInvite(r.Path, repr)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		json.NewEncoder(w).Encode(inviteLink)
	})
}

// AddProject adds a project to the company
func (r RepresentativeHandler) AddProject() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			return
		}

		cookie, _ := req.Cookie("token")
		token, _ := r.AuthHandler.GetToken(cookie)
		reprID := token.Claims.(jwt.MapClaims)["ID"].(string)

		var data ProjectData
		err := json.NewDecoder(req.Body).Decode(&data)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		repr, err := r.RepresentativeService.FindByID(reprID)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		project, err := repr.CreateProject(
			data.Description,
			data.Compensation,
			data.Duration,
			data.Recommendations,
		)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = r.RepresentativeService.CompanyService.AddProject(project)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		json.NewEncoder(w).Encode(project.ID)
	})
}

// Register registers all representative related handlers
func (r RepresentativeHandler) Register() {
	http.Handle(r.Path, LoggingHandler(os.Stdout, r.AuthHandler.Validate("", r.FetchRepresentativeByID())))
	http.Handle("/signup"+r.Path+"invite/", LoggingHandler(os.Stdout, r.SignupRepresentative()))
	http.Handle(r.Path+"invitelink/", LoggingHandler(os.Stdout, r.AuthHandler.Validate(domain.RepresentativeRole, r.MakeInviteLink())))
	http.Handle(r.Path+"invitations/", LoggingHandler(os.Stdout, r.AuthHandler.Validate(domain.RepresentativeRole, r.FetchCreatedInvitations())))
	http.Handle(r.Path+"projects/", LoggingHandler(os.Stdout, r.AuthHandler.Validate(domain.RepresentativeRole, r.AddProject())))
}
