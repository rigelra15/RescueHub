package repository

import (
	"RescueHub/structs"
	"database/sql"
	"errors"
)

func CreateVolunteer(db *sql.DB, volunteer structs.Volunteer) error {
	if !isValidVolunteerStatus(volunteer.Status) {
			return errors.New("invalid volunteer status")
	}

	sqlQuery := `INSERT INTO volunteers (user_id, disaster_id, skill, location, status, created_at, updated_at)
							 VALUES ($1, $2, $3, $4, $5, NOW(), NOW()) RETURNING id, created_at, updated_at`
	err := db.QueryRow(sqlQuery, volunteer.UserID, volunteer.DisasterID, volunteer.Skill, volunteer.Location, volunteer.Status).
			Scan(&volunteer.ID, &volunteer.CreatedAt, &volunteer.UpdatedAt)

	if err != nil {
			return err
	}

	return nil
}

func GetAllVolunteers(db *sql.DB) ([]structs.Volunteer, error) {
	query := `SELECT id, user_id, disaster_id, name, skill, location, status, created_at, updated_at FROM volunteers`
	rows, err := db.Query(query)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var volunteers []structs.Volunteer
	for rows.Next() {
		var volunteer structs.Volunteer
		err := rows.Scan(&volunteer.ID, &volunteer.UserID, &volunteer.DisasterID, &volunteer.Skill, &volunteer.Location, &volunteer.Status, &volunteer.CreatedAt, &volunteer.UpdatedAt)
		if err != nil {
			return nil, err
		}
		volunteers = append(volunteers, volunteer)
	}


	if len(volunteers) == 0 {
		return nil, errors.New("tidak ada daftar relawan yang tersedia")
	}

	return volunteers, nil
}

func GetVolunteerByID(db *sql.DB, id int) (structs.Volunteer, error) {
	query := `SELECT id, user_id, disaster_id, name, skill, location, status, created_at, updated_at FROM volunteers WHERE id = $1`
	var volunteer structs.Volunteer
	err := db.QueryRow(query, id).Scan(&volunteer.ID, &volunteer.UserID, &volunteer.DisasterID, &volunteer.Skill, &volunteer.Location, &volunteer.Status, &volunteer.CreatedAt, &volunteer.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return volunteer, errors.New("volunteer not found")
		}
		return volunteer, err
	}
	return volunteer, nil
}

func UpdateVolunteer(db *sql.DB, volunteer structs.Volunteer) error {
	if !isValidVolunteerStatus(volunteer.Status) {
			return errors.New("invalid volunteer status")
	}

	sqlQuery := `UPDATE volunteers SET user_id=$1, disaster_id=$2, skill=$3, location=$4, status=$5, updated_at=NOW() WHERE id=$6`
	_, err := db.Exec(sqlQuery, volunteer.UserID, volunteer.DisasterID, volunteer.Skill, volunteer.Location, volunteer.Status, volunteer.ID)
	if err != nil {
			return err
	}
	return nil
}

func DeleteVolunteer(db *sql.DB, id int) error {
	sqlQuery := `DELETE FROM volunteers WHERE id=$1`
	_, err := db.Exec(sqlQuery, id)
	if err != nil {
		return err
	}
	return nil
}

func isValidVolunteerStatus(status string) bool {
	validStatuses := []string{"available", "on_mission", "completed"}
	for _, valid := range validStatuses {
		if status == valid {
			return true
		}
	}
	return false
}