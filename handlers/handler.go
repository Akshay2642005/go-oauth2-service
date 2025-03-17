package handlers

import (
	"server/prisma/db"
	"server/services/auth"
)

type Handler struct {
	db   *db.PrismaClient
	auth *auth.AuthService
}

func New(db *db.PrismaClient, auth *auth.AuthService) *Handler {
	return &Handler{
		db:   db,
		auth: auth,
	}
}
