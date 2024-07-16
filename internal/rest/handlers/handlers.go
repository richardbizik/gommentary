package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/richardbizik/gommentary/internal/database"
)

type OApiHandlers struct {
	ValidateUsers bool
	DB            *database.Sqlite
}

func RegisterOApiHandlers(r chi.Router, validateUsers bool, database *database.Sqlite) {
	HandlerFromMux(NewStrictHandler(OApiHandlers{
		ValidateUsers: validateUsers,
		DB:            database,
	}, []StrictMiddlewareFunc{}), r)
}
