package pickupcalendar_service

import core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"

func mapCalendarDaysToDTO(days []core_domain.PickupCalendarDay) []PickupCalendarDayDTO {
	result := make([]PickupCalendarDayDTO, 0, len(days))
	for _, day := range days {
		result = append(result, mapCalendarDayToDTO(day))
	}
	return result
}

func mapCalendarDayToDTO(day core_domain.PickupCalendarDay) PickupCalendarDayDTO {
	var block *PickupCalendarBlockDTO
	if day.Block != nil {
		value := mapBlockToDTO(*day.Block)
		block = &value
	}
	return PickupCalendarDayDTO{
		WarehouseID:   day.WarehouseID,
		Date:          day.Date,
		MaxOrders:     day.MaxOrders,
		CurrentOrders: day.CurrentOrders,
		IsClosed:      day.IsClosed,
		Block:         block,
	}
}

func mapBlockToDTO(block core_domain.PickupCalendarBlock) PickupCalendarBlockDTO {
	return PickupCalendarBlockDTO{
		ID:          block.ID,
		WarehouseID: block.WarehouseID,
		BlockedDate: block.BlockedDate,
		Reason:      block.Reason,
		CreatedBy:   block.CreatedBy,
		CreatedAt:   block.CreatedAt,
	}
}

func mapCapacityToDTO(capacity core_domain.PickupCapacity) PickupCapacityDTO {
	return PickupCapacityDTO{
		ID:            capacity.ID,
		WarehouseID:   capacity.WarehouseID,
		PickupDate:    capacity.PickupDate,
		MaxOrders:     capacity.MaxOrders,
		CurrentOrders: capacity.CurrentOrders,
		IsClosed:      capacity.IsClosed,
	}
}
