package service

import (
	"github.com/autodidaddict/go-shopping/catalog/proto"
	"golang.org/x/net/context"
)

type catalogService struct {
	catalogRepo catalogRepository
}

type catalogRepository interface {
	GetProduct(sku string) (product *catalog.Product, err error)
	GetCategories() (categories []*catalog.ProductCategory, err error)
	GetProductsInCategory(categoryID uint64) (products []*catalog.Product, err error)
	Find(searchTerm string, categories []uint64) (products []*catalog.Product, err error)
}

// NewCatalogService creates a new catalog service
func NewCatalogService(catalogRepo catalogRepository) catalog.CatalogHandler {
	return &catalogService{
		catalogRepo: catalogRepo,
	}
}

func (c *catalogService) GetProductDetails(ctx context.Context, request *catalog.DetailRequest,
	response *catalog.DetailResponse) error {
	return nil
}

func (c *catalogService) GetProductCategories(ctx context.Context, request *catalog.AllCategoriesRequest,
	response *catalog.AllCategoriesResponse) error {
	return nil
}

func (c *catalogService) GetProductsInCategory(ctx context.Context, request *catalog.CategoryProductsRequest,
	response *catalog.CategoryProductsResponse) error {
	return nil
}

func (c *catalogService) ProductSearch(ctx context.Context, request *catalog.SearchRequest,
	response *catalog.SearchResponse) error {
	return nil
}
