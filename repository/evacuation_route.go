package repository

import (
	"database/sql"
	"RescueHub/structs"
	"errors"
)

func isValidEvacuationStatus(status string) bool {
	validStatuses := []string{"safe", "risky", "blocked"}
	for _, valid := range validStatuses {
			if status == valid {
					return true
			}
	}
	return false
}

func CreateEvacuationRoute(db *sql.DB, route structs.EvacuationRoute) error {
	if !isValidEvacuationStatus(route.Status) {
		return errors.New("invalid evacuation route status")
	}

	sqlQuery := `INSERT INTO evacuation_routes (disaster_id, origin, destination, distance, route, status, created_at, updated_at)
	             VALUES ($1, $2, $3, $4, NOW(), NOW()) RETURNING id, created_at, updated_at`
	err := db.QueryRow(sqlQuery, route.DisasterID, route.Origin, route.Destination, route.Distance, route.Route, route.Status).
		Scan(&route.ID, &route.CreatedAt, &route.UpdatedAt)

	if err != nil {
		return err
	}

	return nil
}

func GetAllEvacuationRoutes(db *sql.DB) ([]structs.EvacuationRoute, error) {
	query := `SELECT id, disaster_id, origin, destination, distance, route, status, created_at, updated_at FROM evacuation_routes`
	rows, err := db.Query(query)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var routes []structs.EvacuationRoute
	for rows.Next() {
		var route structs.EvacuationRoute
		err := rows.Scan(&route.ID, &route.DisasterID, &route.Origin, &route.Destination, &route.Distance, &route.Route, &route.Status, &route.CreatedAt, &route.UpdatedAt)
		if err != nil {
			return nil, err
		}
		routes = append(routes, route)
	}
	return routes, nil
}

func GetEvacuationRouteByID(db *sql.DB, id int) (structs.EvacuationRoute, error) {
	query := `SELECT id, disaster_id, start_location, end_location, status, created_at, updated_at FROM evacuation_routes WHERE id = $1`
	var route structs.EvacuationRoute
	err := db.QueryRow(query, id).Scan(&route.ID, &route.DisasterID, &route.Origin, &route.Destination, &route.Distance, &route.Route, &route.Status, &route.CreatedAt, &route.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return route, errors.New("evacuation route not found")
		}
		return route, err
	}
	return route, nil
}

func UpdateEvacuationRoute(db *sql.DB, route structs.EvacuationRoute) error {
	if !isValidEvacuationStatus(route.Status) {
		return errors.New("invalid evacuation route status")
	}

	sqlQuery := `UPDATE evacuation_routes SET disaster_id=$1, origin=$2, destination=$3, distance=$4, route=$5, status=$6, updated_at=NOW() WHERE id=$7`
	_, err := db.Exec(sqlQuery, route.DisasterID, route.Origin, route.Destination, route.Distance, route.Route, route.Status, route.ID)
	if err != nil {
		return err
	}
	return nil
}

func DeleteEvacuationRoute(db *sql.DB, id int) error {
	sqlQuery := `DELETE FROM evacuation_routes WHERE id=$1`
	_, err := db.Exec(sqlQuery, id)
	if err != nil {
		return err
	}
	return nil
}
