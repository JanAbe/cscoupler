package postgres

import (
	"database/sql"

	"github.com/lib/pq"

	d "github.com/janabe/cscoupler/domain"
)

// CompanyRepo struct for postgres database
type CompanyRepo struct {
	DB       *sql.DB
	ReprRepo RepresentativeRepo
}

// Create inserts a company in the DB. It should be used as a single
// unit of work, as it has its own transaction inside.
func (c CompanyRepo) Create(company d.Company) error {
	tx, err := c.DB.Begin()
	if err != nil {
		return err
	}

	const insertCompanyQuery = `INSERT INTO "Company"(company_id, name, information, description) VALUES ($1, $2, $3, $4);`
	_, err = tx.Exec(insertCompanyQuery, company.ID, company.Name, company.Information, company.Description)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	const insertAddressesQuery = `INSERT INTO "Address"(address_id, street, zipcode, city, number, ref_company) VALUES ($1, $2, $3, $4, $5, $6);`
	for _, l := range company.Locations {
		_, err = tx.Exec(insertAddressesQuery,
			l.ID,
			l.Street,
			l.Zipcode,
			l.City,
			l.Number,
			company.ID,
		)

		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	err = c.ReprRepo.CreateTx(tx, company.Representatives[0])
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// FindByID finds a company in the DB based on id. It should be used as a single
// unit of work, as it has its own transaction inside.
func (c CompanyRepo) FindByID(id string) (d.Company, error) {
	tx, err := c.DB.Begin()
	if err != nil {
		return d.Company{}, err
	}

	company, err := c.FindByIDTx(tx, id)
	if err != nil {
		return d.Company{}, err
	}

	err = tx.Commit()
	if err != nil {
		return d.Company{}, err
	}

	return company, nil
}

// FindByName finds a company in the DB based on name. It should be used as a single
// unit of work, as it has its own transaction inside.
func (c CompanyRepo) FindByName(name string) (d.Company, error) {
	tx, err := c.DB.Begin()
	if err != nil {
		return d.Company{}, err
	}

	company, err := c.FindByNameTx(tx, name)
	if err != nil {
		return d.Company{}, err
	}

	err = tx.Commit()
	if err != nil {
		return d.Company{}, err
	}

	return company, nil
}

// FindAll finds all companies in the DB. It should be used as a single
// unit of work, as it has its own transaction inside.
func (c CompanyRepo) FindAll() ([]d.Company, error) {
	tx, err := c.DB.Begin()
	if err != nil {
		return []d.Company{}, err
	}

	companies := []d.Company{}
	const selectIDSQuery = `SELECT company_id FROM "Company";`
	rows, err := tx.Query(selectIDSQuery)
	defer rows.Close()

	ids := []string{}
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			_ = tx.Rollback()
			return []d.Company{}, err
		}
		ids = append(ids, id)
	}

	// added the extra ids slice because it wasn't possible
	// to execute a new query inside the 'select all ids' query.
	for _, id := range ids {
		company, err := c.FindByIDTx(tx, id)
		if err != nil {
			return []d.Company{}, err
		}
		companies = append(companies, company)
	}

	err = tx.Commit()
	if err != nil {
		return []d.Company{}, err
	}

	return companies, nil
}

// AddProject adds a project to the company in the db. It should be used as a
// single unit of work, as it has its own transaction inside.
func (c CompanyRepo) AddProject(p d.Project) error {
	tx, err := c.DB.Begin()
	if err != nil {
		return err
	}

	const insertQuery = `INSERT INTO "Project"(project_id, description, 
	compensation, duration, recommendations, ref_company)
	VALUES($1, $2, $3, $4, $5, $6);`

	_, err = tx.Exec(insertQuery,
		p.ID,
		p.Description,
		p.Compensation,
		p.Duration,
		pq.Array(p.Recommendations),
		p.CompanyID,
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

// Update updates a company in the DB. It should be used as a single unit of work
// as it has its own transaction inside.
func (c CompanyRepo) Update(company d.Company) error {
	tx, err := c.DB.Begin()
	if err != nil {
		return err
	}

	err = c.UpdateTx(tx, company)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// UpdateTx updates a company in the DB. It should be used as PART of
// a unit of work, as a transaction gets passed in but will not be commited.
// This is the responsibility of the caller.
// It will rollback and return an error if something goes wrong.
func (c CompanyRepo) UpdateTx(tx *sql.Tx, company d.Company) error {
	const updateCompanyQuery = `UPDATE "Company" c
	SET name=$1, description=$2, information=$3 WHERE c.company_id=$4;`

	_, err := tx.Exec(updateCompanyQuery,
		company.Name,
		company.Description,
		company.Information,
		company.ID,
	)

	if err != nil {
		_ = tx.Rollback()
		return err
	}

	const updateLocationsQuery = `UPDATE "Address" a
	SET street=$1, zipcode=$2, city=$3, number=$4 WHERE a.address_id=$5;`
	for _, l := range company.Locations {
		_, err := tx.Exec(updateLocationsQuery,
			l.Street,
			l.Zipcode,
			l.City,
			l.Number,
			l.ID,
		)

		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	const updateProjectsQuery = `UPDATE "Project" p
	SET description=$1, compensation=$2, duration=$3, recommendations=$4
	WHERE p.project_id=$5;`
	for _, p := range company.Projects {
		_, err := tx.Exec(updateProjectsQuery,
			p.Description,
			p.Compensation,
			p.Duration,
			p.Recommendations,
			p.ID,
		)

		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	return nil
}

// FindByIDTx finds a company in the DB based on id. It should be used as PART of a
// unit of work, as a transaction gets passed in but will not be committed.
// This is the responsibility of the caller.
// It will rollback and return an error if something goes wrong
func (c CompanyRepo) FindByIDTx(tx *sql.Tx, id string) (d.Company, error) {
	var cID, info, name string
	var street, zip, city, num string
	var pID, desc, comp, dur string
	var recommendations []string
	var rID, jobTitle string
	var uID, fname, lname, email, hash, role string

	const selectCompanyQuery = `
		SELECT c.company_id, c.information, c.name
		FROM "Company" c
		WHERE c.company_id = $1;
	`
	const selectAddressesQuery = `
		SELECT a.street, a.zipcode, a.city, a.number
		FROM "Address" a
		WHERE ref_company = $1;
	`
	const selectProjectsQuery = `
		SELECT p.project_id, p.description, p.compensation, p.duration, p.recommendations
		FROM "Project" p
		WHERE p.ref_company = $1;
	`
	const selectRepresentativesQuery = `
		SELECT 	r.representative_id, r.job_title, u.user_id, u.first_name, u.last_name, u.email, u.hashed_password, u.role
		FROM "Representative" r
		JOIN "User" u on r.ref_user = u.user_id
		WHERE r.ref_company = $1;
	`
	companyResult := tx.QueryRow(selectCompanyQuery, id)
	err := companyResult.Scan(&cID, &info, &name)
	if err != nil {
		_ = tx.Rollback()
		return d.Company{}, err
	}

	addresses := []d.Address{}
	addressRows, err := tx.Query(selectAddressesQuery, id)
	if err != nil {
		_ = tx.Rollback()
		return d.Company{}, err
	}
	defer addressRows.Close()

	for addressRows.Next() {
		if err = addressRows.Scan(&street, &zip, &city, &num); err != nil {
			_ = tx.Rollback()
			return d.Company{}, err
		}
		addresses = append(addresses, d.Address{
			Street:  street,
			Zipcode: zip,
			City:    city,
			Number:  num,
		})
	}

	projects := []d.Project{}
	projectRows, err := tx.Query(selectProjectsQuery, id)
	if err != nil {
		_ = tx.Rollback()
		return d.Company{}, err
	}
	defer projectRows.Close()

	for projectRows.Next() {
		if err = projectRows.Scan(&pID, &desc, &comp, &dur, pq.Array(&recommendations)); err != nil {
			_ = tx.Rollback()
			return d.Company{}, err
		}
		projects = append(projects, d.Project{
			ID:              pID,
			Description:     desc,
			Duration:        dur,
			Compensation:    comp,
			Recommendations: recommendations,
			CompanyID:       cID,
		})
	}

	representatives := []d.Representative{}
	reprRows, err := tx.Query(selectRepresentativesQuery, id)
	if err != nil {
		_ = tx.Rollback()
		return d.Company{}, err
	}
	defer reprRows.Close()

	for reprRows.Next() {
		if err = reprRows.Scan(&rID, &jobTitle, &uID, &fname, &lname, &email, &hash, &role); err != nil {
			_ = tx.Rollback()
			return d.Company{}, err
		}
		representatives = append(representatives, d.Representative{
			ID:        rID,
			JobTitle:  jobTitle,
			CompanyID: cID,
			User: d.User{
				ID:             uID,
				FirstName:      fname,
				LastName:       lname,
				Email:          email,
				HashedPassword: hash,
				Role:           role,
			},
		})
	}

	return d.Company{
		ID:              cID,
		Name:            name,
		Information:     info,
		Locations:       addresses,
		Representatives: representatives,
		Projects:        projects,
	}, nil
}

// FindByNameTx finds a company in the DB based on name. It should be used as PART of a
// unit of work, as a transaction gets passed in but will not be committed.
// This is the responsibility of the caller.
// It will rollback and return an error if something goes wrong
func (c CompanyRepo) FindByNameTx(tx *sql.Tx, name string) (d.Company, error) {
	var cID, info, cName string
	var street, zip, city, num string
	var pID, desc, comp, dur string
	var recommendations []string
	var rID, jobTitle string
	var uID, fname, lname, email, hash, role string

	const selectCompanyQuery = `
		SELECT c.company_id, c.information, c.name
		FROM "Company" c
		WHERE c.name = $1;
	`
	const selectAddressesQuery = `
		SELECT a.street, a.zipcode, a.city, a.number
		FROM "Address" a
		WHERE a.ref_company = $1;
	`
	const selectProjectsQuery = `
		SELECT p.project_id, p.description, p.compensation, p.duration, p.recommendations
		FROM "Project" p
		WHERE p.ref_company = $1;
	`
	const selectRepresentativesQuery = `
		SELECT 	r.representative_id, r.job_title, u.user_id, u.first_name, u.last_name, u.email, u.hashed_password, u.role
		FROM "Representative" r
		JOIN "User" u on r.ref_user = u.user_id
		WHERE r.ref_company = $1;
	`
	companyResult := tx.QueryRow(selectCompanyQuery, name)
	err := companyResult.Scan(&cID, &info, &cName)
	if err != nil {
		_ = tx.Rollback()
		return d.Company{}, err
	}

	addresses := []d.Address{}
	addressRows, err := tx.Query(selectAddressesQuery, cID)
	if err != nil {
		_ = tx.Rollback()
		return d.Company{}, err
	}
	defer addressRows.Close()

	for addressRows.Next() {
		if err = addressRows.Scan(&street, &zip, &city, &num); err != nil {
			_ = tx.Rollback()
			return d.Company{}, err
		}
		addresses = append(addresses, d.Address{
			Street:  street,
			Zipcode: zip,
			City:    city,
			Number:  num,
		})
	}

	projects := []d.Project{}
	projectRows, err := tx.Query(selectProjectsQuery, cID)
	if err != nil {
		_ = tx.Rollback()
		return d.Company{}, err
	}
	defer projectRows.Close()

	for projectRows.Next() {
		if err = projectRows.Scan(&pID, &desc, &comp, &dur, pq.Array(&recommendations)); err != nil {
			_ = tx.Rollback()
			return d.Company{}, err
		}
		projects = append(projects, d.Project{
			ID:              pID,
			Description:     desc,
			Duration:        dur,
			Compensation:    comp,
			Recommendations: recommendations,
			CompanyID:       cID,
		})
	}

	representatives := []d.Representative{}
	reprRows, err := tx.Query(selectRepresentativesQuery, cID)
	if err != nil {
		_ = tx.Rollback()
		return d.Company{}, err
	}
	defer reprRows.Close()

	for reprRows.Next() {
		if err = reprRows.Scan(&rID, &jobTitle, &uID, &fname, &lname, &email, &hash, &role); err != nil {
			_ = tx.Rollback()
			return d.Company{}, err
		}
		representatives = append(representatives, d.Representative{
			ID:        rID,
			JobTitle:  jobTitle,
			CompanyID: cID,
			User: d.User{
				ID:             uID,
				FirstName:      fname,
				LastName:       lname,
				Email:          email,
				HashedPassword: hash,
				Role:           role,
			},
		})
	}

	return d.Company{
		ID:              cID,
		Name:            cName,
		Information:     info,
		Locations:       addresses,
		Representatives: representatives,
		Projects:        projects,
	}, nil
}
