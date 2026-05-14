package warehouses_service

import "context"

func (s *WarehousesService) ListProductTypes(
	ctx context.Context,
) ([]ProductTypeDTO, error) {
	productTypes, err := s.repo.ListProductTypes(ctx)
	if err != nil {
		return nil, err
	}

	return toProductTypeDTOs(productTypes), nil
}

func (s *WarehousesService) ListCargoPlaceTypes(
	ctx context.Context,
) ([]CargoPlaceTypeDTO, error) {
	cargoPlaceTypes, err := s.repo.ListCargoPlaceTypes(ctx)
	if err != nil {
		return nil, err
	}

	return toCargoPlaceTypeDTOs(cargoPlaceTypes), nil
}
