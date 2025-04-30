package logisticsservices

import (
	"context"
	"time"

	"github.com/cloudfresco/sc-ubl/internal/common"
	commonproto "github.com/cloudfresco/sc-ubl/internal/protogen/common/v1"
	logisticsproto "github.com/cloudfresco/sc-ubl/internal/protogen/logistics/v1"
	partyproto "github.com/cloudfresco/sc-ubl/internal/protogen/party/v1"
	partyservice "github.com/cloudfresco/sc-ubl/internal/services/partyservices"
	commonstruct "github.com/cloudfresco/sc-ubl/internal/servicestructs/common/v1"
	logisticsstruct "github.com/cloudfresco/sc-ubl/internal/servicestructs/logistics/v1"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// DespatchService - For accessing Despatch services
type DespatchService struct {
	log               *zap.Logger
	DBService         *common.DBService
	RedisService      *common.RedisService
	UserServiceClient partyproto.UserServiceClient
	logisticsproto.UnimplementedDespatchServiceServer
}

// NewDespatchService - Create Despatch service
func NewDespatchService(log *zap.Logger, dbOpt *common.DBService, redisOpt *common.RedisService, userServiceClient partyproto.UserServiceClient) *DespatchService {
	return &DespatchService{
		log:               log,
		DBService:         dbOpt,
		RedisService:      redisOpt,
		UserServiceClient: userServiceClient,
	}
}

const insertDespatchHeaderSQL = `insert into despatch_headers
	  ( 
    uuid4,
    desph_id,
    document_status_code,
    despatch_advice_type_code,
    note,
    line_count_numeric,
    order_id,
    despatch_supplier_party_id,
    delivery_customer_party_id,
    buyer_customer_party_id,
    seller_supplier_party_id,
    originator_customer_party_id,
    shipment_id,
    issue_date,
    status_code,
    created_by_user_id,
    updated_by_user_id,
    created_at,
    updated_at)
  values (:uuid4,
:desph_id,
:document_status_code,
:despatch_advice_type_code,
:note,
:line_count_numeric,
:order_id,
:despatch_supplier_party_id,
:delivery_customer_party_id,
:buyer_customer_party_id,
:seller_supplier_party_id,
:originator_customer_party_id,
:shipment_id,
:issue_date,
:status_code,
:created_by_user_id,
:updated_by_user_id,
:created_at,
:updated_at);`

const selectDespatchHeadersSQL = `select 
  id,
  uuid4,
  desph_id,
  document_status_code,
  despatch_advice_type_code,
  note,
  line_count_numeric,
  order_id,
  despatch_supplier_party_id,
  delivery_customer_party_id,
  buyer_customer_party_id,
  seller_supplier_party_id,
  originator_customer_party_id,
  shipment_id,
  issue_date,
  status_code,
  created_by_user_id,
  updated_by_user_id,
  created_at,
  updated_at from despatch_headers`

// updateDespatchSQL - update DespatchSQL query
const updateDespatchHeaderSQL = `update despatch_headers set 
  document_status_code= ?,
  despatch_advice_type_code= ?,
  note= ?,
  updated_at = ? where uuid4 = ?;`

// CreateDespatchHeader - Create Despatch Header
func (ds *DespatchService) CreateDespatchHeader(ctx context.Context, in *logisticsproto.CreateDespatchHeaderRequest) (*logisticsproto.CreateDespatchHeaderResponse, error) {
	user, err := partyservice.GetUserWithNewContext(ctx, in.UserId, in.UserEmail, in.RequestId, ds.UserServiceClient)
	if err != nil {
		ds.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	ttime := common.GetTimeDetails()
	tn := common.TimeToTimestamp(ttime)

	issueDate, err := time.Parse(common.Layout, in.IssueDate)
	if err != nil {
		ds.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	despatchHeaderD := logisticsproto.DespatchHeaderD{}
	despatchHeaderD.Uuid4, err = common.GetUUIDBytes()
	if err != nil {
		ds.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	despatchHeaderD.DesphId = in.DesphId
	despatchHeaderD.DocumentStatusCode = in.DocumentStatusCode
	despatchHeaderD.DespatchAdviceTypeCode = in.DespatchAdviceTypeCode
	despatchHeaderD.Note = in.Note
	despatchHeaderD.LineCountNumeric = in.LineCountNumeric
	despatchHeaderD.OrderId = in.OrderId
	despatchHeaderD.DespatchSupplierPartyId = in.DespatchSupplierPartyId
	despatchHeaderD.DeliveryCustomerPartyId = in.DeliveryCustomerPartyId
	despatchHeaderD.BuyerCustomerPartyId = in.BuyerCustomerPartyId
	despatchHeaderD.SellerSupplierPartyId = in.SellerSupplierPartyId
	despatchHeaderD.OriginatorCustomerPartyId = in.OriginatorCustomerPartyId
	despatchHeaderD.ShipmentId = in.OriginatorCustomerPartyId

	despatchHeaderT := logisticsproto.DespatchHeaderT{}
	despatchHeaderT.IssueDate = common.TimeToTimestamp(issueDate.UTC().Truncate(time.Second))

	crUpdUser := commonproto.CrUpdUser{}
	crUpdUser.StatusCode = "active"
	crUpdUser.CreatedByUserId = user.Id
	crUpdUser.UpdatedByUserId = user.Id

	crUpdTime := commonproto.CrUpdTime{}
	crUpdTime.CreatedAt = tn
	crUpdTime.UpdatedAt = tn

	despatchLines := []*logisticsproto.DespatchLine{}
	// we will do for loop on lines which is comes from client form
	for _, line := range in.DespatchLines {
		line.UserId = in.UserId
		line.UserEmail = in.UserEmail
		line.RequestId = in.RequestId
		// we wl call CreateDespatchLine function which wl populate form values to Despatchline struct
		despatchLine, err := ds.ProcessDespatchLineRequest(ctx, line)
		if err != nil {
			ds.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
			return nil, err
		}
		despatchLines = append(despatchLines, despatchLine)
	}

	despatchHeader := logisticsproto.DespatchHeader{DespatchHeaderD: &despatchHeaderD, DespatchHeaderT: &despatchHeaderT, CrUpdUser: &crUpdUser, CrUpdTime: &crUpdTime}

	err = ds.insertDespatchHeader(ctx, insertDespatchHeaderSQL, &despatchHeader, insertDespatchLineSQL, despatchLines, in.GetUserEmail(), in.GetRequestId())
	if err != nil {
		ds.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	despatchHeaderResponse := logisticsproto.CreateDespatchHeaderResponse{}
	despatchHeaderResponse.DespatchHeader = &despatchHeader
	return &despatchHeaderResponse, nil
}

// insertDespatch - Insert Desptach
func (ds *DespatchService) insertDespatchHeader(ctx context.Context, insertDespatchHeaderSQL string, despatchHeader *logisticsproto.DespatchHeader, insertDespatchLineSQL string, despatchLines []*logisticsproto.DespatchLine, userEmail string, requestID string) error {
	despatchHeaderTmp, err := ds.crDespatchHeaderStruct(ctx, despatchHeader, userEmail, requestID)
	if err != nil {
		ds.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
		return err
	}

	err = ds.DBService.InsUpd(ctx, userEmail, requestID, func(tx *sqlx.Tx) error {
		res, err := tx.NamedExecContext(ctx, insertDespatchHeaderSQL, despatchHeaderTmp)
		if err != nil {
			ds.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}

		uID, err := res.LastInsertId()
		if err != nil {
			ds.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}
		despatchHeader.DespatchHeaderD.Id = uint32(uID)
		uuid4Str, err := common.UUIDBytesToStr(despatchHeader.DespatchHeaderD.Uuid4)
		if err != nil {
			ds.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}
		despatchHeader.DespatchHeaderD.IdS = uuid4Str

		for _, despatchLine := range despatchLines {
			despatchLine.DespatchLineD.DespatchHeaderId = despatchHeader.DespatchHeaderD.Id
			despatchLineTmp, err := ds.crDespatchLineStruct(ctx, despatchLine, userEmail, requestID)
			if err != nil {
				ds.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
				return err
			}
			_, err = tx.NamedExecContext(ctx, insertDespatchLineSQL, despatchLineTmp)
			if err != nil {
				ds.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
				return err
			}
		}
		return nil
	})

	if err != nil {
		ds.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
		return err
	}
	return nil
}

// crDespatchHeaderStruct - process DespatchHeader details
func (ds *DespatchService) crDespatchHeaderStruct(ctx context.Context, despatchHeader *logisticsproto.DespatchHeader, userEmail string, requestID string) (*logisticsstruct.DespatchHeader, error) {
	despatchHeaderT := new(logisticsstruct.DespatchHeaderT)
	despatchHeaderT.IssueDate = common.TimestampToTime(despatchHeader.DespatchHeaderT.IssueDate)

	crUpdTime := new(commonstruct.CrUpdTime)
	crUpdTime.CreatedAt = common.TimestampToTime(despatchHeader.CrUpdTime.CreatedAt)
	crUpdTime.UpdatedAt = common.TimestampToTime(despatchHeader.CrUpdTime.UpdatedAt)

	despatchHeaderTmp := logisticsstruct.DespatchHeader{DespatchHeaderD: despatchHeader.DespatchHeaderD, DespatchHeaderT: despatchHeaderT, CrUpdUser: despatchHeader.CrUpdUser, CrUpdTime: crUpdTime}

	return &despatchHeaderTmp, nil
}

// GetDespatchHeaders - Get DespatchHeaders
func (ds *DespatchService) GetDespatchHeaders(ctx context.Context, in *logisticsproto.GetDespatchHeadersRequest) (*logisticsproto.GetDespatchHeadersResponse, error) {
	limit := in.GetLimit()
	nextCursor := in.GetNextCursor()
	if limit == "" {
		limit = ds.DBService.LimitSQLRows
	}
	query := "status_code = ?"
	if nextCursor == "" {
		query = query + " order by id desc " + " limit " + limit + ";"
	} else {
		nextCursor = common.DecodeCursor(nextCursor)
		query = query + " " + "and" + " " + "id <= " + nextCursor + " order by id desc " + " limit " + limit + ";"
	}

	despatchHeaders := []*logisticsproto.DespatchHeader{}

	nselectDespatchHeadersSQL := selectDespatchHeadersSQL + ` where ` + query

	rows, err := ds.DBService.DB.QueryxContext(ctx, nselectDespatchHeadersSQL, "active")
	if err != nil {
		ds.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	for rows.Next() {

		despatchHeaderTmp := logisticsstruct.DespatchHeader{}
		err = rows.StructScan(&despatchHeaderTmp)
		if err != nil {
			ds.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
			return nil, err
		}

		getRequest := commonproto.GetRequest{}
		getRequest.UserEmail = in.UserEmail
		getRequest.RequestId = in.RequestId
		despatchHeader, err := ds.getDespatchHeaderStruct(ctx, &getRequest, despatchHeaderTmp)
		if err != nil {
			ds.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
			return nil, err
		}
		despatchHeaders = append(despatchHeaders, despatchHeader)

	}

	despatchHeadersResponse := logisticsproto.GetDespatchHeadersResponse{}
	if len(despatchHeaders) != 0 {
		next := despatchHeaders[len(despatchHeaders)-1].DespatchHeaderD.Id
		next--
		nextc := common.EncodeCursor(next)
		despatchHeadersResponse = logisticsproto.GetDespatchHeadersResponse{DespatchHeaders: despatchHeaders, NextCursor: nextc}
	} else {
		despatchHeadersResponse = logisticsproto.GetDespatchHeadersResponse{DespatchHeaders: despatchHeaders, NextCursor: "0"}
	}
	return &despatchHeadersResponse, nil
}

// GetDespatchHeader - Get DespatchHeader
func (ds *DespatchService) GetDespatchHeader(ctx context.Context, inReq *logisticsproto.GetDespatchHeaderRequest) (*logisticsproto.GetDespatchHeaderResponse, error) {
	in := inReq.GetRequest
	uuid4byte, err := common.UUIDStrToBytes(in.Id)
	if err != nil {
		ds.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	nselectDespatchHeadersSQL := selectDespatchHeadersSQL + ` where uuid4 = ? and status_code = ?;`
	row := ds.DBService.DB.QueryRowxContext(ctx, nselectDespatchHeadersSQL, uuid4byte, "active")
	despatchHeaderTmp := logisticsstruct.DespatchHeader{}
	err = row.StructScan(&despatchHeaderTmp)
	if err != nil {
		ds.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	despatchHeader, err := ds.getDespatchHeaderStruct(ctx, in, despatchHeaderTmp)
	if err != nil {
		ds.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	despatchHeaderResponse := logisticsproto.GetDespatchHeaderResponse{}
	despatchHeaderResponse.DespatchHeader = despatchHeader
	return &despatchHeaderResponse, nil
}

// GetDespatchHeaderByPk - Get DespatchHeader By Primary key(Id)
func (ds *DespatchService) GetDespatchHeaderByPk(ctx context.Context, inReq *logisticsproto.GetDespatchHeaderByPkRequest) (*logisticsproto.GetDespatchHeaderByPkResponse, error) {
	in := inReq.GetByIdRequest
	nselectDespatchHeadersSQL := selectDespatchHeadersSQL + ` where id = ? and status_code = ?;`
	row := ds.DBService.DB.QueryRowxContext(ctx, nselectDespatchHeadersSQL, in.Id, "active")
	despatchHeaderTmp := logisticsstruct.DespatchHeader{}
	err := row.StructScan(&despatchHeaderTmp)
	if err != nil {
		ds.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	getRequest := commonproto.GetRequest{}
	getRequest.UserEmail = in.UserEmail
	getRequest.RequestId = in.RequestId
	despatchHeader, err := ds.getDespatchHeaderStruct(ctx, &getRequest, despatchHeaderTmp)
	if err != nil {
		ds.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	despatchHeaderResponse := logisticsproto.GetDespatchHeaderByPkResponse{}
	despatchHeaderResponse.DespatchHeader = despatchHeader
	return &despatchHeaderResponse, nil
}

// getDespatchHeaderStruct - Get despatch header
func (ds *DespatchService) getDespatchHeaderStruct(ctx context.Context, in *commonproto.GetRequest, despatchHeaderTmp logisticsstruct.DespatchHeader) (*logisticsproto.DespatchHeader, error) {
	uuid4Str, err := common.UUIDBytesToStr(despatchHeaderTmp.DespatchHeaderD.Uuid4)
	if err != nil {
		ds.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	despatchHeaderTmp.DespatchHeaderD.IdS = uuid4Str

	despatchHeaderT := new(logisticsproto.DespatchHeaderT)
	despatchHeaderT.IssueDate = common.TimeToTimestamp(despatchHeaderTmp.DespatchHeaderT.IssueDate)

	crUpdTime := new(commonproto.CrUpdTime)
	crUpdTime.CreatedAt = common.TimeToTimestamp(despatchHeaderTmp.CrUpdTime.CreatedAt)
	crUpdTime.UpdatedAt = common.TimeToTimestamp(despatchHeaderTmp.CrUpdTime.UpdatedAt)

	despatchHeader := logisticsproto.DespatchHeader{DespatchHeaderD: despatchHeaderTmp.DespatchHeaderD, DespatchHeaderT: despatchHeaderT, CrUpdUser: despatchHeaderTmp.CrUpdUser, CrUpdTime: crUpdTime}

	return &despatchHeader, nil
}

// UpdateDespatchHeader - Update DespatchHeader
func (ds *DespatchService) UpdateDespatchHeader(ctx context.Context, in *logisticsproto.UpdateDespatchHeaderRequest) (*logisticsproto.UpdateDespatchHeaderResponse, error) {
	db := ds.DBService.DB
	tn := common.GetTimeDetails()

	uuid4byte, err := common.UUIDStrToBytes(in.Id)
	if err != nil {
		ds.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	stmt, err := db.PreparexContext(ctx, updateDespatchHeaderSQL)
	if err != nil {
		ds.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	err = ds.DBService.InsUpd(ctx, in.GetUserEmail(), in.GetRequestId(), func(tx *sqlx.Tx) error {
		_, err = tx.StmtxContext(ctx, stmt).ExecContext(ctx,
			in.DocumentStatusCode,
			in.DespatchAdviceTypeCode,
			in.Note,
			tn,
			uuid4byte)
		if err != nil {
			ds.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
			err1 := stmt.Close()
			if err1 != nil {
				ds.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err1))
				return err1
			}
			return err
		}
		return nil
	})

	if err != nil {
		ds.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	err = stmt.Close()
	if err != nil {
		ds.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	return &logisticsproto.UpdateDespatchHeaderResponse{}, nil
}
