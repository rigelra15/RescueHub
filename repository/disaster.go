package repository

import (
	"RescueHub/structs"
	"database/sql"
	"errors"
)

func isDisasterExists(db *sql.DB, id int) bool {
	query := `SELECT id FROM disasters WHERE id = $1`
	err := db.QueryRow(query, id).Scan(&id)
	
	return err == nil
}

func isValidDisasterStatus(status string) bool {
	validStatuses := []string{"active", "resolved", "archived"}
	for _, valid := range validStatuses {
			if status == valid {
					return true
			}
	}
	return false
}

func CreateDisaster(db *sql.DB, disaster structs.Disaster) error {
	if !isValidDisasterStatus(disaster.Status) {
		return errors.New("invalid disaster status")
	}

	if isDisasterExists(db, disaster.ID) {
		return errors.New("disaster already exists")
	}

	sqlQuery := `INSERT INTO disasters (type, location, description, status, reported_by, created_at, updated_at)
							 VALUES ($1, $2, $3, $4, $5, NOW(), NOW()) RETURNING id, created_at, updated_at`
	err := db.QueryRow(sqlQuery, disaster.Type, disaster.Location, disaster.Description, disaster.Status, disaster.ReportedBy).
			Scan(&disaster.ID, &disaster.CreatedAt, &disaster.UpdatedAt)

	if err != nil {
			return err
	}

	return nil
}

func GetAllDisasters(db *sql.DB) ([]structs.Disaster, error) {
	query := `SELECT id, type, location, description, status, reported_by, created_at, updated_at FROM disasters`
	rows, err := db.Query(query)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var disasters []structs.Disaster
	for rows.Next() {
		var disaster structs.Disaster
		err := rows.Scan(&disaster.ID, &disaster.Type, &disaster.Location, &disaster.Description, &disaster.Status, &disaster.ReportedBy, &disaster.CreatedAt, &disaster.UpdatedAt)
		if err != nil {
			return nil, err
		}
		disasters = append(disasters, disaster)
	}
	return disasters, nil
}

func GetDisasterByID(db *sql.DB, id int) (structs.Disaster, error) {
	query := `SELECT id, type, location, description, status, reported_by, created_at, updated_at FROM disasters WHERE id = $1`
	var disaster structs.Disaster
	err := db.QueryRow(query, id).Scan(&disaster.ID, &disaster.Type, &disaster.Location, &disaster.Description, &disaster.Status, &disaster.ReportedBy, &disaster.CreatedAt, &disaster.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return disaster, errors.New("disaster not found")
		}
		return disaster, err
	}
	return disaster, nil
}

func UpdateDisaster(db *sql.DB, disaster structs.Disaster) error {
	if !isValidDisasterStatus(disaster.Status) {
		return errors.New("invalid disaster status")
	}

	if isDisasterExists(db, disaster.ID) {
		return errors.New("disaster not found")
	}

	sqlQuery := `UPDATE disasters SET type = $1, location = $2, description = $3, status = $4, reported_by = $5, updated_at = CURRENT_TIMESTAMP WHERE id = $6`
	_, err := db.Exec(sqlQuery, disaster.Type, disaster.Location, disaster.Description, disaster.Status, disaster.ReportedBy, disaster.ID)
	if err != nil {
			return err
	}
	return nil
}

func DeleteDisaster(db *sql.DB, id int) error {
	sqlQuery := `DELETE FROM disasters WHERE id=$1`
	_, err := db.Exec(sqlQuery, id)
	if err != nil {
		return err
	}
	return nil
}