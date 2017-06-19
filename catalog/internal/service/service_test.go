package service_test

import (
	"github.com/autodidaddict/go-shopping/catalog/internal/service"
	"github.com/autodidaddict/go-shopping/catalog/proto"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestProductSearch(t *testing.T) {
	Convey("Given a catalog service", t, func() {
		svc := service.NewCatalogService(newFakeRepo())
	})
}

type fakeRepo struct {
	shouldFail bool
}

func newFakeRepo() *fakeRepo {
	return &fakeRepo{}
}

func (r *fakeRepo) GetProduct(sku string) (product *catalog.Product, err error) {
	product = &catalog.Product{}
	return
}

func (r *fakeRepo) GetCategories() (categories []*catalog.ProductCategory, err error) {
	return
}

func (r *fakeRepo) GetProductsInCategory(categoryID uint64) (products []*catalog.Product, err error) {
	return
}

func (r *fakeRepo) Find(searchTerm string, categories []uint64) (products []*catalog.Product, err error) {
	return
}
