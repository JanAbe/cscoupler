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
	const insertQuery = `INSERT INTO "Invite_Link"(invite_link_id, url, created_at, expiry_date, used)
	VALUES($1, $2, $3, $4, $5);`
	_, err := tx.Exec(insertQuery,
		inviteLink.ID,
		inviteLink.URL,
		inviteLink.CreatedAt,
		inviteLink.ExpiryDate,
		inviteLink.Used,
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
	var iID, url string
	var createdAt, expiryDate time.Time
	var used bool
	const selectQuery = `SELECT i.invite_link_id, i.url, i.created_at, i.expiry_date, i.used
	FROM "Invite_Link" i WHERE i.invite_link_id=$1;`
	result := tx.QueryRow(selectQuery, id)

	err := result.Scan(&iID, &url, &createdAt, &expiryDate, &used)
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
	}, nil
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
