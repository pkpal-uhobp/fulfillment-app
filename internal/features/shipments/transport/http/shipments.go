package shipments_transport_http

import (
	"net/http"

	core_http_middleware "github.com/pkpal-uhobp/fulfillment-app/internal/core/transport/http/middleware"
	core_http_request "github.com/pkpal-uhobp/fulfillment-app/internal/core/transport/http/request"
	core_http_response "github.com/pkpal-uhobp/fulfillment-app/internal/core/transport/http/response"
	core_http_utils "github.com/pkpal-uhobp/fulfillment-app/internal/core/transport/http/utils"
	shipments_service "github.com/pkpal-uhobp/fulfillment-app/internal/features/shipments/service"
)

func (h *ShipmentsHTTPHandler) CreateShipment(w http.ResponseWriter, r *http.Request) {
	response := core_http_response.NewHTTPResponseHandler(h.log, w)
	user, err := core_http_middleware.CurrentUserFromContext(r.Context())
	if err != nil {
		response.ErrorResponse(err, "get current user")
		return
	}
	var request CreateShipmentRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		response.ErrorResponse(err, "invalid create shipment request")
		return
	}
	shipment, err := h.shipmentsService.CreateShipment(
		r.Context(),
		user.ID,
		user.Role,
		shipments_service.CreateShipmentInput{
			DestinationWarehouseID: request.DestinationWarehouseID,
			GateID:                 request.GateID,
			PlannedDepartureAt:     request.PlannedDepartureAt,
		},
	)
	if err != nil {
		response.ErrorResponse(err, "create shipment")
		return
	}
	response.JSONResponse(ShipmentResponse{Shipment: shipment}, http.StatusCreated)
}

func (h *ShipmentsHTTPHandler) ListShipments(w http.ResponseWriter, r *http.Request) {
	response := core_http_response.NewHTTPResponseHandler(h.log, w)
	user, err := core_http_middleware.CurrentUserFromContext(r.Context())
	if err != nil {
		response.ErrorResponse(err, "get current user")
		return
	}
	destinationWarehouseID, err := core_http_utils.QueryInt64Ptr(r, "destination_warehouse_id")
	if err != nil {
		response.ErrorResponse(err, "invalid destination warehouse id")
		return
	}
	gateID, err := core_http_utils.QueryInt64Ptr(r, "gate_id")
	if err != nil {
		response.ErrorResponse(err, "invalid gate id")
		return
	}
	shipments, err := h.shipmentsService.ListShipments(
		r.Context(),
		user.ID,
		user.Role,
		shipments_service.ShipmentFilter{
			Status:                 core_http_utils.QueryString(r, "status"),
			DestinationWarehouseID: destinationWarehouseID,
			GateID:                 gateID,
			Date:                   core_http_utils.QueryString(r, "date"),
		},
	)
	if err != nil {
		response.ErrorResponse(err, "list shipments")
		return
	}
	response.JSONResponse(ShipmentsResponse{Shipments: shipments}, http.StatusOK)
}

func (h *ShipmentsHTTPHandler) GetShipment(w http.ResponseWriter, r *http.Request) {
	response := core_http_response.NewHTTPResponseHandler(h.log, w)
	user, err := core_http_middleware.CurrentUserFromContext(r.Context())
	if err != nil {
		response.ErrorResponse(err, "get current user")
		return
	}
	shipmentID, err := core_http_utils.PathInt64(r, "id")
	if err != nil {
		response.ErrorResponse(err, "invalid shipment id")
		return
	}
	shipment, err := h.shipmentsService.GetShipment(r.Context(), shipmentID, user.ID, user.Role)
	if err != nil {
		response.ErrorResponse(err, "get shipment")
		return
	}
	response.JSONResponse(ShipmentResponse{Shipment: shipment}, http.StatusOK)
}

func (h *ShipmentsHTTPHandler) AddShipmentItem(w http.ResponseWriter, r *http.Request) {
	response := core_http_response.NewHTTPResponseHandler(h.log, w)
	user, err := core_http_middleware.CurrentUserFromContext(r.Context())
	if err != nil {
		response.ErrorResponse(err, "get current user")
		return
	}
	shipmentID, err := core_http_utils.PathInt64(r, "id")
	if err != nil {
		response.ErrorResponse(err, "invalid shipment id")
		return
	}
	var request AddShipmentItemRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		response.ErrorResponse(err, "invalid add shipment item request")
		return
	}
	shipment, err := h.shipmentsService.AddShipmentItem(
		r.Context(),
		shipmentID,
		user.ID,
		user.Role,
		shipments_service.AddShipmentItemInput{
			CargoItemID: request.CargoItemID,
			Comment:     request.Comment,
		},
	)
	if err != nil {
		response.ErrorResponse(err, "add shipment item")
		return
	}
	response.JSONResponse(ShipmentResponse{Shipment: shipment}, http.StatusOK)
}

func (h *ShipmentsHTTPHandler) RemoveShipmentItem(w http.ResponseWriter, r *http.Request) {
	response := core_http_response.NewHTTPResponseHandler(h.log, w)
	user, err := core_http_middleware.CurrentUserFromContext(r.Context())
	if err != nil {
		response.ErrorResponse(err, "get current user")
		return
	}
	shipmentID, err := core_http_utils.PathInt64(r, "id")
	if err != nil {
		response.ErrorResponse(err, "invalid shipment id")
		return
	}
	cargoItemID, err := core_http_utils.PathInt64(r, "cargo_item_id")
	if err != nil {
		response.ErrorResponse(err, "invalid cargo item id")
		return
	}
	if err := h.shipmentsService.RemoveShipmentItem(r.Context(), shipmentID, cargoItemID, user.ID, user.Role); err != nil {
		response.ErrorResponse(err, "remove shipment item")
		return
	}
	response.NoContentResponse()
}

func (h *ShipmentsHTTPHandler) UpdateShipmentStatus(w http.ResponseWriter, r *http.Request) {
	response := core_http_response.NewHTTPResponseHandler(h.log, w)
	user, err := core_http_middleware.CurrentUserFromContext(r.Context())
	if err != nil {
		response.ErrorResponse(err, "get current user")
		return
	}
	shipmentID, err := core_http_utils.PathInt64(r, "id")
	if err != nil {
		response.ErrorResponse(err, "invalid shipment id")
		return
	}
	var request UpdateShipmentStatusRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		response.ErrorResponse(err, "invalid update shipment status request")
		return
	}
	shipment, err := h.shipmentsService.UpdateShipmentStatus(
		r.Context(),
		shipmentID,
		user.ID,
		user.Role,
		shipments_service.UpdateShipmentStatusInput{
			Status:  request.Status,
			Comment: request.Comment,
		},
	)
	if err != nil {
		response.ErrorResponse(err, "update shipment status")
		return
	}
	response.JSONResponse(ShipmentResponse{Shipment: shipment}, http.StatusOK)
}
