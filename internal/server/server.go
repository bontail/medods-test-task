package server

import (
	"context"
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/golang-migrate/migrate/v4"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	handlersAuth "medods-test-task/internal/handlers/auth"
	handlersUser "medods-test-task/internal/handlers/user"
	customLogger "medods-test-task/internal/logger"
	loggerConfig "medods-test-task/internal/logger/config"
	"medods-test-task/internal/middlewares"
	"medods-test-task/internal/notificator"
	notificatorConfig "medods-test-task/internal/notificator/config"
	serverConfig "medods-test-task/internal/server/config"
	"medods-test-task/internal/storages/postgresql"
)

func runMigrations(dbUrl string) error {
	m, err := migrate.New(
		"file://migrations",
		dbUrl+"?sslmode=disable",
	)
	if err != nil {
		return err
	}
	defer func(m *migrate.Migrate) {
		_, _ = m.Close()
	}(m)

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	return nil
}

func Run(ctx context.Context) {
	serverCfg, err := serverConfig.NewConfig()
	if err != nil {
		panic(err)
	}
	gin.SetMode(serverCfg.GinMode)

	if err = runMigrations(serverCfg.DatabaseUrl); err != nil {
		panic(err)
	}

	pool, err := pgxpool.New(ctx, serverCfg.DatabaseUrl)
	if err != nil {
		panic(err)
	}

	loggerCfg, err := loggerConfig.NewConfig()
	if err != nil {
		panic(err)
	}

	logger, err := customLogger.NewLogger(loggerCfg)
	if err != nil {
		panic(err)
	}

	notificatorCfg, err := notificatorConfig.NewConfig()
	if err != nil {
		panic(err)
	}

	ntf := notificator.Notificator(notificator.NewHTTPNotificator(notificatorCfg, logger))

	userStorage := postgresql.UserStorage(postgresql.NewDefaultUserStorage(pool, logger))
	authStorage := postgresql.AuthStorage(postgresql.NewDefaultAuthStorage(pool, logger))

	jwtMiddleware := middlewares.Middleware(middlewares.NewJWTMiddleware(logger, serverCfg))
	lgrMiddleware := middlewares.Middleware(middlewares.NewLoggerMiddleware(logger, serverCfg))

	authHandler := handlersAuth.NewAuthHandler(userStorage, authStorage, logger, serverCfg, ntf)
	userHandler := handlersUser.NewUserHandler(userStorage, logger, serverCfg, ntf)

	router := gin.New()
	router.Use(lgrMiddleware.MiddlewareFunc())

	authGroup := router.Group("/auth")
	authGroup.POST("/signIn", authHandler.SignIn)
	authGroup.POST("/signOut", jwtMiddleware.MiddlewareFunc(), authHandler.SignOut)
	authGroup.POST("/refresh", jwtMiddleware.MiddlewareFunc(), authHandler.Refresh)

	userGroup := router.Group("/user")
	userGroup.GET("", jwtMiddleware.MiddlewareFunc(), userHandler.GetUser)
	userGroup.POST("/register", userHandler.Register)

	err = router.Run(serverCfg.Host + ":" + strconv.Itoa(serverCfg.Port))
	if err != nil {
		logger.Error(err.Error())
		panic(err)
	}
}
