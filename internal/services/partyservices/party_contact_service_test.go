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

func TestPartyService_GetPartyContact(t *testing.T) {
	err := test.LoadSQL(logUser, dbService)
	if err != nil {
		t.Error(err)
		return
	}

	ctx := LoginUser()

	timeAt, err := common.ConvertTimeToTimestamp(Layout, "2019-07-23T10:04:26Z")
	if err != nil {
		log.Error("Error", zap.Error(err))
	}

	partyService := NewPartyService(log, dbService, redisService, userServiceClient)
	partyContactD := partyproto.PartyContactD{}
	partyContactD.Id = uint32(3)
	partyContactD.Uuid4 = []byte{173, 35, 71, 252, 168, 25, 77, 28, 155, 139, 102, 136, 113, 57, 78, 249}
	partyContactD.IdS = "ad2347fc-a819-4d1c-9b8b-668871394ef9"
	partyContactD.FirstName = "Fred"
	partyContactD.MiddleName = "Churchill"
	partyContactD.Email = "fred@iytcorporation.gov.uk"
	partyContactD.PhoneMobile = "0127 2653214"
	partyContactD.PhoneFax = "0127 2653215"
	partyContactD.CountryCallingCode = "+91"

	crUpdTime := new(commonproto.CrUpdTime)
	crUpdTime.CreatedAt = timeAt
	crUpdTime.UpdatedAt = timeAt

	crUpdUser := new(commonproto.CrUpdUser)
	crUpdUser.StatusCode = "active"
	crUpdUser.CreatedByUserId = "auth0|673c75d516e8adb9e6ffc892"
	crUpdUser.UpdatedByUserId = "auth0|673c75d516e8adb9e6ffc892"

	partyInfo := commonproto.PartyInfo{}
	partyInfo.PartyId = uint32(2)
	partyInfo.PartyName = "IYT Corporation"

	partyContact := partyproto.PartyContact{PartyContactD: &partyContactD, PartyInfo: &partyInfo, CrUpdUser: crUpdUser, CrUpdTime: crUpdTime}

	partyContactResponse := partyproto.GetPartyContactResponse{}
	partyContactResponse.PartyContact = &partyContact

	form := partyproto.GetPartyContactRequest{}
	gform := commonproto.GetRequest{}
	gform.Id = "ad2347fc-a819-4d1c-9b8b-668871394ef9"
	gform.UserEmail = "sprov300@gmail.com"
	gform.RequestId = "bks1m1g91jau4nkks2f0"
	form.GetRequest = &gform

	type args struct {
		ctx context.Context
		in  *partyproto.GetPartyContactRequest
	}
	tests := []struct {
		ps      *PartyService
		args    args
		want    *partyproto.GetPartyContactResponse
		wantErr bool
	}{
		{
			ps: partyService,
			args: args{
				ctx: ctx,
				in:  &form,
			},
			want:    &partyContactResponse,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		partyContactResp, err := tt.ps.GetPartyContact(tt.args.ctx, tt.args.in)
		if (err != nil) != tt.wantErr {
			t.Errorf("PartyService.GetPartyContact() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(partyContactResp, tt.want) {
			t.Errorf("PartyService.GetPartyContact() = %v, want %v", partyContactResp, tt.want)
		}
		assert.NotNil(t, partyContactResp)
		partyContactResult := partyContactResp.PartyContact
		assert.Equal(t, partyContactResult.PartyContactD.FirstName, "Fred", "they should be equal")
		assert.Equal(t, partyContactResult.PartyContactD.MiddleName, "Churchill", "they should be equal")
		assert.Equal(t, partyContactResult.PartyContactD.Email, "fred@iytcorporation.gov.uk", "they should be equal")
		assert.Equal(t, partyContactResult.PartyContactD.PhoneMobile, "0127 2653214", "they should be equal")
		assert.Equal(t, partyContactResult.PartyContactD.PhoneFax, "0127 2653215", "they should be equal")
		assert.Equal(t, partyContactResult.PartyContactD.CountryCallingCode, "+91", "they should be equal")
	}
}
