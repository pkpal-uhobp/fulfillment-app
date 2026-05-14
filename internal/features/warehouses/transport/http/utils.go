package warehouses_transport_http

import (
	"fmt"
	"net/http"
	"strconv"

	core_errors "github.com/pkpal-uhobp/fulfillment-app/internal/core/errors"
)

func pathInt64(r *http.Request, name string) (int64, error) {
	value := r.PathValue(name)
	if value == "" {
		return 0, fmt.Errorf("%w: missing path value %s", core_errors.ErrInvalidArgument, name)
	}

	id, err := strconv.ParseInt(value, 10, 64)
	if err != nil || id <= 0 {
		return 0, fmt.Errorf("%w: invalid path value %s", core_errors.ErrInvalidArgument, name)
	}

	return id, nil
}

func queryInt64(r *http.Request, name string) (int64, error) {
	value := r.URL.Query().Get(name)
	if value == "" {
		return 0, nil
	}

	id, err := strconv.ParseInt(value, 10, 64)
	if err != nil || id <= 0 {
		return 0, fmt.Errorf("%w: invalid query param %s", core_errors.ErrInvalidArgument, name)
	}

	return id, nil
}
