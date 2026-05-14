package cargoitems_transport_http

import (
	"context"
	"net/http"

	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
	core_logger "github.com/pkpal-uhobp/fulfillment-app/internal/core/logger"
	core_http_server "github.com/pkpal-uhobp/fulfillment-app/internal/core/transport/http/server"
	cargoitems_service "github.com/pkpal-uhobp/fulfillment-app/internal/features/cargoitems/service"
)

type CargoItemsHTTPHandler struct {
	log               *core_logger.Logger
	cargoItemsService CargoItemsService
}

type CargoItemsService interface {
	CreateCargoItem(
		ctx context.Context,
		actorID int64,
		actorRole string,
		input cargoitems_service.CreateCargoItemInput,
	) (cargoitems_service.CargoItemDTO, error)

	ListCargoItems(
		ctx context.Context,
		actorID int64,
		actorRole string,
		filter cargoitems_service.CargoItemFilter,
	) ([]cargoitems_service.CargoItemDTO, error)

	GetCargoItem(
		ctx context.Context,
		cargoItemID int64,
		actorID int64,
		actorRole string,
	) (cargoitems_service.CargoItemDTO, error)

	GetCargoItemHistory(
		ctx context.Context,
		cargoItemID int64,
		actorID int64,
		actorRole string,
	) ([]cargoitems_service.CargoStatusHistoryDTO, error)

	UpdateCargoItemStatus(
		ctx context.Context,
		cargoItemID int64,
		actorID int64,
		actorRole string,
		input cargoitems_service.UpdateCargoItemStatusInput,
	) (cargoitems_service.CargoItemDTO, error)

	AssignStorageZone(
		ctx context.Context,
		cargoItemID int64,
		actorID int64,
		actorRole string,
		input cargoitems_service.AssignStorageZoneInput,
	) (cargoitems_service.CargoItemDTO, error)

	AssignGate(
		ctx context.Context,
		cargoItemID int64,
		actorID int64,
		actorRole string,
		input cargoitems_service.AssignGateInput,
	) (cargoitems_service.CargoItemDTO, error)
}

func NewCargoItemsHTTPHandler(
	log *core_logger.Logger,
	cargoItemsService CargoItemsService,
) *CargoItemsHTTPHandler {
	return &CargoItemsHTTPHandler{
		log:               log,
		cargoItemsService: cargoItemsService,
	}
}

func (h *CargoItemsHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		core_http_server.NewRoute(
			http.MethodPost,
			"/orders/{id}/cargo-items",
			h.CreateCargoItem,
			[]string{
				core_domain.RoleWorker.String(),
				core_domain.RoleAdmin.String(),
			},
		),
		core_http_server.NewRoute(
			http.MethodGet,
			"/cargo-items",
			h.ListCargoItems,
			[]string{
				core_domain.RoleClient.String(),
				core_domain.RoleWorker.String(),
				core_domain.RoleLogist.String(),
				core_domain.RoleAdmin.String(),
			},
		),
		core_http_server.NewRoute(
			http.MethodGet,
			"/cargo-items/{id}",
			h.GetCargoItem,
			[]string{
				core_domain.RoleClient.String(),
				core_domain.RoleWorker.String(),
				core_domain.RoleLogist.String(),
				core_domain.RoleAdmin.String(),
			},
		),
		core_http_server.NewRoute(
			http.MethodGet,
			"/cargo-items/{id}/history",
			h.GetCargoItemHistory,
			[]string{
				core_domain.RoleClient.String(),
				core_domain.RoleWorker.String(),
				core_domain.RoleLogist.String(),
				core_domain.RoleAdmin.String(),
			},
		),
		core_http_server.NewRoute(
			http.MethodPatch,
			"/cargo-items/{id}/status",
			h.UpdateCargoItemStatus,
			[]string{
				core_domain.RoleWorker.String(),
				core_domain.RoleLogist.String(),
				core_domain.RoleAdmin.String(),
			},
		),
		core_http_server.NewRoute(
			http.MethodPatch,
			"/cargo-items/{id}/assign-zone",
			h.AssignStorageZone,
			[]string{
				core_domain.RoleLogist.String(),
				core_domain.RoleAdmin.String(),
			},
		),
		core_http_server.NewRoute(
			http.MethodPatch,
			"/cargo-items/{id}/assign-gate",
			h.AssignGate,
			[]string{
				core_domain.RoleLogist.String(),
				core_domain.RoleAdmin.String(),
			},
		),
	}
}
