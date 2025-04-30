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

func TestReceiptAdviceHeaderService_GetReceiptAdviceHeaders(t *testing.T) {
	err := test.LoadSQL(logUser, dbService)
	if err != nil {
		t.Error(err)
		return
	}

	ctx := LoginUser()

	receiptAdviceHeaderService := NewReceiptAdviceHeaderService(log, dbService, redisService, userServiceClient)
	receiptAdviceHeaders := []*logisticsproto.ReceiptAdviceHeader{}

	receiptAdviceHeader, err := GetReceiptAdviceHeader(uint32(1), []byte{35, 79, 213, 102, 148, 81, 78, 62, 131, 24, 75, 113, 62, 104, 137, 96}, "234fd566-9451-4e3e-8318-4b713e688960", "2019-07-23T10:04:26Z", "sample", uint32(1), uint32(1), uint32(2), uint32(1), uint32(1), "2019-07-23T10:04:26Z", "2019-07-23T10:04:26Z", "auth0|673c75d516e8adb9e6ffc892", "auth0|673c75d516e8adb9e6ffc892")
	if err != nil {
		t.Error(err)
		return
	}

	receiptAdviceHeaders = append(receiptAdviceHeaders, receiptAdviceHeader)

	form := logisticsproto.GetReceiptAdviceHeadersRequest{}
	form.Limit = "8"
	form.NextCursor = ""
	form.UserEmail = "sprov300@gmail.com"
	form.RequestId = "bks1m1g91jau4nkks2f0"

	nextc := "MA=="
	receiptAdviceHeadersResponse := logisticsproto.GetReceiptAdviceHeadersResponse{ReceiptAdviceHeaders: receiptAdviceHeaders, NextCursor: nextc}
	type args struct {
		ctx context.Context
		in  *logisticsproto.GetReceiptAdviceHeadersRequest
	}
	tests := []struct {
		rs      *ReceiptAdviceHeaderService
		args    args
		want    *logisticsproto.GetReceiptAdviceHeadersResponse
		wantErr bool
	}{
		{
			rs: receiptAdviceHeaderService,
			args: args{
				ctx: ctx,
				in:  &form,
			},
			want:    &receiptAdviceHeadersResponse,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		receiptAdviceHeadersResp, err := tt.rs.GetReceiptAdviceHeaders(tt.args.ctx, tt.args.in)
		if (err != nil) != tt.wantErr {
			t.Errorf("ReceiptAdviceHeaderService.GetReceiptAdviceHeaders() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(receiptAdviceHeadersResp, tt.want) {
			t.Errorf("ReceiptAdviceHeaderService.GetReceiptAdviceHeaders() = %v, want %v", receiptAdviceHeadersResp, tt.want)
		}
		assert.NotNil(t, receiptAdviceHeadersResp)
		receiptAdviceHeaderResult := receiptAdviceHeadersResp.ReceiptAdviceHeaders[0]
		assert.Equal(t, receiptAdviceHeaderResult.ReceiptAdviceHeaderD.OrderId, uint32(1), "they should be equal")
		assert.Equal(t, receiptAdviceHeaderResult.ReceiptAdviceHeaderD.DespatchId, uint32(1), "they should be equal")
		assert.Equal(t, receiptAdviceHeaderResult.ReceiptAdviceHeaderD.Note, "sample", "they should be equal")
	}
}

func TestReceiptAdviceHeaderService_GetReceiptAdviceHeader(t *testing.T) {
	err := test.LoadSQL(logUser, dbService)
	if err != nil {
		t.Error(err)
		return
	}

	ctx := LoginUser()

	receiptAdviceHeaderService := NewReceiptAdviceHeaderService(log, dbService, redisService, userServiceClient)

	receiptAdviceHeader, err := GetReceiptAdviceHeader(uint32(1), []byte{35, 79, 213, 102, 148, 81, 78, 62, 131, 24, 75, 113, 62, 104, 137, 96}, "234fd566-9451-4e3e-8318-4b713e688960", "2019-07-23T10:04:26Z", "sample", uint32(1), uint32(1), uint32(2), uint32(1), uint32(1), "2019-07-23T10:04:26Z", "2019-07-23T10:04:26Z", "auth0|673c75d516e8adb9e6ffc892", "auth0|673c75d516e8adb9e6ffc892")
	if err != nil {
		t.Error(err)
		return
	}
	receiptAdviceHeaderResponse := logisticsproto.GetReceiptAdviceHeaderResponse{}
	receiptAdviceHeaderResponse.ReceiptAdviceHeader = receiptAdviceHeader

	gform := commonproto.GetRequest{}
	gform.Id = "234fd566-9451-4e3e-8318-4b713e688960"
	gform.UserEmail = "sprov300@gmail.com"
	gform.RequestId = "bks1m1g91jau4nkks2f0"

	form := logisticsproto.GetReceiptAdviceHeaderRequest{}
	form.GetRequest = &gform

	type args struct {
		ctx context.Context
		in  *logisticsproto.GetReceiptAdviceHeaderRequest
	}
	tests := []struct {
		rs      *ReceiptAdviceHeaderService
		args    args
		want    *logisticsproto.GetReceiptAdviceHeaderResponse
		wantErr bool
	}{
		{
			rs: receiptAdviceHeaderService,
			args: args{
				ctx: ctx,
				in:  &form,
			},
			want:    &receiptAdviceHeaderResponse,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		receiptAdviceHeaderResp, err := tt.rs.GetReceiptAdviceHeader(tt.args.ctx, tt.args.in)
		if (err != nil) != tt.wantErr {
			t.Errorf("ReceiptAdviceHeaderService.GetReceiptAdviceHeader() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(receiptAdviceHeaderResp, tt.want) {
			t.Errorf("ReceiptAdviceHeaderService.GetReceiptAdviceHeader() = %v, want %v", receiptAdviceHeaderResp, tt.want)
		}
		assert.NotNil(t, receiptAdviceHeaderResp)
		receiptAdviceHeaderResult := receiptAdviceHeaderResp.ReceiptAdviceHeader
		assert.Equal(t, receiptAdviceHeaderResult.ReceiptAdviceHeaderD.OrderId, uint32(1), "they should be equal")
		assert.Equal(t, receiptAdviceHeaderResult.ReceiptAdviceHeaderD.DespatchId, uint32(1), "they should be equal")
		assert.Equal(t, receiptAdviceHeaderResult.ReceiptAdviceHeaderD.Note, "sample", "they should be equal")
	}
}

func TestReceiptAdviceHeaderService_GetReceiptAdviceHeaderByPk(t *testing.T) {
	err := test.LoadSQL(logUser, dbService)
	if err != nil {
		t.Error(err)
		return
	}

	ctx := LoginUser()

	receiptAdviceHeaderService := NewReceiptAdviceHeaderService(log, dbService, redisService, userServiceClient)

	receiptAdviceHeader, err := GetReceiptAdviceHeader(uint32(1), []byte{35, 79, 213, 102, 148, 81, 78, 62, 131, 24, 75, 113, 62, 104, 137, 96}, "234fd566-9451-4e3e-8318-4b713e688960", "2019-07-23T10:04:26Z", "sample", uint32(1), uint32(1), uint32(2), uint32(1), uint32(1), "2019-07-23T10:04:26Z", "2019-07-23T10:04:26Z", "auth0|673c75d516e8adb9e6ffc892", "auth0|673c75d516e8adb9e6ffc892")
	if err != nil {
		t.Error(err)
		return
	}

	receiptAdviceHeaderResponse := logisticsproto.GetReceiptAdviceHeaderByPkResponse{}
	receiptAdviceHeaderResponse.ReceiptAdviceHeader = receiptAdviceHeader

	gform := commonproto.GetByIdRequest{}
	gform.Id = uint32(1)
	gform.UserEmail = "sprov300@gmail.com"
	gform.RequestId = "bks1m1g91jau4nkks2f0"

	form := logisticsproto.GetReceiptAdviceHeaderByPkRequest{}
	form.GetByIdRequest = &gform

	type args struct {
		ctx context.Context
		in  *logisticsproto.GetReceiptAdviceHeaderByPkRequest
	}
	tests := []struct {
		rs      *ReceiptAdviceHeaderService
		args    args
		want    *logisticsproto.GetReceiptAdviceHeaderByPkResponse
		wantErr bool
	}{
		{
			rs: receiptAdviceHeaderService,
			args: args{
				ctx: ctx,
				in:  &form,
			},
			want:    &receiptAdviceHeaderResponse,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		receiptAdviceHeaderResp, err := tt.rs.GetReceiptAdviceHeaderByPk(tt.args.ctx, tt.args.in)
		if (err != nil) != tt.wantErr {
			t.Errorf("ReceiptAdviceHeaderService.GetReceiptAdviceHeaderByPk() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(receiptAdviceHeaderResp, tt.want) {
			t.Errorf("ReceiptAdviceHeaderService.GetReceiptAdviceHeaderByPk() = %v, want %v", receiptAdviceHeaderResp, tt.want)
		}
		assert.NotNil(t, receiptAdviceHeaderResp)
		receiptAdviceHeaderResult := receiptAdviceHeaderResp.ReceiptAdviceHeader
		assert.Equal(t, receiptAdviceHeaderResult.ReceiptAdviceHeaderD.OrderId, uint32(1), "they should be equal")
		assert.Equal(t, receiptAdviceHeaderResult.ReceiptAdviceHeaderD.DespatchId, uint32(1), "they should be equal")
		assert.Equal(t, receiptAdviceHeaderResult.ReceiptAdviceHeaderD.Note, "sample", "they should be equal")
	}
}

func TestReceiptAdviceHeaderService_CreateReceiptAdviceHeader(t *testing.T) {
	err := test.LoadSQL(logUser, dbService)
	if err != nil {
		t.Error(err)
		return
	}

	ctx := LoginUser()

	receiptAdviceHeaderService := NewReceiptAdviceHeaderService(log, dbService, redisService, userServiceClient)

	receiptAdviceHeader := logisticsproto.CreateReceiptAdviceHeaderRequest{}
	receiptAdviceHeader.IssueDate = "06/20/2022"
	receiptAdviceHeader.ReceiptAdviceTypeCode = "ABCFES"
	receiptAdviceHeader.Note = "sample"
	receiptAdviceHeader.OrderId = uint32(1)
	receiptAdviceHeader.DespatchId = uint32(1)
	receiptAdviceHeader.UserId = "auth0|673c75d516e8adb9e6ffc892"
	receiptAdviceHeader.UserEmail = "sprov300@gmail.com"
	receiptAdviceHeader.RequestId = "bks1m1g91jau4nkks2f0"

	receiptAdviceLine := logisticsproto.CreateReceiptAdviceLineRequest{}
	receiptAdviceLine.Note = "Mrs Green agreed to waive charge"
	receiptAdviceLine.ReceivedQuantity = uint32(2)
	receiptAdviceLine.ShortQuantity = uint32(1)
	receiptAdviceLine.ReceivedDate = "06/25/2022"
	receiptAdviceLine.OrderLineId = uint32(1)
	receiptAdviceLine.DespatchLineId = uint32(1)
	receiptAdviceLine.ItemId = uint32(7)
	receiptAdviceLine.ShipmentId = uint32(1)
	receiptAdviceLine.UserId = "auth0|673c75d516e8adb9e6ffc892"
	receiptAdviceLine.UserEmail = "sprov300@gmail.com"
	receiptAdviceLine.RequestId = "bks1m1g91jau4nkks2f0"

	receiptAdviceLines := []*logisticsproto.CreateReceiptAdviceLineRequest{}
	receiptAdviceLines = append(receiptAdviceLines, &receiptAdviceLine)
	receiptAdviceHeader.ReceiptAdviceLines = receiptAdviceLines
	type args struct {
		ctx context.Context
		in  *logisticsproto.CreateReceiptAdviceHeaderRequest
	}
	tests := []struct {
		rs      *ReceiptAdviceHeaderService
		args    args
		wantErr bool
	}{
		{
			rs: receiptAdviceHeaderService,
			args: args{
				ctx: ctx,
				in:  &receiptAdviceHeader,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		receiptAdviceHeadersResp, err := tt.rs.CreateReceiptAdviceHeader(tt.args.ctx, tt.args.in)
		if (err != nil) != tt.wantErr {
			t.Errorf("ReceiptAdviceHeaderService.CreateReceiptAdviceHeader() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		assert.NotNil(t, receiptAdviceHeadersResp)
		receiptAdviceHeaderResult := receiptAdviceHeadersResp.ReceiptAdviceHeader
		assert.Equal(t, receiptAdviceHeaderResult.ReceiptAdviceHeaderD.OrderId, uint32(1), "they should be equal")
		assert.Equal(t, receiptAdviceHeaderResult.ReceiptAdviceHeaderD.DespatchId, uint32(1), "they should be equal")
		assert.Equal(t, receiptAdviceHeaderResult.ReceiptAdviceHeaderD.Note, "sample", "they should be equal")

	}
}

func GetReceiptAdviceHeader(id uint32, uuid4 []byte, idS string, issueDate string, note string, orderId uint32, despatchId uint32, deliveryCustomerPartyId uint32, despatchSupplierPartyId uint32, shipmentId uint32, createdAt string, updatedAt string, createdByUserId string, updatedByUserId string) (*logisticsproto.ReceiptAdviceHeader, error) {
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
	receiptAdviceHeaderD := logisticsproto.ReceiptAdviceHeaderD{}
	receiptAdviceHeaderD.Id = id
	receiptAdviceHeaderD.Uuid4 = uuid4
	receiptAdviceHeaderD.IdS = idS
	receiptAdviceHeaderD.Note = note
	receiptAdviceHeaderD.OrderId = orderId
	receiptAdviceHeaderD.DespatchId = despatchId
	receiptAdviceHeaderD.DeliveryCustomerPartyId = deliveryCustomerPartyId
	receiptAdviceHeaderD.DespatchSupplierPartyId = despatchSupplierPartyId
	receiptAdviceHeaderD.ShipmentId = shipmentId

	receiptAdviceHeaderT := logisticsproto.ReceiptAdviceHeaderT{}
	receiptAdviceHeaderT.IssueDate = issueDate1

	crUpdTime := new(commonproto.CrUpdTime)
	crUpdTime.CreatedAt = createdAt1
	crUpdTime.UpdatedAt = updatedAt1

	crUpdUser := new(commonproto.CrUpdUser)
	crUpdUser.StatusCode = "active"
	crUpdUser.CreatedByUserId = createdByUserId
	crUpdUser.UpdatedByUserId = updatedByUserId

	receiptAdviceHeader := logisticsproto.ReceiptAdviceHeader{ReceiptAdviceHeaderD: &receiptAdviceHeaderD, ReceiptAdviceHeaderT: &receiptAdviceHeaderT, CrUpdUser: crUpdUser, CrUpdTime: crUpdTime}

	return &receiptAdviceHeader, nil
}

func GetReceiptAdviceLine(id uint32, uuid4 []byte, idS string, note string, receivedQuantity uint32, shortQuantity uint32, receivedDate string, orderLineId uint32, despatchLineId uint32, itemId uint32, shipmentId uint32, receiptAdviceHeaderId uint32, createdAt string, updatedAt string, createdByUserId string, updatedByUserId string) (*logisticsproto.ReceiptAdviceLine, error) {
	createdAt1, err := common.ConvertTimeToTimestamp(Layout, createdAt)
	if err != nil {
		log.Error("Error", zap.Error(err))
	}

	updatedAt1, err := common.ConvertTimeToTimestamp(Layout, updatedAt)
	if err != nil {
		log.Error("Error", zap.Error(err))
	}

	receivedDate1, err := common.ConvertTimeToTimestamp(Layout, receivedDate)
	if err != nil {
		log.Error("Error", zap.Error(err))
	}

	receiptAdviceLineD := logisticsproto.ReceiptAdviceLineD{}
	receiptAdviceLineD.Id = id
	receiptAdviceLineD.Uuid4 = uuid4
	receiptAdviceLineD.IdS = idS
	receiptAdviceLineD.Note = note
	receiptAdviceLineD.ReceivedQuantity = receivedQuantity
	receiptAdviceLineD.ShortQuantity = shortQuantity
	receiptAdviceLineD.OrderLineId = orderLineId
	receiptAdviceLineD.DespatchLineId = despatchLineId
	receiptAdviceLineD.ItemId = itemId
	receiptAdviceLineD.ShipmentId = shipmentId
	receiptAdviceLineD.ReceiptAdviceHeaderId = receiptAdviceHeaderId

	receiptAdviceLineT := logisticsproto.ReceiptAdviceLineT{}
	receiptAdviceLineT.ReceivedDate = receivedDate1

	crUpdTime := new(commonproto.CrUpdTime)
	crUpdTime.CreatedAt = createdAt1
	crUpdTime.UpdatedAt = updatedAt1

	crUpdUser := new(commonproto.CrUpdUser)
	crUpdUser.StatusCode = "active"
	crUpdUser.CreatedByUserId = createdByUserId
	crUpdUser.UpdatedByUserId = updatedByUserId

	receiptAdviceLine := logisticsproto.ReceiptAdviceLine{ReceiptAdviceLineD: &receiptAdviceLineD, ReceiptAdviceLineT: &receiptAdviceLineT, CrUpdUser: crUpdUser, CrUpdTime: crUpdTime}

	return &receiptAdviceLine, nil
}

func TestReceiptAdviceHeaderService_GetReceiptAdviceLines(t *testing.T) {
	err := test.LoadSQL(logUser, dbService)
	if err != nil {
		t.Error(err)
		return
	}

	ctx := LoginUser()

	receiptAdviceHeaderService := NewReceiptAdviceHeaderService(log, dbService, redisService, userServiceClient)

	receiptAdviceLines := []*logisticsproto.ReceiptAdviceLine{}
	receiptAdviceLine, err := GetReceiptAdviceLine(uint32(1), []byte{122, 58, 114, 57, 52, 23, 72, 139, 137, 204, 176, 161, 133, 206, 29, 25}, "7a3a7239-3417-488b-89cc-b0a185ce1d19", "SAMPLE", uint32(90), uint32(10), "2019-07-23T10:04:26Z", uint32(1), uint32(1), uint32(1), uint32(1), uint32(1), "2019-07-23T10:04:26Z", "2019-07-23T10:04:26Z", "auth0|673c75d516e8adb9e6ffc892", "auth0|673c75d516e8adb9e6ffc892")
	if err != nil {
		t.Error(err)
		return
	}

	receiptAdviceLines = append(receiptAdviceLines, receiptAdviceLine)

	receiptAdviceLinesResponse := logisticsproto.GetReceiptAdviceLinesResponse{}
	receiptAdviceLinesResponse.ReceiptAdviceLines = receiptAdviceLines

	gform := commonproto.GetRequest{}
	gform.Id = "234fd566-9451-4e3e-8318-4b713e688960"
	gform.UserEmail = "sprov300@gmail.com"
	gform.RequestId = "bks1m1g91jau4nkks2f0"

	form := logisticsproto.GetReceiptAdviceLinesRequest{}
	form.GetRequest = &gform

	type args struct {
		ctx context.Context
		in  *logisticsproto.GetReceiptAdviceLinesRequest
	}
	tests := []struct {
		rs      *ReceiptAdviceHeaderService
		args    args
		want    *logisticsproto.GetReceiptAdviceLinesResponse
		wantErr bool
	}{
		{
			rs: receiptAdviceHeaderService,
			args: args{
				ctx: ctx,
				in:  &form,
			},
			want:    &receiptAdviceLinesResponse,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		receiptAdviceLinesResp, err := tt.rs.GetReceiptAdviceLines(tt.args.ctx, tt.args.in)
		if (err != nil) != tt.wantErr {
			t.Errorf("ReceiptAdviceHeaderService.GetReceiptAdviceLines() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(receiptAdviceLinesResp, tt.want) {
			t.Errorf("ReceiptAdviceHeaderService.GetReceiptAdviceLines() = %v, want %v", receiptAdviceLinesResp, tt.want)
		}
		assert.NotNil(t, receiptAdviceLinesResp)
		receiptAdviceLineResult := receiptAdviceLinesResp.ReceiptAdviceLines[0]
		assert.Equal(t, receiptAdviceLineResult.ReceiptAdviceLineD.Note, "SAMPLE", "they should be equal")
		assert.Equal(t, receiptAdviceLineResult.ReceiptAdviceLineD.ReceivedQuantity, uint32(90), "they should be equal")
		assert.Equal(t, receiptAdviceLineResult.ReceiptAdviceLineD.ShortQuantity, uint32(10), "they should be equal")
	}
}

func TestReceiptAdviceHeaderService_UpdateReceiptAdviceHeader(t *testing.T) {
	err := test.LoadSQL(logUser, dbService)
	if err != nil {
		t.Error(err)
		return
	}

	ctx := LoginUser()

	receiptAdviceHeaderService := NewReceiptAdviceHeaderService(log, dbService, redisService, userServiceClient)

	form := logisticsproto.UpdateReceiptAdviceHeaderRequest{}
	form.ReceiptAdviceTypeCode = "DFESAB"
	form.Note = "sample1"
	form.LineCountNumeric = uint32(10)
	form.Id = "234fd566-9451-4e3e-8318-4b713e688960"
	form.UserId = "auth0|673c75d516e8adb9e6ffc892"
	form.UserEmail = "sprov300@gmail.com"
	form.RequestId = "bks1m1g91jau4nkks2f0"

	updateResponse := logisticsproto.UpdateReceiptAdviceHeaderResponse{}

	type args struct {
		ctx context.Context
		in  *logisticsproto.UpdateReceiptAdviceHeaderRequest
	}
	tests := []struct {
		rs      *ReceiptAdviceHeaderService
		args    args
		want    *logisticsproto.UpdateReceiptAdviceHeaderResponse
		wantErr bool
	}{
		{
			rs: receiptAdviceHeaderService,
			args: args{
				ctx: ctx,
				in:  &form,
			},
			want:    &updateResponse,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		got, err := tt.rs.UpdateReceiptAdviceHeader(tt.args.ctx, tt.args.in)
		if (err != nil) != tt.wantErr {
			t.Errorf("ReceiptAdviceHeaderService.UpdateReceiptAdviceHeader() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("ReceiptAdviceHeaderService.UpdateReceiptAdviceHeader() = %v, want %v", got, tt.want)
		}
	}
}
