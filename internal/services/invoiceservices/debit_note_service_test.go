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

func TestDebitNoteHeaderService_GetDebitNoteHeaders(t *testing.T) {
	err := test.LoadSQL(logUser, dbService)
	if err != nil {
		t.Error(err)
		return
	}

	ctx := LoginUser()

	debitNoteHeader, err := GetDebitNoteHeader(uint32(1), []byte{147, 15, 136, 6, 219, 36, 69, 98, 184, 199, 114, 223, 117, 81, 131, 85}, "930f8806-db24-4562-b8c7-72df75518355", "2005-12-15T10:04:26Z", ">Ordered in our booth at the convention", "2019-11-30T10:04:26Z", "EUR", "Project cost code 123", "2009-11-01T10:04:26Z", "2009-11-30T10:04:26Z", "2009-12-02T10:04:26Z", "2009-12-03T10:04:26Z", "2009-12-04T10:04:26Z", "2009-12-05T10:04:26Z", "2019-07-23T10:04:26Z", "2019-07-23T10:04:26Z", "auth0|673c75d516e8adb9e6ffc892", "auth0|673c75d516e8adb9e6ffc892")
	if err != nil {
		t.Error(err)
		return
	}

	debitNoteHeaderService := NewDebitNoteHeaderService(log, dbService, redisService, userServiceClient)
	debitNoteHeaders := []*invoiceproto.DebitNoteHeader{}

	debitNoteHeaders = append(debitNoteHeaders, debitNoteHeader)

	form := invoiceproto.GetDebitNoteHeadersRequest{}
	form.Limit = "8"
	form.NextCursor = ""
	form.UserEmail = "sprov300@gmail.com"
	form.RequestId = "bks1m1g91jau4nkks2f0"

	nextc := "MA=="
	debitNoteHeadersResponse := invoiceproto.GetDebitNoteHeadersResponse{DebitNoteHeaders: debitNoteHeaders, NextCursor: nextc}

	type args struct {
		ctx context.Context
		in  *invoiceproto.GetDebitNoteHeadersRequest
	}
	tests := []struct {
		ds      *DebitNoteHeaderService
		args    args
		want    *invoiceproto.GetDebitNoteHeadersResponse
		wantErr bool
	}{
		{
			ds: debitNoteHeaderService,
			args: args{
				ctx: ctx,
				in:  &form,
			},
			want:    &debitNoteHeadersResponse,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		debitNoteHeadersResp, err := tt.ds.GetDebitNoteHeaders(tt.args.ctx, tt.args.in)
		if (err != nil) != tt.wantErr {
			t.Errorf("DebitNoteHeaderService.GetDebitNoteHeaders() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(debitNoteHeadersResp, tt.want) {
			t.Errorf("DebitNoteHeaderService.GetDebitNoteHeaders() = %v, want %v", debitNoteHeadersResp, tt.want)
		}
		assert.NotNil(t, debitNoteHeadersResp)
		debitNoteHeaderResult := debitNoteHeadersResp.DebitNoteHeaders[0]
		assert.Equal(t, debitNoteHeaderResult.DebitNoteHeaderD.Note, ">Ordered in our booth at the convention", "they should be equal")
		assert.Equal(t, debitNoteHeaderResult.DebitNoteHeaderD.DocumentCurrencyCode, "EUR", "they should be equal")
		assert.Equal(t, debitNoteHeaderResult.DebitNoteHeaderD.AccountingCost, "Project cost code 123", "they should be equal")
	}
}

func TestDebitNoteHeaderService_GetDebitNoteHeader(t *testing.T) {
	err := test.LoadSQL(logUser, dbService)
	if err != nil {
		t.Error(err)
		return
	}

	ctx := LoginUser()

	debitNoteHeaderService := NewDebitNoteHeaderService(log, dbService, redisService, userServiceClient)

	debitNoteHeader, err := GetDebitNoteHeader(uint32(1), []byte{147, 15, 136, 6, 219, 36, 69, 98, 184, 199, 114, 223, 117, 81, 131, 85}, "930f8806-db24-4562-b8c7-72df75518355", "2005-12-15T10:04:26Z", ">Ordered in our booth at the convention", "2019-11-30T10:04:26Z", "EUR", "Project cost code 123", "2009-11-01T10:04:26Z", "2009-11-30T10:04:26Z", "2009-12-02T10:04:26Z", "2009-12-03T10:04:26Z", "2009-12-04T10:04:26Z", "2009-12-05T10:04:26Z", "2019-07-23T10:04:26Z", "2019-07-23T10:04:26Z", "auth0|673c75d516e8adb9e6ffc892", "auth0|673c75d516e8adb9e6ffc892")
	if err != nil {
		t.Error(err)
		return
	}

	debitNoteHeaderResponse := invoiceproto.GetDebitNoteHeaderResponse{}
	debitNoteHeaderResponse.DebitNoteHeader = debitNoteHeader

	form := invoiceproto.GetDebitNoteHeaderRequest{}
	gform := commonproto.GetRequest{}
	gform.Id = "930f8806-db24-4562-b8c7-72df75518355"
	gform.UserEmail = "sprov300@gmail.com"
	gform.RequestId = "bks1m1g91jau4nkks2f0"
	form.GetRequest = &gform

	type args struct {
		ctx context.Context
		in  *invoiceproto.GetDebitNoteHeaderRequest
	}
	tests := []struct {
		ds      *DebitNoteHeaderService
		args    args
		want    *invoiceproto.GetDebitNoteHeaderResponse
		wantErr bool
	}{
		{
			ds: debitNoteHeaderService,
			args: args{
				ctx: ctx,
				in:  &form,
			},
			want:    &debitNoteHeaderResponse,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		debitNoteHeaderResp, err := tt.ds.GetDebitNoteHeader(tt.args.ctx, tt.args.in)
		if (err != nil) != tt.wantErr {
			t.Errorf("DebitNoteHeaderService.GetDebitNoteHeader() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(debitNoteHeaderResp, tt.want) {
			t.Errorf("DebitNoteHeaderService.GetDebitNoteHeader() = %v, want %v", debitNoteHeaderResp, tt.want)
		}
		assert.NotNil(t, debitNoteHeaderResp)
		debitNoteHeaderResult := debitNoteHeaderResp.DebitNoteHeader
		assert.Equal(t, debitNoteHeaderResult.DebitNoteHeaderD.Note, ">Ordered in our booth at the convention", "they should be equal")
		assert.Equal(t, debitNoteHeaderResult.DebitNoteHeaderD.DocumentCurrencyCode, "EUR", "they should be equal")
		assert.Equal(t, debitNoteHeaderResult.DebitNoteHeaderD.AccountingCost, "Project cost code 123", "they should be equal")
	}
}

func TestDebitNoteHeaderService_GetDebitNoteHeaderByPk(t *testing.T) {
	err := test.LoadSQL(logUser, dbService)
	if err != nil {
		t.Error(err)
		return
	}

	ctx := LoginUser()

	debitNoteHeader, err := GetDebitNoteHeader(uint32(1), []byte{147, 15, 136, 6, 219, 36, 69, 98, 184, 199, 114, 223, 117, 81, 131, 85}, "930f8806-db24-4562-b8c7-72df75518355", "2005-12-15T10:04:26Z", ">Ordered in our booth at the convention", "2019-11-30T10:04:26Z", "EUR", "Project cost code 123", "2009-11-01T10:04:26Z", "2009-11-30T10:04:26Z", "2009-12-02T10:04:26Z", "2009-12-03T10:04:26Z", "2009-12-04T10:04:26Z", "2009-12-05T10:04:26Z", "2019-07-23T10:04:26Z", "2019-07-23T10:04:26Z", "auth0|673c75d516e8adb9e6ffc892", "auth0|673c75d516e8adb9e6ffc892")
	if err != nil {
		t.Error(err)
		return
	}

	debitNoteHeaderService := NewDebitNoteHeaderService(log, dbService, redisService, userServiceClient)

	debitNoteHeaderResponse := invoiceproto.GetDebitNoteHeaderByPkResponse{}
	debitNoteHeaderResponse.DebitNoteHeader = debitNoteHeader

	form := invoiceproto.GetDebitNoteHeaderByPkRequest{}
	gform := commonproto.GetByIdRequest{}
	gform.Id = uint32(1)
	gform.UserEmail = "sprov300@gmail.com"
	gform.RequestId = "bks1m1g91jau4nkks2f0"
	form.GetByIdRequest = &gform

	type args struct {
		ctx context.Context
		in  *invoiceproto.GetDebitNoteHeaderByPkRequest
	}
	tests := []struct {
		ds      *DebitNoteHeaderService
		args    args
		want    *invoiceproto.GetDebitNoteHeaderByPkResponse
		wantErr bool
	}{
		{
			ds: debitNoteHeaderService,
			args: args{
				ctx: ctx,
				in:  &form,
			},
			want:    &debitNoteHeaderResponse,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		debitNoteHeaderResp, err := tt.ds.GetDebitNoteHeaderByPk(tt.args.ctx, tt.args.in)
		if (err != nil) != tt.wantErr {
			t.Errorf("DebitNoteHeaderService.GetDebitNoteHeaderByPk() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(debitNoteHeaderResp, tt.want) {
			t.Errorf("DebitNoteHeaderService.GetDebitNoteHeaderByPk() = %v, want %v", debitNoteHeaderResp, tt.want)
		}
		assert.NotNil(t, debitNoteHeaderResp)
		debitNoteHeaderResult := debitNoteHeaderResp.DebitNoteHeader
		assert.Equal(t, debitNoteHeaderResult.DebitNoteHeaderD.Note, ">Ordered in our booth at the convention", "they should be equal")
		assert.Equal(t, debitNoteHeaderResult.DebitNoteHeaderD.DocumentCurrencyCode, "EUR", "they should be equal")
		assert.Equal(t, debitNoteHeaderResult.DebitNoteHeaderD.AccountingCost, "Project cost code 123", "they should be equal")
	}
}

func TestDebitNoteHeaderService_GetDebitNoteLines(t *testing.T) {
	err := test.LoadSQL(logUser, dbService)
	if err != nil {
		t.Error(err)
		return
	}

	ctx := LoginUser()

	debitNoteHeaderService := NewDebitNoteHeaderService(log, dbService, redisService, userServiceClient)
	debitNoteLines := []*invoiceproto.DebitNoteLine{}
	debitNoteLine, err := GetDebitNoteLine(uint32(1), []byte{120, 38, 247, 229, 52, 89, 75, 136, 181, 182, 33, 70, 213, 222, 113, 147}, "7826f7e5-3459-4b88-b5b6-2146d5de7193", "Scratch on box", float64(1), float64(1273), "2005-06-21T10:04:26Z", "BookingCode002", uint32(2), "2005-06-18T10:04:26Z", "2005-06-19T10:04:26Z", uint32(1), "2019-07-23T10:04:26Z", "2019-07-23T10:04:26Z", "auth0|673c75d516e8adb9e6ffc892", "auth0|673c75d516e8adb9e6ffc892")
	if err != nil {
		t.Error(err)
		return
	}

	debitNoteLine2, err := GetDebitNoteLine(uint32(2), []byte{69, 8, 4, 218, 212, 211, 65, 69, 163, 166, 34, 158, 239, 225, 55, 161}, "450804da-d4d3-4145-a3a6-229eefe137a1", "Cover is slightly damaged", float64(1), float64(3.96), "2009-06-21T10:04:26Z", "", uint32(3), "2009-11-08T10:04:26Z", "2009-11-09T10:04:26Z", uint32(1), "2019-07-23T10:04:26Z", "2019-07-23T10:04:26Z", "auth0|673c75d516e8adb9e6ffc892", "auth0|673c75d516e8adb9e6ffc892")
	if err != nil {
		t.Error(err)
		return
	}

	debitNoteLine3, err := GetDebitNoteLine(uint32(3), []byte{203, 32, 99, 24, 246, 126, 64, 139, 176, 96, 45, 126, 146, 29, 101, 220}, "cb206318-f67e-408b-b060-2d7e921d65dc", "", float64(2), float64(4.96), "2009-06-21T10:04:26Z", "", uint32(4), "2009-11-08T10:04:26Z", "2009-11-09T10:04:26Z", uint32(1), "2019-07-23T10:04:26Z", "2019-07-23T10:04:26Z", "auth0|673c75d516e8adb9e6ffc892", "auth0|673c75d516e8adb9e6ffc892")
	if err != nil {
		t.Error(err)
		return
	}

	debitNoteLine4, err := GetDebitNoteLine(uint32(4), []byte{101, 12, 230, 48, 19, 229, 75, 231, 168, 85, 146, 245, 90, 166, 66, 39}, "650ce630-13e5-4be7-a855-92f55aa64227", "", float64(1), float64(25), "2009-06-21T10:04:26Z", "", uint32(5), "2009-11-08T10:04:26Z", "2009-11-09T10:04:26Z", uint32(1), "2019-07-23T10:04:26Z", "2019-07-23T10:04:26Z", "auth0|673c75d516e8adb9e6ffc892", "auth0|673c75d516e8adb9e6ffc892")
	if err != nil {
		t.Error(err)
		return
	}

	debitNoteLine5, err := GetDebitNoteLine(uint32(5), []byte{38, 172, 200, 60, 7, 154, 67, 6, 184, 161, 202, 251, 72, 254, 31, 205}, "26acc83c-079a-4306-b8a1-cafb48fe1fcd", "", float64(250), float64(187.5), "2009-06-21T10:04:26Z", "BookingCode002", uint32(6), "2009-11-08T10:04:26Z", "2009-11-09T10:04:26Z", uint32(1), "2019-07-23T10:04:26Z", "2019-07-23T10:04:26Z", "auth0|673c75d516e8adb9e6ffc892", "auth0|673c75d516e8adb9e6ffc892")
	if err != nil {
		t.Error(err)
		return
	}

	debitNoteLines = append(debitNoteLines, debitNoteLine, debitNoteLine2, debitNoteLine3, debitNoteLine4, debitNoteLine5)

	debitNoteLinesResponse := invoiceproto.GetDebitNoteLinesResponse{}
	debitNoteLinesResponse.DebitNoteLines = debitNoteLines

	form := invoiceproto.GetDebitNoteLinesRequest{}
	gform := commonproto.GetRequest{}
	gform.Id = "930f8806-db24-4562-b8c7-72df75518355"
	gform.UserEmail = "sprov300@gmail.com"
	gform.RequestId = "bks1m1g91jau4nkks2f0"
	form.GetRequest = &gform

	type args struct {
		ctx context.Context
		in  *invoiceproto.GetDebitNoteLinesRequest
	}
	tests := []struct {
		ds      *DebitNoteHeaderService
		args    args
		want    *invoiceproto.GetDebitNoteLinesResponse
		wantErr bool
	}{
		{
			ds: debitNoteHeaderService,
			args: args{
				ctx: ctx,
				in:  &form,
			},
			want:    &debitNoteLinesResponse,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		debitNoteLinesResp, err := tt.ds.GetDebitNoteLines(tt.args.ctx, tt.args.in)
		if (err != nil) != tt.wantErr {
			t.Errorf("DebitNoteHeaderService.GetDebitNoteLines() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(debitNoteLinesResp, tt.want) {
			t.Errorf("DebitNoteHeaderService.GetDebitNoteLines() = %v, want %v", debitNoteLinesResp, tt.want)
		}
		assert.NotNil(t, debitNoteLinesResp)
		debitNoteLineResult := debitNoteLinesResp.DebitNoteLines[0]
		assert.Equal(t, debitNoteLineResult.DebitNoteLineD.Note, "Scratch on box", "they should be equal")
		assert.Equal(t, debitNoteLineResult.DebitNoteLineD.DebitedQuantity, float64(1), "they should be equal")
		assert.Equal(t, debitNoteLineResult.DebitNoteLineD.LineExtensionAmount, float64(1273), "they should be equal")
		assert.Equal(t, debitNoteLineResult.DebitNoteLineD.AccountingCost, "BookingCode002", "they should be equal")
	}
}

func TestDebitNoteHeaderService_CreateDebitNoteHeader(t *testing.T) {
	err := test.LoadSQL(logUser, dbService)
	if err != nil {
		t.Error(err)
		return
	}

	ctx := LoginUser()

	debitNoteHeaderService := NewDebitNoteHeaderService(log, dbService, redisService, userServiceClient)
	debitNoteHeader := invoiceproto.CreateDebitNoteHeaderRequest{}
	debitNoteHeader.IssueDate = "11/15/2022"
	debitNoteHeader.TaxPointDate = "11/25/2022"
	debitNoteHeader.Note = "Ordered in our booth at the convention"
	debitNoteHeader.DocumentCurrencyCode = "EUR"
	debitNoteHeader.AccountingCost = "Project cost code 123"
	debitNoteHeader.InvoicePeriodStartDate = "11/15/2022"
	debitNoteHeader.InvoicePeriodEndDate = "12/01/2022"
	debitNoteHeader.TaxExDate = "11/25/2022"
	debitNoteHeader.PricingExDate = "11/25/2022"
	debitNoteHeader.PaymentExDate = "11/25/2022"
	debitNoteHeader.PaymentAltExDate = "11/25/2022"
	debitNoteHeader.ChargeTotalAmount = float64(0)
	debitNoteHeader.PrepaidAmount = float64(0)
	debitNoteHeader.PayableRoundingAmount = float64(0)
	debitNoteHeader.PayableAmount = float64(1250)
	debitNoteHeader.UserId = "auth0|673c75d516e8adb9e6ffc892"
	debitNoteHeader.UserEmail = "sprov300@gmail.com"
	debitNoteHeader.RequestId = "bks1m1g91jau4nkks2f0"

	debitNoteLine := invoiceproto.CreateDebitNoteLineRequest{}
	debitNoteLine.Note = "Scratch on box"
	debitNoteLine.DebitedQuantity = float64(1)
	debitNoteLine.LineExtensionAmount = float64(1250)
	debitNoteLine.TaxPointDate = "11/25/2022"
	debitNoteLine.AccountingCost = "BookingCode002"
	debitNoteLine.ItemId = uint32(7)
	debitNoteLine.PriceValidityPeriodStartDate = "11/15/2022"
	debitNoteLine.PriceValidityPeriodEndDate = "12/01/2022"
	debitNoteLine.UserId = "auth0|673c75d516e8adb9e6ffc892"
	debitNoteLine.UserEmail = "sprov300@gmail.com"
	debitNoteLine.RequestId = "bks1m1g91jau4nkks2f0"

	debitNoteLines := []*invoiceproto.CreateDebitNoteLineRequest{}
	debitNoteLines = append(debitNoteLines, &debitNoteLine)
	debitNoteHeader.DebitNoteLines = debitNoteLines

	type args struct {
		ctx context.Context
		in  *invoiceproto.CreateDebitNoteHeaderRequest
	}
	tests := []struct {
		ds      *DebitNoteHeaderService
		args    args
		wantErr bool
	}{
		{
			ds: debitNoteHeaderService,
			args: args{
				ctx: ctx,
				in:  &debitNoteHeader,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		debitNoteHeaderResp, err := tt.ds.CreateDebitNoteHeader(tt.args.ctx, tt.args.in)
		if (err != nil) != tt.wantErr {
			t.Errorf("DebitNoteHeaderService.CreateDebitNoteHeader() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		assert.NotNil(t, debitNoteHeaderResp)
		debitNoteHeaderResult := debitNoteHeaderResp.DebitNoteHeader
		assert.Equal(t, debitNoteHeaderResult.DebitNoteHeaderD.Note, "Ordered in our booth at the convention", "they should be equal")
		assert.Equal(t, debitNoteHeaderResult.DebitNoteHeaderD.DocumentCurrencyCode, "EUR", "they should be equal")
		assert.Equal(t, debitNoteHeaderResult.DebitNoteHeaderD.AccountingCost, "Project cost code 123", "they should be equal")
	}
}

func GetDebitNoteHeader(id uint32, uuid4 []byte, idS string, issueDate string, note string, taxPointDate string, documentCurrencyCode string, accountingCost string, invoicePeriodStartDate string, invoicePeriodEndDate string, taxExDate string, pricingExDate string, paymentExDate string, paymentAltExDate string, createdAt string, updatedAt string, createdByUserId string, updatedByUserId string) (*invoiceproto.DebitNoteHeader, error) {
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

	debitNoteHeaderD := invoiceproto.DebitNoteHeaderD{}
	debitNoteHeaderD.Id = id
	debitNoteHeaderD.Uuid4 = uuid4
	debitNoteHeaderD.IdS = idS
	debitNoteHeaderD.Note = note
	debitNoteHeaderD.DocumentCurrencyCode = documentCurrencyCode
	debitNoteHeaderD.AccountingCost = accountingCost

	debitNoteHeaderT := invoiceproto.DebitNoteHeaderT{}
	debitNoteHeaderT.IssueDate = issueDate1
	debitNoteHeaderT.TaxPointDate = taxPointDate1
	debitNoteHeaderT.InvoicePeriodStartDate = invoicePeriodStartDate1
	debitNoteHeaderT.InvoicePeriodEndDate = invoicePeriodEndDate1
	debitNoteHeaderT.TaxExDate = taxExDate1
	debitNoteHeaderT.PricingExDate = pricingExDate1
	debitNoteHeaderT.PaymentExDate = paymentExDate1
	debitNoteHeaderT.PaymentAltExDate = paymentAltExDate1

	crUpdTime := new(commonproto.CrUpdTime)
	crUpdTime.CreatedAt = createdAt1
	crUpdTime.UpdatedAt = updatedAt1

	crUpdUser := new(commonproto.CrUpdUser)
	crUpdUser.StatusCode = "active"
	crUpdUser.CreatedByUserId = createdByUserId
	crUpdUser.UpdatedByUserId = updatedByUserId

	debitNoteHeader := invoiceproto.DebitNoteHeader{DebitNoteHeaderD: &debitNoteHeaderD, DebitNoteHeaderT: &debitNoteHeaderT, CrUpdUser: crUpdUser, CrUpdTime: crUpdTime}

	return &debitNoteHeader, nil
}

func GetDebitNoteLine(id uint32, uuid4 []byte, idS string, note string, debitedQuantity float64, lineExtensionAmount float64, taxPointDate string, accountingCost string, itemId uint32, priceValidityPeriodStartDate string, priceValidityPeriodEndDate string, debitNoteHeaderId uint32, createdAt string, updatedAt string, createdByUserId string, updatedByUserId string) (*invoiceproto.DebitNoteLine, error) {
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

	priceValidityPeriodStartDate1, err := common.ConvertTimeToTimestamp(Layout, priceValidityPeriodStartDate)
	if err != nil {
		log.Error("Error", zap.Error(err))
	}

	priceValidityPeriodEndDate1, err := common.ConvertTimeToTimestamp(Layout, priceValidityPeriodEndDate)
	if err != nil {
		log.Error("Error", zap.Error(err))
	}

	debitNoteLineD := invoiceproto.DebitNoteLineD{}
	debitNoteLineD.Id = id
	debitNoteLineD.Uuid4 = uuid4
	debitNoteLineD.IdS = idS
	debitNoteLineD.Note = note
	debitNoteLineD.DebitedQuantity = debitedQuantity
	debitNoteLineD.LineExtensionAmount = lineExtensionAmount
	debitNoteLineD.AccountingCost = accountingCost
	debitNoteLineD.ItemId = itemId
	debitNoteLineD.DebitNoteHeaderId = debitNoteHeaderId

	debitNoteLineT := invoiceproto.DebitNoteLineT{}
	debitNoteLineT.TaxPointDate = taxPointDate1
	debitNoteLineT.PriceValidityPeriodStartDate = priceValidityPeriodStartDate1
	debitNoteLineT.PriceValidityPeriodEndDate = priceValidityPeriodEndDate1

	crUpdTime := new(commonproto.CrUpdTime)
	crUpdTime.CreatedAt = createdAt1
	crUpdTime.UpdatedAt = updatedAt1

	crUpdUser := new(commonproto.CrUpdUser)
	crUpdUser.StatusCode = "active"
	crUpdUser.CreatedByUserId = createdByUserId
	crUpdUser.UpdatedByUserId = updatedByUserId

	debitNoteLine := invoiceproto.DebitNoteLine{DebitNoteLineD: &debitNoteLineD, DebitNoteLineT: &debitNoteLineT, CrUpdUser: crUpdUser, CrUpdTime: crUpdTime}

	return &debitNoteLine, nil
}

func TestDebitNoteHeaderService_UpdateDebitNoteHeader(t *testing.T) {
	err := test.LoadSQL(logUser, dbService)
	if err != nil {
		t.Error(err)
		return
	}

	ctx := LoginUser()

	debitNoteHeaderService := NewDebitNoteHeaderService(log, dbService, redisService, userServiceClient)
	form := invoiceproto.UpdateDebitNoteHeaderRequest{}
	form.Note = "Ordered"
	form.DocumentCurrencyCode = "EUR"
	form.ChargeTotalAmount = float64(100)
	form.PrepaidAmount = float64(60)
	form.PayableRoundingAmount = float64(100)
	form.PayableAmount = float64(300)
	form.Id = "930f8806-db24-4562-b8c7-72df75518355"
	form.UserId = "auth0|673c75d516e8adb9e6ffc892"
	form.UserEmail = "sprov300@gmail.com"
	form.RequestId = "bks1m1g91jau4nkks2f0"

	updateResponse := invoiceproto.UpdateDebitNoteHeaderResponse{}

	type args struct {
		ctx context.Context
		in  *invoiceproto.UpdateDebitNoteHeaderRequest
	}
	tests := []struct {
		ds      *DebitNoteHeaderService
		args    args
		want    *invoiceproto.UpdateDebitNoteHeaderResponse
		wantErr bool
	}{
		{
			ds: debitNoteHeaderService,
			args: args{
				ctx: ctx,
				in:  &form,
			},
			want:    &updateResponse,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		got, err := tt.ds.UpdateDebitNoteHeader(tt.args.ctx, tt.args.in)
		if (err != nil) != tt.wantErr {
			t.Errorf("DebitNoteHeaderService.UpdateDebitNoteHeader() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("DebitNoteHeaderService.UpdateDebitNoteHeader() = %v, want %v", got, tt.want)
		}
	}
}
