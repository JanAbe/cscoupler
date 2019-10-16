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

// CompanyHandler struct containing all
// company related handler funcs
type CompanyHandler struct {
	CompanyService services.CompanyService
	AuthHandler    AuthHandler
	Path           string
}

// CompanyData is a struct that corresponds to incoming company data
type CompanyData struct {
	Name            string               `json:"name"`
	Information     string               `json:"information"`
	Locations       []LocationData       `json:"locations"`
	Representatives []RepresentativeData `json:"representatives"`
	Projects        []ProjectData        `json:"projects"`
}

// LocationData is a struct that corresponds to incoming location data
// of companies
type LocationData struct {
	Street  string `json:"street"`
	Zipcode string `json:"zipcode"`
	City    string `json:"city"`
	Number  string `json:"number"`
}

// ProjectData is a struct that corresponds to incoming project data
type ProjectData struct {
	ID              string   `json:"id"`
	Description     string   `json:"description"`
	Compensation    string   `json:"compensation"`
	Duration        string   `json:"duration"`
	Recommendations []string `json:"recommendations"`
	CompanyID       string   `json:"companyID"`
}

// SignupCompany signs up a company and the main representative
// of this company
func (c CompanyHandler) SignupCompany() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			return
		}

		var data CompanyData

		// check if json is invalid
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		company, err := domain.NewCompany(data.Name, data.Information)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		for _, l := range data.Locations {
			location, err := domain.NewAddress(l.Street, l.Zipcode, l.City, l.Number)
			if err != nil {
				fmt.Println(err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			company.Locations = append(
				company.Locations,
				location,
			)
		}

		// There should be 1 representative sent
		// when creating a company, the main representative.
		if len(data.Representatives) != 1 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		mainRepresentative := data.Representatives[0]
		user, err := domain.NewUser(
			mainRepresentative.UserData.Email,
			mainRepresentative.UserData.Password,
			mainRepresentative.UserData.Firstname,
			mainRepresentative.UserData.Lastname,
			domain.RepresentativeRole,
		)

		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		representative, err := domain.NewRepresentative(
			mainRepresentative.JobTitle,
			company.ID,
			user,
		)

		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		company.Representatives = append(company.Representatives, representative)

		err = c.CompanyService.Register(company)
		if err == e.ErrorEmailAlreadyUsed || err == e.ErrorCompanyNameAlreadyUsed {
			fmt.Println(err)
			w.WriteHeader(http.StatusConflict)
			return
		}

		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// todo: how does a representative gets his/her own id
		fmt.Println(representative.ID)
		json.NewEncoder(w).Encode(company.ID)
	})
}

// FetchCompanyByID fetches a company based on ID
// path = /companies/... where the dots are a company ID
func (c CompanyHandler) FetchCompanyByID() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			return
		}

		id := strings.TrimPrefix(r.URL.Path, c.Path)
		company, err := c.CompanyService.FindByID(id)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		companyData := ToCompanyData(company)

		json.NewEncoder(w).Encode(companyData)
	})
}

// Register registers all company related handlers
func (c CompanyHandler) Register() {
	http.Handle(c.Path, LoggingHandler(os.Stdout, c.AuthHandler.Validate("", c.FetchCompanyByID())))
	http.Handle("/signup/company", LoggingHandler(os.Stdout, c.SignupCompany()))
}
