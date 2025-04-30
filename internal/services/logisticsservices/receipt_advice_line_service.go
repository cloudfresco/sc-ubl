package logisticsservices

import (
	"context"
	"time"

	"github.com/cloudfresco/sc-ubl/internal/common"
	commonproto "github.com/cloudfresco/sc-ubl/internal/protogen/common/v1"
	logisticsproto "github.com/cloudfresco/sc-ubl/internal/protogen/logistics/v1"
	partyservice "github.com/cloudfresco/sc-ubl/internal/services/partyservices"
	commonstruct "github.com/cloudfresco/sc-ubl/internal/servicestructs/common/v1"
	logisticsstruct "github.com/cloudfresco/sc-ubl/internal/servicestructs/logistics/v1"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

const insertReceiptAdviceLineSQL = `insert into receipt_advice_lines
	    (uuid4,
      rcptl_id,
      note,
      received_quantity,
      short_quantity,
      shortage_action_code,
      rejected_quantity,
      reject_reason_code,
      reject_reason,
      reject_action_code,
      quantity_discrepancy_code,
      oversupply_quantity,
      timing_complaint_code,
      timing_complaint,
      order_line_id,
      despatch_line_id,
      item_id,
      shipment_id,
      receipt_advice_header_id,
      received_date,
      status_code,
      created_by_user_id,
      updated_by_user_id,
      created_at,
      updated_at)
  values (:uuid4,
          :rcptl_id,
          :note,
          :received_quantity,
          :short_quantity,
          :shortage_action_code,
          :rejected_quantity,
          :reject_reason_code,
          :reject_reason,
          :reject_action_code,
          :quantity_discrepancy_code,
          :oversupply_quantity,
          :timing_complaint_code,
          :timing_complaint,
          :order_line_id,
          :despatch_line_id,
          :item_id,
          :shipment_id,
          :receipt_advice_header_id,
          :received_date,
          :status_code,
          :created_by_user_id,
          :updated_by_user_id,
          :created_at,
          :updated_at);`

const selectReceiptAdviceLinesSQL = `select 
  id,
  uuid4,
  rcptl_id,
  note,
  received_quantity,
  short_quantity,
  shortage_action_code,
  rejected_quantity,
  reject_reason_code,
  reject_reason,
  reject_action_code,
  quantity_discrepancy_code,
  oversupply_quantity,
  timing_complaint_code,
  timing_complaint,
  order_line_id,
  despatch_line_id,
  item_id,
  shipment_id,
  receipt_advice_header_id,
  received_date,
  status_code,
  created_by_user_id,
  updated_by_user_id,
  created_at,
  updated_at from receipt_advice_lines`

// CreateReceiptAdviceLine - Create Despatch Line
func (rs *ReceiptAdviceHeaderService) CreateReceiptAdviceLine(ctx context.Context, in *logisticsproto.CreateReceiptAdviceLineRequest) (*logisticsproto.CreateReceiptAdviceLineResponse, error) {
	receiptAdviceLine, err := rs.ProcessReceiptAdviceLineRequest(ctx, in)
	if err != nil {
		rs.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	err = rs.insertReceiptAdviceLine(ctx, insertReceiptAdviceLineSQL, receiptAdviceLine, in.GetUserEmail(), in.GetRequestId())
	if err != nil {
		rs.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	receiptAdviceLineResponse := logisticsproto.CreateReceiptAdviceLineResponse{}
	receiptAdviceLineResponse.ReceiptAdviceLine = receiptAdviceLine
	return &receiptAdviceLineResponse, nil
}

// ProcessReceiptAdviceLineRequest - ProcessReceiptAdviceLineRequest
func (rs *ReceiptAdviceHeaderService) ProcessReceiptAdviceLineRequest(ctx context.Context, in *logisticsproto.CreateReceiptAdviceLineRequest) (*logisticsproto.ReceiptAdviceLine, error) {
	user, err := partyservice.GetUserWithNewContext(ctx, in.UserId, in.UserEmail, in.RequestId, rs.UserServiceClient)
	if err != nil {
		rs.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	ttime := common.GetTimeDetails()
	tn := common.TimeToTimestamp(ttime)

	receivedDate, err := time.Parse(common.Layout, in.ReceivedDate)
	if err != nil {
		rs.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	receiptAdviceLineD := logisticsproto.ReceiptAdviceLineD{}
	receiptAdviceLineD.Uuid4, err = common.GetUUIDBytes()
	if err != nil {
		rs.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	receiptAdviceLineD.RcptlId = in.RcptlId
	receiptAdviceLineD.Note = in.Note
	receiptAdviceLineD.ReceivedQuantity = in.ReceivedQuantity
	receiptAdviceLineD.ShortQuantity = in.ShortQuantity
	receiptAdviceLineD.ShortageActionCode = in.ShortageActionCode
	receiptAdviceLineD.RejectedQuantity = in.RejectedQuantity
	receiptAdviceLineD.RejectReasonCode = in.RejectReasonCode
	receiptAdviceLineD.RejectReason = in.RejectReason
	receiptAdviceLineD.RejectActionCode = in.RejectActionCode
	receiptAdviceLineD.QuantityDiscrepancyCode = in.QuantityDiscrepancyCode
	receiptAdviceLineD.OversupplyQuantity = in.OversupplyQuantity
	receiptAdviceLineD.TimingComplaintCode = in.TimingComplaintCode
	receiptAdviceLineD.TimingComplaint = in.TimingComplaint
	receiptAdviceLineD.OrderLineId = in.OrderLineId
	receiptAdviceLineD.DespatchLineId = in.DespatchLineId
	receiptAdviceLineD.ItemId = in.ItemId
	receiptAdviceLineD.ShipmentId = in.ShipmentId
	receiptAdviceLineD.ReceiptAdviceHeaderId = in.ReceiptAdviceHeaderId

	receiptAdviceLineT := logisticsproto.ReceiptAdviceLineT{}
	receiptAdviceLineT.ReceivedDate = common.TimeToTimestamp(receivedDate.UTC().Truncate(time.Second))

	crUpdUser := commonproto.CrUpdUser{}
	crUpdUser.StatusCode = "active"
	crUpdUser.CreatedByUserId = user.Id
	crUpdUser.UpdatedByUserId = user.Id

	crUpdTime := commonproto.CrUpdTime{}
	crUpdTime.CreatedAt = tn
	crUpdTime.UpdatedAt = tn

	receiptAdviceLine := logisticsproto.ReceiptAdviceLine{ReceiptAdviceLineD: &receiptAdviceLineD, ReceiptAdviceLineT: &receiptAdviceLineT, CrUpdUser: &crUpdUser, CrUpdTime: &crUpdTime}

	return &receiptAdviceLine, nil
}

// insertReceiptAdviceLine - Insert Despatch Line
func (rs *ReceiptAdviceHeaderService) insertReceiptAdviceLine(ctx context.Context, insertReceiptAdviceLineSQL string, receiptAdviceLine *logisticsproto.ReceiptAdviceLine, userEmail string, requestID string) error {
	receiptAdviceLineTmp, err := rs.crReceiptAdviceLineStruct(ctx, receiptAdviceLine, userEmail, requestID)
	if err != nil {
		rs.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
		return err
	}
	err = rs.DBService.InsUpd(ctx, userEmail, requestID, func(tx *sqlx.Tx) error {
		res, err := tx.NamedExecContext(ctx, insertReceiptAdviceLineSQL, receiptAdviceLineTmp)
		if err != nil {
			rs.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}

		uID, err := res.LastInsertId()
		if err != nil {
			rs.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}
		receiptAdviceLine.ReceiptAdviceLineD.Id = uint32(uID)
		uuid4Str, err := common.UUIDBytesToStr(receiptAdviceLine.ReceiptAdviceLineD.Uuid4)
		if err != nil {
			rs.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}
		receiptAdviceLine.ReceiptAdviceLineD.IdS = uuid4Str
		return nil
	})

	if err != nil {
		rs.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
		return err
	}
	return nil
}

// crReceiptAdviceLineStruct - process ReceiptAdviceLine details
func (rs *ReceiptAdviceHeaderService) crReceiptAdviceLineStruct(ctx context.Context, receiptAdviceLine *logisticsproto.ReceiptAdviceLine, userEmail string, requestID string) (*logisticsstruct.ReceiptAdviceLine, error) {
	receiptAdviceLineT := new(logisticsstruct.ReceiptAdviceLineT)
	receiptAdviceLineT.ReceivedDate = common.TimestampToTime(receiptAdviceLine.ReceiptAdviceLineT.ReceivedDate)

	crUpdTime := new(commonstruct.CrUpdTime)
	crUpdTime.CreatedAt = common.TimestampToTime(receiptAdviceLine.CrUpdTime.CreatedAt)
	crUpdTime.UpdatedAt = common.TimestampToTime(receiptAdviceLine.CrUpdTime.UpdatedAt)

	receiptAdviceLineTmp := logisticsstruct.ReceiptAdviceLine{ReceiptAdviceLineD: receiptAdviceLine.ReceiptAdviceLineD, ReceiptAdviceLineT: receiptAdviceLineT, CrUpdUser: receiptAdviceLine.CrUpdUser, CrUpdTime: crUpdTime}

	return &receiptAdviceLineTmp, nil
}

// GetReceiptAdviceLines - GetReceiptAdviceLines
func (rs *ReceiptAdviceHeaderService) GetReceiptAdviceLines(ctx context.Context, inReq *logisticsproto.GetReceiptAdviceLinesRequest) (*logisticsproto.GetReceiptAdviceLinesResponse, error) {
	in := inReq.GetRequest
	getRequest := commonproto.GetRequest{}
	getRequest.Id = in.Id
	getRequest.UserEmail = in.UserEmail
	getRequest.RequestId = in.RequestId
	form := logisticsproto.GetReceiptAdviceHeaderRequest{}
	form.GetRequest = &getRequest

	receiptAdviceHeaderResponse, err := rs.GetReceiptAdviceHeader(ctx, &form)
	if err != nil {
		rs.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	receiptAdviceHeader := receiptAdviceHeaderResponse.ReceiptAdviceHeader

	receiptAdviceLines := []*logisticsproto.ReceiptAdviceLine{}

	nselectReceiptAdviceLinesSQL := selectReceiptAdviceLinesSQL + ` where receipt_advice_header_id = ? and status_code = ?;`
	rows, err := rs.DBService.DB.QueryxContext(ctx, nselectReceiptAdviceLinesSQL, receiptAdviceHeader.ReceiptAdviceHeaderD.Id, "active")
	if err != nil {
		rs.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	for rows.Next() {

		receiptAdviceLineTmp := logisticsstruct.ReceiptAdviceLine{}
		err = rows.StructScan(&receiptAdviceLineTmp)
		if err != nil {
			rs.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
			return nil, err
		}

		receiptAdviceLine, err := rs.getReceiptAdviceLineStruct(ctx, in, receiptAdviceLineTmp)
		if err != nil {
			rs.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
			return nil, err
		}

		receiptAdviceLines = append(receiptAdviceLines, receiptAdviceLine)
	}
	receiptAdviceLinesResponse := logisticsproto.GetReceiptAdviceLinesResponse{}
	receiptAdviceLinesResponse.ReceiptAdviceLines = receiptAdviceLines
	return &receiptAdviceLinesResponse, nil
}

// getReceiptAdviceLineStruct - Get ReceiptAdviceHeader
func (rs *ReceiptAdviceHeaderService) getReceiptAdviceLineStruct(ctx context.Context, in *commonproto.GetRequest, receiptAdviceLineTmp logisticsstruct.ReceiptAdviceLine) (*logisticsproto.ReceiptAdviceLine, error) {
	uuid4Str, err := common.UUIDBytesToStr(receiptAdviceLineTmp.ReceiptAdviceLineD.Uuid4)
	if err != nil {
		rs.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	receiptAdviceLineTmp.ReceiptAdviceLineD.IdS = uuid4Str

	receiptAdviceLineT := new(logisticsproto.ReceiptAdviceLineT)
	receiptAdviceLineT.ReceivedDate = common.TimeToTimestamp(receiptAdviceLineTmp.ReceiptAdviceLineT.ReceivedDate)

	crUpdTime := new(commonproto.CrUpdTime)
	crUpdTime.CreatedAt = common.TimeToTimestamp(receiptAdviceLineTmp.CrUpdTime.CreatedAt)
	crUpdTime.UpdatedAt = common.TimeToTimestamp(receiptAdviceLineTmp.CrUpdTime.UpdatedAt)

	receiptAdviceLine := logisticsproto.ReceiptAdviceLine{ReceiptAdviceLineD: receiptAdviceLineTmp.ReceiptAdviceLineD, ReceiptAdviceLineT: receiptAdviceLineT, CrUpdUser: receiptAdviceLineTmp.CrUpdUser, CrUpdTime: crUpdTime}

	return &receiptAdviceLine, nil
}
