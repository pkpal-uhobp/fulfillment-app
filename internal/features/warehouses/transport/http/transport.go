package warehouses_transport_http

import (
	"context"
	"net/http"

	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
	core_logger "github.com/pkpal-uhobp/fulfillment-app/internal/core/logger"
	core_http_server "github.com/pkpal-uhobp/fulfillment-app/internal/core/transport/http/server"
	warehouses_service "github.com/pkpal-uhobp/fulfillment-app/internal/features/warehouses/service"
)

type WarehousesHTTPHandler struct {
	log               *core_logger.Logger
	warehousesService WarehousesService
}

type WarehousesService interface {
	ListWarehouses(
		ctx context.Context,
		filter warehouses_service.WarehouseFilter,
	) ([]warehouses_service.WarehouseDTO, error)

	GetWarehouse(
		ctx context.Context,
		warehouseID int64,
	) (warehouses_service.WarehouseDTO, error)

	CreateWarehouse(
		ctx context.Context,
		input warehouses_service.CreateWarehouseInput,
	) (warehouses_service.WarehouseDTO, error)

	PatchWarehouse(
		ctx context.Context,
		warehouseID int64,
		input warehouses_service.PatchWarehouseInput,
	) (warehouses_service.WarehouseDTO, error)

	DeactivateWarehouse(
		ctx context.Context,
		warehouseID int64,
	) error

	ListStorageZones(
		ctx context.Context,
		warehouseID int64,
	) ([]warehouses_service.StorageZoneDTO, error)

	CreateStorageZone(
		ctx context.Context,
		input warehouses_service.CreateStorageZoneInput,
	) (warehouses_service.StorageZoneDTO, error)

	PatchStorageZone(
		ctx context.Context,
		zoneID int64,
		input warehouses_service.PatchStorageZoneInput,
	) (warehouses_service.StorageZoneDTO, error)

	ListGates(
		ctx context.Context,
		warehouseID int64,
	) ([]warehouses_service.GateDTO, error)

	CreateGate(
		ctx context.Context,
		input warehouses_service.CreateGateInput,
	) (warehouses_service.GateDTO, error)

	PatchGate(
		ctx context.Context,
		gateID int64,
		input warehouses_service.PatchGateInput,
	) (warehouses_service.GateDTO, error)

	ListProductTypes(ctx context.Context) ([]warehouses_service.ProductTypeDTO, error)

	ListCargoPlaceTypes(ctx context.Context) ([]warehouses_service.CargoPlaceTypeDTO, error)
}

func NewWarehousesHTTPHandler(
	log *core_logger.Logger,
	warehousesService WarehousesService,
) *WarehousesHTTPHandler {
	return &WarehousesHTTPHandler{
		log:               log,
		warehousesService: warehousesService,
	}
}

func (h *WarehousesHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		core_http_server.NewRoute(
			http.MethodGet,
			"/warehouses",
			h.ListWarehouses,
			nil,
		),
		core_http_server.NewRoute(
			http.MethodGet,
			"/warehouses/{id}",
			h.GetWarehouse,
			nil,
		),
		core_http_server.NewRoute(
			http.MethodPost,
			"/warehouses",
			h.CreateWarehouse,
			[]string{core_domain.RoleAdmin.String()},
		),
		core_http_server.NewRoute(
			http.MethodPatch,
			"/warehouses/{id}",
			h.PatchWarehouse,
			[]string{core_domain.RoleAdmin.String()},
		),
		core_http_server.NewRoute(
			http.MethodDelete,
			"/warehouses/{id}",
			h.DeactivateWarehouse,
			[]string{core_domain.RoleAdmin.String()},
		),

		core_http_server.NewRoute(
			http.MethodGet,
			"/storage-zones",
			h.ListStorageZones,
			nil,
		),
		core_http_server.NewRoute(
			http.MethodPost,
			"/storage-zones",
			h.CreateStorageZone,
			[]string{core_domain.RoleAdmin.String()},
		),
		core_http_server.NewRoute(
			http.MethodPatch,
			"/storage-zones/{id}",
			h.PatchStorageZone,
			[]string{core_domain.RoleAdmin.String()},
		),

		core_http_server.NewRoute(
			http.MethodGet,
			"/gates",
			h.ListGates,
			nil,
		),
		core_http_server.NewRoute(
			http.MethodPost,
			"/gates",
			h.CreateGate,
			[]string{core_domain.RoleAdmin.String()},
		),
		core_http_server.NewRoute(
			http.MethodPatch,
			"/gates/{id}",
			h.PatchGate,
			[]string{core_domain.RoleAdmin.String()},
		),

		core_http_server.NewRoute(
			http.MethodGet,
			"/product-types",
			h.ListProductTypes,
			nil,
		),
		core_http_server.NewRoute(
			http.MethodGet,
			"/cargo-place-types",
			h.ListCargoPlaceTypes,
			nil,
		),
	}
}
