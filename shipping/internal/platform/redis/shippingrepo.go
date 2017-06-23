package redis

import (
	"fmt"
	"github.com/autodidaddict/go-shopping/shipping/proto"
	"github.com/garyburd/redigo/redis"
	"math/rand"
)

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
	// I don't feel like simulating shipping costs based on weight, distance, and provider so I'm
	// just going to return some bogus values here.
	return []*shipping.ShippingCost{
		&shipping.ShippingCost{
			Price:  2500,
			Method: shipping.ShippingMethod_SM_FEDEX,
		},
		&shipping.ShippingCost{
			Price:  1000,
			Method: shipping.ShippingMethod_SM_RAVEN,
		},
	}, nil
	return
}

// MarkShipped marks a particular product within an order as shipped
func (r *ShippingRepository) MarkShipped(sku string, orderID uint64, note string, shippingMethod shipping.ShippingMethod) (trackingNumber string, err error) {
	c, err := redis.Dial("tcp", r.redisDialString)
	if err != nil {
		return "", err
	}
	defer c.Close()
	itemKey := fmt.Sprintf("order:%d:shippingstatus:%s", orderID, sku)

	itemStatus := redisShippingStatus{
		Shipped:        true,
		TrackingNumber: randStringBytes(6),
		ShippingMethod: uint(shippingMethod),
	}

	_, err = c.Do("HMSET", redis.Args{}.Add(itemKey).AddFlat(&itemStatus)...)
	if err != nil {
		return "", err
	}

	return itemStatus.TrackingNumber, nil
}

// GetShippingStatus queries the shipping method and tracking number for a particular order item. Shipping status
// is stored in the database under order:{id}:shippingstatus:{sku} as a hashmap
func (r *ShippingRepository) GetShippingStatus(orderID uint64, sku string) (shippingStatus *shipping.ShippingStatus, err error) {
	c, err := redis.Dial("tcp", r.redisDialString)
	if err != nil {
		return nil, err
	}
	defer c.Close()
	itemKey := fmt.Sprintf("order:%d:shippingstatus:%s", orderID, sku)

	exists, err := redis.Bool(c.Do("EXISTS", itemKey))
	if err != nil {
		return nil, err
	}
	// If status hasn't been marked, just return a default "no status" (unshipped)
	if !exists {
		shippingStatus = &shipping.ShippingStatus{
			ShippingMethod: shipping.ShippingMethod_SM_NOTSHIPPED,
			TrackingNumber: "-",
			Shipped:        false,
		}
		return shippingStatus, nil
	}

	res, err := redis.Values(c.Do("HGETALL", itemKey))
	var itemStatus redisShippingStatus
	err = redis.ScanStruct(res, &itemStatus)
	if err != nil {
		return nil, err
	}
	shippingStatus = &shipping.ShippingStatus{
		ShippingMethod: shipping.ShippingMethod(itemStatus.ShippingMethod),
		TrackingNumber: itemStatus.TrackingNumber,
		Shipped:        itemStatus.Shipped,
	}
	return shippingStatus, nil
}

// ProductExists indicates whether a product exists
func (r *ShippingRepository) ProductExists(sku string) (exists bool, err error) {
	c, err := redis.Dial("tcp", r.redisDialString)
	if err != nil {
		return false, err
	}
	defer c.Close()
	productKey := fmt.Sprintf("product:%s", sku)
	exists, err = redis.Bool(c.Do("EXISTS", productKey))
	return exists, err
}

// OrderExists indicates whether an order exists
func (r *ShippingRepository) OrderExists(orderID uint64) (exists bool, err error) {
	c, err := redis.Dial("tcp", r.redisDialString)
	if err != nil {
		return false, err
	}
	defer c.Close()
	orderKey := fmt.Sprintf("order:%d", orderID)
	exists, err = redis.Bool(c.Do("EXISTS", orderKey))
	return exists, err
}

type redisShippingStatus struct {
	Shipped        bool   `redis:"shipped"`
	TrackingNumber string `redis:"tracking_number"`
	ShippingMethod uint   `redis:"shipping_method"`
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
