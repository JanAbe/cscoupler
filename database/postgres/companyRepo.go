package postgres

import (
	"database/sql"

	"github.com/lib/pq"

	"github.com/google/uuid"
	d "github.com/janabe/cscoupler/domain"
)

// CompanyRepo ...
type CompanyRepo struct {
	DB       *sql.DB
	reprRepo d.RepresentativeRepository
}

// Create ...
func (c CompanyRepo) Create(company d.Company) error {
	tx, err := c.DB.Begin()
	if err != nil {
		return err
	}

	const insertCompanyQuery = `INSERT INTO "Company"(company_id, name, information) VALUES ($1, $2, $3);`
	_, err = tx.Exec(insertCompanyQuery, company.ID, company.Name, company.Information)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	const insertAddressesQuery = `INSERT INTO "Address"(address_id, street, zipcode, city, number, ref_company) VALUES ($1, $2, $3, $4, $5, $6);`
	for _, c := range company.Locations {
		_, err = tx.Exec(insertAddressesQuery,
			uuid.New().String(),
			c.Street,
			c.Zipcode,
			c.City,
			c.Number,
			company.ID,
		)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	err = c.reprRepo.Create(company.Representatives[0])
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
func (c CompanyRepo) FindByID(id string) (d.Company, error) {
	tx, err := c.DB.Begin()
	if err != nil {
		return d.Company{}, err
	}

	var cID, info, name string
	var street, zip, city, num string
	var pID, desc, comp, dur string
	var recommendations []string
	var rID, jobTitle string
	var uID, fname, lname, email, hash, role string

	const selectCompanyQuery = `
		SELECT c.company_id, c.information, c.name
		FROM "Company" c
		WHERE c.company_id = '$1';
	`
	const selectAddressesQuery = `
		SELECT a.street, a.zipcode, a.city, a.number
		FROM "Address" a
		WHERE ref_company = '$1'
	`
	const selectProjectsQuery = `
		SELECT p.project_id, p.description, p.compensation, p.duration, p.recommendations
		FROM "Project" p
		WHERE p.ref_company = '$1'
	`
	const selectRepresentativesQuery = `
		SELECT 	r.representative_id, r.job_title, u.user_id, u.first_name, u.last_name, u.email, u.hashed_password, u.role
		FROM "Representative" r
		JOIN "User" u on r.ref_user = u.user_id
		WHERE r.ref_company = '$1';
	`
	companyResult := tx.QueryRow(selectCompanyQuery, id)
	err = companyResult.Scan(&cID, &info, &name)
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

// FindByName ...
func (c CompanyRepo) FindByName(name string) (d.Company, error) {
	return d.Company{}, nil
}

// FindAll ...
func (c CompanyRepo) FindAll() ([]d.Company, error) {
	return []d.Company{}, nil
}
