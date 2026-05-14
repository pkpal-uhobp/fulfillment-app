package core_http_utils

import (
	"fmt"
	"net/http"
	"strconv"

	core_errors "github.com/pkpal-uhobp/fulfillment-app/internal/core/errors"
)

func QueryString(r *http.Request, name string) string {
	return r.URL.Query().Get(name)
}

func QueryInt64Ptr(r *http.Request, name string) (*int64, error) {
	value := r.URL.Query().Get(name)
	if value == "" {
		return nil, nil
	}

	id, err := strconv.ParseInt(value, 10, 64)
	if err != nil || id <= 0 {
		return nil, fmt.Errorf("%w: invalid query param %s", core_errors.ErrInvalidArgument, name)
	}

	return &id, nil
}

func QueryInt64(r *http.Request, name string) (int64, error) {
	ptr, err := QueryInt64Ptr(r, name)
	if err != nil {
		return 0, err
	}

	if ptr == nil {
		return 0, nil
	}

	return *ptr, nil
}

