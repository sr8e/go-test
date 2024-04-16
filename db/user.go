package db

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type User struct {
	Id           string
	DisplayName  string
	IconURL      string
	AccessToken  string
	RefreshToken string
	Expire       time.Time
	SecretHash   string
	SecretSalt   string
}

func finishTx(tx *sql.Tx, err error) error {
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("Rollback failed: %w, during handling error %w", rbErr, err)
		}
	} else {
		err = tx.Commit()
	}
	return err
}

func (u *User) Save() (err error) {
	tx, err := dbPool.Begin()
	if err != nil {
		return fmt.Errorf("cannot begin transaction: %w", err)
	}
	defer finishTx(tx, err)

	row := tx.QueryRow("select id from ir_user where id = $1;", u.Id)
	var dummy int
	noRow := row.Scan(&dummy)
	if noRow == nil {
		// user already exists, update
		_, err = tx.Exec(
			`update ir_user set (
				id, display_name, icon_url, access_token,
				refresh_token, expire, secret_hash, secret_salt
			)=($1, $2, $3, $4, $5, $6, $7, $8);`,
			u.Id, u.DisplayName, u.IconURL, u.AccessToken,
			u.RefreshToken, u.Expire, u.SecretHash, u.SecretSalt,
		)
	} else if errors.Is(noRow, sql.ErrNoRows) {
		// user not exist in table, create
		_, err = tx.Exec(
			`insert into ir_user (
				id, display_name, icon_url, access_token,
				refresh_token, expire, secret_hash, secret_salt
			) values ($1, $2, $3, $4, $5, $6, $7, $8);`,
			u.Id, u.DisplayName, u.IconURL, u.AccessToken,
			u.RefreshToken, u.Expire, u.SecretHash, u.SecretSalt,
		)
	} else {
		return noRow
	}
	return err
}

func (u *User) Get() (ok bool, err error) {
	row := dbPool.QueryRow(
		`select id, display_name, icon_url, access_token,
		refresh_token, expire, secret_hash, secret_salt
		from ir_user where id = $1;`,
		u.Id,
	)
	err = row.Scan(
		&u.Id, &u.DisplayName, &u.IconURL, &u.AccessToken,
		&u.RefreshToken, &u.Expire, &u.SecretHash, &u.SecretSalt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}