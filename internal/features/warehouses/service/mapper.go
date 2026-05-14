package warehouses_service

import core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"

func toWarehouseDTO(warehouse core_domain.Warehouse) WarehouseDTO {
	dto := WarehouseDTO{
		ID:            warehouse.ID,
		Name:          warehouse.Name,
		WarehouseType: warehouse.WarehouseType.String(),
		City:          warehouse.City,
		Address:       warehouse.Address,
		IsActive:      warehouse.IsActive,
		CreatedAt:     warehouse.CreatedAt,
	}

	if warehouse.Marketplace != nil {
		dto.Marketplace = *warehouse.Marketplace
	}

	return dto
}

func toWarehouseDTOs(warehouses []core_domain.Warehouse) []WarehouseDTO {
	result := make([]WarehouseDTO, 0, len(warehouses))
	for _, warehouse := range warehouses {
		result = append(result, toWarehouseDTO(warehouse))
	}

	return result
}

func toStorageZoneDTO(zone core_domain.StorageZone) StorageZoneDTO {
	dto := StorageZoneDTO{
		ID:          zone.ID,
		WarehouseID: zone.WarehouseID,
		Name:        zone.Name,
		IsActive:    zone.IsActive,
	}

	if zone.Description != nil {
		dto.Description = *zone.Description
	}

	return dto
}

func toStorageZoneDTOs(zones []core_domain.StorageZone) []StorageZoneDTO {
	result := make([]StorageZoneDTO, 0, len(zones))
	for _, zone := range zones {
		result = append(result, toStorageZoneDTO(zone))
	}

	return result
}

func toGateDTO(gate core_domain.Gate) GateDTO {
	return GateDTO{
		ID:          gate.ID,
		WarehouseID: gate.WarehouseID,
		Name:        gate.Name,
		IsActive:    gate.IsActive,
	}
}

func toGateDTOs(gates []core_domain.Gate) []GateDTO {
	result := make([]GateDTO, 0, len(gates))
	for _, gate := range gates {
		result = append(result, toGateDTO(gate))
	}

	return result
}

func toProductTypeDTO(productType core_domain.ProductType) ProductTypeDTO {
	dto := ProductTypeDTO{
		ID:       productType.ID,
		Name:     productType.Name,
		IsActive: productType.IsActive,
	}

	if productType.Description != nil {
		dto.Description = *productType.Description
	}

	return dto
}

func toProductTypeDTOs(productTypes []core_domain.ProductType) []ProductTypeDTO {
	result := make([]ProductTypeDTO, 0, len(productTypes))
	for _, productType := range productTypes {
		result = append(result, toProductTypeDTO(productType))
	}

	return result
}

func toCargoPlaceTypeDTO(cargoPlaceType core_domain.CargoPlaceType) CargoPlaceTypeDTO {
	dto := CargoPlaceTypeDTO{
		ID:       cargoPlaceType.ID,
		Name:     cargoPlaceType.Name,
		IsActive: cargoPlaceType.IsActive,
	}

	if cargoPlaceType.Description != nil {
		dto.Description = *cargoPlaceType.Description
	}

	return dto
}

func toCargoPlaceTypeDTOs(cargoPlaceTypes []core_domain.CargoPlaceType) []CargoPlaceTypeDTO {
	result := make([]CargoPlaceTypeDTO, 0, len(cargoPlaceTypes))
	for _, cargoPlaceType := range cargoPlaceTypes {
		result = append(result, toCargoPlaceTypeDTO(cargoPlaceType))
	}

	return result
}
