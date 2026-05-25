package model

import "rqms/controller/dataStore/postgres"





type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Note      string `json:"note"`
	Time      string `json:"time"`
	Status    string `json:"status"`
	TableNo   string `json:"table_no"`
}

// CreateUser adds a new queue entry
func (u *User) CreateUser() error {
	if u.Status == "" {
		u.Status = "waiting"
	}
	_, err := postgres.Db.Exec(
		`INSERT INTO users (first_name, last_name, email, phone, note, time, status, table_no)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
		u.FirstName, u.LastName, u.Email, u.Phone, u.Note, u.Time, u.Status, u.TableNo,
	)
	return err
}

// GetUsers returns all queue entries ordered by created_at
func GetUsers() ([]User, error) {
	rows, err := postgres.Db.Query(
		`SELECT id, first_name, last_name, email, phone, note, time,
		        COALESCE(status,'waiting'), COALESCE(table_no,'')
		 FROM users ORDER BY id`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		rows.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.Phone, &u.Note, &u.Time, &u.Status, &u.TableNo)
		users = append(users, u)
	}
	return users, nil
}

// DeleteUser removes a queue entry by id
func DeleteUser(id int) error {
	_, err := postgres.Db.Exec("DELETE FROM users WHERE id=$1", id)
	return err
}

// UpdateQueueStatus changes the status and optionally assigns a table
func UpdateQueueStatus(id int, status, tableNo string) error {
	_, err := postgres.Db.Exec(
		"UPDATE users SET status=$1, table_no=$2 WHERE id=$3",
		status, tableNo, id,
	)
	// If seating, mark the table as occupied
	if status == "seated" && tableNo != "" {
		postgres.Db.Exec("UPDATE tables SET status='occupied' WHERE table_no=$1", tableNo)
	}
	// If cancelled/done, free the table
	if (status == "cancelled" || status == "done") && tableNo != "" {
		postgres.Db.Exec("UPDATE tables SET status='available' WHERE table_no=$1", tableNo)
	}
	return err
}

// QueueStats returns counts per status
type QueueStats struct {
	Waiting   int `json:"waiting"`
	Seated    int `json:"seated"`
	Done      int `json:"done"`
	Cancelled int `json:"cancelled"`
	Total     int `json:"total"`
}

func GetQueueStats() (QueueStats, error) {
	var s QueueStats
	postgres.Db.QueryRow("SELECT COUNT(*) FROM users WHERE status='waiting'").Scan(&s.Waiting)
	postgres.Db.QueryRow("SELECT COUNT(*) FROM users WHERE status='seated'").Scan(&s.Seated)
	postgres.Db.QueryRow("SELECT COUNT(*) FROM users WHERE status='done'").Scan(&s.Done)
	postgres.Db.QueryRow("SELECT COUNT(*) FROM users WHERE status='cancelled'").Scan(&s.Cancelled)
	s.Total = s.Waiting + s.Seated + s.Done + s.Cancelled
	return s, nil
}
