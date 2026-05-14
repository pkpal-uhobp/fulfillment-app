package warehouses_service

import (
	"context"
	"fmt"
	"strings"

	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
	core_errors "github.com/pkpal-uhobp/fulfillment-app/internal/core/errors"
)

func (s *WarehousesService) ListWarehouses(
	ctx context.Context,
	filter WarehouseFilter,
) ([]WarehouseDTO, error) {
	filter.WarehouseType = strings.TrimSpace(filter.WarehouseType)
	filter.Marketplace = strings.TrimSpace(filter.Marketplace)
	filter.City = strings.TrimSpace(filter.City)

	if filter.WarehouseType != "" {
		warehouseType := core_domain.WarehouseType(filter.WarehouseType)
		if !warehouseType.IsValid() {
			return nil, fmt.Errorf("%w: invalid warehouse type", core_errors.ErrInvalidArgument)
		}
	}

	warehouses, err := s.repo.ListWarehouses(ctx, filter)
	if err != nil {
		return nil, err
	}

	return toWarehouseDTOs(warehouses), nil
}

func (s *WarehousesService) GetWarehouse(
	ctx context.Context,
	warehouseID int64,
) (WarehouseDTO, error) {
	if warehouseID <= 0 {
		return WarehouseDTO{}, fmt.Errorf("%w: invalid warehouse id", core_errors.ErrInvalidArgument)
	}

	warehouse, err := s.repo.GetWarehouseByID(ctx, warehouseID)
	if err != nil {
		return WarehouseDTO{}, err
	}

	return toWarehouseDTO(warehouse), nil
}

func (s *WarehousesService) CreateWarehouse(
	ctx context.Context,
	input CreateWarehouseInput,
) (WarehouseDTO, error) {
	var marketplace *string
	if strings.TrimSpace(input.Marketplace) != "" {
		value := strings.TrimSpace(input.Marketplace)
		marketplace = &value
	}

	warehouse, err := core_domain.NewWarehouse(
		input.Name,
		core_domain.WarehouseType(strings.TrimSpace(input.WarehouseType)),
		marketplace,
		input.City,
		input.Address,
	)
	if err != nil {
		return WarehouseDTO{}, err
	}

	created, err := s.repo.CreateWarehouse(ctx, warehouse)
	if err != nil {
		return WarehouseDTO{}, err
	}

	return toWarehouseDTO(created), nil
}

func (s *WarehousesService) PatchWarehouse(
	ctx context.Context,
	warehouseID int64,
	input PatchWarehouseInput,
) (WarehouseDTO, error) {
	if warehouseID <= 0 {
		return WarehouseDTO{}, fmt.Errorf("%w: invalid warehouse id", core_errors.ErrInvalidArgument)
	}

	var patch core_domain.WarehousePatch

	if input.Name != nil {
		value := strings.TrimSpace(*input.Name)
		if value == "" {
			return WarehouseDTO{}, fmt.Errorf("%w: warehouse name is required", core_errors.ErrInvalidArgument)
		}
		patch.Name = &value
	}

	if input.WarehouseType != nil {
		warehouseType := core_domain.WarehouseType(strings.TrimSpace(*input.WarehouseType))
		if !warehouseType.IsValid() {
			return WarehouseDTO{}, fmt.Errorf("%w: invalid warehouse type", core_errors.ErrInvalidArgument)
		}
		patch.WarehouseType = &warehouseType
	}

	if input.Marketplace != nil {
		value := strings.TrimSpace(*input.Marketplace)
		if value != "" {
			patch.Marketplace = &value
		}
	}

	if input.City != nil {
		value := strings.TrimSpace(*input.City)
		if value == "" {
			return WarehouseDTO{}, fmt.Errorf("%w: warehouse city is required", core_errors.ErrInvalidArgument)
		}
		patch.City = &value
	}

	if input.Address != nil {
		value := strings.TrimSpace(*input.Address)
		if value == "" {
			return WarehouseDTO{}, fmt.Errorf("%w: warehouse address is required", core_errors.ErrInvalidArgument)
		}
		patch.Address = &value
	}

	patch.IsActive = input.IsActive

	warehouse, err := s.repo.PatchWarehouse(ctx, warehouseID, patch)
	if err != nil {
		return WarehouseDTO{}, err
	}

	return toWarehouseDTO(warehouse), nil
}

func (s *WarehousesService) DeactivateWarehouse(
	ctx context.Context,
	warehouseID int64,
) error {
	if warehouseID <= 0 {
		return fmt.Errorf("%w: invalid warehouse id", core_errors.ErrInvalidArgument)
	}

	return s.repo.DeactivateWarehouse(ctx, warehouseID)
}
