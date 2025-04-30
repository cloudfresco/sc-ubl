package invoiceservices

import (
	"context"
	"reflect"
	"testing"

	"github.com/cloudfresco/sc-ubl/internal/common"
	commonproto "github.com/cloudfresco/sc-ubl/internal/protogen/common/v1"
	invoiceproto "github.com/cloudfresco/sc-ubl/internal/protogen/invoice/v1"
	"github.com/cloudfresco/sc-ubl/test"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestInvoiceService_GetInvoices(t *testing.T) {
	err := test.LoadSQL(logUser, dbService)
	if err != nil {
		t.Error(err)
		return
	}

	ctx := LoginUser()

	invoiceService := NewInvoiceService(log, dbService, redisService, userServiceClient)
	invoiceHeaders := []*invoiceproto.InvoiceHeader{}

	invoiceHeader, err := GetInvoiceHeader(uint32(1), []byte{131, 77, 213, 134, 161, 184, 78, 63, 137, 207, 191, 210, 94, 221, 154, 96}, "834dd586-a1b8-4e3f-89cf-bfd25edd9a60", "2005-06-21T10:04:26Z", "2005-06-21T10:04:26Z", "2005-06-21T10:04:26Z", "2005-06-26T10:04:26Z", "sample", "SalesInvoice", "2005-06-12T10:04:26Z", "2005-06-25T10:04:26Z", "2005-06-27T10:04:26Z", "2005-06-28T10:04:26Z", "2005-06-29T10:04:26Z",
		"2019-07-23T10:04:26Z", "2019-07-23T10:04:26Z", "auth0|673c75d516e8adb9e6ffc892", "auth0|673c75d516e8adb9e6ffc892")
	if err != nil {
		t.Error(err)
		return
	}

	invoiceHeader2, err := GetInvoiceHeader(uint32(2), []byte{44, 173, 88, 224, 76, 40, 67, 214, 175, 109, 66, 242, 35, 55, 60, 89}, "2cad58e0-4c28-43d6-af6d-42f223373c59", "2009-12-15T10:04:26Z", "2009-12-15T10:04:26Z", "2009-11-30T10:04:26Z", "2009-12-02T10:04:26Z", "Ordered in our booth at the convention.", "", "2009-11-02T10:04:26Z", "2009-12-01T10:04:26Z", "2009-12-03T10:04:26Z", "2009-12-04T10:04:26Z", "2009-12-05T10:04:26Z", "2019-07-23T10:04:26Z", "2019-07-23T10:04:26Z", "auth0|673c75d516e8adb9e6ffc892", "auth0|673c75d516e8adb9e6ffc892")
	if err != nil {
		t.Error(err)
		return
	}

	invoiceHeaders = append(invoiceHeaders, invoiceHeader2, invoiceHeader)

	form := invoiceproto.GetInvoicesRequest{}
	form.Limit = "8"
	form.NextCursor = ""
	form.UserEmail = "sprov300@gmail.com"
	form.RequestId = "bks1m1g91jau4nkks2f0"

	nextc := "MA=="
	invoiceHeadersResponse := invoiceproto.GetInvoicesResponse{InvoiceHeaders: invoiceHeaders, NextCursor: nextc}

	type args struct {
		ctx context.Context
		in  *invoiceproto.GetInvoicesRequest
	}
	tests := []struct {
		is      *InvoiceService
		args    args
		want    *invoiceproto.GetInvoicesResponse
		wantErr bool
	}{
		{
			is: invoiceService,
			args: args{
				ctx: ctx,
				in:  &form,
			},
			want:    &invoiceHeadersResponse,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		invoiceHeadersResp, err := tt.is.GetInvoices(tt.args.ctx, tt.args.in)
		if (err != nil) != tt.wantErr {
			t.Errorf("InvoiceService.GetInvoices() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(invoiceHeadersResp, tt.want) {
			t.Errorf("InvoiceService.GetInvoices() = %v, want %v", invoiceHeadersResp, tt.want)
		}
		assert.NotNil(t, invoiceHeadersResp)
		invoiceHeaderResult := invoiceHeadersResp.InvoiceHeaders[0]
		assert.Equal(t, invoiceHeaderResult.InvoiceHeaderD.Note, "Ordered in our booth at the convention.", "they should be equal")
		assert.Equal(t, invoiceHeaderResult.InvoiceHeaderD.InvoiceTypeCode, "", "they should be equal")
	}
}

func TestInvoiceService_GetInvoice(t *testing.T) {
	err := test.LoadSQL(logUser, dbService)
	if err != nil {
		t.Error(err)
		return
	}

	ctx := LoginUser()

	invoiceService := NewInvoiceService(log, dbService, redisService, userServiceClient)

	invoiceHeader, err := GetInvoiceHeader(uint32(1), []byte{131, 77, 213, 134, 161, 184, 78, 63, 137, 207, 191, 210, 94, 221, 154, 96}, "834dd586-a1b8-4e3f-89cf-bfd25edd9a60", "2005-06-21T10:04:26Z", "2005-06-21T10:04:26Z", "2005-06-21T10:04:26Z", "2005-06-26T10:04:26Z", "sample", "SalesInvoice", "2005-06-12T10:04:26Z", "2005-06-25T10:04:26Z", "2005-06-27T10:04:26Z", "2005-06-28T10:04:26Z", "2005-06-29T10:04:26Z",
		"2019-07-23T10:04:26Z", "2019-07-23T10:04:26Z", "auth0|673c75d516e8adb9e6ffc892", "auth0|673c75d516e8adb9e6ffc892")
	if err != nil {
		t.Error(err)
		return
	}
	invoiceHeaderResponse := invoiceproto.GetInvoiceResponse{}
	invoiceHeaderResponse.InvoiceHeader = invoiceHeader

	form := invoiceproto.GetInvoiceRequest{}
	gform := commonproto.GetRequest{}
	gform.Id = "834dd586-a1b8-4e3f-89cf-bfd25edd9a60"
	gform.UserEmail = "sprov300@gmail.com"
	gform.RequestId = "bks1m1g91jau4nkks2f0"
	form.GetRequest = &gform

	type args struct {
		ctx context.Context
		in  *invoiceproto.GetInvoiceRequest
	}
	tests := []struct {
		is      *InvoiceService
		args    args
		want    *invoiceproto.GetInvoiceResponse
		wantErr bool
	}{
		{
			is: invoiceService,
			args: args{
				ctx: ctx,
				in:  &form,
			},
			want:    &invoiceHeaderResponse,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		invoiceHeaderResp, err := tt.is.GetInvoice(tt.args.ctx, tt.args.in)
		if (err != nil) != tt.wantErr {
			t.Errorf("InvoiceService.GetInvoice() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(invoiceHeaderResp, tt.want) {
			t.Errorf("InvoiceService.GetInvoice() = %v, want %v", invoiceHeaderResp, tt.want)
		}
		assert.NotNil(t, invoiceHeaderResp)
		invoiceHeaderResult := invoiceHeaderResp.InvoiceHeader
		assert.Equal(t, invoiceHeaderResult.InvoiceHeaderD.Note, "sample", "they should be equal")
		assert.Equal(t, invoiceHeaderResult.InvoiceHeaderD.InvoiceTypeCode, "SalesInvoice", "they should be equal")
	}
}

func TestInvoiceService_GetInvoiceByPk(t *testing.T) {
	err := test.LoadSQL(logUser, dbService)
	if err != nil {
		t.Error(err)
		return
	}

	ctx := LoginUser()

	invoiceService := NewInvoiceService(log, dbService, redisService, userServiceClient)

	invoiceHeader, err := GetInvoiceHeader(uint32(1), []byte{131, 77, 213, 134, 161, 184, 78, 63, 137, 207, 191, 210, 94, 221, 154, 96}, "834dd586-a1b8-4e3f-89cf-bfd25edd9a60", "2005-06-21T10:04:26Z", "2005-06-21T10:04:26Z", "2005-06-21T10:04:26Z", "2005-06-26T10:04:26Z", "sample", "SalesInvoice", "2005-06-12T10:04:26Z", "2005-06-25T10:04:26Z", "2005-06-27T10:04:26Z", "2005-06-28T10:04:26Z", "2005-06-29T10:04:26Z",
		"2019-07-23T10:04:26Z", "2019-07-23T10:04:26Z", "auth0|673c75d516e8adb9e6ffc892", "auth0|673c75d516e8adb9e6ffc892")
	if err != nil {
		t.Error(err)
		return
	}

	invoiceHeaderResponse := invoiceproto.GetInvoiceByPkResponse{}
	invoiceHeaderResponse.InvoiceHeader = invoiceHeader

	form := invoiceproto.GetInvoiceByPkRequest{}
	gform := commonproto.GetByIdRequest{}
	gform.Id = uint32(1)
	gform.UserEmail = "sprov300@gmail.com"
	gform.RequestId = "bks1m1g91jau4nkks2f0"
	form.GetByIdRequest = &gform

	type args struct {
		ctx context.Context
		in  *invoiceproto.GetInvoiceByPkRequest
	}
	tests := []struct {
		is      *InvoiceService
		args    args
		want    *invoiceproto.GetInvoiceByPkResponse
		wantErr bool
	}{
		{
			is: invoiceService,
			args: args{
				ctx: ctx,
				in:  &form,
			},
			want:    &invoiceHeaderResponse,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		invoiceHeaderResp, err := tt.is.GetInvoiceByPk(tt.args.ctx, tt.args.in)
		if (err != nil) != tt.wantErr {
			t.Errorf("InvoiceService.GetInvoiceByPk() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(invoiceHeaderResp, tt.want) {
			t.Errorf("InvoiceService.GetInvoiceByPk() = %v, want %v", invoiceHeaderResp, tt.want)
		}
		assert.NotNil(t, invoiceHeaderResp)
		invoiceHeaderResult := invoiceHeaderResp.InvoiceHeader
		assert.Equal(t, invoiceHeaderResult.InvoiceHeaderD.Note, "sample", "they should be equal")
		assert.Equal(t, invoiceHeaderResult.InvoiceHeaderD.InvoiceTypeCode, "SalesInvoice", "they should be equal")
	}
}

func TestInvoiceService_GetInvoiceLines(t *testing.T) {
	err := test.LoadSQL(logUser, dbService)
	if err != nil {
		t.Error(err)
		return
	}

	ctx := LoginUser()
	invoiceService := NewInvoiceService(log, dbService, redisService, userServiceClient)
	invoiceLine, err := GetInvoiceLine(uint32(1), []byte{97, 193, 98, 66, 182, 31, 71, 158, 185, 55, 14, 147, 71, 225, 142, 69}, "61c16242-b61f-479e-b937-0e9347e18e45", float64(100), float64(100), uint32(2), "2005-06-20T10:04:26Z", "2005-06-12T10:04:26Z", "2005-06-25T10:04:26Z", "2005-06-30T10:04:26Z", "2005-07-01T10:04:26Z", uint32(1), "2019-07-23T10:04:26Z", "2019-07-23T10:04:26Z", "auth0|673c75d516e8adb9e6ffc892", "auth0|673c75d516e8adb9e6ffc892")
	if err != nil {
		t.Error(err)
	}
	invoiceLines := []*invoiceproto.InvoiceLine{}
	invoiceLines = append(invoiceLines, invoiceLine)

	invoiceLinesResponse := invoiceproto.GetInvoiceLinesResponse{}
	invoiceLinesResponse.InvoiceLines = invoiceLines

	form := invoiceproto.GetInvoiceLinesRequest{}
	gform := commonproto.GetRequest{}
	gform.Id = "834dd586-a1b8-4e3f-89cf-bfd25edd9a60"
	gform.UserEmail = "sprov300@gmail.com"
	gform.RequestId = "bks1m1g91jau4nkks2f0"
	form.GetRequest = &gform

	type args struct {
		ctx context.Context
		in  *invoiceproto.GetInvoiceLinesRequest
	}
	tests := []struct {
		is      *InvoiceService
		args    args
		want    *invoiceproto.GetInvoiceLinesResponse
		wantErr bool
	}{
		{
			is: invoiceService,
			args: args{
				ctx: ctx,
				in:  &form,
			},
			want:    &invoiceLinesResponse,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		invoiceLinesResp, err := tt.is.GetInvoiceLines(tt.args.ctx, tt.args.in)
		if (err != nil) != tt.wantErr {
			t.Errorf("InvoiceService.GetInvoiceLines() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(invoiceLinesResp, tt.want) {
			t.Errorf("InvoiceService.GetInvoiceLines() = %v, want %v", invoiceLinesResp, tt.want)
		}
		assert.NotNil(t, invoiceLinesResp)
		invoiceLineResult := invoiceLinesResp.InvoiceLines[0]
		assert.Equal(t, invoiceLineResult.InvoiceLineD.InvoicedQuantity, float64(100), "they should be equal")
		assert.Equal(t, invoiceLineResult.InvoiceLineD.LineExtensionAmount, float64(100), "they should be equal")
	}
}

func TestInvoiceService_CreateInvoice(t *testing.T) {
	err := test.LoadSQL(logUser, dbService)
	if err != nil {
		t.Error(err)
		return
	}

	ctx := LoginUser()

	invoiceService := NewInvoiceService(log, dbService, redisService, userServiceClient)

	invoiceHeader := invoiceproto.CreateInvoiceRequest{}
	invoiceHeader.IssueDate = "06/21/2022"
	invoiceHeader.DueDate = "06/25/2022"
	invoiceHeader.TaxPointDate = "07/25/2022"
	invoiceHeader.InvoiceTypeCode = "SalesInvoice"
	invoiceHeader.Note = "sample"
	invoiceHeader.InvoicePeriodStartDate = "06/21/2022"
	invoiceHeader.InvoicePeriodEndDate = "07/01/2022"
	invoiceHeader.TaxExDate = "06/21/2022"
	invoiceHeader.PricingExDate = "06/22/2022"
	invoiceHeader.PaymentExDate = "06/23/2022"
	invoiceHeader.PaymentAltExDate = "06/24/2022"
	invoiceHeader.ChargeTotalAmount = float64(0)
	invoiceHeader.PrepaidAmount = float64(0)
	invoiceHeader.PayableRoundingAmount = float64(0)
	invoiceHeader.PayableAmount = float64(200)
	invoiceHeader.UserId = "auth0|673c75d516e8adb9e6ffc892"
	invoiceHeader.UserEmail = "sprov300@gmail.com"
	invoiceHeader.RequestId = "bks1m1g91jau4nkks2f0"

	invoiceLine := invoiceproto.CreateInvoiceLineRequest{}
	invoiceLine.Note = "Ordered in our booth at the convention."
	invoiceLine.InvoicedQuantity = float64(1)
	invoiceLine.LineExtensionAmount = float64(200)
	invoiceLine.TaxPointDate = "11/25/2022"
	invoiceLine.AccountingCost = "Code002"
	invoiceLine.ItemId = uint32(7)
	invoiceLine.InvoicePeriodStartDate = "06/21/2022"
	invoiceLine.InvoicePeriodEndDate = "07/01/2022"
	invoiceLine.PriceValidityPeriodStartDate = "07/01/2022"
	invoiceLine.PriceValidityPeriodEndDate = "08/01/2022"
	invoiceLine.UserId = "auth0|673c75d516e8adb9e6ffc892"
	invoiceLine.UserEmail = "sprov300@gmail.com"
	invoiceLine.RequestId = "bks1m1g91jau4nkks2f0"

	invoiceLines := []*invoiceproto.CreateInvoiceLineRequest{}
	invoiceLines = append(invoiceLines, &invoiceLine)
	invoiceHeader.InvoiceLines = invoiceLines

	type args struct {
		ctx context.Context
		in  *invoiceproto.CreateInvoiceRequest
	}
	tests := []struct {
		is      *InvoiceService
		args    args
		wantErr bool
	}{
		{
			is: invoiceService,
			args: args{
				ctx: ctx,
				in:  &invoiceHeader,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		invoiceHeaderResp, err := tt.is.CreateInvoice(tt.args.ctx, tt.args.in)
		if (err != nil) != tt.wantErr {
			t.Errorf("InvoiceService.CreateInvoice() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		assert.NotNil(t, invoiceHeaderResp)
		invoiceHeaderResult := invoiceHeaderResp.InvoiceHeader
		assert.Equal(t, invoiceHeaderResult.InvoiceHeaderD.Note, "sample", "they should be equal")
		assert.Equal(t, invoiceHeaderResult.InvoiceHeaderD.InvoiceTypeCode, "SalesInvoice", "they should be equal")
	}
}

func GetInvoiceHeader(id uint32, uuid4 []byte, idS string, issueDate string, dueDate string, taxPointDate string, taxExDate string, note string, invoiceTypeCode string, invoicePeriodStartDate string, invoicePeriodEndDate string, pricingExDate string, paymentExDate string, paymentAltExDate string, createdAt string, updatedAt string, createdByUserId string, updatedByUserId string) (*invoiceproto.InvoiceHeader, error) {
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

	dueDate1, err := common.ConvertTimeToTimestamp(Layout, dueDate)
	if err != nil {
		log.Error("Error", zap.Error(err))
	}

	taxExDate1, err := common.ConvertTimeToTimestamp(Layout, taxExDate)
	if err != nil {
		log.Error("Error", zap.Error(err))
	}

	taxPointDate1, err := common.ConvertTimeToTimestamp(Layout, taxPointDate)
	if err != nil {
		log.Error("Error", zap.Error(err))
	}

	invoicePeriodStartDate1, err := common.ConvertTimeToTimestamp(Layout, invoicePeriodStartDate)
	if err != nil {
		log.Error("Error", zap.Error(err))
	}

	invoicePeriodEndDate1, err := common.ConvertTimeToTimestamp(Layout, invoicePeriodEndDate)
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

	paymentAltExDate1, err := common.ConvertTimeToTimestamp(Layout, paymentAltExDate)
	if err != nil {
		log.Error("Error", zap.Error(err))
	}

	invoiceHeaderD := invoiceproto.InvoiceHeaderD{}
	invoiceHeaderD.Id = id
	invoiceHeaderD.Uuid4 = uuid4
	invoiceHeaderD.IdS = idS
	invoiceHeaderD.Note = note
	invoiceHeaderD.InvoiceTypeCode = invoiceTypeCode

	invoiceHeaderT := invoiceproto.InvoiceHeaderT{}
	invoiceHeaderT.IssueDate = issueDate1
	invoiceHeaderT.DueDate = dueDate1
	invoiceHeaderT.TaxPointDate = taxPointDate1
	invoiceHeaderT.TaxExDate = taxExDate1
	invoiceHeaderT.InvoicePeriodStartDate = invoicePeriodStartDate1
	invoiceHeaderT.InvoicePeriodEndDate = invoicePeriodEndDate1
	invoiceHeaderT.PricingExDate = pricingExDate1
	invoiceHeaderT.PaymentExDate = paymentExDate1
	invoiceHeaderT.PaymentAltExDate = paymentAltExDate1

	crUpdTime := new(commonproto.CrUpdTime)
	crUpdTime.CreatedAt = createdAt1
	crUpdTime.UpdatedAt = updatedAt1

	crUpdUser := new(commonproto.CrUpdUser)
	crUpdUser.StatusCode = "active"
	crUpdUser.CreatedByUserId = createdByUserId
	crUpdUser.UpdatedByUserId = updatedByUserId

	invoiceHeader := invoiceproto.InvoiceHeader{InvoiceHeaderD: &invoiceHeaderD, InvoiceHeaderT: &invoiceHeaderT, CrUpdUser: crUpdUser, CrUpdTime: crUpdTime}

	return &invoiceHeader, nil
}

func GetInvoiceLine(id uint32, uuid4 []byte, idS string, invoicedQuantity float64, lineExtensionAmount float64, itemId uint32, taxPointDate string, invoicePeriodStartDate string, invoicePeriodEndDate string, priceValidityPeriodStartDate string, priceValidityPeriodEndDate string, invoiceHeaderId uint32, createdAt string, updatedAt string, createdByUserId string, updatedByUserId string) (*invoiceproto.InvoiceLine, error) {
	createdAt1, err := common.ConvertTimeToTimestamp(Layout, createdAt)
	if err != nil {
		log.Error("Error", zap.Error(err))
	}

	updatedAt1, err := common.ConvertTimeToTimestamp(Layout, updatedAt)
	if err != nil {
		log.Error("Error", zap.Error(err))
	}

	taxPointDate1, err := common.ConvertTimeToTimestamp(Layout, taxPointDate)
	if err != nil {
		log.Error("Error", zap.Error(err))
	}

	invoicePeriodStartDate1, err := common.ConvertTimeToTimestamp(Layout, invoicePeriodStartDate)
	if err != nil {
		log.Error("Error", zap.Error(err))
	}

	invoicePeriodEndDate1, err := common.ConvertTimeToTimestamp(Layout, invoicePeriodEndDate)
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

	invoiceLineD := invoiceproto.InvoiceLineD{}
	invoiceLineD.Id = id
	invoiceLineD.Uuid4 = uuid4
	invoiceLineD.IdS = idS
	invoiceLineD.InvoicedQuantity = invoicedQuantity
	invoiceLineD.LineExtensionAmount = lineExtensionAmount
	invoiceLineD.ItemId = itemId
	invoiceLineD.InvoiceHeaderId = invoiceHeaderId

	invoiceLineT := invoiceproto.InvoiceLineT{}
	invoiceLineT.TaxPointDate = taxPointDate1
	invoiceLineT.InvoicePeriodStartDate = invoicePeriodStartDate1
	invoiceLineT.InvoicePeriodEndDate = invoicePeriodEndDate1
	invoiceLineT.PriceValidityPeriodStartDate = priceValidityPeriodStartDate1
	invoiceLineT.PriceValidityPeriodEndDate = priceValidityPeriodEndDate1

	crUpdTime := new(commonproto.CrUpdTime)
	crUpdTime.CreatedAt = createdAt1
	crUpdTime.UpdatedAt = updatedAt1

	crUpdUser := new(commonproto.CrUpdUser)
	crUpdUser.StatusCode = "active"
	crUpdUser.CreatedByUserId = createdByUserId
	crUpdUser.UpdatedByUserId = updatedByUserId

	invoiceLine := invoiceproto.InvoiceLine{InvoiceLineD: &invoiceLineD, InvoiceLineT: &invoiceLineT, CrUpdUser: crUpdUser, CrUpdTime: crUpdTime}

	return &invoiceLine, nil
}

func TestInvoiceService_UpdateInvoice(t *testing.T) {
	err := test.LoadSQL(logUser, dbService)
	if err != nil {
		t.Error(err)
		return
	}

	ctx := LoginUser()

	invoiceService := NewInvoiceService(log, dbService, redisService, userServiceClient)

	form := invoiceproto.UpdateInvoiceRequest{}
	form.Note = "Ordered"
	form.InvoiceTypeCode = "Sales"
	form.ChargeTotalAmount = float64(200)
	form.PrepaidAmount = float64(100)
	form.PayableRoundingAmount = float64(150)
	form.PayableAmount = float64(400)
	form.Id = "834dd586-a1b8-4e3f-89cf-bfd25edd9a60"
	form.UserId = "auth0|673c75d516e8adb9e6ffc892"
	form.UserEmail = "sprov300@gmail.com"
	form.RequestId = "bks1m1g91jau4nkks2f0"

	updateResponse := invoiceproto.UpdateInvoiceResponse{}

	type args struct {
		ctx context.Context
		in  *invoiceproto.UpdateInvoiceRequest
	}
	tests := []struct {
		is      *InvoiceService
		args    args
		want    *invoiceproto.UpdateInvoiceResponse
		wantErr bool
	}{
		{
			is: invoiceService,
			args: args{
				ctx: ctx,
				in:  &form,
			},
			want:    &updateResponse,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		got, err := tt.is.UpdateInvoice(tt.args.ctx, tt.args.in)
		if (err != nil) != tt.wantErr {
			t.Errorf("InvoiceService.UpdateInvoice() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("InvoiceService.UpdateInvoice() = %v, want %v", got, tt.want)
		}
	}
}
