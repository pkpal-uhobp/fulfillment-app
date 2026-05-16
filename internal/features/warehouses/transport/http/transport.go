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
	ActivateWarehouse(
		ctx context.Context,
		warehouseID int64,
	) error
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
	ActivateStorageZone(
		ctx context.Context,
		zoneID int64,
	) error
	DeactivateStorageZone(
		ctx context.Context,
		zoneID int64,
	) error
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
	ActivateGate(
		ctx context.Context,
		gateID int64,
	) error
	DeactivateGate(
		ctx context.Context,
		gateID int64,
	) error
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
	adminRoles := []string{core_domain.RoleAdmin.String()}
	warehouseStructureReadRoles := []string{
		core_domain.RoleWorker.String(),
		core_domain.RoleLogist.String(),
		core_domain.RoleAdmin.String(),
	}

	return []core_http_server.Route{
		core_http_server.NewRoute(http.MethodGet, "/warehouses", h.ListWarehouses, nil),
		core_http_server.NewRoute(http.MethodGet, "/warehouses/{id}", h.GetWarehouse, nil),
		core_http_server.NewRoute(http.MethodPost, "/warehouses", h.CreateWarehouse, adminRoles),
		core_http_server.NewRoute(http.MethodPatch, "/warehouses/{id}", h.PatchWarehouse, adminRoles),
		core_http_server.NewRoute(http.MethodPatch, "/warehouses/{id}/activate", h.ActivateWarehouse, adminRoles),
		core_http_server.NewRoute(http.MethodDelete, "/warehouses/{id}", h.DeactivateWarehouse, adminRoles),

		core_http_server.NewRoute(http.MethodGet, "/storage-zones", h.ListStorageZones, warehouseStructureReadRoles),
		core_http_server.NewRoute(http.MethodPost, "/storage-zones", h.CreateStorageZone, adminRoles),
		core_http_server.NewRoute(http.MethodPatch, "/storage-zones/{id}", h.PatchStorageZone, adminRoles),
		core_http_server.NewRoute(http.MethodPatch, "/storage-zones/{id}/activate", h.ActivateStorageZone, adminRoles),
		core_http_server.NewRoute(http.MethodDelete, "/storage-zones/{id}", h.DeactivateStorageZone, adminRoles),

		core_http_server.NewRoute(http.MethodGet, "/gates", h.ListGates, warehouseStructureReadRoles),
		core_http_server.NewRoute(http.MethodPost, "/gates", h.CreateGate, adminRoles),
		core_http_server.NewRoute(http.MethodPatch, "/gates/{id}", h.PatchGate, adminRoles),
		core_http_server.NewRoute(http.MethodPatch, "/gates/{id}/activate", h.ActivateGate, adminRoles),
		core_http_server.NewRoute(http.MethodDelete, "/gates/{id}", h.DeactivateGate, adminRoles),

		core_http_server.NewRoute(http.MethodGet, "/product-types", h.ListProductTypes, nil),
		core_http_server.NewRoute(http.MethodGet, "/cargo-place-types", h.ListCargoPlaceTypes, nil),
	}
}
