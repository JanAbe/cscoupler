package postgres

import (
	"database/sql"
	"time"

	d "github.com/janabe/cscoupler/domain"
)

// InviteLinkRepo struct for postgres database
type InviteLinkRepo struct {
	DB *sql.DB
}

// Create inserts an InviteLink in the DB. It should be used as a single
// unit of work, as it has its own transaction inside.
func (i InviteLinkRepo) Create(inviteLink d.InviteLink) error {
	tx, err := i.DB.Begin()
	if err != nil {
		return err
	}

	err = i.CreateTx(tx, inviteLink)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// FindByID finds an InviteLink in the DB based on id. It should be used as a single
// unit of work, as it has its own transaction inside.
func (i InviteLinkRepo) FindByID(id string) (d.InviteLink, error) {
	tx, err := i.DB.Begin()
	if err != nil {
		return d.InviteLink{}, err
	}

	inviteLink, err := i.FindByIDTx(tx, id)
	if err != nil {
		return d.InviteLink{}, err
	}

	err = tx.Commit()
	if err != nil {
		return d.InviteLink{}, err
	}

	return inviteLink, nil
}

// FindByCreator finds all inviteLinks in the DB that are created by the provided
// representativeID.
func (i InviteLinkRepo) FindByCreator(representativeID string) ([]d.InviteLink, error) {
	tx, err := i.DB.Begin()
	if err != nil {
		return nil, err
	}

	invitations, err := i.FindByCreatorTx(tx, representativeID)
	if err != nil {
		return []d.InviteLink{}, err
	}

	err = tx.Commit()
	if err != nil {
		return []d.InviteLink{}, err
	}

	return invitations, nil
}

// Update updates an InviteLink in the DB. It should be used as a single
// unit of work, as it has its own transaction inside.
func (i InviteLinkRepo) Update(inviteLink d.InviteLink) error {
	tx, err := i.DB.Begin()
	if err != nil {
		return err
	}

	err = i.UpdateTx(tx, inviteLink)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// CreateTx inserts an InviteLink in the DB. It should be used as PART of a
// unit of work, as a transaction gets passed in but will not be committed.
// This is the responsibility of the caller.
// It will rollback and return an error if something goes wrong
func (i InviteLinkRepo) CreateTx(tx *sql.Tx, inviteLink d.InviteLink) error {
	const insertQuery = `INSERT INTO "Invite_Link"(invite_link_id, url, created_at, expiry_date, used, ref_representative)
	VALUES($1, $2, $3, $4, $5, $6);`
	_, err := tx.Exec(insertQuery,
		inviteLink.ID,
		inviteLink.URL,
		inviteLink.CreatedAt,
		inviteLink.ExpiryDate,
		inviteLink.Used,
		inviteLink.CreatedBy,
	)

	if err != nil {
		_ = tx.Rollback()
		return err
	}

	return nil
}

// FindByIDTx finds an InviteLink in the DB based on id. It should be used as PART of a
// unit of work, as a transaction gets passed in but will not be committed.
// This is the responsibility of the caller.
// It will rollback and return an error if something goes wrong
func (i InviteLinkRepo) FindByIDTx(tx *sql.Tx, id string) (d.InviteLink, error) {
	var iID, url, createdBy string
	var createdAt, expiryDate time.Time
	var used bool
	const selectQuery = `SELECT i.invite_link_id, i.url, i.created_at, i.expiry_date, i.used, i.ref_representative
	FROM "Invite_Link" i WHERE i.invite_link_id=$1;`
	result := tx.QueryRow(selectQuery, id)

	err := result.Scan(&iID, &url, &createdAt, &expiryDate, &used, &createdBy)
	if err != nil {
		_ = tx.Rollback()
		return d.InviteLink{}, err
	}

	return d.InviteLink{
		ID:         iID,
		CreatedAt:  createdAt,
		ExpiryDate: expiryDate,
		URL:        url,
		Used:       used,
		CreatedBy:  createdBy,
	}, nil
}

// FindByCreatorTx finds all inviteLinks in the DB that are created by the representative with the provided
// representativeID. It should be used as PART of a unit of work, as a transaction gets passed in but will not be committed.
// This is the responsibility of the caller.
// It will rollback and return an error if something goes wrong.
func (i InviteLinkRepo) FindByCreatorTx(tx *sql.Tx, representativeID string) ([]d.InviteLink, error) {
	const selectQuery = `SELECT i.invite_link_id, i.url, i.created_at, i.expiry_date, i.used, i.ref_representative
	FROM "Invite_Link" i WHERE i.ref_representative=$1;`

	rows, err := tx.Query(selectQuery, representativeID)
	if err != nil {
		_ = tx.Rollback()
		return []d.InviteLink{}, err
	}
	defer rows.Close()

	invitations := []d.InviteLink{}
	for rows.Next() {
		var iID, url, createdBy string
		var createdAt, expiryDate time.Time
		var used bool

		if err := rows.Scan(&iID, &url, &createdAt, &expiryDate, &used, &createdBy); err != nil {
			_ = tx.Rollback()
			return []d.InviteLink{}, err
		}

		invitations = append(invitations, d.InviteLink{
			ID:         iID,
			URL:        url,
			CreatedBy:  createdBy,
			CreatedAt:  createdAt,
			Used:       used,
			ExpiryDate: expiryDate,
		})
	}

	return invitations, nil
}

// UpdateTx updates an InviteLink in the DB. It should be used as PART of a
// unit of work, as a transaction gets passed in but will not be committed.
// This is the responsibility of the caller.
// It will rollback and return an error if something goes wrong
func (i InviteLinkRepo) UpdateTx(tx *sql.Tx, inviteLink d.InviteLink) error {
	// maybe extend so the expiry date can be postponed?
	const updateQuery = `UPDATE "Invite_Link" i SET used=$1 WHERE i.invite_link_id=$2;`
	_, err := tx.Exec(updateQuery, inviteLink.Used, inviteLink.ID)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	return nil
}
