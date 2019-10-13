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
	return nil
}

// FindByID ...
func (s StudentRepo) FindByID(id string) (d.Student, error) {

	return d.Student{}, nil
}

// FindAll ...
func (s StudentRepo) FindAll() ([]d.Student, error) {
	return []d.Student{}, nil
}
