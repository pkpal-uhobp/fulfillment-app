package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"go.uber.org/zap"

	core_logger "github.com/pkpal-uhobp/fulfillment-app/internal/core/logger"
	core_postgres_pool "github.com/pkpal-uhobp/fulfillment-app/internal/core/repository/pool"
	core_postgres_tx "github.com/pkpal-uhobp/fulfillment-app/internal/core/repository/tx"
	core_http_middleware "github.com/pkpal-uhobp/fulfillment-app/internal/core/transport/http/middleware"
	core_http_response "github.com/pkpal-uhobp/fulfillment-app/internal/core/transport/http/response"
	core_http_server "github.com/pkpal-uhobp/fulfillment-app/internal/core/transport/http/server"

	auth_postgres "github.com/pkpal-uhobp/fulfillment-app/internal/features/auth/repository/postgres"
	auth_service "github.com/pkpal-uhobp/fulfillment-app/internal/features/auth/service"
	auth_http "github.com/pkpal-uhobp/fulfillment-app/internal/features/auth/transport/http"
)

func main() {
	_ = godotenv.Load()

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	logConfig := core_logger.NewConfigMust()

	log, err := core_logger.NewLogger(logConfig)
	if err != nil {
		panic(err)
	}
	defer log.Close()

	postgresConfig := core_postgres_pool.NewConfigMust()

	db, err := core_postgres_pool.NewConnectionPool(ctx, postgresConfig)
	if err != nil {
		log.Fatal("create postgres connection pool", zap.Error(err))
	}
	defer db.Close()

	txManager := core_postgres_tx.NewTx(db)

	authRepo := auth_postgres.NewAuthRepository(txManager)

	authConfig := auth_service.NewConfigMust()

	authService := auth_service.NewAuthService(
		txManager,
		authRepo,
		authConfig,
	)

	tokenVerifier := auth_http.NewAccessTokenVerifier(authService)

	authMiddleware := core_http_middleware.Auth(tokenVerifier)

	v1 := core_http_server.NewAPIVersionRouter(
		core_http_server.ApiVersion1,
	)

	v1.SetRoleMiddleware(
		core_http_middleware.RequireRoles(tokenVerifier),
	)

	registerHealthRoute(v1, log)

	authHTTPHandler := auth_http.NewAuthHTTPHandler(
		log,
		authService,
		authMiddleware,
	)

	v1.RegisterRoutes(authHTTPHandler.Routes()...)

	httpConfig := core_http_server.NewConfigMust()

	httpServer := core_http_server.NewHTTPServer(
		httpConfig,
		log,
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(log),
		core_http_middleware.Panic(),
		core_http_middleware.Trace(),
	)

	httpServer.RegisterAPIRouters(v1)

	if err := httpServer.Run(ctx); err != nil {
		log.Fatal("run HTTP server", zap.Error(err))
	}
}

func registerHealthRoute(
	router *core_http_server.APIVersionRouter,
	log *core_logger.Logger,
) {
	router.RegisterRoutes(
		core_http_server.NewRoute(
			http.MethodGet,
			"/health",
			func(w http.ResponseWriter, r *http.Request) {
				response := core_http_response.NewHTTPResponseHandler(log, w)

				response.JSONResponse(
					map[string]string{
						"status": "ok",
					},
					http.StatusOK,
				)
			},
			nil,
		),
	)
}
