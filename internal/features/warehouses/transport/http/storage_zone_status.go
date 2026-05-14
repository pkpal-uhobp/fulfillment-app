package warehouses_transport_http

import (
	"net/http"

	core_http_response "github.com/pkpal-uhobp/fulfillment-app/internal/core/transport/http/response"
)

func (h *WarehousesHTTPHandler) ActivateStorageZone(w http.ResponseWriter, r *http.Request) {
	response := core_http_response.NewHTTPResponseHandler(h.log, w)

	zoneID, err := pathInt64(r, "id")
	if err != nil {
		response.ErrorResponse(err, "invalid storage zone id")
		return
	}

	if err := h.warehousesService.ActivateStorageZone(r.Context(), zoneID); err != nil {
		response.ErrorResponse(err, "activate storage zone")
		return
	}

	response.NoContentResponse()
}

func (h *WarehousesHTTPHandler) DeactivateStorageZone(w http.ResponseWriter, r *http.Request) {
	response := core_http_response.NewHTTPResponseHandler(h.log, w)

	zoneID, err := pathInt64(r, "id")
	if err != nil {
		response.ErrorResponse(err, "invalid storage zone id")
		return
	}

	if err := h.warehousesService.DeactivateStorageZone(r.Context(), zoneID); err != nil {
		response.ErrorResponse(err, "deactivate storage zone")
		return
	}

	response.NoContentResponse()
}
