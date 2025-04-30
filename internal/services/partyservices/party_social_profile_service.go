package partyservices

import (
	"context"

	"github.com/cloudfresco/sc-ubl/internal/common"
	commonproto "github.com/cloudfresco/sc-ubl/internal/protogen/common/v1"
	partyproto "github.com/cloudfresco/sc-ubl/internal/protogen/party/v1"
	commonstruct "github.com/cloudfresco/sc-ubl/internal/servicestructs/common/v1"
	partystruct "github.com/cloudfresco/sc-ubl/internal/servicestructs/party/v1"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

const insertPartySocialProfileSQL = `insert into party_social_profiles
	  ( 
    uuid4,
    social_profle_name,
    social_media_type_code,
    uri,
    party_id,
    party_name,
    status_code,
    created_by_user_id,
    updated_by_user_id,
    created_at,
    updated_at)
  values (:uuid4,
:social_profle_name,
:social_media_type_code,
:uri,
:party_id,
:party_name,
:status_code,
:created_by_user_id,
:updated_by_user_id,
:created_at,
:updated_at);`

const selectPartySocialProfilesSQL = `select 
  id,
  uuid4,
  social_profle_name,
  social_media_type_code,
  uri,
  party_id,
  party_name,
  status_code,
  created_by_user_id,
  updated_by_user_id,
  created_at,
  updated_at from party_social_profiles`

// CreatePartySocialProfile - Create PartySocialProfile
func (ps *PartyService) CreatePartySocialProfile(ctx context.Context, in *partyproto.CreatePartySocialProfileRequest) (*partyproto.CreatePartySocialProfileResponse, error) {
	user, err := GetUserWithNewContext(ctx, in.UserId, in.UserEmail, in.RequestId, ps.UserServiceClient)
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	ttime := common.GetTimeDetails()
	tn := common.TimeToTimestamp(ttime)

	getByIDRequest := commonproto.GetByIdRequest{}
	getByIDRequest.Id = in.PartyId
	getByIDRequest.UserEmail = in.UserEmail
	getByIDRequest.RequestId = in.RequestId
	form := partyproto.GetPartyByPkRequest{}
	form.GetByIdRequest = &getByIDRequest
	partyResponse, err := ps.GetPartyByPk(ctx, &form)
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	party := partyResponse.Party

	partySocialProfileD := partyproto.PartySocialProfileD{}
	partySocialProfileD.Uuid4, err = common.GetUUIDBytes()
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	partySocialProfileD.SocialProfleName = in.SocialProfleName
	partySocialProfileD.SocialMediaTypeCode = in.SocialMediaTypeCode
	partySocialProfileD.Uri = in.Uri
	/*  Common  */
	crUpdUser := commonproto.CrUpdUser{}
	crUpdUser.StatusCode = "active"
	crUpdUser.CreatedByUserId = user.Id
	crUpdUser.UpdatedByUserId = user.Id

	crUpdTime := commonproto.CrUpdTime{}
	crUpdTime.CreatedAt = tn
	crUpdTime.UpdatedAt = tn

	/* PartyInfo  */
	partyInfo := commonproto.PartyInfo{}
	partyInfo.PartyId = party.PartyD.Id
	partyInfo.PartyName = party.PartyD.PartyName

	partySocialProfile := partyproto.PartySocialProfile{PartySocialProfileD: &partySocialProfileD, PartyInfo: &partyInfo, CrUpdUser: &crUpdUser, CrUpdTime: &crUpdTime}

	err = ps.insertPartySocialProfile(ctx, insertPartySocialProfileSQL, &partySocialProfile, in.GetUserEmail(), in.GetRequestId())
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	partySocialProfileResponse := partyproto.CreatePartySocialProfileResponse{}
	partySocialProfileResponse.PartySocialProfile = &partySocialProfile
	return &partySocialProfileResponse, nil
}

// insertPartySocialProfile - insert PartySocialProfile
func (ps *PartyService) insertPartySocialProfile(ctx context.Context, insertPartySocialProfileSQL string, partySocialProfile *partyproto.PartySocialProfile, userEmail string, requestID string) error {
	partySocialProfileTmp, err := ps.crPartySocialProfileStruct(ctx, partySocialProfile, userEmail, requestID)
	if err != nil {
		ps.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
		return err
	}

	err = ps.DBService.InsUpd(ctx, userEmail, requestID, func(tx *sqlx.Tx) error {
		res, err := tx.NamedExecContext(ctx, insertPartySocialProfileSQL, partySocialProfileTmp)
		if err != nil {
			ps.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}

		uID, err := res.LastInsertId()
		if err != nil {
			ps.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}
		partySocialProfile.PartySocialProfileD.Id = uint32(uID)
		uuid4Str, err := common.UUIDBytesToStr(partySocialProfile.PartySocialProfileD.Uuid4)
		if err != nil {
			ps.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}
		partySocialProfile.PartySocialProfileD.IdS = uuid4Str
		return nil
	})

	if err != nil {
		ps.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
		return err
	}
	return nil
}

// crPartySocialProfileStruct - process PartySocialProfile details
func (ps *PartyService) crPartySocialProfileStruct(ctx context.Context, partySocialProfile *partyproto.PartySocialProfile, userEmail string, requestID string) (*partystruct.PartySocialProfile, error) {
	crUpdTime := new(commonstruct.CrUpdTime)
	crUpdTime.CreatedAt = common.TimestampToTime(partySocialProfile.CrUpdTime.CreatedAt)
	crUpdTime.UpdatedAt = common.TimestampToTime(partySocialProfile.CrUpdTime.UpdatedAt)

	partySocialProfileTmp := partystruct.PartySocialProfile{PartySocialProfileD: partySocialProfile.PartySocialProfileD, PartyInfo: partySocialProfile.PartyInfo, CrUpdUser: partySocialProfile.CrUpdUser, CrUpdTime: crUpdTime}

	return &partySocialProfileTmp, nil
}

// GetPartySocialProfile - Get Party Social Profile
func (ps *PartyService) GetPartySocialProfile(ctx context.Context, inReq *partyproto.GetPartySocialProfileRequest) (*partyproto.GetPartySocialProfileResponse, error) {
	in := inReq.GetRequest
	uuid4byte, err := common.UUIDStrToBytes(in.Id)
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	nselectPartySocialProfilesSQL := selectPartySocialProfilesSQL + ` where uuid4 = ? and status_code = ?;`

	row := ps.DBService.DB.QueryRowxContext(ctx, nselectPartySocialProfilesSQL, uuid4byte, "active")

	partySocialProfileTmp := partystruct.PartySocialProfile{}
	err = row.StructScan(&partySocialProfileTmp)
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	partySocialProfile, err := ps.getPartySocialProfileStruct(ctx, in, partySocialProfileTmp)
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	partySocialProfileResponse := partyproto.GetPartySocialProfileResponse{}
	partySocialProfileResponse.PartySocialProfile = partySocialProfile
	return &partySocialProfileResponse, nil
}

// getPartySocialProfileStruct - Get PartySocialProfile Struct
func (ps *PartyService) getPartySocialProfileStruct(ctx context.Context, in *commonproto.GetRequest, partySocialProfileTmp partystruct.PartySocialProfile) (*partyproto.PartySocialProfile, error) {
	uuid4Str, err := common.UUIDBytesToStr(partySocialProfileTmp.PartySocialProfileD.Uuid4)
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	partySocialProfileTmp.PartySocialProfileD.IdS = uuid4Str

	crUpdTime := new(commonproto.CrUpdTime)
	crUpdTime.CreatedAt = common.TimeToTimestamp(partySocialProfileTmp.CrUpdTime.CreatedAt)
	crUpdTime.UpdatedAt = common.TimeToTimestamp(partySocialProfileTmp.CrUpdTime.UpdatedAt)

	partySocialProfile := partyproto.PartySocialProfile{PartySocialProfileD: partySocialProfileTmp.PartySocialProfileD, PartyInfo: partySocialProfileTmp.PartyInfo, CrUpdUser: partySocialProfileTmp.CrUpdUser, CrUpdTime: crUpdTime}

	return &partySocialProfile, nil
}
