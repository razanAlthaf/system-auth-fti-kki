package main

import (
	"log"
	"net/http"
	"os"

	"auth-system/internal/config"
	"auth-system/internal/handlers"
	"auth-system/internal/middleware"
	"auth-system/internal/models"

	"github.com/joho/godotenv"
)

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	godotenv.Load()
	config.ConnectDB()

	models.SeedAdmin(config.DB)

	mux := http.NewServeMux()
	mux.HandleFunc("/api/health", handlers.HealthCheck)
	mux.HandleFunc("/api/register", handlers.Register)
	mux.HandleFunc("/api/login", handlers.Login)
	mux.HandleFunc("/api/logout", middleware.AuthMiddleware(handlers.Logout))
	mux.HandleFunc("/api/profile", middleware.AuthMiddleware(handlers.GetProfile))
	mux.HandleFunc("/api/users", middleware.AuthMiddleware(handlers.GetUsers))
	mux.HandleFunc("/api/force-logout", middleware.AuthMiddleware(handlers.ForceLogout))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Backend berjalan di port %s...", port)
	log.Fatal(http.ListenAndServe(":"+port, enableCORS(mux)))
}