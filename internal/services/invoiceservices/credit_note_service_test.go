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

func TestCreditNoteHeaderService_GetCreditNoteHeaders(t *testing.T) {
	err := test.LoadSQL(logUser, dbService)
	if err != nil {
		t.Error(err)
		return
	}

	ctx := LoginUser()

	creditNoteHeaderService := NewCreditNoteHeaderService(log, dbService, redisService, userServiceClient)
	creditNoteHeaders := []*invoiceproto.CreditNoteHeader{}
	creditNoteHeader, err := GetCreditNoteHeader(uint32(1), []byte{88, 185, 127, 200, 184, 142, 72, 33, 146, 248, 145, 38, 92, 77, 194, 252}, "58b97fc8-b88e-4821-92f8-91265c4dc2fc", "2005-06-25T10:04:26Z", "2005-06-25T10:04:26Z", "2005-06-21T10:04:26Z", "sample", "GBP", "2005-06-12T10:04:26Z", "2005-06-13T10:04:26Z", "2005-06-14T10:04:26Z", "2005-06-15T10:04:26Z", "2005-06-16T10:04:26Z", "2005-06-17T10:04:26Z", float64(0), float64(0), float64(0), float64(107.5), "2005-06-23T10:04:26Z", "2005-06-23T10:04:26Z", "auth0|673c75d516e8adb9e6ffc892", "auth0|673c75d516e8adb9e6ffc892")
	if err != nil {
		t.Error(err)
		return
	}

	creditNoteHeader2, err := GetCreditNoteHeader(uint32(2), []byte{240, 33, 191, 102, 37, 200, 71, 166, 168, 247, 97, 102, 197, 81, 160, 23}, "f021bf66-25c8-47a6-a8f7-6166c551a017", "2009-12-15T10:04:26Z", "2009-12-15T10:04:26Z", "2009-12-15T10:04:26Z", "Ordered in our booth at the convention", "EUR", "2009-11-02T10:04:26Z", "2009-11-03T10:04:26Z", "2009-11-04T10:04:26Z", "2009-11-05T10:04:26Z", "2009-11-06T10:04:26Z", "2009-11-07T10:04:26Z", float64(100), float64(1000), float64(0.3), float64(729), "2009-11-01T10:04:26Z", "2009-11-01T10:04:26Z", "auth0|673c75d516e8adb9e6ffc892", "auth0|673c75d516e8adb9e6ffc892")
	if err != nil {
		t.Error(err)
		return
	}

	creditNoteHeader3, err := GetCreditNoteHeader(uint32(3), []byte{131, 155, 98, 222, 250, 229, 76, 183, 170, 88, 60, 60, 224, 188, 139, 9}, "839b62de-fae5-4cb7-aa58-3c3ce0bc8b09", "2011-06-01T10:04:26Z", "2011-06-01T10:04:26Z", "2011-06-01T10:04:26Z", "", "GBP", "2011-06-01T10:04:26Z", "2011-06-02T10:04:26Z", "2011-06-03T10:04:26Z", "2011-06-04T10:04:26Z", "2011-06-05T10:04:26Z", "2011-06-06T10:04:26Z", float64(0), float64(0), float64(0), float64(20), "2011-06-01T10:04:26Z", "2011-06-01T10:04:26Z", "auth0|673c75d516e8adb9e6ffc892", "auth0|673c75d516e8adb9e6ffc892")
	if err != nil {
		t.Error(err)
		return
	}

	creditNoteHeaders = append(creditNoteHeaders, creditNoteHeader3, creditNoteHeader2, creditNoteHeader)

	form := invoiceproto.GetCreditNoteHeadersRequest{}
	form.Limit = "8"
	form.NextCursor = ""
	form.UserEmail = "sprov300@gmail.com"
	form.RequestId = "bks1m1g91jau4nkks2f0"

	nextc := "MA=="
	creditNoteHeadersResponse := invoiceproto.GetCreditNoteHeadersResponse{CreditNoteHeaders: creditNoteHeaders, NextCursor: nextc}

	type args struct {
		ctx context.Context
		in  *invoiceproto.GetCreditNoteHeadersRequest
	}
	tests := []struct {
		cs      *CreditNoteHeaderService
		args    args
		want    *invoiceproto.GetCreditNoteHeadersResponse
		wantErr bool
	}{
		{
			cs: creditNoteHeaderService,
			args: args{
				ctx: ctx,
				in:  &form,
			},
			want:    &creditNoteHeadersResponse,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		creditNoteHeadersResp, err := tt.cs.GetCreditNoteHeaders(tt.args.ctx, tt.args.in)
		if (err != nil) != tt.wantErr {
			t.Errorf("CreditNoteHeaderService.GetCreditNoteHeaders() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(creditNoteHeadersResp, tt.want) {
			t.Errorf("CreditNoteHeaderService.GetCreditNoteHeaders() = %v, want %v", creditNoteHeadersResp, tt.want)
		}
		assert.NotNil(t, creditNoteHeadersResp)
		creditNoteHeaderResult := creditNoteHeadersResp.CreditNoteHeaders[2]
		assert.Equal(t, creditNoteHeaderResult.CreditNoteHeaderD.Note, "sample", "they should be equal")
		assert.Equal(t, creditNoteHeaderResult.CreditNoteHeaderD.TaxCurrencyCode, "GBP", "they should be equal")
		assert.Equal(t, creditNoteHeaderResult.CreditNoteHeaderD.PayableAmount, float64(107.5), "they should be equal")
	}
}

func TestCreditNoteHeaderService_GetCreditNoteHeader(t *testing.T) {
	err := test.LoadSQL(logUser, dbService)
	if err != nil {
		t.Error(err)
		return
	}

	ctx := LoginUser()

	creditNoteHeaderService := NewCreditNoteHeaderService(log, dbService, redisService, userServiceClient)

	creditNoteHeader2, err := GetCreditNoteHeader(uint32(2), []byte{240, 33, 191, 102, 37, 200, 71, 166, 168, 247, 97, 102, 197, 81, 160, 23}, "f021bf66-25c8-47a6-a8f7-6166c551a017", "2009-12-15T10:04:26Z", "2009-12-15T10:04:26Z", "2009-12-15T10:04:26Z", "Ordered in our booth at the convention", "EUR", "2009-11-02T10:04:26Z", "2009-11-03T10:04:26Z", "2009-11-04T10:04:26Z", "2009-11-05T10:04:26Z", "2009-11-06T10:04:26Z", "2009-11-07T10:04:26Z", float64(100), float64(1000), float64(0.3), float64(729), "2009-11-01T10:04:26Z", "2009-11-01T10:04:26Z", "auth0|673c75d516e8adb9e6ffc892", "auth0|673c75d516e8adb9e6ffc892")
	if err != nil {
		t.Error(err)
		return
	}

	creditNoteHeaderResponse := invoiceproto.GetCreditNoteHeaderResponse{}
	creditNoteHeaderResponse.CreditNoteHeader = creditNoteHeader2

	form := invoiceproto.GetCreditNoteHeaderRequest{}
	gform := commonproto.GetRequest{}
	gform.Id = "f021bf66-25c8-47a6-a8f7-6166c551a017"
	gform.UserEmail = "sprov300@gmail.com"
	gform.RequestId = "bks1m1g91jau4nkks2f0"
	form.GetRequest = &gform

	type args struct {
		ctx context.Context
		in  *invoiceproto.GetCreditNoteHeaderRequest
	}
	tests := []struct {
		cs      *CreditNoteHeaderService
		args    args
		want    *invoiceproto.GetCreditNoteHeaderResponse
		wantErr bool
	}{
		{
			cs: creditNoteHeaderService,
			args: args{
				ctx: ctx,
				in:  &form,
			},
			want:    &creditNoteHeaderResponse,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		creditNoteHeadersResp, err := tt.cs.GetCreditNoteHeader(tt.args.ctx, tt.args.in)
		if (err != nil) != tt.wantErr {
			t.Errorf("CreditNoteHeaderService.GetCreditNoteHeader() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(creditNoteHeadersResp, tt.want) {
			t.Errorf("CreditNoteHeaderService.GetCreditNoteHeader() = %v, want %v", creditNoteHeadersResp, tt.want)
		}
		assert.NotNil(t, creditNoteHeadersResp)
		creditNoteHeaderResult := creditNoteHeadersResp.CreditNoteHeader
		assert.Equal(t, creditNoteHeaderResult.CreditNoteHeaderD.Note, "Ordered in our booth at the convention", "they should be equal")
		assert.Equal(t, creditNoteHeaderResult.CreditNoteHeaderD.TaxCurrencyCode, "EUR", "they should be equal")
		assert.Equal(t, creditNoteHeaderResult.CreditNoteHeaderD.ChargeTotalAmount, float64(100), "they should be equal")
		assert.Equal(t, creditNoteHeaderResult.CreditNoteHeaderD.PrepaidAmount, float64(1000), "they should be equal")
		assert.Equal(t, creditNoteHeaderResult.CreditNoteHeaderD.PayableRoundingAmount, float64(0.3), "they should be equal")
		assert.Equal(t, creditNoteHeaderResult.CreditNoteHeaderD.PayableAmount, float64(729), "they should be equal")
	}
}

func TestCreditNoteHeaderService_GetCreditNoteHeaderByPk(t *testing.T) {
	err := test.LoadSQL(logUser, dbService)
	if err != nil {
		t.Error(err)
		return
	}

	ctx := LoginUser()

	creditNoteHeaderService := NewCreditNoteHeaderService(log, dbService, redisService, userServiceClient)

	creditNoteHeader2, err := GetCreditNoteHeader(uint32(2), []byte{240, 33, 191, 102, 37, 200, 71, 166, 168, 247, 97, 102, 197, 81, 160, 23}, "f021bf66-25c8-47a6-a8f7-6166c551a017", "2009-12-15T10:04:26Z", "2009-12-15T10:04:26Z", "2009-12-15T10:04:26Z", "Ordered in our booth at the convention", "EUR", "2009-11-02T10:04:26Z", "2009-11-03T10:04:26Z", "2009-11-04T10:04:26Z", "2009-11-05T10:04:26Z", "2009-11-06T10:04:26Z", "2009-11-07T10:04:26Z", float64(100), float64(1000), float64(0.3), float64(729), "2009-11-01T10:04:26Z", "2009-11-01T10:04:26Z", "auth0|673c75d516e8adb9e6ffc892", "auth0|673c75d516e8adb9e6ffc892")
	if err != nil {
		t.Error(err)
		return
	}

	creditNoteHeaderResponse := invoiceproto.GetCreditNoteHeaderByPkResponse{}
	creditNoteHeaderResponse.CreditNoteHeader = creditNoteHeader2

	form := invoiceproto.GetCreditNoteHeaderByPkRequest{}
	gform := commonproto.GetByIdRequest{}
	gform.Id = uint32(2)
	gform.UserEmail = "sprov300@gmail.com"
	gform.RequestId = "bks1m1g91jau4nkks2f0"

	form.GetByIdRequest = &gform

	type args struct {
		ctx context.Context
		in  *invoiceproto.GetCreditNoteHeaderByPkRequest
	}
	tests := []struct {
		cs      *CreditNoteHeaderService
		args    args
		want    *invoiceproto.GetCreditNoteHeaderByPkResponse
		wantErr bool
	}{
		{
			cs: creditNoteHeaderService,
			args: args{
				ctx: ctx,
				in:  &form,
			},
			want:    &creditNoteHeaderResponse,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		creditNoteHeadersResp, err := tt.cs.GetCreditNoteHeaderByPk(tt.args.ctx, tt.args.in)
		if (err != nil) != tt.wantErr {
			t.Errorf("CreditNoteHeaderService.GetCreditNoteHeaderByPk() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(creditNoteHeadersResp, tt.want) {
			t.Errorf("CreditNoteHeaderService.GetCreditNoteHeaderByPk() = %v, want %v", creditNoteHeadersResp, tt.want)
		}
		assert.NotNil(t, creditNoteHeadersResp)
		creditNoteHeaderResult := creditNoteHeadersResp.CreditNoteHeader
		assert.Equal(t, creditNoteHeaderResult.CreditNoteHeaderD.Note, "Ordered in our booth at the convention", "they should be equal")
		assert.Equal(t, creditNoteHeaderResult.CreditNoteHeaderD.TaxCurrencyCode, "EUR", "they should be equal")
		assert.Equal(t, creditNoteHeaderResult.CreditNoteHeaderD.ChargeTotalAmount, float64(100), "they should be equal")
		assert.Equal(t, creditNoteHeaderResult.CreditNoteHeaderD.PrepaidAmount, float64(1000), "they should be equal")
		assert.Equal(t, creditNoteHeaderResult.CreditNoteHeaderD.PayableRoundingAmount, float64(0.3), "they should be equal")
		assert.Equal(t, creditNoteHeaderResult.CreditNoteHeaderD.PayableAmount, float64(729), "they should be equal")
	}
}

func TestCreditNoteHeaderService_GetCreditNoteLines(t *testing.T) {
	err := test.LoadSQL(logUser, dbService)
	if err != nil {
		t.Error(err)
		return
	}

	ctx := LoginUser()

	creditNoteHeaderService := NewCreditNoteHeaderService(log, dbService, redisService, userServiceClient)
	creditNoteLines := []*invoiceproto.CreditNoteLine{}
	creditNoteLine, err := GetCreditNoteLine(uint32(1), []byte{246, 110, 121, 136, 32, 189, 68, 63, 131, 39, 244, 234, 209, 15, 27, 75}, "f66e7988-20bd-443f-8327-f4ead10f1b4b", "as agreed on phone, the invoice should have been cancelled earlier, apologies", float64(100), float64(100), "2005-06-21T10:04:26Z", uint32(1), "2005-06-12T10:04:26Z", "2005-06-13T10:04:26Z", "2005-06-18T10:04:26Z", "2005-06-19T10:04:26Z", uint32(1), "2019-07-23T10:04:26Z", "2019-07-23T10:04:26Z", "auth0|673c75d516e8adb9e6ffc892", "auth0|673c75d516e8adb9e6ffc892")
	if err != nil {
		t.Error(err)
	}

	creditNoteLines = append(creditNoteLines, creditNoteLine)

	creditNoteLinesResponse := invoiceproto.GetCreditNoteLinesResponse{}
	creditNoteLinesResponse.CreditNoteLines = creditNoteLines

	form := invoiceproto.GetCreditNoteLinesRequest{}
	gform := commonproto.GetRequest{}
	gform.Id = "58b97fc8-b88e-4821-92f8-91265c4dc2fc"
	gform.UserEmail = "sprov300@gmail.com"
	gform.RequestId = "bks1m1g91jau4nkks2f0"
	form.GetRequest = &gform

	type args struct {
		ctx context.Context
		in  *invoiceproto.GetCreditNoteLinesRequest
	}
	tests := []struct {
		cs      *CreditNoteHeaderService
		args    args
		want    *invoiceproto.GetCreditNoteLinesResponse
		wantErr bool
	}{
		{
			cs: creditNoteHeaderService,
			args: args{
				ctx: ctx,
				in:  &form,
			},
			want:    &creditNoteLinesResponse,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		creditNoteLinesResp, err := tt.cs.GetCreditNoteLines(tt.args.ctx, tt.args.in)
		if (err != nil) != tt.wantErr {
			t.Errorf("CreditNoteHeaderService.GetCreditNoteLines() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(creditNoteLinesResp, tt.want) {
			t.Errorf("CreditNoteHeaderService.GetCreditNoteLines() = %v, want %v", creditNoteLinesResp, tt.want)
		}
		assert.NotNil(t, creditNoteLinesResp)
		creditNoteLineResult := creditNoteLinesResp.CreditNoteLines[0]
		assert.Equal(t, creditNoteLineResult.CreditNoteLineD.Note, "as agreed on phone, the invoice should have been cancelled earlier, apologies", "they should be equal")
		assert.Equal(t, creditNoteLineResult.CreditNoteLineD.CreditedQuantity, float64(100), "they should be equal")
		assert.Equal(t, creditNoteLineResult.CreditNoteLineD.LineExtensionAmount, float64(100), "they should be equal")
	}
}

func TestCreditNoteHeaderService_CreateCreditNoteHeader(t *testing.T) {
	err := test.LoadSQL(logUser, dbService)
	if err != nil {
		t.Error(err)
		return
	}
	ctx := LoginUser()

	creditNoteHeaderService := NewCreditNoteHeaderService(log, dbService, redisService, userServiceClient)
	creditNoteHeader := invoiceproto.CreateCreditNoteHeaderRequest{}
	creditNoteHeader.IssueDate = "06/01/2022"
	creditNoteHeader.DueDate = "08/21/2022"
	creditNoteHeader.TaxPointDate = "08/01/2022"
	creditNoteHeader.Note = "Ordered in our booth at the convention"
	creditNoteHeader.DocumentCurrencyCode = ""
	creditNoteHeader.TaxCurrencyCode = "GBP"
	creditNoteHeader.InvoicePeriodStartDate = "08/01/2022"
	creditNoteHeader.InvoicePeriodEndDate = "09/01/2022"
	creditNoteHeader.TaxExDate = "08/10/2022"
	creditNoteHeader.PricingExDate = "08/10/2022"
	creditNoteHeader.PaymentExDate = "08/11/2022"
	creditNoteHeader.PaymentAltExDate = "08/10/2022"
	creditNoteHeader.ChargeTotalAmount = float64(0)
	creditNoteHeader.PrepaidAmount = float64(0)
	creditNoteHeader.PayableRoundingAmount = float64(0)
	creditNoteHeader.PayableAmount = float64(20)
	creditNoteHeader.UserId = "auth0|673c75d516e8adb9e6ffc892"
	creditNoteHeader.UserEmail = "sprov300@gmail.com"
	creditNoteHeader.RequestId = "bks1m1g91jau4nkks2f0"

	creditNoteLine := invoiceproto.CreateCreditNoteLineRequest{}
	creditNoteLine.Note = "Ordered in our booth at the convention"
	creditNoteLine.CreditedQuantity = float64(2)
	creditNoteLine.LineExtensionAmount = float64(20)
	creditNoteLine.TaxPointDate = "08/01/2022"
	creditNoteLine.InvoicePeriodStartDate = "08/01/2022"
	creditNoteLine.InvoicePeriodEndDate = "09/01/2022"
	creditNoteLine.ItemId = uint32(7)
	creditNoteLine.PriceValidityPeriodStartDate = "09/01/2022"
	creditNoteLine.PriceValidityPeriodEndDate = "09/02/2022"
	creditNoteLine.UserId = "auth0|673c75d516e8adb9e6ffc892"
	creditNoteLine.UserEmail = "sprov300@gmail.com"
	creditNoteLine.RequestId = "bks1m1g91jau4nkks2f0"

	creditNoteLines := []*invoiceproto.CreateCreditNoteLineRequest{}
	creditNoteLines = append(creditNoteLines, &creditNoteLine)
	creditNoteHeader.CreditNoteLines = creditNoteLines

	type args struct {
		ctx context.Context
		in  *invoiceproto.CreateCreditNoteHeaderRequest
	}
	tests := []struct {
		cs      *CreditNoteHeaderService
		args    args
		wantErr bool
	}{
		{
			cs: creditNoteHeaderService,
			args: args{
				ctx: ctx,
				in:  &creditNoteHeader,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		creditNoteHeadersResp, err := tt.cs.CreateCreditNoteHeader(tt.args.ctx, tt.args.in)
		if (err != nil) != tt.wantErr {
			t.Errorf("CreditNoteHeaderService.CreateCreditNoteHeader() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		assert.NotNil(t, creditNoteHeadersResp)
		creditNoteHeaderResult := creditNoteHeadersResp.CreditNoteHeader
		assert.Equal(t, creditNoteHeaderResult.CreditNoteHeaderD.Note, "Ordered in our booth at the convention", "they should be equal")
		assert.Equal(t, creditNoteHeaderResult.CreditNoteHeaderD.TaxCurrencyCode, "GBP", "they should be equal")
		assert.Equal(t, creditNoteHeaderResult.CreditNoteHeaderD.PayableAmount, float64(20), "they should be equal")
	}
}

func GetCreditNoteHeader(id uint32, uuid4 []byte, idS string, issueDate string, dueDate string, taxPointDate string, note string, taxCurrencyCode string, invoicePeriodStartDate string, invoicePeriodEndDate string, taxExDate string, pricingExDate string, paymentExDate string, paymentAltExDate string, chargeTotalAmount float64, prepaidAmount float64, payableRoundingAmount float64, payableAmount float64, createdAt string, updatedAt string, createdByUserId string, updatedByUserId string) (*invoiceproto.CreditNoteHeader, error) {
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

	paymentAltExDate1, err := common.ConvertTimeToTimestamp(Layout, paymentAltExDate)
	if err != nil {
		log.Error("Error", zap.Error(err))
	}

	creditNoteHeaderD := invoiceproto.CreditNoteHeaderD{}
	creditNoteHeaderD.Id = id
	creditNoteHeaderD.Uuid4 = uuid4
	creditNoteHeaderD.IdS = idS
	creditNoteHeaderD.Note = note
	creditNoteHeaderD.TaxCurrencyCode = taxCurrencyCode
	creditNoteHeaderD.ChargeTotalAmount = chargeTotalAmount
	creditNoteHeaderD.PrepaidAmount = prepaidAmount
	creditNoteHeaderD.PayableRoundingAmount = payableRoundingAmount
	creditNoteHeaderD.PayableAmount = payableAmount

	creditNoteHeaderT := invoiceproto.CreditNoteHeaderT{}
	creditNoteHeaderT.IssueDate = issueDate1
	creditNoteHeaderT.DueDate = dueDate1
	creditNoteHeaderT.TaxPointDate = taxPointDate1
	creditNoteHeaderT.InvoicePeriodStartDate = invoicePeriodStartDate1
	creditNoteHeaderT.InvoicePeriodEndDate = invoicePeriodEndDate1
	creditNoteHeaderT.TaxExDate = taxExDate1
	creditNoteHeaderT.PricingExDate = pricingExDate1
	creditNoteHeaderT.PaymentExDate = paymentExDate1
	creditNoteHeaderT.PaymentAltExDate = paymentAltExDate1

	crUpdTime := new(commonproto.CrUpdTime)
	crUpdTime.CreatedAt = createdAt1
	crUpdTime.UpdatedAt = updatedAt1

	crUpdUser := new(commonproto.CrUpdUser)
	crUpdUser.StatusCode = "active"
	crUpdUser.CreatedByUserId = createdByUserId
	crUpdUser.UpdatedByUserId = updatedByUserId

	creditNoteHeader := invoiceproto.CreditNoteHeader{CreditNoteHeaderD: &creditNoteHeaderD, CreditNoteHeaderT: &creditNoteHeaderT, CrUpdUser: crUpdUser, CrUpdTime: crUpdTime}

	return &creditNoteHeader, nil
}

func GetCreditNoteLine(id uint32, uuid4 []byte, idS string, note string, creditedQuantity float64, lineExtensionAmount float64, taxPointDate string, itemId uint32, invoicePeriodStartDate string, invoicePeriodEndDate string, priceValidityPeriodStartDate string, priceValidityPeriodEndDate string, creditNoteHeaderId uint32, createdAt string, updatedAt string, createdByUserId string, updatedByUserId string) (*invoiceproto.CreditNoteLine, error) {
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

	creditNoteLineD := invoiceproto.CreditNoteLineD{}
	creditNoteLineD.Id = id
	creditNoteLineD.Uuid4 = uuid4
	creditNoteLineD.IdS = idS
	creditNoteLineD.Note = note
	creditNoteLineD.CreditedQuantity = creditedQuantity
	creditNoteLineD.LineExtensionAmount = lineExtensionAmount
	creditNoteLineD.ItemId = itemId
	creditNoteLineD.CreditNoteHeaderId = uint32(1)

	creditNoteLineT := invoiceproto.CreditNoteLineT{}
	creditNoteLineT.TaxPointDate = taxPointDate1
	creditNoteLineT.InvoicePeriodStartDate = invoicePeriodStartDate1
	creditNoteLineT.InvoicePeriodEndDate = invoicePeriodEndDate1
	creditNoteLineT.PriceValidityPeriodStartDate = priceValidityPeriodStartDate1
	creditNoteLineT.PriceValidityPeriodEndDate = priceValidityPeriodEndDate1

	crUpdTime := new(commonproto.CrUpdTime)
	crUpdTime.CreatedAt = createdAt1
	crUpdTime.UpdatedAt = updatedAt1

	crUpdUser := new(commonproto.CrUpdUser)
	crUpdUser.StatusCode = "active"
	crUpdUser.CreatedByUserId = createdByUserId
	crUpdUser.UpdatedByUserId = updatedByUserId

	creditNoteLine := invoiceproto.CreditNoteLine{CreditNoteLineD: &creditNoteLineD, CreditNoteLineT: &creditNoteLineT, CrUpdUser: crUpdUser, CrUpdTime: crUpdTime}

	return &creditNoteLine, nil
}

func TestCreditNoteHeaderService_UpdateCreditNoteHeader(t *testing.T) {
	err := test.LoadSQL(logUser, dbService)
	if err != nil {
		t.Error(err)
		return
	}

	ctx := LoginUser()

	creditNoteHeaderService := NewCreditNoteHeaderService(log, dbService, redisService, userServiceClient)

	form := invoiceproto.UpdateCreditNoteHeaderRequest{}
	form.Note = "Ordered"
	form.TaxCurrencyCode = "GBP"
	form.ChargeTotalAmount = float64(10)
	form.PrepaidAmount = float64(20)
	form.PayableRoundingAmount = float64(100)
	form.PayableAmount = float64(200)
	form.Id = "58b97fc8-b88e-4821-92f8-91265c4dc2fc"
	form.UserId = "auth0|673c75d516e8adb9e6ffc892"
	form.UserEmail = "sprov300@gmail.com"
	form.RequestId = "bks1m1g91jau4nkks2f0"

	updateResponse := invoiceproto.UpdateCreditNoteHeaderResponse{}

	type args struct {
		ctx context.Context
		in  *invoiceproto.UpdateCreditNoteHeaderRequest
	}
	tests := []struct {
		ds      *CreditNoteHeaderService
		args    args
		want    *invoiceproto.UpdateCreditNoteHeaderResponse
		wantErr bool
	}{
		{
			ds: creditNoteHeaderService,
			args: args{
				ctx: ctx,
				in:  &form,
			},
			want:    &updateResponse,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		got, err := tt.ds.UpdateCreditNoteHeader(tt.args.ctx, tt.args.in)
		if (err != nil) != tt.wantErr {
			t.Errorf("CreditNoteHeaderService.UpdateCreditNoteHeader() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("CreditNoteHeaderService.UpdateCreditNoteHeader() = %v, want %v", got, tt.want)
		}
	}
}
