package users_transport_http

import (
	"context"
	"net/http"

	core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"
	core_logger "github.com/pkpal-uhobp/fulfillment-app/internal/core/logger"
	core_http_server "github.com/pkpal-uhobp/fulfillment-app/internal/core/transport/http/server"
	users_service "github.com/pkpal-uhobp/fulfillment-app/internal/features/users/service"
)

type UsersHTTPHandler struct {
	log          *core_logger.Logger
	usersService UsersService
}

type UsersService interface {
	ListUsers(ctx context.Context, actorRole string, filter users_service.UserFilter) ([]users_service.UserDTO, error)
	CreateUser(ctx context.Context, actorRole string, input users_service.CreateUserInput) (users_service.UserDTO, error)
	PatchUser(ctx context.Context, actorID int64, actorRole string, userID int64, input users_service.PatchUserInput) (users_service.UserDTO, error)
	BlockUser(ctx context.Context, actorID int64, actorRole string, userID int64, input users_service.BlockUserInput) error
	DeactivateUser(ctx context.Context, actorID int64, actorRole string, userID int64) error
}

func NewUsersHTTPHandler(
	log *core_logger.Logger,
	usersService UsersService,
) *UsersHTTPHandler {
	return &UsersHTTPHandler{
		log:          log,
		usersService: usersService,
	}
}

func (h *UsersHTTPHandler) Routes() []core_http_server.Route {
	adminRoles := []string{core_domain.RoleAdmin.String()}
	return []core_http_server.Route{
		core_http_server.NewRoute(http.MethodGet, "/users", h.ListUsers, adminRoles),
		core_http_server.NewRoute(http.MethodPost, "/users", h.CreateUser, adminRoles),
		core_http_server.NewRoute(http.MethodPatch, "/users/{id}", h.PatchUser, adminRoles),
		core_http_server.NewRoute(http.MethodPatch, "/users/{id}/block", h.BlockUser, adminRoles),
		core_http_server.NewRoute(http.MethodDelete, "/users/{id}", h.DeactivateUser, adminRoles),
	}
}
