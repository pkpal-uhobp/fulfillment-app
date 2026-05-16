package shipments_service

import (
	"context"
	"strings"

	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
)

func (s *ShipmentsService) ListShipments(
	ctx context.Context,
	actorID int64,
	actorRole string,
	filter ShipmentFilter,
) ([]ShipmentDTO, error) {
	if err := requireLogistOrAdmin(actorRole); err != nil {
		return nil, err
	}
	filter.Status = strings.TrimSpace(filter.Status)
	if filter.Status != "" {
		if _, err := validateShipmentStatus(filter.Status); err != nil {
			return nil, err
		}
	}
	shipments, err := s.repo.ListShipments(ctx, core_domain.ShipmentFilter{
		Status:                 filter.Status,
		DestinationWarehouseID: filter.DestinationWarehouseID,
		GateID:                 filter.GateID,
		Date:                   strings.TrimSpace(filter.Date),
		Page:                   filter.Page,
		Limit:                  filter.Limit,
	})
	if err != nil {
		return nil, err
	}
	return mapShipmentsToDTO(shipments), nil
}
