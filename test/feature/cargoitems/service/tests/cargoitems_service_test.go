package cargoitems_service_tests

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
	core_errors "github.com/pkpal-uhobp/fulfillment-app/internal/core/errors"
	cargoitems_service "github.com/pkpal-uhobp/fulfillment-app/internal/features/cargoitems/service"
)

func TestCreateCargoItemSuccessWithManualQRCode(t *testing.T) {
	now := time.Now().UTC()
	repo := &fakeCargoItemsRepository{
		getOrderCargoPlaceForOrderFn: func(ctx context.Context, orderID int64, orderCargoPlaceID int64) (core_domain.OrderCargoPlace, core_domain.OrderStatus, error) {
			if orderID != 10 {
				t.Fatalf("orderID = %d, want 10", orderID)
			}
			if orderCargoPlaceID != 20 {
				t.Fatalf("orderCargoPlaceID = %d, want 20", orderCargoPlaceID)
			}
			return core_domain.OrderCargoPlace{
				ID:               orderCargoPlaceID,
				OrderID:          orderID,
				CargoPlaceTypeID: 3,
				Quantity:         2,
			}, core_domain.OrderStatusReceived, nil
		},
		countCargoItemsByOrderCargoPlaceFn: func(ctx context.Context, orderCargoPlaceID int64) (int, error) {
			if orderCargoPlaceID != 20 {
				t.Fatalf("orderCargoPlaceID = %d, want 20", orderCargoPlaceID)
			}
			return 1, nil
		},
		createCargoItemFn: func(ctx context.Context, item core_domain.CargoItem, changedBy int64) (core_domain.CargoItem, error) {
			if changedBy != 5 {
				t.Fatalf("changedBy = %d, want 5", changedBy)
			}
			if item.OrderID != 10 || item.OrderCargoPlaceID != 20 || item.CargoPlaceTypeID != 3 {
				t.Fatalf("unexpected item ids: %+v", item)
			}
			if item.QRCode != "QR-001" {
				t.Fatalf("QRCode = %q, want QR-001", item.QRCode)
			}
			if item.Status != core_domain.CargoItemStatusAccepted {
				t.Fatalf("Status = %q, want accepted", item.Status)
			}
			if item.Comment == nil || *item.Comment != "Принято без повреждений" {
				t.Fatalf("Comment = %v, want trimmed comment", item.Comment)
			}
			item.ID = 100
			item.CreatedAt = now
			item.UpdatedAt = now
			return item, nil
		},
	}
	service := cargoitems_service.NewCargoItemsService(repo)

	dto, err := service.CreateCargoItem(context.Background(), 5, core_domain.RoleWorker.String(), cargoitems_service.CreateCargoItemInput{
		OrderID:           10,
		OrderCargoPlaceID: 20,
		QRCode:            ptr(" QR-001 "),
		Comment:           ptr(" Принято без повреждений "),
	})
	if err != nil {
		t.Fatalf("CreateCargoItem returned error: %v", err)
	}
	if dto.ID != 100 {
		t.Fatalf("ID = %d, want 100", dto.ID)
	}
	if dto.Status != core_domain.CargoItemStatusAccepted.String() {
		t.Fatalf("Status = %q, want accepted", dto.Status)
	}
}

func TestCreateCargoItemGeneratesQRCodeWhenMissing(t *testing.T) {
	repo := &fakeCargoItemsRepository{
		getOrderCargoPlaceForOrderFn: func(ctx context.Context, orderID int64, orderCargoPlaceID int64) (core_domain.OrderCargoPlace, core_domain.OrderStatus, error) {
			return core_domain.OrderCargoPlace{ID: orderCargoPlaceID, OrderID: orderID, CargoPlaceTypeID: 3, Quantity: 1}, core_domain.OrderStatusReceived, nil
		},
		countCargoItemsByOrderCargoPlaceFn: func(ctx context.Context, orderCargoPlaceID int64) (int, error) {
			return 0, nil
		},
		createCargoItemFn: func(ctx context.Context, item core_domain.CargoItem, changedBy int64) (core_domain.CargoItem, error) {
			if !strings.HasPrefix(item.QRCode, "CI-") {
				t.Fatalf("QRCode = %q, want CI- prefix", item.QRCode)
			}
			if len(item.QRCode) != 19 {
				t.Fatalf("QRCode len = %d, want 19", len(item.QRCode))
			}
			item.ID = 1
			return item, nil
		},
	}
	service := cargoitems_service.NewCargoItemsService(repo)

	_, err := service.CreateCargoItem(context.Background(), 5, core_domain.RoleWorker.String(), cargoitems_service.CreateCargoItemInput{
		OrderID:           10,
		OrderCargoPlaceID: 20,
	})
	if err != nil {
		t.Fatalf("CreateCargoItem returned error: %v", err)
	}
}

func TestCreateCargoItemRejectsForbiddenRole(t *testing.T) {
	service := cargoitems_service.NewCargoItemsService(&fakeCargoItemsRepository{})

	_, err := service.CreateCargoItem(context.Background(), 5, core_domain.RoleClient.String(), cargoitems_service.CreateCargoItemInput{
		OrderID:           10,
		OrderCargoPlaceID: 20,
	})
	if !errors.Is(err, core_errors.ErrForbidden) {
		t.Fatalf("err = %v, want ErrForbidden", err)
	}
}

func TestCreateCargoItemRejectsExceededDeclaredQuantity(t *testing.T) {
	repo := &fakeCargoItemsRepository{
		getOrderCargoPlaceForOrderFn: func(ctx context.Context, orderID int64, orderCargoPlaceID int64) (core_domain.OrderCargoPlace, core_domain.OrderStatus, error) {
			return core_domain.OrderCargoPlace{ID: orderCargoPlaceID, OrderID: orderID, CargoPlaceTypeID: 3, Quantity: 2}, core_domain.OrderStatusReceived, nil
		},
		countCargoItemsByOrderCargoPlaceFn: func(ctx context.Context, orderCargoPlaceID int64) (int, error) {
			return 2, nil
		},
	}
	service := cargoitems_service.NewCargoItemsService(repo)

	_, err := service.CreateCargoItem(context.Background(), 5, core_domain.RoleWorker.String(), cargoitems_service.CreateCargoItemInput{
		OrderID:           10,
		OrderCargoPlaceID: 20,
	})
	if !errors.Is(err, core_errors.ErrConflict) {
		t.Fatalf("err = %v, want ErrConflict", err)
	}
}

func TestCreateCargoItemRejectsTerminalOrder(t *testing.T) {
	repo := &fakeCargoItemsRepository{
		getOrderCargoPlaceForOrderFn: func(ctx context.Context, orderID int64, orderCargoPlaceID int64) (core_domain.OrderCargoPlace, core_domain.OrderStatus, error) {
			return core_domain.OrderCargoPlace{ID: orderCargoPlaceID, OrderID: orderID, CargoPlaceTypeID: 3, Quantity: 1}, core_domain.OrderStatusCancelled, nil
		},
	}
	service := cargoitems_service.NewCargoItemsService(repo)

	_, err := service.CreateCargoItem(context.Background(), 5, core_domain.RoleWorker.String(), cargoitems_service.CreateCargoItemInput{
		OrderID:           10,
		OrderCargoPlaceID: 20,
	})
	if !errors.Is(err, core_errors.ErrConflict) {
		t.Fatalf("err = %v, want ErrConflict", err)
	}
}

func TestListCargoItemsNormalizesFilterAndRestrictsClient(t *testing.T) {
	repo := &fakeCargoItemsRepository{
		listCargoItemsFn: func(ctx context.Context, filter core_domain.CargoItemFilter) ([]core_domain.CargoItem, error) {
			if filter.ClientID == nil || *filter.ClientID != 42 {
				t.Fatalf("ClientID = %v, want 42", filter.ClientID)
			}
			if filter.Status != core_domain.CargoItemStatusStored.String() {
				t.Fatalf("Status = %q, want stored", filter.Status)
			}
			if filter.QRCode != "QR-777" {
				t.Fatalf("QRCode = %q, want QR-777", filter.QRCode)
			}
			if filter.Page != 1 || filter.Limit != 100 {
				t.Fatalf("Page/Limit = %d/%d, want 1/100", filter.Page, filter.Limit)
			}
			return []core_domain.CargoItem{{ID: 7, QRCode: "QR-777", Status: core_domain.CargoItemStatusStored}}, nil
		},
	}
	service := cargoitems_service.NewCargoItemsService(repo)

	items, err := service.ListCargoItems(context.Background(), 42, core_domain.RoleClient.String(), cargoitems_service.CargoItemFilter{
		Status: " stored ",
		QRCode: " QR-777 ",
		Limit:  999,
	})
	if err != nil {
		t.Fatalf("ListCargoItems returned error: %v", err)
	}
	if len(items) != 1 || items[0].ID != 7 {
		t.Fatalf("unexpected items: %+v", items)
	}
}

func TestListCargoItemsRejectsInvalidStatus(t *testing.T) {
	service := cargoitems_service.NewCargoItemsService(&fakeCargoItemsRepository{})

	_, err := service.ListCargoItems(context.Background(), 1, core_domain.RoleAdmin.String(), cargoitems_service.CargoItemFilter{Status: "wrong"})
	if !errors.Is(err, core_errors.ErrInvalidArgument) {
		t.Fatalf("err = %v, want ErrInvalidArgument", err)
	}
}

func TestGetCargoItemRejectsForeignClient(t *testing.T) {
	repo := &fakeCargoItemsRepository{
		getCargoItemFn: func(ctx context.Context, cargoItemID int64) (core_domain.CargoItem, error) {
			return core_domain.CargoItem{ID: cargoItemID, Status: core_domain.CargoItemStatusAccepted}, nil
		},
		clientOwnsCargoItemFn: func(ctx context.Context, cargoItemID int64, clientID int64) (bool, error) {
			if cargoItemID != 7 || clientID != 42 {
				t.Fatalf("cargoItemID/clientID = %d/%d, want 7/42", cargoItemID, clientID)
			}
			return false, nil
		},
	}
	service := cargoitems_service.NewCargoItemsService(repo)

	_, err := service.GetCargoItem(context.Background(), 7, 42, core_domain.RoleClient.String())
	if !errors.Is(err, core_errors.ErrForbidden) {
		t.Fatalf("err = %v, want ErrForbidden", err)
	}
}

func TestUpdateCargoItemStatusSuccessAcceptedToStored(t *testing.T) {
	zoneID := int64(9)
	repo := &fakeCargoItemsRepository{
		getCargoItemFn: func(ctx context.Context, cargoItemID int64) (core_domain.CargoItem, error) {
			return core_domain.CargoItem{ID: cargoItemID, Status: core_domain.CargoItemStatusAccepted, StorageZoneID: &zoneID}, nil
		},
		updateCargoItemStatusFn: func(ctx context.Context, cargoItemID int64, status string, changedBy int64, comment *string) (core_domain.CargoItem, error) {
			if cargoItemID != 7 || status != core_domain.CargoItemStatusStored.String() || changedBy != 5 {
				t.Fatalf("unexpected update args: id=%d status=%q changedBy=%d", cargoItemID, status, changedBy)
			}
			if comment == nil || *comment != "Перемещено в зону" {
				t.Fatalf("comment = %v, want trimmed comment", comment)
			}
			return core_domain.CargoItem{ID: cargoItemID, Status: core_domain.CargoItemStatusStored, StorageZoneID: &zoneID}, nil
		},
	}
	service := cargoitems_service.NewCargoItemsService(repo)

	dto, err := service.UpdateCargoItemStatus(context.Background(), 7, 5, core_domain.RoleWorker.String(), cargoitems_service.UpdateCargoItemStatusInput{
		Status:  "stored",
		Comment: ptr(" Перемещено в зону "),
	})
	if err != nil {
		t.Fatalf("UpdateCargoItemStatus returned error: %v", err)
	}
	if dto.Status != core_domain.CargoItemStatusStored.String() {
		t.Fatalf("Status = %q, want stored", dto.Status)
	}
}

func TestUpdateCargoItemStatusRejectsMissingStorageZone(t *testing.T) {
	repo := &fakeCargoItemsRepository{
		getCargoItemFn: func(ctx context.Context, cargoItemID int64) (core_domain.CargoItem, error) {
			return core_domain.CargoItem{ID: cargoItemID, Status: core_domain.CargoItemStatusAccepted}, nil
		},
	}
	service := cargoitems_service.NewCargoItemsService(repo)

	_, err := service.UpdateCargoItemStatus(context.Background(), 7, 5, core_domain.RoleWorker.String(), cargoitems_service.UpdateCargoItemStatusInput{Status: "stored"})
	if !errors.Is(err, core_errors.ErrConflict) {
		t.Fatalf("err = %v, want ErrConflict", err)
	}
}

func TestUpdateCargoItemStatusRejectsInvalidTransition(t *testing.T) {
	gateID := int64(4)
	repo := &fakeCargoItemsRepository{
		getCargoItemFn: func(ctx context.Context, cargoItemID int64) (core_domain.CargoItem, error) {
			return core_domain.CargoItem{ID: cargoItemID, Status: core_domain.CargoItemStatusAccepted, GateID: &gateID}, nil
		},
	}
	service := cargoitems_service.NewCargoItemsService(repo)

	_, err := service.UpdateCargoItemStatus(context.Background(), 7, 5, core_domain.RoleWorker.String(), cargoitems_service.UpdateCargoItemStatusInput{Status: "shipped"})
	if !errors.Is(err, core_errors.ErrConflict) {
		t.Fatalf("err = %v, want ErrConflict", err)
	}
}

func TestAssignStorageZoneSuccess(t *testing.T) {
	zoneID := int64(9)
	repo := &fakeCargoItemsRepository{
		getCargoItemFn: func(ctx context.Context, cargoItemID int64) (core_domain.CargoItem, error) {
			return core_domain.CargoItem{ID: cargoItemID, Status: core_domain.CargoItemStatusAccepted}, nil
		},
		storageZoneBelongsToCargoOrderFn: func(ctx context.Context, cargoItemID int64, storageZoneID int64) (bool, error) {
			if cargoItemID != 7 || storageZoneID != zoneID {
				t.Fatalf("unexpected zone check args: %d/%d", cargoItemID, storageZoneID)
			}
			return true, nil
		},
		assignStorageZoneFn: func(ctx context.Context, cargoItemID int64, storageZoneID int64, changedBy int64, comment *string) (core_domain.CargoItem, error) {
			if cargoItemID != 7 || storageZoneID != zoneID || changedBy != 5 {
				t.Fatalf("unexpected assign args: %d/%d/%d", cargoItemID, storageZoneID, changedBy)
			}
			return core_domain.CargoItem{ID: cargoItemID, Status: core_domain.CargoItemStatusAccepted, StorageZoneID: &storageZoneID}, nil
		},
	}
	service := cargoitems_service.NewCargoItemsService(repo)

	dto, err := service.AssignStorageZone(context.Background(), 7, 5, core_domain.RoleLogist.String(), cargoitems_service.AssignStorageZoneInput{StorageZoneID: zoneID})
	if err != nil {
		t.Fatalf("AssignStorageZone returned error: %v", err)
	}
	if dto.StorageZoneID == nil || *dto.StorageZoneID != zoneID {
		t.Fatalf("StorageZoneID = %v, want %d", dto.StorageZoneID, zoneID)
	}
}

func TestAssignStorageZoneRejectsForeignZone(t *testing.T) {
	repo := &fakeCargoItemsRepository{
		getCargoItemFn: func(ctx context.Context, cargoItemID int64) (core_domain.CargoItem, error) {
			return core_domain.CargoItem{ID: cargoItemID, Status: core_domain.CargoItemStatusAccepted}, nil
		},
		storageZoneBelongsToCargoOrderFn: func(ctx context.Context, cargoItemID int64, storageZoneID int64) (bool, error) {
			return false, nil
		},
	}
	service := cargoitems_service.NewCargoItemsService(repo)

	_, err := service.AssignStorageZone(context.Background(), 7, 5, core_domain.RoleLogist.String(), cargoitems_service.AssignStorageZoneInput{StorageZoneID: 9})
	if !errors.Is(err, core_errors.ErrInvalidArgument) {
		t.Fatalf("err = %v, want ErrInvalidArgument", err)
	}
}

func TestAssignGateSuccess(t *testing.T) {
	gateID := int64(4)
	repo := &fakeCargoItemsRepository{
		getCargoItemFn: func(ctx context.Context, cargoItemID int64) (core_domain.CargoItem, error) {
			return core_domain.CargoItem{ID: cargoItemID, Status: core_domain.CargoItemStatusStored}, nil
		},
		gateBelongsToCargoOrderFn: func(ctx context.Context, cargoItemID int64, gateID int64) (bool, error) {
			if cargoItemID != 7 || gateID != 4 {
				t.Fatalf("unexpected gate check args: %d/%d", cargoItemID, gateID)
			}
			return true, nil
		},
		assignGateFn: func(ctx context.Context, cargoItemID int64, gateID int64, changedBy int64, comment *string) (core_domain.CargoItem, error) {
			return core_domain.CargoItem{ID: cargoItemID, Status: core_domain.CargoItemStatusStored, GateID: &gateID}, nil
		},
	}
	service := cargoitems_service.NewCargoItemsService(repo)

	dto, err := service.AssignGate(context.Background(), 7, 5, core_domain.RoleLogist.String(), cargoitems_service.AssignGateInput{GateID: gateID})
	if err != nil {
		t.Fatalf("AssignGate returned error: %v", err)
	}
	if dto.GateID == nil || *dto.GateID != gateID {
		t.Fatalf("GateID = %v, want %d", dto.GateID, gateID)
	}
}

func TestAssignGateRejectsWrongStatus(t *testing.T) {
	repo := &fakeCargoItemsRepository{
		getCargoItemFn: func(ctx context.Context, cargoItemID int64) (core_domain.CargoItem, error) {
			return core_domain.CargoItem{ID: cargoItemID, Status: core_domain.CargoItemStatusAccepted}, nil
		},
	}
	service := cargoitems_service.NewCargoItemsService(repo)

	_, err := service.AssignGate(context.Background(), 7, 5, core_domain.RoleLogist.String(), cargoitems_service.AssignGateInput{GateID: 4})
	if !errors.Is(err, core_errors.ErrConflict) {
		t.Fatalf("err = %v, want ErrConflict", err)
	}
}

func TestGetCargoItemHistoryChecksAccessAndMapsDTO(t *testing.T) {
	oldStatus := core_domain.CargoItemStatusAccepted
	now := time.Now().UTC()
	repo := &fakeCargoItemsRepository{
		getCargoItemFn: func(ctx context.Context, cargoItemID int64) (core_domain.CargoItem, error) {
			return core_domain.CargoItem{ID: cargoItemID, Status: core_domain.CargoItemStatusStored}, nil
		},
		clientOwnsCargoItemFn: func(ctx context.Context, cargoItemID int64, clientID int64) (bool, error) {
			return true, nil
		},
		listCargoStatusHistoryFn: func(ctx context.Context, cargoItemID int64) ([]core_domain.CargoStatusHistory, error) {
			if cargoItemID != 7 {
				t.Fatalf("cargoItemID = %d, want 7", cargoItemID)
			}
			return []core_domain.CargoStatusHistory{{
				ID:          1,
				CargoItemID: cargoItemID,
				OldStatus:   &oldStatus,
				NewStatus:   core_domain.CargoItemStatusStored,
				ChangedBy:   5,
				ChangedAt:   now,
			}}, nil
		},
	}
	service := cargoitems_service.NewCargoItemsService(repo)

	history, err := service.GetCargoItemHistory(context.Background(), 7, 42, core_domain.RoleClient.String())
	if err != nil {
		t.Fatalf("GetCargoItemHistory returned error: %v", err)
	}
	if len(history) != 1 || history[0].OldStatus == nil || *history[0].OldStatus != "accepted" || history[0].NewStatus != "stored" {
		t.Fatalf("unexpected history: %+v", history)
	}
}

func TestScanCargoItemByQRCodeSuccessScopesClient(t *testing.T) {
	now := time.Now().UTC()
	repo := &fakeCargoItemsRepository{
		listCargoItemsFn: func(ctx context.Context, filter core_domain.CargoItemFilter) ([]core_domain.CargoItem, error) {
			if filter.QRCode != "QR-SCAN-001" {
				t.Fatalf("QRCode = %q, want QR-SCAN-001", filter.QRCode)
			}
			if filter.ClientID == nil || *filter.ClientID != 42 {
				t.Fatalf("ClientID = %v, want 42", filter.ClientID)
			}
			if filter.Page != 1 || filter.Limit != 1 {
				t.Fatalf("page/limit = %d/%d, want 1/1", filter.Page, filter.Limit)
			}
			return []core_domain.CargoItem{{ID: 7, OrderID: 10, QRCode: "QR-SCAN-001", Status: core_domain.CargoItemStatusStored, CreatedAt: now, UpdatedAt: now}}, nil
		},
	}
	service := cargoitems_service.NewCargoItemsService(repo)

	dto, err := service.ScanCargoItem(context.Background(), 42, core_domain.RoleClient.String(), " QR-SCAN-001 ")
	if err != nil {
		t.Fatalf("ScanCargoItem returned error: %v", err)
	}
	if dto.ID != 7 || dto.QRCode != "QR-SCAN-001" {
		t.Fatalf("unexpected dto: %+v", dto)
	}
}

func TestScanCargoItemRejectsEmptyQRCode(t *testing.T) {
	service := cargoitems_service.NewCargoItemsService(&fakeCargoItemsRepository{})

	_, err := service.ScanCargoItem(context.Background(), 42, core_domain.RoleWorker.String(), "   ")
	if !errors.Is(err, core_errors.ErrInvalidArgument) {
		t.Fatalf("err = %v, want ErrInvalidArgument", err)
	}
}

func TestScanCargoItemReturnsNotFound(t *testing.T) {
	repo := &fakeCargoItemsRepository{
		listCargoItemsFn: func(ctx context.Context, filter core_domain.CargoItemFilter) ([]core_domain.CargoItem, error) {
			return nil, nil
		},
	}
	service := cargoitems_service.NewCargoItemsService(repo)

	_, err := service.ScanCargoItem(context.Background(), 42, core_domain.RoleWorker.String(), "QR-MISSING")
	if !errors.Is(err, core_errors.ErrNotFound) {
		t.Fatalf("err = %v, want ErrNotFound", err)
	}
}

func TestGetCargoItemLabelSuccess(t *testing.T) {
	now := time.Now().UTC()
	zoneID := int64(9)
	gateID := int64(4)
	repo := &fakeCargoItemsRepository{
		getCargoItemFn: func(ctx context.Context, cargoItemID int64) (core_domain.CargoItem, error) {
			if cargoItemID != 7 {
				t.Fatalf("cargoItemID = %d, want 7", cargoItemID)
			}
			return core_domain.CargoItem{
				ID:                cargoItemID,
				OrderID:           10,
				OrderCargoPlaceID: 20,
				CargoPlaceTypeID:  3,
				QRCode:            "QR-LABEL-001",
				Status:            core_domain.CargoItemStatusStored,
				StorageZoneID:     &zoneID,
				GateID:            &gateID,
				ReceivedAt:        &now,
				CreatedAt:         now,
				UpdatedAt:         now,
			}, nil
		},
	}
	service := cargoitems_service.NewCargoItemsService(repo)

	label, err := service.GetCargoItemLabel(context.Background(), 7, 5, core_domain.RoleWorker.String())
	if err != nil {
		t.Fatalf("GetCargoItemLabel returned error: %v", err)
	}
	if label.CargoItemID != 7 || label.QRCodeValue != "QR-LABEL-001" || !strings.Contains(label.LabelText, "ORDER-10") {
		t.Fatalf("unexpected label: %+v", label)
	}
}
