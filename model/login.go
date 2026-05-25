package model

import (
	"database/sql"
	"errors"
	"rqms/dataStore/postgres"
)

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthUser struct {
	Email     string
	FirstName string
	LastName  string
}

// Authenticate checks customer credentials, returns full user info
func (l *Login) Authenticate() (*AuthUser, error) {
	var u AuthUser
	err := postgres.Db.QueryRow(
		"SELECT email, first_name, last_name FROM signup WHERE email=$1 AND password=$2",
		l.Email, l.Password,
	).Scan(&u.Email, &u.FirstName, &u.LastName)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("invalid email or password")
		}
		return nil, err
	}
	return &u, nil
}