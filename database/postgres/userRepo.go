package postgres

import (
	"database/sql"

	d "github.com/janabe/cscoupler/domain"
)

// UserRepo ..
type UserRepo struct {
	DB *sql.DB
}

// Create ...
func (u UserRepo) Create(user d.User) error {
	tx, err := u.DB.Begin()
	if err != nil {
		return err
	}

	const insertQuery = `INSERT INTO "User"(user_id, first_name, last_name, email, 
		hashed_password, role) VALUES($1, $2, $3, $4, $5, $6);`
	_, err = tx.Exec(insertQuery,
		user.ID,
		user.FirstName,
		user.LastName,
		user.Email,
		user.HashedPassword,
		user.Role,
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
func (u UserRepo) FindByID(id string) (d.User, error) {
	tx, err := u.DB.Begin()
	if err != nil {
		return d.User{}, err
	}

	var uID, fname, lname, email, hash, role string
	const selectQuery = `SELECT user_id, first_name, last_name, email, role FROM "User" WHERE user_id = $1;`
	result := tx.QueryRow(selectQuery, id)

	err = result.Scan(&uID, &fname, &lname, &email, &hash, &role)
	if err != nil {
		_ = tx.Rollback()
		return d.User{}, err
	}

	user := d.User{
		ID:             uID,
		Email:          email,
		HashedPassword: hash,
		FirstName:      fname,
		LastName:       lname,
		Role:           role,
	}

	err = tx.Commit()
	if err != nil {
		return d.User{}, err
	}

	return user, nil
}

// FindByEmail ...
func (u UserRepo) FindByEmail(email string) (d.User, error) {
	tx, err := u.DB.Begin()
	if err != nil {
		return d.User{}, err
	}

	var uID, fname, lname, uEmail, hash, role string
	const selectQuery = `SELECT user_id, first_name, last_name, email, hashed_password, role FROM "User" WHERE email = $1;`
	result := tx.QueryRow(selectQuery, email)

	err = result.Scan(&uID, &fname, &lname, &uEmail, &hash, &role)
	if err != nil {
		_ = tx.Rollback()
		return d.User{}, err
	}

	user := d.User{
		ID:             uID,
		Email:          uEmail,
		HashedPassword: hash,
		FirstName:      fname,
		LastName:       lname,
		Role:           role,
	}

	err = tx.Commit()
	if err != nil {
		return d.User{}, err
	}

	return user, nil
}
