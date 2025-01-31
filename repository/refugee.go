package repository

import (
	"RescueHub/structs"
	"database/sql"
	"errors"
	"strconv"
	"strings"
)


func CreateRefugee(db *sql.DB, refugee *structs.Refugee) error {
	sqlQuery := `INSERT INTO refugees (name, age, condition, needs, shelter_id, disaster_id, created_at, updated_at)
	             VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW()) RETURNING id, created_at, updated_at`
	err := db.QueryRow(sqlQuery, refugee.Name, refugee.Age, refugee.Condition, refugee.Needs, refugee.ShelterID, refugee.DisasterID).
		Scan(&refugee.ID, &refugee.CreatedAt, &refugee.UpdatedAt)

	if err != nil {
		return err
	}

	return nil
}


func GetAllRefugees(db *sql.DB) ([]structs.Refugee, error) {
	query := `SELECT id, name, age, condition, needs, shelter_id, disaster_id, created_at, updated_at FROM refugees`
	rows, err := db.Query(query)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var refugees []structs.Refugee
	for rows.Next() {
		var refugee structs.Refugee
		err := rows.Scan(&refugee.ID, &refugee.Name, &refugee.Age, &refugee.Condition, &refugee.Needs, &refugee.ShelterID, &refugee.DisasterID, &refugee.CreatedAt, &refugee.UpdatedAt)
		if err != nil {
			return nil, err
		}
		refugees = append(refugees, refugee)
	}
	return refugees, nil
}


func GetRefugeeByID(db *sql.DB, id int) (structs.Refugee, error) {
	query := `SELECT id, name, age, condition, needs, shelter_id, disaster_id, created_at, updated_at FROM refugees WHERE id = $1`
	var refugee structs.Refugee
	err := db.QueryRow(query, id).Scan(&refugee.ID, &refugee.Name, &refugee.Age, &refugee.Condition, &refugee.Needs, &refugee.ShelterID, &refugee.DisasterID, &refugee.CreatedAt, &refugee.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return refugee, errors.New("refugee not found")
		}
		return refugee, err
	}
	return refugee, nil
}

func isRefugeeExists(db *sql.DB, id int) bool {
	query := `SELECT id FROM refugees WHERE id = $1`
	err := db.QueryRow(query, id).Scan(&id)
	return err == nil
}

func UpdateRefugee(db *sql.DB, refugee structs.Refugee) error {
	if !isRefugeeExists(db, refugee.ID) {
		return errors.New("refugee not found")
	}

	var updateFields []string
	var values []interface{}
	counter := 1

	if refugee.Name != "" {
		updateFields = append(updateFields, "name = $"+strconv.Itoa(counter))
		values = append(values, refugee.Name)
		counter++
	}
	if refugee.Age != 0 {
		updateFields = append(updateFields, "age = $"+strconv.Itoa(counter))
		values = append(values, refugee.Age)
		counter++
	}
	if refugee.Condition != "" {
		updateFields = append(updateFields, "condition = $"+strconv.Itoa(counter))
		values = append(values, refugee.Condition)
		counter++
	}
	if refugee.Needs != "" {
		updateFields = append(updateFields, "needs = $"+strconv.Itoa(counter))
		values = append(values, refugee.Needs)
		counter++
	}
	if refugee.ShelterID != nil {
		updateFields = append(updateFields, "shelter_id = $"+strconv.Itoa(counter))
		values = append(values, refugee.ShelterID)
		counter++
	}
	if refugee.DisasterID != nil {
		updateFields = append(updateFields, "disaster_id = $"+strconv.Itoa(counter))
		values = append(values, refugee.DisasterID)
		counter++
	}

	if len(updateFields) == 0 {
		return errors.New("tidak ada field yang dapat diperbarui")
	}

	updateFields = append(updateFields, "updated_at = NOW()")
	query := "UPDATE refugees SET " + strings.Join(updateFields, ", ") + " WHERE id = $" + strconv.Itoa(counter)
	values = append(values, refugee.ID)

	_, err := db.Exec(query, values...)
	if err != nil {
		return err
	}
	return nil
}

func DeleteRefugee(db *sql.DB, id int) error {
	sqlQuery := `DELETE FROM refugees WHERE id=$1`
	_, err := db.Exec(sqlQuery, id)
	if err != nil {
		return err
	}
	return nil
}

func GetDistributionLogsByRefugeeID(db *sql.DB, refugeeID int) ([]structs.DistributionLog, error) {
	var logs []structs.DistributionLog
	query := `SELECT id, logistic_id, origin, destination, sender_name, recipient_name, quantity_sent, sent_at, created_at, updated_at 
	          FROM distribution_logs WHERE recipient_name = (SELECT name FROM refugees WHERE id = $1)`
	rows, err := db.Query(query, refugeeID)
	if err != nil {
		return logs, err
	}
	defer rows.Close()

	for rows.Next() {
		var log structs.DistributionLog
		err := rows.Scan(&log.ID, &log.LogisticID, &log.Origin, &log.Destination, &log.SenderName, &log.RecipientName, &log.QuantitySent, &log.SentAt, &log.CreatedAt, &log.UpdatedAt)
		if err != nil {
			return logs, err
		}
		logs = append(logs, log)
	}
	return logs, nil
}