package repository

import (
	"RescueHub/structs"
	"database/sql"
	"errors"
	"strconv"
	"strings"
)

func isValidDonationStatus(status string) bool {
	validStatuses := []string{"pending", "confirmed", "rejected"}
	for _, valid := range validStatuses {
		if status == valid {
			return true
		}
	}
	return false
}

func CreateDonation(db *sql.DB, donation *structs.Donation) error {
	if !isValidDonationStatus(donation.Status) {
		return errors.New("invalid donation status")
	}

	sqlQuery := `INSERT INTO donations (donor_id, disaster_id, amount, item_name, status, created_at, updated_at)
	             VALUES ($1, $2, $3, $4, $5, NOW(), NOW()) RETURNING id, created_at, updated_at`
	err := db.QueryRow(sqlQuery, donation.DonorID, donation.DisasterID, donation.Amount, donation.ItemName, donation.Status).
		Scan(&donation.ID, &donation.CreatedAt, &donation.UpdatedAt)

	if err != nil {
		return err
	}

	return nil
}

func GetAllDonations(db *sql.DB) ([]structs.Donation, error) {
	query := `SELECT id, donor_id, disaster_id, amount, item_name, status, created_at, updated_at FROM donations`
	rows, err := db.Query(query)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var donations []structs.Donation
	for rows.Next() {
		var donation structs.Donation
		err := rows.Scan(&donation.ID, &donation.DonorID, &donation.DisasterID, &donation.Amount, &donation.ItemName, &donation.Status, &donation.CreatedAt, &donation.UpdatedAt)
		if err != nil {
			return nil, err
		}
		donations = append(donations, donation)
	}

	if len(donations) == 0 {
		return nil, errors.New("tidak ada daftar donasi yang tersedia")
	}

	return donations, nil
}

func GetDonationByID(db *sql.DB, id int) (structs.Donation, error) {
	query := `SELECT id, donor_id, disaster_id, amount, item_name, status, created_at, updated_at FROM donations WHERE id = $1`
	var donation structs.Donation
	err := db.QueryRow(query, id).Scan(&donation.ID, &donation.DonorID, &donation.DisasterID, &donation.Amount, &donation.ItemName, &donation.Status, &donation.CreatedAt, &donation.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return donation, errors.New("donation not found")
		}
		return donation, err
	}
	return donation, nil
}

func UpdateDonation(db *sql.DB, donation structs.Donation) error {
	if donation.Status != "" && !isValidDonationStatus(donation.Status) {
		return errors.New("invalid donation status")
	}

	var updateFields []string
	var values []interface{}
	counter := 1

	if donation.DonorID != nil {
		updateFields = append(updateFields, "donor_id = $"+strconv.Itoa(counter))
		values = append(values, donation.DonorID)
		counter++
	}
	if donation.DisasterID != nil {
		updateFields = append(updateFields, "disaster_id = $"+strconv.Itoa(counter))
		values = append(values, donation.DisasterID)
		counter++
	}
	if donation.Amount != 0 {
		updateFields = append(updateFields, "amount = $"+strconv.Itoa(counter))
		values = append(values, donation.Amount)
		counter++
	}
	if donation.ItemName != "" {
		updateFields = append(updateFields, "item_name = $"+strconv.Itoa(counter))
		values = append(values, donation.ItemName)
		counter++
	}
	if donation.Status != "" {
		updateFields = append(updateFields, "status = $"+strconv.Itoa(counter))
		values = append(values, donation.Status)
		counter++
	}

	if len(updateFields) == 0 {
		return errors.New("tidak ada field yang dapat diperbarui")
	}

	updateFields = append(updateFields, "updated_at = NOW()")
	query := "UPDATE donations SET " + strings.Join(updateFields, ", ") + " WHERE id = $" + strconv.Itoa(counter)
	values = append(values, donation.ID)

	_, err := db.Exec(query, values...)
	if err != nil {
		return err
	}
	return nil
}


func DeleteDonation(db *sql.DB, id int) error {
	sqlQuery := `DELETE FROM donations WHERE id=$1`
	_, err := db.Exec(sqlQuery, id)
	if err != nil {
		return err
	}
	return nil
}

