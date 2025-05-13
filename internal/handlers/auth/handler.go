package handlers

import (
	"log/slog"

	"medods-test-task/internal/handlers"
	"medods-test-task/internal/notificator"
	"medods-test-task/internal/server/config"
	"medods-test-task/internal/storages/postgresql"
)

type AuthHandler struct {
	handlers.BaseHandler
	UserStorage postgresql.UserStorage
	AuthStorage postgresql.AuthStorage
}

func NewAuthHandler(userStorage postgresql.UserStorage, authStorage postgresql.AuthStorage, lgr *slog.Logger, cfg *config.Config, ntf notificator.Notificator) *AuthHandler {
	handler := &AuthHandler{
		BaseHandler: handlers.BaseHandler{
			Lgr: lgr,
			Cfg: cfg,
			Ntf: ntf,
		},
		UserStorage: userStorage,
		AuthStorage: authStorage,
	}
	return handler
}
