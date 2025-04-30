package taxservices

import (
	"context"
	"reflect"
	"testing"

	"github.com/cloudfresco/sc-ubl/internal/common"
	commonproto "github.com/cloudfresco/sc-ubl/internal/protogen/common/v1"
	taxproto "github.com/cloudfresco/sc-ubl/internal/protogen/tax/v1"
	"github.com/cloudfresco/sc-ubl/test"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestTaxService_GetTaxSchemes(t *testing.T) {
	err := test.LoadSQL(logUser, dbService)
	if err != nil {
		t.Error(err)
		return
	}

	ctx := LoginUser()

	taxService := NewTaxService(log, dbService, redisService, userServiceClient)

	taxSchemes := []*taxproto.TaxScheme{}

	taxScheme1, err := GetTaxScheme(uint32(13), []byte{244, 58, 127, 61, 55, 60, 77, 197, 141, 9, 187, 93, 78, 65, 99, 63}, "f43a7f3d-373c-4dc5-8d09-bb5d4e41633f", "VAT", "", "VAT", "", "2019-07-23T10:04:26Z", "2019-07-23T10:04:26Z", "auth0|673c75d516e8adb9e6ffc892", "auth0|673c75d516e8adb9e6ffc892")
	if err != nil {
		t.Error(err)
		return
	}

	taxScheme2, err := GetTaxScheme(uint32(14), []byte{223, 155, 40, 67, 127, 30, 76, 169, 151, 80, 98, 232, 239, 239, 125, 222}, "df9b2843-7f1e-4ca9-9750-62e8efef7dde", "UK VAT", "", "VAT", "", "2019-07-23T10:04:26Z", "2019-07-23T10:04:26Z", "auth0|673c75d516e8adb9e6ffc892", "auth0|673c75d516e8adb9e6ffc892")
	if err != nil {
		t.Error(err)
		return
	}

	taxSchemes = append(taxSchemes, taxScheme2, taxScheme1)

	form := taxproto.GetTaxSchemesRequest{}
	form.Limit = "2"
	form.NextCursor = ""
	form.UserEmail = "sprov300@gmail.com"
	form.RequestId = "bks1m1g91jau4nkks2f0"

	nextc := "MTI="
	taxSchemesResponse := taxproto.GetTaxSchemesResponse{TaxSchemes: taxSchemes, NextCursor: nextc}

	type args struct {
		ctx context.Context
		in  *taxproto.GetTaxSchemesRequest
	}
	tests := []struct {
		ts      *TaxService
		args    args
		want    *taxproto.GetTaxSchemesResponse
		wantErr bool
	}{
		{
			ts: taxService,
			args: args{
				ctx: ctx,
				in:  &form,
			},
			want:    &taxSchemesResponse,
			wantErr: false,
		},
	}
	for _, tt := range tests {

		taxSchemeResp, err := tt.ts.GetTaxSchemes(tt.args.ctx, tt.args.in)
		if (err != nil) != tt.wantErr {
			t.Errorf("TaxService.GetTaxSchemes() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(taxSchemeResp, tt.want) {
			t.Errorf("TaxService.GetTaxSchemes() = %v, want %v", taxSchemeResp, tt.want)
		}
		assert.NotNil(t, taxSchemeResp)
		taxSchemeResult := taxSchemeResp.TaxSchemes[0]
		assert.Equal(t, taxSchemeResult.TaxSchemeD.TsId, "UK VAT", "they should be equal")
		assert.Equal(t, taxSchemeResult.TaxSchemeD.TaxTypeCode, "VAT", "they should be equal")

	}
}

func TestTaxService_GetTaxScheme(t *testing.T) {
	err := test.LoadSQL(logUser, dbService)
	if err != nil {
		t.Error(err)
		return
	}

	ctx := LoginUser()

	taxService := NewTaxService(log, dbService, redisService, userServiceClient)

	taxScheme, err := GetTaxScheme(uint32(11), []byte{81, 166, 138, 130, 149, 6, 72, 81, 157, 218, 209, 103, 152, 168, 139, 135}, "51a68a82-9506-4851-9dda-d16798a88b87", "VAT", "", "VAT", "", "2019-07-23T10:04:26Z", "2019-07-23T10:04:26Z", "auth0|673c75d516e8adb9e6ffc892", "auth0|673c75d516e8adb9e6ffc892")
	if err != nil {
		t.Error(err)
		return
	}

	form := taxproto.GetTaxSchemeRequest{}
	gform := commonproto.GetRequest{}
	gform.Id = "51a68a82-9506-4851-9dda-d16798a88b87"
	gform.UserEmail = "sprov300@gmail.com"
	gform.RequestId = "bks1m1g91jau4nkks2f0"
	form.GetRequest = &gform

	taxSchemeResponse := taxproto.GetTaxSchemeResponse{}
	taxSchemeResponse.TaxScheme = taxScheme

	type args struct {
		ctx   context.Context
		inReq *taxproto.GetTaxSchemeRequest
	}

	tests := []struct {
		ts      *TaxService
		args    args
		want    *taxproto.GetTaxSchemeResponse
		wantErr bool
	}{
		{
			ts: taxService,
			args: args{
				ctx:   ctx,
				inReq: &form,
			},
			want:    &taxSchemeResponse,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		taxSchemeResp, err := tt.ts.GetTaxScheme(tt.args.ctx, tt.args.inReq)
		if (err != nil) != tt.wantErr {
			t.Errorf("TaxService.GetTaxScheme() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(taxSchemeResp, tt.want) {
			t.Errorf("TaxService.GetTaxScheme() = %v, want %v", taxSchemeResp, tt.want)
		}
		assert.NotNil(t, taxSchemeResp)
		taxSchemeResult := taxSchemeResp.TaxScheme
		assert.Equal(t, taxSchemeResult.TaxSchemeD.TsId, "VAT", "they should be equal")
		assert.Equal(t, taxSchemeResult.TaxSchemeD.TaxTypeCode, "VAT", "they should be equal")
	}
}

func TestTaxService_CreateTaxScheme(t *testing.T) {
	err := test.LoadSQL(logUser, dbService)
	if err != nil {
		t.Error(err)
		return
	}

	ctx := LoginUser()

	taxService := NewTaxService(log, dbService, redisService, userServiceClient)

	taxScheme := taxproto.CreateTaxSchemeRequest{}
	taxScheme.TsId = "UK VAT"
	taxScheme.TaxSchemeName = "TaxSchemeName"
	taxScheme.TaxTypeCode = "VAT"
	taxScheme.CurrencyCode = "EUR"
	taxScheme.UserId = "auth0|673c75d516e8adb9e6ffc892"
	taxScheme.UserEmail = "sprov300@gmail.com"
	taxScheme.RequestId = "bks1m1g91jau4nkks2f0"

	type args struct {
		ctx context.Context
		in  *taxproto.CreateTaxSchemeRequest
	}
	tests := []struct {
		ts      *TaxService
		args    args
		wantErr bool
	}{
		{
			ts: taxService,
			args: args{
				ctx: ctx,
				in:  &taxScheme,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		taxSchemeResp, err := tt.ts.CreateTaxScheme(tt.args.ctx, tt.args.in)
		if (err != nil) != tt.wantErr {
			t.Errorf("TaxService.CreateTaxScheme() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		assert.NotNil(t, taxSchemeResp)
		taxSchemeResult := taxSchemeResp.TaxScheme
		assert.Equal(t, taxSchemeResult.TaxSchemeD.TsId, "UK VAT", "they should be equal")
		assert.Equal(t, taxSchemeResult.TaxSchemeD.TaxTypeCode, "VAT", "they should be equal")
		assert.Equal(t, taxSchemeResult.TaxSchemeD.CurrencyCode, "EUR", "they should be equal")
	}
}

func TestTaxService_UpdateTaxScheme(t *testing.T) {
	err := test.LoadSQL(logUser, dbService)
	if err != nil {
		t.Error(err)
		return
	}

	ctx := LoginUser()

	taxService := NewTaxService(log, dbService, redisService, userServiceClient)

	taxScheme := taxproto.UpdateTaxSchemeRequest{}
	taxScheme.TaxSchemeName = "TaxSchemeName"
	taxScheme.TaxTypeCode = "VAT"
	taxScheme.CurrencyCode = "EUR"
	taxScheme.Id = "15fd8632-8674-462d-a08a-de1560d2d9e8"
	taxScheme.UserId = "auth0|673c75d516e8adb9e6ffc892"
	taxScheme.UserEmail = "sprov300@gmail.com"
	taxScheme.RequestId = "bks1m1g91jau4nkks2f0"

	updateResponse := taxproto.UpdateTaxSchemeResponse{}

	type args struct {
		ctx context.Context
		in  *taxproto.UpdateTaxSchemeRequest
	}
	tests := []struct {
		ts      *TaxService
		args    args
		want    *taxproto.UpdateTaxSchemeResponse
		wantErr bool
	}{
		{
			ts: taxService,
			args: args{
				ctx: ctx,
				in:  &taxScheme,
			},
			want:    &updateResponse,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		got, err := tt.ts.UpdateTaxScheme(tt.args.ctx, tt.args.in)
		if (err != nil) != tt.wantErr {
			t.Errorf("TaxService.UpdateTaxScheme() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("TaxService.UpdateTaxScheme() = %v, want %v", got, tt.want)
		}
	}
}

func TestTaxService_CreateTaxCategory(t *testing.T) {
	err := test.LoadSQL(logUser, dbService)
	if err != nil {
		t.Error(err)
		return
	}

	ctx := LoginUser()

	taxService := NewTaxService(log, dbService, redisService, userServiceClient)

	taxCategory := taxproto.CreateTaxCategoryRequest{}
	taxCategory.TcId = ""
	taxCategory.TaxCategoryName = "TaxCategory"
	taxCategory.Percent = float32(20)
	taxCategory.BaseUnitMeasure = "EUR"
	taxCategory.PerUnitAmount = float64(10)
	taxCategory.TaxExemptionReasonCode = ""
	taxCategory.TaxExemptionReason = ""
	taxCategory.TierRange = ""
	taxCategory.TierRatePercent = float32(10)
	taxCategory.TaxSchemeId = "1e275aa8-8ff1-46e4-b5c2-a8dbbbb1a231"
	taxCategory.UserId = "auth0|673c75d516e8adb9e6ffc892"
	taxCategory.UserEmail = "sprov300@gmail.com"
	taxCategory.RequestId = "bks1m1g91jau4nkks2f0"

	type args struct {
		ctx context.Context
		in  *taxproto.CreateTaxCategoryRequest
	}
	tests := []struct {
		ts      *TaxService
		args    args
		wantErr bool
	}{
		{
			ts: taxService,
			args: args{
				ctx: ctx,
				in:  &taxCategory,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		taxCategoryResp, err := tt.ts.CreateTaxCategory(tt.args.ctx, tt.args.in)
		if (err != nil) != tt.wantErr {
			t.Errorf("TaxService.CreateTaxCategory() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		assert.NotNil(t, taxCategoryResp)
		taxCategoryResult := taxCategoryResp.TaxCategory
		assert.Equal(t, taxCategoryResult.TaxCategoryD.Percent, float32(20), "they should be equal")
		assert.Equal(t, taxCategoryResult.TaxCategoryD.PerUnitAmount, float64(10), "they should be equal")
	}
}

func TestTaxService_UpdateTaxCategory(t *testing.T) {
	err := test.LoadSQL(logUser, dbService)
	if err != nil {
		t.Error(err)
		return
	}

	ctx := LoginUser()

	taxService := NewTaxService(log, dbService, redisService, userServiceClient)

	taxCategory := taxproto.UpdateTaxCategoryRequest{}
	taxCategory.TaxCategoryName = "TaxCategory1"
	taxCategory.Percent = float32(10)
	taxCategory.BaseUnitMeasure = "EUR"
	taxCategory.PerUnitAmount = float64(10)
	taxCategory.TaxExemptionReasonCode = ""
	taxCategory.TaxExemptionReason = ""
	taxCategory.Id = "7036e24c-0ec5-48ac-bc97-4c2cbe75a54c"
	taxCategory.UserId = "auth0|673c75d516e8adb9e6ffc892"
	taxCategory.UserEmail = "sprov300@gmail.com"
	taxCategory.RequestId = "bks1m1g91jau4nkks2f0"

	updateResponse := taxproto.UpdateTaxCategoryResponse{}

	type args struct {
		ctx context.Context
		in  *taxproto.UpdateTaxCategoryRequest
	}
	tests := []struct {
		ts      *TaxService
		args    args
		want    *taxproto.UpdateTaxCategoryResponse
		wantErr bool
	}{
		{
			ts: taxService,
			args: args{
				ctx: ctx,
				in:  &taxCategory,
			},
			want:    &updateResponse,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		got, err := tt.ts.UpdateTaxCategory(tt.args.ctx, tt.args.in)
		if (err != nil) != tt.wantErr {
			t.Errorf("TaxService.UpdateTaxCategory() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("TaxService.UpdateTaxCategory() = %v, want %v", got, tt.want)
		}

	}
}

func TestTaxService_CreateTaxTotal(t *testing.T) {
	err := test.LoadSQL(logUser, dbService)
	if err != nil {
		t.Error(err)
		return
	}

	ctx := LoginUser()

	taxService := NewTaxService(log, dbService, redisService, userServiceClient)

	taxTotal := taxproto.CreateTaxTotalRequest{}
	taxTotal.TaxAmount = float64(17.5)
	taxTotal.RoundingAmount = float64(18)
	taxTotal.TaxEvidenceIndicator = false
	taxTotal.TaxIncludedIndicator = false
	taxTotal.MasterFlag = "CNL"
	taxTotal.MasterId = uint32(1)
	taxTotal.TaxCategoryId = "7036e24c-0ec5-48ac-bc97-4c2cbe75a54c"
	taxTotal.UserId = "auth0|673c75d516e8adb9e6ffc892"
	taxTotal.UserEmail = "sprov300@gmail.com"
	taxTotal.RequestId = "bks1m1g91jau4nkks2f0"

	type args struct {
		ctx context.Context
		in  *taxproto.CreateTaxTotalRequest
	}
	tests := []struct {
		ts      *TaxService
		args    args
		wantErr bool
	}{
		{
			ts: taxService,
			args: args{
				ctx: ctx,
				in:  &taxTotal,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {

		taxTotalResp, err := tt.ts.CreateTaxTotal(tt.args.ctx, tt.args.in)
		if (err != nil) != tt.wantErr {
			t.Errorf("TaxService.CreateTaxTotal() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		assert.NotNil(t, taxTotalResp)
		taxTotalResult := taxTotalResp.TaxTotal
		assert.Equal(t, taxTotalResult.TaxTotalD.TaxAmount, float64(17.5), "they should be equal")
		assert.False(t, taxTotalResult.TaxTotalD.TaxEvidenceIndicator, "Its False")
		assert.Equal(t, taxTotalResult.TaxTotalD.MasterFlag, "CNL", "they should be equal")
	}
}

func TestTaxService_UpdateTaxTotal(t *testing.T) {
	err := test.LoadSQL(logUser, dbService)
	if err != nil {
		t.Error(err)
		return
	}

	ctx := LoginUser()

	taxService := NewTaxService(log, dbService, redisService, userServiceClient)

	taxTotal := taxproto.UpdateTaxTotalRequest{}
	taxTotal.TaxAmount = float64(18.5)
	taxTotal.RoundingAmount = float64(19)
	taxTotal.Id = "a62007af-a514-45a5-967d-270ad8c34f91"
	taxTotal.UserId = "auth0|673c75d516e8adb9e6ffc892"
	taxTotal.UserEmail = "sprov300@gmail.com"
	taxTotal.RequestId = "bks1m1g91jau4nkks2f0"

	updateResponse := taxproto.UpdateTaxTotalResponse{}
	type args struct {
		ctx context.Context
		in  *taxproto.UpdateTaxTotalRequest
	}
	tests := []struct {
		ts      *TaxService
		args    args
		want    *taxproto.UpdateTaxTotalResponse
		wantErr bool
	}{
		{
			ts: taxService,
			args: args{
				ctx: ctx,
				in:  &taxTotal,
			},
			want:    &updateResponse,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		got, err := tt.ts.UpdateTaxTotal(tt.args.ctx, tt.args.in)
		if (err != nil) != tt.wantErr {
			t.Errorf("TaxService.UpdateTaxTotal() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("TaxService.UpdateTaxTotal() = %v, want %v", got, tt.want)
		}
	}
}

func TestTaxService_CreateTaxSubTotal(t *testing.T) {
	err := test.LoadSQL(logUser, dbService)
	if err != nil {
		t.Error(err)
		return
	}

	ctx := LoginUser()

	taxService := NewTaxService(log, dbService, redisService, userServiceClient)

	taxSubTotal := taxproto.CreateTaxSubTotalRequest{}
	taxSubTotal.TaxableAmount = float64(100)
	taxSubTotal.TaxAmount = float64(17.5)
	taxSubTotal.CalculationSequenceNumeric = uint32(0)
	taxSubTotal.TransactionCurrencyTaxAmount = float64(0)
	taxSubTotal.Percent = float32(10)
	taxSubTotal.BaseUnitMeasure = "EUR"
	taxSubTotal.PerUnitAmount = float64(10)
	taxSubTotal.TierRange = ""
	taxSubTotal.TierRatePercent = float64(10)
	taxSubTotal.TaxCategoryId = uint32(1)
	taxSubTotal.TaxTotalId = "ddee759b-936d-4935-a60c-b11c2109ffe9"
	taxSubTotal.UserId = "auth0|673c75d516e8adb9e6ffc892"
	taxSubTotal.UserEmail = "sprov300@gmail.com"
	taxSubTotal.RequestId = "bks1m1g91jau4nkks2f0"

	type args struct {
		ctx context.Context
		in  *taxproto.CreateTaxSubTotalRequest
	}
	tests := []struct {
		ts      *TaxService
		args    args
		wantErr bool
	}{
		{
			ts: taxService,
			args: args{
				ctx: ctx,
				in:  &taxSubTotal,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		taxSubTotalResp, err := tt.ts.CreateTaxSubTotal(tt.args.ctx, tt.args.in)
		if (err != nil) != tt.wantErr {
			t.Errorf("TaxService.CreateTaxSubTotal() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		assert.NotNil(t, taxSubTotalResp)
		taxSubTotalResult := taxSubTotalResp.TaxSubTotal
		assert.Equal(t, taxSubTotalResult.TaxSubTotalD.TaxableAmount, float64(100), "they should be equal")
		assert.Equal(t, taxSubTotalResult.TaxSubTotalD.TaxAmount, float64(17.5), "they should be equal")
	}
}

func TestTaxService_UpdateTaxSubTotal(t *testing.T) {
	err := test.LoadSQL(logUser, dbService)
	if err != nil {
		t.Error(err)
		return
	}

	ctx := LoginUser()

	taxService := NewTaxService(log, dbService, redisService, userServiceClient)

	taxSubTotal := taxproto.UpdateTaxSubTotalRequest{}
	taxSubTotal.TaxableAmount = float64(200)
	taxSubTotal.TaxAmount = float64(18.5)
	taxSubTotal.CalculationSequenceNumeric = uint32(0)
	taxSubTotal.TransactionCurrencyTaxAmount = float64(0)
	taxSubTotal.Percent = float32(15)
	taxSubTotal.BaseUnitMeasure = "EUR"
	taxSubTotal.PerUnitAmount = float64(5)
	taxSubTotal.Id = "f7da0aaa-8068-495c-9176-b7681e651226"
	taxSubTotal.UserId = "auth0|673c75d516e8adb9e6ffc892"
	taxSubTotal.UserEmail = "sprov300@gmail.com"
	taxSubTotal.RequestId = "bks1m1g91jau4nkks2f0"

	updateResponse := taxproto.UpdateTaxSubTotalResponse{}

	type args struct {
		ctx context.Context
		in  *taxproto.UpdateTaxSubTotalRequest
	}
	tests := []struct {
		ts      *TaxService
		args    args
		want    *taxproto.UpdateTaxSubTotalResponse
		wantErr bool
	}{
		{
			ts: taxService,
			args: args{
				ctx: ctx,
				in:  &taxSubTotal,
			},
			want:    &updateResponse,
			wantErr: false,
		},
	}
	for _, tt := range tests {

		got, err := tt.ts.UpdateTaxSubTotal(tt.args.ctx, tt.args.in)
		if (err != nil) != tt.wantErr {
			t.Errorf("TaxService.UpdateTaxSubTotal() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("TaxService.UpdateTaxSubTotal() = %v, want %v", got, tt.want)
		}

	}
}

func GetTaxScheme(id uint32, uuid4 []byte, idS string, tsId string, taxSchemeName string, taxTypeCode string, currencyCode string, createdAt string, updatedAt string, createdByUserId string, updatedByUserId string) (*taxproto.TaxScheme, error) {
	createdAt1, err := common.ConvertTimeToTimestamp(Layout, createdAt)
	if err != nil {
		log.Error("Error", zap.Error(err))
	}

	updatedAt1, err := common.ConvertTimeToTimestamp(Layout, updatedAt)
	if err != nil {
		log.Error("Error", zap.Error(err))
	}

	taxSchemeD := taxproto.TaxSchemeD{}
	taxSchemeD.Id = id
	taxSchemeD.Uuid4 = uuid4
	taxSchemeD.IdS = idS
	taxSchemeD.TsId = tsId
	taxSchemeD.TaxSchemeName = taxSchemeName
	taxSchemeD.TaxTypeCode = taxTypeCode
	taxSchemeD.CurrencyCode = currencyCode

	crUpdTime := commonproto.CrUpdTime{}
	crUpdTime.CreatedAt = createdAt1
	crUpdTime.UpdatedAt = updatedAt1

	crUpdUser := commonproto.CrUpdUser{}
	crUpdUser.StatusCode = "active"
	crUpdUser.CreatedByUserId = createdByUserId
	crUpdUser.UpdatedByUserId = updatedByUserId

	taxScheme := taxproto.TaxScheme{TaxSchemeD: &taxSchemeD, CrUpdUser: &crUpdUser, CrUpdTime: &crUpdTime}

	return &taxScheme, nil
}
