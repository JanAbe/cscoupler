package services

import (
	"github.com/janabe/cscoupler/domain"
	e "github.com/janabe/cscoupler/errors"
)

// RepresentativeService struct, containing all features
// the app supports regarding representatives
type RepresentativeService struct {
	RepresentativeRepo domain.RepresentativeRepository
	CompanyService     CompanyService
	UserService        UserService
}

// Register registers a representive with the provided data
func (r RepresentativeService) Register(representative domain.Representative) error {
	if !r.CompanyService.Exists(representative.CompanyID) {
		return e.ErrorEntityNotFound
	}

	err := r.RepresentativeRepo.Create(representative)
	if err != nil {
		return err
	}

	return nil
}

// FindByID finds a representative based on id
func (r RepresentativeService) FindByID(id string) (domain.Representative, error) {
	repr, err := r.RepresentativeRepo.FindByID(id)
	if err != nil {
		return domain.Representative{}, err
	}

	return repr, nil
}
