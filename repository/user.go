package repository

import (
	"RescueHub/structs"
	"database/sql"
	"errors"
	"time"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func SaveOTP(db *sql.DB, userID int, otp string) error {
	sqlQuery := `UPDATE users SET otp_code = $1, otp_expiry = $2 WHERE id = $3`
	res, err := db.Exec(sqlQuery, otp, time.Now().Add(5*time.Minute), userID)
	if err != nil {
			fmt.Println("Error SaveOTP:", err)
			return err
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
			return errors.New("user ID tidak ditemukan")
	}

	return nil
}

func ValidateOTP(db *sql.DB, userID int, otp string) (bool, error) {
	var storedOTP string
	var expiry time.Time

	err := db.QueryRow("SELECT otp_code, otp_expiry FROM users WHERE id = $1", userID).Scan(&storedOTP, &expiry)
	if err != nil {
		return false, err
	}

	if storedOTP != otp || time.Now().After(expiry) {
		return false, errors.New("OTP salah atau telah kedaluwarsa")
	}

	_, _ = db.Exec("UPDATE users SET otp_code = NULL, otp_expiry = NULL WHERE id = $1", userID)

	return true, nil
}

func Enable2FA(db *sql.DB, email string, isEnabled bool) error {
	sqlQuery := `UPDATE users SET is_2fa = $1 WHERE email = $2`
	_, err := db.Exec(sqlQuery, isEnabled, email)
	return err
}

func isValidUserRole(role string) bool {
	validRoles := []string{"admin", "donor", "user"}
	for _, valid := range validRoles {
		if role == valid {
			return true
		}
	}
	return false
}

func IsEmailExists(db *sql.DB, email string) (bool, error) {
	var exists bool
	err := db.QueryRow("SELECT EXISTS (SELECT 1 FROM users WHERE email = $1)", email).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
	

func CreateUser(db *sql.DB, user structs.User) error {
	emailExists, err := IsEmailExists(db, user.Email)
	if err != nil {
		return err
	}

	if !isValidUserRole(user.Role) {
		return errors.New("invalid user role")
	}

	if emailExists {
		return errors.New("email sudah terdaftar")
	}

	sqlQuery := `INSERT INTO users (name, email, password, role, contact) VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at, updated_at`
	_, err = db.Exec(sqlQuery, user.Name, user.Email, user.Password, user.Role, user.Contact)
	if err != nil {
		return err
	}

	return nil
}

func GetUserByID(db *sql.DB, id int) (structs.User, error) {
	var user structs.User
	err := db.QueryRow("SELECT id, name, email, role, contact, created_at, updated_at FROM users WHERE id = $1", id).Scan(&user.ID, &user.Name, &user.Email, &user.Role, &user.Contact, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return user, err
	}

	return user, nil
}

func UpdateUser(db *sql.DB, user structs.User) error {
	emailExists, err := IsEmailExists(db, user.Email)
	if err != nil {
		return err
	}

	if !isValidUserRole(user.Role) {
		return errors.New("invalid user role")
	}
	
	if emailExists {
		return errors.New("email sudah terdaftar")
	}

	sqlQuery := `UPDATE users SET name = $1, email = $2, role = $3, contact = $4, updated_at = CURRENT_TIMESTAMP WHERE id = $5`
	_, err = db.Exec(sqlQuery, user.Name, user.Email, user.Role, user.Contact, user.ID)
	if err != nil {
		return err
	}

	return nil
}

func UpdateUserInfoWithoutEmail(db *sql.DB, user structs.User) error {
	sqlQuery := `UPDATE users SET name = $1, role = $2, contact = $3, updated_at = CURRENT_TIMESTAMP WHERE id = $4`
	_, err := db.Exec(sqlQuery, user.Name, user.Role, user.Contact, user.ID)
	if err != nil {
		return err
	}

	return nil
}

func DeleteUser(db *sql.DB, id int) error {
	sqlQuery := `DELETE FROM users WHERE id = $1`
	_, err := db.Exec(sqlQuery, id)
	if err != nil {
		return err
	}

	return nil
}

func GetAllUsers(db *sql.DB) ([]structs.User, error) {
	var users []structs.User
	rows, err := db.Query("SELECT id, name, email, role, contact, created_at, updated_at FROM users")
	if err != nil {
		return users, err
	}
	defer rows.Close()

	for rows.Next() {
		var user structs.User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Role, &user.Contact, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return users, err
		}

		users = append(users, user)
	}

	return users, nil
}

func GetAllUsersByRole(db *sql.DB, role string) ([]structs.User, error) {
	var users []structs.User
	rows, err := db.Query("SELECT id, name, email, role, contact FROM users WHERE role = $1", role)
	if err != nil {
		return users, err
	}
	defer rows.Close()

	for rows.Next() {
		var user structs.User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Role, &user.Contact)
		if err != nil {
			return users, err
		}

		users = append(users, user)
	}

	return users, nil
}

func GetUserByEmail(db *sql.DB, email string) (structs.User, error) {
	var user structs.User
	sqlQuery := "SELECT id, name, email, password, role, contact, created_at, updated_at FROM users WHERE email = $1"
	err := db.QueryRow(sqlQuery, email).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role, &user.Contact, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, errors.New("user tidak ditemukan")
		}
		return user, err
	}
	return user, nil
}

func CheckPasswordHash(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}