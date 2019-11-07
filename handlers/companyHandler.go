package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
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
	ID              string               `json:"id"`
	Name            string               `json:"name"`
	Description     string               `json:"description"`
	Information     string               `json:"information"`
	Locations       []LocationData       `json:"locations"`
	Representatives []RepresentativeData `json:"representatives"`
	Projects        []ProjectData        `json:"projects"`
}

// LocationData is a struct that corresponds to incoming location data
// of companies
type LocationData struct {
	ID      string `json:"id"`
	Street  string `json:"street"`
	Zipcode string `json:"zipcode"`
	City    string `json:"city"`
	Number  string `json:"number"`
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

		company, err := domain.NewCompany(uuid.New().String(), data.Name, data.Information, data.Description)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		for _, l := range data.Locations {
			location, err := domain.NewAddress(uuid.New().String(), l.Street, l.Zipcode, l.City, l.Number)
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
			uuid.New().String(),
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

		json.NewEncoder(w).Encode(company.ID)
	})
}

// EditCompany edits the company account
func (c CompanyHandler) EditCompany() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			return
		}

		cookie, _ := r.Cookie("token")
		token, _ := c.AuthHandler.GetToken(cookie)
		reprID := token.Claims.(jwt.MapClaims)["ID"].(string)

		companyID := strings.TrimPrefix(r.URL.Path, c.Path+"edit/")
		company, err := c.CompanyService.FindByID(companyID)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		// Validate that the representative editing the company data
		// is an employee of this specific company
		var worksForCompany bool
		for _, r := range company.Representatives {
			if r.ID == reprID {
				worksForCompany = true
				break
			}
		}

		if !worksForCompany {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		var updatedCompanyData CompanyData
		err = json.NewDecoder(r.Body).Decode(&updatedCompanyData)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		updatedCompany, err := domain.NewCompany(companyID, updatedCompanyData.Name, updatedCompanyData.Information, updatedCompanyData.Description)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		for _, l := range updatedCompanyData.Locations {
			location, err := domain.NewAddress(l.ID, l.Street, l.Zipcode, l.City, l.Number)
			if err != nil {
				fmt.Println(err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			updatedCompany.Locations = append(
				updatedCompany.Locations,
				location,
			)
		}

		for _, p := range updatedCompanyData.Projects {
			project, err := domain.NewProject(p.Description, p.Compensation, p.Duration, p.CompanyID, p.Recommendations)
			if err != nil {
				fmt.Println(err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			updatedCompany.Projects = append(updatedCompany.Projects, project)
		}

		err = c.CompanyService.Edit(updatedCompany)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		json.NewEncoder(w).Encode(updatedCompany.ID)
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

// FetchCompanyNameByID fetches the name of the company
// based on ID
func (c CompanyHandler) FetchCompanyNameByID() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			return
		}

		id := strings.TrimPrefix(r.URL.Path, c.Path+"name/")
		company, err := c.CompanyService.FindByID(id)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		companyData := ToCompanyData(company)

		json.NewEncoder(w).Encode(companyData.Name)
	})
}

// FetchAllCompanies fetches all the companies
func (c CompanyHandler) FetchAllCompanies() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" && r.Method != "OPTIONS" {
			return
		}

		companies, err := c.CompanyService.FindAll()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Println(err)
			return
		}

		var companiesData []CompanyData
		for _, c := range companies {
			companiesData = append(companiesData, ToCompanyData(c))
		}

		json.NewEncoder(w).Encode(companiesData)
	})
}

// Register registers all company related handlers
func (c CompanyHandler) Register() {
	http.Handle(c.Path, LoggingHandler(os.Stdout, c.AuthHandler.Validate("", c.FetchCompanyByID())))
	http.Handle(c.Path+"all", LoggingHandler(os.Stdout, c.AuthHandler.Validate("", c.FetchAllCompanies())))
	http.Handle("/companies/name/", LoggingHandler(os.Stdout, c.FetchCompanyNameByID()))
	http.Handle("/signup/company", LoggingHandler(os.Stdout, c.SignupCompany()))
	http.Handle(c.Path+"edit/", LoggingHandler(os.Stdout, c.AuthHandler.Validate("representative", c.EditCompany())))
}
