package model

import (
	"database/sql"
	"errors"
	"rqms/dataStore/postgres"
)

// Staff represents a restaurant staff member / admin
type Staff struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password,omitempty"`
	Role      string `json:"role"`
}

// RegisterStaff inserts a new staff record
func (s *Staff) RegisterStaff() error {
	_, err := postgres.Db.Exec(
		"INSERT INTO staff(first_name, last_name, email, password, role) VALUES($1,$2,$3,$4,$5)",
		s.FirstName, s.LastName, s.Email, s.Password, s.Role,
	)
	return err
}

// AuthenticateStaff checks staff credentials and returns role
func (s *Staff) AuthenticateStaff() (string, error) {
	var role string
	err := postgres.Db.QueryRow(
		"SELECT role FROM staff WHERE email=$1 AND password=$2",
		s.Email, s.Password,
	).Scan(&role)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("invalid email or password")
		}
		return "", err
	}
	return role, nil
}

// GetAllStaff returns all staff members
func GetAllStaff() ([]Staff, error) {
	rows, err := postgres.Db.Query(
		"SELECT id, first_name, last_name, email, role FROM staff ORDER BY id",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var staffList []Staff
	for rows.Next() {
		var s Staff
		rows.Scan(&s.ID, &s.FirstName, &s.LastName, &s.Email, &s.Role)
		staffList = append(staffList, s)
	}
	return staffList, nil
}
