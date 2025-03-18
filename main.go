package main

import (
	"fmt"
	"log"
	"net/http"
	"server/config"
	"server/handlers"
	"server/prisma/db"
	"server/services/auth"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize Prisma client
	prismaClient := db.NewClient()
	if err := prismaClient.Prisma.Connect(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer prismaClient.Prisma.Disconnect()

	// Initialize session store
	sessionStore := auth.NewCookieStore(auth.SessionOptions{
		CookiesKey: config.Envs.CookiesAuthSecret,
		MaxAge:     config.Envs.CookiesAuthAgeInSeconds,
		Secure:     config.Envs.CookiesAuthIsSecure,
		HttpOnly:   config.Envs.CookiesAuthIsHttpOnly,
		SameSite:   "None",
	})

	// Initialize AuthService
	authService := auth.NewAuthService(sessionStore)

	// Initialize Router
	router := mux.NewRouter()
	handler := handlers.New(prismaClient, authService)

	// Define Auth Routes
	router.HandleFunc("/auth/{provider}", handler.HandleProviderLogin).Methods("GET")
	router.HandleFunc("/auth/{provider}/callback", handler.HandleAuthCallbackFunction).Methods("GET")
	router.HandleFunc("/auth/logout/{provider}", handler.HandleLogout).Methods("GET")
	router.HandleFunc("/auth/session", handler.HandleSessionUser).Methods("GET", "OPTIONS")

	// ✅ Wrap the router with CORS middleware
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodOptions},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})
	routerWithCORS := c.Handler(router)

	// Start Server
	serverAddr := fmt.Sprintf(":%s", config.Envs.Port)
	log.Printf("Server: Listening on %s\n", serverAddr)
	log.Fatal(http.ListenAndServe(serverAddr, routerWithCORS))
}
