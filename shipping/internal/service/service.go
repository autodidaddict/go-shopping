package service

import (
	"github.com/autodidaddict/go-shopping/shipping/proto"
	"github.com/micro/go-micro/errors"
	"golang.org/x/net/context"
	"time"
)

type shippingService struct {
	repo           shippingRepository
	eventPublisher shippedEventPublisher
}

type shippingRepository interface {
	GetShippingCosts(sku string, zipCode string) (costs []*shipping.ShippingCost, err error)
	MarkShipped(sku string, orderID uint64, note string, shippingMethod shipping.ShippingMethod) (trackingNumber string, err error)
	GetShippingStatus(orderID uint64, sku string) (shippingStatus *shipping.ShippingStatus, err error)
	ProductExists(sku string) (exists bool, err error)
	OrderExists(orderID uint64) (exists bool, err error)
}

type shippedEventPublisher interface {
	PublishItemShippedEvent(event *shipping.ItemShippedEvent) (err error)
}

// NewShippingService creates a new shipping service
func NewShippingService(repo shippingRepository, publisher shippedEventPublisher) shipping.ShippingHandler {
	return &shippingService{repo: repo, eventPublisher: publisher}
}

func (s *shippingService) GetShippingCost(ctx context.Context, request *shipping.ShippingCostRequest,
	response *shipping.ShippingCostResponse) error {

	if request == nil {
		return errors.BadRequest("", "Missing shipping cost request")
	}
	exists, err := s.repo.ProductExists(request.Sku)
	if err != nil {
		return errors.InternalServerError("", "Failed to check product existence: %s", err)
	}
	if !exists {
		return errors.NotFound(request.Sku, "No such product")
	}

	shippingCosts, err := s.repo.GetShippingCosts(request.Sku, request.ZipCode)
	if err != nil {
		return errors.InternalServerError("", "Failed to retrieve shipping cost: %s", err)
	}
	response.ShippingCosts = shippingCosts
	return nil
}

func (s *shippingService) MarkItemShipped(ctx context.Context, request *shipping.MarkShippedRequest,
	response *shipping.MarkShippedResponse) error {

	if request == nil {
		return errors.BadRequest("", "Missing mark shipped request")
	}
	if request.ShippingMethod == shipping.ShippingMethod_SM_UNKNOWN {
		return errors.BadRequest("", "Must supply a valid shipping method")
	}
	exists, err := s.repo.OrderExists(request.OrderId)
	if err != nil {
		return errors.InternalServerError("", "Failed to check order existence: %s", err.Error())
	}
	if !exists {
		return errors.NotFound(string(request.OrderId), "No such order")
	}
	tracking, err := s.repo.MarkShipped(request.Sku, request.OrderId, request.Note, request.ShippingMethod)
	if err != nil {
		return errors.InternalServerError(string(request.OrderId), "Failed to mark item as shipped: %s", err.Error())
	}
	response.TrackingNumber = tracking

	err = s.eventPublisher.PublishItemShippedEvent(&shipping.ItemShippedEvent{
		TrackingNumber: tracking,
		OrderId:        request.OrderId,
		Note:           request.Note,
		ShippingMethod: request.ShippingMethod,
		Sku:            request.Sku,
		Timestamp:      time.Now().UTC().Unix(),
	})
	response.Success = err == nil

	return nil
}

func (s *shippingService) GetShippingStatus(ctx context.Context, request *shipping.ShippingStatusRequest,
	response *shipping.ShippingStatusResponse) error {

	if request == nil {
		return errors.BadRequest("", "Missing shipping status request")
	}
	exists, err := s.repo.OrderExists(request.OrderId)
	if err != nil {
		return errors.InternalServerError("", "Failed to check order existence: %s", err)
	}
	if !exists {
		return errors.NotFound(string(request.OrderId), "No such order")
	}

	status, err := s.repo.GetShippingStatus(request.OrderId, request.Sku)
	if err != nil {
		return errors.InternalServerError("", "Failed to query shipping status: %s", err)
	}
	response.ShippingStatus = status
	return nil
}
