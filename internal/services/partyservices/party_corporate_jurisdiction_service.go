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

const insertPartyCorporateJurisdictionSQL = `insert into party_corporate_jurisdictions
	  ( 
  uuid4,
  address_id,
  party_id,
  party_name,
  status_code,
  created_by_user_id,
  updated_by_user_id,
  created_at,
  updated_at)
  values (:uuid4,
:address_id,
:party_id,
:party_name,
:status_code,
:created_by_user_id,
:updated_by_user_id,
:created_at,
:updated_at);`

const selectPartyCorporateJurisdictionsSQL = `select 
  id,
  uuid4,
  address_id,
  party_id,
  party_name,
  status_code,
  created_by_user_id,
  updated_by_user_id,
  created_at,
  updated_at from party_corporate_jurisdictions`

// CreatePartyCorporateJurisdiction - Create PartyCorporateJurisdiction
func (ps *PartyService) CreatePartyCorporateJurisdiction(ctx context.Context, in *partyproto.CreatePartyCorporateJurisdictionRequest) (*partyproto.CreatePartyCorporateJurisdictionResponse, error) {
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

	partyCorporateJurisdictionD := partyproto.PartyCorporateJurisdictionD{}
	partyCorporateJurisdictionD.Uuid4, err = common.GetUUIDBytes()
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	address := commonproto.Address{}
	/*  Addr  */
	address.AddrListAgencyId = in.AddrListAgencyId
	address.AddrListId = in.AddrListAgencyId
	address.AddrListVersionId = in.AddrListAgencyId
	address.AddressTypeCode = in.AddressTypeCode
	address.AddressFormatCode = in.AddressFormatCode
	address.Postbox = in.Postbox
	address.Floor1 = in.Floor
	address.Room = in.Room
	address.StreetName = in.StreetName
	address.AdditionalStreetName = in.AdditionalStreetName
	address.BlockName = in.BlockName
	address.BuildingName = in.BuildingName
	address.BuildingNumber = in.BuildingNumber
	address.InhouseMail = in.InhouseMail
	address.Department = in.Department
	address.MarkAttention = in.MarkAttention
	address.MarkCare = in.MarkCare
	address.PlotIdentification = in.PlotIdentification
	address.CitySubdivisionName = in.CitySubdivisionName
	address.CityName = in.CityName
	address.PostalZone = in.PostalZone
	address.CountrySubentity = in.CountrySubentity
	address.CountrySubentityCode = in.CountrySubentityCode
	address.Region = in.Region
	address.District = in.District
	address.TimezoneOffset = in.TimezoneOffset
	address.CountryIdCode = in.CountryIdCode
	address.CountryName = in.CountryName
	address.LocationCoordLat = in.LocationCoordLat
	address.LocationCoordLon = in.LocationCoordLon
	address.Note = in.Note
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

	partyCorporateJurisdiction := partyproto.PartyCorporateJurisdiction{PartyCorporateJurisdictionD: &partyCorporateJurisdictionD, PartyInfo: &partyInfo, CrUpdUser: &crUpdUser, CrUpdTime: &crUpdTime}

	err = ps.insertPartyCorporateJurisdiction(ctx, insertPartyCorporateJurisdictionSQL, &partyCorporateJurisdiction, &address, in.GetUserEmail(), in.GetRequestId())
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	partyCorporateJurisdictionResponse := partyproto.CreatePartyCorporateJurisdictionResponse{}
	partyCorporateJurisdictionResponse.PartyCorporateJurisdiction = &partyCorporateJurisdiction
	return &partyCorporateJurisdictionResponse, nil
}

// insertPartyCorporateJurisdiction - insert PartyCorporateJurisdiction
func (ps *PartyService) insertPartyCorporateJurisdiction(ctx context.Context, insertPartyCorporateJurisdictionSQL string, partyCorporateJurisdiction *partyproto.PartyCorporateJurisdiction, address *commonproto.Address, userEmail string, requestID string) error {
	partyCorporateJurisdictionTmp, err := ps.crPartyCorporateJurisdictionStruct(ctx, partyCorporateJurisdiction, userEmail, requestID)
	if err != nil {
		ps.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
		return err
	}

	err = ps.DBService.InsUpd(ctx, userEmail, requestID, func(tx *sqlx.Tx) error {
		addr, err := common.InsertAddress(ctx, tx, address, userEmail, requestID)
		if err != nil {
			ps.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}
		partyCorporateJurisdictionTmp.PartyCorporateJurisdictionD.AddressId = addr.Id
		partyCorporateJurisdiction.PartyCorporateJurisdictionD.AddressId = addr.Id
		res, err := tx.NamedExecContext(ctx, insertPartyCorporateJurisdictionSQL, partyCorporateJurisdictionTmp)
		if err != nil {
			ps.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}
		uID, err := res.LastInsertId()
		if err != nil {
			ps.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}
		partyCorporateJurisdiction.PartyCorporateJurisdictionD.Id = uint32(uID)
		partyCorporateJurisdiction.PartyCorporateJurisdictionD.AddressId = addr.Id
		uuid4Str, err := common.UUIDBytesToStr(partyCorporateJurisdiction.PartyCorporateJurisdictionD.Uuid4)
		if err != nil {
			ps.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}
		partyCorporateJurisdiction.PartyCorporateJurisdictionD.IdS = uuid4Str
		return nil
	})

	if err != nil {
		ps.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
		return err
	}
	return nil
}

// crPartyCorporateJurisdictionStruct - process PartyCorporateJurisdiction details
func (ps *PartyService) crPartyCorporateJurisdictionStruct(ctx context.Context, partyCorporateJurisdiction *partyproto.PartyCorporateJurisdiction, userEmail string, requestID string) (*partystruct.PartyCorporateJurisdiction, error) {
	crUpdTime := new(commonstruct.CrUpdTime)
	crUpdTime.CreatedAt = common.TimestampToTime(partyCorporateJurisdiction.CrUpdTime.CreatedAt)
	crUpdTime.UpdatedAt = common.TimestampToTime(partyCorporateJurisdiction.CrUpdTime.UpdatedAt)

	partyCorporateJurisdictionTmp := partystruct.PartyCorporateJurisdiction{PartyCorporateJurisdictionD: partyCorporateJurisdiction.PartyCorporateJurisdictionD, PartyInfo: partyCorporateJurisdiction.PartyInfo, CrUpdUser: partyCorporateJurisdiction.CrUpdUser, CrUpdTime: crUpdTime}

	return &partyCorporateJurisdictionTmp, nil
}

// GetPartyCorporateJurisdiction - Get Party Corporate Jurisdiction
func (ps *PartyService) GetPartyCorporateJurisdiction(ctx context.Context, inReq *partyproto.GetPartyCorporateJurisdictionRequest) (*partyproto.GetPartyCorporateJurisdictionResponse, error) {
	in := inReq.GetRequest
	uuid4byte, err := common.UUIDStrToBytes(in.Id)
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	nselectPartyCorporateJurisdictionsSQL := selectPartyCorporateJurisdictionsSQL + ` where uuid4 = ? and status_code = ?;`
	row := ps.DBService.DB.QueryRowxContext(ctx, nselectPartyCorporateJurisdictionsSQL, uuid4byte, "active")

	partyCorporateJurisdictionTmp := partystruct.PartyCorporateJurisdiction{}
	err = row.StructScan(&partyCorporateJurisdictionTmp)
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	partyCorporateJurisdiction, err := ps.getPartyCorporateJurisdictionStruct(ctx, in, partyCorporateJurisdictionTmp)
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	partyCorporateJurisdictionResponse := partyproto.GetPartyCorporateJurisdictionResponse{}
	partyCorporateJurisdictionResponse.PartyCorporateJurisdiction = partyCorporateJurisdiction
	return &partyCorporateJurisdictionResponse, nil
}

// getPartyCorporateJurisdictionStruct - Get PartyCorporateJurisdiction Struct
func (ps *PartyService) getPartyCorporateJurisdictionStruct(ctx context.Context, in *commonproto.GetRequest, partyCorporateJurisdictionTmp partystruct.PartyCorporateJurisdiction) (*partyproto.PartyCorporateJurisdiction, error) {
	uuid4Str, err := common.UUIDBytesToStr(partyCorporateJurisdictionTmp.PartyCorporateJurisdictionD.Uuid4)
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	partyCorporateJurisdictionTmp.PartyCorporateJurisdictionD.IdS = uuid4Str

	crUpdTime := new(commonproto.CrUpdTime)
	crUpdTime.CreatedAt = common.TimeToTimestamp(partyCorporateJurisdictionTmp.CrUpdTime.CreatedAt)
	crUpdTime.UpdatedAt = common.TimeToTimestamp(partyCorporateJurisdictionTmp.CrUpdTime.UpdatedAt)

	partyCorporateJurisdiction := partyproto.PartyCorporateJurisdiction{PartyCorporateJurisdictionD: partyCorporateJurisdictionTmp.PartyCorporateJurisdictionD, PartyInfo: partyCorporateJurisdictionTmp.PartyInfo, CrUpdUser: partyCorporateJurisdictionTmp.CrUpdUser, CrUpdTime: crUpdTime}

	return &partyCorporateJurisdiction, nil
}
