package repository

import (
	"RescueHub/structs"
	"database/sql"
	"errors"
	"strings"
	"strconv"
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

func isDistributionLogExists(db *sql.DB, id int) bool {
	query := `SELECT EXISTS(SELECT 1 FROM distribution_logs WHERE id = $1)`
	var exists bool
	err := db.QueryRow(query, id).Scan(&exists)
	if err != nil {
		return false
	}
	return exists
}

func UpdateDistributionLog(db *sql.DB, log structs.DistributionLog) error {
	if !isDistributionLogExists(db, log.ID) {
		return errors.New("distribution log not found")
	}

	var updateFields []string
	var values []interface{}
	counter := 1

	if log.LogisticID != nil {
		updateFields = append(updateFields, "logistic_id = $"+strconv.Itoa(counter))
		values = append(values, log.LogisticID)
		counter++
	}
	if log.Origin != "" {
		updateFields = append(updateFields, "origin = $"+strconv.Itoa(counter))
		values = append(values, log.Origin)
		counter++
	}
	if log.Destination != "" {
		updateFields = append(updateFields, "destination = $"+strconv.Itoa(counter))
		values = append(values, log.Destination)
		counter++
	}
	if log.Distance != 0 {
		updateFields = append(updateFields, "distance = $"+strconv.Itoa(counter))
		values = append(values, log.Distance)
		counter++
	}
	if log.SenderName != "" {
		updateFields = append(updateFields, "sender_name = $"+strconv.Itoa(counter))
		values = append(values, log.SenderName)
		counter++
	}
	if log.RecipientName != "" {
		updateFields = append(updateFields, "recipient_name = $"+strconv.Itoa(counter))
		values = append(values, log.RecipientName)
		counter++
	}
	if log.QuantitySent != 0 {
		updateFields = append(updateFields, "quantity_sent = $"+strconv.Itoa(counter))
		values = append(values, log.QuantitySent)
		counter++
	}
	if !log.SentAt.IsZero() {
		updateFields = append(updateFields, "sent_at = $"+strconv.Itoa(counter))
		values = append(values, log.SentAt)
		counter++
	}

	if len(updateFields) == 0 {
		return errors.New("tidak ada field yang dapat diperbarui")
	}

	updateFields = append(updateFields, "updated_at = NOW()")
	query := "UPDATE distribution_logs SET " + strings.Join(updateFields, ", ") + " WHERE id = $" + strconv.Itoa(counter)
	values = append(values, log.ID)

	_, err := db.Exec(query, values...)
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
