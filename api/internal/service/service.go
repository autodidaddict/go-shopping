package service

import (
	stderrors "errors"
	"github.com/autodidaddict/go-shopping/catalog/proto"
	"github.com/autodidaddict/go-shopping/shipping/proto"
	"github.com/autodidaddict/go-shopping/warehouse/proto"
	"github.com/emicklei/go-restful"
	"github.com/micro/go-log"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/errors"
	"golang.org/x/net/context"
	"net/http"
)

const (
	catalogService   = "go.shopping.srv.catalog"
	shippingService  = "go.shopping.srv.shipping"
	warehouseService = "go.shopping.srv.warehouse"
)

type CommerceService struct {
	warehouseClient warehouse.WarehouseClient
	shippingClient  shipping.ShippingClient
	catalogClient   catalog.CatalogClient
}

type catalogResults struct {
	catalogResponse *catalog.DetailResponse
	err             error
}

type warehouseResults struct {
	warehouseResponse *warehouse.DetailsResponse
	err               error
}

func NewCommerceService(c client.Client) *CommerceService {
	return &CommerceService{
		warehouseClient: warehouse.NewWarehouseClient(warehouseService, c),
		shippingClient:  shipping.NewShippingClient(shippingService, c),
		catalogClient:   catalog.NewCatalogClient(catalogService, c),
	}
}

func (cs *CommerceService) GetProductDetails(request *restful.Request, response *restful.Response) {

	sku := request.PathParameter("sku")
	log.Logf("Received request for product details: %s", sku)
	ctx := context.Background()
	catalogCh := cs.getCatalogDetails(ctx, sku)
	warehouseCh := cs.getWarehouseDetails(ctx, sku)

	catalogReply := <-catalogCh
	if catalogReply.err != nil {
		writeError(response, catalogReply.err)
		return
	}

	warehouseReply := <-warehouseCh
	if warehouseReply.err != nil {
		writeError(response, warehouseReply.err)
		return
	}
	product := catalogReply.catalogResponse.Product

	details := productDetails{
		SKU:            product.SKU,
		StockRemaining: warehouseReply.warehouseResponse.Details.StockRemaining,
		Manufacturer:   product.Manufacturer,
		Price:          product.Price,
		Model:          product.Model,
		Name:           product.Name,
		Description:    product.Description,
	}
	response.WriteEntity(details)
}

func (cs *CommerceService) getCatalogDetails(ctx context.Context, sku string) chan catalogResults {
	ch := make(chan catalogResults, 1)

	go func() {
		res, err := cs.catalogClient.GetProductDetails(ctx, &catalog.DetailRequest{SKU: sku})
		ch <- catalogResults{catalogResponse: res, err: err}
	}()

	return ch
}

func (cs *CommerceService) getWarehouseDetails(ctx context.Context, sku string) chan warehouseResults {
	ch := make(chan warehouseResults, 1)

	go func() {
		res, err := cs.warehouseClient.GetWarehouseDetails(ctx, &warehouse.DetailsRequest{Sku: sku})
		ch <- warehouseResults{warehouseResponse: res, err: err}
	}()

	return ch
}

func writeError(response *restful.Response, err error) {
	realError := errors.Parse(err.Error())
	if realError != nil {
		response.WriteError(int(realError.Code), stderrors.New(realError.Detail))
		return
	}
	response.WriteError(http.StatusInternalServerError, err)

}
