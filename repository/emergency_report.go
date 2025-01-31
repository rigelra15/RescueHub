package repository

import (
	"RescueHub/structs"
	"database/sql"
	"errors"
	"strconv"
	"strings"
)

func CreateEmergencyReport(db *sql.DB, report *structs.EmergencyReport) error {
	sqlQuery := `INSERT INTO emergency_reports (user_id, disaster_id, description, location, created_at, updated_at)
							VALUES ($1, $2, $3, $4, NOW(), NOW()) RETURNING id, created_at, updated_at`
	err := db.QueryRow(sqlQuery, report.UserID, report.DisasterID, report.Description, report.Location).
		Scan(&report.ID, &report.CreatedAt, &report.UpdatedAt)

	if err != nil {
		return err
	}

	return nil
}

func GetAllEmergencyReports(db *sql.DB) ([]structs.EmergencyReport, error) {
	query := `SELECT id, user_id, disaster_id, description, location, created_at, updated_at FROM emergency_reports`
	rows, err := db.Query(query)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reports []structs.EmergencyReport
	for rows.Next() {
		var report structs.EmergencyReport
		err := rows.Scan(&report.ID, &report.UserID, &report.DisasterID, &report.Description, &report.Location, &report.CreatedAt, &report.UpdatedAt)
		if err != nil {
			return nil, err
		}
		reports = append(reports, report)
	}
	return reports, nil
}

func GetEmergencyReportByID(db *sql.DB, id int) (structs.EmergencyReport, error) {
	query := `SELECT id, user_id, disaster_id, description, location, created_at, updated_at FROM emergency_reports WHERE id = $1`
	var report structs.EmergencyReport
	err := db.QueryRow(query, id).Scan(&report.ID, &report.UserID, &report.DisasterID, &report.Description, &report.Location, &report.CreatedAt, &report.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return report, errors.New("emergency report not found")
		}
		return report, err
	}
	return report, nil
}

func isEmergencyReportExists(db *sql.DB, id int) bool {
	query := `SELECT EXISTS(SELECT 1 FROM emergency_reports WHERE id = $1)`
	var exists bool
	db.QueryRow(query, id).Scan(&exists)
	return exists
}

func UpdateEmergencyReport(db *sql.DB, report structs.EmergencyReport) error {
	if !isEmergencyReportExists(db, report.ID) {
		return errors.New("emergency report not found")
	}

	var updateFields []string
	var values []interface{}
	counter := 1

	if report.UserID != nil {
		updateFields = append(updateFields, "user_id = $"+strconv.Itoa(counter))
		values = append(values, report.UserID)
		counter++
	}
	if report.DisasterID != nil {
		updateFields = append(updateFields, "disaster_id = $"+strconv.Itoa(counter))
		values = append(values, report.DisasterID)
		counter++
	}
	if report.Description != "" {
		updateFields = append(updateFields, "description = $"+strconv.Itoa(counter))
		values = append(values, report.Description)
		counter++
	}
	if report.Location != "" {
		updateFields = append(updateFields, "location = $"+strconv.Itoa(counter))
		values = append(values, report.Location)
		counter++
	}

	if len(updateFields) == 0 {
		return errors.New("tidak ada field yang dapat diperbarui")
	}

	updateFields = append(updateFields, "updated_at = NOW()")
	query := "UPDATE emergency_reports SET " + strings.Join(updateFields, ", ") + " WHERE id = $" + strconv.Itoa(counter)
	values = append(values, report.ID)

	_, err := db.Exec(query, values...)
	if err != nil {
		return err
	}
	return nil
}

func DeleteEmergencyReport(db *sql.DB, id int) error {
	sqlQuery := `DELETE FROM emergency_reports WHERE id=$1`
	_, err := db.Exec(sqlQuery, id)
	if err != nil {
		return err
	}
	return nil
}
