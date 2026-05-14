package cargoitems_service

import (
	"context"
)

func (s *CargoItemsService) GetCargoItemHistory(
	ctx context.Context,
	cargoItemID int64,
	actorID int64,
	actorRole string,
) ([]CargoStatusHistoryDTO, error) {
	if _, err := s.getCargoItemWithAccess(ctx, cargoItemID, actorID, actorRole); err != nil {
		return nil, err
	}

	history, err := s.repo.ListCargoStatusHistory(ctx, cargoItemID)
	if err != nil {
		return nil, err
	}

	return mapCargoStatusHistoryListToDTO(history), nil
}
