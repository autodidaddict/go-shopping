package service_test

import (
	. "github.com/smartystreets/goconvey/convey"
	"golang.org/x/net/context"
	"testing"

	stderrors "errors"
	"github.com/autodidaddict/go-shopping/shipping/proto"
	"github.com/autodidaddict/go-shopping/warehouse/internal/service"
	"github.com/autodidaddict/go-shopping/warehouse/proto"
	"github.com/micro/go-micro/errors"
	"net/http"
)

func TestWarehouseService_GetWarehouseDetails(t *testing.T) {
	Convey("Given a warehouse service", t, func() {
		ctx := context.Background()
		stockChan := make(chan string)
		repo := &fakeRepo{stockChan: stockChan}
		shippedChannel := make(chan *shipping.ItemShippedEvent)
		svc := service.NewWarehouseService(repo, shippedChannel)

		Convey("requesting warehouse details should invoke the repository", func() {
			repo.shouldFail = false
			var resp warehouse.DetailsResponse
			err := svc.GetWarehouseDetails(ctx, &warehouse.DetailsRequest{Sku: "111111"}, &resp)
			So(err, ShouldBeNil)
			So(resp.Details, ShouldNotBeNil)
			So(resp.Details.Manufacturer, ShouldEqual, "TOSHIBA")
			So(resp.Details.StockRemaining, ShouldEqual, 42)
			So(resp.Details.ModelNumber, ShouldEqual, "T-1000")
		})

		Convey("requesting warehouse details should fail when the repo fails", func() {
			repo.shouldFail = true
			var resp warehouse.DetailsResponse
			err := svc.GetWarehouseDetails(ctx, &warehouse.DetailsRequest{Sku: "111111"}, &resp)
			So(err, ShouldNotBeNil)
			realError := errors.Parse(err.Error())
			So(realError, ShouldNotBeNil)
			So(realError.Code, ShouldEqual, http.StatusInternalServerError)
		})

		Convey("requesting warehouse details for a non-existent sku should fail with a 404", func() {
			repo.shouldFail = false
			var resp warehouse.DetailsResponse
			err := svc.GetWarehouseDetails(ctx, &warehouse.DetailsRequest{Sku: "nevergonnahappen"}, &resp)
			So(err, ShouldNotBeNil)
			realError := errors.Parse(err.Error())
			So(realError, ShouldNotBeNil)
			So(realError.Code, ShouldEqual, http.StatusNotFound)
		})

		Convey("requesting warehouse details with a nil request should fail", func() {
			repo.shouldFail = false
			var resp warehouse.DetailsResponse
			err := svc.GetWarehouseDetails(ctx, nil, &resp)
			So(err, ShouldNotBeNil)
			realError := errors.Parse(err.Error())
			So(realError, ShouldNotBeNil)
			So(realError.Code, ShouldEqual, http.StatusBadRequest)
		})

		Convey("requesting warehouse details with a bad SKU should fail", func() {
			repo.shouldFail = false
			var resp warehouse.DetailsResponse
			err := svc.GetWarehouseDetails(ctx, &warehouse.DetailsRequest{Sku: "1111"}, &resp) // SKU is too short
			So(err, ShouldNotBeNil)
			realError := errors.Parse(err.Error())
			So(realError, ShouldNotBeNil)
			So(realError.Code, ShouldEqual, http.StatusBadRequest)
		})

		Convey("when an item shipped event is received, it should ask the repo to decrement stock", func() {
			repo.shouldFail = false
			evt := &shipping.ItemShippedEvent{
				Sku:            "111111",
				ShippingMethod: shipping.ShippingMethod_SM_FEDEX,
				TrackingNumber: "abc1233",
			}
			shippedChannel <- evt
			s := <-stockChan
			So(s, ShouldEqual, "111111")
		})
	})
}

type fakeRepo struct {
	shouldFail bool
	stockChan  chan string
}

func (r *fakeRepo) GetWarehouseDetails(sku string) (details *warehouse.WarehouseDetails, err error) {
	if r.shouldFail {
		return nil, stderrors.New("Faily Fail")
	}
	return &warehouse.WarehouseDetails{
		ModelNumber:    "T-1000",
		StockRemaining: 42,
		Manufacturer:   "TOSHIBA",
		Sku:            "111111",
	}, nil
}

func (r *fakeRepo) DecrementStock(sku string) (err error) {
	r.stockChan <- sku
	return nil
}

func (r *fakeRepo) SkuExists(sku string) (exists bool, err error) {
	return sku == "111111", nil
}
