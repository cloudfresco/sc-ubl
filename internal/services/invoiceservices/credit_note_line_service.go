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

const insertCreditNoteLineSQL = `insert into credit_note_lines
	  (uuid4,
cnl_id,
note,
credited_quantity,
line_extension_amount,
accounting_cost_code,
accounting_cost,
payment_purpose_code,
free_of_charge_indicator,
discrepancy_response,
order_line_id,
despatch_line_id,
receipt_line_id,
billing_id,
originator_party_id,
item_id,
price_amount,
price_base_quantity,
price_change_reason,
price_type_code,
price_type,
orderable_unit_factor_rate,
price_list_id,
credit_note_header_id,
tax_point_date,
invoice_period_start_date,
invoice_period_end_date,
price_validity_period_start_date,
price_validity_period_end_date,
status_code,
created_by_user_id,
updated_by_user_id,
created_at,
updated_at)
  values (:uuid4,
:cnl_id,
:note,
:credited_quantity,
:line_extension_amount,
:accounting_cost_code,
:accounting_cost,
:payment_purpose_code,
:free_of_charge_indicator,
:discrepancy_response,
:order_line_id,
:despatch_line_id,
:receipt_line_id,
:billing_id,
:originator_party_id,
:item_id,
:price_amount,
:price_base_quantity,
:price_change_reason,
:price_type_code,
:price_type,
:orderable_unit_factor_rate,
:price_list_id,
:credit_note_header_id,
:tax_point_date,
:invoice_period_start_date,
:invoice_period_end_date,
:price_validity_period_start_date,
:price_validity_period_end_date,
:status_code,
:created_by_user_id,
:updated_by_user_id,
:created_at,
:updated_at);`

const selectCreditNoteLinesSQL = `select 
 id,
uuid4,
cnl_id,
note,
credited_quantity,
line_extension_amount,
accounting_cost_code,
accounting_cost,
payment_purpose_code,
free_of_charge_indicator,
discrepancy_response,
order_line_id,
despatch_line_id,
receipt_line_id,
billing_id,
originator_party_id,
item_id,
price_amount,
price_base_quantity,
price_change_reason,
price_type_code,
price_type,
orderable_unit_factor_rate,
price_list_id,
credit_note_header_id,
tax_point_date,
invoice_period_start_date,
invoice_period_end_date,
price_validity_period_start_date,
price_validity_period_end_date,
status_code,
created_by_user_id,
updated_by_user_id,
created_at,
updated_at from credit_note_lines`

// CreateCreditNoteLine - Create CreditNoteLine
func (cs *CreditNoteHeaderService) CreateCreditNoteLine(ctx context.Context, in *invoiceproto.CreateCreditNoteLineRequest) (*invoiceproto.CreateCreditNoteLineResponse, error) {
	creditNoteLine, err := cs.ProcessCreditNoteLineRequest(ctx, in)
	if err != nil {
		cs.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	err = cs.insertCreditNoteLine(ctx, insertCreditNoteLineSQL, creditNoteLine, in.GetUserEmail(), in.GetRequestId())
	if err != nil {
		cs.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	creditNoteLineResponse := invoiceproto.CreateCreditNoteLineResponse{}
	creditNoteLineResponse.CreditNoteLine = creditNoteLine
	return &creditNoteLineResponse, nil
}

func (cs *CreditNoteHeaderService) ProcessCreditNoteLineRequest(ctx context.Context, in *invoiceproto.CreateCreditNoteLineRequest) (*invoiceproto.CreditNoteLine, error) {
	getRequest := commonproto.GetRequest{}
	getRequest.Id = in.UserId
	getRequest.UserEmail = in.UserEmail
	getRequest.RequestId = in.RequestId
	ctxNew, err := common.CreateCtxJWT(ctx)
	if err != nil {
		cs.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	form := partyproto.GetUserRequest{}
	form.GetRequest = &getRequest
	userResponse, err := cs.UserServiceClient.GetUser(ctxNew, &form)
	if err != nil {
		cs.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	user := userResponse.User

	taxPointDate, err := time.Parse(common.Layout, in.TaxPointDate)
	if err != nil {
		cs.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	invoicePeriodStartDate, err := time.Parse(common.Layout, in.InvoicePeriodStartDate)
	if err != nil {
		cs.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	invoicePeriodEndDate, err := time.Parse(common.Layout, in.InvoicePeriodEndDate)
	if err != nil {
		cs.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	priceValidityPeriodStartDate, err := time.Parse(common.Layout, in.PriceValidityPeriodStartDate)
	if err != nil {
		cs.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	priceValidityPeriodEndDate, err := time.Parse(common.Layout, in.PriceValidityPeriodEndDate)
	if err != nil {
		cs.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	ttime := common.GetTimeDetails()
	tn := common.TimeToTimestamp(ttime)

	creditNoteLineD := invoiceproto.CreditNoteLineD{}
	creditNoteLineD.Uuid4, err = common.GetUUIDBytes()
	if err != nil {
		cs.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	creditNoteLineD.CnlId = in.CnlId
	creditNoteLineD.Note = in.Note
	creditNoteLineD.CreditedQuantity = in.CreditedQuantity
	creditNoteLineD.LineExtensionAmount = in.LineExtensionAmount
	creditNoteLineD.AccountingCostCode = in.AccountingCostCode
	creditNoteLineD.AccountingCost = in.AccountingCost
	creditNoteLineD.PaymentPurposeCode = in.PaymentPurposeCode
	creditNoteLineD.FreeOfChargeIndicator = in.FreeOfChargeIndicator
	creditNoteLineD.DiscrepancyResponse = in.DiscrepancyResponse
	creditNoteLineD.OrderLineId = in.OrderLineId
	creditNoteLineD.DespatchLineId = in.DespatchLineId
	creditNoteLineD.ReceiptLineId = in.ReceiptLineId
	creditNoteLineD.BillingId = in.BillingId
	creditNoteLineD.OriginatorPartyId = in.OriginatorPartyId
	creditNoteLineD.ItemId = in.ItemId
	creditNoteLineD.PriceAmount = in.PriceAmount
	creditNoteLineD.PriceBaseQuantity = in.PriceBaseQuantity
	creditNoteLineD.PriceChangeReason = in.PriceChangeReason
	creditNoteLineD.PriceTypeCode = in.PriceTypeCode
	creditNoteLineD.PriceType = in.PriceType
	creditNoteLineD.OrderableUnitFactorRate = in.OrderableUnitFactorRate
	creditNoteLineD.PriceListId = in.PriceListId
	creditNoteLineD.CreditNoteHeaderId = in.CreditNoteHeaderId

	creditNoteLineT := invoiceproto.CreditNoteLineT{}
	creditNoteLineT.TaxPointDate = common.TimeToTimestamp(taxPointDate.UTC().Truncate(time.Second))
	creditNoteLineT.InvoicePeriodStartDate = common.TimeToTimestamp(invoicePeriodStartDate.UTC().Truncate(time.Second))
	creditNoteLineT.InvoicePeriodEndDate = common.TimeToTimestamp(invoicePeriodEndDate.UTC().Truncate(time.Second))
	creditNoteLineT.PriceValidityPeriodStartDate = common.TimeToTimestamp(priceValidityPeriodStartDate.UTC().Truncate(time.Second))
	creditNoteLineT.PriceValidityPeriodEndDate = common.TimeToTimestamp(priceValidityPeriodEndDate.UTC().Truncate(time.Second))

	crUpdUser := commonproto.CrUpdUser{}
	crUpdUser.StatusCode = "active"
	crUpdUser.CreatedByUserId = user.Id
	crUpdUser.UpdatedByUserId = user.Id

	crUpdTime := commonproto.CrUpdTime{}
	crUpdTime.CreatedAt = tn
	crUpdTime.UpdatedAt = tn

	creditNoteLine := invoiceproto.CreditNoteLine{CreditNoteLineD: &creditNoteLineD, CreditNoteLineT: &creditNoteLineT, CrUpdUser: &crUpdUser, CrUpdTime: &crUpdTime}

	err = cs.insertCreditNoteLine(ctx, insertCreditNoteLineSQL, &creditNoteLine, in.GetUserEmail(), in.GetRequestId())
	if err != nil {
		cs.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	return &creditNoteLine, nil
}

// insertCreditNoteLine - Insert CreditNoteLine details into database
func (cs *CreditNoteHeaderService) insertCreditNoteLine(ctx context.Context, insertCreditNoteLineSQL string, creditNoteLine *invoiceproto.CreditNoteLine, userEmail string, requestID string) error {
	creditNoteLineTmp, err := cs.crCreditNoteLineStruct(ctx, creditNoteLine, userEmail, requestID)
	if err != nil {
		cs.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
		return err
	}
	err = cs.DBService.InsUpd(ctx, userEmail, requestID, func(tx *sqlx.Tx) error {
		res, err := tx.NamedExecContext(ctx, insertCreditNoteLineSQL, creditNoteLineTmp)
		if err != nil {
			cs.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}

		uID, err := res.LastInsertId()
		if err != nil {
			cs.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}
		creditNoteLine.CreditNoteLineD.Id = uint32(uID)
		uuid4Str, err := common.UUIDBytesToStr(creditNoteLine.CreditNoteLineD.Uuid4)
		if err != nil {
			cs.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}
		creditNoteLine.CreditNoteLineD.IdS = uuid4Str
		return nil
	})

	if err != nil {
		cs.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
		return err
	}
	return nil
}

// crCreditNoteLineStruct - process CreditNoteLine details
func (cs *CreditNoteHeaderService) crCreditNoteLineStruct(ctx context.Context, creditNoteLine *invoiceproto.CreditNoteLine, userEmail string, requestID string) (*invoicestruct.CreditNoteLine, error) {
	crUpdTime := new(commonstruct.CrUpdTime)
	crUpdTime.CreatedAt = common.TimestampToTime(creditNoteLine.CrUpdTime.CreatedAt)
	crUpdTime.UpdatedAt = common.TimestampToTime(creditNoteLine.CrUpdTime.UpdatedAt)

	creditNoteLineT := new(invoicestruct.CreditNoteLineT)
	creditNoteLineT.TaxPointDate = common.TimestampToTime(creditNoteLine.CreditNoteLineT.TaxPointDate)
	creditNoteLineT.InvoicePeriodStartDate = common.TimestampToTime(creditNoteLine.CreditNoteLineT.InvoicePeriodStartDate)
	creditNoteLineT.InvoicePeriodEndDate = common.TimestampToTime(creditNoteLine.CreditNoteLineT.InvoicePeriodEndDate)
	creditNoteLineT.PriceValidityPeriodStartDate = common.TimestampToTime(creditNoteLine.CreditNoteLineT.PriceValidityPeriodStartDate)
	creditNoteLineT.PriceValidityPeriodEndDate = common.TimestampToTime(creditNoteLine.CreditNoteLineT.PriceValidityPeriodEndDate)

	creditNoteLineTmp := invoicestruct.CreditNoteLine{CreditNoteLineD: creditNoteLine.CreditNoteLineD, CreditNoteLineT: creditNoteLineT, CrUpdUser: creditNoteLine.CrUpdUser, CrUpdTime: crUpdTime}

	return &creditNoteLineTmp, nil
}

// GetCreditNoteLines - GetCreditNoteLines
func (cs *CreditNoteHeaderService) GetCreditNoteLines(ctx context.Context, inReq *invoiceproto.GetCreditNoteLinesRequest) (*invoiceproto.GetCreditNoteLinesResponse, error) {
	in := inReq.GetRequest
	getRequest := commonproto.GetRequest{}
	getRequest.Id = in.Id
	getRequest.UserEmail = in.UserEmail
	getRequest.RequestId = in.RequestId
	form := invoiceproto.GetCreditNoteHeaderRequest{}
	form.GetRequest = &getRequest

	creditNoteHeaderResponse, err := cs.GetCreditNoteHeader(ctx, &form)
	if err != nil {
		cs.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	creditNoteHeader := creditNoteHeaderResponse.CreditNoteHeader
	creditNoteLines := []*invoiceproto.CreditNoteLine{}

	nselectCreditNoteLinesSQL := selectCreditNoteLinesSQL + ` where credit_note_header_id = ? and status_code = ?;`
	rows, err := cs.DBService.DB.QueryxContext(ctx, nselectCreditNoteLinesSQL, creditNoteHeader.CreditNoteHeaderD.Id, "active")
	if err != nil {
		cs.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	for rows.Next() {

		creditNoteLineTmp := invoicestruct.CreditNoteLine{}
		err = rows.StructScan(&creditNoteLineTmp)
		if err != nil {
			cs.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
			return nil, err
		}

		creditNoteLine, err := cs.getCreditNoteLineStruct(ctx, &getRequest, creditNoteLineTmp)
		if err != nil {
			cs.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
			return nil, err
		}

		creditNoteLines = append(creditNoteLines, creditNoteLine)
	}

	creditNoteLinesResponse := invoiceproto.GetCreditNoteLinesResponse{}
	creditNoteLinesResponse.CreditNoteLines = creditNoteLines
	return &creditNoteLinesResponse, nil
}

// getCreditNoteLineStruct - Get credit note Line
func (cs *CreditNoteHeaderService) getCreditNoteLineStruct(ctx context.Context, in *commonproto.GetRequest, creditNoteLineTmp invoicestruct.CreditNoteLine) (*invoiceproto.CreditNoteLine, error) {
	creditNoteLineT := new(invoiceproto.CreditNoteLineT)
	creditNoteLineT.TaxPointDate = common.TimeToTimestamp(creditNoteLineTmp.CreditNoteLineT.TaxPointDate)
	creditNoteLineT.InvoicePeriodStartDate = common.TimeToTimestamp(creditNoteLineTmp.CreditNoteLineT.InvoicePeriodStartDate)
	creditNoteLineT.InvoicePeriodEndDate = common.TimeToTimestamp(creditNoteLineTmp.CreditNoteLineT.InvoicePeriodEndDate)
	creditNoteLineT.PriceValidityPeriodStartDate = common.TimeToTimestamp(creditNoteLineTmp.CreditNoteLineT.PriceValidityPeriodStartDate)
	creditNoteLineT.PriceValidityPeriodEndDate = common.TimeToTimestamp(creditNoteLineTmp.CreditNoteLineT.PriceValidityPeriodEndDate)

	uuid4Str, err := common.UUIDBytesToStr(creditNoteLineTmp.CreditNoteLineD.Uuid4)
	if err != nil {
		cs.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	creditNoteLineTmp.CreditNoteLineD.IdS = uuid4Str

	crUpdTime := new(commonproto.CrUpdTime)
	crUpdTime.CreatedAt = common.TimeToTimestamp(creditNoteLineTmp.CrUpdTime.CreatedAt)
	crUpdTime.UpdatedAt = common.TimeToTimestamp(creditNoteLineTmp.CrUpdTime.UpdatedAt)

	creditNoteLine := invoiceproto.CreditNoteLine{CreditNoteLineD: creditNoteLineTmp.CreditNoteLineD, CreditNoteLineT: creditNoteLineT, CrUpdUser: creditNoteLineTmp.CrUpdUser, CrUpdTime: crUpdTime}

	return &creditNoteLine, nil
}
