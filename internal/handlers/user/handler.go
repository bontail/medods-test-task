package handlers

import (
	"log/slog"

	"medods-test-task/internal/handlers"
	"medods-test-task/internal/notificator"
	"medods-test-task/internal/server/config"
	"medods-test-task/internal/storages/postgresql"
)

type UserHandler struct {
	handlers.BaseHandler
	UserStorage postgresql.UserStorage
}

func NewUserHandler(userStorage postgresql.UserStorage, logger *slog.Logger, config *config.Config, ntf notificator.Notificator) *UserHandler {
	handler := &UserHandler{
		BaseHandler: handlers.BaseHandler{
			Lgr: logger,
			Cfg: config,
			Ntf: ntf,
		},
		UserStorage: userStorage,
	}
	return handler
}
