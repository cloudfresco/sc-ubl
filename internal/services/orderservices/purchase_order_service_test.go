package orderservices

import (
	"context"
	"reflect"
	"testing"

	"github.com/cloudfresco/sc-ubl/internal/common"
	commonproto "github.com/cloudfresco/sc-ubl/internal/protogen/common/v1"
	orderproto "github.com/cloudfresco/sc-ubl/internal/protogen/order/v1"
	"github.com/cloudfresco/sc-ubl/test"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestPurchaseOrderHeaderService_GetPurchaseOrderHeaders(t *testing.T) {
	err := test.LoadSQL(logUser, dbService)
	if err != nil {
		t.Error(err)
		return
	}

	ctx := LoginUser()

	purchaseOrderService := NewPurchaseOrderHeaderService(log, dbService, redisService, userServiceClient)
	purchaseOrders := []*orderproto.PurchaseOrderHeader{}

	purchaseOrderHeader, err := GetPurchaseOrderHeader(uint32(1), []byte{65, 58, 64, 181, 95, 123, 64, 197, 187, 175, 214, 224, 37, 84, 63, 222}, "413a40b5-5f7b-40c5-bbaf-d6e025543fde", "2005-06-20T10:04:26Z", "2005-06-21T10:04:26Z", uint32(2), uint32(1), uint32(2), float64(100), float64(100), "2010-12-02T10:04:26Z", "2010-12-03T10:04:26Z", "2010-12-04T10:04:26Z", "2019-07-23T10:04:26Z", "2019-07-23T10:04:26Z", "auth0|673c75d516e8adb9e6ffc892", "auth0|673c75d516e8adb9e6ffc892")
	if err != nil {
		t.Error(err)
		return
	}

	purchaseOrderHeader2, err := GetPurchaseOrderHeader(uint32(2), []byte{139, 244, 179, 212, 149, 255, 72, 147, 172, 12, 131, 140, 90, 27, 25, 89}, "8bf4b3d4-95ff-4893-ac0c-838c5a1b1959", "2010-01-20T10:04:26Z", "2010-01-21T10:04:26Z", uint32(2), uint32(1), uint32(2), float64(6225), float64(6225), "2010-12-02T10:04:26Z", "2010-12-03T10:04:26Z", "2010-12-04T10:04:26Z", "2019-07-23T10:04:26Z", "2019-07-23T10:04:26Z", "auth0|673c75d516e8adb9e6ffc892", "auth0|673c75d516e8adb9e6ffc892")
	if err != nil {
		t.Error(err)
		return
	}

	purchaseOrders = append(purchaseOrders, purchaseOrderHeader2, purchaseOrderHeader)

	form := orderproto.GetPurchaseOrderHeadersRequest{}
	form.Limit = "8"
	form.NextCursor = ""
	form.UserEmail = "sprov300@gmail.com"
	form.RequestId = "bks1m1g91jau4nkks2f0"

	nextc := "MA=="
	purchaseOrderHeadersResponse := orderproto.GetPurchaseOrderHeadersResponse{PurchaseOrderHeaders: purchaseOrders, NextCursor: nextc}

	type args struct {
		ctx context.Context
		in  *orderproto.GetPurchaseOrderHeadersRequest
	}
	tests := []struct {
		ps      *PurchaseOrderHeaderService
		args    args
		want    *orderproto.GetPurchaseOrderHeadersResponse
		wantErr bool
	}{
		{
			ps: purchaseOrderService,
			args: args{
				ctx: ctx,
				in:  &form,
			},
			want:    &purchaseOrderHeadersResponse,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		purchaseOrderHeadersResp, err := tt.ps.GetPurchaseOrderHeaders(tt.args.ctx, tt.args.in)
		if (err != nil) != tt.wantErr {
			t.Errorf("PurchaseOrderHeaderService.GetPurchaseOrderHeaders() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(purchaseOrderHeadersResp, tt.want) {
			t.Errorf("PurchaseOrderHeaderService.GetPurchaseOrderHeaders() = %v, want %v", purchaseOrderHeadersResp, tt.want)
		}
		assert.NotNil(t, purchaseOrderHeadersResp)
		purchaseOrderHeaderResult := purchaseOrderHeadersResp.PurchaseOrderHeaders[0]
		assert.Equal(t, purchaseOrderHeaderResult.PurchaseOrderHeaderD.SellerSupplierPartyId, uint32(1), "they should be equal")
		assert.Equal(t, purchaseOrderHeaderResult.PurchaseOrderHeaderD.OriginatorCustomerPartyId, uint32(2), "they should be equal")
		assert.Equal(t, purchaseOrderHeaderResult.PurchaseOrderHeaderD.LineExtensionAmount, float64(6225), "they should be equal")
		assert.Equal(t, purchaseOrderHeaderResult.PurchaseOrderHeaderD.PayableAmount, float64(6225), "they should be equal")
	}
}

func TestPurchaseOrderHeaderService_GetPurchaseOrderHeader(t *testing.T) {
	err := test.LoadSQL(logUser, dbService)
	if err != nil {
		t.Error(err)
		return
	}

	ctx := LoginUser()

	purchaseOrderService := NewPurchaseOrderHeaderService(log, dbService, redisService, userServiceClient)

	purchaseOrderHeader, err := GetPurchaseOrderHeader(uint32(1), []byte{65, 58, 64, 181, 95, 123, 64, 197, 187, 175, 214, 224, 37, 84, 63, 222}, "413a40b5-5f7b-40c5-bbaf-d6e025543fde", "2005-06-20T10:04:26Z", "2005-06-21T10:04:26Z", uint32(2), uint32(1), uint32(2), float64(100), float64(100), "2010-12-02T10:04:26Z", "2010-12-03T10:04:26Z", "2010-12-04T10:04:26Z", "2019-07-23T10:04:26Z", "2019-07-23T10:04:26Z", "auth0|673c75d516e8adb9e6ffc892", "auth0|673c75d516e8adb9e6ffc892")
	if err != nil {
		t.Error(err)
		return
	}
	purchaseOrderHeaderResponse := orderproto.GetPurchaseOrderHeaderResponse{}
	purchaseOrderHeaderResponse.PurchaseOrderHeader = purchaseOrderHeader

	form := orderproto.GetPurchaseOrderHeaderRequest{}
	gform := commonproto.GetRequest{}
	gform.Id = "413a40b5-5f7b-40c5-bbaf-d6e025543fde"
	gform.UserEmail = "sprov300@gmail.com"
	gform.RequestId = "bks1m1g91jau4nkks2f0"
	form.GetRequest = &gform

	type args struct {
		ctx context.Context
		in  *orderproto.GetPurchaseOrderHeaderRequest
	}
	tests := []struct {
		ps      *PurchaseOrderHeaderService
		args    args
		want    *orderproto.GetPurchaseOrderHeaderResponse
		wantErr bool
	}{
		{
			ps: purchaseOrderService,
			args: args{
				ctx: ctx,
				in:  &form,
			},
			want:    &purchaseOrderHeaderResponse,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		purchaseOrderHeaderResp, err := tt.ps.GetPurchaseOrderHeader(tt.args.ctx, tt.args.in)
		if (err != nil) != tt.wantErr {
			t.Errorf("PurchaseOrderHeaderService.GetPurchaseOrderHeader() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(purchaseOrderHeaderResp, tt.want) {
			t.Errorf("PurchaseOrderHeaderService.GetPurchaseOrderHeader() = %v, want %v", purchaseOrderHeaderResp, tt.want)
		}
		assert.NotNil(t, purchaseOrderHeaderResp)
		purchaseOrderHeaderResult := purchaseOrderHeaderResp.PurchaseOrderHeader
		assert.Equal(t, purchaseOrderHeaderResult.PurchaseOrderHeaderD.SellerSupplierPartyId, uint32(1), "they should be equal")
		assert.Equal(t, purchaseOrderHeaderResult.PurchaseOrderHeaderD.OriginatorCustomerPartyId, uint32(2), "they should be equal")
		assert.Equal(t, purchaseOrderHeaderResult.PurchaseOrderHeaderD.LineExtensionAmount, float64(100), "they should be equal")
		assert.Equal(t, purchaseOrderHeaderResult.PurchaseOrderHeaderD.PayableAmount, float64(100), "they should be equal")
	}
}

func TestPurchaseOrderHeaderService_GetPurchaseOrderHeaderByPk(t *testing.T) {
	err := test.LoadSQL(logUser, dbService)
	if err != nil {
		t.Error(err)
		return
	}

	ctx := LoginUser()

	purchaseOrderService := NewPurchaseOrderHeaderService(log, dbService, redisService, userServiceClient)
	purchaseOrderHeader, err := GetPurchaseOrderHeader(uint32(1), []byte{65, 58, 64, 181, 95, 123, 64, 197, 187, 175, 214, 224, 37, 84, 63, 222}, "413a40b5-5f7b-40c5-bbaf-d6e025543fde", "2005-06-20T10:04:26Z", "2005-06-21T10:04:26Z", uint32(2), uint32(1), uint32(2), float64(100), float64(100), "2010-12-02T10:04:26Z", "2010-12-03T10:04:26Z", "2010-12-04T10:04:26Z", "2019-07-23T10:04:26Z", "2019-07-23T10:04:26Z", "auth0|673c75d516e8adb9e6ffc892", "auth0|673c75d516e8adb9e6ffc892")
	if err != nil {
		t.Error(err)
		return
	}

	purchaseOrderHeaderResponse := orderproto.GetPurchaseOrderHeaderByPkResponse{}
	purchaseOrderHeaderResponse.PurchaseOrderHeader = purchaseOrderHeader

	form := orderproto.GetPurchaseOrderHeaderByPkRequest{}

	gform := commonproto.GetByIdRequest{}
	gform.Id = uint32(1)
	gform.UserEmail = "sprov300@gmail.com"
	gform.RequestId = "bks1m1g91jau4nkks2f0"
	form.GetByIdRequest = &gform

	type args struct {
		ctx context.Context
		in  *orderproto.GetPurchaseOrderHeaderByPkRequest
	}
	tests := []struct {
		ps      *PurchaseOrderHeaderService
		args    args
		want    *orderproto.GetPurchaseOrderHeaderByPkResponse
		wantErr bool
	}{
		{
			ps: purchaseOrderService,
			args: args{
				ctx: ctx,
				in:  &form,
			},
			want:    &purchaseOrderHeaderResponse,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		purchaseOrderHeaderResp, err := tt.ps.GetPurchaseOrderHeaderByPk(tt.args.ctx, tt.args.in)
		if (err != nil) != tt.wantErr {
			t.Errorf("PurchaseOrderHeaderService.GetPurchaseOrderHeaderByPk() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(purchaseOrderHeaderResp, tt.want) {
			t.Errorf("PurchaseOrderHeaderService.GetPurchaseOrderHeaderByPk() = %v, want %v", purchaseOrderHeaderResp, tt.want)
		}
		assert.NotNil(t, purchaseOrderHeaderResp)
		purchaseOrderHeaderResult := purchaseOrderHeaderResp.PurchaseOrderHeader
		assert.Equal(t, purchaseOrderHeaderResult.PurchaseOrderHeaderD.SellerSupplierPartyId, uint32(1), "they should be equal")
		assert.Equal(t, purchaseOrderHeaderResult.PurchaseOrderHeaderD.OriginatorCustomerPartyId, uint32(2), "they should be equal")
		assert.Equal(t, purchaseOrderHeaderResult.PurchaseOrderHeaderD.LineExtensionAmount, float64(100), "they should be equal")
		assert.Equal(t, purchaseOrderHeaderResult.PurchaseOrderHeaderD.PayableAmount, float64(100), "they should be equal")
	}
}

func TestPurchaseOrderHeaderService_GetPurchaseOrderLines(t *testing.T) {
	err := test.LoadSQL(logUser, dbService)
	if err != nil {
		t.Error(err)
		return
	}

	ctx := LoginUser()

	purchaseOrderService := NewPurchaseOrderHeaderService(log, dbService, redisService, userServiceClient)
	purchaseLines := []*orderproto.PurchaseOrderLine{}

	purchaseLine1, err := GetPurchaseOrderLine(uint32(3), []byte{83, 33, 47, 15, 236, 76, 74, 8, 128, 194, 1, 217, 197, 110, 254, 19}, "53212f0f-ec4c-4a08-80c2-01d9c56efe13", float64(15), float64(225), float64(10), uint32(2), uint32(3), float64(15), float64(1), "2009-12-06T10:04:26Z", "2009-12-07T10:04:26Z", uint32(2), "2019-07-23T10:04:26Z", "2019-07-23T10:04:26Z", "auth0|673c75d516e8adb9e6ffc892", "auth0|673c75d516e8adb9e6ffc892")
	if err != nil {
		t.Error(err)
		return
	}
	purchaseLine2, err := GetPurchaseOrderLine(uint32(2), []byte{71, 103, 92, 161, 252, 5, 69, 81, 136, 203, 228, 82, 5, 98, 102, 148}, "47675ca1-fc05-4551-88cb-e45205626694", float64(120), float64(6000), float64(10), uint32(2), uint32(2), float64(50), float64(1), "2009-12-06T10:04:26Z", "2009-12-07T10:04:26Z", uint32(2), "2019-07-23T10:04:26Z", "2019-07-23T10:04:26Z", "auth0|673c75d516e8adb9e6ffc892", "auth0|673c75d516e8adb9e6ffc892")
	if err != nil {
		t.Error(err)
		return
	}

	purchaseLines = append(purchaseLines, purchaseLine2, purchaseLine1)

	purchaseOrderLinesResponse := orderproto.GetPurchaseOrderLinesResponse{}
	purchaseOrderLinesResponse.PurchaseOrderLines = purchaseLines

	form := orderproto.GetPurchaseOrderLinesRequest{}
	gform := commonproto.GetRequest{}
	gform.Id = "8bf4b3d4-95ff-4893-ac0c-838c5a1b1959"
	gform.UserEmail = "sprov300@gmail.com"
	gform.RequestId = "bks1m1g91jau4nkks2f0"
	form.GetRequest = &gform

	type args struct {
		ctx context.Context
		in  *orderproto.GetPurchaseOrderLinesRequest
	}
	tests := []struct {
		ps      *PurchaseOrderHeaderService
		args    args
		want    *orderproto.GetPurchaseOrderLinesResponse
		wantErr bool
	}{
		{
			ps: purchaseOrderService,
			args: args{
				ctx: ctx,
				in:  &form,
			},
			want:    &purchaseOrderLinesResponse,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		purchaseOrderLinesResp, err := tt.ps.GetPurchaseOrderLines(tt.args.ctx, tt.args.in)
		if (err != nil) != tt.wantErr {
			t.Errorf("PurchaseOrderHeaderService.GetPurchaseOrderLines() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(purchaseOrderLinesResp, tt.want) {
			t.Errorf("PurchaseOrderHeaderService.GetPurchaseOrderLines() = %v, want %v", purchaseOrderLinesResp, tt.want)
		}
		assert.NotNil(t, purchaseOrderLinesResp)
		purchaseOrderLineResult := purchaseOrderLinesResp.PurchaseOrderLines[0]
		assert.Equal(t, purchaseOrderLineResult.PurchaseOrderLineD.Quantity, float64(120), "they should be equal")
		assert.Equal(t, purchaseOrderLineResult.PurchaseOrderLineD.LineExtensionAmount, float64(6000), "they should be equal")
		assert.Equal(t, purchaseOrderLineResult.PurchaseOrderLineD.TotalTaxAmount, float64(10), "they should be equal")
		assert.Equal(t, purchaseOrderLineResult.PurchaseOrderLineD.PriceAmount, float64(50), "they should be equal")
		assert.Equal(t, purchaseOrderLineResult.PurchaseOrderLineD.PriceBaseQuantity, float64(1), "they should be equal")
	}
}

func TestPurchaseOrderHeaderService_CreatePurchaseOrderHeader(t *testing.T) {
	err := test.LoadSQL(logUser, dbService)
	if err != nil {
		t.Error(err)
		return
	}

	ctx := LoginUser()

	purchaseOrderService := NewPurchaseOrderHeaderService(log, dbService, redisService, userServiceClient)

	purchaseOrderHeader := orderproto.CreatePurchaseOrderHeaderRequest{}
	purchaseOrderHeader.IssueDate = "02/12/2022"
	purchaseOrderHeader.ValidityPeriod = "02/22/2022"
	purchaseOrderHeader.OrderTypeCode = "ABCFES"
	purchaseOrderHeader.Note = "Sample"
	purchaseOrderHeader.LineExtensionAmount = float64(100)
	purchaseOrderHeader.PayableAmount = float64(100)
	purchaseOrderHeader.TaxExDate = "08/10/2022"
	purchaseOrderHeader.PricingExDate = "08/10/2022"
	purchaseOrderHeader.PaymentExDate = "08/11/2022"
	purchaseOrderHeader.UserId = "auth0|673c75d516e8adb9e6ffc892"
	purchaseOrderHeader.UserEmail = "sprov300@gmail.com"
	purchaseOrderHeader.RequestId = "bks1m1g91jau4nkks2f0"

	purchaseOrderLine := orderproto.CreatePurchaseOrderLineRequest{}
	purchaseOrderLine.Note = "Mrs Green agreed to waive charge"
	purchaseOrderLine.LineStatusCode = "ABCFES"
	purchaseOrderLine.OriginatorPartyId = uint32(1)
	purchaseOrderLine.Quantity = float64(100)
	purchaseOrderLine.LineExtensionAmount = float64(100)
	purchaseOrderLine.TotalTaxAmount = float64(17.5)
	purchaseOrderLine.ItemId = uint32(7)
	purchaseOrderLine.PriceAmount = float64(17.5)
	purchaseOrderLine.PriceBaseQuantity = float64(100)
	purchaseOrderLine.PriceValidityPeriodStartDate = "02/22/2022"
	purchaseOrderLine.PriceValidityPeriodEndDate = "02/27/2022"
	purchaseOrderLine.UserId = "auth0|673c75d516e8adb9e6ffc892"
	purchaseOrderLine.UserEmail = "sprov300@gmail.com"
	purchaseOrderLine.RequestId = "bks1m1g91jau4nkks2f0"

	purchaseOrderLines := []*orderproto.CreatePurchaseOrderLineRequest{}
	purchaseOrderLines = append(purchaseOrderLines, &purchaseOrderLine)
	purchaseOrderHeader.PurchaseOrderLines = purchaseOrderLines

	type args struct {
		ctx context.Context
		in  *orderproto.CreatePurchaseOrderHeaderRequest
	}
	tests := []struct {
		ps      *PurchaseOrderHeaderService
		args    args
		wantErr bool
	}{
		{
			ps: purchaseOrderService,
			args: args{
				ctx: ctx,
				in:  &purchaseOrderHeader,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		purchaseOrderHeaderResp, err := tt.ps.CreatePurchaseOrderHeader(tt.args.ctx, tt.args.in)
		if (err != nil) != tt.wantErr {
			t.Errorf("PurchaseOrderHeaderService.CreatePurchaseOrderHeader() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		assert.NotNil(t, purchaseOrderHeaderResp)
		purchaseOrderHeaderResult := purchaseOrderHeaderResp.PurchaseOrderHeader
		assert.Equal(t, purchaseOrderHeaderResult.PurchaseOrderHeaderD.OrderTypeCode, "ABCFES", "they should be equal")
		assert.Equal(t, purchaseOrderHeaderResult.PurchaseOrderHeaderD.LineExtensionAmount, float64(100), "they should be equal")
		assert.Equal(t, purchaseOrderHeaderResult.PurchaseOrderHeaderD.PayableAmount, float64(100), "they should be equal")
	}
}

func GetPurchaseOrderHeader(id uint32, uuid4 []byte, idS string, issueDate string, validityPeriod string, buyerCustomerPartyId uint32, sellerSupplierPartyId uint32, originatorCustomerPartyId uint32, lineExtensionAmount float64, payableAmount float64, taxExDate string, pricingExDate string, paymentExDate string, createdAt string, updatedAt string, createdByUserId string, updatedByUserId string) (*orderproto.PurchaseOrderHeader, error) {
	createdAt1, err := common.ConvertTimeToTimestamp(Layout, createdAt)
	if err != nil {
		log.Error("Error", zap.Error(err))
	}

	updatedAt1, err := common.ConvertTimeToTimestamp(Layout, updatedAt)
	if err != nil {
		log.Error("Error", zap.Error(err))
	}

	issueDate1, err := common.ConvertTimeToTimestamp(Layout, issueDate)
	if err != nil {
		log.Error("Error", zap.Error(err))
	}

	validityPeriod1, err := common.ConvertTimeToTimestamp(Layout, validityPeriod)
	if err != nil {
		log.Error("Error", zap.Error(err))
	}

	taxExDate1, err := common.ConvertTimeToTimestamp(Layout, taxExDate)
	if err != nil {
		log.Error("Error", zap.Error(err))
	}

	pricingExDate1, err := common.ConvertTimeToTimestamp(Layout, pricingExDate)
	if err != nil {
		log.Error("Error", zap.Error(err))
	}

	paymentExDate1, err := common.ConvertTimeToTimestamp(Layout, paymentExDate)
	if err != nil {
		log.Error("Error", zap.Error(err))
	}

	purchaseOrderHeaderD := orderproto.PurchaseOrderHeaderD{}
	purchaseOrderHeaderD.Id = id
	purchaseOrderHeaderD.Uuid4 = uuid4
	purchaseOrderHeaderD.IdS = idS
	purchaseOrderHeaderD.BuyerCustomerPartyId = buyerCustomerPartyId
	purchaseOrderHeaderD.SellerSupplierPartyId = sellerSupplierPartyId
	purchaseOrderHeaderD.OriginatorCustomerPartyId = originatorCustomerPartyId
	purchaseOrderHeaderD.LineExtensionAmount = lineExtensionAmount
	purchaseOrderHeaderD.PayableAmount = payableAmount

	purchaseOrderHeaderT := orderproto.PurchaseOrderHeaderT{}
	purchaseOrderHeaderT.IssueDate = issueDate1
	purchaseOrderHeaderT.ValidityPeriod = validityPeriod1
	purchaseOrderHeaderT.TaxExDate = taxExDate1
	purchaseOrderHeaderT.PricingExDate = pricingExDate1
	purchaseOrderHeaderT.PaymentExDate = paymentExDate1

	crUpdTime := commonproto.CrUpdTime{}
	crUpdTime.CreatedAt = createdAt1
	crUpdTime.UpdatedAt = updatedAt1

	crUpdUser := commonproto.CrUpdUser{}
	crUpdUser.StatusCode = "active"
	crUpdUser.CreatedByUserId = createdByUserId
	crUpdUser.UpdatedByUserId = updatedByUserId

	purchaseOrderHeader := orderproto.PurchaseOrderHeader{PurchaseOrderHeaderD: &purchaseOrderHeaderD, PurchaseOrderHeaderT: &purchaseOrderHeaderT, CrUpdUser: &crUpdUser, CrUpdTime: &crUpdTime}

	return &purchaseOrderHeader, nil
}

func GetPurchaseOrderLine(id uint32, uuid4 []byte, idS string, quantity float64, lineExtensionAmount float64, totalTaxAmount float64, originatorPartyId uint32, itemId uint32, priceAmount float64, priceBaseQuantity float64, priceValidityPeriodStartDate string, priceValidityPeriodEndDate string, purchaseOrderHeaderId uint32, createdAt string, updatedAt string, createdByUserId string, updatedByUserId string) (*orderproto.PurchaseOrderLine, error) {
	createdAt1, err := common.ConvertTimeToTimestamp(Layout, createdAt)
	if err != nil {
		log.Error("Error", zap.Error(err))
	}

	updatedAt1, err := common.ConvertTimeToTimestamp(Layout, updatedAt)
	if err != nil {
		log.Error("Error", zap.Error(err))
	}

	priceValidityPeriodStartDate1, err := common.ConvertTimeToTimestamp(Layout, priceValidityPeriodStartDate)
	if err != nil {
		log.Error("Error", zap.Error(err))
	}

	priceValidityPeriodEndDate1, err := common.ConvertTimeToTimestamp(Layout, priceValidityPeriodEndDate)
	if err != nil {
		log.Error("Error", zap.Error(err))
	}

	purchaseOrderLineD := orderproto.PurchaseOrderLineD{}
	purchaseOrderLineD.Id = id
	purchaseOrderLineD.Uuid4 = uuid4
	purchaseOrderLineD.IdS = idS
	purchaseOrderLineD.Quantity = quantity
	purchaseOrderLineD.LineExtensionAmount = lineExtensionAmount
	purchaseOrderLineD.TotalTaxAmount = totalTaxAmount
	purchaseOrderLineD.OriginatorPartyId = originatorPartyId
	purchaseOrderLineD.ItemId = itemId
	purchaseOrderLineD.PriceAmount = priceAmount
	purchaseOrderLineD.PriceBaseQuantity = priceBaseQuantity
	purchaseOrderLineD.PurchaseOrderHeaderId = purchaseOrderHeaderId

	purchaseOrderLineT := orderproto.PurchaseOrderLineT{}
	purchaseOrderLineT.PriceValidityPeriodStartDate = priceValidityPeriodStartDate1
	purchaseOrderLineT.PriceValidityPeriodEndDate = priceValidityPeriodEndDate1

	crUpdTime := commonproto.CrUpdTime{}
	crUpdTime.CreatedAt = createdAt1
	crUpdTime.UpdatedAt = updatedAt1

	crUpdUser := commonproto.CrUpdUser{}
	crUpdUser.StatusCode = "active"
	crUpdUser.CreatedByUserId = createdByUserId
	crUpdUser.UpdatedByUserId = updatedByUserId

	purchaseOrderLine := orderproto.PurchaseOrderLine{PurchaseOrderLineD: &purchaseOrderLineD, PurchaseOrderLineT: &purchaseOrderLineT, CrUpdUser: &crUpdUser, CrUpdTime: &crUpdTime}

	return &purchaseOrderLine, nil
}

func TestPurchaseOrderHeaderService_UpdatePurchaseOrderHeader(t *testing.T) {
	err := test.LoadSQL(logUser, dbService)
	if err != nil {
		t.Error(err)
		return
	}

	ctx := LoginUser()

	purchaseOrderService := NewPurchaseOrderHeaderService(log, dbService, redisService, userServiceClient)

	form := orderproto.UpdatePurchaseOrderHeaderRequest{}
	form.OrderTypeCode = "BCFESD"
	form.Note = "sample2"
	form.RequestedInvoiceCurrencyCode = "EUR"
	form.DocumentCurrencyCode = "EUR"
	form.PricingCurrencyCode = "EUR"
	form.TaxCurrencyCode = "EUR"
	form.AccountingCostCode = ""
	form.AccountingCost = "BookingCode001"
	form.Id = "413a40b5-5f7b-40c5-bbaf-d6e025543fde"
	form.UserId = "auth0|673c75d516e8adb9e6ffc892"
	form.UserEmail = "sprov300@gmail.com"
	form.RequestId = "bks1m1g91jau4nkks2f0"

	updateResponse := orderproto.UpdatePurchaseOrderHeaderResponse{}

	type args struct {
		ctx context.Context
		in  *orderproto.UpdatePurchaseOrderHeaderRequest
	}
	tests := []struct {
		ps      *PurchaseOrderHeaderService
		args    args
		want    *orderproto.UpdatePurchaseOrderHeaderResponse
		wantErr bool
	}{
		{
			ps: purchaseOrderService,
			args: args{
				ctx: ctx,
				in:  &form,
			},
			want:    &updateResponse,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		got, err := tt.ps.UpdatePurchaseOrderHeader(tt.args.ctx, tt.args.in)
		if (err != nil) != tt.wantErr {
			t.Errorf("PurchaseOrderHeaderService.UpdatePurchaseOrderHeader() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("PurchaseOrderHeaderService.UpdatePurchaseOrderHeader() = %v, want %v", got, tt.want)
		}
	}
}
