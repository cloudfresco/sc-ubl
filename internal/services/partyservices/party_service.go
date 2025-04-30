package partyservices

import (
	"context"
	"net"
	"os"
	"time"

	"github.com/cloudfresco/sc-ubl/internal/common"
	"github.com/cloudfresco/sc-ubl/internal/config"
	commonproto "github.com/cloudfresco/sc-ubl/internal/protogen/common/v1"
	partyproto "github.com/cloudfresco/sc-ubl/internal/protogen/party/v1"
	taxproto "github.com/cloudfresco/sc-ubl/internal/protogen/tax/v1"
	commonstruct "github.com/cloudfresco/sc-ubl/internal/servicestructs/common/v1"
	partystruct "github.com/cloudfresco/sc-ubl/internal/servicestructs/party/v1"
	taxstruct "github.com/cloudfresco/sc-ubl/internal/servicestructs/tax/v1"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	timestamp "google.golang.org/protobuf/types/known/timestamppb"
)

// PartyService - For accessing Party services
type PartyService struct {
	log               *zap.Logger
	DBService         *common.DBService
	RedisService      *common.RedisService
	UserServiceClient partyproto.UserServiceClient
	partyproto.UnimplementedPartyServiceServer
}

// NewPartyService - Create Party service
func NewPartyService(log *zap.Logger, dbOpt *common.DBService, redisOpt *common.RedisService, userServiceClient partyproto.UserServiceClient) *PartyService {
	return &PartyService{
		log:               log,
		DBService:         dbOpt,
		RedisService:      redisOpt,
		UserServiceClient: userServiceClient,
	}
}

// StartPartyServer - Start Party server
func StartPartyServer(log *zap.Logger, isTest bool, pwd string, dbOpt *config.DBOptions, redisOpt *config.RedisOptions, mailerOpt *config.MailerOptions, grpcServerOpt *config.GrpcServerOptions, jwtOpt *config.JWTOptions, oauthOpt *config.OauthOptions, userOpt *config.UserOptions, uptraceOpt *config.UptraceOptions, dbService *common.DBService, redisService *common.RedisService, mailerService common.MailerIntf) {
	common.SetJWTOpt(jwtOpt)

	creds, err := common.GetSrvCred(log, isTest, pwd, grpcServerOpt)
	if err != nil {
		os.Exit(1)
	}

	userCreds, err := common.GetClientCred(log, isTest, pwd, grpcServerOpt)
	if err != nil {
		os.Exit(1)
	}

	var srvOpts []grpc.ServerOption

	userConn, err := grpc.NewClient(grpcServerOpt.GrpcUserServerPort, grpc.WithTransportCredentials(userCreds), grpc.WithStatsHandler(otelgrpc.NewClientHandler()))
	if err != nil {
		log.Error("Error", zap.Error(err))
		os.Exit(1)
	}

	srvOpts = append(srvOpts, grpc.Creds(creds))

	srvOpts = append(srvOpts, grpc.StatsHandler(otelgrpc.NewServerHandler()))

	uc := partyproto.NewUserServiceClient(userConn)
	partyService := NewPartyService(log, dbService, redisService, uc)
	financialinstitutionService := NewFinancialInstitutionService(log, dbService, redisService, uc)

	lis, err := net.Listen("tcp", grpcServerOpt.GrpcPartyServerPort)
	if err != nil {
		log.Error("Error", zap.Error(err))
		os.Exit(1)
	}

	srv := grpc.NewServer(srvOpts...)
	partyproto.RegisterPartyServiceServer(srv, partyService)
	partyproto.RegisterFinancialInstitutionServiceServer(srv, financialinstitutionService)

	if err := srv.Serve(lis); err != nil {
		log.Error("Error", zap.Error(err))
		os.Exit(1)
	}
}

const insertPartySQL = `insert into parties
	  ( 
    uuid4,
    party_endpoint_id,
    party_endpoint_scheme_id,
    party_name,
    party_desc,
    party_type,
    level_p,
    parent_id,
    num_chd,
    leaf,
    tax_reference1,
    tax_reference2,
    public_key,
    address_id,
    registration_name,
    company_id,
    company_legal_form_code,
    company_legal_form,
    sole_proprietorship_indicator,
    company_liquidation_status_code,
    corporate_stock_amount,
    fully_paid_shares_indicator,
    corporate_registration_id,
    corporate_registration_name,
    corporate_registration_type_code,
    tax_level_code,
    exemption_reason_code,
    exemption_reason,
    tax_scheme_id,
    registration_date,
    registration_expiration_date,
    status_code,
    created_by_user_id,
    updated_by_user_id,
    created_at,
    updated_at)
  values (:uuid4,
    :party_endpoint_id,
    :party_endpoint_scheme_id,
    :party_name,
    :party_desc,
    :party_type,
    :level_p,
    :parent_id,
    :num_chd,
    :leaf,
    :tax_reference1,
    :tax_reference2,
    :public_key,
    :address_id,
    :registration_name,
    :company_id,
    :company_legal_form_code,
    :company_legal_form,
    :sole_proprietorship_indicator,
    :company_liquidation_status_code,
    :corporate_stock_amount,
    :fully_paid_shares_indicator,
    :corporate_registration_id,
    :corporate_registration_name,
    :corporate_registration_type_code,
    :tax_level_code,
    :exemption_reason_code,
    :exemption_reason,
    :tax_scheme_id,
    :registration_date,
    :registration_expiration_date,
    :status_code,
    :created_by_user_id,
    :updated_by_user_id,
    :created_at,
    :updated_at);`

const insertChildPartySQL = `insert into party_chds
	  ( 
    uuid4,
    party_id,
    party_chd_id,
    status_code,
    created_by_user_id,
    updated_by_user_id,
    created_at,
    updated_at)
  values (:uuid4,
          :party_id,
          :party_chd_id,
          :status_code,
          :created_by_user_id,
          :updated_by_user_id,
          :created_at,
          :updated_at);`

const updateNumChildrenSQL = `update parties set 
				  num_chd = ?,
          leaf = ?,
				  updated_at = ? where id = ? and status_code = ?;`

const updatePartySQL = `update parties set 
		  party_name = ?,
      party_desc = ?,
			updated_at = ? where id = ? and status_code = ?;`

const deletePartySQL = `update parties set 
		  status_code = ?,
			updated_at = ? where uuid4= ?;`

const selectPartiesSQL = `select 
      id, 
      uuid4,
      party_endpoint_id,
      party_endpoint_scheme_id,
      party_name,
      party_desc,
      party_type,
      level_p,
      parent_id,
      num_chd,
      leaf,
      tax_reference1,
      tax_reference2,
      public_key,
      address_id,
      registration_name,
      company_id,
      company_legal_form_code,
      company_legal_form,
      sole_proprietorship_indicator,
      company_liquidation_status_code,
      corporate_stock_amount,
      fully_paid_shares_indicator,
      corporate_registration_id,
      corporate_registration_name,
      corporate_registration_type_code,
      tax_level_code,
      exemption_reason_code,
      exemption_reason,
      tax_scheme_id,
      registration_date,
      registration_expiration_date,
      status_code,
      created_by_user_id,
      updated_by_user_id,
      created_at,
      updated_at from parties`

const selectChildPartiesSQL = `select 
	  p.id,
    p.uuid4,
    p.party_endpoint_id,
    p.party_endpoint_scheme_id,
    p.party_name,
    p.party_desc,
    p.party_type,
    p.level_p,
    p.parent_id,
    p.num_chd,
    p.leaf,
    p.tax_reference1,
    p.tax_reference2,
    p.public_key,
    p.address_id,
    p.registration_name,
    p.company_id,
    p.company_legal_form_code,
    p.company_legal_form,
    p.sole_proprietorship_indicator,
    p.company_liquidation_status_code,
    p.corporate_stock_amount,
    p.fully_paid_shares_indicator,
    p.corporate_registration_id,
    p.corporate_registration_name,
    p.corporate_registration_type_code,
    p.tax_level_code,
    p.exemption_reason_code,
    p.exemption_reason,
    p.tax_scheme_id,
    p.registration_date,
    p.registration_expiration_date,
    p.status_code,
    p.created_by_user_id,
    p.updated_by_user_id,
    p.created_at,
    p.updated_at from parties p inner join party_chds ch on (p.id = ch.party_chd_id)`

const insertTaxSchemeSQL = `insert into tax_schemes
	  (
    uuid4,
    ts_id,
    tax_scheme_name,
    tax_type_code,
    currency_code,
    status_code,
    created_by_user_id,
    updated_by_user_id,
    created_at,
    updated_at)
  values (:uuid4,
    :ts_id,
    :tax_scheme_name,
    :tax_type_code,
    :currency_code,
    :status_code,
    :created_by_user_id,
    :updated_by_user_id,
    :created_at,
    :updated_at);`

// CreateParty - Create Party
func (ps *PartyService) CreateParty(ctx context.Context, in *partyproto.CreatePartyRequest) (*partyproto.CreatePartyResponse, error) {
	user, err := GetUserWithNewContext(ctx, in.UserId, in.UserEmail, in.RequestId, ps.UserServiceClient)
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	ttime := common.GetTimeDetails()
	tn := common.TimeToTimestamp(ttime)

	partyD := partyproto.PartyD{}
	partyD.Uuid4, err = common.GetUUIDBytes()
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	partyD.PartyEndpointId = in.PartyEndpointId
	partyD.PartyEndpointSchemeId = in.PartyEndpointSchemeId
	partyD.PartyName = in.PartyName
	partyD.PartyDesc = in.PartyDesc
	partyD.PartyType = in.PartyType
	partyD.LevelP = 0
	partyD.ParentId = uint32(0)
	partyD.NumChd = 0
	partyD.Leaf = false
	partyD.TaxReference1 = in.TaxReference1
	partyD.TaxReference2 = in.TaxReference2
	partyD.PublicKey = in.PublicKey

	/*  PartyLegalEntity  */
	partyLegalEntity, err := ps.processPartyLegalEntity(ctx, in)
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	address, err := ps.processAddress(ctx, in)
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	crUpdUser := commonproto.CrUpdUser{}
	crUpdUser.StatusCode = "active"
	crUpdUser.CreatedByUserId = user.Id
	crUpdUser.UpdatedByUserId = user.Id

	crUpdTime := commonproto.CrUpdTime{}
	crUpdTime.CreatedAt = tn
	crUpdTime.UpdatedAt = tn

	taxScheme, err := ps.processTaxScheme(ctx, in, user.Id, tn)
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	party := partyproto.Party{PartyD: &partyD, PartyLegalEntityD: partyLegalEntity.PartyLegalEntityD, PartyLegalEntityT: partyLegalEntity.PartyLegalEntityT, CrUpdUser: &crUpdUser, CrUpdTime: &crUpdTime}

	err = ps.insertParty(ctx, insertPartySQL, insertTaxSchemeSQL, &party, address, taxScheme, in.GetUserEmail(), in.GetRequestId())
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	partyResponse := partyproto.CreatePartyResponse{}
	partyResponse.Party = &party
	return &partyResponse, nil
}

// insertParty - Insert party details into database
func (ps *PartyService) insertParty(ctx context.Context, insertPartySQL string, insertTaxSchemeSQL string, party *partyproto.Party, address *commonproto.Address, taxScheme *taxstruct.TaxScheme, userEmail string, requestID string) error {
	partyTmp, err := ps.crPartyStruct(ctx, party, userEmail, requestID)
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

		res1, err := tx.NamedExecContext(ctx, insertTaxSchemeSQL, taxScheme)
		if err != nil {
			ps.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}

		tID, err := res1.LastInsertId()
		if err != nil {
			ps.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}

		party.PartyLegalEntityD.TaxSchemeId = uint32(tID)
		partyTmp.PartyLegalEntityD.TaxSchemeId = uint32(tID)

		party.PartyD.AddressId = addr.Id
		partyTmp.PartyD.AddressId = addr.Id

		res, err := tx.NamedExecContext(ctx, insertPartySQL, partyTmp)
		if err != nil {
			ps.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}

		uID, err := res.LastInsertId()
		if err != nil {
			ps.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}
		party.PartyD.Id = uint32(uID)
		uuid4Str, err := common.UUIDBytesToStr(party.PartyD.Uuid4)
		if err != nil {
			ps.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}
		party.PartyD.IdS = uuid4Str
		return nil
	})

	if err != nil {
		ps.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
		return err
	}
	return nil
}

// crPartyStruct - process Party details
func (ps *PartyService) crPartyStruct(ctx context.Context, party *partyproto.Party, userEmail string, requestID string) (*partystruct.Party, error) {
	partyLegalEntityT := new(commonstruct.PartyLegalEntityT)
	partyLegalEntityT.RegistrationDate = common.TimestampToTime(party.PartyLegalEntityT.RegistrationDate)
	partyLegalEntityT.RegistrationExpirationDate = common.TimestampToTime(party.PartyLegalEntityT.RegistrationExpirationDate)

	crUpdTime := new(commonstruct.CrUpdTime)
	crUpdTime.CreatedAt = common.TimestampToTime(party.CrUpdTime.CreatedAt)
	crUpdTime.UpdatedAt = common.TimestampToTime(party.CrUpdTime.UpdatedAt)

	partyTmp := partystruct.Party{PartyD: party.PartyD, PartyLegalEntityD: party.PartyLegalEntityD, PartyLegalEntityT: partyLegalEntityT, CrUpdUser: party.CrUpdUser, CrUpdTime: crUpdTime}

	return &partyTmp, nil
}

// CreateChild - Create Child Party
func (ps *PartyService) CreateChild(ctx context.Context, inReq *partyproto.CreateChildRequest) (*partyproto.CreateChildResponse, error) {
	in := inReq.CreatePartyRequest
	db := ps.DBService.DB
	updateNumChildrenStmt, err := db.PreparexContext(ctx, updateNumChildrenSQL)
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	user, err := GetUserWithNewContext(ctx, in.UserId, in.UserEmail, in.RequestId, ps.UserServiceClient)
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	getByIDRequest := commonproto.GetByIdRequest{}
	getByIDRequest.Id = in.ParentId
	getByIDRequest.UserEmail = in.GetUserEmail()
	getByIDRequest.RequestId = in.GetRequestId()

	form := partyproto.GetPartyByPkRequest{}
	form.GetByIdRequest = &getByIDRequest
	parentPartyResponse, err := ps.GetPartyByPk(ctx, &form)
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	parent := parentPartyResponse.Party

	ttime := common.GetTimeDetails()
	tn := common.TimeToTimestamp(ttime)

	partyD := partyproto.PartyD{}
	partyD.Uuid4, err = common.GetUUIDBytes()
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	partyD.PartyEndpointId = in.PartyEndpointId
	partyD.PartyEndpointSchemeId = in.PartyEndpointSchemeId
	partyD.PartyName = in.PartyName
	partyD.PartyDesc = in.PartyDesc
	partyD.PartyType = in.PartyType
	partyD.LevelP = parent.PartyD.LevelP + 1
	partyD.ParentId = parent.PartyD.Id
	partyD.NumChd = 0
	partyD.Leaf = false
	partyD.TaxReference1 = in.TaxReference1
	partyD.TaxReference2 = in.TaxReference2
	partyD.PublicKey = in.PublicKey

	partyLegalEntity, err := ps.processPartyLegalEntity(ctx, in)
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	address, err := ps.processAddress(ctx, in)
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	taxScheme, err := ps.processTaxScheme(ctx, in, user.Id, tn)
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	crUpdUser := commonproto.CrUpdUser{}
	crUpdUser.StatusCode = "active"
	crUpdUser.CreatedByUserId = user.Id
	crUpdUser.UpdatedByUserId = user.Id

	crUpdTime := commonproto.CrUpdTime{}
	crUpdTime.CreatedAt = tn
	crUpdTime.UpdatedAt = tn

	party := partyproto.Party{PartyD: &partyD, PartyLegalEntityD: partyLegalEntity.PartyLegalEntityD, PartyLegalEntityT: partyLegalEntity.PartyLegalEntityT, CrUpdUser: &crUpdUser, CrUpdTime: &crUpdTime}

	err = ps.insertChild(ctx, insertPartySQL, insertChildPartySQL, updateNumChildrenSQL, updateNumChildrenStmt, insertTaxSchemeSQL, parent, &party, address, taxScheme, in.GetUserEmail(), in.GetRequestId())
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	partyResponse := partyproto.CreateChildResponse{}
	partyResponse.Party = &party
	return &partyResponse, nil
}

// insertChild - insert Child Party
func (ps *PartyService) insertChild(ctx context.Context, insertPartySQL string, insertChildPartySQL string, updateNumChildrenSQL string, updateNumChildrenStmt *sqlx.Stmt, insertTaxSchemeSQL string, parent *partyproto.Party, party *partyproto.Party, address *commonproto.Address, taxScheme *taxstruct.TaxScheme, userEmail string, requestID string) error {
	partyTmp, err := ps.crPartyStruct(ctx, party, userEmail, requestID)
	if err != nil {
		ps.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
		return err
	}

	partyChdD := partyproto.PartyChdD{}
	partyChdD.Uuid4, err = common.GetUUIDBytes()
	if err != nil {
		ps.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
		return err
	}
	partyChdD.PartyId = parent.PartyD.Id

	err = ps.DBService.InsUpd(ctx, userEmail, requestID, func(tx *sqlx.Tx) error {
		addr, err := common.InsertAddress(ctx, tx, address, userEmail, requestID)
		if err != nil {
			ps.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}
		res1, err := tx.NamedExecContext(ctx, insertTaxSchemeSQL, taxScheme)
		if err != nil {
			ps.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}

		tID, err := res1.LastInsertId()
		if err != nil {
			ps.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}
		party.PartyLegalEntityD.TaxSchemeId = uint32(tID)
		partyTmp.PartyLegalEntityD.TaxSchemeId = uint32(tID)

		party.PartyD.AddressId = addr.Id
		partyTmp.PartyD.AddressId = addr.Id
		res, err := tx.NamedExecContext(ctx, insertPartySQL, partyTmp)
		if err != nil {
			ps.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}

		uID, err := res.LastInsertId()
		if err != nil {
			ps.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}
		party.PartyD.Id = uint32(uID)
		uuid4Str, err := common.UUIDBytesToStr(party.PartyD.Uuid4)
		if err != nil {
			ps.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}
		party.PartyD.IdS = uuid4Str

		partyChdD.PartyChdId = party.PartyD.Id

		partyChdTmp := partystruct.PartyChd{PartyChdD: &partyChdD, CrUpdUser: partyTmp.CrUpdUser, CrUpdTime: partyTmp.CrUpdTime}

		_, err = tx.NamedExecContext(ctx, insertChildPartySQL, partyChdTmp)

		if err != nil {
			ps.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}

		numchd := parent.PartyD.NumChd + 1
		tn1 := common.GetTimeDetails()
		_, err = tx.StmtxContext(ctx, updateNumChildrenStmt).ExecContext(ctx,
			numchd,
			true,
			tn1,
			parent.PartyD.Id,
			"active")

		if err != nil {
			ps.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}
		return nil
	})
	if err != nil {
		ps.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
		return err
	}
	return nil
}

// GetParties - Get Parties
func (ps *PartyService) GetParties(ctx context.Context, in *partyproto.GetPartiesRequest) (*partyproto.GetPartiesResponse, error) {
	limit := in.GetLimit()
	nextCursor := in.GetNextCursor()
	if limit == "" {
		limit = ps.DBService.LimitSQLRows
	}
	query := "level_p = ? and status_code = ?"
	if nextCursor == "" {
		query = query + " order by id desc " + " limit " + limit + ";"
	} else {
		nextCursor = common.DecodeCursor(nextCursor)
		query = query + " " + "and" + " " + "id <= " + nextCursor + " order by id desc " + " limit " + limit + ";"
	}

	parties := []*partyproto.Party{}

	nselectPartiesSQL := selectPartiesSQL + ` where ` + query

	rows, err := ps.DBService.DB.QueryxContext(ctx, nselectPartiesSQL, 0, "active")
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	for rows.Next() {

		partyTmp := partystruct.Party{}
		err = rows.StructScan(&partyTmp)
		if err != nil {
			ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
			return nil, err
		}
		getRequest := commonproto.GetRequest{}
		getRequest.UserEmail = in.UserEmail
		getRequest.RequestId = in.RequestId
		party, err := ps.getPartyStruct(ctx, &getRequest, partyTmp)
		if err != nil {
			ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
			return nil, err
		}
		parties = append(parties, party)

	}

	partiesResponse := partyproto.GetPartiesResponse{}
	if len(parties) != 0 {
		next := parties[len(parties)-1].PartyD.Id
		next--
		nextc := common.EncodeCursor(next)
		partiesResponse = partyproto.GetPartiesResponse{Parties: parties, NextCursor: nextc}
	} else {
		partiesResponse = partyproto.GetPartiesResponse{Parties: parties, NextCursor: "0"}
	}
	return &partiesResponse, nil
}

// GetParty - Get Party
func (ps *PartyService) GetParty(ctx context.Context, inReq *partyproto.GetPartyRequest) (*partyproto.GetPartyResponse, error) {
	in := inReq.GetRequest
	uuid4byte, err := common.UUIDStrToBytes(in.Id)
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	nselectPartiesSQL := selectPartiesSQL + ` where uuid4 = ? and status_code = ?;`
	row := ps.DBService.DB.QueryRowxContext(ctx, nselectPartiesSQL, uuid4byte, "active")

	partyTmp := partystruct.Party{}
	err = row.StructScan(&partyTmp)
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	party, err := ps.getPartyStruct(ctx, in, partyTmp)
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	partyResponse := partyproto.GetPartyResponse{}
	partyResponse.Party = party
	return &partyResponse, nil
}

// GetPartyByPk - Get Party By Primary key(Id)
func (ps *PartyService) GetPartyByPk(ctx context.Context, inReq *partyproto.GetPartyByPkRequest) (*partyproto.GetPartyByPkResponse, error) {
	in := inReq.GetByIdRequest
	nselectPartiesSQL := selectPartiesSQL + ` where id = ? and status_code = ?;`
	row := ps.DBService.DB.QueryRowxContext(ctx, nselectPartiesSQL, in.Id, "active")

	partyTmp := partystruct.Party{}
	err := row.StructScan(&partyTmp)
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	getRequest := commonproto.GetRequest{}
	getRequest.UserEmail = in.UserEmail
	getRequest.RequestId = in.RequestId
	party, err := ps.getPartyStruct(ctx, &getRequest, partyTmp)
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	partyResponse := partyproto.GetPartyByPkResponse{}
	partyResponse.Party = party
	return &partyResponse, nil
}

// GetTopLevelParties - Get top level parties
func (ps *PartyService) GetTopLevelParties(ctx context.Context, in *partyproto.GetTopLevelPartiesRequest) (*partyproto.GetTopLevelPartiesResponse, error) {
	parties := []*partyproto.Party{}

	nselectPartiesSQL := selectPartiesSQL + ` where level_p = ? and status_code = ?;`

	rows, err := ps.DBService.DB.QueryxContext(ctx, nselectPartiesSQL, 0, "active")
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	for rows.Next() {

		partyTmp := partystruct.Party{}
		err = rows.StructScan(&partyTmp)
		if err != nil {
			ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
			return nil, err
		}
		getRequest := commonproto.GetRequest{}
		getRequest.UserEmail = in.UserEmail
		getRequest.RequestId = in.RequestId
		party, err := ps.getPartyStruct(ctx, &getRequest, partyTmp)
		if err != nil {
			ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
			return nil, err
		}
		parties = append(parties, party)
	}
	partiesResponse := partyproto.GetTopLevelPartiesResponse{}
	partiesResponse.Parties = parties
	return &partiesResponse, nil
}

// GetChildParties - Get child parties
func (ps *PartyService) GetChildParties(ctx context.Context, inReq *partyproto.GetChildPartiesRequest) (*partyproto.GetChildPartiesResponse, error) {
	in := inReq.GetRequest
	form := partyproto.GetPartyRequest{}
	form.GetRequest = in
	partyResponse, err := ps.GetParty(ctx, &form)
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	party := partyResponse.Party
	parties := []*partyproto.Party{}

	nselectChildPartiesSQL := selectChildPartiesSQL + ` where ch.party_id = ?`

	rows, err := ps.DBService.DB.QueryxContext(ctx, nselectChildPartiesSQL, party.PartyD.Id)
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	for rows.Next() {
		partyTmp := partystruct.Party{}
		err = rows.StructScan(&partyTmp)
		if err != nil {
			ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
			return nil, err
		}
		party, err := ps.getPartyStruct(ctx, in, partyTmp)
		if err != nil {
			ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
			return nil, err
		}

		parties = append(parties, party)
	}

	partiesResponse := partyproto.GetChildPartiesResponse{}
	partiesResponse.Parties = parties
	return &partiesResponse, nil
}

// GetParentParty - Get Parent Party
func (ps *PartyService) GetParentParty(ctx context.Context, inReq *partyproto.GetParentPartyRequest) (*partyproto.GetParentPartyResponse, error) {
	in := inReq.GetRequest
	form := partyproto.GetPartyRequest{}
	form.GetRequest = in
	partyResponse, err := ps.GetParty(ctx, &form)
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	childparty := partyResponse.Party

	nselectPartiesSQL := selectPartiesSQL + ` where id = ? and status_code = ?;`
	row := ps.DBService.DB.QueryRowxContext(ctx, nselectPartiesSQL, childparty.PartyD.ParentId, "active")

	partyTmp := partystruct.Party{}
	err = row.StructScan(&partyTmp)
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	party, err := ps.getPartyStruct(ctx, in, partyTmp)
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	parentPartyResponse := partyproto.GetParentPartyResponse{}
	parentPartyResponse.Party = party
	return &parentPartyResponse, nil
}

// getPartyStruct - Get Party Struct
func (ps *PartyService) getPartyStruct(ctx context.Context, in *commonproto.GetRequest, partyTmp partystruct.Party) (*partyproto.Party, error) {
	uuid4Str, err := common.UUIDBytesToStr(partyTmp.PartyD.Uuid4)
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	partyTmp.PartyD.IdS = uuid4Str

	crUpdTime := new(commonproto.CrUpdTime)
	crUpdTime.CreatedAt = common.TimeToTimestamp(partyTmp.CrUpdTime.CreatedAt)
	crUpdTime.UpdatedAt = common.TimeToTimestamp(partyTmp.CrUpdTime.UpdatedAt)

	partyLegalEntityT := new(commonproto.PartyLegalEntityT)
	partyLegalEntityT.RegistrationDate = common.TimeToTimestamp(partyTmp.PartyLegalEntityT.RegistrationDate)
	partyLegalEntityT.RegistrationExpirationDate = common.TimeToTimestamp(partyTmp.PartyLegalEntityT.RegistrationExpirationDate)

	party := partyproto.Party{PartyD: partyTmp.PartyD, PartyLegalEntityD: partyTmp.PartyLegalEntityD, PartyLegalEntityT: partyLegalEntityT, CrUpdUser: partyTmp.CrUpdUser, CrUpdTime: crUpdTime}

	return &party, nil
}

// UpdateParty - Update party
func (ps *PartyService) UpdateParty(ctx context.Context, in *partyproto.UpdatePartyRequest) (*partyproto.UpdatePartyResponse, error) {
	getRequest := commonproto.GetRequest{}
	getRequest.Id = in.Id
	getRequest.UserEmail = in.UserEmail
	getRequest.RequestId = in.RequestId
	form := partyproto.GetPartyRequest{}
	form.GetRequest = &getRequest
	partyResponse, err := ps.GetParty(ctx, &form)
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	party := partyResponse.Party
	db := ps.DBService.DB
	tn := common.GetTimeDetails()
	stmt, err := db.PreparexContext(ctx, updatePartySQL)
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	err = ps.DBService.InsUpd(ctx, in.GetUserEmail(), in.GetRequestId(), func(tx *sqlx.Tx) error {
		_, err = tx.StmtxContext(ctx, stmt).ExecContext(ctx,
			in.PartyName,
			in.PartyDesc,
			tn,
			party.PartyD.Id,
			"active")
		if err != nil {
			ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
			err1 := stmt.Close()
			if err1 != nil {
				ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err1))
				return err1
			}
			return err
		}
		return nil
	})

	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	err = stmt.Close()
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	return &partyproto.UpdatePartyResponse{}, nil
}

// DeleteParty - Delete party
func (ps *PartyService) DeleteParty(ctx context.Context, inReq *partyproto.DeletePartyRequest) (*partyproto.DeletePartyResponse, error) {
	in := inReq.GetRequest
	uuid4byte, err := common.UUIDStrToBytes(in.Id)
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	db := ps.DBService.DB
	tn := common.GetTimeDetails()
	stmt, err := db.PreparexContext(ctx, deletePartySQL)
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	err = ps.DBService.InsUpd(ctx, in.GetUserEmail(), in.GetRequestId(), func(tx *sqlx.Tx) error {
		_, err = tx.StmtxContext(ctx, stmt).ExecContext(ctx,
			"inactive",
			tn,
			uuid4byte)

		if err != nil {
			ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
			err1 := stmt.Close()
			if err != nil {
				ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err1))
				return err1
			}
			return err
		}
		return nil
	})

	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	err = stmt.Close()
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	return &partyproto.DeletePartyResponse{}, nil
}

// processAddress - process Address
func (ps *PartyService) processAddress(ctx context.Context, in *partyproto.CreatePartyRequest) (*commonproto.Address, error) {
	var err error
	address := commonproto.Address{}
	address.Uuid4, err = common.GetUUIDBytes()
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	/*  Addr  */
	address.Name1 = in.Name1
	address.AddrListAgencyId = in.AddrListAgencyId
	address.AddrListId = in.AddrListAgencyId
	address.AddrListVersionId = in.AddrListAgencyId
	address.AddressTypeCode = in.AddressTypeCode
	address.AddressFormatCode = in.AddressFormatCode
	address.Postbox = in.Postbox
	address.Floor1 = in.Floor1
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
	return &address, nil
}

// processPartyLegalEntity - process Party Legal Entity
func (ps *PartyService) processPartyLegalEntity(ctx context.Context, in *partyproto.CreatePartyRequest) (*commonproto.PartyLegalEntity, error) {
	registrationDate, err := time.Parse(common.Layout, in.RegistrationDate)
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	registrationExpirationDate, err := time.Parse(common.Layout, in.RegistrationExpirationDate)
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	partyLegalEntityD := commonproto.PartyLegalEntityD{}
	partyLegalEntityD.RegistrationName = in.RegistrationName
	partyLegalEntityD.CompanyId = in.CompanyId
	partyLegalEntityD.CompanyLegalFormCode = in.CompanyLegalFormCode
	partyLegalEntityD.CompanyLegalForm = in.CompanyLegalForm
	partyLegalEntityD.SoleProprietorshipIndicator = in.SoleProprietorshipIndicator
	partyLegalEntityD.CompanyLiquidationStatusCode = in.CompanyLiquidationStatusCode
	partyLegalEntityD.CorporateStockAmount = in.CorporateStockAmount
	partyLegalEntityD.FullyPaidSharesIndicator = in.FullyPaidSharesIndicator
	partyLegalEntityD.CorporateRegistrationId = in.CorporateRegistrationId
	partyLegalEntityD.CorporateRegistrationName = in.CorporateRegistrationName
	partyLegalEntityD.CorporateRegistrationTypeCode = in.CorporateRegistrationTypeCode
	partyLegalEntityD.TaxLevelCode = in.TaxLevelCode
	partyLegalEntityD.ExemptionReasonCode = in.ExemptionReasonCode
	partyLegalEntityD.ExemptionReason = in.ExemptionReason

	partyLegalEntityT := commonproto.PartyLegalEntityT{}
	partyLegalEntityT.RegistrationDate = common.TimeToTimestamp(registrationDate.UTC().Truncate(time.Second))
	partyLegalEntityT.RegistrationExpirationDate = common.TimeToTimestamp(registrationExpirationDate.UTC().Truncate(time.Second))

	partyLegalEntityTmp := commonproto.PartyLegalEntity{PartyLegalEntityD: &partyLegalEntityD, PartyLegalEntityT: &partyLegalEntityT}

	return &partyLegalEntityTmp, nil
}

// processTaxScheme - process Tax Scheme
func (ps *PartyService) processTaxScheme(ctx context.Context, in *partyproto.CreatePartyRequest, userId string, tn *timestamp.Timestamp) (*taxstruct.TaxScheme, error) {
	var err error
	taxSchemeD := taxproto.TaxSchemeD{}
	taxSchemeD.Uuid4, err = common.GetUUIDBytes()
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	taxSchemeD.TsId = in.TsId
	taxSchemeD.TaxSchemeName = in.TaxSchemeName
	taxSchemeD.TaxTypeCode = in.TaxTypeCode
	taxSchemeD.CurrencyCode = in.CurrencyCode

	crUpdUser := commonproto.CrUpdUser{}
	crUpdUser.StatusCode = "active"
	crUpdUser.CreatedByUserId = userId
	crUpdUser.UpdatedByUserId = userId

	crUpdTime := commonstruct.CrUpdTime{}
	crUpdTime.CreatedAt = common.TimestampToTime(tn)
	crUpdTime.UpdatedAt = common.TimestampToTime(tn)

	taxSchemeTmp := taxstruct.TaxScheme{TaxSchemeD: &taxSchemeD, CrUpdUser: &crUpdUser, CrUpdTime: &crUpdTime}

	return &taxSchemeTmp, nil
}
