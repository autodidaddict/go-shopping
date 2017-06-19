package service

import (
	"github.com/autodidaddict/go-shopping/catalog/proto"
	"golang.org/x/net/context"
)

type catalogService struct{}

// NewCatalogService creates a new catalog service
func NewCatalogService() catalog.CatalogHandler {
	return &catalogService{}
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
