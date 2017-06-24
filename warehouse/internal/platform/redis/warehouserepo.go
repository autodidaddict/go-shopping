package redis

import (
	"fmt"
	"github.com/autodidaddict/go-shopping/warehouse/proto"
	"github.com/garyburd/redigo/redis"
)

// WarehouseRepository represents a redis repository over warehouse data
type WarehouseRepository struct {
	redisDialString string
}

// NewWarehouseRepository creates a new warehouse repo
func NewWarehouseRepository(redisDialString string) *WarehouseRepository {
	return &WarehouseRepository{redisDialString: redisDialString}
}

// GetWarehouseDetails queries the information about physical inventory in the warehouse for a given SKU
func (r *WarehouseRepository) GetWarehouseDetails(sku string) (details *warehouse.WarehouseDetails, err error) {
	c, err := redis.Dial("tcp", r.redisDialString)
	if err != nil {
		return nil, err
	}
	defer c.Close()
	warehouseKey := fmt.Sprintf("warehouse:%s", sku)
	res, err := redis.Values(c.Do("HGETALL", warehouseKey))
	var itemDetails redisWarehouseDetails
	err = redis.ScanStruct(res, &itemDetails)
	if err != nil {
		return nil, err
	}
	stockKey := fmt.Sprintf("warehouse:%s:stock", sku)
	stockCount, err := redis.Int(c.Do("GET", stockKey))
	if err != nil {
		return nil, err
	}
	details = &warehouse.WarehouseDetails{
		Sku:            itemDetails.SKU,
		Manufacturer:   itemDetails.Manufacturer,
		ModelNumber:    itemDetails.ModelNumber,
		StockRemaining: uint32(stockCount),
	}
	return details, nil
}

// SkuExists indicates whether the SKU exists in the warehouse inventory (regardless of in-stock quantity)
func (r *WarehouseRepository) SkuExists(sku string) (exists bool, err error) {
	c, err := redis.Dial("tcp", r.redisDialString)
	if err != nil {
		return false, err
	}
	defer c.Close()
	warehouseKey := fmt.Sprintf("warehouse:%s", sku)
	exists, err = redis.Bool(c.Do("EXISTS", warehouseKey))
	return exists, err
}

// DecrementStock will reduce the on-hand quantity of a SKU by 1
func (r *WarehouseRepository) DecrementStock(sku string) (err error) {
	c, err := redis.Dial("tcp", r.redisDialString)
	if err != nil {
		return err
	}
	defer c.Close()
	warehouseKey := fmt.Sprintf("warehouse:%s:stock", sku)
	_, err = c.Do("INCRBY", warehouseKey, "-1")
	if err != nil {
		return err
	}

	return nil
}

type redisWarehouseDetails struct {
	SKU          string `redis:"sku"`
	Manufacturer string `redis:"mfr"`
	ModelNumber  string `redis:"model"`
}
