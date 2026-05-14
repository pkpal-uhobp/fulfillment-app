package core_http_utils

import (
	"fmt"
	"net/http"
	"strconv"

	core_errors "github.com/pkpal-uhobp/fulfillment-app/internal/core/errors"
)

func PathInt64(r *http.Request, name string) (int64, error) {
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
