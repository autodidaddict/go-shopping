package redis

import "github.com/autodidaddict/go-shopping/shipping/proto"

// ShippingRepository represents a Redis shipping repo implementation
type ShippingRepository struct {
	redisDialString string
}

// NewRedisRepository creates a new CatalogRepository
func NewRedisRepository(redisDialString string) *ShippingRepository {
	return &ShippingRepository{redisDialString: redisDialString}
}

// GetShippingCosts returns the shipping cost options for a particular product to a particular zip code
func (r *ShippingRepository) GetShippingCosts(sku string, zipCode string) (costs []*shipping.ShippingCost, err error) {
	return
}

// MarkShipped marks a particular product within an order as shipped
func (r *ShippingRepository) MarkShipped(sku string, orderID uint64, note string, shippingMethod shipping.ShippingMethod) (trackingNumber string, err error) {
	return
}

// GetShippingStatus queries the shipping method and tracking number for a particular order item
func (r *ShippingRepository) GetShippingStatus(orderID uint64) (shippingStatus *shipping.ShippingStatus, err error) {
	return
}

// ProductExists indicates whether a product exists
func (r *ShippingRepository) ProductExists(sku string) (exists bool, err error) {
	return
}

// OrderExists indicates whether an order exists
func (r *ShippingRepository) OrderExists(orderID uint64) (exists bool, err error) {
	return
}
