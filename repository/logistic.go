package repository

import (
	"database/sql"
	"RescueHub/structs"
	"errors"
)

func isValidLogisticsStatus(status string) bool {
	validStatuses := []string{"available", "distributed", "out_of_stock"}
	for _, valid := range validStatuses {
			if status == valid {
					return true
			}
	}
	return false
}

func CreateLogistic(db *sql.DB, logistics *structs.Logistic) error {
	if !isValidLogisticsStatus(logistics.Status) {
			return errors.New("invalid logistics status")
	}

	sqlQuery := `INSERT INTO logistics (type, quantity, status, disaster_id, created_at, updated_at)
							VALUES ($1, $2, $3, $4, NOW(), NOW()) RETURNING id, created_at, updated_at`
	err := db.QueryRow(sqlQuery, logistics.Type, logistics.Quantity, logistics.Status, logistics.DisasterID).
			Scan(&logistics.ID, &logistics.CreatedAt, &logistics.UpdatedAt)

	if err != nil {
			return err
	}

	return nil
}

func GetAllLogistics(db *sql.DB) ([]structs.Logistic, error) {
	query := `SELECT id, type, quantity, status, disaster_id, created_at, updated_at FROM logistics`
	rows, err := db.Query(query)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logistics []structs.Logistic
	for rows.Next() {
		var logistic structs.Logistic
		err := rows.Scan(&logistic.ID, &logistic.Type, &logistic.Quantity, &logistic.Status, &logistic.DisasterID, &logistic.CreatedAt, &logistic.UpdatedAt)
		if err != nil {
			return nil, err
		}
		logistics = append(logistics, logistic)
	}
	return logistics, nil
}

func GetLogisticByID(db *sql.DB, id int) (structs.Logistic, error) {
	query := `SELECT id, type, quantity, status, disaster_id, created_at, updated_at FROM logistics WHERE id = $1`
	var logistic structs.Logistic
	err := db.QueryRow(query, id).Scan(&logistic.ID, &logistic.Type, &logistic.Quantity, &logistic.Status, &logistic.DisasterID, &logistic.CreatedAt, &logistic.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return logistic, errors.New("logistic not found")
		}
		return logistic, err
	}
	return logistic, nil
}

func UpdateLogistic(db *sql.DB, logistics structs.Logistic) error {
	if !isValidLogisticsStatus(logistics.Status) {
			return errors.New("invalid logistics status")
	}

	sqlQuery := `UPDATE logistics SET type = $1, quantity = $2, status = $3, disaster_id = $4, updated_at = CURRENT_TIMESTAMP WHERE id = $5`
	_, err := db.Exec(sqlQuery, logistics.Type, logistics.Quantity, logistics.Status, logistics.DisasterID, logistics.ID)
	if err != nil {
			return err
	}
	return nil
}

func DeleteLogistic(db *sql.DB, id int) error {
	sqlQuery := `DELETE FROM logistics WHERE id=$1`
	_, err := db.Exec(sqlQuery, id)
	if err != nil {
		return err
	}
	return nil
}
