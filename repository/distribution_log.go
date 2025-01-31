package repository

import (
	"database/sql"
	"RescueHub/structs"
	"errors"
)

func CreateDistributionLog(db *sql.DB, log *structs.DistributionLog) error {
	sqlQuery := `INSERT INTO distribution_logs (logistic_id, origin, destination, distance, sender_name, recipient_name, quantity_sent, sent_at, created_at, updated_at)
							VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW()) RETURNING id, created_at, updated_at`
	err := db.QueryRow(sqlQuery, log.LogisticID, log.Origin, log.Destination, log.Distance, log.SenderName, log.RecipientName, log.QuantitySent, log.SentAt).
		Scan(&log.ID, &log.CreatedAt, &log.UpdatedAt)

	if err != nil {
		return err
	}

	return nil
}

func GetAllDistributionLogs(db *sql.DB) ([]structs.DistributionLog, error) {
	query := `SELECT id, logistic_id, origin, destination, distance, sender_name, recipient_name, quantity_sent, sent_at, created_at, updated_at FROM distribution_logs`
	rows, err := db.Query(query)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []structs.DistributionLog
	for rows.Next() {
		var log structs.DistributionLog
		err := rows.Scan(&log.ID, &log.LogisticID, &log.Origin, &log.Destination, &log.Distance, &log.SenderName, &log.RecipientName, &log.QuantitySent, &log.SentAt, &log.CreatedAt, &log.UpdatedAt)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}
	return logs, nil
}

func GetDistributionLogByID(db *sql.DB, id int) (structs.DistributionLog, error) {
	query := `SELECT id, logistic_id, origin, destination, distance, sender_name, recipient_name, quantity_sent, sent_at, created_at, updated_at FROM distribution_logs WHERE id = $1`
	var log structs.DistributionLog
	err := db.QueryRow(query, id).Scan(&log.ID, &log.LogisticID, &log.Origin, &log.Destination, &log.Distance, &log.SenderName, &log.RecipientName, &log.QuantitySent, &log.SentAt, &log.CreatedAt, &log.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return log, errors.New("distribution log not found")
		}
		return log, err
	}
	return log, nil
}

func UpdateDistributionLog(db *sql.DB, log structs.DistributionLog) error {
	sqlQuery := `UPDATE distribution_logs SET logistic_id=$1, origin=$2, destination=$3, distance=$4, sender_name=$5, recipient_name=$6, quantity_sent=$7, sent_at=$8, updated_at=CURRENT_TIMESTAMP WHERE id=$9`
	_, err := db.Exec(sqlQuery, log.LogisticID, log.Origin, log.Destination, log.Distance, log.SenderName, log.RecipientName, log.QuantitySent, log.SentAt, log.ID)
	if err != nil {
		return err
	}
	return nil
}

func DeleteDistributionLog(db *sql.DB, id int) error {
	sqlQuery := `DELETE FROM distribution_logs WHERE id=$1`
	_, err := db.Exec(sqlQuery, id)
	if err != nil {
		return err
	}
	return nil
}
