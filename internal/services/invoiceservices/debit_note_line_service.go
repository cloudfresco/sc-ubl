package invoiceservices

import (
	"context"
	"time"

	"github.com/cloudfresco/sc-ubl/internal/common"
	commonproto "github.com/cloudfresco/sc-ubl/internal/protogen/common/v1"
	invoiceproto "github.com/cloudfresco/sc-ubl/internal/protogen/invoice/v1"
	partyproto "github.com/cloudfresco/sc-ubl/internal/protogen/party/v1"
	commonstruct "github.com/cloudfresco/sc-ubl/internal/servicestructs/common/v1"
	invoicestruct "github.com/cloudfresco/sc-ubl/internal/servicestructs/invoice/v1"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

const insertDebitNoteLineSQL = `insert into debit_note_lines
	  (uuid4,
dnl_id,
note,
debited_quantity,
line_extension_amount,
accounting_cost_code,
accounting_cost,
payment_purpose_code,
discrepancy_response,
despatch_line_id,
receipt_line_id,
billing_id,
item_id,
price_amount,
price_base_quantity,
price_change_reason,
price_type_code,
price_type,
orderable_unit_factor_rate,
price_list_id,
debit_note_header_id,
tax_point_date,
price_validity_period_start_date,
price_validity_period_end_date,
status_code,
created_by_user_id,
updated_by_user_id,
created_at,
updated_at)
  values (:uuid4,
:dnl_id,
:note,
:debited_quantity,
:line_extension_amount,
:accounting_cost_code,
:accounting_cost,
:payment_purpose_code,
:discrepancy_response,
:despatch_line_id,
:receipt_line_id,
:billing_id,
:item_id,
:price_amount,
:price_base_quantity,
:price_change_reason,
:price_type_code,
:price_type,
:orderable_unit_factor_rate,
:price_list_id,
:debit_note_header_id,
:tax_point_date,
:price_validity_period_start_date,
:price_validity_period_end_date,
:status_code,
:created_by_user_id,
:updated_by_user_id,
:created_at,
:updated_at);`

const selectDebitNoteLinesSQL = `select 
id,
uuid4,
dnl_id,
note,
debited_quantity,
line_extension_amount,
accounting_cost_code,
accounting_cost,
payment_purpose_code,
discrepancy_response,
despatch_line_id,
receipt_line_id,
billing_id,
item_id,
price_amount,
price_base_quantity,
price_change_reason,
price_type_code,
price_type,
orderable_unit_factor_rate,
price_list_id,
debit_note_header_id,
tax_point_date,
price_validity_period_start_date,
price_validity_period_end_date,
status_code,
created_by_user_id,
updated_by_user_id,
created_at,
updated_at from debit_note_lines`

// CreateDebitNoteLine - Create DebitNoteLine
func (ds *DebitNoteHeaderService) CreateDebitNoteLine(ctx context.Context, in *invoiceproto.CreateDebitNoteLineRequest) (*invoiceproto.CreateDebitNoteLineResponse, error) {
	debitNoteLine, err := ds.ProcessDebitNoteLineRequest(ctx, in)
	if err != nil {
		ds.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	err = ds.insertDebitNoteLine(ctx, insertDebitNoteLineSQL, debitNoteLine, in.GetUserEmail(), in.GetRequestId())
	if err != nil {
		ds.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	debitNoteLineResponse := invoiceproto.CreateDebitNoteLineResponse{}
	debitNoteLineResponse.DebitNoteLine = debitNoteLine
	return &debitNoteLineResponse, nil
}

func (ds *DebitNoteHeaderService) ProcessDebitNoteLineRequest(ctx context.Context, in *invoiceproto.CreateDebitNoteLineRequest) (*invoiceproto.DebitNoteLine, error) {
	getRequest := commonproto.GetRequest{}
	getRequest.Id = in.UserId
	getRequest.UserEmail = in.UserEmail
	getRequest.RequestId = in.RequestId
	ctxNew, err := common.CreateCtxJWT(ctx)
	if err != nil {
		ds.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	form := partyproto.GetUserRequest{}
	form.GetRequest = &getRequest
	userResponse, err := ds.UserServiceClient.GetUser(ctxNew, &form)
	if err != nil {
		ds.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	user := userResponse.User

	taxPointDate, err := time.Parse(common.Layout, in.TaxPointDate)
	if err != nil {
		ds.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	priceValidityPeriodStartDate, err := time.Parse(common.Layout, in.PriceValidityPeriodStartDate)
	if err != nil {
		ds.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	priceValidityPeriodEndDate, err := time.Parse(common.Layout, in.PriceValidityPeriodEndDate)
	if err != nil {
		ds.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	ttime := common.GetTimeDetails()
	tn := common.TimeToTimestamp(ttime)

	debitNoteLineD := invoiceproto.DebitNoteLineD{}
	debitNoteLineD.Uuid4, err = common.GetUUIDBytes()
	if err != nil {
		ds.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	debitNoteLineD.DnlId = in.DnlId
	debitNoteLineD.Note = in.Note
	debitNoteLineD.DebitedQuantity = in.DebitedQuantity
	debitNoteLineD.LineExtensionAmount = in.LineExtensionAmount
	debitNoteLineD.AccountingCostCode = in.AccountingCostCode
	debitNoteLineD.AccountingCost = in.AccountingCost
	debitNoteLineD.PaymentPurposeCode = in.PaymentPurposeCode
	debitNoteLineD.DiscrepancyResponse = in.DiscrepancyResponse
	debitNoteLineD.DespatchLineId = in.DespatchLineId
	debitNoteLineD.ReceiptLineId = in.ReceiptLineId
	debitNoteLineD.BillingId = in.BillingId
	debitNoteLineD.ItemId = in.ItemId
	debitNoteLineD.PriceAmount = in.PriceAmount
	debitNoteLineD.PriceBaseQuantity = in.PriceBaseQuantity
	debitNoteLineD.PriceChangeReason = in.PriceChangeReason
	debitNoteLineD.PriceTypeCode = in.PriceTypeCode
	debitNoteLineD.PriceType = in.PriceType
	debitNoteLineD.OrderableUnitFactorRate = in.OrderableUnitFactorRate
	debitNoteLineD.PriceListId = in.PriceListId
	debitNoteLineD.DebitNoteHeaderId = in.DebitNoteHeaderId

	debitNoteLineT := invoiceproto.DebitNoteLineT{}
	debitNoteLineT.TaxPointDate = common.TimeToTimestamp(taxPointDate.UTC().Truncate(time.Second))
	debitNoteLineT.PriceValidityPeriodStartDate = common.TimeToTimestamp(priceValidityPeriodStartDate.UTC().Truncate(time.Second))
	debitNoteLineT.PriceValidityPeriodEndDate = common.TimeToTimestamp(priceValidityPeriodEndDate.UTC().Truncate(time.Second))

	crUpdUser := commonproto.CrUpdUser{}
	crUpdUser.StatusCode = "active"
	crUpdUser.CreatedByUserId = user.Id
	crUpdUser.UpdatedByUserId = user.Id

	crUpdTime := commonproto.CrUpdTime{}
	crUpdTime.CreatedAt = tn
	crUpdTime.UpdatedAt = tn

	debitNoteLine := invoiceproto.DebitNoteLine{DebitNoteLineD: &debitNoteLineD, DebitNoteLineT: &debitNoteLineT, CrUpdUser: &crUpdUser, CrUpdTime: &crUpdTime}

	return &debitNoteLine, nil
}

// insertDebitNoteLine - Insert DebitNoteLine details into database
func (ds *DebitNoteHeaderService) insertDebitNoteLine(ctx context.Context, insertDebitNoteLineSQL string, debitNoteLine *invoiceproto.DebitNoteLine, userEmail string, requestID string) error {
	debitNoteLineTmp, err := ds.crDebitNoteLineStruct(ctx, debitNoteLine, userEmail, requestID)
	if err != nil {
		ds.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
		return err
	}
	err = ds.DBService.InsUpd(ctx, userEmail, requestID, func(tx *sqlx.Tx) error {
		res, err := tx.NamedExecContext(ctx, insertDebitNoteLineSQL, debitNoteLineTmp)
		if err != nil {
			ds.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}

		uID, err := res.LastInsertId()
		if err != nil {
			ds.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}
		debitNoteLine.DebitNoteLineD.Id = uint32(uID)
		uuid4Str, err := common.UUIDBytesToStr(debitNoteLine.DebitNoteLineD.Uuid4)
		if err != nil {
			ds.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}
		debitNoteLine.DebitNoteLineD.IdS = uuid4Str
		return nil
	})

	if err != nil {
		ds.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
		return err
	}
	return nil
}

// crDebitNoteLineStruct - process DebitNoteLine details
func (ds *DebitNoteHeaderService) crDebitNoteLineStruct(ctx context.Context, debitNoteLine *invoiceproto.DebitNoteLine, userEmail string, requestID string) (*invoicestruct.DebitNoteLine, error) {
	crUpdTime := new(commonstruct.CrUpdTime)
	crUpdTime.CreatedAt = common.TimestampToTime(debitNoteLine.CrUpdTime.CreatedAt)
	crUpdTime.UpdatedAt = common.TimestampToTime(debitNoteLine.CrUpdTime.UpdatedAt)

	debitNoteLineT := new(invoicestruct.DebitNoteLineT)
	debitNoteLineT.TaxPointDate = common.TimestampToTime(debitNoteLine.DebitNoteLineT.TaxPointDate)
	debitNoteLineT.PriceValidityPeriodStartDate = common.TimestampToTime(debitNoteLine.DebitNoteLineT.PriceValidityPeriodStartDate)
	debitNoteLineT.PriceValidityPeriodEndDate = common.TimestampToTime(debitNoteLine.DebitNoteLineT.PriceValidityPeriodEndDate)

	debitNoteLineTmp := invoicestruct.DebitNoteLine{DebitNoteLineD: debitNoteLine.DebitNoteLineD, DebitNoteLineT: debitNoteLineT, CrUpdUser: debitNoteLine.CrUpdUser, CrUpdTime: crUpdTime}

	return &debitNoteLineTmp, nil
}

// GetDebitNoteLines - GetDebitNoteLines
func (ds *DebitNoteHeaderService) GetDebitNoteLines(ctx context.Context, inReq *invoiceproto.GetDebitNoteLinesRequest) (*invoiceproto.GetDebitNoteLinesResponse, error) {
	in := inReq.GetRequest
	getRequest := commonproto.GetRequest{}
	getRequest.Id = in.Id
	getRequest.UserEmail = in.UserEmail
	getRequest.RequestId = in.RequestId
	form := invoiceproto.GetDebitNoteHeaderRequest{}
	form.GetRequest = &getRequest

	debitNoteHeaderResponse, err := ds.GetDebitNoteHeader(ctx, &form)
	if err != nil {
		ds.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	debitNoteHeader := debitNoteHeaderResponse.DebitNoteHeader
	debitNoteLines := []*invoiceproto.DebitNoteLine{}

	nselectDebitNoteLinesSQL := selectDebitNoteLinesSQL + ` where debit_note_header_id = ? and status_code = ?;`
	rows, err := ds.DBService.DB.QueryxContext(ctx, nselectDebitNoteLinesSQL, debitNoteHeader.DebitNoteHeaderD.Id, "active")
	if err != nil {
		ds.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	for rows.Next() {

		debitNoteLineTmp := invoicestruct.DebitNoteLine{}
		err = rows.StructScan(&debitNoteLineTmp)
		if err != nil {
			ds.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
			return nil, err
		}
		debitNoteLine, err := ds.getDebitNoteLineStruct(ctx, &getRequest, debitNoteLineTmp)
		if err != nil {
			ds.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
			return nil, err
		}

		debitNoteLines = append(debitNoteLines, debitNoteLine)
	}

	debitNoteLinesResponse := invoiceproto.GetDebitNoteLinesResponse{}
	debitNoteLinesResponse.DebitNoteLines = debitNoteLines
	return &debitNoteLinesResponse, nil
}

// getDebitNoteLineStruct - Get debit note Line
func (ds *DebitNoteHeaderService) getDebitNoteLineStruct(ctx context.Context, in *commonproto.GetRequest, debitNoteLineTmp invoicestruct.DebitNoteLine) (*invoiceproto.DebitNoteLine, error) {
	debitNoteLineT := new(invoiceproto.DebitNoteLineT)
	debitNoteLineT.TaxPointDate = common.TimeToTimestamp(debitNoteLineTmp.DebitNoteLineT.TaxPointDate)
	debitNoteLineT.PriceValidityPeriodStartDate = common.TimeToTimestamp(debitNoteLineTmp.DebitNoteLineT.PriceValidityPeriodStartDate)
	debitNoteLineT.PriceValidityPeriodEndDate = common.TimeToTimestamp(debitNoteLineTmp.DebitNoteLineT.PriceValidityPeriodEndDate)

	uuid4Str, err := common.UUIDBytesToStr(debitNoteLineTmp.DebitNoteLineD.Uuid4)
	if err != nil {
		ds.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	debitNoteLineTmp.DebitNoteLineD.IdS = uuid4Str

	crUpdTime := new(commonproto.CrUpdTime)
	crUpdTime.CreatedAt = common.TimeToTimestamp(debitNoteLineTmp.CrUpdTime.CreatedAt)
	crUpdTime.UpdatedAt = common.TimeToTimestamp(debitNoteLineTmp.CrUpdTime.UpdatedAt)

	debitNoteLine := invoiceproto.DebitNoteLine{DebitNoteLineD: debitNoteLineTmp.DebitNoteLineD, DebitNoteLineT: debitNoteLineT, CrUpdUser: debitNoteLineTmp.CrUpdUser, CrUpdTime: crUpdTime}

	return &debitNoteLine, nil
}
