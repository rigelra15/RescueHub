package repository

import (
	"RescueHub/structs"
	"database/sql"
	"errors"
	"strconv"
	"strings"
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

func CreateEvacuationRoute(db *sql.DB, route *structs.EvacuationRoute) error {
	if !isValidEvacuationStatus(route.Status) {
		return errors.New("invalid evacuation route status")
	}

	sqlQuery := `INSERT INTO evacuation_routes (disaster_id, origin, destination, distance, route, status, created_at, updated_at)
							VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW()) RETURNING id, created_at, updated_at`
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
	query := `SELECT id, disaster_id, origin, destination, distance, route, status, created_at, updated_at FROM evacuation_routes WHERE id = $1`
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

func isEvacuationRouteExists(db *sql.DB, id int) bool {
	query := `SELECT EXISTS(SELECT 1 FROM evacuation_routes WHERE id = $1)`
	var exists bool
	db.QueryRow(query, id).Scan(&exists)
	return exists
}

func UpdateEvacuationRoute(db *sql.DB, route structs.EvacuationRoute) error {
	if !isEvacuationRouteExists(db, route.ID) {
		return errors.New("evacuation route not found")
	}

	if route.Status != "" && !isValidEvacuationStatus(route.Status) {
		return errors.New("invalid evacuation route status")
	}

	var updateFields []string
	var values []interface{}
	counter := 1

	if route.DisasterID != nil {
		updateFields = append(updateFields, "disaster_id = $"+strconv.Itoa(counter))
		values = append(values, route.DisasterID)
		counter++
	}
	if route.Origin != "" {
		updateFields = append(updateFields, "origin = $"+strconv.Itoa(counter))
		values = append(values, route.Origin)
		counter++
	}
	if route.Destination != "" {
		updateFields = append(updateFields, "destination = $"+strconv.Itoa(counter))
		values = append(values, route.Destination)
		counter++
	}
	if route.Distance != 0 {
		updateFields = append(updateFields, "distance = $"+strconv.Itoa(counter))
		values = append(values, route.Distance)
		counter++
	}
	if route.Route != "" {
		updateFields = append(updateFields, "route = $"+strconv.Itoa(counter))
		values = append(values, route.Route)
		counter++
	}
	if route.Status != "" {
		updateFields = append(updateFields, "status = $"+strconv.Itoa(counter))
		values = append(values, route.Status)
		counter++
	}

	if len(updateFields) == 0 {
		return errors.New("tidak ada field yang dapat diperbarui")
	}

	updateFields = append(updateFields, "updated_at = NOW()")
	query := "UPDATE evacuation_routes SET " + strings.Join(updateFields, ", ") + " WHERE id = $" + strconv.Itoa(counter)
	values = append(values, route.ID)

	_, err := db.Exec(query, values...)
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
