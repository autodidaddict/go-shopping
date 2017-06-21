package service

import (
	"github.com/autodidaddict/go-shopping/catalog/internal/platform/errors"
	"github.com/autodidaddict/go-shopping/catalog/proto"
	"net/http"
)

func validateSearchTerm(searchTerm string) bool {
	if len(searchTerm) < 3 {
		return false
	}
	return true
}

func generateServiceError(err error) *catalog.Error {
	var catalogError *catalog.Error
	switch err {
	case errors.BadSearchTerm:
		{
			catalogError = &catalog.Error{HttpHint: http.StatusBadRequest,
				Code:    catalog.ErrorCode_BADSEARCHREQUEST,
				Message: err.Error()}

		}
	case errors.NoSuchCategory:
		{
			catalogError = &catalog.Error{HttpHint: http.StatusNotFound,
				Code:    catalog.ErrorCode_NOSUCHCATEGORY,
				Message: err.Error()}
		}
	case errors.NoSuchProduct:
		{
			catalogError = &catalog.Error{HttpHint: http.StatusNotFound,
				Code:    catalog.ErrorCode_NOSUCHSKU,
				Message: err.Error()}
		}
	default:
		{
			catalogError = &catalog.Error{HttpHint: http.StatusInternalServerError,
				Code:    catalog.ErrorCode_UNKNOWN,
				Message: err.Error()}
		}
	}
	return catalogError
}
