package users_transport_http

import (
	"fmt"
	"net/http"
	"strconv"

	core_errors "github.com/pkpal-uhobp/fulfillment-app/internal/core/errors"
	core_http_middleware "github.com/pkpal-uhobp/fulfillment-app/internal/core/transport/http/middleware"
	core_http_request "github.com/pkpal-uhobp/fulfillment-app/internal/core/transport/http/request"
	core_http_response "github.com/pkpal-uhobp/fulfillment-app/internal/core/transport/http/response"
	core_http_utils "github.com/pkpal-uhobp/fulfillment-app/internal/core/transport/http/utils"
	users_service "github.com/pkpal-uhobp/fulfillment-app/internal/features/users/service"
)

func (h *UsersHTTPHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	response := core_http_response.NewHTTPResponseHandler(h.log, w)

	actor, err := core_http_middleware.CurrentUserFromContext(r.Context())
	if err != nil {
		response.ErrorResponse(err, "get current user")
		return
	}

	isActive, err := queryBoolPtr(r, "is_active")
	if err != nil {
		response.ErrorResponse(err, "invalid is_active")
		return
	}

	isBlocked, err := queryBoolPtr(r, "is_blocked")
	if err != nil {
		response.ErrorResponse(err, "invalid is_blocked")
		return
	}

	page, err := queryInt(r, "page")
	if err != nil {
		response.ErrorResponse(err, "invalid page")
		return
	}

	limit, err := queryInt(r, "limit")
	if err != nil {
		response.ErrorResponse(err, "invalid limit")
		return
	}

	users, err := h.usersService.ListUsers(
		r.Context(),
		actor.Role,
		users_service.UserFilter{
			Role:      core_http_utils.QueryString(r, "role"),
			IsActive:  isActive,
			IsBlocked: isBlocked,
			Search:    core_http_utils.QueryString(r, "search"),
			Page:      page,
			Limit:     limit,
		},
	)
	if err != nil {
		response.ErrorResponse(err, "list users")
		return
	}

	response.JSONResponse(UsersResponse{Users: users}, http.StatusOK)
}

func (h *UsersHTTPHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	response := core_http_response.NewHTTPResponseHandler(h.log, w)

	actor, err := core_http_middleware.CurrentUserFromContext(r.Context())
	if err != nil {
		response.ErrorResponse(err, "get current user")
		return
	}

	var request CreateUserRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		response.ErrorResponse(err, "invalid create user request")
		return
	}

	user, err := h.usersService.CreateUser(
		r.Context(),
		actor.Role,
		users_service.CreateUserInput{
			Email:    request.Email,
			Password: request.Password,
			FullName: request.FullName,
			Phone:    request.Phone,
			Role:     request.Role,
		},
	)
	if err != nil {
		response.ErrorResponse(err, "create user")
		return
	}

	response.JSONResponse(UserResponse{User: user}, http.StatusCreated)
}

func (h *UsersHTTPHandler) PatchUser(w http.ResponseWriter, r *http.Request) {
	response := core_http_response.NewHTTPResponseHandler(h.log, w)

	actor, err := core_http_middleware.CurrentUserFromContext(r.Context())
	if err != nil {
		response.ErrorResponse(err, "get current user")
		return
	}

	userID, err := core_http_utils.PathInt64(r, "id")
	if err != nil {
		response.ErrorResponse(err, "invalid user id")
		return
	}

	var request PatchUserRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		response.ErrorResponse(err, "invalid patch user request")
		return
	}

	user, err := h.usersService.PatchUser(
		r.Context(),
		actor.ID,
		actor.Role,
		userID,
		users_service.PatchUserInput{
			Email:     request.Email,
			Password:  request.Password,
			FullName:  request.FullName,
			Phone:     request.Phone,
			Role:      request.Role,
			IsActive:  request.IsActive,
			IsBlocked: request.IsBlocked,
		},
	)
	if err != nil {
		response.ErrorResponse(err, "patch user")
		return
	}

	response.JSONResponse(UserResponse{User: user}, http.StatusOK)
}

func (h *UsersHTTPHandler) BlockUser(w http.ResponseWriter, r *http.Request) {
	response := core_http_response.NewHTTPResponseHandler(h.log, w)

	actor, err := core_http_middleware.CurrentUserFromContext(r.Context())
	if err != nil {
		response.ErrorResponse(err, "get current user")
		return
	}

	userID, err := core_http_utils.PathInt64(r, "id")
	if err != nil {
		response.ErrorResponse(err, "invalid user id")
		return
	}

	var request BlockUserRequest
	if r.Body != http.NoBody {
		if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
			response.ErrorResponse(err, "invalid block user request")
			return
		}
	}

	if err := h.usersService.BlockUser(
		r.Context(),
		actor.ID,
		actor.Role,
		userID,
		users_service.BlockUserInput{Reason: request.Reason},
	); err != nil {
		response.ErrorResponse(err, "block user")
		return
	}

	response.NoContentResponse()
}

func (h *UsersHTTPHandler) DeactivateUser(w http.ResponseWriter, r *http.Request) {
	response := core_http_response.NewHTTPResponseHandler(h.log, w)

	actor, err := core_http_middleware.CurrentUserFromContext(r.Context())
	if err != nil {
		response.ErrorResponse(err, "get current user")
		return
	}

	userID, err := core_http_utils.PathInt64(r, "id")
	if err != nil {
		response.ErrorResponse(err, "invalid user id")
		return
	}

	if err := h.usersService.DeactivateUser(
		r.Context(),
		actor.ID,
		actor.Role,
		userID,
	); err != nil {
		response.ErrorResponse(err, "deactivate user")
		return
	}

	response.NoContentResponse()
}

func queryBoolPtr(r *http.Request, name string) (*bool, error) {
	value := r.URL.Query().Get(name)
	if value == "" {
		return nil, nil
	}

	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid query param %s", core_errors.ErrInvalidArgument, name)
	}

	return &parsed, nil
}

func queryInt(r *http.Request, name string) (int, error) {
	value, err := core_http_utils.QueryInt64Ptr(r, name)
	if err != nil {
		return 0, err
	}
	if value == nil {
		return 0, nil
	}

	return int(*value), nil
}
