package domain

import (
	"fmt"
	"strings"
	"time"

	core_errors "github.com/pkpal-uhobp/fulfillment-app/internal/core/errors"
)

type WarehouseType string

const (
	WarehouseTypeReceiving   WarehouseType = "receiving"
	WarehouseTypeDestination WarehouseType = "destination"
	WarehouseTypeBoth        WarehouseType = "both"
)

func (t WarehouseType) String() string {
	return string(t)
}

func (t WarehouseType) IsValid() bool {
	switch t {
	case WarehouseTypeReceiving, WarehouseTypeDestination, WarehouseTypeBoth:
		return true
	default:
		return false
	}
}

type WarehouseFilter struct {
	WarehouseType string
	Marketplace   string
	City          string
}

type Warehouse struct {
	ID            int64
	Name          string
	WarehouseType WarehouseType
	Marketplace   *string
	City          string
	Address       string
	IsActive      bool
	CreatedAt     time.Time
}

type WarehousePatch struct {
	Name          *string
	WarehouseType *WarehouseType
	Marketplace   *string
	City          *string
	Address       *string
	IsActive      *bool
}

func NewWarehouse(
	name string,
	warehouseType WarehouseType,
	marketplace *string,
	city string,
	address string,
) (Warehouse, error) {
	warehouse := Warehouse{
		Name:          strings.TrimSpace(name),
		WarehouseType: warehouseType,
		Marketplace:   normalizeOptionalString(marketplace),
		City:          strings.TrimSpace(city),
		Address:       strings.TrimSpace(address),
		IsActive:      true,
	}

	if err := warehouse.Validate(); err != nil {
		return Warehouse{}, err
	}

	return warehouse, nil
}

func (w Warehouse) Validate() error {
	if strings.TrimSpace(w.Name) == "" {
		return fmt.Errorf("%w: warehouse name is required", core_errors.ErrInvalidArgument)
	}

	if !w.WarehouseType.IsValid() {
		return fmt.Errorf("%w: invalid warehouse type", core_errors.ErrInvalidArgument)
	}

	if strings.TrimSpace(w.City) == "" {
		return fmt.Errorf("%w: warehouse city is required", core_errors.ErrInvalidArgument)
	}

	if strings.TrimSpace(w.Address) == "" {
		return fmt.Errorf("%w: warehouse address is required", core_errors.ErrInvalidArgument)
	}

	return nil
}

type StorageZone struct {
	ID          int64
	WarehouseID int64
	Name        string
	Description *string
	IsActive    bool
}

type StorageZonePatch struct {
	Name        *string
	Description *string
	IsActive    *bool
}

func NewStorageZone(
	warehouseID int64,
	name string,
	description *string,
) (StorageZone, error) {
	zone := StorageZone{
		WarehouseID: warehouseID,
		Name:        strings.TrimSpace(name),
		Description: normalizeOptionalString(description),
		IsActive:    true,
	}

	if err := zone.Validate(); err != nil {
		return StorageZone{}, err
	}

	return zone, nil
}

func (z StorageZone) Validate() error {
	if z.WarehouseID <= 0 {
		return fmt.Errorf("%w: warehouse id is required", core_errors.ErrInvalidArgument)
	}

	if strings.TrimSpace(z.Name) == "" {
		return fmt.Errorf("%w: storage zone name is required", core_errors.ErrInvalidArgument)
	}

	return nil
}

type Gate struct {
	ID          int64
	WarehouseID int64
	Name        string
	IsActive    bool
}

type GatePatch struct {
	Name     *string
	IsActive *bool
}

func NewGate(
	warehouseID int64,
	name string,
) (Gate, error) {
	gate := Gate{
		WarehouseID: warehouseID,
		Name:        strings.TrimSpace(name),
		IsActive:    true,
	}

	if err := gate.Validate(); err != nil {
		return Gate{}, err
	}

	return gate, nil
}

func (g Gate) Validate() error {
	if g.WarehouseID <= 0 {
		return fmt.Errorf("%w: warehouse id is required", core_errors.ErrInvalidArgument)
	}

	if strings.TrimSpace(g.Name) == "" {
		return fmt.Errorf("%w: gate name is required", core_errors.ErrInvalidArgument)
	}

	return nil
}

type ProductType struct {
	ID          int64
	Name        string
	Description *string
	IsActive    bool
}

type CargoPlaceType struct {
	ID          int64
	Name        string
	Description *string
	IsActive    bool
}

func normalizeOptionalString(value *string) *string {
	if value == nil {
		return nil
	}

	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}

	return &trimmed
}
