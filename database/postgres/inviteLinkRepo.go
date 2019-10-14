package postgres

import (
	"database/sql"
	"time"

	d "github.com/janabe/cscoupler/domain"
)

// InviteLinkRepo ...
type InviteLinkRepo struct {
	DB *sql.DB
}

// Create ...
func (i InviteLinkRepo) Create(inviteLink d.InviteLink) error {
	tx, err := i.DB.Begin()
	if err != nil {
		return err
	}

	const insertQuery = `INSERT INTO "Invite_Link"(invite_link_id, url, created_at, expiry_date, used)
	VALUES($1, $2, $3, $4, $5);`
	_, err = tx.Exec(insertQuery,
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

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// FindByID ...
func (i InviteLinkRepo) FindByID(id string) (d.InviteLink, error) {
	tx, err := i.DB.Begin()
	if err != nil {
		return d.InviteLink{}, err
	}

	var iID, url string
	var createdAt, expiryDate time.Time
	var used bool
	const selectQuery = `SELECT i.invite_link_id, i.url, i.created_at, i.expiry_date, i.used
	FROM "Invite_Link" i WHERE i.invite_link_id=$1;`
	result := tx.QueryRow(selectQuery, id)

	err = result.Scan(&iID, &url, &createdAt, &expiryDate, &used)
	if err != nil {
		_ = tx.Rollback()
		return d.InviteLink{}, err
	}

	inviteLink := d.InviteLink{
		ID:         iID,
		CreatedAt:  createdAt,
		ExpiryDate: expiryDate,
		URL:        url,
		Used:       used,
	}

	err = tx.Commit()
	if err != nil {
		return d.InviteLink{}, err
	}

	return inviteLink, nil
}

// Update ...
func (i InviteLinkRepo) Update(inviteLink d.InviteLink) error {
	tx, err := i.DB.Begin()
	if err != nil {
		return err
	}

	// maybe extend so the expiry date can be postponed?
	const updateQuery = `UPDATE "Invite_Link" i SET used=$1 WHERE i.invite_link_id=$2;`
	_, err = tx.Exec(updateQuery, inviteLink.Used, inviteLink.ID)
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
