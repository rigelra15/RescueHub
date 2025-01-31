package repository

import (
	"RescueHub/structs"
	"database/sql"
	"errors"
	"strconv"
	"strings"
)

func CreateVolunteer(db *sql.DB, volunteer *structs.Volunteer) error {
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
	query := `SELECT id, user_id, disaster_id, skill, location, status, created_at, updated_at FROM volunteers`
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
	query := `SELECT id, user_id, disaster_id, skill, location, status, created_at, updated_at FROM volunteers WHERE id = $1`
	var volunteer structs.Volunteer

	err := db.QueryRow(query, id).Scan(
		&volunteer.ID,
		&volunteer.UserID,
		&volunteer.DisasterID,
		&volunteer.Skill,
		&volunteer.Location,
		&volunteer.Status,
		&volunteer.CreatedAt,
		&volunteer.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return volunteer, errors.New("volunteer not found")
		}
		return volunteer, err
	}
	return volunteer, nil
}

func UpdateVolunteer(db *sql.DB, volunteer structs.Volunteer) error {
	var updateFields []string
	var values []interface{}
	counter := 1

	if volunteer.UserID != nil {
		updateFields = append(updateFields, "user_id = $"+strconv.Itoa(counter))
		values = append(values, volunteer.UserID)
		counter++
	}
	if volunteer.DisasterID != nil {
		updateFields = append(updateFields, "disaster_id = $"+strconv.Itoa(counter))
		values = append(values, volunteer.DisasterID)
		counter++
	}
	if volunteer.Skill != "" {
		updateFields = append(updateFields, "skill = $"+strconv.Itoa(counter))
		values = append(values, volunteer.Skill)
		counter++
	}
	if volunteer.Location != "" {
		updateFields = append(updateFields, "location = $"+strconv.Itoa(counter))
		values = append(values, volunteer.Location)
		counter++
	}
	if volunteer.Status != "" {
		if !isValidVolunteerStatus(volunteer.Status) {
			return errors.New("invalid volunteer status")
		}
		updateFields = append(updateFields, "status = $"+strconv.Itoa(counter))
		values = append(values, volunteer.Status)
		counter++
	}

	if len(updateFields) == 0 {
		return errors.New("tidak ada field yang dapat diperbarui")
	}

	updateFields = append(updateFields, "updated_at = NOW()")
	query := "UPDATE volunteers SET " + strings.Join(updateFields, ", ") + " WHERE id = $" + strconv.Itoa(counter)
	values = append(values, volunteer.ID)

	_, err := db.Exec(query, values...)
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