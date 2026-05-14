//go:build integration

package cargoitems_repository_postgres_tests

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
	core_errors "github.com/pkpal-uhobp/fulfillment-app/internal/core/errors"
	core_postgres_pool "github.com/pkpal-uhobp/fulfillment-app/internal/core/repository/pool"
	core_postgres_tx "github.com/pkpal-uhobp/fulfillment-app/internal/core/repository/tx"
	cargoitems_repository_postgres "github.com/pkpal-uhobp/fulfillment-app/internal/features/cargoitems/repository/postgres"
)

type cargoItemsRepositoryFixture struct {
	ClientID               int64
	WorkerID               int64
	ReceivingWarehouseID   int64
	DestinationWarehouseID int64
	StorageZoneID          int64
	GateID                 int64
	ProductTypeID          int64
	CargoPlaceTypeID       int64
	OrderID                int64
	OrderCargoPlaceID      int64
}

func newIntegrationRepository(t *testing.T) (*cargoitems_repository_postgres.CargoItemsRepository, *core_postgres_pool.ConnectionPool) {
	t.Helper()

	requiredEnv := []string{
		"POSTGRES_HOST",
		"POSTGRES_USER",
		"POSTGRES_PASSWORD",
		"POSTGRES_DB",
	}
	for _, envName := range requiredEnv {
		if os.Getenv(envName) == "" {
			t.Skipf("skip integration test: %s is not set", envName)
		}
	}

	config := core_postgres_pool.Config{
		Host:              os.Getenv("POSTGRES_HOST"),
		Port:              envOrDefault("POSTGRES_PORT", "5432"),
		User:              os.Getenv("POSTGRES_USER"),
		Password:          os.Getenv("POSTGRES_PASSWORD"),
		Database:          os.Getenv("POSTGRES_DB"),
		SSLMode:           envOrDefault("POSTGRES_SSL_MODE", "disable"),
		MaxConns:          4,
		MinConns:          1,
		MaxConnLifetime:   time.Hour,
		MaxConnIdleTime:   30 * time.Minute,
		HealthCheckPeriod: time.Minute,
		ConnectTimeout:    5 * time.Second,
		QueryTimeout:      5 * time.Second,
	}

	pool, err := core_postgres_pool.NewConnectionPool(context.Background(), config)
	if err != nil {
		t.Fatalf("connect postgres: %v", err)
	}
	t.Cleanup(pool.Close)

	_, err = pool.Exec(context.Background(), `
		TRUNCATE TABLE
			cargo_status_history,
			cargo_items,
			order_cargo_places,
			orders,
			storage_zones,
			gates,
			worker_profiles,
			warehouses,
			product_types,
			cargo_place_types,
			issued_tokens,
			users
		RESTART IDENTITY CASCADE;
	`)
	if err != nil {
		t.Fatalf("truncate cargoitems integration tables: %v", err)
	}

	tx := core_postgres_tx.NewTx(pool)
	repo := cargoitems_repository_postgres.NewCargoItemsRepository(tx)

	return repo, pool
}

func envOrDefault(name string, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}

func seedCargoItemsFixture(t *testing.T, pool *core_postgres_pool.ConnectionPool) cargoItemsRepositoryFixture {
	t.Helper()

	ctx := context.Background()
	suffix := uuid.NewString()

	var fixture cargoItemsRepositoryFixture

	err := pool.QueryRow(ctx, `
		INSERT INTO users (email, password_hash, full_name, phone, role)
		VALUES ($1, 'hash', 'Client User', NULL, 'client')
		RETURNING id;
	`, fmt.Sprintf("client-%s@example.com", suffix)).Scan(&fixture.ClientID)
	if err != nil {
		t.Fatalf("insert client user: %v", err)
	}

	err = pool.QueryRow(ctx, `
		INSERT INTO users (email, password_hash, full_name, phone, role)
		VALUES ($1, 'hash', 'Worker User', NULL, 'worker')
		RETURNING id;
	`, fmt.Sprintf("worker-%s@example.com", suffix)).Scan(&fixture.WorkerID)
	if err != nil {
		t.Fatalf("insert worker user: %v", err)
	}

	err = pool.QueryRow(ctx, `
		INSERT INTO warehouses (name, warehouse_type, marketplace, city, address)
		VALUES ($1, 'both', NULL, 'Moscow', 'Test address 1')
		RETURNING id;
	`, "Receiving Warehouse "+suffix[:8]).Scan(&fixture.ReceivingWarehouseID)
	if err != nil {
		t.Fatalf("insert receiving warehouse: %v", err)
	}

	err = pool.QueryRow(ctx, `
		INSERT INTO warehouses (name, warehouse_type, marketplace, city, address)
		VALUES ($1, 'destination', 'Ozon', 'Moscow', 'Test address 2')
		RETURNING id;
	`, "Destination Warehouse "+suffix[:8]).Scan(&fixture.DestinationWarehouseID)
	if err != nil {
		t.Fatalf("insert destination warehouse: %v", err)
	}

	err = pool.QueryRow(ctx, `
		INSERT INTO storage_zones (warehouse_id, name, description)
		VALUES ($1, $2, 'Test zone')
		RETURNING id;
	`, fixture.ReceivingWarehouseID, "A-"+suffix[:8]).Scan(&fixture.StorageZoneID)
	if err != nil {
		t.Fatalf("insert storage zone: %v", err)
	}

	err = pool.QueryRow(ctx, `
		INSERT INTO gates (warehouse_id, name)
		VALUES ($1, $2)
		RETURNING id;
	`, fixture.ReceivingWarehouseID, "G-"+suffix[:8]).Scan(&fixture.GateID)
	if err != nil {
		t.Fatalf("insert gate: %v", err)
	}

	err = pool.QueryRow(ctx, `
		INSERT INTO product_types (name, description)
		VALUES ($1, 'Test product type')
		RETURNING id;
	`, "Product-"+suffix[:8]).Scan(&fixture.ProductTypeID)
	if err != nil {
		t.Fatalf("insert product type: %v", err)
	}

	err = pool.QueryRow(ctx, `
		INSERT INTO cargo_place_types (name, description)
		VALUES ($1, 'Test cargo place type')
		RETURNING id;
	`, "Box-"+suffix[:8]).Scan(&fixture.CargoPlaceTypeID)
	if err != nil {
		t.Fatalf("insert cargo place type: %v", err)
	}

	err = pool.QueryRow(ctx, `
		INSERT INTO orders (
			client_id,
			receiving_warehouse_id,
			destination_warehouse_id,
			product_type_id,
			handover_type,
			self_delivery_date,
			status,
			comment
		)
		VALUES ($1, $2, $3, $4, 'self_delivery', CURRENT_DATE, 'received', 'integration order')
		RETURNING id;
	`, fixture.ClientID, fixture.ReceivingWarehouseID, fixture.DestinationWarehouseID, fixture.ProductTypeID).Scan(&fixture.OrderID)
	if err != nil {
		t.Fatalf("insert order: %v", err)
	}

	err = pool.QueryRow(ctx, `
		INSERT INTO order_cargo_places (
			order_id,
			cargo_place_type_id,
			quantity,
			weight_per_place_kg,
			comment
		)
		VALUES ($1, $2, 2, 10.50, 'declared boxes')
		RETURNING id;
	`, fixture.OrderID, fixture.CargoPlaceTypeID).Scan(&fixture.OrderCargoPlaceID)
	if err != nil {
		t.Fatalf("insert order cargo place: %v", err)
	}

	return fixture
}

func createRepositoryCargoItem(
	t *testing.T,
	repo *cargoitems_repository_postgres.CargoItemsRepository,
	fixture cargoItemsRepositoryFixture,
	qrCode string,
) core_domain.CargoItem {
	t.Helper()

	comment := "accepted by integration test"
	item, err := repo.CreateCargoItem(context.Background(), core_domain.CargoItem{
		OrderID:           fixture.OrderID,
		OrderCargoPlaceID: fixture.OrderCargoPlaceID,
		CargoPlaceTypeID:  fixture.CargoPlaceTypeID,
		QRCode:            qrCode,
		Status:            core_domain.CargoItemStatusAccepted,
		Comment:           &comment,
	}, fixture.WorkerID)
	if err != nil {
		t.Fatalf("CreateCargoItem returned error: %v", err)
	}
	return item
}

func TestCreateCargoItemPersistsItemAndHistory(t *testing.T) {
	repo, pool := newIntegrationRepository(t)
	fixture := seedCargoItemsFixture(t, pool)
	ctx := context.Background()

	created := createRepositoryCargoItem(t, repo, fixture, "QR-"+uuid.NewString())

	if created.ID == 0 {
		t.Fatalf("created cargo item ID must not be zero")
	}
	if created.Status != core_domain.CargoItemStatusAccepted {
		t.Fatalf("created.Status = %q, want accepted", created.Status)
	}
	if created.ReceivedBy == nil || *created.ReceivedBy != fixture.WorkerID {
		t.Fatalf("created.ReceivedBy = %v, want %d", created.ReceivedBy, fixture.WorkerID)
	}
	if created.ReceivedAt == nil {
		t.Fatalf("created.ReceivedAt must not be nil")
	}

	found, err := repo.GetCargoItem(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetCargoItem returned error: %v", err)
	}
	if found.QRCode != created.QRCode {
		t.Fatalf("found.QRCode = %q, want %q", found.QRCode, created.QRCode)
	}

	history, err := repo.ListCargoStatusHistory(ctx, created.ID)
	if err != nil {
		t.Fatalf("ListCargoStatusHistory returned error: %v", err)
	}
	if len(history) != 1 {
		t.Fatalf("len(history) = %d, want 1", len(history))
	}
	if history[0].OldStatus != nil {
		t.Fatalf("history[0].OldStatus = %v, want nil", history[0].OldStatus)
	}
	if history[0].NewStatus != core_domain.CargoItemStatusAccepted {
		t.Fatalf("history[0].NewStatus = %q, want accepted", history[0].NewStatus)
	}
}

func TestCreateCargoItemRejectsDuplicateQRCode(t *testing.T) {
	repo, pool := newIntegrationRepository(t)
	fixture := seedCargoItemsFixture(t, pool)
	ctx := context.Background()
	qrCode := "QR-" + uuid.NewString()

	createRepositoryCargoItem(t, repo, fixture, qrCode)

	_, err := repo.CreateCargoItem(ctx, core_domain.CargoItem{
		OrderID:           fixture.OrderID,
		OrderCargoPlaceID: fixture.OrderCargoPlaceID,
		CargoPlaceTypeID:  fixture.CargoPlaceTypeID,
		QRCode:            qrCode,
		Status:            core_domain.CargoItemStatusAccepted,
	}, fixture.WorkerID)
	if !errors.Is(err, core_errors.ErrConflict) {
		t.Fatalf("err = %v, want ErrConflict", err)
	}
}

func TestGetOrderCargoPlaceForOrderAndCount(t *testing.T) {
	repo, pool := newIntegrationRepository(t)
	fixture := seedCargoItemsFixture(t, pool)
	ctx := context.Background()

	cargoPlace, orderStatus, err := repo.GetOrderCargoPlaceForOrder(ctx, fixture.OrderID, fixture.OrderCargoPlaceID)
	if err != nil {
		t.Fatalf("GetOrderCargoPlaceForOrder returned error: %v", err)
	}
	if cargoPlace.ID != fixture.OrderCargoPlaceID {
		t.Fatalf("cargoPlace.ID = %d, want %d", cargoPlace.ID, fixture.OrderCargoPlaceID)
	}
	if orderStatus != core_domain.OrderStatusReceived {
		t.Fatalf("orderStatus = %q, want received", orderStatus)
	}

	count, err := repo.CountCargoItemsByOrderCargoPlace(ctx, fixture.OrderCargoPlaceID)
	if err != nil {
		t.Fatalf("CountCargoItemsByOrderCargoPlace returned error: %v", err)
	}
	if count != 0 {
		t.Fatalf("count before create = %d, want 0", count)
	}

	createRepositoryCargoItem(t, repo, fixture, "QR-"+uuid.NewString())

	count, err = repo.CountCargoItemsByOrderCargoPlace(ctx, fixture.OrderCargoPlaceID)
	if err != nil {
		t.Fatalf("CountCargoItemsByOrderCargoPlace after create returned error: %v", err)
	}
	if count != 1 {
		t.Fatalf("count after create = %d, want 1", count)
	}
}

func TestListCargoItemsAndClientOwnership(t *testing.T) {
	repo, pool := newIntegrationRepository(t)
	fixture := seedCargoItemsFixture(t, pool)
	ctx := context.Background()
	qrCode := "QR-" + uuid.NewString()
	created := createRepositoryCargoItem(t, repo, fixture, qrCode)

	items, err := repo.ListCargoItems(ctx, core_domain.CargoItemFilter{
		OrderID:  &fixture.OrderID,
		QRCode:   qrCode,
		ClientID: &fixture.ClientID,
		Page:     1,
		Limit:    10,
	})
	if err != nil {
		t.Fatalf("ListCargoItems returned error: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("len(items) = %d, want 1", len(items))
	}
	if items[0].ID != created.ID {
		t.Fatalf("items[0].ID = %d, want %d", items[0].ID, created.ID)
	}

	owns, err := repo.ClientOwnsCargoItem(ctx, created.ID, fixture.ClientID)
	if err != nil {
		t.Fatalf("ClientOwnsCargoItem returned error: %v", err)
	}
	if !owns {
		t.Fatalf("ClientOwnsCargoItem = false, want true")
	}

	owns, err = repo.ClientOwnsCargoItem(ctx, created.ID, fixture.ClientID+999)
	if err != nil {
		t.Fatalf("ClientOwnsCargoItem for another client returned error: %v", err)
	}
	if owns {
		t.Fatalf("ClientOwnsCargoItem for another client = true, want false")
	}
}

func TestAssignZoneAssignGateAndUpdateStatus(t *testing.T) {
	repo, pool := newIntegrationRepository(t)
	fixture := seedCargoItemsFixture(t, pool)
	ctx := context.Background()
	created := createRepositoryCargoItem(t, repo, fixture, "QR-"+uuid.NewString())

	zoneOk, err := repo.StorageZoneBelongsToCargoOrder(ctx, created.ID, fixture.StorageZoneID)
	if err != nil {
		t.Fatalf("StorageZoneBelongsToCargoOrder returned error: %v", err)
	}
	if !zoneOk {
		t.Fatalf("StorageZoneBelongsToCargoOrder = false, want true")
	}

	zoneComment := "zone assigned by integration test"
	withZone, err := repo.AssignStorageZone(ctx, created.ID, fixture.StorageZoneID, fixture.WorkerID, &zoneComment)
	if err != nil {
		t.Fatalf("AssignStorageZone returned error: %v", err)
	}
	if withZone.StorageZoneID == nil || *withZone.StorageZoneID != fixture.StorageZoneID {
		t.Fatalf("StorageZoneID = %v, want %d", withZone.StorageZoneID, fixture.StorageZoneID)
	}

	storedComment := "stored by integration test"
	stored, err := repo.UpdateCargoItemStatus(ctx, created.ID, core_domain.CargoItemStatusStored.String(), fixture.WorkerID, &storedComment)
	if err != nil {
		t.Fatalf("UpdateCargoItemStatus stored returned error: %v", err)
	}
	if stored.Status != core_domain.CargoItemStatusStored {
		t.Fatalf("stored.Status = %q, want stored", stored.Status)
	}

	gateOk, err := repo.GateBelongsToCargoOrder(ctx, created.ID, fixture.GateID)
	if err != nil {
		t.Fatalf("GateBelongsToCargoOrder returned error: %v", err)
	}
	if !gateOk {
		t.Fatalf("GateBelongsToCargoOrder = false, want true")
	}

	gateComment := "gate assigned by integration test"
	withGate, err := repo.AssignGate(ctx, created.ID, fixture.GateID, fixture.WorkerID, &gateComment)
	if err != nil {
		t.Fatalf("AssignGate returned error: %v", err)
	}
	if withGate.GateID == nil || *withGate.GateID != fixture.GateID {
		t.Fatalf("GateID = %v, want %d", withGate.GateID, fixture.GateID)
	}

	readyComment := "ready to ship by integration test"
	readyToShip, err := repo.UpdateCargoItemStatus(ctx, created.ID, core_domain.CargoItemStatusReadyToShip.String(), fixture.WorkerID, &readyComment)
	if err != nil {
		t.Fatalf("UpdateCargoItemStatus ready_to_ship returned error: %v", err)
	}
	if readyToShip.Status != core_domain.CargoItemStatusReadyToShip {
		t.Fatalf("readyToShip.Status = %q, want ready_to_ship", readyToShip.Status)
	}

	history, err := repo.ListCargoStatusHistory(ctx, created.ID)
	if err != nil {
		t.Fatalf("ListCargoStatusHistory returned error: %v", err)
	}
	if len(history) < 5 {
		t.Fatalf("len(history) = %d, want at least 5", len(history))
	}
}
