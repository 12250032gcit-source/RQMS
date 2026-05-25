package model

import "rqms/dataStore/postgres"

// Table represents a restaurant table
type Table struct {
	ID       int    `json:"id"`
	TableNo  string `json:"table_no"`
	Capacity int    `json:"capacity"`
	Status   string `json:"status"`
}

// GetTables returns all tables
func GetTables() ([]Table, error) {
	rows, err := postgres.Db.Query(
		"SELECT id, table_no, capacity, status FROM tables ORDER BY table_no",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []Table
	for rows.Next() {
		var t Table
		rows.Scan(&t.ID, &t.TableNo, &t.Capacity, &t.Status)
		tables = append(tables, t)
	}
	return tables, nil
}

// UpdateTableStatus changes a table's status
func UpdateTableStatus(tableNo, status string) error {
	_, err := postgres.Db.Exec(
		"UPDATE tables SET status=$1 WHERE table_no=$2",
		status, tableNo,
	)
	return err
}

// AddTable inserts a new table
func AddTable(tableNo string, capacity int) error {
	_, err := postgres.Db.Exec(
		"INSERT INTO tables (table_no, capacity, status) VALUES ($1, $2, 'available')",
		tableNo, capacity,
	)
	return err
}

// DeleteTable removes a table (only if not occupied)
func DeleteTable(tableNo string) error {
	_, err := postgres.Db.Exec("DELETE FROM tables WHERE table_no=$1", tableNo)
	return err
}
