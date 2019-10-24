package postgres

import (
	"database/sql"

	d "github.com/janabe/cscoupler/domain"
	e "github.com/janabe/cscoupler/errors"
	"github.com/lib/pq"
)

// UserRepo struct for postgres db
type UserRepo struct {
	DB *sql.DB
}

// Create inserts a user in the DB. It should be used as a single
// unit of work, as it has its own transaction inside.
func (u UserRepo) Create(user d.User) error {
	tx, err := u.DB.Begin()
	if err != nil {
		return err
	}

	err = u.CreateTx(tx, user)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// FindRoleID finds the id of the role the user has
// So if the user is a student, FindRoleID will find the id
// of the student that is associated with the provided user
// account. It should be used as a single unit of work,
// as it has its own transaction inside
func (u UserRepo) FindRoleID(user d.User) (string, error) {
	tx, err := u.DB.Begin()
	if err != nil {
		return "", err
	}

	var roleID string

	if user.Role == d.StudentRole {
		const query = `SELECT student_id FROM "Student" WHERE ref_user=$1;`
		result := tx.QueryRow(query, user.ID)
		err = result.Scan(&roleID)
		if err != nil {
			_ = tx.Rollback()
			return "", err
		}
	} else if user.Role == d.RepresentativeRole {
		const query = `SELECT representative_id FROM "Representative" WHERE ref_user=$1;`
		result := tx.QueryRow(query, user.ID)
		err = result.Scan(&roleID)
		if err != nil {
			_ = tx.Rollback()
			return "", err
		}
	}

	err = tx.Commit()
	if err != nil {
		return "", err
	}

	return roleID, nil
}

// FindByID finds a user in the DB based on id. It should be used as a single
// unit of work, as it has its own transaction inside.
func (u UserRepo) FindByID(id string) (d.User, error) {
	tx, err := u.DB.Begin()
	if err != nil {
		return d.User{}, err
	}

	user, err := u.FindByIDTx(tx, id)
	if err != nil {
		return d.User{}, err
	}

	err = tx.Commit()
	if err != nil {
		return d.User{}, err
	}

	return user, nil
}

// FindByEmail finds a user in the DB based on email. It should be used as a single
// unit of work, as it has its own transaction inside.
func (u UserRepo) FindByEmail(email string) (d.User, error) {
	tx, err := u.DB.Begin()
	if err != nil {
		return d.User{}, err
	}

	user, err := u.FindByEmailTx(tx, email)
	if err != nil {
		return d.User{}, err
	}

	err = tx.Commit()
	if err != nil {
		return d.User{}, err
	}

	return user, nil
}

// CreateTx inserts a user in the DB. It should be used as PART of a
// unit of work, as a transaction gets passed in but will not be committed.
// This is the responsibility of the caller.
// It will rollback and return an error if something goes wrong
func (u UserRepo) CreateTx(tx *sql.Tx, user d.User) error {
	const insertQuery = `INSERT INTO "User"(user_id, first_name, last_name, email, 
		hashed_password, role) VALUES($1, $2, $3, $4, $5, $6);`
	_, err := tx.Exec(insertQuery,
		user.ID,
		user.FirstName,
		user.LastName,
		user.Email,
		user.HashedPassword,
		user.Role,
	)

	if err != nil {
		if err.(*pq.Error).Code.Name() == "unique_violation" {
			_ = tx.Rollback()
			return e.ErrorEmailAlreadyUsed
		}
		_ = tx.Rollback()
		return err
	}

	return nil
}

// FindByIDTx finds a user in the DB based on id. It should be used as PART of a
// unit of work, as a transaction gets passed in but will not be committed.
// This is the responsibility of the caller.
// It will rollback and return an error if something goes wrong
func (u UserRepo) FindByIDTx(tx *sql.Tx, id string) (d.User, error) {
	var uID, fname, lname, email, hash, role string
	const selectQuery = `SELECT user_id, first_name, last_name, email, role FROM "User" WHERE user_id = $1;`
	result := tx.QueryRow(selectQuery, id)

	err := result.Scan(&uID, &fname, &lname, &email, &hash, &role)
	if err != nil {
		_ = tx.Rollback()
		return d.User{}, err
	}

	return d.User{
		ID:             uID,
		Email:          email,
		HashedPassword: hash,
		FirstName:      fname,
		LastName:       lname,
		Role:           role,
	}, nil
}

// FindByEmailTx finds a user in the DB based on email. It should be used as PART of a
// unit of work, as a transaction gets passed in but will not be committed.
// This is the responsibility of the caller.
// It will rollback and return an error if something goes wrong
func (u UserRepo) FindByEmailTx(tx *sql.Tx, email string) (d.User, error) {
	var uID, fname, lname, uEmail, hash, role string
	const selectQuery = `SELECT user_id, first_name, last_name, email, hashed_password, role FROM "User" WHERE email=$1;`
	result := tx.QueryRow(selectQuery, email)

	err := result.Scan(&uID, &fname, &lname, &uEmail, &hash, &role)
	if err != nil {
		_ = tx.Rollback()
		return d.User{}, err
	}

	return d.User{
		ID:             uID,
		Email:          uEmail,
		HashedPassword: hash,
		FirstName:      fname,
		LastName:       lname,
		Role:           role,
	}, nil
}
