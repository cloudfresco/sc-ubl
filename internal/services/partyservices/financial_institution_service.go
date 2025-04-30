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

// FinancialInstitutionService - For accessing FinancialInstitution services
type FinancialInstitutionService struct {
	log               *zap.Logger
	DBService         *common.DBService
	RedisService      *common.RedisService
	UserServiceClient partyproto.UserServiceClient
	partyproto.UnimplementedFinancialInstitutionServiceServer
}

// NewFinancialInstitutionService - Create FinancialInstitution service
func NewFinancialInstitutionService(log *zap.Logger, dbOpt *common.DBService, redisOpt *common.RedisService, userServiceClient partyproto.UserServiceClient) *FinancialInstitutionService {
	return &FinancialInstitutionService{
		log:               log,
		DBService:         dbOpt,
		RedisService:      redisOpt,
		UserServiceClient: userServiceClient,
	}
}

const insertFinancialInstitutionSQL = `insert into financial_institutions
	  ( 
  uuid4,
  fi_id,
  name1,
  address_id,
  status_code,
  created_by_user_id,
  updated_by_user_id,
  created_at,
  updated_at)
  values (:uuid4,
:fi_id,
:name1,
:address_id,
:status_code,
:created_by_user_id,
:updated_by_user_id,
:created_at,
:updated_at);`

const insertFinancialInstitutionBranchSQL = `insert into financial_institution_branches
	  ( 
    uuid4,
    fb_id,
    name1,
    financial_institution_id,
    address_id,
  status_code,
  created_by_user_id,
  updated_by_user_id,
  created_at,
  updated_at)
  values (:uuid4,
    :fb_id,
    :name1,
    :financial_institution_id,
    :address_id,
    :status_code,
    :created_by_user_id,
    :updated_by_user_id,
    :created_at,
    :updated_at);`

const updateFinancialInstitutionSQL = `update financial_institutions set 
		  fi_id = ?,
      name1 = ?,
      status_code = ?,
			updated_at = ? where id = ? and status_code = ?;`

const deleteFinancialInstitutionSQL = `update financial_institutions set 
		  status_code = ?,
			updated_at = ? where uuid4= ?;`

const updateFinancialInstitutionBranchSQL = `update financial_institution_branches set 
  fb_id  = ?,
  name1  = ?, 
  status_code = ?,
  updated_at  = ? where id = ? and status_code = ?;`

const deleteFinancialInstitutionBranchSQL = `update financial_institution_branches set 
		  status_code = ?,
			updated_at = ? where uuid4= ?;`

const selectFinancialInstitutionsSQL = `select 
      id, 
      uuid4,
      fi_id,
      name1,
      address_id,
      status_code,
      created_by_user_id,
      updated_by_user_id,
      created_at,
      updated_at from financial_institutions`

const selectFinancialInstitutionBranchesSQL = `select 
	  id,
    uuid4,
    fb_id,
    name1,
    financial_institution_id,
    address_id,
   status_code,
  created_by_user_id,
  updated_by_user_id,
  created_at,
  updated_at from financial_institution_branches`

const selectFinancialInstitutionWithBranchesSQL = `select 
	  p.id,
    p.uuid4,
    p.fi_id,
    p.name1,
    p.address_id,
    p.status_code,
    p.created_at,
    p.updated_at,
    p.created_by_user_id,
    p.updated_by_user_id,
    m.id,
    m.uuid4,
    m.fb_id,
    m.name1,
    m.financial_institution_id,
    m.address_id,
    m.status_code,
    m.created_by_user_id,
    m.updated_by_user_id,
    m.created_at,
    m.updated_at from financial_institutions p inner join financial_institution_branches m on (p.id = m.financial_institution_id)`

// CreateFinancialInstitution - Create FinancialInstitution
func (fis *FinancialInstitutionService) CreateFinancialInstitution(ctx context.Context, in *partyproto.CreateFinancialInstitutionRequest) (*partyproto.CreateFinancialInstitutionResponse, error) {
	user, err := GetUserWithNewContext(ctx, in.UserId, in.UserEmail, in.RequestId, fis.UserServiceClient)
	if err != nil {
		fis.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	ttime := common.GetTimeDetails()
	tn := common.TimeToTimestamp(ttime)

	financialInstitutionD := partyproto.FinancialInstitutionD{}
	financialInstitutionD.Uuid4, err = common.GetUUIDBytes()
	if err != nil {
		fis.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	financialInstitutionD.FiId = in.FiId
	financialInstitutionD.Name1 = in.Name1

	crUpdUser := commonproto.CrUpdUser{}
	crUpdUser.StatusCode = "active"
	crUpdUser.CreatedByUserId = user.Id
	crUpdUser.UpdatedByUserId = user.Id

	crUpdTime := commonproto.CrUpdTime{}
	crUpdTime.CreatedAt = tn
	crUpdTime.UpdatedAt = tn

	financialInstitution := partyproto.FinancialInstitution{FinancialInstitutionD: &financialInstitutionD, CrUpdUser: &crUpdUser, CrUpdTime: &crUpdTime}

	err = fis.insertFinancialInstitution(ctx, insertFinancialInstitutionSQL, &financialInstitution, in.GetUserEmail(), in.GetRequestId())
	if err != nil {
		fis.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	financialInstitutionResponse := partyproto.CreateFinancialInstitutionResponse{}
	financialInstitutionResponse.FinancialInstitution = &financialInstitution
	return &financialInstitutionResponse, nil
}

// insertFinancialInstitution - Insert financialInstitution details into database
func (fis *FinancialInstitutionService) insertFinancialInstitution(ctx context.Context, insertFinancialInstitutionSQL string, financialInstitution *partyproto.FinancialInstitution, userEmail string, requestID string) error {
	financialInstitutionTmp, err := fis.crFinancialInstitutionStruct(ctx, financialInstitution, userEmail, requestID)
	if err != nil {
		fis.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
		return err
	}

	err = fis.DBService.InsUpd(ctx, userEmail, requestID, func(tx *sqlx.Tx) error {
		res, err := tx.NamedExecContext(ctx, insertFinancialInstitutionSQL, financialInstitutionTmp)
		if err != nil {
			fis.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}

		uID, err := res.LastInsertId()
		if err != nil {
			fis.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}
		financialInstitution.FinancialInstitutionD.Id = uint32(uID)
		uuid4Str, err := common.UUIDBytesToStr(financialInstitution.FinancialInstitutionD.Uuid4)
		if err != nil {
			fis.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}
		financialInstitution.FinancialInstitutionD.IdS = uuid4Str
		return nil
	})

	if err != nil {
		fis.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
		return err
	}
	return nil
}

// crFinancialInstitutionStruct - process FinancialInstitution details
func (fis FinancialInstitutionService) crFinancialInstitutionStruct(ctx context.Context, financialInstitution *partyproto.FinancialInstitution, userEmail string, requestID string) (*partystruct.FinancialInstitution, error) {
	crUpdTime := new(commonstruct.CrUpdTime)
	crUpdTime.CreatedAt = common.TimestampToTime(financialInstitution.CrUpdTime.CreatedAt)
	crUpdTime.UpdatedAt = common.TimestampToTime(financialInstitution.CrUpdTime.UpdatedAt)

	financialInstitutionTmp := partystruct.FinancialInstitution{FinancialInstitutionD: financialInstitution.FinancialInstitutionD, CrUpdUser: financialInstitution.CrUpdUser, CrUpdTime: crUpdTime}

	return &financialInstitutionTmp, nil
}

// CreateFinancialInstitutionBranch - Create FinancialInstitutionBranch
func (fis *FinancialInstitutionService) CreateFinancialInstitutionBranch(ctx context.Context, in *partyproto.CreateFinancialInstitutionBranchRequest) (*partyproto.CreateFinancialInstitutionBranchResponse, error) {
	user, err := GetUserWithNewContext(ctx, in.UserId, in.UserEmail, in.RequestId, fis.UserServiceClient)
	if err != nil {
		fis.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	ttime := common.GetTimeDetails()
	tn := common.TimeToTimestamp(ttime)

	financialInstitutionBranchD := partyproto.FinancialInstitutionBranchD{}
	financialInstitutionBranchD.Uuid4, err = common.GetUUIDBytes()
	if err != nil {
		fis.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	financialInstitutionBranchD.FbId = in.FbId
	financialInstitutionBranchD.Name1 = in.Name1
	financialInstitutionBranchD.FinancialInstitutionId = in.FinancialInstitutionId

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

	crUpdUser := commonproto.CrUpdUser{}
	crUpdUser.StatusCode = "active"
	crUpdUser.CreatedByUserId = user.Id
	crUpdUser.UpdatedByUserId = user.Id

	crUpdTime := commonproto.CrUpdTime{}
	crUpdTime.CreatedAt = tn
	crUpdTime.UpdatedAt = tn

	financialInstitutionBranch := partyproto.FinancialInstitutionBranch{FinancialInstitutionBranchD: &financialInstitutionBranchD, CrUpdUser: &crUpdUser, CrUpdTime: &crUpdTime}

	err = fis.insertFinancialInstitutionBranch(ctx, insertFinancialInstitutionBranchSQL, &financialInstitutionBranch, &address, in.GetUserEmail(), in.GetRequestId())
	if err != nil {
		fis.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	financialInstitutionBranchResponse := partyproto.CreateFinancialInstitutionBranchResponse{}
	financialInstitutionBranchResponse.FinancialInstitutionBranch = &financialInstitutionBranch
	return &financialInstitutionBranchResponse, nil
}

// insertFinancialInstitutionBranch - Create insertFinancialInstitutionBranch
func (fis *FinancialInstitutionService) insertFinancialInstitutionBranch(ctx context.Context, insertFinancialInstitutionBranchSQL string, financialInstitutionBranch *partyproto.FinancialInstitutionBranch, address *commonproto.Address, userEmail string, requestID string) error {
	financialInstitutionBranchTmp, err := fis.crFinancialInstitutionBranchStruct(ctx, financialInstitutionBranch, userEmail, requestID)
	if err != nil {
		fis.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
		return err
	}

	err = fis.DBService.InsUpd(ctx, userEmail, requestID, func(tx *sqlx.Tx) error {
		addr, err := common.InsertAddress(ctx, tx, address, userEmail, requestID)
		if err != nil {
			fis.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}

		financialInstitutionBranchTmp.FinancialInstitutionBranchD.AddressId = addr.Id
		financialInstitutionBranch.FinancialInstitutionBranchD.AddressId = addr.Id
		res, err := tx.NamedExecContext(ctx, insertFinancialInstitutionBranchSQL, financialInstitutionBranchTmp)
		if err != nil {
			fis.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}

		uID, err := res.LastInsertId()
		if err != nil {
			fis.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}
		financialInstitutionBranch.FinancialInstitutionBranchD.Id = uint32(uID)
		uuid4Str, err := common.UUIDBytesToStr(financialInstitutionBranch.FinancialInstitutionBranchD.Uuid4)
		if err != nil {
			fis.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}
		financialInstitutionBranch.FinancialInstitutionBranchD.IdS = uuid4Str
		return nil
	})
	if err != nil {
		fis.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
		return err
	}
	return nil
}

// crFinancialInstitutionBranchStruct - process FinancialInstitutionBranch details
func (fis FinancialInstitutionService) crFinancialInstitutionBranchStruct(ctx context.Context, financialInstitutionBranch *partyproto.FinancialInstitutionBranch, userEmail string, requestID string) (*partystruct.FinancialInstitutionBranch, error) {
	crUpdTime := new(commonstruct.CrUpdTime)
	crUpdTime.CreatedAt = common.TimestampToTime(financialInstitutionBranch.CrUpdTime.CreatedAt)
	crUpdTime.UpdatedAt = common.TimestampToTime(financialInstitutionBranch.CrUpdTime.UpdatedAt)

	financialInstitutionBranchTmp := partystruct.FinancialInstitutionBranch{FinancialInstitutionBranchD: financialInstitutionBranch.FinancialInstitutionBranchD, CrUpdUser: financialInstitutionBranch.CrUpdUser, CrUpdTime: crUpdTime}

	return &financialInstitutionBranchTmp, nil
}

// GetFinancialInstitutions - GetFinancialInstitutions
func (fis *FinancialInstitutionService) GetFinancialInstitutions(ctx context.Context, in *partyproto.GetFinancialInstitutionsRequest) (*partyproto.GetFinancialInstitutionsResponse, error) {
	limit := in.GetLimit()
	nextCursor := in.GetNextCursor()
	if limit == "" {
		limit = fis.DBService.LimitSQLRows
	}
	query := "levelp = ? and status_code = ?"
	if nextCursor == "" {
		query = query + " order by id desc " + " limit " + limit + ";"
	} else {
		nextCursor = common.DecodeCursor(nextCursor)
		query = query + " " + "and" + " " + "id <= " + nextCursor + " order by id desc " + " limit " + limit + ";"
	}

	financialInstitutions := []*partyproto.FinancialInstitution{}

	nselectFinancialInstitutionSQL := selectFinancialInstitutionsSQL + ` where ` + query
	rows, err := fis.DBService.DB.QueryxContext(ctx, nselectFinancialInstitutionSQL, 0, "active")
	if err != nil {
		fis.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	for rows.Next() {
		financialInstitutionTmp := partystruct.FinancialInstitution{}
		err = rows.StructScan(&financialInstitutionTmp)

		if err != nil {
			fis.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
			return nil, err
		}

		getRequest := commonproto.GetRequest{}
		getRequest.UserEmail = in.UserEmail
		getRequest.RequestId = in.RequestId
		financialInstitution, err := fis.getFinancialInstitutionStruct(ctx, &getRequest, financialInstitutionTmp)
		if err != nil {
			fis.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
			return nil, err
		}

		financialInstitutions = append(financialInstitutions, financialInstitution)
	}

	financialInstitutionsResponse := partyproto.GetFinancialInstitutionsResponse{}
	if len(financialInstitutions) != 0 {
		next := financialInstitutions[len(financialInstitutions)-1].FinancialInstitutionD.Id
		next--
		nextc := common.EncodeCursor(next)
		financialInstitutionsResponse = partyproto.GetFinancialInstitutionsResponse{FinancialInstitutions: financialInstitutions, NextCursor: nextc}
	} else {
		financialInstitutionsResponse = partyproto.GetFinancialInstitutionsResponse{FinancialInstitutions: financialInstitutions, NextCursor: "0"}
	}
	return &financialInstitutionsResponse, nil
}

// GetFinancialInstitution - Get FinancialInstitution
func (fis *FinancialInstitutionService) GetFinancialInstitution(ctx context.Context, inReq *partyproto.GetFinancialInstitutionRequest) (*partyproto.GetFinancialInstitutionResponse, error) {
	in := inReq.GetRequest
	uuid4byte, err := common.UUIDStrToBytes(in.Id)
	if err != nil {
		fis.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	nselectFinancialInstitutionsSQL := selectFinancialInstitutionsSQL + ` where uuid4 = ? and status_code = ?;`
	row := fis.DBService.DB.QueryRowxContext(ctx, nselectFinancialInstitutionsSQL, uuid4byte, "active")

	financialInstitutionTmp := partystruct.FinancialInstitution{}
	err = row.StructScan(&financialInstitutionTmp)
	if err != nil {
		fis.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	getRequest := commonproto.GetRequest{}
	getRequest.UserEmail = in.UserEmail
	getRequest.RequestId = in.RequestId
	financialInstitution, err := fis.getFinancialInstitutionStruct(ctx, &getRequest, financialInstitutionTmp)
	if err != nil {
		fis.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	financialInstitutionResponse := partyproto.GetFinancialInstitutionResponse{}
	financialInstitutionResponse.FinancialInstitution = financialInstitution
	return &financialInstitutionResponse, nil
}

// GetFinancialInstitutionWithBranches - Get Financial Institution With Branches
func (fis *FinancialInstitutionService) GetFinancialInstitutionWithBranches(ctx context.Context, inReq *partyproto.GetFinancialInstitutionWithBranchesRequest) (*partyproto.GetFinancialInstitutionWithBranchesResponse, error) {
	in := inReq.GetRequest
	uuid4byte, err := common.UUIDStrToBytes(in.Id)
	if err != nil {
		fis.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	financialInstitution := partyproto.FinancialInstitution{}
	financialInstitutionD := partyproto.FinancialInstitutionD{}
	nselectFinancialInstitutionWithBranchesSQL := selectFinancialInstitutionWithBranchesSQL + ` where p.uuid4 = ?;`
	rows, err := fis.DBService.DB.QueryxContext(ctx, nselectFinancialInstitutionWithBranchesSQL, uuid4byte)
	if err != nil {
		fis.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	for rows.Next() {

		financialInstitutionBranchD := partyproto.FinancialInstitutionBranchD{}
		crUpdTimeTmp := commonstruct.CrUpdTime{}
		crUpdUserTmp := commonproto.CrUpdUser{}

		branchcrUpdTimeTmp := commonstruct.CrUpdTime{}
		branchcrUpdUserTmp := commonproto.CrUpdUser{}

		err := rows.Scan(
			&financialInstitutionD.Id,
			&financialInstitutionD.Uuid4,
			&financialInstitutionD.FiId,
			&financialInstitutionD.Name1,
			&crUpdUserTmp.StatusCode,
			&crUpdUserTmp.CreatedByUserId,
			&crUpdUserTmp.UpdatedByUserId,
			&crUpdTimeTmp.CreatedAt,
			&crUpdTimeTmp.UpdatedAt,
			&financialInstitutionBranchD.Id,
			&financialInstitutionBranchD.Uuid4,
			&financialInstitutionBranchD.FbId,
			&financialInstitutionBranchD.Name1,
			&financialInstitutionBranchD.FinancialInstitutionId,
			&financialInstitutionBranchD.AddressId,
			&branchcrUpdUserTmp.StatusCode,
			&branchcrUpdUserTmp.CreatedByUserId,
			&branchcrUpdUserTmp.UpdatedByUserId,
			&branchcrUpdTimeTmp.CreatedAt,
			&branchcrUpdTimeTmp.UpdatedAt)
		if err != nil {
			fis.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Int("fiBranchnum", 5315), zap.Error(err))
			return nil, err
		}

		crUpdTime := new(commonproto.CrUpdTime)
		crUpdTime.CreatedAt = common.TimeToTimestamp(crUpdTimeTmp.CreatedAt)
		crUpdTime.UpdatedAt = common.TimeToTimestamp(crUpdTimeTmp.UpdatedAt)

		uuid4Str1, err := common.UUIDBytesToStr(financialInstitutionD.Uuid4)
		if err != nil {
			fis.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Int("fiBranchnum", 5316), zap.Error(err))
			return nil, err
		}
		financialInstitutionD.IdS = uuid4Str1

		financialInstitution = partyproto.FinancialInstitution{FinancialInstitutionD: &financialInstitutionD, CrUpdUser: &crUpdUserTmp, CrUpdTime: crUpdTime}

		branchcrUpdTime := new(commonproto.CrUpdTime)
		branchcrUpdTime.CreatedAt = common.TimeToTimestamp(branchcrUpdTimeTmp.CreatedAt)
		branchcrUpdTime.UpdatedAt = common.TimeToTimestamp(branchcrUpdTimeTmp.UpdatedAt)

		uuid4Str, err := common.UUIDBytesToStr(financialInstitutionBranchD.Uuid4)
		if err != nil {
			fis.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Int("fiBranchnum", 5317), zap.Error(err))
			return nil, err
		}
		financialInstitutionBranchD.IdS = uuid4Str

		financialInstitutionBranch := partyproto.FinancialInstitutionBranch{FinancialInstitutionBranchD: &financialInstitutionBranchD, CrUpdUser: &branchcrUpdUserTmp, CrUpdTime: branchcrUpdTime}

		financialInstitution.FinancialInstitutionBranches = append(financialInstitution.FinancialInstitutionBranches, &financialInstitutionBranch)
	}

	financialInstitutionBranchResponse := partyproto.GetFinancialInstitutionWithBranchesResponse{}
	financialInstitutionBranchResponse.FinancialInstitution = &financialInstitution
	return &financialInstitutionBranchResponse, nil
}

// getFinancialInstitutionStruct - Get FinancialInstitution Struct
func (fis *FinancialInstitutionService) getFinancialInstitutionStruct(ctx context.Context, in *commonproto.GetRequest, financialInstitutionTmp partystruct.FinancialInstitution) (*partyproto.FinancialInstitution, error) {
	uuid4Str, err := common.UUIDBytesToStr(financialInstitutionTmp.FinancialInstitutionD.Uuid4)
	if err != nil {
		fis.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	financialInstitutionTmp.FinancialInstitutionD.IdS = uuid4Str

	crUpdTime := new(commonproto.CrUpdTime)
	crUpdTime.CreatedAt = common.TimeToTimestamp(financialInstitutionTmp.CrUpdTime.CreatedAt)
	crUpdTime.UpdatedAt = common.TimeToTimestamp(financialInstitutionTmp.CrUpdTime.UpdatedAt)

	financialInstitution := partyproto.FinancialInstitution{FinancialInstitutionD: financialInstitutionTmp.FinancialInstitutionD, CrUpdUser: financialInstitutionTmp.CrUpdUser, CrUpdTime: crUpdTime}

	return &financialInstitution, nil
}

// UpdateFinancialInstitution - Update financialInstitution
func (fis *FinancialInstitutionService) UpdateFinancialInstitution(ctx context.Context, in *partyproto.UpdateFinancialInstitutionRequest) (*partyproto.UpdateFinancialInstitutionResponse, error) {
	getRequest := commonproto.GetRequest{}
	getRequest.Id = in.FiId
	getRequest.UserEmail = in.UserEmail
	getRequest.RequestId = in.RequestId

	form := partyproto.GetFinancialInstitutionRequest{}
	form.GetRequest = &getRequest
	financialInstitutionResp, err := fis.GetFinancialInstitution(ctx, &form)
	if err != nil {
		fis.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	financialInstitution := financialInstitutionResp.FinancialInstitution
	db := fis.DBService.DB
	tn := common.GetTimeDetails()
	stmt, err := db.PreparexContext(ctx, updateFinancialInstitutionSQL)
	if err != nil {
		fis.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	err = fis.DBService.InsUpd(ctx, in.GetUserEmail(), in.GetRequestId(), func(tx *sqlx.Tx) error {
		_, err = tx.StmtxContext(ctx, stmt).ExecContext(ctx,
			in.FiId,
			in.Name1,
			tn,
			financialInstitution.FinancialInstitutionD.Id,
			"active")
		if err != nil {
			fis.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
			err1 := stmt.Close()
			if err1 != nil {
				fis.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err1))
				return err1
			}
			return err
		}
		return nil
	})

	if err != nil {
		fis.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	err = stmt.Close()
	if err != nil {
		fis.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	return &partyproto.UpdateFinancialInstitutionResponse{}, nil
}

// DeleteFinancialInstitution - Delete financialInstitution
func (fis *FinancialInstitutionService) DeleteFinancialInstitution(ctx context.Context, inReq *partyproto.DeleteFinancialInstitutionRequest) (*partyproto.DeleteFinancialInstitutionResponse, error) {
	in := inReq.GetRequest
	uuid4byte, err := common.UUIDStrToBytes(in.Id)
	if err != nil {
		fis.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	db := fis.DBService.DB
	tn := common.GetTimeDetails()
	stmt, err := db.PreparexContext(ctx, deleteFinancialInstitutionSQL)
	if err != nil {
		fis.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	err = fis.DBService.InsUpd(ctx, in.GetUserEmail(), in.GetRequestId(), func(tx *sqlx.Tx) error {
		_, err = tx.StmtxContext(ctx, stmt).ExecContext(ctx,
			"inactive",
			tn,
			uuid4byte)

		if err != nil {
			fis.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
			err1 := stmt.Close()
			if err != nil {
				fis.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err1))
				return err1
			}
			return err
		}
		return nil
	})

	if err != nil {
		fis.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	err = stmt.Close()
	if err != nil {
		fis.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	return &partyproto.DeleteFinancialInstitutionResponse{}, nil
}

// GetFinancialInstitutionBranch - Get FinancialInstitution Branch
func (fis *FinancialInstitutionService) GetFinancialInstitutionBranch(ctx context.Context, inReq *partyproto.GetFinancialInstitutionBranchRequest) (*partyproto.GetFinancialInstitutionBranchResponse, error) {
	in := inReq.GetRequest
	uuid4byte, err := common.UUIDStrToBytes(in.Id)
	if err != nil {
		fis.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	nselectFinancialInstitutionBranchesSQL := selectFinancialInstitutionBranchesSQL + ` where uuid4 = ? and status_code = ?;`

	row := fis.DBService.DB.QueryRowxContext(ctx, nselectFinancialInstitutionBranchesSQL, uuid4byte, "active")

	financialInstitutionBranchTmp := partystruct.FinancialInstitutionBranch{}
	err = row.StructScan(&financialInstitutionBranchTmp)
	if err != nil {
		fis.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	getRequest := commonproto.GetRequest{}
	getRequest.UserEmail = in.UserEmail
	getRequest.RequestId = in.RequestId
	financialInstitutionBranch, err := fis.getFinancialInstitutionBranchStruct(ctx, &getRequest, financialInstitutionBranchTmp)
	if err != nil {
		fis.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	financialInstitutionBranchResponse := partyproto.GetFinancialInstitutionBranchResponse{}
	financialInstitutionBranchResponse.FinancialInstitutionBranch = financialInstitutionBranch
	return &financialInstitutionBranchResponse, nil
}

// getFinancialInstitutionBranchStruct - Get FinancialInstitutionBranch Struct
func (fis *FinancialInstitutionService) getFinancialInstitutionBranchStruct(ctx context.Context, in *commonproto.GetRequest, financialInstitutionBranchTmp partystruct.FinancialInstitutionBranch) (*partyproto.FinancialInstitutionBranch, error) {
	uuid4Str, err := common.UUIDBytesToStr(financialInstitutionBranchTmp.FinancialInstitutionBranchD.Uuid4)
	if err != nil {
		fis.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	financialInstitutionBranchTmp.FinancialInstitutionBranchD.IdS = uuid4Str

	crUpdTime := new(commonproto.CrUpdTime)
	crUpdTime.CreatedAt = common.TimeToTimestamp(financialInstitutionBranchTmp.CrUpdTime.CreatedAt)
	crUpdTime.UpdatedAt = common.TimeToTimestamp(financialInstitutionBranchTmp.CrUpdTime.UpdatedAt)

	financialInstitutionBranch := partyproto.FinancialInstitutionBranch{FinancialInstitutionBranchD: financialInstitutionBranchTmp.FinancialInstitutionBranchD, CrUpdUser: financialInstitutionBranchTmp.CrUpdUser, CrUpdTime: crUpdTime}

	return &financialInstitutionBranch, nil
}

// UpdateFinancialInstitutionBranch -Update Financial Institution Branch
func (fis *FinancialInstitutionService) UpdateFinancialInstitutionBranch(ctx context.Context, in *partyproto.UpdateFinancialInstitutionBranchRequest) (*partyproto.UpdateFinancialInstitutionBranchResponse, error) {
	getRequest := commonproto.GetRequest{}
	getRequest.Id = in.FinancialInstitutionBranchId
	getRequest.UserEmail = in.UserEmail
	getRequest.RequestId = in.RequestId

	form := partyproto.GetFinancialInstitutionBranchRequest{}
	form.GetRequest = &getRequest

	financialInstitutionBranchResp, err := fis.GetFinancialInstitutionBranch(ctx, &form)
	if err != nil {
		fis.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	financialInstitutionBranch := financialInstitutionBranchResp.FinancialInstitutionBranch

	db := fis.DBService.DB
	tn := common.GetTimeDetails()
	stmt, err := db.PreparexContext(ctx, updateFinancialInstitutionBranchSQL)
	if err != nil {
		fis.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	err = fis.DBService.InsUpd(ctx, in.GetUserEmail(), in.GetRequestId(), func(tx *sqlx.Tx) error {
		_, err = tx.StmtxContext(ctx, stmt).ExecContext(ctx,
			in.FbId,
			in.Name1,
			"active",
			tn,
			financialInstitutionBranch.FinancialInstitutionBranchD.Id,
			"active")
		if err != nil {
			fis.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
			err1 := stmt.Close()
			if err1 != nil {
				fis.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err1))
				return err1
			}
			return err
		}
		return nil
	})

	if err != nil {
		fis.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	err = stmt.Close()
	if err != nil {
		fis.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	return &partyproto.UpdateFinancialInstitutionBranchResponse{}, nil
}

// DeleteFinancialInstitutionBranch - Delete financialInstitution Branch
func (fis *FinancialInstitutionService) DeleteFinancialInstitutionBranch(ctx context.Context, inReq *partyproto.DeleteFinancialInstitutionBranchRequest) (*partyproto.DeleteFinancialInstitutionBranchResponse, error) {
	in := inReq.GetRequest
	uuid4byte, err := common.UUIDStrToBytes(in.Id)
	if err != nil {
		fis.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	db := fis.DBService.DB
	tn := common.GetTimeDetails()
	stmt, err := db.PreparexContext(ctx, deleteFinancialInstitutionBranchSQL)
	if err != nil {
		fis.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	err = fis.DBService.InsUpd(ctx, in.GetUserEmail(), in.GetRequestId(), func(tx *sqlx.Tx) error {
		_, err = tx.StmtxContext(ctx, stmt).ExecContext(ctx,
			"inactive",
			tn,
			uuid4byte)

		if err != nil {
			fis.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
			err1 := stmt.Close()
			if err1 != nil {
				fis.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err1))
				return err1
			}
			return err
		}

		return nil
	})

	err = stmt.Close()
	if err != nil {
		fis.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	return &partyproto.DeleteFinancialInstitutionBranchResponse{}, nil
}
