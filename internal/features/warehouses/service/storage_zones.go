package warehouses_service

import (
	"context"
	"fmt"
	"strings"

	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
	core_errors "github.com/pkpal-uhobp/fulfillment-app/internal/core/errors"
)

func (s *WarehousesService) ListStorageZones(
	ctx context.Context,
	warehouseID int64,
) ([]StorageZoneDTO, error) {
	if warehouseID < 0 {
		return nil, fmt.Errorf("%w: invalid warehouse id", core_errors.ErrInvalidArgument)
	}

	zones, err := s.repo.ListStorageZones(ctx, warehouseID)
	if err != nil {
		return nil, err
	}

	return toStorageZoneDTOs(zones), nil
}

func (s *WarehousesService) CreateStorageZone(
	ctx context.Context,
	input CreateStorageZoneInput,
) (StorageZoneDTO, error) {
	var description *string
	if strings.TrimSpace(input.Description) != "" {
		value := strings.TrimSpace(input.Description)
		description = &value
	}

	zone, err := core_domain.NewStorageZone(
		input.WarehouseID,
		input.Name,
		description,
	)
	if err != nil {
		return StorageZoneDTO{}, err
	}

	created, err := s.repo.CreateStorageZone(ctx, zone)
	if err != nil {
		return StorageZoneDTO{}, err
	}

	return toStorageZoneDTO(created), nil
}

func (s *WarehousesService) PatchStorageZone(
	ctx context.Context,
	zoneID int64,
	input PatchStorageZoneInput,
) (StorageZoneDTO, error) {
	if zoneID <= 0 {
		return StorageZoneDTO{}, fmt.Errorf("%w: invalid storage zone id", core_errors.ErrInvalidArgument)
	}

	var patch core_domain.StorageZonePatch

	if input.Name != nil {
		value := strings.TrimSpace(*input.Name)
		if value == "" {
			return StorageZoneDTO{}, fmt.Errorf("%w: storage zone name is required", core_errors.ErrInvalidArgument)
		}
		patch.Name = &value
	}

	if input.Description != nil {
		value := strings.TrimSpace(*input.Description)
		if value != "" {
			patch.Description = &value
		}
	}

	patch.IsActive = input.IsActive

	zone, err := s.repo.PatchStorageZone(ctx, zoneID, patch)
	if err != nil {
		return StorageZoneDTO{}, err
	}

	return toStorageZoneDTO(zone), nil
}
