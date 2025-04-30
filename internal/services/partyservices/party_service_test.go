package partyservices

import (
	"context"
	"reflect"
	"testing"

	"github.com/cloudfresco/sc-ubl/internal/common"
	commonproto "github.com/cloudfresco/sc-ubl/internal/protogen/common/v1"
	partyproto "github.com/cloudfresco/sc-ubl/internal/protogen/party/v1"
	"github.com/cloudfresco/sc-ubl/test"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestPartyService_GetParties(t *testing.T) {
	err := test.LoadSQL(logUser, dbService)
	if err != nil {
		t.Error(err)
		return
	}

	ctx := LoginUser()

	partyService := NewPartyService(log, dbService, redisService, userServiceClient)
	parties := []*partyproto.Party{}

	party1, err := GetParty(uint32(1), []byte{204, 171, 244, 233, 201, 146, 76, 145, 176, 8, 228, 162, 97, 56, 221, 28}, "ccabf4e9-c992-4c91-b008-e4a26138dd1c", "Consortial", uint32(1), true, uint32(0), uint32(0), uint32(1), "", "", "2023-03-24T10:04:26Z", "2024-07-25T10:04:26Z", "active", "2019-07-23T10:04:26Z", "2019-07-23T10:04:26Z", "auth0|673c75d516e8adb9e6ffc892", "auth0|673c75d516e8adb9e6ffc892")
	if err != nil {
		t.Error(err)
		return
	}

	party2, err := GetParty(uint32(2), []byte{187, 161, 132, 9, 243, 71, 68, 255, 138, 42, 243, 222, 249, 139, 227, 204}, "bba18409-f347-44ff-8a2a-f3def98be3cc", "IYT Corporation", uint32(0), false, uint32(0), uint32(0), uint32(2), "", "", "2023-03-24T10:04:26Z", "2024-07-25T10:04:26Z", "active", "2019-07-23T10:04:26Z", "2019-07-23T10:04:26Z", "auth0|673c75d516e8adb9e6ffc892", "auth0|673c75d516e8adb9e6ffc892")
	if err != nil {
		t.Error(err)
		return
	}

	party3, err := GetParty(uint32(3), []byte{50, 20, 74, 126, 243, 122, 65, 175, 141, 246, 206, 190, 64, 216, 130, 82}, "32144a7e-f37a-41af-8df6-cebe40d88252", "Salescompany ltd", uint32(0), false, uint32(0), uint32(0), uint32(3), "The Sellercompany Incorporated", "5402697509", "2023-03-24T10:04:26Z", "2024-07-25T10:04:26Z", "active", "2019-07-23T10:04:26Z", "2019-07-23T10:04:26Z", "auth0|673c75d516e8adb9e6ffc892", "auth0|673c75d516e8adb9e6ffc892")
	if err != nil {
		t.Error(err)
		return
	}

	party4, err := GetParty(uint32(4), []byte{247, 30, 240, 166, 216, 120, 79, 47, 175, 126, 228, 36, 254, 110, 244, 208}, "f71ef0a6-d878-4f2f-af7e-e424fe6ef4d0", "Buyercompany ltd", uint32(0), false, uint32(0), uint32(0), uint32(4), "The buyercompany inc.", "5645342123", "2023-03-24T10:04:26Z", "2024-07-25T10:04:26Z", "active", "2019-07-23T10:04:26Z", "2019-07-23T10:04:26Z", "auth0|673c75d516e8adb9e6ffc892", "auth0|673c75d516e8adb9e6ffc892")
	if err != nil {
		t.Error(err)
		return
	}

	party5, err := GetParty(uint32(5), []byte{247, 32, 104, 111, 133, 211, 67, 90, 136, 101, 6, 224, 58, 88, 24, 169}, "f720686f-85d3-435a-8865-06e03a5818a9", "Ebeneser Scrooge Inc", uint32(0), false, uint32(0), uint32(0), uint32(5), "Ebeneser Scrooge Inc.", "6411982340", "2023-03-24T10:04:26Z", "2024-07-25T10:04:26Z", "active", "2019-07-23T10:04:26Z", "2019-07-23T10:04:26Z", "auth0|673c75d516e8adb9e6ffc892", "auth0|673c75d516e8adb9e6ffc892")
	if err != nil {
		t.Error(err)
		return
	}

	party6, err := GetParty(uint32(6), []byte{151, 110, 134, 181, 169, 89, 71, 23, 160, 10, 37, 167, 23, 241, 85, 195}, "976e86b5-a959-4717-a00a-25a717f155c3", "Test supplier", uint32(0), false, uint32(0), uint32(0), uint32(6), "", "", "2023-03-24T10:04:26Z", "2024-07-25T10:04:26Z", "active", "2019-07-23T10:04:26Z", "2019-07-23T10:04:26Z", "auth0|673c75d516e8adb9e6ffc892", "auth0|673c75d516e8adb9e6ffc892")
	if err != nil {
		t.Error(err)
		return
	}

	parties = append(parties, party6, party5, party4, party3, party2, party1)

	nextc := "MA=="
	partiesResponse := partyproto.GetPartiesResponse{Parties: parties, NextCursor: nextc}

	form := partyproto.GetPartiesRequest{}
	form.Limit = "8"
	form.NextCursor = ""
	form.UserEmail = "sprov300@gmail.com"
	form.RequestId = "bks1m1g91jau4nkks2f0"

	type args struct {
		ctx context.Context
		in  *partyproto.GetPartiesRequest
	}
	tests := []struct {
		ps      *PartyService
		args    args
		want    *partyproto.GetPartiesResponse
		wantErr bool
	}{
		{
			ps: partyService,
			args: args{
				ctx: ctx,
				in:  &form,
			},
			want:    &partiesResponse,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		partiesResp, err := tt.ps.GetParties(tt.args.ctx, tt.args.in)
		if (err != nil) != tt.wantErr {
			t.Errorf("PartyService.GetParties() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(partiesResp, tt.want) {
			t.Errorf("PartyService.GetParties() = %v, want %v", partiesResp, tt.want)
		}
		assert.NotNil(t, partiesResp)
		partyResult := partiesResp.Parties[1]
		assert.Equal(t, partyResult.PartyD.PartyName, "Ebeneser Scrooge Inc", "they should be equal")
		assert.Equal(t, partyResult.PartyLegalEntityD.CompanyId, "6411982340", "they should be equal")
		assert.Equal(t, partyResult.PartyLegalEntityD.RegistrationName, "Ebeneser Scrooge Inc.", "they should be equal")
	}
}

func TestPartyService_GetParty(t *testing.T) {
	err := test.LoadSQL(logUser, dbService)
	if err != nil {
		t.Error(err)
		return
	}

	ctx := LoginUser()

	partyService := NewPartyService(log, dbService, redisService, userServiceClient)
	party, err := GetParty(uint32(1), []byte{204, 171, 244, 233, 201, 146, 76, 145, 176, 8, 228, 162, 97, 56, 221, 28}, "ccabf4e9-c992-4c91-b008-e4a26138dd1c", "Consortial", uint32(1), true, uint32(0), uint32(0), uint32(1), "", "", "2023-03-24T10:04:26Z", "2024-07-25T10:04:26Z", "active", "2019-07-23T10:04:26Z", "2019-07-23T10:04:26Z", "auth0|673c75d516e8adb9e6ffc892", "auth0|673c75d516e8adb9e6ffc892")
	if err != nil {
		t.Error(err)
		return
	}

	partyResponse := partyproto.GetPartyResponse{}
	partyResponse.Party = party

	form := partyproto.GetPartyRequest{}
	gform := commonproto.GetRequest{}
	gform.Id = "ccabf4e9-c992-4c91-b008-e4a26138dd1c"
	gform.UserEmail = "sprov300@gmail.com"
	gform.RequestId = "bks1m1g91jau4nkks2f0"
	form.GetRequest = &gform

	type args struct {
		ctx context.Context
		in  *partyproto.GetPartyRequest
	}
	tests := []struct {
		ps      *PartyService
		args    args
		want    *partyproto.GetPartyResponse
		wantErr bool
	}{
		{
			ps: partyService,
			args: args{
				ctx: ctx,
				in:  &form,
			},
			want:    &partyResponse,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		partyResp, err := tt.ps.GetParty(tt.args.ctx, tt.args.in)
		if (err != nil) != tt.wantErr {
			t.Errorf("PartyService.GetParty() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(partyResp, tt.want) {
			t.Errorf("PartyService.GetParty() = %v, want %v", partyResp, tt.want)
		}
		assert.NotNil(t, partyResp)
		partyResult := partyResp.Party
		assert.Equal(t, partyResult.PartyD.PartyName, "Consortial", "they should be equal")
		assert.Equal(t, partyResult.PartyD.NumChd, uint32(1), "they should be equal")
		assert.True(t, partyResult.PartyD.Leaf, "Its true")
	}
}

func TestPartyService_GetPartyByPk(t *testing.T) {
	err := test.LoadSQL(logUser, dbService)
	if err != nil {
		t.Error(err)
		return
	}

	ctx := LoginUser()

	partyService := NewPartyService(log, dbService, redisService, userServiceClient)
	party, err := GetParty(uint32(1), []byte{204, 171, 244, 233, 201, 146, 76, 145, 176, 8, 228, 162, 97, 56, 221, 28}, "ccabf4e9-c992-4c91-b008-e4a26138dd1c", "Consortial", uint32(1), true, uint32(0), uint32(0), uint32(1), "", "", "2023-03-24T10:04:26Z", "2024-07-25T10:04:26Z", "active", "2019-07-23T10:04:26Z", "2019-07-23T10:04:26Z", "auth0|673c75d516e8adb9e6ffc892", "auth0|673c75d516e8adb9e6ffc892")
	if err != nil {
		t.Error(err)
		return
	}

	partyResponse := partyproto.GetPartyByPkResponse{}
	partyResponse.Party = party

	form := partyproto.GetPartyByPkRequest{}
	gform := commonproto.GetByIdRequest{}
	gform.Id = uint32(1)
	gform.UserEmail = "sprov300@gmail.com"
	gform.RequestId = "bks1m1g91jau4nkks2f0"
	form.GetByIdRequest = &gform

	type args struct {
		ctx context.Context
		in  *partyproto.GetPartyByPkRequest
	}
	tests := []struct {
		ps      *PartyService
		args    args
		want    *partyproto.GetPartyByPkResponse
		wantErr bool
	}{
		{
			ps: partyService,
			args: args{
				ctx: ctx,
				in:  &form,
			},
			want:    &partyResponse,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		partyResp, err := tt.ps.GetPartyByPk(tt.args.ctx, tt.args.in)
		if (err != nil) != tt.wantErr {
			t.Errorf("PartyService.GetPartyByPk() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(partyResp, tt.want) {
			t.Errorf("PartyService.GetPartyByPk() = %v, want %v", partyResp, tt.want)
		}
		assert.NotNil(t, partyResp)
		partyResult := partyResp.Party
		assert.Equal(t, partyResult.PartyD.PartyName, "Consortial", "they should be equal")
		assert.Equal(t, partyResult.PartyD.NumChd, uint32(1), "they should be equal")
		assert.True(t, partyResult.PartyD.Leaf, "Its true")
	}
}

func TestPartyService_GetTopLevelParties(t *testing.T) {
	err := test.LoadSQL(logUser, dbService)
	if err != nil {
		t.Error(err)
		return
	}

	ctx := LoginUser()

	partyService := NewPartyService(log, dbService, redisService, userServiceClient)

	party1, err := GetParty(uint32(1), []byte{204, 171, 244, 233, 201, 146, 76, 145, 176, 8, 228, 162, 97, 56, 221, 28}, "ccabf4e9-c992-4c91-b008-e4a26138dd1c", "Consortial", uint32(1), true, uint32(0), uint32(0), uint32(1), "", "", "2023-03-24T10:04:26Z", "2024-07-25T10:04:26Z", "active", "2019-07-23T10:04:26Z", "2019-07-23T10:04:26Z", "auth0|673c75d516e8adb9e6ffc892", "auth0|673c75d516e8adb9e6ffc892")
	if err != nil {
		t.Error(err)
		return
	}

	party2, err := GetParty(uint32(2), []byte{187, 161, 132, 9, 243, 71, 68, 255, 138, 42, 243, 222, 249, 139, 227, 204}, "bba18409-f347-44ff-8a2a-f3def98be3cc", "IYT Corporation", uint32(0), false, uint32(0), uint32(0), uint32(2), "", "", "2023-03-24T10:04:26Z", "2024-07-25T10:04:26Z", "active", "2019-07-23T10:04:26Z", "2019-07-23T10:04:26Z", "auth0|673c75d516e8adb9e6ffc892", "auth0|673c75d516e8adb9e6ffc892")
	if err != nil {
		t.Error(err)
		return
	}

	party3, err := GetParty(uint32(3), []byte{50, 20, 74, 126, 243, 122, 65, 175, 141, 246, 206, 190, 64, 216, 130, 82}, "32144a7e-f37a-41af-8df6-cebe40d88252", "Salescompany ltd", uint32(0), false, uint32(0), uint32(0), uint32(3), "The Sellercompany Incorporated", "5402697509", "2023-03-24T10:04:26Z", "2024-07-25T10:04:26Z", "active", "2019-07-23T10:04:26Z", "2019-07-23T10:04:26Z", "auth0|673c75d516e8adb9e6ffc892", "auth0|673c75d516e8adb9e6ffc892")
	if err != nil {
		t.Error(err)
		return
	}

	party4, err := GetParty(uint32(4), []byte{247, 30, 240, 166, 216, 120, 79, 47, 175, 126, 228, 36, 254, 110, 244, 208}, "f71ef0a6-d878-4f2f-af7e-e424fe6ef4d0", "Buyercompany ltd", uint32(0), false, uint32(0), uint32(0), uint32(4), "The buyercompany inc.", "5645342123", "2023-03-24T10:04:26Z", "2024-07-25T10:04:26Z", "active", "2019-07-23T10:04:26Z", "2019-07-23T10:04:26Z", "auth0|673c75d516e8adb9e6ffc892", "auth0|673c75d516e8adb9e6ffc892")
	if err != nil {
		t.Error(err)
		return
	}

	party5, err := GetParty(uint32(5), []byte{247, 32, 104, 111, 133, 211, 67, 90, 136, 101, 6, 224, 58, 88, 24, 169}, "f720686f-85d3-435a-8865-06e03a5818a9", "Ebeneser Scrooge Inc", uint32(0), false, uint32(0), uint32(0), uint32(5), "Ebeneser Scrooge Inc.", "6411982340", "2023-03-24T10:04:26Z", "2024-07-25T10:04:26Z", "active", "2019-07-23T10:04:26Z", "2019-07-23T10:04:26Z", "auth0|673c75d516e8adb9e6ffc892", "auth0|673c75d516e8adb9e6ffc892")
	if err != nil {
		t.Error(err)
		return
	}

	party6, err := GetParty(uint32(6), []byte{151, 110, 134, 181, 169, 89, 71, 23, 160, 10, 37, 167, 23, 241, 85, 195}, "976e86b5-a959-4717-a00a-25a717f155c3", "Test supplier", uint32(0), false, uint32(0), uint32(0), uint32(6), "", "", "2023-03-24T10:04:26Z", "2024-07-25T10:04:26Z", "active", "2019-07-23T10:04:26Z", "2019-07-23T10:04:26Z", "auth0|673c75d516e8adb9e6ffc892", "auth0|673c75d516e8adb9e6ffc892")
	if err != nil {
		t.Error(err)
		return
	}
	parties := []*partyproto.Party{}
	parties = append(parties, party1, party2, party3, party4, party5, party6)

	form := partyproto.GetTopLevelPartiesRequest{}
	form.UserEmail = "sprov300@gmail.com"
	form.RequestId = "bks1m1g91jau4nkks2f0"

	partiesResponse := partyproto.GetTopLevelPartiesResponse{Parties: parties}

	type args struct {
		ctx context.Context
		in  *partyproto.GetTopLevelPartiesRequest
	}
	tests := []struct {
		ps      *PartyService
		args    args
		want    *partyproto.GetTopLevelPartiesResponse
		wantErr bool
	}{
		{
			ps: partyService,
			args: args{
				ctx: ctx,
				in:  &form,
			},
			want:    &partiesResponse,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		partiesResp, err := tt.ps.GetTopLevelParties(tt.args.ctx, tt.args.in)
		if (err != nil) != tt.wantErr {
			t.Errorf("PartyService.GetTopLevelParties() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(partiesResp, tt.want) {
			t.Errorf("PartyService.GetTopLevelParties() = %v, want %v", partiesResp, tt.want)
		}
		assert.NotNil(t, partiesResp)
		partyResult := partiesResp.Parties[3]
		assert.Equal(t, partyResult.PartyD.PartyName, "Buyercompany ltd", "they should be equal")
		assert.Equal(t, partyResult.PartyLegalEntityD.CompanyId, "5645342123", "they should be equal")
		assert.Equal(t, partyResult.PartyLegalEntityD.RegistrationName, "The buyercompany inc.", "they should be equal")
	}
}

func TestPartyService_GetChildParties(t *testing.T) {
	err := test.LoadSQL(logUser, dbService)
	if err != nil {
		t.Error(err)
		return
	}

	ctx := LoginUser()

	partyService := NewPartyService(log, dbService, redisService, userServiceClient)

	party, err := GetParty(uint32(7), []byte{238, 12, 95, 123, 28, 62, 66, 135, 179, 175, 7, 180, 201, 112, 127, 155}, "ee0c5f7b-1c3e-4287-b3af-07b4c9707f9b", "Test customer", uint32(0), false, uint32(1), uint32(1), uint32(7), "", "", "2023-03-24T10:04:26Z", "2024-07-25T10:04:26Z", "active", "2019-07-23T10:04:26Z", "2019-07-23T10:04:26Z", "auth0|673c75d516e8adb9e6ffc892", "auth0|673c75d516e8adb9e6ffc892")
	if err != nil {
		t.Error(err)
		return
	}

	parties := []*partyproto.Party{}
	parties = append(parties, party)
	partiesResponse := partyproto.GetChildPartiesResponse{Parties: parties}

	form := partyproto.GetChildPartiesRequest{}
	gform := commonproto.GetRequest{}
	gform.Id = "ccabf4e9-c992-4c91-b008-e4a26138dd1c"
	gform.UserEmail = "sprov300@gmail.com"
	gform.RequestId = "bks1m1g91jau4nkks2f0"
	form.GetRequest = &gform

	type args struct {
		ctx context.Context
		in  *partyproto.GetChildPartiesRequest
	}
	tests := []struct {
		ps      *PartyService
		args    args
		want    *partyproto.GetChildPartiesResponse
		wantErr bool
	}{
		{
			ps: partyService,
			args: args{
				ctx: ctx,
				in:  &form,
			},
			want:    &partiesResponse,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		partiesResp, err := tt.ps.GetChildParties(tt.args.ctx, tt.args.in)
		if (err != nil) != tt.wantErr {
			t.Errorf("PartyService.GetChildParties() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(partiesResp, tt.want) {
			t.Errorf("PartyService.GetChildParties() = %v, want %v", partiesResp, tt.want)
		}
		assert.NotNil(t, partiesResp)
		partyResult := partiesResp.Parties[0]
		assert.Equal(t, partyResult.PartyD.PartyName, "Test customer", "they should be equal")
		assert.Equal(t, partyResult.PartyD.LevelP, uint32(1), "they should be equal")
		assert.Equal(t, partyResult.PartyD.NumChd, uint32(0), "they should be equal")
	}
}

func TestPartyService_GetParentParty(t *testing.T) {
	err := test.LoadSQL(logUser, dbService)
	if err != nil {
		t.Error(err)
		return
	}

	ctx := LoginUser()

	partyService := NewPartyService(log, dbService, redisService, userServiceClient)
	party, err := GetParty(uint32(1), []byte{204, 171, 244, 233, 201, 146, 76, 145, 176, 8, 228, 162, 97, 56, 221, 28}, "ccabf4e9-c992-4c91-b008-e4a26138dd1c", "Consortial", uint32(1), true, uint32(0), uint32(0), uint32(1), "", "", "2023-03-24T10:04:26Z", "2024-07-25T10:04:26Z", "active", "2019-07-23T10:04:26Z", "2019-07-23T10:04:26Z", "auth0|673c75d516e8adb9e6ffc892", "auth0|673c75d516e8adb9e6ffc892")
	if err != nil {
		t.Error(err)
		return
	}

	partyResponse := partyproto.GetParentPartyResponse{}
	partyResponse.Party = party

	form := partyproto.GetParentPartyRequest{}
	gform := commonproto.GetRequest{}
	gform.Id = "ee0c5f7b-1c3e-4287-b3af-07b4c9707f9b"
	gform.UserEmail = "sprov300@gmail.com"
	gform.RequestId = "bks1m1g91jau4nkks2f0"
	form.GetRequest = &gform

	type args struct {
		ctx context.Context
		in  *partyproto.GetParentPartyRequest
	}
	tests := []struct {
		ps      *PartyService
		args    args
		want    *partyproto.GetParentPartyResponse
		wantErr bool
	}{
		{
			ps: partyService,
			args: args{
				ctx: ctx,
				in:  &form,
			},
			want:    &partyResponse,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		partyResp, err := tt.ps.GetParentParty(tt.args.ctx, tt.args.in)
		if (err != nil) != tt.wantErr {
			t.Errorf("PartyService.GetParentParty() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(partyResp, tt.want) {
			t.Errorf("PartyService.GetParentParty() = %v, want %v", partyResp, tt.want)
		}
		assert.NotNil(t, partyResp)
		partyResult := partyResp.Party
		assert.Equal(t, partyResult.PartyD.PartyName, "Consortial", "they should be equal")
		assert.Equal(t, partyResult.PartyD.NumChd, uint32(1), "they should be equal")
		assert.True(t, partyResult.PartyD.Leaf, "Its true")

	}
}

func GetParty(id uint32, uuid4 []byte, idS string, partyName string, numChd uint32, leaf bool, levelP uint32, parentId uint32, addressId uint32, registrationName string, companyId string, registrationDate string, registrationExpirationDate string, statusCode string, createdAt string, updatedAt string, createdByUserId string, updatedByUserId string) (*partyproto.Party, error) {
	createdAt1, err := common.ConvertTimeToTimestamp(Layout, createdAt)
	if err != nil {
		log.Error("Error", zap.Error(err))
	}

	updatedAt1, err := common.ConvertTimeToTimestamp(Layout, updatedAt)
	if err != nil {
		log.Error("Error", zap.Error(err))
	}

	registrationDate1, err := common.ConvertTimeToTimestamp(Layout, registrationDate)
	if err != nil {
		log.Error("Error", zap.Error(err))
	}

	registrationExpirationDate1, err := common.ConvertTimeToTimestamp(Layout, registrationExpirationDate)
	if err != nil {
		log.Error("Error", zap.Error(err))
	}

	partyD := partyproto.PartyD{}
	partyD.Id = id
	partyD.Uuid4 = uuid4
	partyD.IdS = idS
	partyD.PartyName = partyName
	partyD.NumChd = numChd
	partyD.Leaf = leaf
	partyD.LevelP = levelP
	partyD.ParentId = parentId
	partyD.AddressId = addressId

	partyLegalEntityD := commonproto.PartyLegalEntityD{}
	partyLegalEntityD.RegistrationName = registrationName
	partyLegalEntityD.CompanyId = companyId

	partyLegalEntityT := commonproto.PartyLegalEntityT{}
	partyLegalEntityT.RegistrationDate = registrationDate1
	partyLegalEntityT.RegistrationExpirationDate = registrationExpirationDate1

	crUpdTime := new(commonproto.CrUpdTime)
	crUpdTime.CreatedAt = createdAt1
	crUpdTime.UpdatedAt = updatedAt1

	crUpdUser := new(commonproto.CrUpdUser)
	crUpdUser.StatusCode = "active"
	crUpdUser.CreatedByUserId = createdByUserId
	crUpdUser.UpdatedByUserId = updatedByUserId

	party := partyproto.Party{PartyD: &partyD, PartyLegalEntityD: &partyLegalEntityD, PartyLegalEntityT: &partyLegalEntityT, CrUpdUser: crUpdUser, CrUpdTime: crUpdTime}
	return &party, nil
}

func TestPartyService_UpdateParty(t *testing.T) {
	err := test.LoadSQL(logUser, dbService)
	if err != nil {
		t.Error(err)
		return
	}

	ctx := LoginUser()

	partyService := NewPartyService(log, dbService, redisService, userServiceClient)
	form := partyproto.UpdatePartyRequest{}
	form.PartyName = "Consortial New"
	form.PartyDesc = "Consortial Description"
	form.Id = "ccabf4e9-c992-4c91-b008-e4a26138dd1c"
	form.UserEmail = "sprov300@gmail.com"
	form.RequestId = "bks1m1g91jau4nkks2f0"

	partyResponse := partyproto.UpdatePartyResponse{}

	type args struct {
		ctx context.Context
		in  *partyproto.UpdatePartyRequest
	}
	tests := []struct {
		ps      *PartyService
		args    args
		want    *partyproto.UpdatePartyResponse
		wantErr bool
	}{
		{
			ps: partyService,
			args: args{
				ctx: ctx,
				in:  &form,
			},
			want:    &partyResponse,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		got, err := tt.ps.UpdateParty(tt.args.ctx, tt.args.in)
		if (err != nil) != tt.wantErr {
			t.Errorf("PartyService.UpdateParty() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("PartyService.UpdateParty() = %v, want %v", got, tt.want)
		}
	}
}

func TestPartyService_DeleteParty(t *testing.T) {
	err := test.LoadSQL(logUser, dbService)
	if err != nil {
		t.Error(err)
		return
	}

	ctx := LoginUser()

	partyService := NewPartyService(log, dbService, redisService, userServiceClient)

	form := partyproto.DeletePartyRequest{}
	gform := commonproto.GetRequest{}
	gform.Id = "ccabf4e9-c992-4c91-b008-e4a26138dd1c"
	gform.UserEmail = "sprov300@gmail.com"
	gform.RequestId = "bks1m1g91jau4nkks2f0"
	form.GetRequest = &gform

	partyResponse := partyproto.DeletePartyResponse{}

	type args struct {
		ctx context.Context
		in  *partyproto.DeletePartyRequest
	}
	tests := []struct {
		ps      *PartyService
		args    args
		want    *partyproto.DeletePartyResponse
		wantErr bool
	}{
		{
			ps: partyService,
			args: args{
				ctx: ctx,
				in:  &form,
			},
			want:    &partyResponse,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		got, err := tt.ps.DeleteParty(tt.args.ctx, tt.args.in)
		if (err != nil) != tt.wantErr {
			t.Errorf("PartyService.DeleteParty() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("PartyService.DeleteParty() = %v, want %v", got, tt.want)
		}
	}
}

func TestPartyService_CreateParty(t *testing.T) {
	err := test.LoadSQL(logUser, dbService)
	if err != nil {
		t.Error(err)
		return
	}

	ctx := LoginUser()

	partyService := NewPartyService(log, dbService, redisService, userServiceClient)

	party := partyproto.CreatePartyRequest{}
	party.PartyName = "Ebeneser Scrooge Inc"
	party.PublicKey = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6IkFzaW"
	party.RegistrationName = "Ebeneser Scrooge Inc."
	party.CompanyId = "6411982340"
	party.RegistrationDate = "03/24/2023"
	party.RegistrationExpirationDate = "07/25/2024"
	party.Name1 = "Ebeneser Scrooge Inc"
	party.AddrListAgencyId = "9"
	party.AddrListId = "1234567890123"
	party.AddrListVersionId = "GLN"
	party.Postbox = "5467"
	party.StreetName = "Main street"
	party.AdditionalStreetName = "Suite 123"
	party.BuildingName = "Back door"
	party.BuildingNumber = "1234"
	party.Department = "Revenue department"
	party.CityName = "Big city"
	party.PostalZone = "54321"
	party.CountryIdCode = "RegionA"
	party.CountryName = "GB"
	party.UserId = "auth0|673c75d516e8adb9e6ffc892"
	party.UserEmail = "sprov300@gmail.com"
	party.RequestId = "bks1m1g91jau4nkks2f0"

	type args struct {
		ctx context.Context
		in  *partyproto.CreatePartyRequest
	}
	tests := []struct {
		ps      *PartyService
		args    args
		wantErr bool
	}{
		{
			ps: partyService,
			args: args{
				ctx: ctx,
				in:  &party,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		partyResp, err := tt.ps.CreateParty(tt.args.ctx, tt.args.in)
		if (err != nil) != tt.wantErr {
			t.Errorf("PartyService.CreateParty() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		assert.NotNil(t, partyResp)
		partyResult := partyResp.Party
		assert.Equal(t, partyResult.PartyD.PartyName, "Ebeneser Scrooge Inc", "they should be equal")
	}
}
