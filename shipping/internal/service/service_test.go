package service_test

import (
	. "github.com/smartystreets/goconvey/convey"
	"golang.org/x/net/context"
	"testing"

	stderrors "errors"
	"github.com/autodidaddict/go-shopping/shipping/internal/service"
	"github.com/autodidaddict/go-shopping/shipping/proto"
	"github.com/micro/go-micro/errors"
	"net/http"
)

func TestShippingService_GetShippingCost(t *testing.T) {
	Convey("Given a shipping service", t, func() {
		ctx := context.Background()
		repo := &fakeRepo{}
		pub := &fakePublisher{}
		svc := service.NewShippingService(repo, pub)

		Convey("requesting shipping cost should invoke repository", func() {
			repo.shouldFail = false
			var resp shipping.ShippingCostResponse
			err := svc.GetShippingCost(ctx, &shipping.ShippingCostRequest{Sku: "8675309", ZipCode: "90210"}, &resp)
			So(err, ShouldBeNil)
			So(len(resp.ShippingCosts), ShouldEqual, 2)
			So(resp.ShippingCosts[0].Method, ShouldEqual, shipping.ShippingMethod_SM_FEDEX)
			So(resp.ShippingCosts[0].Price, ShouldEqual, 2500)
		})

		Convey("requesting a shipping cost for a non-existent sku should give us an appropriate error", func() {
			repo.shouldFail = false
			var resp shipping.ShippingCostResponse
			err := svc.GetShippingCost(ctx, &shipping.ShippingCostRequest{Sku: "notarealsku", ZipCode: "90210"}, &resp)
			So(err, ShouldNotBeNil)
			realError := errors.Parse(err.Error())
			So(realError, ShouldNotBeNil)
			So(realError.Code, ShouldEqual, http.StatusNotFound)
		})

		Convey("requesting a shipping cost should fail when the repo fails", func() {
			repo.shouldFail = true
			var resp shipping.ShippingCostResponse
			err := svc.GetShippingCost(ctx, &shipping.ShippingCostRequest{Sku: "8675309", ZipCode: "90210"}, &resp)
			So(err, ShouldNotBeNil)
			realError := errors.Parse(err.Error())
			So(realError, ShouldNotBeNil)
			So(realError.Code, ShouldEqual, http.StatusInternalServerError)
			So(realError.Detail, ShouldEqual, "Failed to retrieve shipping cost: Faily Fail")
		})

		Convey("requesting shipping cost with a null request should fail", func() {
			var resp shipping.ShippingCostResponse
			err := svc.GetShippingCost(ctx, nil, &resp)
			So(err, ShouldNotBeNil)
			realError := errors.Parse(err.Error())
			So(realError, ShouldNotBeNil)
			So(realError.Code, ShouldEqual, http.StatusBadRequest)
		})
	})
}

func TestShippingService_GetShippingStatus(t *testing.T) {
	Convey("Given a shipping service", t, func() {
		ctx := context.Background()
		repo := &fakeRepo{}
		pub := &fakePublisher{}
		svc := service.NewShippingService(repo, pub)

		Convey("requesting shipping status should invoke repository", func() {
			repo.shouldFail = false
			var resp shipping.ShippingStatusResponse
			err := svc.GetShippingStatus(ctx, &shipping.ShippingStatusRequest{OrderId: 42}, &resp)
			So(err, ShouldBeNil)
			So(resp.ShippingStatus, ShouldNotBeNil)
			So(resp.ShippingStatus.ShippingMethod, ShouldEqual, shipping.ShippingMethod_SM_RAVEN)
			So(resp.ShippingStatus.TrackingNumber, ShouldEqual, "111111")
		})

		Convey("requesting shipping status for non-existent order should fail", func() {
			repo.shouldFail = false
			var resp shipping.ShippingStatusResponse
			err := svc.GetShippingStatus(ctx, &shipping.ShippingStatusRequest{OrderId: 1}, &resp)
			So(err, ShouldNotBeNil)
			realError := errors.Parse(err.Error())
			So(realError, ShouldNotBeNil)
			So(realError.Code, ShouldEqual, http.StatusNotFound)
		})

		Convey("requesting shipping status when repo fails should fail", func() {
			repo.shouldFail = true
			var resp shipping.ShippingStatusResponse
			err := svc.GetShippingStatus(ctx, &shipping.ShippingStatusRequest{OrderId: 42}, &resp)
			So(err, ShouldNotBeNil)
			realError := errors.Parse(err.Error())
			So(realError, ShouldNotBeNil)
			So(realError.Detail, ShouldEqual, "Failed to query shipping status: Faily Fail")
			So(realError.Code, ShouldEqual, http.StatusInternalServerError)
		})

		Convey("requesting shipping status with a nil request should fail", func() {
			var resp shipping.ShippingStatusResponse
			err := svc.GetShippingStatus(ctx, nil, &resp)
			So(err, ShouldNotBeNil)
			realError := errors.Parse(err.Error())
			So(realError, ShouldNotBeNil)
			So(realError.Code, ShouldEqual, http.StatusBadRequest)
		})
	})
}

func TestShippingService_MarkItemShipped(t *testing.T) {
	Convey("Given a shipping service", t, func() {
		ctx := context.Background()
		repo := &fakeRepo{}
		pub := &fakePublisher{}
		svc := service.NewShippingService(repo, pub)

		Convey("marking an item as shipped should invoke repository", func() {
			repo.shouldFail = false
			pub.publishCount = 0
			var resp shipping.MarkShippedResponse
			err := svc.MarkItemShipped(ctx, &shipping.MarkShippedRequest{OrderId: 42, ShippingMethod: shipping.ShippingMethod_SM_UPS, Sku: "8675309"}, &resp)
			So(err, ShouldBeNil)
			So(resp.TrackingNumber, ShouldEqual, "111111")
			So(resp.Success, ShouldEqual, true)
			So(pub.publishCount, ShouldEqual, 1)
		})

		Convey("marking an item as shipped on non-existent order should fail", func() {
			repo.shouldFail = false
			var resp shipping.MarkShippedResponse
			err := svc.MarkItemShipped(ctx, &shipping.MarkShippedRequest{OrderId: 1, ShippingMethod: shipping.ShippingMethod_SM_UPS, Sku: "8675309"}, &resp)
			So(err, ShouldNotBeNil)
			realError := errors.Parse(err.Error())
			So(realError, ShouldNotBeNil)
			So(realError.Code, ShouldEqual, http.StatusNotFound)
		})

		Convey("marking an item as shipped should fail when repo fails", func() {
			repo.shouldFail = true
			var resp shipping.MarkShippedResponse
			err := svc.MarkItemShipped(ctx, &shipping.MarkShippedRequest{OrderId: 42, ShippingMethod: shipping.ShippingMethod_SM_UPS, Sku: "8675309"}, &resp)
			So(err, ShouldNotBeNil)
			realError := errors.Parse(err.Error())
			So(realError, ShouldNotBeNil)
			So(realError.Code, ShouldEqual, http.StatusInternalServerError)
		})

		Convey("marking an item as shipped should fail with a nil request", func() {
			var resp shipping.MarkShippedResponse
			err := svc.MarkItemShipped(ctx, nil, &resp)
			So(err, ShouldNotBeNil)
			realError := errors.Parse(err.Error())
			So(realError, ShouldNotBeNil)
			So(realError.Code, ShouldEqual, http.StatusBadRequest)
		})

		Convey("should get a bad request when we don't supply a valid shipping method", func() {
			repo.shouldFail = false
			var resp shipping.MarkShippedResponse
			err := svc.MarkItemShipped(ctx, &shipping.MarkShippedRequest{OrderId: 42, ShippingMethod: shipping.ShippingMethod_SM_UNKNOWN, Sku: "8675309"}, &resp)
			So(err, ShouldNotBeNil)
			realError := errors.Parse(err.Error())
			So(realError, ShouldNotBeNil)
			So(realError.Code, ShouldEqual, http.StatusBadRequest)
		})
	})
}

type fakeRepo struct {
	shouldFail bool
}

func (r *fakeRepo) GetShippingCosts(sku string, zipCode string) (costs []*shipping.ShippingCost, err error) {
	if r.shouldFail {
		return nil, stderrors.New("Faily Fail")
	}
	return []*shipping.ShippingCost{
		&shipping.ShippingCost{
			Price:  2500,
			Method: shipping.ShippingMethod_SM_FEDEX,
		},
		&shipping.ShippingCost{
			Price:  1000,
			Method: shipping.ShippingMethod_SM_RAVEN,
		},
	}, nil
}

func (r *fakeRepo) MarkShipped(sku string, orderID uint64, note string, shippingMethod shipping.ShippingMethod) (trackingNumber string, err error) {
	if r.shouldFail {
		return "", stderrors.New("Faily Fail")
	}
	return "111111", nil
}

func (r *fakeRepo) GetShippingStatus(orderID uint64, sku string) (shippingStatus *shipping.ShippingStatus, err error) {
	if r.shouldFail {
		return nil, stderrors.New("Faily Fail")
	}

	if orderID == 42 {
		return &shipping.ShippingStatus{
			TrackingNumber: "111111",
			ShippingMethod: shipping.ShippingMethod_SM_RAVEN,
		}, nil
	}
	return
}

func (r *fakeRepo) ProductExists(sku string) (exists bool, err error) {
	return sku == "8675309", nil
}

func (r *fakeRepo) OrderExists(orderID uint64) (exists bool, err error) {
	return orderID == 42, nil
}

type fakePublisher struct {
	shouldFail   bool
	publishCount int
}

func (p *fakePublisher) PublishItemShippedEvent(event *shipping.ItemShippedEvent) (err error) {
	if p.shouldFail {
		return stderrors.New("Faily Fail")
	}
	p.publishCount++
	return nil
}
