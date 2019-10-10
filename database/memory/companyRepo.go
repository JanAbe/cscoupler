package memory

import (
	"errors"

	"github.com/janabe/cscoupler/domain"
)

// CompanyRepo ..
type CompanyRepo struct {
	DB map[string]domain.Company
}

// Create ...
func (c CompanyRepo) Create(company domain.Company) error {
	c.DB[company.ID] = company
	return nil
}

// FindByID ...
func (c CompanyRepo) FindByID(id string) (domain.Company, error) {
	if company, ok := c.DB[id]; ok {
		return company, nil
	}

	return domain.Company{}, errors.New("no company with id: " + id)
}

// FindByName ...
func (c CompanyRepo) FindByName(name string) (domain.Company, error) {
	for _, company := range c.DB {
		if company.Name == name {
			return company, nil
		}
	}

	return domain.Company{}, errors.New("no company with name: " + name)
}

// FindAll ...
func (c CompanyRepo) FindAll() ([]domain.Company, error) {
	companies := []domain.Company{}
	for _, company := range c.DB {
		companies = append(companies, company)
	}

	return companies, nil
}
