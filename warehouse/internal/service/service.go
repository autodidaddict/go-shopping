package service

import (
	"github.com/autodidaddict/go-shopping/warehouse/proto"
	"golang.org/x/net/context"
)

type warehouseService struct{}

// NewWarehouseService returns an instance of a warehouse handler
func NewWarehouseService() warehouse.WarehouseHandler {
	return &warehouseService{}
}

func (w *warehouseService) GetWarehouseDetails(ctx context.Context, request *warehouse.DetailsRequest,
	response *warehouse.DetailsResponse) error {
	response.Manufacturer = "TOSHIBA"
	response.ModelNumber = "T-1000"
	response.SKU = request.SKU
	response.StockRemaining = 35
	return nil
}
