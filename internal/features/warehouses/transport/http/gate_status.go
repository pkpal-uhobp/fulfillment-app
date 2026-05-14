package warehouses_transport_http

import (
	"net/http"

	core_http_response "github.com/pkpal-uhobp/fulfillment-app/internal/core/transport/http/response"
)

func (h *WarehousesHTTPHandler) ActivateGate(w http.ResponseWriter, r *http.Request) {
	response := core_http_response.NewHTTPResponseHandler(h.log, w)

	gateID, err := pathInt64(r, "id")
	if err != nil {
		response.ErrorResponse(err, "invalid gate id")
		return
	}

	if err := h.warehousesService.ActivateGate(r.Context(), gateID); err != nil {
		response.ErrorResponse(err, "activate gate")
		return
	}

	response.NoContentResponse()
}

func (h *WarehousesHTTPHandler) DeactivateGate(w http.ResponseWriter, r *http.Request) {
	response := core_http_response.NewHTTPResponseHandler(h.log, w)

	gateID, err := pathInt64(r, "id")
	if err != nil {
		response.ErrorResponse(err, "invalid gate id")
		return
	}

	if err := h.warehousesService.DeactivateGate(r.Context(), gateID); err != nil {
		response.ErrorResponse(err, "deactivate gate")
		return
	}

	response.NoContentResponse()
}
