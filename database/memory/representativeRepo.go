package memory

import (
	"errors"

	"github.com/janabe/cscoupler/domain"
)

// RepresentativeRepo ...
type RepresentativeRepo struct {
	DB map[string]domain.Representative
}

// Create ...
func (r RepresentativeRepo) Create(repr domain.Representative) error {
	r.DB[repr.ID] = repr
	return nil
}

// FindByID ...
func (r RepresentativeRepo) FindByID(id string) (domain.Representative, error) {
	if repr, ok := r.DB[id]; ok {
		return repr, nil
	}

	return domain.Representative{}, errors.New("no representative with id: " + id)
}
