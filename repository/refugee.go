package repository

import (
	"database/sql"
	"RescueHub/structs"
	"errors"
)

func CreateRefugee(db *sql.DB, refugee structs.Refugee) error {
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

func UpdateRefugee(db *sql.DB, refugee structs.Refugee) error {
	sqlQuery := `UPDATE refugees SET name=$1, age=$2, condition=$3, needs=$4, shelter_id=$5, disaster_id=$6, updated_at=NOW() WHERE id=$7`
	_, err := db.Exec(sqlQuery, refugee.Name, refugee.Age, refugee.Condition, refugee.Needs, refugee.ShelterID, refugee.DisasterID, refugee.ID)
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
