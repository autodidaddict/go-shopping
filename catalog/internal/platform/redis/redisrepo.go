package redis

import "github.com/autodidaddict/go-shopping/catalog/proto"

// CatalogRepository is a Redis-backed product catalog repository
type CatalogRepository struct {
}

// NewRedisRepository creates a new CatalogRepository
func NewRedisRepository() *CatalogRepository {
	return &CatalogRepository{}
}

// GetProduct retrieves a single product from the repository
func (r *CatalogRepository) GetProduct(sku string) (product *catalog.Product, err error) {
	return
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
