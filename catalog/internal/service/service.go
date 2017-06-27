package service

import (
	"github.com/autodidaddict/go-shopping/catalog/proto"
	"github.com/micro/go-micro/errors"
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

	if request == nil {
		return errors.BadRequest("", "Missing detail request")
	}
	exists, err := c.catalogRepo.ProductExists(request.Sku)
	if err != nil {
		return errors.InternalServerError("", "Failed to check product existence: %s", err.Error())
	}
	if !exists {
		return errors.NotFound(request.Sku, "No such product")
	}
	results, err := c.catalogRepo.GetProduct(request.Sku)
	if err != nil {
		return errors.InternalServerError(request.Sku, "Failed to fetch product: %s", err.Error())
	}

	response.Product = results
	return nil
}

func (c *catalogService) GetProductCategories(ctx context.Context, request *catalog.AllCategoriesRequest,
	response *catalog.AllCategoriesResponse) error {

	if request == nil {
		return errors.BadRequest("", "Missing categories request")
	}
	results, err := c.catalogRepo.GetCategories()
	if err != nil {
		return errors.InternalServerError("", "Failed to load categories: %s", err.Error())
	}
	response.Categories = results
	return nil
}

func (c *catalogService) GetProductsInCategory(ctx context.Context, request *catalog.CategoryProductsRequest,
	response *catalog.CategoryProductsResponse) error {

	if request == nil {
		return errors.BadRequest("", "Missing category products request")
	}
	exists, err := c.catalogRepo.CategoryExists(request.CategoryId)
	if err != nil {
		return errors.InternalServerError("", "Failed to check category existence: %s", err.Error())
	}
	if !exists {
		return errors.NotFound(string(request.CategoryId), "No such category")
	}

	results, err := c.catalogRepo.GetProductsInCategory(request.CategoryId)
	if err != nil {
		return errors.InternalServerError("", "Failed to load products in category: %s", err.Error())
	}
	response.Products = results
	return nil
}

func (c *catalogService) ProductSearch(ctx context.Context, request *catalog.SearchRequest,
	response *catalog.SearchResponse) error {

	if request == nil {
		return errors.BadRequest("", "Missing product search request")
	}
	if !validateSearchTerm(request.SearchTerm) {
		return errors.BadRequest("", "Invalid search term")
	}
	results, repoErr := c.catalogRepo.Find(request.SearchTerm, request.Categories)
	if repoErr != nil {
		return errors.InternalServerError("", "Failed to perform search: %s", repoErr.Error())
	}
	response.SearchResults = results

	return nil
}

func validateSearchTerm(term string) bool {
	return len(term) > 2
}
