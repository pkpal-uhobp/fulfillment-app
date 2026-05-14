package warehouses_transport_http

import (
	"net/http"

	core_http_response "github.com/pkpal-uhobp/fulfillment-app/internal/core/transport/http/response"
)

func (h *WarehousesHTTPHandler) ActivateWarehouse(w http.ResponseWriter, r *http.Request) {
	response := core_http_response.NewHTTPResponseHandler(h.log, w)

	warehouseID, err := pathInt64(r, "id")
	if err != nil {
		response.ErrorResponse(err, "invalid warehouse id")
		return
	}

	if err := h.warehousesService.ActivateWarehouse(r.Context(), warehouseID); err != nil {
		response.ErrorResponse(err, "activate warehouse")
		return
	}

	response.NoContentResponse()
}
