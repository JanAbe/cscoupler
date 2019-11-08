package postgres

import (
	"database/sql"

	d "github.com/janabe/cscoupler/domain"
)

// RepresentativeRepo struct for postgres database
type RepresentativeRepo struct {
	DB       *sql.DB
	UserRepo UserRepo
}

// Create inserts a representative in the DB. It should be used as a single
// unit of work, as it has its own transaction inside.
func (r RepresentativeRepo) Create(repr d.Representative) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}

	err = r.UserRepo.CreateTx(tx, repr.User)
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

// FindByID finds a representative in the DB based on id. It should be used as a single
// unit of work, as it has its own transaction inside.
func (r RepresentativeRepo) FindByID(id string) (d.Representative, error) {
	tx, err := r.DB.Begin()
	if err != nil {
		return d.Representative{}, err
	}

	representative, err := r.FindByIDTx(tx, id)
	if err != nil {
		return d.Representative{}, err
	}

	err = tx.Commit()
	if err != nil {
		return d.Representative{}, err
	}

	return representative, nil
}

// Update updates a representative in the DB. It should be used as a single unit of work,
// as it has its own transaction inside.
func (r RepresentativeRepo) Update(representative d.Representative) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}

	err = r.UpdateTx(tx, representative)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// CreateTx inserts a representative in the DB. It should be used as PART of a
// unit of work, as a transaction gets passed in but will not be committed.
// This is the responsibility of the caller.
// It will rollback and return an error if something goes wrong
func (r RepresentativeRepo) CreateTx(tx *sql.Tx, repr d.Representative) error {
	err := r.UserRepo.CreateTx(tx, repr.User)
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

	return nil
}

// UpdateTx updates a representative in the DB. It should be used as PART of
// a unit of work, as a transaction gets passed in but will not be commited.
// This is the responsibility of the caller.
// It will rollback and return an error if something goes wrong.
func (r RepresentativeRepo) UpdateTx(tx *sql.Tx, repr d.Representative) error {
	const updateRepresentativeQuery = `UPDATE "Representative" r
	SET job_title=$1 WHERE r.representative_id=$2;`

	_, err := tx.Exec(updateRepresentativeQuery,
		repr.JobTitle,
		repr.ID,
	)

	if err != nil {
		_ = tx.Rollback()
		return err
	}

	const updateUserQuery = `UPDATE "User" u SET first_name=$1, last_name=$2, email=$3
	WHERE u.user_id=(SELECT ref_user FROM "Student" WHERE student_id=$4);`
	_, err = tx.Exec(updateUserQuery,
		repr.User.FirstName,
		repr.User.LastName,
		repr.User.Email,
		repr.ID,
	)

	if err != nil {
		_ = tx.Rollback()
		return err
	}

	return nil
}

// FindByIDTx finds a representative in the DB based on id. It should be used as PART of a
// unit of work, as a transaction gets passed in but will not be committed.
// This is the responsibility of the caller.
// It will rollback and return an error if something goes wrong
func (r RepresentativeRepo) FindByIDTx(tx *sql.Tx, id string) (d.Representative, error) {
	var rID, title, cID, uID, fname, lname, email, hash, role string
	const selectQuery = `SELECT r.representative_id, r.job_title, r.ref_company,
	u.user_id, u.first_name, u.last_name, u.email, u.hashed_password, u.role
	FROM "Representative" r JOIN "User" u ON r.ref_user = u.user_id 
	WHERE r.representative_id = $1;`
	result := tx.QueryRow(selectQuery, id)
	err := result.Scan(&rID, &title, &cID, &uID, &fname, &lname, &email, &hash, &role)
	if err != nil {
		_ = tx.Rollback()
		return d.Representative{}, err
	}

	return d.Representative{
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
	}, nil
}
