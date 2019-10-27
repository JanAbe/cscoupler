package postgres

import (
	"database/sql"

	"github.com/lib/pq"

	d "github.com/janabe/cscoupler/domain"
)

// StudentRepo struct for postgres database
type StudentRepo struct {
	DB       *sql.DB
	UserRepo UserRepo
}

// Create inserts a student in the DB. It should be used as a single
// unit of work, as it has its own transaction inside.
func (s StudentRepo) Create(student d.Student) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}

	err = s.CreateTx(tx, student)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// Update updates a student in the DB. It should be used as a single
// unit of work, as it has its own transaction inside.
func (s StudentRepo) Update(student d.Student) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}

	err = s.UpdateTx(tx, student)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// FindByID finds a student in the DB based on id. It should be used as a single
// unit of work, as it has its own transaction inside.
func (s StudentRepo) FindByID(id string) (d.Student, error) {
	tx, err := s.DB.Begin()
	if err != nil {
		return d.Student{}, err
	}

	student, err := s.FindByIDTx(tx, id)
	if err != nil {
		return d.Student{}, err
	}

	err = tx.Commit()
	if err != nil {
		return d.Student{}, err
	}

	return student, nil
}

// FindAll finds all the students in the DB based. It should be used as a single
// unit of work, as it has its own transaction inside.
func (s StudentRepo) FindAll() ([]d.Student, error) {
	tx, err := s.DB.Begin()
	if err != nil {
		return []d.Student{}, err
	}

	const selectQuery = `SELECT s.student_id, s.university, s.skills, s.experiences, s.short_experiences, 
	s.wishes, s.status, s.resume, u.user_id, u.first_name, u.last_name, u.email, u.role 
	FROM "Student" s JOIN "User" u ON s.ref_user = u.user_id`

	rows, err := tx.Query(selectQuery)
	if err != nil {
		_ = tx.Rollback()
		return []d.Student{}, err
	}
	defer rows.Close()

	students := []d.Student{}
	for rows.Next() {
		var (
			sID, uni, resume, wishes              string
			uID, fname, lname, email, role        string
			skills, experiences, shortExperiences []string
			status                                d.Status
		)

		if err := rows.Scan(&sID, &uni, pq.Array(&skills),
			pq.Array(&experiences), pq.Array(&shortExperiences), &wishes,
			&status, &resume, &uID, &fname, &lname, &email, &role); err != nil {
			_ = tx.Rollback()
			return []d.Student{}, err
		}

		students = append(students, d.Student{
			ID:               sID,
			University:       uni,
			Skills:           skills,
			Experiences:      experiences,
			ShortExperiences: shortExperiences,
			Wishes:           wishes,
			Status:           status,
			Resume:           resume,
			User: d.User{
				ID:        uID,
				Email:     email,
				FirstName: fname,
				LastName:  lname,
				Role:      role,
			},
		})
	}

	err = tx.Commit()
	if err != nil {
		return []d.Student{}, err
	}

	return students, nil
}

// CreateTx inserts a student in the DB. It should be used as PART of a
// unit of work, as a transaction gets passed in but will not be committed.
// This is the responsibility of the caller.
// It will rollback and return an error if something goes wrong
func (s StudentRepo) CreateTx(tx *sql.Tx, student d.Student) error {
	err := s.UserRepo.CreateTx(tx, student.User)
	if err != nil {
		return err
	}

	const insertQuery = `INSERT INTO "Student"(student_id, university, skills, experiences, short_experiences, wishes, status, resume, ref_user) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);`
	_, err = tx.Exec(insertQuery,
		student.ID,
		student.University,
		pq.Array(student.Skills),
		pq.Array(student.Experiences),
		pq.Array(student.ShortExperiences),
		student.Wishes,
		student.Status,
		student.Resume,
		student.User.ID,
	)

	if err != nil {
		_ = tx.Rollback()
		return err
	}

	return nil
}

// UpdateTx udpates a student in the DB. It should be used as PART of a
// unit of work, as a transaction gets passed in but will not be committed.
// This is the responsibility of the caller.
// It will rollback and return an error if something goes wrong
func (s StudentRepo) UpdateTx(tx *sql.Tx, student d.Student) error {
	const updateStudentQuery = `UPDATE "Student" s 
	SET university=$1, skills=$2, experiences=$3, short_experiences=$4, wishes=$5, status=$6, resume=$7 WHERE s.student_id=$8;`
	_, err := tx.Exec(updateStudentQuery,
		student.University,
		pq.Array(student.Skills),
		pq.Array(student.Experiences),
		pq.Array(student.ShortExperiences),
		student.Wishes,
		student.Status,
		student.Resume,
		student.ID,
	)

	if err != nil {
		_ = tx.Rollback()
		return err
	}

	// todo: move this code to UserRepo and call it Update()
	const updateUserQuery = `UPDATE "User" u SET first_name=$1, last_name=$2, email=$3
	WHERE u.user_id=(SELECT ref_user FROM "Student" WHERE student_id=$4);`
	_, err = tx.Exec(updateUserQuery,
		student.User.FirstName,
		student.User.LastName,
		student.User.Email,
		student.ID,
	)

	if err != nil {
		_ = tx.Rollback()
		return err
	}

	return nil
}

// FindByIDTx finds a student in the DB based on id. It should be used as PART of a
// unit of work, as a transaction gets passed in but will not be committed.
// This is the responsibility of the caller.
// It will rollback and return an error if something goes wrong
func (s StudentRepo) FindByIDTx(tx *sql.Tx, id string) (d.Student, error) {
	var uID, fname, lname, email, role, resume string
	var sID, uni, wishes string
	var skills, exp, shortExp []string
	var status d.Status

	const selectQuery = `SELECT student_id, s.university, s.skills, s.experiences, s.short_experiences, s.wishes, s.status, s.resume,
	user_id, u.first_name, u.last_name, u.email, u.role FROM "Student" s JOIN "User" u ON s.ref_user = u.user_id
	WHERE student_id=$1;`
	result := tx.QueryRow(selectQuery, id)

	err := result.Scan(&sID, &uni, pq.Array(&skills), pq.Array(&exp), pq.Array(&shortExp), &wishes, &status, &resume, &uID, &fname, &lname, &email, &role)
	if err != nil {
		_ = tx.Rollback()
		return d.Student{}, err
	}

	return d.Student{
		ID:               sID,
		University:       uni,
		Skills:           skills,
		Experiences:      exp,
		ShortExperiences: shortExp,
		Wishes:           wishes,
		Status:           status,
		Resume:           resume,
		User: d.User{
			ID:        uID,
			Email:     email,
			FirstName: fname,
			LastName:  lname,
			Role:      role,
		},
	}, nil
}
