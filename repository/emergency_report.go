package repository

import (
	"database/sql"
	"RescueHub/structs"
	"errors"
)

func CreateEmergencyReport(db *sql.DB, report structs.EmergencyReport) error {
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

func UpdateEmergencyReport(db *sql.DB, report structs.EmergencyReport) error {
	sqlQuery := `UPDATE emergency_reports SET user_id=$1, disaster_id=$2, description=$3, location=$4, updated_at=NOW() WHERE id=$5`
	_, err := db.Exec(sqlQuery, report.UserID, report.DisasterID, report.Description, report.Location, report.ID)
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
