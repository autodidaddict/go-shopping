package service

import (
	"github.com/autodidaddict/go-shopping/warehouse/proto"
	"github.com/micro/go-micro/errors"
	"golang.org/x/net/context"
)

type warehouseService struct {
	repo warehouseRepository
}

type warehouseRepository interface {
	GetWarehouseDetails(sku string) (details *warehouse.WarehouseDetails, err error)
	SkuExists(sku string) (exists bool, err error)
}

// NewWarehouseService returns an instance of a warehouse handler
func NewWarehouseService(repo warehouseRepository) warehouse.WarehouseHandler {
	return &warehouseService{repo: repo}
}

func (w *warehouseService) GetWarehouseDetails(ctx context.Context, request *warehouse.DetailsRequest,
	response *warehouse.DetailsResponse) error {

	if request == nil {
		return errors.BadRequest("", "Missing details request")
	}
	if len(request.Sku) < 6 {
		return errors.BadRequest("", "Invalid SKU")
	}
	exists, err := w.repo.SkuExists(request.Sku)
	if err != nil {
		return errors.InternalServerError(request.Sku, "Failed to check for SKU existence: %s", err)
	}
	if !exists {
		return errors.NotFound(request.Sku, "No such SKU")
	}

	details, err := w.repo.GetWarehouseDetails(request.Sku)
	if err != nil {
		return errors.InternalServerError(request.Sku, "Failed to query warehouse details: %s", err)
	}

	response.Details = details

	return nil
}
