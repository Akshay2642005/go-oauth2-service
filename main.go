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
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	prismaClient := db.NewClient()
	if err := prismaClient.Prisma.Connect(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer prismaClient.Prisma.Disconnect()

	sessionStore := auth.NewCookieStore(auth.SessionOptions{
		CookiesKey: config.Envs.CookiesAuthSecret,
		MaxAge:     config.Envs.CookiesAuthAgeInSeconds,
		Secure:     config.Envs.CookiesAuthIsSecure,
		HttpOnly:   config.Envs.CookiesAuthIsHttpOnly,
	})
	authService := auth.NewAuthService(sessionStore)

	router := mux.NewRouter()

	handler := handlers.New(prismaClient, authService)

	// Auth
	router.HandleFunc("/auth/{provider}", handler.HandleProviderLogin).Methods("GET")
	router.HandleFunc("/auth/{provider}/callback", handler.HandleAuthCallbackFunction).Methods("GET")
	router.HandleFunc("/auth/logout/{provider}", handler.HandleLogout).Methods("GET")

	log.Printf("Server: Listening on %s:%s\n", config.Envs.PublicHost, config.Envs.Port)
	log.Fatalln(http.ListenAndServe(fmt.Sprintf(":%s", config.Envs.Port), router))
}
