package pickupcalendar_transport_http

import (
	"context"
	"net/http"

	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
	core_logger "github.com/pkpal-uhobp/fulfillment-app/internal/core/logger"
	core_http_server "github.com/pkpal-uhobp/fulfillment-app/internal/core/transport/http/server"
	pickupcalendar_service "github.com/pkpal-uhobp/fulfillment-app/internal/features/pickupcalendar/service"
)

type PickupCalendarHTTPHandler struct {
	log                   *core_logger.Logger
	pickupCalendarService PickupCalendarService
}

type PickupCalendarService interface {
	GetCalendar(
		ctx context.Context,
		actorID int64,
		actorRole string,
		filter pickupcalendar_service.CalendarFilter,
	) ([]pickupcalendar_service.PickupCalendarDayDTO, error)
	BlockDate(
		ctx context.Context,
		actorID int64,
		actorRole string,
		input pickupcalendar_service.BlockDateInput,
	) (pickupcalendar_service.PickupCalendarBlockDTO, error)
	UnblockDate(
		ctx context.Context,
		actorID int64,
		actorRole string,
		blockID int64,
	) error
	SetCapacity(
		ctx context.Context,
		actorID int64,
		actorRole string,
		input pickupcalendar_service.SetCapacityInput,
	) (pickupcalendar_service.PickupCapacityDTO, error)
}

func NewPickupCalendarHTTPHandler(
	log *core_logger.Logger,
	pickupCalendarService PickupCalendarService,
) *PickupCalendarHTTPHandler {
	return &PickupCalendarHTTPHandler{
		log:                   log,
		pickupCalendarService: pickupCalendarService,
	}
}

func (h *PickupCalendarHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		core_http_server.NewRoute(
			http.MethodGet,
			"/pickup-calendar",
			h.GetCalendar,
			[]string{
				core_domain.RoleClient.String(),
				core_domain.RoleLogist.String(),
				core_domain.RoleAdmin.String(),
			},
		),
		core_http_server.NewRoute(
			http.MethodPost,
			"/pickup-calendar/blocks",
			h.BlockDate,
			[]string{
				core_domain.RoleLogist.String(),
				core_domain.RoleAdmin.String(),
			},
		),
		core_http_server.NewRoute(
			http.MethodDelete,
			"/pickup-calendar/blocks/{id}",
			h.UnblockDate,
			[]string{
				core_domain.RoleLogist.String(),
				core_domain.RoleAdmin.String(),
			},
		),
		core_http_server.NewRoute(
			http.MethodPatch,
			"/pickup-calendar/capacity",
			h.SetCapacity,
			[]string{
				core_domain.RoleLogist.String(),
				core_domain.RoleAdmin.String(),
			},
		),
	}
}
