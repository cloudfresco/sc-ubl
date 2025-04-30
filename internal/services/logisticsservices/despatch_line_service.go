package logisticsservices

import (
	"context"

	"github.com/cloudfresco/sc-ubl/internal/common"
	commonproto "github.com/cloudfresco/sc-ubl/internal/protogen/common/v1"
	logisticsproto "github.com/cloudfresco/sc-ubl/internal/protogen/logistics/v1"
	partyservice "github.com/cloudfresco/sc-ubl/internal/services/partyservices"
	commonstruct "github.com/cloudfresco/sc-ubl/internal/servicestructs/common/v1"
	logisticsstruct "github.com/cloudfresco/sc-ubl/internal/servicestructs/logistics/v1"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

const insertDespatchLineSQL = `insert into despatch_lines
	  (uuid4,
    despl_id,
    note,
    line_status_code,
    delivered_quantity,
    backorder_quantity,
    backorder_reason,
    outstanding_quantity,
    outstanding_reason,
    oversupply_quantity,
    order_line_id,
    item_id,
    shipment_id,
    despatch_header_id,
    status_code,
    created_by_user_id,
    updated_by_user_id,
    created_at,
    updated_at)
  values (:uuid4,
          :despl_id,
          :note,
          :line_status_code,
          :delivered_quantity,
          :backorder_quantity,
          :backorder_reason,
          :outstanding_quantity,
          :outstanding_reason,
          :oversupply_quantity,
          :order_line_id,
          :item_id,
          :shipment_id,
          :despatch_header_id,
          :status_code,
          :created_by_user_id,
          :updated_by_user_id,
          :created_at,
          :updated_at);`

const selectDespatchLinesSQL = `select 
  id,
  uuid4,
  despl_id,
  note,
  line_status_code,
  delivered_quantity,
  backorder_quantity,
  backorder_reason,
  outstanding_quantity,
  outstanding_reason,
  oversupply_quantity,
  order_line_id,
  item_id,
  shipment_id,
  despatch_header_id,
  status_code,
  created_by_user_id,
  updated_by_user_id,
  created_at,
  updated_at from despatch_lines`

// CreateDespatchLine - Create Despatch Line
func (ds *DespatchService) CreateDespatchLine(ctx context.Context, in *logisticsproto.CreateDespatchLineRequest) (*logisticsproto.CreateDespatchLineResponse, error) {
	despatchLine, err := ds.ProcessDespatchLineRequest(ctx, in)
	if err != nil {
		ds.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	err = ds.insertDespatchLine(ctx, insertDespatchLineSQL, despatchLine, in.GetUserEmail(), in.GetRequestId())
	if err != nil {
		ds.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	despatchLineResponse := logisticsproto.CreateDespatchLineResponse{}
	despatchLineResponse.DespatchLine = despatchLine
	return &despatchLineResponse, nil
}

// ProcessDespatchLineRequest - ProcessDespatchLineRequest
func (ds *DespatchService) ProcessDespatchLineRequest(ctx context.Context, in *logisticsproto.CreateDespatchLineRequest) (*logisticsproto.DespatchLine, error) {
	user, err := partyservice.GetUserWithNewContext(ctx, in.UserId, in.UserEmail, in.RequestId, ds.UserServiceClient)
	if err != nil {
		ds.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	ttime := common.GetTimeDetails()
	tn := common.TimeToTimestamp(ttime)

	despatchLineD := logisticsproto.DespatchLineD{}
	despatchLineD.Uuid4, err = common.GetUUIDBytes()
	if err != nil {
		ds.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	despatchLineD.DesplId = in.DesplId
	despatchLineD.Note = in.Note
	despatchLineD.LineStatusCode = in.LineStatusCode
	despatchLineD.DeliveredQuantity = in.DeliveredQuantity
	despatchLineD.BackorderQuantity = in.BackorderQuantity
	despatchLineD.BackorderReason = in.BackorderReason
	despatchLineD.OutstandingQuantity = in.OutstandingQuantity
	despatchLineD.OutstandingReason = in.OutstandingReason
	despatchLineD.OversupplyQuantity = in.OversupplyQuantity
	despatchLineD.OrderLineId = in.OrderLineId
	despatchLineD.ItemId = in.ItemId
	despatchLineD.ShipmentId = in.ShipmentId
	despatchLineD.DespatchHeaderId = in.DespatchHeaderId

	crUpdUser := commonproto.CrUpdUser{}
	crUpdUser.StatusCode = "active"
	crUpdUser.CreatedByUserId = user.Id
	crUpdUser.UpdatedByUserId = user.Id

	crUpdTime := commonproto.CrUpdTime{}
	crUpdTime.CreatedAt = tn
	crUpdTime.UpdatedAt = tn

	despatchLine := logisticsproto.DespatchLine{DespatchLineD: &despatchLineD, CrUpdUser: &crUpdUser, CrUpdTime: &crUpdTime}

	return &despatchLine, nil
}

// insertDespatchLine - Insert Despatch Line
func (ds *DespatchService) insertDespatchLine(ctx context.Context, insertDespatchLineSQL string, despatchLine *logisticsproto.DespatchLine, userEmail string, requestID string) error {
	despatchLineTmp, err := ds.crDespatchLineStruct(ctx, despatchLine, userEmail, requestID)
	if err != nil {
		ds.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
		return err
	}
	err = ds.DBService.InsUpd(ctx, userEmail, requestID, func(tx *sqlx.Tx) error {
		res, err := tx.NamedExecContext(ctx, insertDespatchLineSQL, despatchLineTmp)
		if err != nil {
			ds.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}

		uID, err := res.LastInsertId()
		if err != nil {
			ds.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}
		despatchLine.DespatchLineD.Id = uint32(uID)
		uuid4Str, err := common.UUIDBytesToStr(despatchLine.DespatchLineD.Uuid4)
		if err != nil {
			ds.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}
		despatchLine.DespatchLineD.IdS = uuid4Str
		return nil
	})

	if err != nil {
		ds.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
		return err
	}
	return nil
}

// crDespatchLineStruct - process DespatchLine details
func (ds *DespatchService) crDespatchLineStruct(ctx context.Context, despatchLine *logisticsproto.DespatchLine, userEmail string, requestID string) (*logisticsstruct.DespatchLine, error) {
	crUpdTime := new(commonstruct.CrUpdTime)
	crUpdTime.CreatedAt = common.TimestampToTime(despatchLine.CrUpdTime.CreatedAt)
	crUpdTime.UpdatedAt = common.TimestampToTime(despatchLine.CrUpdTime.UpdatedAt)

	despatchLineTmp := logisticsstruct.DespatchLine{DespatchLineD: despatchLine.DespatchLineD, CrUpdUser: despatchLine.CrUpdUser, CrUpdTime: crUpdTime}

	return &despatchLineTmp, nil
}

// GetDespatchLines - GetDespatchLines
func (ds *DespatchService) GetDespatchLines(ctx context.Context, inReq *logisticsproto.GetDespatchLinesRequest) (*logisticsproto.GetDespatchLinesResponse, error) {
	in := inReq.GetRequest
	getRequest := commonproto.GetRequest{}
	getRequest.Id = in.Id
	getRequest.UserEmail = in.UserEmail
	getRequest.RequestId = in.RequestId
	form := logisticsproto.GetDespatchHeaderRequest{}
	form.GetRequest = &getRequest

	despatchHeaderResponse, err := ds.GetDespatchHeader(ctx, &form)
	if err != nil {
		ds.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	despatchHeader := despatchHeaderResponse.DespatchHeader

	despatchLines := []*logisticsproto.DespatchLine{}

	nselectDespatchLinesSQL := selectDespatchLinesSQL + ` where despatch_header_id = ? and status_code = ?;`
	rows, err := ds.DBService.DB.QueryxContext(ctx, nselectDespatchLinesSQL, despatchHeader.DespatchHeaderD.Id, "active")
	if err != nil {
		ds.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	for rows.Next() {

		despatchLineTmp := logisticsstruct.DespatchLine{}
		err = rows.StructScan(&despatchLineTmp)
		if err != nil {
			ds.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
			return nil, err
		}
		despatchLine, err := ds.getDespatchLineStruct(ctx, in, despatchLineTmp)
		if err != nil {
			ds.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
			return nil, err
		}
		despatchLines = append(despatchLines, despatchLine)
	}

	despatchLinesResponse := logisticsproto.GetDespatchLinesResponse{}
	despatchLinesResponse.DespatchLines = despatchLines
	return &despatchLinesResponse, nil
}

// getDespatchLineStruct - Get DespatchLine header
func (ds *DespatchService) getDespatchLineStruct(ctx context.Context, in *commonproto.GetRequest, despatchLineTmp logisticsstruct.DespatchLine) (*logisticsproto.DespatchLine, error) {
	uuid4Str, err := common.UUIDBytesToStr(despatchLineTmp.DespatchLineD.Uuid4)
	if err != nil {
		ds.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	despatchLineTmp.DespatchLineD.IdS = uuid4Str

	crUpdTime := new(commonproto.CrUpdTime)
	crUpdTime.CreatedAt = common.TimeToTimestamp(despatchLineTmp.CrUpdTime.CreatedAt)
	crUpdTime.UpdatedAt = common.TimeToTimestamp(despatchLineTmp.CrUpdTime.UpdatedAt)

	despatchLine := logisticsproto.DespatchLine{DespatchLineD: despatchLineTmp.DespatchLineD, CrUpdUser: despatchLineTmp.CrUpdUser, CrUpdTime: crUpdTime}

	return &despatchLine, nil
}
