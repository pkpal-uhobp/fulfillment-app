//go:build integration

package orders_repository_postgres_tests

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
	core_errors "github.com/pkpal-uhobp/fulfillment-app/internal/core/errors"
	core_postgres_pool "github.com/pkpal-uhobp/fulfillment-app/internal/core/repository/pool"
	core_postgres_tx "github.com/pkpal-uhobp/fulfillment-app/internal/core/repository/tx"
	orders_repository_postgres "github.com/pkpal-uhobp/fulfillment-app/internal/features/orders/repository/postgres"
)

type seedRefs struct {
	clientID               int64
	otherClientID          int64
	logistID               int64
	receivingWarehouseID   int64
	destinationWarehouseID int64
	productTypeID          int64
	cargoPlaceTypeID       int64
}

func newIntegrationRepository(t *testing.T) (*orders_repository_postgres.OrdersRepository, *core_postgres_pool.ConnectionPool) {
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

	_, err = pool.Exec(
		context.Background(),
		`TRUNCATE TABLE
			users,
			warehouses,
			product_types,
			cargo_place_types,
			orders,
			order_cargo_places,
			pickup_requests,
			order_status_history
		RESTART IDENTITY CASCADE`,
	)
	if err != nil {
		t.Fatalf("truncate orders test tables: %v", err)
	}

	tx := core_postgres_tx.NewTx(pool)
	repo := orders_repository_postgres.NewOrdersRepository(tx)
	return repo, pool
}

func envOrDefault(name string, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}

func seedBaseData(t *testing.T, pool *core_postgres_pool.ConnectionPool) seedRefs {
	t.Helper()
	ctx := context.Background()
	suffix := time.Now().UnixNano()

	clientID := insertUser(t, pool, fmt.Sprintf("client-%d@test.com", suffix), core_domain.RoleClient.String())
	otherClientID := insertUser(t, pool, fmt.Sprintf("other-client-%d@test.com", suffix), core_domain.RoleClient.String())
	logistID := insertUser(t, pool, fmt.Sprintf("logist-%d@test.com", suffix), core_domain.RoleLogist.String())

	var receivingWarehouseID int64
	if err := pool.QueryRow(
		ctx,
		`INSERT INTO warehouses (name, warehouse_type, marketplace, city, address)
		 VALUES ($1, 'receiving', NULL, 'Perm', 'Perm, Lenina 1')
		 RETURNING id`,
		fmt.Sprintf("Receiving Test Warehouse %d", suffix),
	).Scan(&receivingWarehouseID); err != nil {
		t.Fatalf("insert receiving warehouse: %v", err)
	}

	var destinationWarehouseID int64
	if err := pool.QueryRow(
		ctx,
		`INSERT INTO warehouses (name, warehouse_type, marketplace, city, address)
		 VALUES ($1, 'destination', 'Ozon', 'Moscow', 'Moscow, Tverskaya 1')
		 RETURNING id`,
		fmt.Sprintf("Destination Test Warehouse %d", suffix),
	).Scan(&destinationWarehouseID); err != nil {
		t.Fatalf("insert destination warehouse: %v", err)
	}

	var productTypeID int64
	if err := pool.QueryRow(
		ctx,
		`INSERT INTO product_types (name, description) VALUES ($1, 'test product type') RETURNING id`,
		fmt.Sprintf("Product Type %d", suffix),
	).Scan(&productTypeID); err != nil {
		t.Fatalf("insert product type: %v", err)
	}

	var cargoPlaceTypeID int64
	if err := pool.QueryRow(
		ctx,
		`INSERT INTO cargo_place_types (name, description) VALUES ($1, 'test cargo place type') RETURNING id`,
		fmt.Sprintf("Box %d", suffix),
	).Scan(&cargoPlaceTypeID); err != nil {
		t.Fatalf("insert cargo place type: %v", err)
	}

	return seedRefs{
		clientID:               clientID,
		otherClientID:          otherClientID,
		logistID:               logistID,
		receivingWarehouseID:   receivingWarehouseID,
		destinationWarehouseID: destinationWarehouseID,
		productTypeID:          productTypeID,
		cargoPlaceTypeID:       cargoPlaceTypeID,
	}
}

func insertUser(t *testing.T, pool *core_postgres_pool.ConnectionPool, email string, role string) int64 {
	t.Helper()

	var id int64
	if err := pool.QueryRow(
		context.Background(),
		`INSERT INTO users (email, password_hash, full_name, phone, role)
		 VALUES ($1, 'hashed-password', 'Test User', NULL, $2)
		 RETURNING id`,
		email,
		role,
	).Scan(&id); err != nil {
		t.Fatalf("insert user %s: %v", role, err)
	}
	return id
}

func createSelfDeliveryOrder(
	t *testing.T,
	repo *orders_repository_postgres.OrdersRepository,
	refs seedRefs,
	clientID int64,
	status core_domain.OrderStatus,
) core_domain.OrderDetails {
	t.Helper()

	selfDeliveryDate := "2026-06-01"
	details, err := repo.CreateOrder(
		context.Background(),
		core_domain.Order{
			ClientID:               clientID,
			ReceivingWarehouseID:   refs.receivingWarehouseID,
			DestinationWarehouseID: refs.destinationWarehouseID,
			ProductTypeID:          refs.productTypeID,
			HandoverType:           core_domain.HandoverTypeSelfDelivery,
			SelfDeliveryDate:       &selfDeliveryDate,
			Status:                 status,
		},
		[]core_domain.OrderCargoPlace{
			{CargoPlaceTypeID: refs.cargoPlaceTypeID, Quantity: 2},
		},
		nil,
		clientID,
	)
	if err != nil {
		t.Fatalf("CreateOrder self_delivery returned error: %v", err)
	}
	return details
}

func TestCreateGetAndHistorySelfDeliveryOrder(t *testing.T) {
	repo, pool := newIntegrationRepository(t)
	refs := seedBaseData(t, pool)

	comment := "Самостоятельная сдача на склад"
	selfDeliveryDate := "2026-06-01"
	weight := 12.5

	created, err := repo.CreateOrder(
		context.Background(),
		core_domain.Order{
			ClientID:               refs.clientID,
			ReceivingWarehouseID:   refs.receivingWarehouseID,
			DestinationWarehouseID: refs.destinationWarehouseID,
			ProductTypeID:          refs.productTypeID,
			HandoverType:           core_domain.HandoverTypeSelfDelivery,
			SelfDeliveryDate:       &selfDeliveryDate,
			Status:                 core_domain.OrderStatusCreated,
			Comment:                &comment,
		},
		[]core_domain.OrderCargoPlace{
			{CargoPlaceTypeID: refs.cargoPlaceTypeID, Quantity: 2, WeightPerPlaceKG: &weight},
		},
		nil,
		refs.clientID,
	)
	if err != nil {
		t.Fatalf("CreateOrder returned error: %v", err)
	}
	if created.Order.ID == 0 {
		t.Fatalf("created order ID must not be zero")
	}
	if len(created.CargoPlaces) != 1 || created.CargoPlaces[0].Quantity != 2 {
		t.Fatalf("unexpected cargo places: %+v", created.CargoPlaces)
	}
	if created.Pickup != nil {
		t.Fatalf("pickup = %+v, want nil", created.Pickup)
	}

	loaded, err := repo.GetOrder(context.Background(), created.Order.ID)
	if err != nil {
		t.Fatalf("GetOrder returned error: %v", err)
	}
	if loaded.Order.ID != created.Order.ID || loaded.Order.ClientID != refs.clientID {
		t.Fatalf("unexpected loaded order: %+v", loaded.Order)
	}

	history, err := repo.ListOrderStatusHistory(context.Background(), created.Order.ID)
	if err != nil {
		t.Fatalf("ListOrderStatusHistory returned error: %v", err)
	}
	if len(history) != 1 || history[0].NewStatus != core_domain.OrderStatusCreated {
		t.Fatalf("unexpected history: %+v", history)
	}
}

func TestCreatePickupOrderCreatesPickupRequest(t *testing.T) {
	repo, pool := newIntegrationRepository(t)
	refs := seedBaseData(t, pool)

	pickupDate := "2026-06-02"
	timeFrom := "09:00"
	timeTo := "11:00"
	contactName := "Test Client"
	contactPhone := "+79990000000"

	created, err := repo.CreateOrder(
		context.Background(),
		core_domain.Order{
			ClientID:               refs.clientID,
			ReceivingWarehouseID:   refs.receivingWarehouseID,
			DestinationWarehouseID: refs.destinationWarehouseID,
			ProductTypeID:          refs.productTypeID,
			HandoverType:           core_domain.HandoverTypePickup,
			Status:                 core_domain.OrderStatusCreated,
		},
		[]core_domain.OrderCargoPlace{
			{CargoPlaceTypeID: refs.cargoPlaceTypeID, Quantity: 1},
		},
		&core_domain.PickupRequest{
			PickupAddress:  "Perm, Lenina 10",
			PickupDate:     pickupDate,
			PickupTimeFrom: &timeFrom,
			PickupTimeTo:   &timeTo,
			ContactName:    &contactName,
			ContactPhone:   &contactPhone,
		},
		refs.clientID,
	)
	if err != nil {
		t.Fatalf("CreateOrder pickup returned error: %v", err)
	}
	if created.Pickup == nil {
		t.Fatalf("created.Pickup = nil, want pickup request")
	}
	if created.Pickup.PickupAddress != "Perm, Lenina 10" || created.Pickup.PickupDate != pickupDate {
		t.Fatalf("unexpected pickup: %+v", created.Pickup)
	}
}

func TestListOrdersFiltersByClientAndStatus(t *testing.T) {
	repo, pool := newIntegrationRepository(t)
	refs := seedBaseData(t, pool)

	first := createSelfDeliveryOrder(t, repo, refs, refs.clientID, core_domain.OrderStatusCreated)
	_ = createSelfDeliveryOrder(t, repo, refs, refs.otherClientID, core_domain.OrderStatusCreated)

	orders, err := repo.ListOrders(
		context.Background(),
		core_domain.OrderFilter{
			ClientID: &refs.clientID,
			Status:   core_domain.OrderStatusCreated.String(),
			Page:     1,
			Limit:    10,
		},
	)
	if err != nil {
		t.Fatalf("ListOrders returned error: %v", err)
	}
	if len(orders) != 1 {
		t.Fatalf("len(orders) = %d, want 1: %+v", len(orders), orders)
	}
	if orders[0].Order.ID != first.Order.ID {
		t.Fatalf("orders[0].ID = %d, want %d", orders[0].Order.ID, first.Order.ID)
	}
}

func TestUpdateStatusAndCancelCreateHistory(t *testing.T) {
	repo, pool := newIntegrationRepository(t)
	refs := seedBaseData(t, pool)
	created := createSelfDeliveryOrder(t, repo, refs, refs.clientID, core_domain.OrderStatusCreated)

	comment := "Ожидает самостоятельной сдачи"
	if err := repo.UpdateOrderStatus(
		context.Background(),
		created.Order.ID,
		core_domain.OrderStatusWaitingDelivery.String(),
		refs.logistID,
		&comment,
	); err != nil {
		t.Fatalf("UpdateOrderStatus returned error: %v", err)
	}

	updated, err := repo.GetOrder(context.Background(), created.Order.ID)
	if err != nil {
		t.Fatalf("GetOrder after update returned error: %v", err)
	}
	if updated.Order.Status != core_domain.OrderStatusWaitingDelivery {
		t.Fatalf("Status = %q, want waiting_delivery", updated.Order.Status)
	}

	cancelComment := "Клиент отменил заявку"
	if err := repo.CancelOrder(context.Background(), created.Order.ID, refs.clientID, &cancelComment); err != nil {
		t.Fatalf("CancelOrder returned error: %v", err)
	}

	cancelled, err := repo.GetOrder(context.Background(), created.Order.ID)
	if err != nil {
		t.Fatalf("GetOrder after cancel returned error: %v", err)
	}
	if cancelled.Order.Status != core_domain.OrderStatusCancelled {
		t.Fatalf("Status = %q, want cancelled", cancelled.Order.Status)
	}

	history, err := repo.ListOrderStatusHistory(context.Background(), created.Order.ID)
	if err != nil {
		t.Fatalf("ListOrderStatusHistory returned error: %v", err)
	}
	if len(history) != 3 {
		t.Fatalf("len(history) = %d, want 3: %+v", len(history), history)
	}
	if history[1].OldStatus == nil || *history[1].OldStatus != core_domain.OrderStatusCreated || history[1].NewStatus != core_domain.OrderStatusWaitingDelivery {
		t.Fatalf("unexpected update history item: %+v", history[1])
	}
	if history[2].NewStatus != core_domain.OrderStatusCancelled {
		t.Fatalf("unexpected cancel history item: %+v", history[2])
	}
}

func TestGetOrderReturnsNotFound(t *testing.T) {
	repo, _ := newIntegrationRepository(t)

	_, err := repo.GetOrder(context.Background(), 999)
	if !errors.Is(err, core_errors.ErrNotFound) {
		t.Fatalf("err = %v, want ErrNotFound", err)
	}
}
