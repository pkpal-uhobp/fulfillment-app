package warehouses_transport_http

import (
	"net/http"

	core_http_request "github.com/pkpal-uhobp/fulfillment-app/internal/core/transport/http/request"
	core_http_response "github.com/pkpal-uhobp/fulfillment-app/internal/core/transport/http/response"
	warehouses_service "github.com/pkpal-uhobp/fulfillment-app/internal/features/warehouses/service"
)

func (h *WarehousesHTTPHandler) ListGates(w http.ResponseWriter, r *http.Request) {
	response := core_http_response.NewHTTPResponseHandler(h.log, w)

	warehouseID, err := queryInt64(r, "warehouse_id")
	if err != nil {
		response.ErrorResponse(err, "invalid warehouse id")
		return
	}

	gates, err := h.warehousesService.ListGates(r.Context(), warehouseID)
	if err != nil {
		response.ErrorResponse(err, "list gates")
		return
	}

	response.JSONResponse(GatesResponse{Gates: gates}, http.StatusOK)
}

func (h *WarehousesHTTPHandler) CreateGate(w http.ResponseWriter, r *http.Request) {
	response := core_http_response.NewHTTPResponseHandler(h.log, w)

	var request CreateGateRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		response.ErrorResponse(err, "invalid create gate request")
		return
	}

	gate, err := h.warehousesService.CreateGate(
		r.Context(),
		warehouses_service.CreateGateInput{
			WarehouseID: request.WarehouseID,
			Name:        request.Name,
		},
	)
	if err != nil {
		response.ErrorResponse(err, "create gate")
		return
	}

	response.JSONResponse(GateResponse{Gate: gate}, http.StatusCreated)
}

func (h *WarehousesHTTPHandler) PatchGate(w http.ResponseWriter, r *http.Request) {
	response := core_http_response.NewHTTPResponseHandler(h.log, w)

	gateID, err := pathInt64(r, "id")
	if err != nil {
		response.ErrorResponse(err, "invalid gate id")
		return
	}

	var request PatchGateRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		response.ErrorResponse(err, "invalid patch gate request")
		return
	}

	gate, err := h.warehousesService.PatchGate(
		r.Context(),
		gateID,
		warehouses_service.PatchGateInput{
			Name:     request.Name,
			IsActive: request.IsActive,
		},
	)
	if err != nil {
		response.ErrorResponse(err, "patch gate")
		return
	}

	response.JSONResponse(GateResponse{Gate: gate}, http.StatusOK)
}
