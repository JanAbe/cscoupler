package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/janabe/cscoupler/domain"
	e "github.com/janabe/cscoupler/errors"
	"github.com/janabe/cscoupler/services"
)

// RepresentativeHandler struct containing all
// representative related handler funcs
type RepresentativeHandler struct {
	RepresentativeService services.RepresentativeService
	AuthHandler           AuthHandler
	Path                  string
}

// RepresentativeData is a struct that corresponds to incoming
// representative data
type RepresentativeData struct {
	Position  string   `json:"position"`
	CompanyID string   `json:"companyID"`
	UserData  UserData `json:"user"`
}

// SignupRepresentative signs up a representative and binds
// it to the companyID present in the invite-link URL
// Format for invite-links: /signup/representatives/invite/[companyID]
func (r RepresentativeHandler) SignupRepresentative() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			return
		}

		companyID := strings.TrimPrefix(req.URL.Path, "/signup"+r.Path+"invite/")
		var data RepresentativeData

		// check if json is invalid
		err := json.NewDecoder(req.Body).Decode(&data)
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
		)

		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		representative, err := domain.NewRepresentative(
			data.Position,
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

// RegisterHandlers registers all representative related handlers
func (r RepresentativeHandler) RegisterHandlers() {
	http.Handle(r.Path, LoggingHandler(os.Stdout, r.AuthHandler.Validate(r.FetchRepresentativeByID())))
	http.Handle("/signup"+r.Path+"invite/", LoggingHandler(os.Stdout, r.SignupRepresentative()))
}
