package warehouses_service

import (
	"context"
	"fmt"
	"strings"

	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
	core_errors "github.com/pkpal-uhobp/fulfillment-app/internal/core/errors"
)

func (s *WarehousesService) ListGates(
	ctx context.Context,
	warehouseID int64,
) ([]GateDTO, error) {
	if warehouseID < 0 {
		return nil, fmt.Errorf("%w: invalid warehouse id", core_errors.ErrInvalidArgument)
	}

	gates, err := s.repo.ListGates(ctx, warehouseID)
	if err != nil {
		return nil, err
	}

	return toGateDTOs(gates), nil
}

func (s *WarehousesService) CreateGate(
	ctx context.Context,
	input CreateGateInput,
) (GateDTO, error) {
	gate, err := core_domain.NewGate(
		input.WarehouseID,
		input.Name,
	)
	if err != nil {
		return GateDTO{}, err
	}

	created, err := s.repo.CreateGate(ctx, gate)
	if err != nil {
		return GateDTO{}, err
	}

	return toGateDTO(created), nil
}

func (s *WarehousesService) PatchGate(
	ctx context.Context,
	gateID int64,
	input PatchGateInput,
) (GateDTO, error) {
	if gateID <= 0 {
		return GateDTO{}, fmt.Errorf("%w: invalid gate id", core_errors.ErrInvalidArgument)
	}

	var patch core_domain.GatePatch

	if input.Name != nil {
		value := strings.TrimSpace(*input.Name)
		if value == "" {
			return GateDTO{}, fmt.Errorf("%w: gate name is required", core_errors.ErrInvalidArgument)
		}
		patch.Name = &value
	}

	patch.IsActive = input.IsActive

	gate, err := s.repo.PatchGate(ctx, gateID, patch)
	if err != nil {
		return GateDTO{}, err
	}

	return toGateDTO(gate), nil
}
