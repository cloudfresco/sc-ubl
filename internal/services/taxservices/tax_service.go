package taxservices

import (
	"context"
	"net"
	"os"

	"github.com/cloudfresco/sc-ubl/internal/common"
	"github.com/cloudfresco/sc-ubl/internal/config"
	commonproto "github.com/cloudfresco/sc-ubl/internal/protogen/common/v1"
	partyproto "github.com/cloudfresco/sc-ubl/internal/protogen/party/v1"
	taxproto "github.com/cloudfresco/sc-ubl/internal/protogen/tax/v1"
	partyservice "github.com/cloudfresco/sc-ubl/internal/services/partyservices"
	commonstruct "github.com/cloudfresco/sc-ubl/internal/servicestructs/common/v1"
	taxstruct "github.com/cloudfresco/sc-ubl/internal/servicestructs/tax/v1"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// TaxService - For accessing Tax services
type TaxService struct {
	log               *zap.Logger
	DBService         *common.DBService
	RedisService      *common.RedisService
	UserServiceClient partyproto.UserServiceClient
	taxproto.UnimplementedTaxServiceServer
}

// NewTaxService - Create Tax service
func NewTaxService(log *zap.Logger, dbOpt *common.DBService, redisOpt *common.RedisService, userServiceClient partyproto.UserServiceClient) *TaxService {
	return &TaxService{
		log:               log,
		DBService:         dbOpt,
		RedisService:      redisOpt,
		UserServiceClient: userServiceClient,
	}
}

// StartTaxServer - Start Tax server
func StartTaxServer(log *zap.Logger, isTest bool, pwd string, dbOpt *config.DBOptions, redisOpt *config.RedisOptions, mailerOpt *config.MailerOptions, grpcServerOpt *config.GrpcServerOptions, jwtOpt *config.JWTOptions, oauthOpt *config.OauthOptions, userOpt *config.UserOptions, uptraceOpt *config.UptraceOptions, dbService *common.DBService, redisService *common.RedisService, mailerService common.MailerIntf) {
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
	taxService := NewTaxService(log, dbService, redisService, uc)
	lis, err := net.Listen("tcp", grpcServerOpt.GrpcTaxServerPort)
	if err != nil {
		log.Error("Error", zap.Error(err))
		os.Exit(1)
	}

	srv := grpc.NewServer(srvOpts...)
	taxproto.RegisterTaxServiceServer(srv, taxService)

	if err := srv.Serve(lis); err != nil {
		log.Error("Error", zap.Error(err))
		os.Exit(1)
	}
}

const selectTaxCategoriesSQL = `select 
  id,
  uuid4,
  tc_id,
  tax_category_name,
  percent,
  base_unit_measure,
  per_unit_amount,
  tax_exemption_reason_code,
  tax_exemption_reason,
  tier_range,
  tier_rate_percent,
  tax_scheme_id,
  status_code,
  created_by_user_id,
  updated_by_user_id,
  created_at,
  updated_at from tax_categories`

const insertTaxCategorySQL = `insert into tax_categories
	  (
uuid4,
tc_id,
tax_category_name,
percent,
base_unit_measure,
per_unit_amount,
tax_exemption_reason_code,
tax_exemption_reason,
tier_range,
tier_rate_percent,
tax_scheme_id,
status_code,
created_by_user_id,
updated_by_user_id,
created_at,
updated_at)
  values (:uuid4,
:tc_id,
:tax_category_name,
:percent,
:base_unit_measure,
:per_unit_amount,
:tax_exemption_reason_code,
:tax_exemption_reason,
:tier_range,
:tier_rate_percent,
:tax_scheme_id,
:status_code,
:created_by_user_id,
:updated_by_user_id,
:created_at,
:updated_at);`

// updateTaxCategorySQL - update TaxCategorySQL query
const updateTaxCategorySQL = `update tax_categories set 
  tax_category_name = ?,
  percent = ?,
  base_unit_measure = ?,
  per_unit_amount = ?,
  tax_exemption_reason_code = ?,
  tax_exemption_reason = ?,
  updated_at = ? where uuid4 = ?;`

const selectTaxSchemesSQL = `select 
  id,
  uuid4,
  ts_id,
  tax_scheme_name,
  tax_type_code,
  currency_code,
  status_code,
  created_by_user_id,
  updated_by_user_id,
  created_at,
  updated_at from tax_schemes`

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

// updateTaxSchemeSQL - update TaxSchemeSQL query
const updateTaxSchemeSQL = `update tax_schemes set 
  tax_scheme_name = ?,
  tax_type_code = ?,
  currency_code = ?,
  updated_at = ? where uuid4 = ?;`

const insertTaxSchemeJurisdictionSQL = `insert into tax_scheme_jurisdictions
	  (
uuid4,
address_id,
tax_scheme_id,
tax_scheme_name,
status_code,
created_by_user_id,
updated_by_user_id,
created_at,
updated_at)
  values (:uuid4,
:address_id,
:tax_scheme_id,
:tax_scheme_name,
:status_code,
:created_by_user_id,
:updated_by_user_id,
:created_at,
:updated_at);`

// updateTaxSchemeJurisdictionSQL - update TaxSchemeJurisdictionSQL query
const updateTaxSchemeJurisdictionSQL = `update tax_scheme_jurisdictions set 
  tax_scheme_name = ?,
  tax_scheme_id = ?,
  updated_at = ? where uuid4 = ?;`

const insertTaxSubTotalSQL = `insert into tax_sub_totals
	  (
uuid4,
taxable_amount,
tax_amount,
calculation_sequence_numeric,
transaction_currency_tax_amount,
percent,
base_unit_measure,
per_unit_amount,
tier_range,
tier_rate_percent,
tax_category_id,
tax_total_id,
status_code,
created_by_user_id,
updated_by_user_id,
created_at,
updated_at)
  values (:uuid4,
:taxable_amount,
:tax_amount,
:calculation_sequence_numeric,
:transaction_currency_tax_amount,
:percent,
:base_unit_measure,
:per_unit_amount,
:tier_range,
:tier_rate_percent,
:tax_category_id,
:tax_total_id,
:status_code,
:created_by_user_id,
:updated_by_user_id,
:created_at,
:updated_at);`

// updateTaxSubTotalSQL - update TaxSubTotalSQL query
const updateTaxSubTotalSQL = `update tax_sub_totals set 
  taxable_amount = ?,
  tax_amount = ?,
  calculation_sequence_numeric = ?,
  transaction_currency_tax_amount = ?,
  percent = ?,
  base_unit_measure = ?,
  per_unit_amount = ?,
  updated_at = ? where uuid4 = ?;`

const selectTaxTotalsSQL = `select 
  id,
  uuid4,
  tax_amount,
  rounding_amount,
  tax_evidence_indicator,
  tax_included_indicator,
  master_flag,
  master_id,
  tax_category_id,
  status_code,
  created_by_user_id,
  updated_by_user_id,
  created_at,
  updated_at from tax_totals`

const insertTaxTotalSQL = `insert into tax_totals
	  (uuid4,
tax_amount,
rounding_amount,
tax_evidence_indicator,
tax_included_indicator,
master_flag,
master_id,
tax_category_id,
status_code,
created_by_user_id,
updated_by_user_id,
created_at,
updated_at)
  values (:uuid4,
:tax_amount,
:rounding_amount,
:tax_evidence_indicator,
:tax_included_indicator,
:master_flag,
:master_id,
:tax_category_id,
:status_code,
:created_by_user_id,
:updated_by_user_id,
:created_at,
:updated_at);`

// updateTaxTotalSQL - update TaxTotalSQL query
const updateTaxTotalSQL = `update tax_totals set 
  tax_amount = ?,
  rounding_amount = ?,
  updated_at = ? where uuid4 = ?;`

// GetTaxSchemes - Get TaxSchemes
func (ts *TaxService) GetTaxSchemes(ctx context.Context, in *taxproto.GetTaxSchemesRequest) (*taxproto.GetTaxSchemesResponse, error) {
	limit := in.GetLimit()
	nextCursor := in.GetNextCursor()
	if limit == "" {
		limit = ts.DBService.LimitSQLRows
	}
	query := "status_code = ?"
	if nextCursor == "" {
		query = query + " order by id desc " + " limit " + limit + ";"
	} else {
		nextCursor = common.DecodeCursor(nextCursor)
		query = query + " " + "and" + " " + "id <= " + nextCursor + " order by id desc " + " limit " + limit + ";"
	}

	taxSchemes := []*taxproto.TaxScheme{}

	nselectTaxSchemesSQL := selectTaxSchemesSQL + ` where ` + query

	rows, err := ts.DBService.DB.QueryxContext(ctx, nselectTaxSchemesSQL, "active")
	if err != nil {
		ts.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	for rows.Next() {

		taxSchemeTmp := taxstruct.TaxScheme{}
		err = rows.StructScan(&taxSchemeTmp)
		if err != nil {
			ts.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
			return nil, err
		}

		getRequest := commonproto.GetRequest{}
		getRequest.UserEmail = in.UserEmail
		getRequest.RequestId = in.RequestId
		taxScheme, err := ts.getTaxSchemeStruct(ctx, &getRequest, &taxSchemeTmp)
		if err != nil {
			ts.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
			return nil, err
		}
		taxSchemes = append(taxSchemes, taxScheme)

	}

	taxSchemesResponse := taxproto.GetTaxSchemesResponse{}
	if len(taxSchemes) != 0 {
		next := taxSchemes[len(taxSchemes)-1].TaxSchemeD.Id
		next--
		nextc := common.EncodeCursor(next)
		taxSchemesResponse = taxproto.GetTaxSchemesResponse{TaxSchemes: taxSchemes, NextCursor: nextc}
	} else {
		taxSchemesResponse = taxproto.GetTaxSchemesResponse{TaxSchemes: taxSchemes, NextCursor: "0"}
	}
	return &taxSchemesResponse, nil
}

// GetTaxScheme - Get TaxScheme
func (ts *TaxService) GetTaxScheme(ctx context.Context, inReq *taxproto.GetTaxSchemeRequest) (*taxproto.GetTaxSchemeResponse, error) {
	in := inReq.GetRequest
	uuid4byte, err := common.UUIDStrToBytes(in.Id)
	if err != nil {
		ts.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	nselectTaxSchemesSQL := selectTaxSchemesSQL + ` where uuid4 = ? and status_code = ?;`
	row := ts.DBService.DB.QueryRowxContext(ctx, nselectTaxSchemesSQL, uuid4byte, "active")
	taxSchemeTmp := taxstruct.TaxScheme{}
	err = row.StructScan(&taxSchemeTmp)
	if err != nil {
		ts.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	taxScheme, err := ts.getTaxSchemeStruct(ctx, in, &taxSchemeTmp)
	if err != nil {
		ts.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	taxSchemeResponse := taxproto.GetTaxSchemeResponse{}
	taxSchemeResponse.TaxScheme = taxScheme
	return &taxSchemeResponse, nil
}

// getTaxSchemeStruct - Get TaxScheme
func (ts *TaxService) getTaxSchemeStruct(ctx context.Context, in *commonproto.GetRequest, taxSchemeTmp *taxstruct.TaxScheme) (*taxproto.TaxScheme, error) {
	crUpdTime := new(commonproto.CrUpdTime)
	crUpdTime.CreatedAt = common.TimeToTimestamp(taxSchemeTmp.CrUpdTime.CreatedAt)
	crUpdTime.UpdatedAt = common.TimeToTimestamp(taxSchemeTmp.CrUpdTime.UpdatedAt)

	uuid4Str, err := common.UUIDBytesToStr(taxSchemeTmp.TaxSchemeD.Uuid4)
	if err != nil {
		ts.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	taxSchemeTmp.TaxSchemeD.IdS = uuid4Str

	taxScheme := taxproto.TaxScheme{TaxSchemeD: taxSchemeTmp.TaxSchemeD, CrUpdUser: taxSchemeTmp.CrUpdUser, CrUpdTime: crUpdTime}

	return &taxScheme, nil
}

// CreateTaxScheme - Create TaxScheme
func (ts *TaxService) CreateTaxScheme(ctx context.Context, in *taxproto.CreateTaxSchemeRequest) (*taxproto.CreateTaxSchemeResponse, error) {
	user, err := partyservice.GetUserWithNewContext(ctx, in.UserId, in.UserEmail, in.RequestId, ts.UserServiceClient)
	if err != nil {
		ts.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	ttime := common.GetTimeDetails()
	tn := common.TimeToTimestamp(ttime)

	taxSchemeD := taxproto.TaxSchemeD{}
	taxSchemeD.Uuid4, err = common.GetUUIDBytes()
	if err != nil {
		ts.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	taxSchemeD.TsId = in.TsId
	taxSchemeD.TaxSchemeName = in.TaxSchemeName
	taxSchemeD.TaxTypeCode = in.TaxTypeCode
	taxSchemeD.CurrencyCode = in.CurrencyCode

	crUpdUser := commonproto.CrUpdUser{}
	crUpdUser.StatusCode = "active"
	crUpdUser.CreatedByUserId = user.Id
	crUpdUser.UpdatedByUserId = user.Id

	crUpdTime := commonproto.CrUpdTime{}
	crUpdTime.CreatedAt = tn
	crUpdTime.UpdatedAt = tn

	taxScheme := taxproto.TaxScheme{TaxSchemeD: &taxSchemeD, CrUpdUser: &crUpdUser, CrUpdTime: &crUpdTime}

	err = ts.insertTaxScheme(ctx, insertTaxSchemeSQL, &taxScheme, in.GetUserEmail(), in.GetRequestId())
	if err != nil {
		ts.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))

		return nil, err
	}

	taxSchemeResponse := taxproto.CreateTaxSchemeResponse{}
	taxSchemeResponse.TaxScheme = &taxScheme
	return &taxSchemeResponse, nil
}

// insertTaxScheme - Insert tax scheme details into database
func (ts *TaxService) insertTaxScheme(ctx context.Context, insertTaxSchemeSQL string, taxScheme *taxproto.TaxScheme, userEmail string, requestID string) error {
	taxSchemeTmp, err := ts.crTaxSchemeStruct(ctx, taxScheme, userEmail, requestID)
	if err != nil {
		ts.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
		return err
	}

	err = ts.DBService.InsUpd(ctx, userEmail, requestID, func(tx *sqlx.Tx) error {
		res, err := tx.NamedExecContext(ctx, insertTaxSchemeSQL, taxSchemeTmp)
		if err != nil {
			ts.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}

		uID, err := res.LastInsertId()
		if err != nil {
			ts.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}
		taxScheme.TaxSchemeD.Id = uint32(uID)
		uuid4Str, err := common.UUIDBytesToStr(taxScheme.TaxSchemeD.Uuid4)
		if err != nil {
			ts.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}
		taxScheme.TaxSchemeD.IdS = uuid4Str
		return nil
	})

	if err != nil {
		ts.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
		return err
	}
	return nil
}

// crTaxSchemeStruct - process TaxScheme details
func (ts *TaxService) crTaxSchemeStruct(ctx context.Context, taxScheme *taxproto.TaxScheme, userEmail string, requestID string) (*taxstruct.TaxScheme, error) {
	crUpdTime := new(commonstruct.CrUpdTime)
	crUpdTime.CreatedAt = common.TimestampToTime(taxScheme.CrUpdTime.CreatedAt)
	crUpdTime.UpdatedAt = common.TimestampToTime(taxScheme.CrUpdTime.UpdatedAt)

	taxSchemeTmp := taxstruct.TaxScheme{TaxSchemeD: taxScheme.TaxSchemeD, CrUpdUser: taxScheme.CrUpdUser, CrUpdTime: crUpdTime}

	return &taxSchemeTmp, nil
}

// UpdateTaxScheme - Update TaxScheme
func (ts *TaxService) UpdateTaxScheme(ctx context.Context, in *taxproto.UpdateTaxSchemeRequest) (*taxproto.UpdateTaxSchemeResponse, error) {
	db := ts.DBService.DB
	tn := common.GetTimeDetails()

	uuid4byte, err := common.UUIDStrToBytes(in.Id)
	if err != nil {
		ts.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	stmt, err := db.PreparexContext(ctx, updateTaxSchemeSQL)
	if err != nil {
		ts.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	err = ts.DBService.InsUpd(ctx, in.GetUserEmail(), in.GetRequestId(), func(tx *sqlx.Tx) error {
		_, err = tx.StmtxContext(ctx, stmt).ExecContext(ctx,
			in.TaxSchemeName,
			in.TaxTypeCode,
			in.CurrencyCode,
			tn,
			uuid4byte)
		if err != nil {
			ts.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
			err1 := stmt.Close()
			if err1 != nil {
				ts.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err1))
				return err1
			}
			return err
		}
		return nil
	})

	if err != nil {
		ts.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	err = stmt.Close()
	if err != nil {
		ts.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	return &taxproto.UpdateTaxSchemeResponse{}, nil
}

// GetTaxCategory - Get TaxCategory
func (ts *TaxService) GetTaxCategory(ctx context.Context, inReq *taxproto.GetTaxCategoryRequest) (*taxproto.GetTaxCategoryResponse, error) {
	in := inReq.GetRequest
	uuid4byte, err := common.UUIDStrToBytes(in.Id)
	if err != nil {
		ts.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	nselectTaxCategoriesSQL := selectTaxCategoriesSQL + ` where uuid4 = ? and status_code = ?;`
	row := ts.DBService.DB.QueryRowxContext(ctx, nselectTaxCategoriesSQL, uuid4byte, "active")
	taxCategoryTmp := taxstruct.TaxCategory{}
	err = row.StructScan(&taxCategoryTmp)
	if err != nil {
		ts.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	taxCategory, err := ts.getTaxCategoryStruct(ctx, in, &taxCategoryTmp)
	if err != nil {
		ts.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	taxCategoryResponse := taxproto.GetTaxCategoryResponse{}
	taxCategoryResponse.TaxCategory = taxCategory
	return &taxCategoryResponse, nil
}

// getTaxCategoryStruct - Get TaxCategory
func (ts *TaxService) getTaxCategoryStruct(ctx context.Context, in *commonproto.GetRequest, taxCategoryTmp *taxstruct.TaxCategory) (*taxproto.TaxCategory, error) {
	crUpdTime := new(commonproto.CrUpdTime)
	crUpdTime.CreatedAt = common.TimeToTimestamp(taxCategoryTmp.CrUpdTime.CreatedAt)
	crUpdTime.UpdatedAt = common.TimeToTimestamp(taxCategoryTmp.CrUpdTime.UpdatedAt)

	uuid4Str, err := common.UUIDBytesToStr(taxCategoryTmp.TaxCategoryD.Uuid4)
	if err != nil {
		ts.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	taxCategoryTmp.TaxCategoryD.IdS = uuid4Str

	taxCategory := taxproto.TaxCategory{TaxCategoryD: taxCategoryTmp.TaxCategoryD, CrUpdUser: taxCategoryTmp.CrUpdUser, CrUpdTime: crUpdTime}

	return &taxCategory, nil
}

// CreateTaxCategory - Create TaxCategory
func (ts *TaxService) CreateTaxCategory(ctx context.Context, in *taxproto.CreateTaxCategoryRequest) (*taxproto.CreateTaxCategoryResponse, error) {
	user, err := partyservice.GetUserWithNewContext(ctx, in.UserId, in.UserEmail, in.RequestId, ts.UserServiceClient)
	if err != nil {
		ts.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	form := taxproto.GetTaxSchemeRequest{}
	getRequest := commonproto.GetRequest{}
	getRequest.Id = in.TaxSchemeId
	getRequest.UserEmail = in.UserEmail
	getRequest.RequestId = in.RequestId
	form.GetRequest = &getRequest

	taxSchemeResponse, err := ts.GetTaxScheme(ctx, &form)
	if err != nil {
		ts.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	taxScheme := taxSchemeResponse.TaxScheme

	ttime := common.GetTimeDetails()
	tn := common.TimeToTimestamp(ttime)

	taxCategoryD := taxproto.TaxCategoryD{}
	taxCategoryD.Uuid4, err = common.GetUUIDBytes()
	if err != nil {
		ts.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	taxCategoryD.TcId = in.TcId
	taxCategoryD.TaxCategoryName = in.TaxCategoryName
	taxCategoryD.Percent = in.Percent
	taxCategoryD.BaseUnitMeasure = in.BaseUnitMeasure
	taxCategoryD.PerUnitAmount = in.PerUnitAmount
	taxCategoryD.TaxExemptionReasonCode = in.TaxExemptionReasonCode
	taxCategoryD.TaxExemptionReason = in.TaxExemptionReason
	taxCategoryD.TierRange = in.TierRange
	taxCategoryD.TierRatePercent = in.TierRatePercent
	taxCategoryD.TaxSchemeId = taxScheme.TaxSchemeD.Id

	crUpdUser := commonproto.CrUpdUser{}
	crUpdUser.StatusCode = "active"
	crUpdUser.CreatedByUserId = user.Id
	crUpdUser.UpdatedByUserId = user.Id

	crUpdTime := commonproto.CrUpdTime{}
	crUpdTime.CreatedAt = tn
	crUpdTime.UpdatedAt = tn

	taxCategory := taxproto.TaxCategory{TaxCategoryD: &taxCategoryD, CrUpdUser: &crUpdUser, CrUpdTime: &crUpdTime}

	err = ts.insertTaxCategory(ctx, insertTaxCategorySQL, &taxCategory, in.GetUserEmail(), in.GetRequestId())
	if err != nil {
		ts.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	taxCategoryResponse := taxproto.CreateTaxCategoryResponse{}
	taxCategoryResponse.TaxCategory = &taxCategory
	return &taxCategoryResponse, nil
}

// insertTaxCategory - Insert tax category details into database
func (ts *TaxService) insertTaxCategory(ctx context.Context, insertTaxCategorySQL string, taxCategory *taxproto.TaxCategory, userEmail string, requestID string) error {
	taxCategoryTmp, err := ts.crTaxCategoryStruct(ctx, taxCategory, userEmail, requestID)
	if err != nil {
		ts.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
		return err
	}

	err = ts.DBService.InsUpd(ctx, userEmail, requestID, func(tx *sqlx.Tx) error {
		res, err := tx.NamedExecContext(ctx, insertTaxCategorySQL, taxCategoryTmp)
		if err != nil {
			ts.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}

		uID, err := res.LastInsertId()
		if err != nil {
			ts.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}
		taxCategory.TaxCategoryD.Id = uint32(uID)
		uuid4Str, err := common.UUIDBytesToStr(taxCategory.TaxCategoryD.Uuid4)
		if err != nil {
			ts.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}
		taxCategory.TaxCategoryD.IdS = uuid4Str
		return nil
	})

	if err != nil {
		ts.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
		return err
	}
	return nil
}

// crTaxCategoryStruct - process TaxCategory details
func (ts *TaxService) crTaxCategoryStruct(ctx context.Context, taxCategory *taxproto.TaxCategory, userEmail string, requestID string) (*taxstruct.TaxCategory, error) {
	crUpdTime := new(commonstruct.CrUpdTime)
	crUpdTime.CreatedAt = common.TimestampToTime(taxCategory.CrUpdTime.CreatedAt)
	crUpdTime.UpdatedAt = common.TimestampToTime(taxCategory.CrUpdTime.UpdatedAt)

	taxCategoryTmp := taxstruct.TaxCategory{TaxCategoryD: taxCategory.TaxCategoryD, CrUpdUser: taxCategory.CrUpdUser, CrUpdTime: crUpdTime}

	return &taxCategoryTmp, nil
}

// UpdateTaxCategory - Update TaxCategory
func (ts *TaxService) UpdateTaxCategory(ctx context.Context, in *taxproto.UpdateTaxCategoryRequest) (*taxproto.UpdateTaxCategoryResponse, error) {
	db := ts.DBService.DB
	tn := common.GetTimeDetails()

	uuid4byte, err := common.UUIDStrToBytes(in.Id)
	if err != nil {
		ts.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	stmt, err := db.PreparexContext(ctx, updateTaxCategorySQL)
	if err != nil {
		ts.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	err = ts.DBService.InsUpd(ctx, in.GetUserEmail(), in.GetRequestId(), func(tx *sqlx.Tx) error {
		_, err = tx.StmtxContext(ctx, stmt).ExecContext(ctx,
			in.TaxCategoryName,
			in.Percent,
			in.BaseUnitMeasure,
			in.PerUnitAmount,
			in.TaxExemptionReasonCode,
			in.TaxExemptionReason,
			tn,
			uuid4byte)
		if err != nil {
			ts.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
			err1 := stmt.Close()
			if err1 != nil {
				ts.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err1))
				return err1
			}
			return err
		}
		return nil
	})

	if err != nil {
		ts.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	err = stmt.Close()
	if err != nil {
		ts.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	return &taxproto.UpdateTaxCategoryResponse{}, nil
}

// CreateTaxSchemeJurisdiction - Create TaxSchemeJurisdiction
func (ts *TaxService) CreateTaxSchemeJurisdiction(ctx context.Context, in *taxproto.CreateTaxSchemeJurisdictionRequest) (*taxproto.CreateTaxSchemeJurisdictionResponse, error) {
	user, err := partyservice.GetUserWithNewContext(ctx, in.UserId, in.UserEmail, in.RequestId, ts.UserServiceClient)
	if err != nil {
		ts.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	form := taxproto.GetTaxSchemeRequest{}
	getRequest := commonproto.GetRequest{}
	getRequest.Id = in.TaxSchemeId
	getRequest.UserEmail = in.UserEmail
	getRequest.RequestId = in.RequestId
	form.GetRequest = &getRequest

	taxSchemeResponse, err := ts.GetTaxScheme(ctx, &form)
	if err != nil {
		ts.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	taxScheme := taxSchemeResponse.TaxScheme

	ttime := common.GetTimeDetails()
	tn := common.TimeToTimestamp(ttime)

	taxSchemeJurisdictionD := taxproto.TaxSchemeJurisdictionD{}
	taxSchemeJurisdictionD.Uuid4, err = common.GetUUIDBytes()
	if err != nil {
		ts.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	address := commonproto.Address{}
	address.AddrListAgencyId = in.AddrListAgencyId
	address.AddrListId = in.AddrListAgencyId
	address.AddrListVersionId = in.AddrListVersionId
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

	taxSchemeInfo := commonproto.TaxSchemeInfo{}
	taxSchemeInfo.TaxSchemeId = taxScheme.TaxSchemeD.Id
	taxSchemeInfo.TaxSchemeName = in.TaxSchemeName

	crUpdUser := commonproto.CrUpdUser{}
	crUpdUser.StatusCode = "active"
	crUpdUser.CreatedByUserId = user.Id
	crUpdUser.UpdatedByUserId = user.Id

	crUpdTime := commonproto.CrUpdTime{}
	crUpdTime.CreatedAt = tn
	crUpdTime.UpdatedAt = tn

	taxSchemeJurisdiction := taxproto.TaxSchemeJurisdiction{TaxSchemeJurisdictionD: &taxSchemeJurisdictionD, TaxSchemeInfo: &taxSchemeInfo, CrUpdUser: &crUpdUser, CrUpdTime: &crUpdTime}

	err = ts.insertTaxSchemeJurisdiction(ctx, insertTaxSchemeJurisdictionSQL, &taxSchemeJurisdiction, &address, in.GetUserEmail(), in.GetRequestId())
	if err != nil {
		ts.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))

		return nil, err
	}

	taxSchemeJurisdictionResponse := taxproto.CreateTaxSchemeJurisdictionResponse{}
	taxSchemeJurisdictionResponse.TaxSchemeJurisdiction = &taxSchemeJurisdiction
	return &taxSchemeJurisdictionResponse, nil
}

// insertTaxSchemeJurisdiction - Insert tax scheme details into database
func (ts *TaxService) insertTaxSchemeJurisdiction(ctx context.Context, insertTaxSchemeJurisdictionSQL string, taxSchemeJurisdiction *taxproto.TaxSchemeJurisdiction, address *commonproto.Address, userEmail string, requestID string) error {
	taxSchemeJurisdictionTmp, err := ts.crTaxSchemeJurisdictionStruct(ctx, taxSchemeJurisdiction, userEmail, requestID)
	if err != nil {
		ts.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
		return err
	}

	err = ts.DBService.InsUpd(ctx, userEmail, requestID, func(tx *sqlx.Tx) error {
		addr, err := common.InsertAddress(ctx, tx, address, userEmail, requestID)
		if err != nil {
			ts.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}
		taxSchemeJurisdiction.TaxSchemeJurisdictionD.AddressId = addr.Id
		taxSchemeJurisdictionTmp.TaxSchemeJurisdictionD.AddressId = addr.Id
		res, err := tx.NamedExecContext(ctx, insertTaxSchemeJurisdictionSQL, taxSchemeJurisdictionTmp)
		if err != nil {
			ts.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}

		uID, err := res.LastInsertId()
		if err != nil {
			ts.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}
		taxSchemeJurisdiction.TaxSchemeJurisdictionD.Id = uint32(uID)
		uuid4Str, err := common.UUIDBytesToStr(taxSchemeJurisdiction.TaxSchemeJurisdictionD.Uuid4)
		if err != nil {
			ts.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}
		taxSchemeJurisdiction.TaxSchemeJurisdictionD.IdS = uuid4Str
		return nil
	})

	if err != nil {
		ts.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
		return err
	}
	return nil
}

// crTaxSchemeJurisdictionStruct - process TaxSchemeJurisdiction details
func (ts *TaxService) crTaxSchemeJurisdictionStruct(ctx context.Context, taxSchemeJurisdiction *taxproto.TaxSchemeJurisdiction, userEmail string, requestID string) (*taxstruct.TaxSchemeJurisdiction, error) {
	crUpdTime := new(commonstruct.CrUpdTime)
	crUpdTime.CreatedAt = common.TimestampToTime(taxSchemeJurisdiction.CrUpdTime.CreatedAt)
	crUpdTime.UpdatedAt = common.TimestampToTime(taxSchemeJurisdiction.CrUpdTime.UpdatedAt)

	taxSchemeJurisdictionTmp := taxstruct.TaxSchemeJurisdiction{TaxSchemeJurisdictionD: taxSchemeJurisdiction.TaxSchemeJurisdictionD, TaxSchemeInfo: taxSchemeJurisdiction.TaxSchemeInfo, CrUpdUser: taxSchemeJurisdiction.CrUpdUser, CrUpdTime: crUpdTime}

	return &taxSchemeJurisdictionTmp, nil
}

// UpdateTaxSchemeJurisdiction - Update TaxSchemeJurisdiction
func (ts *TaxService) UpdateTaxSchemeJurisdiction(ctx context.Context, in *taxproto.UpdateTaxSchemeJurisdictionRequest) (*taxproto.UpdateTaxSchemeJurisdictionResponse, error) {
	db := ts.DBService.DB
	tn := common.GetTimeDetails()

	uuid4byte, err := common.UUIDStrToBytes(in.Id)
	if err != nil {
		ts.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	stmt, err := db.PreparexContext(ctx, updateTaxSchemeJurisdictionSQL)
	if err != nil {
		ts.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	err = ts.DBService.InsUpd(ctx, in.GetUserEmail(), in.GetRequestId(), func(tx *sqlx.Tx) error {
		_, err = tx.StmtxContext(ctx, stmt).ExecContext(ctx,
			in.TaxSchemeName,
			in.TaxSchemeId,
			tn,
			uuid4byte)
		if err != nil {
			ts.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
			err1 := stmt.Close()
			if err1 != nil {
				ts.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err1))
				return err1
			}
			return err
		}
		return nil
	})

	if err != nil {
		ts.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	err = stmt.Close()
	if err != nil {
		ts.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	return &taxproto.UpdateTaxSchemeJurisdictionResponse{}, nil
}

// GetTaxTotal - Get TaxTotal
func (ts *TaxService) GetTaxTotal(ctx context.Context, inReq *taxproto.GetTaxTotalRequest) (*taxproto.GetTaxTotalResponse, error) {
	in := inReq.GetRequest
	uuid4byte, err := common.UUIDStrToBytes(in.Id)
	if err != nil {
		ts.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	nselectTaxTotalsSQL := selectTaxTotalsSQL + ` where uuid4 = ? and status_code = ?;`
	row := ts.DBService.DB.QueryRowxContext(ctx, nselectTaxTotalsSQL, uuid4byte, "active")
	taxTotalTmp := taxstruct.TaxTotal{}
	err = row.StructScan(&taxTotalTmp)
	if err != nil {
		ts.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	taxTotal, err := ts.getTaxTotalStruct(ctx, in, &taxTotalTmp)
	if err != nil {
		ts.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	taxTotalResponse := taxproto.GetTaxTotalResponse{}
	taxTotalResponse.TaxTotal = taxTotal
	return &taxTotalResponse, nil
}

// getTaxTotalStruct - Get TaxTotal
func (ts *TaxService) getTaxTotalStruct(ctx context.Context, in *commonproto.GetRequest, taxTotalTmp *taxstruct.TaxTotal) (*taxproto.TaxTotal, error) {
	crUpdTime := new(commonproto.CrUpdTime)
	crUpdTime.CreatedAt = common.TimeToTimestamp(taxTotalTmp.CrUpdTime.CreatedAt)
	crUpdTime.UpdatedAt = common.TimeToTimestamp(taxTotalTmp.CrUpdTime.UpdatedAt)

	uuid4Str, err := common.UUIDBytesToStr(taxTotalTmp.TaxTotalD.Uuid4)
	if err != nil {
		ts.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	taxTotalTmp.TaxTotalD.IdS = uuid4Str

	taxTotal := taxproto.TaxTotal{TaxTotalD: taxTotalTmp.TaxTotalD, CrUpdUser: taxTotalTmp.CrUpdUser, CrUpdTime: crUpdTime}

	return &taxTotal, nil
}

// CreateTaxTotal - Create TaxTotal
func (ts *TaxService) CreateTaxTotal(ctx context.Context, in *taxproto.CreateTaxTotalRequest) (*taxproto.CreateTaxTotalResponse, error) {
	user, err := partyservice.GetUserWithNewContext(ctx, in.UserId, in.UserEmail, in.RequestId, ts.UserServiceClient)
	if err != nil {
		ts.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	form := taxproto.GetTaxCategoryRequest{}
	getRequest := commonproto.GetRequest{}
	getRequest.Id = in.TaxCategoryId
	getRequest.UserEmail = in.UserEmail
	getRequest.RequestId = in.RequestId
	form.GetRequest = &getRequest

	taxCategoryResponse, err := ts.GetTaxCategory(ctx, &form)
	if err != nil {
		ts.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	taxCategory := taxCategoryResponse.TaxCategory

	ttime := common.GetTimeDetails()
	tn := common.TimeToTimestamp(ttime)

	taxTotalD := taxproto.TaxTotalD{}
	taxTotalD.Uuid4, err = common.GetUUIDBytes()
	if err != nil {
		ts.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	taxTotalD.TaxAmount = in.TaxAmount
	taxTotalD.RoundingAmount = in.RoundingAmount
	taxTotalD.TaxEvidenceIndicator = in.TaxEvidenceIndicator
	taxTotalD.TaxIncludedIndicator = in.TaxEvidenceIndicator
	taxTotalD.MasterFlag = in.MasterFlag
	taxTotalD.MasterId = in.MasterId
	taxTotalD.TaxCategoryId = taxCategory.TaxCategoryD.Id

	crUpdUser := commonproto.CrUpdUser{}
	crUpdUser.StatusCode = "active"
	crUpdUser.CreatedByUserId = user.Id
	crUpdUser.UpdatedByUserId = user.Id

	crUpdTime := commonproto.CrUpdTime{}
	crUpdTime.CreatedAt = tn
	crUpdTime.UpdatedAt = tn

	taxTotal := taxproto.TaxTotal{TaxTotalD: &taxTotalD, CrUpdUser: &crUpdUser, CrUpdTime: &crUpdTime}

	err = ts.insertTaxTotal(ctx, insertTaxTotalSQL, &taxTotal, in.GetUserEmail(), in.GetRequestId())
	if err != nil {
		ts.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	taxTotalResponse := taxproto.CreateTaxTotalResponse{}
	taxTotalResponse.TaxTotal = &taxTotal
	return &taxTotalResponse, nil
}

// insertTaxTotal - Insert tax scheme details into database
func (ts *TaxService) insertTaxTotal(ctx context.Context, insertTaxTotalSQL string, taxTotal *taxproto.TaxTotal, userEmail string, requestID string) error {
	taxTotalTmp, err := ts.crTaxTotalStruct(ctx, taxTotal, userEmail, requestID)
	if err != nil {
		ts.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
		return err
	}
	err = ts.DBService.InsUpd(ctx, userEmail, requestID, func(tx *sqlx.Tx) error {
		res, err := tx.NamedExecContext(ctx, insertTaxTotalSQL, taxTotalTmp)
		if err != nil {
			ts.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}

		uID, err := res.LastInsertId()
		if err != nil {
			ts.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}
		taxTotal.TaxTotalD.Id = uint32(uID)
		uuid4Str, err := common.UUIDBytesToStr(taxTotal.TaxTotalD.Uuid4)
		if err != nil {
			ts.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}
		taxTotal.TaxTotalD.IdS = uuid4Str
		return nil
	})

	if err != nil {
		ts.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
		return err
	}
	return nil
}

// crTaxTotalStruct - process TaxTotal details
func (ts *TaxService) crTaxTotalStruct(ctx context.Context, taxTotal *taxproto.TaxTotal, userEmail string, requestID string) (*taxstruct.TaxTotal, error) {
	crUpdTime := new(commonstruct.CrUpdTime)
	crUpdTime.CreatedAt = common.TimestampToTime(taxTotal.CrUpdTime.CreatedAt)
	crUpdTime.UpdatedAt = common.TimestampToTime(taxTotal.CrUpdTime.UpdatedAt)

	taxTotalTmp := taxstruct.TaxTotal{TaxTotalD: taxTotal.TaxTotalD, CrUpdUser: taxTotal.CrUpdUser, CrUpdTime: crUpdTime}

	return &taxTotalTmp, nil
}

// UpdateTaxTotal - Update TaxTotal
func (ts *TaxService) UpdateTaxTotal(ctx context.Context, in *taxproto.UpdateTaxTotalRequest) (*taxproto.UpdateTaxTotalResponse, error) {
	db := ts.DBService.DB
	tn := common.GetTimeDetails()

	uuid4byte, err := common.UUIDStrToBytes(in.Id)
	if err != nil {
		ts.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	stmt, err := db.PreparexContext(ctx, updateTaxTotalSQL)
	if err != nil {
		ts.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	err = ts.DBService.InsUpd(ctx, in.GetUserEmail(), in.GetRequestId(), func(tx *sqlx.Tx) error {
		_, err = tx.StmtxContext(ctx, stmt).ExecContext(ctx,
			in.TaxAmount,
			in.RoundingAmount,
			tn,
			uuid4byte)
		if err != nil {
			ts.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
			err1 := stmt.Close()
			if err1 != nil {
				ts.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err1))
				return err1
			}
			return err
		}
		return nil
	})

	if err != nil {
		ts.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	err = stmt.Close()
	if err != nil {
		ts.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	return &taxproto.UpdateTaxTotalResponse{}, nil
}

// CreateTaxSubTotal - Create TaxSubTotal
func (ts *TaxService) CreateTaxSubTotal(ctx context.Context, in *taxproto.CreateTaxSubTotalRequest) (*taxproto.CreateTaxSubTotalResponse, error) {
	user, err := partyservice.GetUserWithNewContext(ctx, in.UserId, in.UserEmail, in.RequestId, ts.UserServiceClient)
	if err != nil {
		ts.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	form := taxproto.GetTaxTotalRequest{}
	getRequest := commonproto.GetRequest{}
	getRequest.Id = in.TaxTotalId
	getRequest.UserEmail = in.UserEmail
	getRequest.RequestId = in.RequestId
	form.GetRequest = &getRequest

	taxTotalResponse, err := ts.GetTaxTotal(ctx, &form)
	if err != nil {
		ts.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	taxTotal := taxTotalResponse.TaxTotal

	ttime := common.GetTimeDetails()
	tn := common.TimeToTimestamp(ttime)

	taxSubTotalD := taxproto.TaxSubTotalD{}
	taxSubTotalD.Uuid4, err = common.GetUUIDBytes()
	if err != nil {
		ts.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	taxSubTotalD.TaxableAmount = in.TaxableAmount
	taxSubTotalD.TaxAmount = in.TaxAmount
	taxSubTotalD.CalculationSequenceNumeric = in.CalculationSequenceNumeric
	taxSubTotalD.TransactionCurrencyTaxAmount = in.TransactionCurrencyTaxAmount
	taxSubTotalD.Percent = in.Percent
	taxSubTotalD.BaseUnitMeasure = in.BaseUnitMeasure
	taxSubTotalD.PerUnitAmount = in.PerUnitAmount
	taxSubTotalD.TierRange = in.TierRange
	taxSubTotalD.TierRatePercent = in.TierRatePercent
	taxSubTotalD.TaxCategoryId = in.TaxCategoryId
	taxSubTotalD.TaxTotalId = taxTotal.TaxTotalD.Id

	crUpdUser := commonproto.CrUpdUser{}
	crUpdUser.StatusCode = "active"
	crUpdUser.CreatedByUserId = user.Id
	crUpdUser.UpdatedByUserId = user.Id

	crUpdTime := commonproto.CrUpdTime{}
	crUpdTime.CreatedAt = tn
	crUpdTime.UpdatedAt = tn

	taxSubTotal := taxproto.TaxSubTotal{TaxSubTotalD: &taxSubTotalD, CrUpdUser: &crUpdUser, CrUpdTime: &crUpdTime}

	err = ts.insertTaxSubTotal(ctx, insertTaxSubTotalSQL, &taxSubTotal, in.GetUserEmail(), in.GetRequestId())
	if err != nil {
		ts.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	taxSubTotalResponse := taxproto.CreateTaxSubTotalResponse{}
	taxSubTotalResponse.TaxSubTotal = &taxSubTotal
	return &taxSubTotalResponse, nil
}

// insertTaxSubTotal - Insert tax scheme details into database
func (ts *TaxService) insertTaxSubTotal(ctx context.Context, insertTaxSubTotalSQL string, taxSubTotal *taxproto.TaxSubTotal, userEmail string, requestID string) error {
	taxSubTotalTmp, err := ts.crTaxSubTotalStruct(ctx, taxSubTotal, userEmail, requestID)
	if err != nil {
		ts.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
		return err
	}

	err = ts.DBService.InsUpd(ctx, userEmail, requestID, func(tx *sqlx.Tx) error {
		res, err := tx.NamedExecContext(ctx, insertTaxSubTotalSQL, taxSubTotalTmp)
		if err != nil {
			ts.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}

		uID, err := res.LastInsertId()
		if err != nil {
			ts.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}
		taxSubTotal.TaxSubTotalD.Id = uint32(uID)

		uuid4Str, err := common.UUIDBytesToStr(taxSubTotal.TaxSubTotalD.Uuid4)
		if err != nil {
			ts.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}
		taxSubTotal.TaxSubTotalD.IdS = uuid4Str

		return nil
	})

	if err != nil {
		ts.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
		return err
	}
	return nil
}

// crTaxSubTotalStruct - process TaxSubTotal details
func (ts *TaxService) crTaxSubTotalStruct(ctx context.Context, taxSubTotal *taxproto.TaxSubTotal, userEmail string, requestID string) (*taxstruct.TaxSubTotal, error) {
	crUpdTime := new(commonstruct.CrUpdTime)
	crUpdTime.CreatedAt = common.TimestampToTime(taxSubTotal.CrUpdTime.CreatedAt)
	crUpdTime.UpdatedAt = common.TimestampToTime(taxSubTotal.CrUpdTime.UpdatedAt)

	taxSubTotalTmp := taxstruct.TaxSubTotal{TaxSubTotalD: taxSubTotal.TaxSubTotalD, CrUpdUser: taxSubTotal.CrUpdUser, CrUpdTime: crUpdTime}

	return &taxSubTotalTmp, nil
}

// UpdateTaxSubTotal - Update TaxSubTotal
func (ts *TaxService) UpdateTaxSubTotal(ctx context.Context, in *taxproto.UpdateTaxSubTotalRequest) (*taxproto.UpdateTaxSubTotalResponse, error) {
	db := ts.DBService.DB
	tn := common.GetTimeDetails()

	uuid4byte, err := common.UUIDStrToBytes(in.Id)
	if err != nil {
		ts.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	stmt, err := db.PreparexContext(ctx, updateTaxSubTotalSQL)
	if err != nil {
		ts.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	err = ts.DBService.InsUpd(ctx, in.GetUserEmail(), in.GetRequestId(), func(tx *sqlx.Tx) error {
		_, err = tx.StmtxContext(ctx, stmt).ExecContext(ctx,
			in.TaxableAmount,
			in.TaxAmount,
			in.CalculationSequenceNumeric,
			in.TransactionCurrencyTaxAmount,
			in.Percent,
			in.BaseUnitMeasure,
			in.PerUnitAmount,
			tn,
			uuid4byte)
		if err != nil {
			ts.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
			err1 := stmt.Close()
			if err1 != nil {
				ts.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err1))
				return err1
			}
			return err
		}
		return nil
	})

	if err != nil {
		ts.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	err = stmt.Close()
	if err != nil {
		ts.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	return &taxproto.UpdateTaxSubTotalResponse{}, nil
}
