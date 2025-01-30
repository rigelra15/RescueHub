package repository

import (
	"database/sql"
	"RescueHub/structs"
	"errors"
)

func CreateShelter(db *sql.DB, shelter structs.Shelter) error {
	sqlQuery := `INSERT INTO shelters (name, location, capacity_total, capacity_remaining, emergency_needs, disaster_id, created_at, updated_at)
	             VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW()) RETURNING id, created_at, updated_at`
	err := db.QueryRow(sqlQuery, shelter.Name, shelter.Location, shelter.CapacityTotal, shelter.CapacityRemaining, shelter.EmergencyNeeds, shelter.DisasterID).
		Scan(&shelter.ID, &shelter.CreatedAt, &shelter.UpdatedAt)

	if err != nil {
		return err
	}

	return nil
}

func GetAllShelters(db *sql.DB) ([]structs.Shelter, error) {
	query := `SELECT id, name, location, capacity_total, capacity_remaining, emergency_needs, disaster_id, created_at, updated_at FROM shelters`
	rows, err := db.Query(query)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var shelters []structs.Shelter
	for rows.Next() {
		var shelter structs.Shelter
		err := rows.Scan(&shelter.ID, &shelter.Name, &shelter.Location, &shelter.CapacityTotal, &shelter.CapacityRemaining, &shelter.EmergencyNeeds, &shelter.DisasterID, &shelter.CreatedAt, &shelter.UpdatedAt)
		if err != nil {
			return nil, err
		}
		shelters = append(shelters, shelter)
	}
	return shelters, nil
}

func GetShelterByID(db *sql.DB, id int) (structs.Shelter, error) {
	query := `SELECT id, name, location, capacity_total, capacity_remaining, emergency_needs, disaster_id, created_at, updated_at FROM shelters WHERE id = $1`
	var shelter structs.Shelter
	err := db.QueryRow(query, id).Scan(&shelter.ID, &shelter.Name, &shelter.Location, &shelter.CapacityTotal, &shelter.CapacityRemaining, &shelter.EmergencyNeeds, &shelter.DisasterID, &shelter.CreatedAt, &shelter.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return shelter, errors.New("shelter not found")
		}
		return shelter, err
	}
	return shelter, nil
}

func UpdateShelter(db *sql.DB, shelter structs.Shelter) error {
	sqlQuery := `UPDATE shelters SET name=$1, location=$2, capacity_total=$3, capacity_remaining=$4, emergency_needs=$5, disaster_id=$6, updated_at=NOW() WHERE id=$7`
	_, err := db.Exec(sqlQuery, shelter.Name, shelter.Location, shelter.CapacityTotal, shelter.CapacityRemaining, shelter.EmergencyNeeds, shelter.DisasterID, shelter.ID)
	if err != nil {
		return err
	}
	return nil
}

func DeleteShelter(db *sql.DB, id int) error {
	sqlQuery := `DELETE FROM shelters WHERE id=$1`
	_, err := db.Exec(sqlQuery, id)
	if err != nil {
		return err
	}
	return nil
}
