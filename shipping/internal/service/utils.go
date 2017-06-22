package service

import (
	"github.com/autodidaddict/go-shopping/shipping/internal/platform/errors"
	"github.com/autodidaddict/go-shopping/shipping/proto"
	"net/http"
)

func generateServiceError(err error) *shipping.Error {
	var shippingError *shipping.Error
	switch err {

	case errors.NoSuchOrder:
		{
			shippingError = &shipping.Error{HttpHint: http.StatusNotFound,
				Code:    shipping.ErrorCode_NOSUCHORDER,
				Message: err.Error()}
		}
	case errors.NoSuchSKU:
		{
			shippingError = &shipping.Error{HttpHint: http.StatusNotFound,
				Code:    shipping.ErrorCode_NOSUCHSKU,
				Message: err.Error()}
		}
	default:
		{
			shippingError = &shipping.Error{HttpHint: http.StatusInternalServerError,
				Code:    shipping.ErrorCode_UNKNOWN,
				Message: err.Error()}
		}
	}
	return shippingError
}
