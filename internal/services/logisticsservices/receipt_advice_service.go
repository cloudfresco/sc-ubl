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

// ReceiptAdviceHeaderService - For accessing receipt advice header services
type ReceiptAdviceHeaderService struct {
	log               *zap.Logger
	DBService         *common.DBService
	RedisService      *common.RedisService
	UserServiceClient partyproto.UserServiceClient
	logisticsproto.UnimplementedReceiptAdviceHeaderServiceServer
}

// NewReceiptAdviceHeaderService - Create Receipt Advice Header service
func NewReceiptAdviceHeaderService(log *zap.Logger, dbOpt *common.DBService, redisOpt *common.RedisService, userServiceClient partyproto.UserServiceClient) *ReceiptAdviceHeaderService {
	return &ReceiptAdviceHeaderService{
		log:               log,
		DBService:         dbOpt,
		RedisService:      redisOpt,
		UserServiceClient: userServiceClient,
	}
}

const insertReceiptAdviceHeaderSQL = `insert into receipt_advice_headers
	  ( 
    uuid4,
    rcpth_id,
    receipt_advice_type_code,
    note,
    line_count_numeric,
    order_id,
    despatch_id,
    delivery_customer_party_id,
    despatch_supplier_party_id,
    buyer_customer_party_id,
    seller_supplier_party_id,
    shipment_id,
    issue_date,
    status_code,
    created_by_user_id,
    updated_by_user_id,
    created_at,
    updated_at)
  values (:uuid4,
:rcpth_id,
:receipt_advice_type_code,
:note,
:line_count_numeric,
:order_id,
:despatch_id,
:delivery_customer_party_id,
:despatch_supplier_party_id,
:buyer_customer_party_id,
:seller_supplier_party_id,
:shipment_id,
:issue_date,
:status_code,
:created_by_user_id,
:updated_by_user_id,
:created_at,
:updated_at);`

const selectReceiptAdviceHeadersSQL = `select 
  id,
  uuid4,
  rcpth_id,
  receipt_advice_type_code,
  note,
  line_count_numeric,
  order_id,
  despatch_id,
  delivery_customer_party_id,
  despatch_supplier_party_id,
  buyer_customer_party_id,
  seller_supplier_party_id,
  shipment_id,
  issue_date,
  status_code,
  created_by_user_id,
  updated_by_user_id,
  created_at,
  updated_at from receipt_advice_headers`

// updateReceiptAdviceHeaderSQL - update ReceiptAdviceHeaderSQL query
const updateReceiptAdviceHeaderSQL = `update receipt_advice_headers set 
  receipt_advice_type_code= ?,
  note= ?,
  line_count_numeric= ?,
  updated_at = ? where uuid4 = ?;`

// CreateReceiptAdviceHeader - Create ReceiptAdviceHeader
func (rs *ReceiptAdviceHeaderService) CreateReceiptAdviceHeader(ctx context.Context, in *logisticsproto.CreateReceiptAdviceHeaderRequest) (*logisticsproto.CreateReceiptAdviceHeaderResponse, error) {
	user, err := partyservice.GetUserWithNewContext(ctx, in.UserId, in.UserEmail, in.RequestId, rs.UserServiceClient)
	if err != nil {
		rs.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	ttime := common.GetTimeDetails()
	tn := common.TimeToTimestamp(ttime)

	issueDate, err := time.Parse(common.Layout, in.IssueDate)
	if err != nil {
		rs.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	receiptAdviceHeaderD := logisticsproto.ReceiptAdviceHeaderD{}
	receiptAdviceHeaderD.Uuid4, err = common.GetUUIDBytes()
	if err != nil {
		rs.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	receiptAdviceHeaderD.RcpthId = in.RcpthId
	receiptAdviceHeaderD.ReceiptAdviceTypeCode = in.ReceiptAdviceTypeCode
	receiptAdviceHeaderD.Note = in.Note
	receiptAdviceHeaderD.LineCountNumeric = in.LineCountNumeric
	receiptAdviceHeaderD.OrderId = in.OrderId
	receiptAdviceHeaderD.DespatchId = in.DespatchId
	receiptAdviceHeaderD.DeliveryCustomerPartyId = in.DeliveryCustomerPartyId
	receiptAdviceHeaderD.DespatchSupplierPartyId = in.DespatchSupplierPartyId
	receiptAdviceHeaderD.BuyerCustomerPartyId = in.BuyerCustomerPartyId
	receiptAdviceHeaderD.SellerSupplierPartyId = in.SellerSupplierPartyId
	receiptAdviceHeaderD.ShipmentId = in.ShipmentId

	receiptAdviceHeaderT := logisticsproto.ReceiptAdviceHeaderT{}
	receiptAdviceHeaderT.IssueDate = common.TimeToTimestamp(issueDate.UTC().Truncate(time.Second))

	crUpdUser := commonproto.CrUpdUser{}
	crUpdUser.StatusCode = "active"
	crUpdUser.CreatedByUserId = user.Id
	crUpdUser.UpdatedByUserId = user.Id

	crUpdTime := commonproto.CrUpdTime{}
	crUpdTime.CreatedAt = tn
	crUpdTime.UpdatedAt = tn

	receiptAdviceLines := []*logisticsproto.ReceiptAdviceLine{}
	// we will do for loop on lines which is comes from client form
	for _, line := range in.ReceiptAdviceLines {
		line.UserId = in.UserId
		line.UserEmail = in.UserEmail
		line.RequestId = in.RequestId
		// we wl call CreateReceiptAdviceLine function which wl populate form values to receiptAdviceline struct
		receiptAdviceLine, err := rs.ProcessReceiptAdviceLineRequest(ctx, line)
		if err != nil {
			rs.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
			return nil, err
		}
		receiptAdviceLines = append(receiptAdviceLines, receiptAdviceLine)
	}

	receiptAdviceHeader := logisticsproto.ReceiptAdviceHeader{ReceiptAdviceHeaderD: &receiptAdviceHeaderD, ReceiptAdviceHeaderT: &receiptAdviceHeaderT, CrUpdUser: &crUpdUser, CrUpdTime: &crUpdTime}

	err = rs.insertReceiptAdviceHeader(ctx, insertReceiptAdviceHeaderSQL, &receiptAdviceHeader, insertReceiptAdviceLineSQL, receiptAdviceLines, in.GetUserEmail(), in.GetRequestId())
	if err != nil {
		rs.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	receiptAdviceHeaderResponse := logisticsproto.CreateReceiptAdviceHeaderResponse{}
	receiptAdviceHeaderResponse.ReceiptAdviceHeader = &receiptAdviceHeader
	return &receiptAdviceHeaderResponse, nil
}

// insertReceiptAdviceHeader - Insert ReceiptAdviceHeader
func (rs *ReceiptAdviceHeaderService) insertReceiptAdviceHeader(ctx context.Context, insertReceiptAdviceHeaderSQL string, receiptAdviceHeader *logisticsproto.ReceiptAdviceHeader, insertReceiptAdviceLineSQL string, receiptAdviceLines []*logisticsproto.ReceiptAdviceLine, userEmail string, requestID string) error {
	receiptAdviceHeaderTmp, err := rs.crReceiptAdviceHeaderStruct(ctx, receiptAdviceHeader, userEmail, requestID)
	if err != nil {
		rs.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
		return err
	}
	err = rs.DBService.InsUpd(ctx, userEmail, requestID, func(tx *sqlx.Tx) error {
		res, err := tx.NamedExecContext(ctx, insertReceiptAdviceHeaderSQL, receiptAdviceHeaderTmp)
		if err != nil {
			rs.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}

		uID, err := res.LastInsertId()
		if err != nil {
			rs.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}
		receiptAdviceHeader.ReceiptAdviceHeaderD.Id = uint32(uID)
		uuid4Str, err := common.UUIDBytesToStr(receiptAdviceHeader.ReceiptAdviceHeaderD.Uuid4)
		if err != nil {
			rs.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}
		receiptAdviceHeader.ReceiptAdviceHeaderD.IdS = uuid4Str
		for _, receiptAdviceLine := range receiptAdviceLines {
			receiptAdviceLine.ReceiptAdviceLineD.ReceiptAdviceHeaderId = receiptAdviceHeader.ReceiptAdviceHeaderD.Id
			receiptAdviceLineTmp, err := rs.crReceiptAdviceLineStruct(ctx, receiptAdviceLine, userEmail, requestID)
			if err != nil {
				rs.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
				return err
			}
			_, err = tx.NamedExecContext(ctx, insertReceiptAdviceLineSQL, receiptAdviceLineTmp)
			if err != nil {
				rs.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
				return err
			}
		}

		return nil
	})

	if err != nil {
		rs.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
		return err
	}
	return nil
}

// crReceiptAdviceHeaderStruct - process ReceiptAdviceHeader details
func (rs *ReceiptAdviceHeaderService) crReceiptAdviceHeaderStruct(ctx context.Context, receiptAdviceHeader *logisticsproto.ReceiptAdviceHeader, userEmail string, requestID string) (*logisticsstruct.ReceiptAdviceHeader, error) {
	receiptAdviceHeaderT := new(logisticsstruct.ReceiptAdviceHeaderT)
	receiptAdviceHeaderT.IssueDate = common.TimestampToTime(receiptAdviceHeader.ReceiptAdviceHeaderT.IssueDate)

	crUpdTime := new(commonstruct.CrUpdTime)
	crUpdTime.CreatedAt = common.TimestampToTime(receiptAdviceHeader.CrUpdTime.CreatedAt)
	crUpdTime.UpdatedAt = common.TimestampToTime(receiptAdviceHeader.CrUpdTime.UpdatedAt)

	receiptAdviceHeaderTmp := logisticsstruct.ReceiptAdviceHeader{ReceiptAdviceHeaderD: receiptAdviceHeader.ReceiptAdviceHeaderD, ReceiptAdviceHeaderT: receiptAdviceHeaderT, CrUpdUser: receiptAdviceHeader.CrUpdUser, CrUpdTime: crUpdTime}

	return &receiptAdviceHeaderTmp, nil
}

// GetReceiptAdviceHeaders - Get ReceiptAdviceHeaders
func (rs *ReceiptAdviceHeaderService) GetReceiptAdviceHeaders(ctx context.Context, in *logisticsproto.GetReceiptAdviceHeadersRequest) (*logisticsproto.GetReceiptAdviceHeadersResponse, error) {
	limit := in.GetLimit()
	nextCursor := in.GetNextCursor()
	if limit == "" {
		limit = rs.DBService.LimitSQLRows
	}
	query := "status_code = ?"
	if nextCursor == "" {
		query = query + " order by id desc " + " limit " + limit + ";"
	} else {
		nextCursor = common.DecodeCursor(nextCursor)
		query = query + " " + "and" + " " + "id <= " + nextCursor + " order by id desc " + " limit " + limit + ";"
	}

	receiptAdviceHeaders := []*logisticsproto.ReceiptAdviceHeader{}

	nselectReceiptAdviceHeadersSQL := selectReceiptAdviceHeadersSQL + ` where ` + query

	rows, err := rs.DBService.DB.QueryxContext(ctx, nselectReceiptAdviceHeadersSQL, "active")
	if err != nil {
		rs.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	for rows.Next() {

		receiptAdviceHeaderTmp := logisticsstruct.ReceiptAdviceHeader{}
		err = rows.StructScan(&receiptAdviceHeaderTmp)
		if err != nil {
			rs.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
			return nil, err
		}
		getRequest := commonproto.GetRequest{}
		getRequest.UserEmail = in.UserEmail
		getRequest.RequestId = in.RequestId
		receiptAdviceHeader, err := rs.getReceiptAdviceHeaderStruct(ctx, &getRequest, receiptAdviceHeaderTmp)
		if err != nil {
			rs.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
			return nil, err
		}

		receiptAdviceHeaders = append(receiptAdviceHeaders, receiptAdviceHeader)

	}

	receiptAdviceHeadersResponse := logisticsproto.GetReceiptAdviceHeadersResponse{}
	if len(receiptAdviceHeaders) != 0 {
		next := receiptAdviceHeaders[len(receiptAdviceHeaders)-1].ReceiptAdviceHeaderD.Id
		next--
		nextc := common.EncodeCursor(next)
		receiptAdviceHeadersResponse = logisticsproto.GetReceiptAdviceHeadersResponse{ReceiptAdviceHeaders: receiptAdviceHeaders, NextCursor: nextc}
	} else {
		receiptAdviceHeadersResponse = logisticsproto.GetReceiptAdviceHeadersResponse{ReceiptAdviceHeaders: receiptAdviceHeaders, NextCursor: "0"}
	}
	return &receiptAdviceHeadersResponse, nil
}

// GetReceiptAdviceHeader - Get ReceiptAdviceHeader
func (rs *ReceiptAdviceHeaderService) GetReceiptAdviceHeader(ctx context.Context, inReq *logisticsproto.GetReceiptAdviceHeaderRequest) (*logisticsproto.GetReceiptAdviceHeaderResponse, error) {
	in := inReq.GetRequest
	uuid4byte, err := common.UUIDStrToBytes(in.Id)
	if err != nil {
		rs.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	nselectReceiptAdviceHeadersSQL := selectReceiptAdviceHeadersSQL + ` where uuid4 = ? and status_code = ?;`
	row := rs.DBService.DB.QueryRowxContext(ctx, nselectReceiptAdviceHeadersSQL, uuid4byte, "active")
	receiptAdviceHeaderTmp := logisticsstruct.ReceiptAdviceHeader{}
	err = row.StructScan(&receiptAdviceHeaderTmp)
	if err != nil {
		rs.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	receiptAdviceHeader, err := rs.getReceiptAdviceHeaderStruct(ctx, in, receiptAdviceHeaderTmp)
	if err != nil {
		rs.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	receiptAdviceHeaderResponse := logisticsproto.GetReceiptAdviceHeaderResponse{}
	receiptAdviceHeaderResponse.ReceiptAdviceHeader = receiptAdviceHeader
	return &receiptAdviceHeaderResponse, nil
}

// GetReceiptAdviceHeaderByPk - Get ReceiptAdviceHeader By Primary key(Id)
func (rs *ReceiptAdviceHeaderService) GetReceiptAdviceHeaderByPk(ctx context.Context, inReq *logisticsproto.GetReceiptAdviceHeaderByPkRequest) (*logisticsproto.GetReceiptAdviceHeaderByPkResponse, error) {
	in := inReq.GetByIdRequest
	nselectReceiptAdviceHeadersSQL := selectReceiptAdviceHeadersSQL + ` where id = ? and status_code = ?;`
	row := rs.DBService.DB.QueryRowxContext(ctx, nselectReceiptAdviceHeadersSQL, in.Id, "active")
	receiptAdviceHeaderTmp := logisticsstruct.ReceiptAdviceHeader{}
	err := row.StructScan(&receiptAdviceHeaderTmp)
	if err != nil {
		rs.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	getRequest := commonproto.GetRequest{}
	getRequest.UserEmail = in.UserEmail
	getRequest.RequestId = in.RequestId
	receiptAdviceHeader, err := rs.getReceiptAdviceHeaderStruct(ctx, &getRequest, receiptAdviceHeaderTmp)
	if err != nil {
		rs.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	receiptAdviceHeaderResponse := logisticsproto.GetReceiptAdviceHeaderByPkResponse{}
	receiptAdviceHeaderResponse.ReceiptAdviceHeader = receiptAdviceHeader
	return &receiptAdviceHeaderResponse, nil
}

// getReceiptAdviceHeaderStruct - Get ReceiptAdviceHeader
func (rs *ReceiptAdviceHeaderService) getReceiptAdviceHeaderStruct(ctx context.Context, in *commonproto.GetRequest, receiptAdviceHeaderTmp logisticsstruct.ReceiptAdviceHeader) (*logisticsproto.ReceiptAdviceHeader, error) {
	uuid4Str, err := common.UUIDBytesToStr(receiptAdviceHeaderTmp.ReceiptAdviceHeaderD.Uuid4)
	if err != nil {
		rs.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	receiptAdviceHeaderTmp.ReceiptAdviceHeaderD.IdS = uuid4Str

	receiptAdviceHeaderT := new(logisticsproto.ReceiptAdviceHeaderT)
	receiptAdviceHeaderT.IssueDate = common.TimeToTimestamp(receiptAdviceHeaderTmp.ReceiptAdviceHeaderT.IssueDate)

	crUpdTime := new(commonproto.CrUpdTime)
	crUpdTime.CreatedAt = common.TimeToTimestamp(receiptAdviceHeaderTmp.CrUpdTime.CreatedAt)
	crUpdTime.UpdatedAt = common.TimeToTimestamp(receiptAdviceHeaderTmp.CrUpdTime.UpdatedAt)

	receiptAdviceHeader := logisticsproto.ReceiptAdviceHeader{ReceiptAdviceHeaderD: receiptAdviceHeaderTmp.ReceiptAdviceHeaderD, ReceiptAdviceHeaderT: receiptAdviceHeaderT, CrUpdUser: receiptAdviceHeaderTmp.CrUpdUser, CrUpdTime: crUpdTime}

	return &receiptAdviceHeader, nil
}

// UpdateReceiptAdviceHeader - Update ReceiptAdviceHeader
func (rs *ReceiptAdviceHeaderService) UpdateReceiptAdviceHeader(ctx context.Context, in *logisticsproto.UpdateReceiptAdviceHeaderRequest) (*logisticsproto.UpdateReceiptAdviceHeaderResponse, error) {
	db := rs.DBService.DB
	tn := common.GetTimeDetails()

	uuid4byte, err := common.UUIDStrToBytes(in.Id)
	if err != nil {
		rs.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	stmt, err := db.PreparexContext(ctx, updateReceiptAdviceHeaderSQL)
	if err != nil {
		rs.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	err = rs.DBService.InsUpd(ctx, in.GetUserEmail(), in.GetRequestId(), func(tx *sqlx.Tx) error {
		_, err = tx.StmtxContext(ctx, stmt).ExecContext(ctx,
			in.ReceiptAdviceTypeCode,
			in.Note,
			in.LineCountNumeric,
			tn,
			uuid4byte)
		if err != nil {
			rs.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
			err1 := stmt.Close()
			if err1 != nil {
				rs.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err1))
				return err1
			}
			return err
		}
		return nil
	})

	if err != nil {
		rs.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	err = stmt.Close()
	if err != nil {
		rs.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	return &logisticsproto.UpdateReceiptAdviceHeaderResponse{}, nil
}
