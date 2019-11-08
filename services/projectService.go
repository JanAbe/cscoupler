package services

import (
	"github.com/janabe/cscoupler/domain"
)

// ProjectService struct, containing all features
// the app support regarding just projects
type ProjectService struct {
	ProjectRepo domain.ProjectRepository
}

// FindByID finds a project by ID
func (p ProjectService) FindByID(id string) (domain.Project, error) {
	project, err := p.ProjectRepo.FindByID(id)
	if err != nil {
		return domain.Project{}, nil
	}

	return project, nil
}

// Delete deletes a project
func (p ProjectService) Delete(id string) error {
	err := p.ProjectRepo.Delete(id)
	if err != nil {
		return err
	}

	return nil
}

// FetchAll fetches all projects
func (p ProjectService) FetchAll() ([]domain.Project, error) {
	projects, err := p.ProjectRepo.FindAll()
	if err != nil {
		return []domain.Project{}, err
	}

	return projects, nil
}
