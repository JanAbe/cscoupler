package postgres

import (
	"database/sql"

	"github.com/lib/pq"

	d "github.com/janabe/cscoupler/domain"
)

// StudentRepo ...
type StudentRepo struct {
	DB *sql.DB
}

// Create ...
func (s StudentRepo) Create(student d.Student) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}

	const insertQuery = `INSERT INTO "Student"(student_id, university, skills, experience, status, ref_user) VALUES ($1, $2, $3, $4, $5, $6);`
	_, err = tx.Exec(insertQuery,
		student.ID,
		student.University,
		pq.Array(student.Skills),
		pq.Array(student.Experience),
		student.Status,
		student.User.ID,
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

// Update ...
func (s StudentRepo) Update(student d.Student) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}

	const updateStudentQuery = `UPDATE "Student" s SET university=$1, skills=$2, experience=$3, 
	status=$4 WHERE s.student_id=$5;`
	_, err = tx.Exec(updateStudentQuery,
		student.University,
		pq.Array(student.Skills),
		pq.Array(student.Experience),
		student.Status,
		student.ID,
	)

	if err != nil {
		_ = tx.Rollback()
		return err
	}

	const updateUserQuery = `UPDATE "User" u SET first_name=$1, last_name=$2, email=$3
	WHERE u.user_id=$4;`
	_, err = tx.Exec(updateUserQuery,
		student.User.FirstName,
		student.User.LastName,
		student.User.Email,
		student.User.ID,
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
func (s StudentRepo) FindByID(id string) (d.Student, error) {
	tx, err := s.DB.Begin()
	if err != nil {
		return d.Student{}, err
	}

	var uID, fname, lname, email, role string
	var sID, uni string
	var skills, exp []string
	var status d.Status

	const selectQuery = `SELECT student_id, s.university, s.skills, s.experience, s.status, 
	user_id, u.first_name, u.last_name, u.email, u.role FROM "Student" s JOIN "User" u ON s.ref_user = u.user_id
	WHERE student_id=$1;`
	result := tx.QueryRow(selectQuery, id)

	err = result.Scan(&sID, &uni, pq.Array(&skills), pq.Array(&exp), &status, &uID, &fname, &lname, &email, &role)
	if err != nil {
		return d.Student{}, err
	}

	student := d.Student{
		ID:         sID,
		University: uni,
		Skills:     skills,
		Experience: exp,
		Status:     status,
		User: d.User{
			ID:        uID,
			Email:     email,
			FirstName: fname,
			LastName:  lname,
			Role:      role,
		},
	}

	return student, nil
}

// FindAll ...
func (s StudentRepo) FindAll() ([]d.Student, error) {
	return []d.Student{}, nil
}
