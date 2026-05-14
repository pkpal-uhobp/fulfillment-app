package warehouses_transport_http

import (
	"net/http"

	core_http_request "github.com/pkpal-uhobp/fulfillment-app/internal/core/transport/http/request"
	core_http_response "github.com/pkpal-uhobp/fulfillment-app/internal/core/transport/http/response"
	warehouses_service "github.com/pkpal-uhobp/fulfillment-app/internal/features/warehouses/service"
)

func (h *WarehousesHTTPHandler) ListStorageZones(w http.ResponseWriter, r *http.Request) {
	response := core_http_response.NewHTTPResponseHandler(h.log, w)

	warehouseID, err := queryInt64(r, "warehouse_id")
	if err != nil {
		response.ErrorResponse(err, "invalid warehouse id")
		return
	}

	zones, err := h.warehousesService.ListStorageZones(r.Context(), warehouseID)
	if err != nil {
		response.ErrorResponse(err, "list storage zones")
		return
	}

	response.JSONResponse(StorageZonesResponse{StorageZones: zones}, http.StatusOK)
}

func (h *WarehousesHTTPHandler) CreateStorageZone(w http.ResponseWriter, r *http.Request) {
	response := core_http_response.NewHTTPResponseHandler(h.log, w)

	var request CreateStorageZoneRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		response.ErrorResponse(err, "invalid create storage zone request")
		return
	}

	zone, err := h.warehousesService.CreateStorageZone(
		r.Context(),
		warehouses_service.CreateStorageZoneInput{
			WarehouseID: request.WarehouseID,
			Name:        request.Name,
			Description: request.Description,
		},
	)
	if err != nil {
		response.ErrorResponse(err, "create storage zone")
		return
	}

	response.JSONResponse(StorageZoneResponse{StorageZone: zone}, http.StatusCreated)
}

func (h *WarehousesHTTPHandler) PatchStorageZone(w http.ResponseWriter, r *http.Request) {
	response := core_http_response.NewHTTPResponseHandler(h.log, w)

	zoneID, err := pathInt64(r, "id")
	if err != nil {
		response.ErrorResponse(err, "invalid storage zone id")
		return
	}

	var request PatchStorageZoneRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		response.ErrorResponse(err, "invalid patch storage zone request")
		return
	}

	zone, err := h.warehousesService.PatchStorageZone(
		r.Context(),
		zoneID,
		warehouses_service.PatchStorageZoneInput{
			Name:        request.Name,
			Description: request.Description,
			IsActive:    request.IsActive,
		},
	)
	if err != nil {
		response.ErrorResponse(err, "patch storage zone")
		return
	}

	response.JSONResponse(StorageZoneResponse{StorageZone: zone}, http.StatusOK)
}
