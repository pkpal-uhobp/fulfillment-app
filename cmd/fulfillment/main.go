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
	cargoitems_postgres "github.com/pkpal-uhobp/fulfillment-app/internal/features/cargoitems/repository/postgres"
	cargoitems_service "github.com/pkpal-uhobp/fulfillment-app/internal/features/cargoitems/service"
	cargoitems_http "github.com/pkpal-uhobp/fulfillment-app/internal/features/cargoitems/transport/http"
	orders_postgres "github.com/pkpal-uhobp/fulfillment-app/internal/features/orders/repository/postgres"
	orders_service "github.com/pkpal-uhobp/fulfillment-app/internal/features/orders/service"
	orders_http "github.com/pkpal-uhobp/fulfillment-app/internal/features/orders/transport/http"
	pickupcalendar_postgres "github.com/pkpal-uhobp/fulfillment-app/internal/features/pickupcalendar/repository/postgres"
	pickupcalendar_service "github.com/pkpal-uhobp/fulfillment-app/internal/features/pickupcalendar/service"
	pickupcalendar_http "github.com/pkpal-uhobp/fulfillment-app/internal/features/pickupcalendar/transport/http"
	shipments_postgres "github.com/pkpal-uhobp/fulfillment-app/internal/features/shipments/repository/postgres"
	shipments_service "github.com/pkpal-uhobp/fulfillment-app/internal/features/shipments/service"
	shipments_http "github.com/pkpal-uhobp/fulfillment-app/internal/features/shipments/transport/http"
	warehouses_postgres "github.com/pkpal-uhobp/fulfillment-app/internal/features/warehouses/repository/postgres"
	warehouses_service "github.com/pkpal-uhobp/fulfillment-app/internal/features/warehouses/service"
	warehouses_http "github.com/pkpal-uhobp/fulfillment-app/internal/features/warehouses/transport/http"
)

func main() {
	_ = godotenv.Load()
	setTestDefaults()

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

	warehousesRepo := warehouses_postgres.NewWarehousesRepository(
		txManager,
	)
	warehousesService := warehouses_service.NewWarehousesService(
		warehousesRepo,
	)
	warehousesHTTPHandler := warehouses_http.NewWarehousesHTTPHandler(
		log,
		warehousesService,
	)
	v1.RegisterRoutes(warehousesHTTPHandler.Routes()...)

	ordersRepo := orders_postgres.NewOrdersRepository(
		txManager,
	)
	ordersService := orders_service.NewOrdersService(
		ordersRepo,
	)
	ordersHTTPHandler := orders_http.NewOrdersHTTPHandler(
		log,
		ordersService,
	)
	v1.RegisterRoutes(ordersHTTPHandler.Routes()...)

	pickupCalendarRepo := pickupcalendar_postgres.NewPickupCalendarRepository(
		txManager,
	)
	pickupCalendarService := pickupcalendar_service.NewPickupCalendarService(
		pickupCalendarRepo,
	)
	pickupCalendarHTTPHandler := pickupcalendar_http.NewPickupCalendarHTTPHandler(
		log,
		pickupCalendarService,
	)
	v1.RegisterRoutes(pickupCalendarHTTPHandler.Routes()...)

	shipmentsRepo := shipments_postgres.NewShipmentsRepository(
		txManager,
	)
	shipmentsService := shipments_service.NewShipmentsService(
		shipmentsRepo,
	)
	shipmentsHTTPHandler := shipments_http.NewShipmentsHTTPHandler(
		log,
		shipmentsService,
	)
	v1.RegisterRoutes(shipmentsHTTPHandler.Routes()...)

	cargoItemsRepo := cargoitems_postgres.NewCargoItemsRepository(
		txManager,
	)
	cargoItemsService := cargoitems_service.NewCargoItemsService(
		cargoItemsRepo,
	)
	cargoItemsHTTPHandler := cargoitems_http.NewCargoItemsHTTPHandler(
		log,
		cargoItemsService,
	)
	v1.RegisterRoutes(cargoItemsHTTPHandler.Routes()...)

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

func setTestDefaults() {
	setDefaultEnv("HTTP_ADDR", ":8080")
	setDefaultEnv("POSTGRES_HOST", "127.0.0.1")
	setDefaultEnv("POSTGRES_PORT", "5433")
	setDefaultEnv("POSTGRES_USER", "postgres")
	setDefaultEnv("POSTGRES_PASSWORD", "postgres")
	setDefaultEnv("POSTGRES_DB", "fulfillment-app")
	setDefaultEnv("POSTGRES_SSL_MODE", "disable")
	setDefaultEnv("POSTGRES_QUERY_TIMEOUT", "5s")
	setDefaultEnv("JWT_SECRET", "dev-secret-for-test")
	setDefaultEnv("JWT_ACCESS_TTL", "15m")
	setDefaultEnv("JWT_REFRESH_TTL", "720h")
}

func setDefaultEnv(key string, value string) {
	if os.Getenv(key) == "" {
		_ = os.Setenv(key, value)
	}
}
