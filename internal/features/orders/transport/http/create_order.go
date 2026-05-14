package orders_transport_http

import (
	"net/http"

	core_http_middleware "github.com/pkpal-uhobp/fulfillment-app/internal/core/transport/http/middleware"
	core_http_request "github.com/pkpal-uhobp/fulfillment-app/internal/core/transport/http/request"
	core_http_response "github.com/pkpal-uhobp/fulfillment-app/internal/core/transport/http/response"
	orders_service "github.com/pkpal-uhobp/fulfillment-app/internal/features/orders/service"
)

func (h *OrdersHTTPHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	response := core_http_response.NewHTTPResponseHandler(h.log, w)

	user, err := core_http_middleware.CurrentUserFromContext(r.Context())
	if err != nil {
		response.ErrorResponse(err, "get current user")
		return
	}

	var request CreateOrderRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		response.ErrorResponse(err, "invalid create order request")
		return
	}

	cargoPlaces := make([]orders_service.CargoPlaceInput, 0, len(request.CargoPlaces))
	for _, cargoPlace := range request.CargoPlaces {
		cargoPlaces = append(cargoPlaces, orders_service.CargoPlaceInput{
			CargoPlaceTypeID: cargoPlace.CargoPlaceTypeID,
			Quantity:         cargoPlace.Quantity,
			WeightPerPlaceKG: cargoPlace.WeightPerPlaceKG,
			LengthCM:         cargoPlace.LengthCM,
			WidthCM:          cargoPlace.WidthCM,
			HeightCM:         cargoPlace.HeightCM,
			Comment:          cargoPlace.Comment,
		})
	}

	var pickup *orders_service.PickupInput
	if request.Pickup != nil {
		pickup = &orders_service.PickupInput{
			PickupAddress:  request.Pickup.PickupAddress,
			PickupDate:     request.Pickup.PickupDate,
			PickupTimeFrom: request.Pickup.PickupTimeFrom,
			PickupTimeTo:   request.Pickup.PickupTimeTo,
			ContactName:    request.Pickup.ContactName,
			ContactPhone:   request.Pickup.ContactPhone,
			Comment:        request.Pickup.Comment,
		}
	}

	order, err := h.ordersService.CreateOrder(
		r.Context(),
		user.ID,
		user.Role,
		orders_service.CreateOrderInput{
			ReceivingWarehouseID:   request.ReceivingWarehouseID,
			DestinationWarehouseID: request.DestinationWarehouseID,
			ProductTypeID:          request.ProductTypeID,
			HandoverType:           request.HandoverType,
			SelfDeliveryDate:       request.SelfDeliveryDate,
			SelfDeliveryTimeFrom:   request.SelfDeliveryTimeFrom,
			SelfDeliveryTimeTo:     request.SelfDeliveryTimeTo,
			Comment:                request.Comment,
			CargoPlaces:            cargoPlaces,
			Pickup:                 pickup,
		},
	)
	if err != nil {
		response.ErrorResponse(err, "create order")
		return
	}

	response.JSONResponse(
		OrderResponse{
			Order: order,
		},
		http.StatusCreated,
	)
}
