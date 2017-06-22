package redis

import (
	"fmt"
	"github.com/autodidaddict/go-shopping/catalog/proto"
	"github.com/garyburd/redigo/redis"
)

// CatalogRepository is a Redis-backed product catalog repository
type CatalogRepository struct {
	redisDialString string
}

// NewRedisRepository creates a new CatalogRepository
func NewRedisRepository(redisDialString string) *CatalogRepository {
	return &CatalogRepository{redisDialString: redisDialString}
}

// GetProduct retrieves a single product from the repository
func (r *CatalogRepository) GetProduct(sku string) (product *catalog.Product, err error) {

	c, err := redis.Dial("tcp", r.redisDialString)
	if err != nil {
		return nil, err
	}
	productKey := fmt.Sprintf("product:%s", sku)
	v, err := redis.Values(c.Do("HGETALL", productKey))
	if err != nil {
		return nil, err
	}
	var p redisProduct
	err = redis.ScanStruct(v, &p)
	if err != nil {
		return nil, err
	}

	product = &catalog.Product{
		p.SKU,
		p.Name,
		p.Description,
		p.Manufacturer,
		p.Model,
		p.Price}

	return product, nil
}

// GetCategories retrieves a list of product categories
func (r *CatalogRepository) GetCategories() (categories []*catalog.ProductCategory, err error) {
	return
}

// GetProductsInCategory retrieves a list of products within a given category
func (r *CatalogRepository) GetProductsInCategory(categoryID uint64) (products []*catalog.Product, err error) {
	return
}

// Find searches for `searchTerm` within the given list of categories.
func (r *CatalogRepository) Find(searchTerm string, categories []uint64) (products []*catalog.Product, err error) {
	return
}

// CategoryExists indicates whether a given category exists
func (r *CatalogRepository) CategoryExists(categoryID uint64) (exists bool, err error) {
	return
}

// ProductExists indicates whether a product exists
func (r *CatalogRepository) ProductExists(sku string) (exists bool, err error) {
	c, err := redis.Dial("tcp", r.redisDialString)
	if err != nil {
		return false, err
	}
	productKey := fmt.Sprintf("product:%s", sku)
	exists, err = redis.Bool(c.Do("EXISTS", productKey))
	return exists, err
}
