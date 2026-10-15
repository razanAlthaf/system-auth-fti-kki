package models

import (
	"database/sql"
	"errors"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
}

func SeedAdmin(db *sql.DB) {
	email := os.Getenv("ADMIN_EMAIL")
	username := os.Getenv("ADMIN_USERNAME")
	password := os.Getenv("ADMIN_PASSWORD")

	if email == "" || username == "" || password == "" {
		return
	}

	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE role = 'admin')").Scan(&exists)
	if err != nil || exists {
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return
	}

	_, err = db.Exec(
		"INSERT INTO users (email, username, password_hash, role) VALUES ($1, $2, $3, 'admin')",
		email, username, string(hashedPassword),
	)
	if err != nil {
		errors.New("gagal membuat akun admin pertama")
	}
}