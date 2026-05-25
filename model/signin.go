package model

import "rqms/dataStore/postgres"









type Sigin struct {
	FName    string `json:"first_name"`
	Lname    string `json:"last_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

const queryInsertUser = "INSERT INTO signup(first_name, last_name, email, password) VALUES($1,$2,$3,$4)"

func (s *Sigin) Adduser() error {
	_, err := postgres.Db.Exec(queryInsertUser, s.FName, s.Lname, s.Email, s.Password)
	return err
}
