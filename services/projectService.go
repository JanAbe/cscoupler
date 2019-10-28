package services

import (
	"github.com/janabe/cscoupler/domain"
)

// ProjectService struct, containing all features
// the app support regarding just projects
type ProjectService struct {
	ProjectRepo domain.ProjectRepository
}

// FetchAll fetches all projects
func (p ProjectService) FetchAll() ([]domain.Project, error) {
	projects, err := p.ProjectRepo.FindAll()
	if err != nil {
		return []domain.Project{}, err
	}

	return projects, nil
}
