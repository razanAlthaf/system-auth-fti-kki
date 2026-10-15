package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	"auth-system/internal/config"
	"auth-system/internal/models"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

// Contoh handler test sederhana
func HealthCheck(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{
        "status": "success",
        "message": "Backend server is running smoothly",
    })
}

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Email    string `json:"email"`
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var userID string
	err = config.DB.QueryRow(
		"INSERT INTO users (email, username, password_hash, role) VALUES ($1, $2, $3, 'user') RETURNING id",
		req.Email, req.Username, string(hashedPassword),
	).Scan(&userID)

	if err != nil {
		http.Error(w, "Email atau username sudah terdaftar", http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Registrasi berhasil"})
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Identifier string `json:"identifier"` // email atau username
		Password   string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	var user models.User
	err := config.DB.QueryRow(
		"SELECT id, email, username, password_hash, role, created_at FROM users WHERE email = $1 OR username = $1",
		req.Identifier,
	).Scan(&user.ID, &user.Email, &user.Username, &user.PasswordHash, &user.Role, &user.CreatedAt)

	if err != nil {
		http.Error(w, "Kredensial tidak valid", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		http.Error(w, "Kredensial tidak valid", http.StatusUnauthorized)
		return
	}

	// Aturan 2 Sesi Maksimal: Hapus sesi tertua jika sudah >= 2
	rows, err := config.DB.Query("SELECT id FROM sessions WHERE user_id = $1 ORDER BY created_at ASC", user.ID)
	if err == nil {
		var sessionIDs []string
		for rows.Next() {
			var sid string
			rows.Scan(&sid)
			sessionIDs = append(sessionIDs, sid)
		}
		rows.Close()

		if len(sessionIDs) >= 2 {
			// Hapus sesi paling lama (index 0 hingga surplus)
			for i := 0; i <= len(sessionIDs)-2; i++ {
				config.DB.Exec("DELETE FROM sessions WHERE id = $1", sessionIDs[i])
			}
		}
	}

	// Buat JWT Token
	tokenClaims := jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Simpan hash token sesi ke database
	tokenHash := hashToken(tokenString)
	_, err = config.DB.Exec(
		"INSERT INTO sessions (user_id, token_hash, expires_at) VALUES ($1, $2, $3)",
		user.ID, tokenHash, time.Now().Add(24*time.Hour),
	)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": tokenString, "role": user.Role})
}

func Logout(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		w.WriteHeader(http.StatusOK)
		return
	}
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	tokenHash := hashToken(tokenString)

	config.DB.Exec("DELETE FROM sessions WHERE token_hash = $1", tokenHash)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Berhasil keluar"})
}

func GetProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	var user models.User
	err := config.DB.QueryRow(
		"SELECT id, email, username, role, created_at FROM users WHERE id = $1",
		userID,
	).Scan(&user.ID, &user.Email, &user.Username, &user.Role, &user.CreatedAt)

	if err != nil {
		http.Error(w, "User tidak ditemukan", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	role := r.Context().Value("role").(string)
	if role != "admin" {
		http.Error(w, "Akses ditolak", http.StatusForbidden)
		return
	}

	rows, err := config.DB.Query("SELECT id, email, username, role, created_at FROM users")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		rows.Scan(&u.ID, &u.Email, &u.Username, &u.Role, &u.CreatedAt)
		users = append(users, u)
	}

	json.NewEncoder(w).Encode(users)
}

func ForceLogout(w http.ResponseWriter, r *http.Request) {
	role := r.Context().Value("role").(string)
	if role != "admin" {
		http.Error(w, "Akses ditolak", http.StatusForbidden)
		return
	}

	targetUserID := r.URL.Query().Get("user_id")
	if targetUserID == "" {
		http.Error(w, "Parameter user_id diperlukan", http.StatusBadRequest)
		return
	}

	_, err := config.DB.Exec("DELETE FROM sessions WHERE user_id = $1", targetUserID)
	if err != nil {
		http.Error(w, "Gagal melakukan force logout", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Pengguna berhasil dipaksa keluar dari seluruh sesi"})
}