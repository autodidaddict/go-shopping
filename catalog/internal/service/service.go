package service

import (
	"github.com/autodidaddict/go-shopping/catalog/internal/platform/errors"
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
	CategoryExists(categoryID uint64) (bool, error)
	ProductExists(sku string) (bool, error)
}

// NewCatalogService creates a new catalog service
func NewCatalogService(catalogRepo catalogRepository) catalog.CatalogHandler {
	return &catalogService{
		catalogRepo: catalogRepo,
	}
}

func (c *catalogService) GetProductDetails(ctx context.Context, request *catalog.DetailRequest,
	response *catalog.DetailResponse) error {

	exists, err := c.catalogRepo.ProductExists(request.SKU)
	if err != nil {
		return err
	}
	if !exists {
		response.Error = generateServiceError(errors.NoSuchProduct)
		return nil
	}
	results, err := c.catalogRepo.GetProduct(request.SKU)
	if err != nil {
		return err
	}

	response.Product = results
	return nil
}

func (c *catalogService) GetProductCategories(ctx context.Context, request *catalog.AllCategoriesRequest,
	response *catalog.AllCategoriesResponse) error {

	results, err := c.catalogRepo.GetCategories()
	if err != nil {
		return err
	}
	response.Categories = results
	return nil
}

func (c *catalogService) GetProductsInCategory(ctx context.Context, request *catalog.CategoryProductsRequest,
	response *catalog.CategoryProductsResponse) error {

	exists, err := c.catalogRepo.CategoryExists(request.CategoryID)
	if err != nil {
		return err
	}
	if !exists {
		response.Error = generateServiceError(errors.NoSuchCategory)
		return nil
	}

	results, err := c.catalogRepo.GetProductsInCategory(request.CategoryID)
	if err != nil {
		return err
	}
	response.Products = results
	return nil
}

func (c *catalogService) ProductSearch(ctx context.Context, request *catalog.SearchRequest,
	response *catalog.SearchResponse) error {

	if !validateSearchTerm(request.SearchTerm) {
		response.Error = generateServiceError(errors.BadSearchTerm)
	}
	results, repoErr := c.catalogRepo.Find(request.SearchTerm, request.Categories)
	if repoErr != nil {
		return repoErr
	}
	response.SearchResults = results

	return nil
}
