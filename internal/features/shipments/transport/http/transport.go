package shipments_transport_http

import (
	"context"
	"net/http"

	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
	core_logger "github.com/pkpal-uhobp/fulfillment-app/internal/core/logger"
	core_http_server "github.com/pkpal-uhobp/fulfillment-app/internal/core/transport/http/server"
	shipments_service "github.com/pkpal-uhobp/fulfillment-app/internal/features/shipments/service"
)

type ShipmentsHTTPHandler struct {
	log              *core_logger.Logger
	shipmentsService ShipmentsService
}

type ShipmentsService interface {
	CreateShipment(ctx context.Context, actorID int64, actorRole string, input shipments_service.CreateShipmentInput) (shipments_service.ShipmentDTO, error)
	ListShipments(ctx context.Context, actorID int64, actorRole string, filter shipments_service.ShipmentFilter) ([]shipments_service.ShipmentDTO, error)
	GetShipment(ctx context.Context, shipmentID int64, actorID int64, actorRole string) (shipments_service.ShipmentDTO, error)
	AddShipmentItem(ctx context.Context, shipmentID int64, actorID int64, actorRole string, input shipments_service.AddShipmentItemInput) (shipments_service.ShipmentDTO, error)
	RemoveShipmentItem(ctx context.Context, shipmentID int64, cargoItemID int64, actorID int64, actorRole string) error
	UpdateShipmentStatus(ctx context.Context, shipmentID int64, actorID int64, actorRole string, input shipments_service.UpdateShipmentStatusInput) (shipments_service.ShipmentDTO, error)
}

func NewShipmentsHTTPHandler(
	log *core_logger.Logger,
	shipmentsService ShipmentsService,
) *ShipmentsHTTPHandler {
	return &ShipmentsHTTPHandler{
		log:              log,
		shipmentsService: shipmentsService,
	}
}

func (h *ShipmentsHTTPHandler) Routes() []core_http_server.Route {
	roles := []string{core_domain.RoleLogist.String(), core_domain.RoleAdmin.String()}
	return []core_http_server.Route{
		core_http_server.NewRoute(http.MethodPost, "/shipments", h.CreateShipment, roles),
		core_http_server.NewRoute(http.MethodGet, "/shipments", h.ListShipments, roles),
		core_http_server.NewRoute(http.MethodGet, "/shipments/{id}", h.GetShipment, roles),
		core_http_server.NewRoute(http.MethodPost, "/shipments/{id}/items", h.AddShipmentItem, roles),
		core_http_server.NewRoute(http.MethodDelete, "/shipments/{id}/items/{cargo_item_id}", h.RemoveShipmentItem, roles),
		core_http_server.NewRoute(http.MethodPatch, "/shipments/{id}/status", h.UpdateShipmentStatus, roles),
	}
}
