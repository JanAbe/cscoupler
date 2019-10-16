package services

import (
	"github.com/janabe/cscoupler/domain"
	e "github.com/janabe/cscoupler/errors"
)

// CompanyService struct, containing all features
// the app supports regarding companies
type CompanyService struct {
	CompanyRepo domain.CompanyRepository
	ReprService *RepresentativeService
}

// Register registers a new company and their main representative
func (c CompanyService) Register(company domain.Company) error {
	if c.NameAlreadyUsed(company.Name) {
		return e.ErrorCompanyNameAlreadyUsed
	}

	if c.ReprService.UserService.EmailAlreadyUsed(
		company.Representatives[0].User.Email) {
		return e.ErrorEmailAlreadyUsed
	}

	err := c.CompanyRepo.Create(company)
	if err != nil {
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

// AddProject adds a project to the company
func (c CompanyService) AddProject(p domain.Project) error {
	err := c.CompanyRepo.AddProject(p)
	if err != nil {
		return err
	}

	return nil
}
