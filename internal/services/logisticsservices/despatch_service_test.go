package logisticsservices

import (
	"context"
	"reflect"
	"testing"

	"github.com/cloudfresco/sc-ubl/internal/common"
	commonproto "github.com/cloudfresco/sc-ubl/internal/protogen/common/v1"
	logisticsproto "github.com/cloudfresco/sc-ubl/internal/protogen/logistics/v1"
	"github.com/cloudfresco/sc-ubl/test"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestDespatchService_GetDespatchHeaders(t *testing.T) {
	err := test.LoadSQL(logUser, dbService)
	if err != nil {
		t.Error(err)
		return
	}

	ctx := LoginUser()

	despatchHeaderService := NewDespatchService(log, dbService, redisService, userServiceClient)

	despatchHeader, err := GetDespatchHeader(uint32(1), []byte{255, 36, 65, 61, 133, 23, 65, 226, 155, 73, 27, 179, 133, 145, 150, 148}, "ff24413d-8517-41e2-9b49-1bb385919694", "2005-06-20T10:04:26Z", "NoStatus", "delivery", "sample", uint32(1), "2019-07-23T10:04:26Z", "2019-07-23T10:04:26Z", "auth0|673c75d516e8adb9e6ffc892", "auth0|673c75d516e8adb9e6ffc892")
	if err != nil {
		t.Error(err)
		return
	}
	despatchHeaders := []*logisticsproto.DespatchHeader{}
	despatchHeaders = append(despatchHeaders, despatchHeader)

	form := logisticsproto.GetDespatchHeadersRequest{}
	form.Limit = "8"
	form.NextCursor = ""
	form.UserEmail = "sprov300@gmail.com"
	form.RequestId = "bks1m1g91jau4nkks2f0"

	nextc := "MA=="
	despatchHeadersResponse := logisticsproto.GetDespatchHeadersResponse{DespatchHeaders: despatchHeaders, NextCursor: nextc}

	type args struct {
		ctx context.Context
		in  *logisticsproto.GetDespatchHeadersRequest
	}
	tests := []struct {
		ds      *DespatchService
		args    args
		want    *logisticsproto.GetDespatchHeadersResponse
		wantErr bool
	}{
		{
			ds: despatchHeaderService,
			args: args{
				ctx: ctx,
				in:  &form,
			},
			want:    &despatchHeadersResponse,
			wantErr: false,
		},
	}
	for _, tt := range tests {

		despatchHeadersResp, err := tt.ds.GetDespatchHeaders(tt.args.ctx, tt.args.in)
		if (err != nil) != tt.wantErr {
			t.Errorf("DespatchService.GetDespatchHeaders() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(despatchHeadersResp, tt.want) {
			t.Errorf("DespatchService.GetDespatchHeaders() = %v, want %v", despatchHeadersResp, tt.want)
		}

		assert.NotNil(t, despatchHeadersResp)
		despatchHeaderResult := despatchHeadersResp.DespatchHeaders[0]
		assert.Equal(t, despatchHeaderResult.DespatchHeaderD.DocumentStatusCode, "NoStatus", "they should be equal")
		assert.Equal(t, despatchHeaderResult.DespatchHeaderD.DespatchAdviceTypeCode, "delivery", "they should be equal")
		assert.Equal(t, despatchHeaderResult.DespatchHeaderD.Note, "sample", "they should be equal")
	}
}

func TestDespatchService_GetDespatchHeader(t *testing.T) {
	err := test.LoadSQL(logUser, dbService)
	if err != nil {
		t.Error(err)
		return
	}
	ctx := LoginUser()
	despatchHeaderService := NewDespatchService(log, dbService, redisService, userServiceClient)

	despatchHeader, err := GetDespatchHeader(uint32(1), []byte{255, 36, 65, 61, 133, 23, 65, 226, 155, 73, 27, 179, 133, 145, 150, 148}, "ff24413d-8517-41e2-9b49-1bb385919694", "2005-06-20T10:04:26Z", "NoStatus", "delivery", "sample", uint32(1), "2019-07-23T10:04:26Z", "2019-07-23T10:04:26Z", "auth0|673c75d516e8adb9e6ffc892", "auth0|673c75d516e8adb9e6ffc892")
	if err != nil {
		t.Error(err)
		return
	}
	despatchHeaderResponse := logisticsproto.GetDespatchHeaderResponse{}
	despatchHeaderResponse.DespatchHeader = despatchHeader
	gform := commonproto.GetRequest{}
	gform.Id = "ff24413d-8517-41e2-9b49-1bb385919694"
	gform.UserEmail = "sprov300@gmail.com"
	gform.RequestId = "bks1m1g91jau4nkks2f0"

	form := logisticsproto.GetDespatchHeaderRequest{}
	form.GetRequest = &gform

	type args struct {
		ctx context.Context
		in  *logisticsproto.GetDespatchHeaderRequest
	}
	tests := []struct {
		ds      *DespatchService
		args    args
		want    *logisticsproto.GetDespatchHeaderResponse
		wantErr bool
	}{
		{
			ds: despatchHeaderService,
			args: args{
				ctx: ctx,
				in:  &form,
			},
			want:    &despatchHeaderResponse,
			wantErr: false,
		},
	}
	for _, tt := range tests {

		despatchHeaderResp, err := tt.ds.GetDespatchHeader(tt.args.ctx, tt.args.in)
		if (err != nil) != tt.wantErr {
			t.Errorf("DespatchService.GetDespatchHeader() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(despatchHeaderResp, tt.want) {
			t.Errorf("DespatchService.GetDespatchHeader() = %v, want %v", despatchHeaderResp, tt.want)
		}
		assert.NotNil(t, despatchHeaderResp)
		despatchHeaderResult := despatchHeaderResp.DespatchHeader
		assert.Equal(t, despatchHeaderResult.DespatchHeaderD.DocumentStatusCode, "NoStatus", "they should be equal")
		assert.Equal(t, despatchHeaderResult.DespatchHeaderD.DespatchAdviceTypeCode, "delivery", "they should be equal")
		assert.Equal(t, despatchHeaderResult.DespatchHeaderD.Note, "sample", "they should be equal")

	}
}

func TestDespatchService_GetDespatchHeaderByPk(t *testing.T) {
	err := test.LoadSQL(logUser, dbService)
	if err != nil {
		t.Error(err)
		return
	}

	ctx := LoginUser()

	despatchHeaderService := NewDespatchService(log, dbService, redisService, userServiceClient)

	despatchHeader, err := GetDespatchHeader(uint32(1), []byte{255, 36, 65, 61, 133, 23, 65, 226, 155, 73, 27, 179, 133, 145, 150, 148}, "ff24413d-8517-41e2-9b49-1bb385919694", "2005-06-20T10:04:26Z", "NoStatus", "delivery", "sample", uint32(1), "2019-07-23T10:04:26Z", "2019-07-23T10:04:26Z", "auth0|673c75d516e8adb9e6ffc892", "auth0|673c75d516e8adb9e6ffc892")
	if err != nil {
		t.Error(err)
		return
	}

	despatchHeaderResponse := logisticsproto.GetDespatchHeaderByPkResponse{}
	despatchHeaderResponse.DespatchHeader = despatchHeader

	gform := commonproto.GetByIdRequest{}
	gform.Id = uint32(1)
	gform.UserEmail = "sprov300@gmail.com"
	gform.RequestId = "bks1m1g91jau4nkks2f0"

	form := logisticsproto.GetDespatchHeaderByPkRequest{}
	form.GetByIdRequest = &gform

	type args struct {
		ctx context.Context
		in  *logisticsproto.GetDespatchHeaderByPkRequest
	}
	tests := []struct {
		ds      *DespatchService
		args    args
		want    *logisticsproto.GetDespatchHeaderByPkResponse
		wantErr bool
	}{
		{
			ds: despatchHeaderService,
			args: args{
				ctx: ctx,
				in:  &form,
			},
			want:    &despatchHeaderResponse,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		despatchHeaderResp, err := tt.ds.GetDespatchHeaderByPk(tt.args.ctx, tt.args.in)
		if (err != nil) != tt.wantErr {
			t.Errorf("DespatchService.GetDespatchHeaderByPk() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(despatchHeaderResp, tt.want) {
			t.Errorf("DespatchService.GetDespatchHeaderByPk() = %v, want %v", despatchHeaderResp, tt.want)
		}
		assert.NotNil(t, despatchHeaderResp)
		despatchHeaderResult := despatchHeaderResp.DespatchHeader
		assert.Equal(t, despatchHeaderResult.DespatchHeaderD.DocumentStatusCode, "NoStatus", "they should be equal")
		assert.Equal(t, despatchHeaderResult.DespatchHeaderD.DespatchAdviceTypeCode, "delivery", "they should be equal")
		assert.Equal(t, despatchHeaderResult.DespatchHeaderD.Note, "sample", "they should be equal")
	}
}

func TestDespatchService_CreateDespatchHeader(t *testing.T) {
	err := test.LoadSQL(logUser, dbService)
	if err != nil {
		t.Error(err)
		return
	}

	ctx := LoginUser()

	despatchHeaderService := NewDespatchService(log, dbService, redisService, userServiceClient)
	despatchHeader := logisticsproto.CreateDespatchHeaderRequest{}
	despatchHeader.IssueDate = "06/20/2022"
	despatchHeader.DocumentStatusCode = "NoStatus"
	despatchHeader.DespatchAdviceTypeCode = "delivery"
	despatchHeader.Note = "sample"
	despatchHeader.OrderId = uint32(1)
	despatchHeader.UserId = "auth0|673c75d516e8adb9e6ffc892"
	despatchHeader.UserEmail = "sprov300@gmail.com"
	despatchHeader.RequestId = "bks1m1g91jau4nkks2f0"

	despatchLine := logisticsproto.CreateDespatchLineRequest{}
	despatchLine.Note = "Mrs Green agreed to waive charge"
	despatchLine.LineStatusCode = "NoStatus"
	despatchLine.DeliveredQuantity = float64(2)
	despatchLine.BackorderQuantity = float64(1)
	despatchLine.ItemId = uint32(7)
	despatchLine.UserId = "auth0|673c75d516e8adb9e6ffc892"
	despatchLine.UserEmail = "sprov300@gmail.com"
	despatchLine.RequestId = "bks1m1g91jau4nkks2f0"

	despatchLines := []*logisticsproto.CreateDespatchLineRequest{}
	despatchLines = append(despatchLines, &despatchLine)
	despatchHeader.DespatchLines = despatchLines
	type args struct {
		ctx context.Context
		in  *logisticsproto.CreateDespatchHeaderRequest
	}
	tests := []struct {
		ds      *DespatchService
		args    args
		wantErr bool
	}{
		{
			ds: despatchHeaderService,
			args: args{
				ctx: ctx,
				in:  &despatchHeader,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		despatchHeaderResp, err := tt.ds.CreateDespatchHeader(tt.args.ctx, tt.args.in)
		if (err != nil) != tt.wantErr {
			t.Errorf("DespatchService.CreateDespatchHeader() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		assert.NotNil(t, despatchHeaderResp)
		despatchHeaderResult := despatchHeaderResp.DespatchHeader
		assert.Equal(t, despatchHeaderResult.DespatchHeaderD.DocumentStatusCode, "NoStatus", "they should be equal")
		assert.Equal(t, despatchHeaderResult.DespatchHeaderD.DespatchAdviceTypeCode, "delivery", "they should be equal")
		assert.Equal(t, despatchHeaderResult.DespatchHeaderD.Note, "sample", "they should be equal")
	}
}

func GetDespatchHeader(id uint32, uuid4 []byte, idS string, issueDate string, documentStatusCode string, despatchAdviceTypeCode string, note string, orderId uint32, createdAt string, updatedAt string, createdByUserId string, updatedByUserId string) (*logisticsproto.DespatchHeader, error) {
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

	despatchHeaderD := logisticsproto.DespatchHeaderD{}
	despatchHeaderD.Id = id
	despatchHeaderD.Uuid4 = uuid4
	despatchHeaderD.IdS = idS
	despatchHeaderD.DocumentStatusCode = documentStatusCode
	despatchHeaderD.DespatchAdviceTypeCode = despatchAdviceTypeCode
	despatchHeaderD.Note = note
	despatchHeaderD.OrderId = orderId

	despatchHeaderT := logisticsproto.DespatchHeaderT{}
	despatchHeaderT.IssueDate = issueDate1

	crUpdTime := new(commonproto.CrUpdTime)
	crUpdTime.CreatedAt = createdAt1
	crUpdTime.UpdatedAt = updatedAt1

	crUpdUser := new(commonproto.CrUpdUser)
	crUpdUser.StatusCode = "active"
	crUpdUser.CreatedByUserId = createdByUserId
	crUpdUser.UpdatedByUserId = updatedByUserId

	despatchHeader := logisticsproto.DespatchHeader{DespatchHeaderD: &despatchHeaderD, DespatchHeaderT: &despatchHeaderT, CrUpdUser: crUpdUser, CrUpdTime: crUpdTime}

	return &despatchHeader, nil
}

func TestDespatchService_UpdateDespatchHeader(t *testing.T) {
	err := test.LoadSQL(logUser, dbService)
	if err != nil {
		t.Error(err)
		return
	}

	ctx := LoginUser()

	despatchHeaderService := NewDespatchService(log, dbService, redisService, userServiceClient)

	form := logisticsproto.UpdateDespatchHeaderRequest{}
	form.DocumentStatusCode = "Status"
	form.DespatchAdviceTypeCode = "delivery"
	form.Note = "sample1"
	form.Id = "234fd566-9451-4e3e-8318-4b713e688960"
	form.UserId = "auth0|673c75d516e8adb9e6ffc892"
	form.UserEmail = "sprov300@gmail.com"
	form.RequestId = "bks1m1g91jau4nkks2f0"

	updateResponse := logisticsproto.UpdateDespatchHeaderResponse{}

	type args struct {
		ctx context.Context
		in  *logisticsproto.UpdateDespatchHeaderRequest
	}
	tests := []struct {
		ds      *DespatchService
		args    args
		want    *logisticsproto.UpdateDespatchHeaderResponse
		wantErr bool
	}{
		{
			ds: despatchHeaderService,
			args: args{
				ctx: ctx,
				in:  &form,
			},
			want:    &updateResponse,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		got, err := tt.ds.UpdateDespatchHeader(tt.args.ctx, tt.args.in)
		if (err != nil) != tt.wantErr {
			t.Errorf("DespatchService.UpdateDespatchHeader() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("DespatchService.UpdateDespatchHeader() = %v, want %v", got, tt.want)
		}
	}
}
