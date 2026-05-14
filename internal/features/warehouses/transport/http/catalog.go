package warehouses_transport_http

import (
	"net/http"

	core_http_response "github.com/pkpal-uhobp/fulfillment-app/internal/core/transport/http/response"
)

func (h *WarehousesHTTPHandler) ListProductTypes(w http.ResponseWriter, r *http.Request) {
	response := core_http_response.NewHTTPResponseHandler(h.log, w)

	productTypes, err := h.warehousesService.ListProductTypes(r.Context())
	if err != nil {
		response.ErrorResponse(err, "list product types")
		return
	}

	response.JSONResponse(ProductTypesResponse{ProductTypes: productTypes}, http.StatusOK)
}

func (h *WarehousesHTTPHandler) ListCargoPlaceTypes(w http.ResponseWriter, r *http.Request) {
	response := core_http_response.NewHTTPResponseHandler(h.log, w)

	cargoPlaceTypes, err := h.warehousesService.ListCargoPlaceTypes(r.Context())
	if err != nil {
		response.ErrorResponse(err, "list cargo place types")
		return
	}

	response.JSONResponse(CargoPlaceTypesResponse{CargoPlaceTypes: cargoPlaceTypes}, http.StatusOK)
}
