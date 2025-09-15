package services

import (
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"os"
	"otp-login-service/repository"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET_KEY"))

var (
	ErrRateLimitExceeded = errors.New("rate limit exceeded, please try again later")
	ErrInvalidOTP        = errors.New("invalid or expired OTP")
)

func RequestOTP(phone string) error {
	if !repository.CheckRateLimit(phone) {
		return ErrRateLimitExceeded
	}

	otp, err := generateRandomOTP(6)
	if err != nil {
		return fmt.Errorf("could not generate OTP: %w", err)
	}

	if err := repository.StoreOTP(phone, otp); err != nil {
		return fmt.Errorf("could not store OTP: %w", err)
	}

	log.Printf("OTP for %s: %s (expires in 2 minutes)\n", phone, otp)
	return nil
}

func VerifyOTP(phone, otp string) (string, error) {
	storedOTP, err := repository.GetOTP(phone)
	if err != nil {
		return "", fmt.Errorf("could not verify OTP from store: %w", err)
	}

	if storedOTP == "" || storedOTP != otp {
		return "", ErrInvalidOTP
	}

	repository.DeleteOTP(phone)

	user, err := repository.FindByPhoneNumber(phone)
	if err != nil {
		log.Printf("User with phone %s not found. Creating new user.", phone)
		user, err = repository.Create(phone)
		if err != nil {
			return "", fmt.Errorf("failed to create user: %w", err)
		}
	}

	token, err := generateUserToken(user.ID)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return token, nil
}

func generateRandomOTP(length int) (string, error) {
	buffer := make([]byte, length)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", err
	}
	for i := 0; i < length; i++ {
		buffer[i] = (buffer[i] % 10) + '0'
	}
	return string(buffer), nil
}

func generateUserToken(userID int64) (string, error) {
	if len(jwtSecret) == 0 {
		log.Println("Warning: JWT_SECRET_KEY is not set. Using a default insecure key.")
		jwtSecret = []byte("default-insecure-secret-key")
	}
	claims := jwt.MapClaims{
		"sub": userID,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
