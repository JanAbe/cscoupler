package postgres

import (
	"database/sql"

	d "github.com/janabe/cscoupler/domain"
)

// CompanyRepo ...
type CompanyRepo struct {
	DB *sql.DB
}

// Create ...
func (c CompanyRepo) Create(company d.Company) error {
	return nil
}

// FindByID ...
func (c CompanyRepo) FindByID(id string) (d.Company, error) {
	return d.Company{}, nil
}

// FindByName ...
func (c CompanyRepo) FindByName(name string) (d.Company, error) {
	return d.Company{}, nil
}

// FindAll ...
func (c CompanyRepo) FindAll() ([]d.Company, error) {
	return []d.Company{}, nil
}
