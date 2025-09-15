package repository

import (
	"database/sql"
	"fmt"
	"log"
	db "otp-login-service/DB"
	"otp-login-service/models"
	"time"
)

func Create(phone string) (*models.User, error) {
	createdat := time.Now()
	query := `INSERT INTO users (phone_number, created_at) VALUES (?, ?)`
	res, err := db.DB.Exec(query, phone, createdat)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &models.User{
		ID:          id,
		PhoneNumber: phone,
		CreatedAt:   createdat,
	}, nil
}

func FindByPhoneNumber(phone string) (*models.User, error) {
	var user models.User
	query := `SELECT id, phone_number, created_at FROM users WHERE phone_number = ?`
	err := db.DB.QueryRow(query, phone).Scan(&user.ID, &user.PhoneNumber, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func FindByID(id int64) (*models.User, error) {
	var user models.User
	query := `SELECT id, phone_number, created_at FROM users WHERE id = ?`
	err := db.DB.QueryRow(query, id).Scan(&user.ID, &user.PhoneNumber, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func FindAll(page, limit int, search string) ([]*models.User, int) {
	var users []*models.User
	args := []interface{}{}

	countQuery := `SELECT COUNT(*) FROM users`
	if search != "" {
		countQuery += ` WHERE phone_number LIKE ?`
		args = append(args, "%"+search+"%")
	}

	var total int
	db.DB.QueryRow(countQuery, args...).Scan(&total)

	query := `SELECT id, phone_number, created_at FROM users`
	if search != "" {
		query += ` WHERE phone_number LIKE ?`
	}
	query += ` ORDER BY created_at DESC LIMIT ? OFFSET ?`

	paginatedArgs := append(args, limit, (page-1)*limit)
	rows, err := db.DB.Query(query, paginatedArgs...)
	if err != nil {
		return []*models.User{}, 0
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.PhoneNumber, &user.CreatedAt); err != nil {
			log.Printf("Error scanning user row: %v", err)
			return []*models.User{}, 0
		}
		users = append(users, &user)
	}

	return users, total
}
