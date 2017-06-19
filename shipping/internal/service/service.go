package service

import (
	"github.com/autodidaddict/go-shopping/shipping/proto"
	"golang.org/x/net/context"
)

type shippingService struct{}

// NewShippingService creates a new shipping service
func NewShippingService() shipping.ShippingHandler {
	return &shippingService{}
}

func (s *shippingService) GetShippingCost(ctx context.Context, request *shipping.ShippingCostRequest,
	response *shipping.ShippingCostResponse) error {
	return nil
}

func (s *shippingService) MarkItemShipped(ctx context.Context, request *shipping.MarkShippedRequest,
	response *shipping.MarkShippedResponse) error {
	return nil
}

func (s *shippingService) GetShippingStatus(ctx context.Context, request *shipping.ShippingStatusRequest,
	response *shipping.ShippingStatusResponse) error {
	return nil
}
