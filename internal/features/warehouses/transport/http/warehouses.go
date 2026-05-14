package warehouses_transport_http

import (
	"net/http"

	core_http_request "github.com/pkpal-uhobp/fulfillment-app/internal/core/transport/http/request"
	core_http_response "github.com/pkpal-uhobp/fulfillment-app/internal/core/transport/http/response"
	core_http_utils "github.com/pkpal-uhobp/fulfillment-app/internal/core/transport/http/utils"
	warehouses_service "github.com/pkpal-uhobp/fulfillment-app/internal/features/warehouses/service"
)

func (h *WarehousesHTTPHandler) ListWarehouses(w http.ResponseWriter, r *http.Request) {
	response := core_http_response.NewHTTPResponseHandler(h.log, w)

	warehouses, err := h.warehousesService.ListWarehouses(
		r.Context(),
		warehouses_service.WarehouseFilter{
			WarehouseType: r.URL.Query().Get("warehouse_type"),
			Marketplace:   r.URL.Query().Get("marketplace"),
			City:          r.URL.Query().Get("city"),
		},
	)
	if err != nil {
		response.ErrorResponse(err, "list warehouses")
		return
	}

	response.JSONResponse(WarehousesResponse{Warehouses: warehouses}, http.StatusOK)
}

func (h *WarehousesHTTPHandler) GetWarehouse(w http.ResponseWriter, r *http.Request) {
	response := core_http_response.NewHTTPResponseHandler(h.log, w)

	warehouseID, err := core_http_utils.PathInt64(r, "id")
	if err != nil {
		response.ErrorResponse(err, "invalid warehouse id")
		return
	}

	warehouse, err := h.warehousesService.GetWarehouse(r.Context(), warehouseID)
	if err != nil {
		response.ErrorResponse(err, "get warehouse")
		return
	}

	response.JSONResponse(WarehouseResponse{Warehouse: warehouse}, http.StatusOK)
}

func (h *WarehousesHTTPHandler) CreateWarehouse(w http.ResponseWriter, r *http.Request) {
	response := core_http_response.NewHTTPResponseHandler(h.log, w)

	var request CreateWarehouseRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		response.ErrorResponse(err, "invalid create warehouse request")
		return
	}

	warehouse, err := h.warehousesService.CreateWarehouse(
		r.Context(),
		warehouses_service.CreateWarehouseInput{
			Name:          request.Name,
			WarehouseType: request.WarehouseType,
			Marketplace:   request.Marketplace,
			City:          request.City,
			Address:       request.Address,
		},
	)
	if err != nil {
		response.ErrorResponse(err, "create warehouse")
		return
	}

	response.JSONResponse(WarehouseResponse{Warehouse: warehouse}, http.StatusCreated)
}

func (h *WarehousesHTTPHandler) PatchWarehouse(w http.ResponseWriter, r *http.Request) {
	response := core_http_response.NewHTTPResponseHandler(h.log, w)

	warehouseID, err := core_http_utils.PathInt64(r, "id")
	if err != nil {
		response.ErrorResponse(err, "invalid warehouse id")
		return
	}

	var request PatchWarehouseRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		response.ErrorResponse(err, "invalid patch warehouse request")
		return
	}

	warehouse, err := h.warehousesService.PatchWarehouse(
		r.Context(),
		warehouseID,
		warehouses_service.PatchWarehouseInput{
			Name:          request.Name,
			WarehouseType: request.WarehouseType,
			Marketplace:   request.Marketplace,
			City:          request.City,
			Address:       request.Address,
			IsActive:      request.IsActive,
		},
	)
	if err != nil {
		response.ErrorResponse(err, "patch warehouse")
		return
	}

	response.JSONResponse(WarehouseResponse{Warehouse: warehouse}, http.StatusOK)
}

func (h *WarehousesHTTPHandler) DeactivateWarehouse(w http.ResponseWriter, r *http.Request) {
	response := core_http_response.NewHTTPResponseHandler(h.log, w)

	warehouseID, err := core_http_utils.PathInt64(r, "id")
	if err != nil {
		response.ErrorResponse(err, "invalid warehouse id")
		return
	}

	if err := h.warehousesService.DeactivateWarehouse(r.Context(), warehouseID); err != nil {
		response.ErrorResponse(err, "deactivate warehouse")
		return
	}

	response.NoContentResponse()
}
