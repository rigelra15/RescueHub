package repository

import (
	"RescueHub/structs"
	"database/sql"
	"errors"
	"strings"
	"strconv"
)

func isDisasterExists(db *sql.DB, id int) bool {
	query := `SELECT id FROM disasters WHERE id = $1`
	err := db.QueryRow(query, id).Scan(&id)
	
	return err == nil
}

func isValidDisasterStatus(status string) bool {
	validStatuses := []string{"active", "resolved", "archived"}
	for _, valid := range validStatuses {
			if status == valid {
					return true
			}
	}
	return false
}

func CreateDisaster(db *sql.DB, disaster *structs.Disaster) error {
	if !isValidDisasterStatus(disaster.Status) {
		return errors.New("invalid disaster status")
	}

	if isDisasterExists(db, disaster.ID) {
		return errors.New("disaster already exists")
	}

	sqlQuery := `INSERT INTO disasters (type, location, description, status, reported_by, created_at, updated_at)
							VALUES ($1, $2, $3, $4, $5, NOW(), NOW()) RETURNING id, created_at, updated_at`
	
	err := db.QueryRow(sqlQuery, disaster.Type, disaster.Location, disaster.Description, disaster.Status, disaster.ReportedBy).
			Scan(&disaster.ID, &disaster.CreatedAt, &disaster.UpdatedAt)

	if err != nil {
		return err
	}

	return nil
}

func GetAllDisasters(db *sql.DB) ([]structs.Disaster, error) {
	query := `SELECT id, type, location, description, status, reported_by, created_at, updated_at FROM disasters`
	rows, err := db.Query(query)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var disasters []structs.Disaster
	for rows.Next() {
		var disaster structs.Disaster
		err := rows.Scan(&disaster.ID, &disaster.Type, &disaster.Location, &disaster.Description, &disaster.Status, &disaster.ReportedBy, &disaster.CreatedAt, &disaster.UpdatedAt)
		if err != nil {
			return nil, err
		}
		disasters = append(disasters, disaster)
	}
	return disasters, nil
}

func GetDisasterByID(db *sql.DB, id int) (structs.Disaster, error) {
	query := `SELECT id, type, location, description, status, reported_by, created_at, updated_at FROM disasters WHERE id = $1`
	var disaster structs.Disaster
	err := db.QueryRow(query, id).Scan(&disaster.ID, &disaster.Type, &disaster.Location, &disaster.Description, &disaster.Status, &disaster.ReportedBy, &disaster.CreatedAt, &disaster.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return disaster, errors.New("disaster not found")
		}
		return disaster, err
	}
	return disaster, nil
}

func UpdateDisaster(db *sql.DB, disaster structs.Disaster) error {
	if disaster.Status != "" && !isValidDisasterStatus(disaster.Status) {
		return errors.New("invalid disaster status")
	}

	if !isDisasterExists(db, disaster.ID) {
		return errors.New("disaster not found")
	}

	var updateFields []string
	var values []interface{}
	counter := 1

	if disaster.Type != "" {
		updateFields = append(updateFields, "type = $"+strconv.Itoa(counter))
		values = append(values, disaster.Type)
		counter++
	}
	if disaster.Location != "" {
		updateFields = append(updateFields, "location = $"+strconv.Itoa(counter))
		values = append(values, disaster.Location)
		counter++
	}
	if disaster.Description != "" {
		updateFields = append(updateFields, "description = $"+strconv.Itoa(counter))
		values = append(values, disaster.Description)
		counter++
	}
	if disaster.Status != "" {
		updateFields = append(updateFields, "status = $"+strconv.Itoa(counter))
		values = append(values, disaster.Status)
		counter++
	}
	if disaster.ReportedBy != 0 {
		updateFields = append(updateFields, "reported_by = $"+strconv.Itoa(counter))
		values = append(values, disaster.ReportedBy)
		counter++
	}

	if len(updateFields) == 0 {
		return errors.New("tidak ada field yang dapat diperbarui")
	}

	updateFields = append(updateFields, "updated_at = NOW()")
	query := "UPDATE disasters SET " + strings.Join(updateFields, ", ") + " WHERE id = $" + strconv.Itoa(counter)
	values = append(values, disaster.ID)

	_, err := db.Exec(query, values...)
	if err != nil {
		return err
	}
	return nil
}


func DeleteDisaster(db *sql.DB, id int) error {
	sqlQuery := `DELETE FROM disasters WHERE id=$1`
	_, err := db.Exec(sqlQuery, id)
	if err != nil {
		return err
	}
	return nil
}

func GetSheltersByDisasterID(db *sql.DB, disasterID int) ([]structs.Shelter, error) {
	query := `SELECT id, disaster_id, name, location, capacity_total, capacity_remaining, emergency_needs, created_at, updated_at FROM shelters WHERE disaster_id = $1`
	rows, err := db.Query(query, disasterID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var shelters []structs.Shelter
	for rows.Next() {
		var shelter structs.Shelter
		err := rows.Scan(&shelter.ID, &shelter.DisasterID, &shelter.Name, &shelter.Location, &shelter.CapacityTotal, &shelter.CapacityRemaining, &shelter.EmergencyNeeds, &shelter.CreatedAt, &shelter.UpdatedAt)
		if err != nil {
			return nil, err
		}
		shelters = append(shelters, shelter)
	}
	return shelters, nil
}

func GetVolunteersByDisasterID(db *sql.DB, disasterID int) ([]structs.Volunteer, error) {
	var volunteers []structs.Volunteer
	query := `SELECT id, donor_id, disaster_id, skill, location, status, created_at, updated_at 
	          FROM volunteers WHERE disaster_id = $1`
	rows, err := db.Query(query, disasterID)
	if err != nil {
		return volunteers, err
	}
	defer rows.Close()

	for rows.Next() {
		var volunteer structs.Volunteer
		err := rows.Scan(&volunteer.ID, &volunteer.UserID, &volunteer.DisasterID, &volunteer.Skill, &volunteer.Location, &volunteer.Status, &volunteer.CreatedAt, &volunteer.UpdatedAt)
		if err != nil {
			return volunteers, err
		}
		volunteers = append(volunteers, volunteer)
	}
	return volunteers, nil
}

func GetLogisticsByDisasterID(db *sql.DB, disasterID int) ([]structs.Logistic, error) {
	var logistics []structs.Logistic
	query := `SELECT id, type, quantity, status, disaster_id, created_at, updated_at 
	          FROM logistics WHERE disaster_id = $1`
	rows, err := db.Query(query, disasterID)
	if err != nil {
		return logistics, err
	}
	defer rows.Close()

	for rows.Next() {
		var logistic structs.Logistic
		err := rows.Scan(&logistic.ID, &logistic.Type, &logistic.Quantity, &logistic.Status, &logistic.DisasterID, &logistic.CreatedAt, &logistic.UpdatedAt)
		if err != nil {
			return logistics, err
		}
		logistics = append(logistics, logistic)
	}
	return logistics, nil
}

func GetEmergencyReportsByDisasterID(db *sql.DB, disasterID int) ([]structs.EmergencyReport, error) {
	var reports []structs.EmergencyReport
	query := `SELECT id, user_id, disaster_id, description, location, created_at, updated_at 
	          FROM emergency_reports WHERE disaster_id = $1`
	rows, err := db.Query(query, disasterID)
	if err != nil {
		return reports, err
	}
	defer rows.Close()

	for rows.Next() {
		var report structs.EmergencyReport
		err := rows.Scan(&report.ID, &report.UserID, &report.DisasterID, &report.Description, &report.Location, &report.CreatedAt, &report.UpdatedAt)
		if err != nil {
			return reports, err
		}
		reports = append(reports, report)
	}
	return reports, nil
}

func GetEvacuationRoutesByDisasterID(db *sql.DB, disasterID int) ([]structs.EvacuationRoute, error) {
	var routes []structs.EvacuationRoute
	query := `SELECT id, disaster_id, origin, destination, distance, route, status, created_at, updated_at 
	          FROM evacuation_routes WHERE disaster_id = $1`
	rows, err := db.Query(query, disasterID)
	if err != nil {
		return routes, err
	}
	defer rows.Close()

	for rows.Next() {
		var route structs.EvacuationRoute
		err := rows.Scan(&route.ID, &route.DisasterID, &route.Origin, &route.Destination, &route.Distance, &route.Route, &route.Status, &route.CreatedAt, &route.UpdatedAt)
		if err != nil {
			return routes, err
		}
		routes = append(routes, route)
	}
	return routes, nil
}

func GetRefugeesByDisasterID(db *sql.DB, disasterID int) ([]structs.Refugee, error) {
	query := `SELECT id, name, age, condition, needs, shelter_id, created_at, updated_at 
	          FROM refugees WHERE disaster_id = $1`
	rows, err := db.Query(query, disasterID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var refugees []structs.Refugee
	for rows.Next() {
		var refugee structs.Refugee
		err := rows.Scan(
			&refugee.ID, &refugee.Name, &refugee.Age, &refugee.Condition, 
			&refugee.Needs, &refugee.ShelterID, &refugee.CreatedAt, &refugee.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		refugees = append(refugees, refugee)
	}

	if len(refugees) == 0 {
		return nil, errors.New("tidak ada daftar pengungsi yang tersedia")
	}

	return refugees, nil
}