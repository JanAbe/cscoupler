package postgres

import (
	"database/sql"

	d "github.com/janabe/cscoupler/domain"
)

// RepresentativeRepo ...
type RepresentativeRepo struct {
	DB *sql.DB
}

// Create ...
func (r RepresentativeRepo) Create(repr d.Representative) error {
	return nil
}

// FindByID ...
func (r RepresentativeRepo) FindByID(id string) (d.Representative, error) {
	return d.Representative{}, nil
}
