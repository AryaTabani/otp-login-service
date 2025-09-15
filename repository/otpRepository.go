package repository

import (
	"database/sql"
	db "otp-login-service/DB"

	"sync"
	"time"
)

type ratelimitentry struct {
	Reqs      int
	resettime time.Time
}

var (
	rateLimits = make(map[string]ratelimitentry)
	mu         sync.Mutex
)

func StoreOTP(phone, code string) error {
	expiresAt := time.Now().Add(2 * time.Minute)
	query := `INSERT OR REPLACE INTO otps (phone_number, code, expires_at) VALUES (?, ?, ?)`
	_, err := db.DB.Exec(query, phone, code, expiresAt)
	return err
}

func GetOTP(phone string) (string, error) {
	var code string
	query := `SELECT code FROM otps WHERE phone_number = ? AND expires_at > ?`
	err := db.DB.QueryRow(query, phone, time.Now()).Scan(&code)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}
	return code, nil
}

func DeleteOTP(phone string) error {
	query := `DELETE FROM otps WHERE phone_number = ?`
	_, err := db.DB.Exec(query, phone)
	return err
}

func CheckRateLimit(phone string) bool {
	mu.Lock()
	defer mu.Unlock()

	entry, found := rateLimits[phone]
	if !found || time.Now().After(entry.resettime) {
		rateLimits[phone] = ratelimitentry{
			Reqs:      1,
			resettime: time.Now().Add(10 * time.Minute),
		}
		return true
	}

	if entry.Reqs < 3 {
		entry.Reqs++
		rateLimits[phone] = entry
		return true
	}

	return false
}
