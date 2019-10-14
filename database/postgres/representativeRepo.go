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
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}

	const insertQuery = `INSERT INTO "Representative"(representative_id, job_title, ref_user, ref_company)
	VALUES ($1, $2, $3, $4);`
	_, err = tx.Exec(insertQuery,
		repr.ID,
		repr.JobTitle,
		repr.User.ID,
		repr.CompanyID,
	)

	if err != nil {
		_ = tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// FindByID ...
func (r RepresentativeRepo) FindByID(id string) (d.Representative, error) {
	tx, err := r.DB.Begin()
	if err != nil {
		return d.Representative{}, err
	}

	var rID, title, cID, uID, fname, lname, email, hash, role string
	const selectQuery = `SELECT r.representative_id, r.job_title, r.ref_company,
	u.user_id, u.first_name, u.last_name, u.email, u.hashed_password, u.role
	FROM "Representative" r JOIN "User" u ON r.ref_user = u.user_id 
	WHERE r.representative_id = $1;`
	result := tx.QueryRow(selectQuery, id)
	err = result.Scan(&rID, &title, &cID, &uID, &fname, &lname, &email, &hash, &role)
	if err != nil {
		_ = tx.Rollback()
		return d.Representative{}, err
	}

	representative := d.Representative{
		ID:        rID,
		JobTitle:  title,
		CompanyID: cID,
		User: d.User{
			ID:             uID,
			FirstName:      fname,
			LastName:       lname,
			Email:          email,
			HashedPassword: hash,
			Role:           role,
		},
	}

	return representative, nil
}
