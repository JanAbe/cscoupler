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

// FindByID finds a project in the DB based on id. It should be used as a single
// unit of work, as it has its own transaction inside.
func (p ProjectRepo) FindByID(id string) (domain.Project, error) {
	tx, err := p.DB.Begin()
	if err != nil {
		return domain.Project{}, err
	}

	project, err := p.FindByIDTx(tx, id)
	if err != nil {
		return domain.Project{}, err
	}

	err = tx.Commit()
	if err != nil {
		return domain.Project{}, err
	}

	return project, nil
}

// Delete deletes a project in the DB based on id. It should be used as a single
// unit of work, as it has its own transaction inside.
func (p ProjectRepo) Delete(id string) error {
	tx, err := p.DB.Begin()
	if err != nil {
		return err
	}

	err = p.DeleteTx(tx, id)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// FindByIDTx finds a project in the DB based on id. It should be used as PART of a
// unit of work, as a transaction gets passed in but will not be committed.
// This is the responsibility of the caller.
// It will rollback and return an error if something goes wrong
func (p ProjectRepo) FindByIDTx(tx *sql.Tx, id string) (domain.Project, error) {
	var (
		pID, descr, comp, dur, cID string
		recomms                    []string
	)

	const selectQuery = `
	SELECT project_id, description, duration, compensation, recommendations, ref_company 
	FROM "Project" WHERE project_id=$1;
	`
	result := tx.QueryRow(selectQuery, id)
	err := result.Scan(&pID, &descr, &dur, &comp, pq.Array(&recomms), &cID)
	if err != nil {
		_ = tx.Rollback()
		return domain.Project{}, err
	}

	return domain.Project{
		ID:              pID,
		Description:     descr,
		Compensation:    comp,
		Duration:        dur,
		Recommendations: recomms,
		CompanyID:       cID,
	}, nil
}

// DeleteTx deletes a project in the DB based on id. It should be used as PART of a
// unit of work, as a transaction gets passed in but will not be committed.
// This is the responsibility of the caller.
// It will rollback and return an error if something goes wrong
func (p ProjectRepo) DeleteTx(tx *sql.Tx, id string) error {
	const deleteQuery = `DELETE FROM "Project" WHERE project_id=$1;`
	_, err := tx.Exec(deleteQuery, id)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	return nil
}

// FindAll finds all the projects in the database
func (p ProjectRepo) FindAll() ([]domain.Project, error) {
	tx, err := p.DB.Begin()
	if err != nil {
		return []domain.Project{}, err
	}

	const selectQuery = `
	SELECT project_id, description, duration, compensation, recommendations, ref_company 
	FROM "Project" ORDER BY RANDOM();
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
