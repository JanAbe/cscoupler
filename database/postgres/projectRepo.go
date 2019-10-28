package postgres

import (
	"database/sql"

	"github.com/janabe/cscoupler/domain"
	"github.com/lib/pq"
)

// ProjectRepo struct for postgres database
type ProjectRepo struct {
	DB *sql.DB
}

// FindAll finds all the projects in the database
func (p ProjectRepo) FindAll() ([]domain.Project, error) {
	tx, err := p.DB.Begin()
	if err != nil {
		return []domain.Project{}, err
	}

	const selectQuery = `
	SELECT project_id, description, duration, compensation, recommendations, ref_company 
	FROM "Project";
	`

	rows, err := tx.Query(selectQuery)
	if err != nil {
		_ = tx.Rollback()
		return []domain.Project{}, err
	}

	projects := []domain.Project{}
	for rows.Next() {
		var (
			pID, descr, comp, dur, cID string
			recomms                    []string
		)

		if err := rows.Scan(&pID, &descr, &dur, &comp, pq.Array(&recomms), &cID); err != nil {
			_ = tx.Rollback()
			return []domain.Project{}, err
		}

		projects = append(projects, domain.Project{
			ID:              pID,
			Description:     descr,
			Compensation:    comp,
			Duration:        dur,
			Recommendations: recomms,
			CompanyID:       cID,
		})
	}

	err = tx.Commit()
	if err != nil {
		return []domain.Project{}, err
	}

	return projects, nil
}
