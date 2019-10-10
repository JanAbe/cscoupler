package services

import (
	"fmt"

	"github.com/janabe/cscoupler/domain"
	e "github.com/janabe/cscoupler/errors"
)

// CompanyService struct, containing all features
// the app supports regarding companies
type CompanyService struct {
	CompanyRepo           domain.CompanyRepository
	RepresentativeService *RepresentativeService
}

// Register registers a new company and their main representative
func (c CompanyService) Register(company domain.Company) error {

	// !!!!!
	// todo: bekijk waar ik wil controleren of een bedrijf/student
	// zijn naam/email al geascocieerd is met een account
	// Momenteel doe ik dit in de handlers, maar is het niet logischer
	// en beter om dat in de services te doen?
	// moet ook kijken naar studentHandler, userHandler + de services
	// het probleem dat ik alleen heb, is dat als ik de checks hier wil uitvoeren
	// (wat beter is vgm), ik niet weet hoe ik specifieke http.Statussen kan
	// returnen in de handler.
	// nvm found an answer, create custom error values.
	// compare returned errors with these erros to determine what went wrong
	if c.NameAlreadyUsed(company.Name) {
		return e.ErrorCompanyNameAlreadyUsed
	}

	err := c.CompanyRepo.Create(company)
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = c.RepresentativeService.Register(company.Representatives[0])
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// FindByID finds a company based on ID
func (c CompanyService) FindByID(id string) (domain.Company, error) {
	company, err := c.CompanyRepo.FindByID(id)
	if err != nil {
		return domain.Company{}, err
	}

	return company, nil
}

// Exists checks if a company exists with the provided id
func (c CompanyService) Exists(id string) bool {
	_, err := c.FindByID(id)
	if err != nil {
		return false
	}

	return true
}

// NameAlreadyUsed checks if a company name already exists or not
func (c CompanyService) NameAlreadyUsed(name string) bool {
	_, err := c.CompanyRepo.FindByName(name)
	if err != nil {
		return false
	}

	return true
}
