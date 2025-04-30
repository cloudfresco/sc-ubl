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

const insertInvoiceLineSQL = `insert into invoice_lines
	  (uuid4,
il_id,
note,
invoiced_quantity,
line_extension_amount,
accounting_cost_code,
accounting_cost,
payment_purpose_code,
free_of_charge_indicator,
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
invoice_header_id,
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
  values (
:uuid4,
:il_id,
:note,
:invoiced_quantity,
:line_extension_amount,
:accounting_cost_code,
:accounting_cost,
:payment_purpose_code,
:free_of_charge_indicator,
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
:invoice_header_id,
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

const selectInvoiceLinesSQL = `select 
 id,
uuid4,
il_id,
note,
invoiced_quantity,
line_extension_amount,
accounting_cost_code,
accounting_cost,
payment_purpose_code,
free_of_charge_indicator,
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
invoice_header_id,
tax_point_date,
invoice_period_start_date,
invoice_period_end_date,
price_validity_period_start_date,
price_validity_period_end_date,
status_code,
created_by_user_id,
updated_by_user_id,
created_at,
updated_at from invoice_lines`

// CreateInvoiceLine - Create InvoiceLine
func (is *InvoiceService) CreateInvoiceLine(ctx context.Context, in *invoiceproto.CreateInvoiceLineRequest) (*invoiceproto.CreateInvoiceLineResponse, error) {
	invoiceLine, err := is.ProcessInvoiceLineRequest(ctx, in)
	if err != nil {
		is.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	err = is.insertInvoiceLine(ctx, insertInvoiceLineSQL, invoiceLine, in.GetUserEmail(), in.GetRequestId())
	if err != nil {
		is.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	invoiceLineResponse := invoiceproto.CreateInvoiceLineResponse{}
	invoiceLineResponse.InvoiceLine = invoiceLine
	return &invoiceLineResponse, nil
}

func (is *InvoiceService) ProcessInvoiceLineRequest(ctx context.Context, in *invoiceproto.CreateInvoiceLineRequest) (*invoiceproto.InvoiceLine, error) {
	getRequest := commonproto.GetRequest{}
	getRequest.Id = in.UserId
	getRequest.UserEmail = in.UserEmail
	getRequest.RequestId = in.RequestId
	ctxNew, err := common.CreateCtxJWT(ctx)
	if err != nil {
		is.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	form := partyproto.GetUserRequest{}
	form.GetRequest = &getRequest
	userResponse, err := is.UserServiceClient.GetUser(ctxNew, &form)
	if err != nil {
		is.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	user := userResponse.User

	taxPointDate, err := time.Parse(common.Layout, in.TaxPointDate)
	if err != nil {
		is.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	invoicePeriodStartDate, err := time.Parse(common.Layout, in.InvoicePeriodStartDate)
	if err != nil {
		is.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	invoicePeriodEndDate, err := time.Parse(common.Layout, in.InvoicePeriodEndDate)
	if err != nil {
		is.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	priceValidityPeriodStartDate, err := time.Parse(common.Layout, in.PriceValidityPeriodStartDate)
	if err != nil {
		is.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	priceValidityPeriodEndDate, err := time.Parse(common.Layout, in.PriceValidityPeriodEndDate)
	if err != nil {
		is.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	ttime := common.GetTimeDetails()
	tn := common.TimeToTimestamp(ttime)

	invoiceLineD := invoiceproto.InvoiceLineD{}
	invoiceLineD.Uuid4, err = common.GetUUIDBytes()
	if err != nil {
		is.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	invoiceLineD.IlId = in.IlId
	invoiceLineD.Note = in.Note
	invoiceLineD.InvoicedQuantity = in.InvoicedQuantity
	invoiceLineD.LineExtensionAmount = in.LineExtensionAmount
	invoiceLineD.AccountingCostCode = in.AccountingCostCode
	invoiceLineD.AccountingCost = in.AccountingCost
	invoiceLineD.PaymentPurposeCode = in.PaymentPurposeCode
	invoiceLineD.FreeOfChargeIndicator = in.FreeOfChargeIndicator
	invoiceLineD.OrderLineId = in.OrderLineId
	invoiceLineD.DespatchLineId = in.DespatchLineId
	invoiceLineD.ReceiptLineId = in.ReceiptLineId
	invoiceLineD.BillingId = in.BillingId
	invoiceLineD.OriginatorPartyId = in.OriginatorPartyId
	invoiceLineD.ItemId = in.ItemId
	invoiceLineD.PriceAmount = in.PriceAmount
	invoiceLineD.PriceBaseQuantity = in.PriceBaseQuantity
	invoiceLineD.PriceChangeReason = in.PriceChangeReason
	invoiceLineD.PriceTypeCode = in.PriceTypeCode
	invoiceLineD.PriceType = in.PriceType
	invoiceLineD.OrderableUnitFactorRate = in.OrderableUnitFactorRate
	invoiceLineD.PriceListId = in.PriceListId
	invoiceLineD.InvoiceHeaderId = in.InvoiceHeaderId

	invoiceLineT := invoiceproto.InvoiceLineT{}
	invoiceLineT.TaxPointDate = common.TimeToTimestamp(taxPointDate.UTC().Truncate(time.Second))
	invoiceLineT.InvoicePeriodStartDate = common.TimeToTimestamp(invoicePeriodStartDate.UTC().Truncate(time.Second))
	invoiceLineT.InvoicePeriodEndDate = common.TimeToTimestamp(invoicePeriodEndDate.UTC().Truncate(time.Second))
	invoiceLineT.PriceValidityPeriodStartDate = common.TimeToTimestamp(priceValidityPeriodStartDate.UTC().Truncate(time.Second))
	invoiceLineT.PriceValidityPeriodEndDate = common.TimeToTimestamp(priceValidityPeriodEndDate.UTC().Truncate(time.Second))

	crUpdUser := commonproto.CrUpdUser{}
	crUpdUser.StatusCode = "active"
	crUpdUser.CreatedByUserId = user.Id
	crUpdUser.UpdatedByUserId = user.Id
	crUpdTime := commonproto.CrUpdTime{}
	crUpdTime.CreatedAt = tn
	crUpdTime.UpdatedAt = tn

	invoiceLine := invoiceproto.InvoiceLine{InvoiceLineD: &invoiceLineD, InvoiceLineT: &invoiceLineT, CrUpdUser: &crUpdUser, CrUpdTime: &crUpdTime}

	return &invoiceLine, nil
}

func (is *InvoiceService) insertInvoiceLine(ctx context.Context, insertInvoiceLineSQL string, invoiceLine *invoiceproto.InvoiceLine, userEmail string, requestID string) error {
	invoiceLineTmp, err := is.crInvoiceLineStruct(ctx, invoiceLine, userEmail, requestID)
	if err != nil {
		is.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
		return err
	}
	err = is.DBService.InsUpd(ctx, userEmail, requestID, func(tx *sqlx.Tx) error {
		res, err := tx.NamedExecContext(ctx, insertInvoiceLineSQL, invoiceLineTmp)
		if err != nil {
			is.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}

		uID, err := res.LastInsertId()
		if err != nil {
			is.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}
		invoiceLine.InvoiceLineD.Id = uint32(uID)
		uuid4Str, err := common.UUIDBytesToStr(invoiceLine.InvoiceLineD.Uuid4)
		if err != nil {
			is.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}
		invoiceLine.InvoiceLineD.IdS = uuid4Str
		return nil
	})

	if err != nil {
		is.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
		return err
	}
	return nil
}

// crInvoiceLineStruct - process InvoiceLine details
func (is *InvoiceService) crInvoiceLineStruct(ctx context.Context, invoiceLine *invoiceproto.InvoiceLine, userEmail string, requestID string) (*invoicestruct.InvoiceLine, error) {
	crUpdTime := new(commonstruct.CrUpdTime)
	crUpdTime.CreatedAt = common.TimestampToTime(invoiceLine.CrUpdTime.CreatedAt)
	crUpdTime.UpdatedAt = common.TimestampToTime(invoiceLine.CrUpdTime.UpdatedAt)

	invoiceLineT := new(invoicestruct.InvoiceLineT)
	invoiceLineT.TaxPointDate = common.TimestampToTime(invoiceLine.InvoiceLineT.TaxPointDate)
	invoiceLineT.InvoicePeriodStartDate = common.TimestampToTime(invoiceLine.InvoiceLineT.InvoicePeriodStartDate)
	invoiceLineT.InvoicePeriodEndDate = common.TimestampToTime(invoiceLine.InvoiceLineT.InvoicePeriodEndDate)
	invoiceLineT.PriceValidityPeriodStartDate = common.TimestampToTime(invoiceLine.InvoiceLineT.PriceValidityPeriodStartDate)
	invoiceLineT.PriceValidityPeriodEndDate = common.TimestampToTime(invoiceLine.InvoiceLineT.PriceValidityPeriodEndDate)

	invoiceLineTmp := invoicestruct.InvoiceLine{InvoiceLineD: invoiceLine.InvoiceLineD, InvoiceLineT: invoiceLineT, CrUpdUser: invoiceLine.CrUpdUser, CrUpdTime: crUpdTime}

	return &invoiceLineTmp, nil
}

// GetInvoiceLines - GetInvoiceLines
func (is *InvoiceService) GetInvoiceLines(ctx context.Context, inReq *invoiceproto.GetInvoiceLinesRequest) (*invoiceproto.GetInvoiceLinesResponse, error) {
	in := inReq.GetRequest
	getRequest := commonproto.GetRequest{}
	getRequest.Id = in.Id
	getRequest.UserEmail = in.UserEmail
	getRequest.RequestId = in.RequestId
	form := invoiceproto.GetInvoiceRequest{}
	form.GetRequest = &getRequest

	invoiceHeaderResponse, err := is.GetInvoice(ctx, &form)
	if err != nil {
		is.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	invoiceHeader := invoiceHeaderResponse.InvoiceHeader

	invoiceLines := []*invoiceproto.InvoiceLine{}

	nselectInvoiceLinesSQL := selectInvoiceLinesSQL + ` where invoice_header_id = ? and status_code = ?;`
	rows, err := is.DBService.DB.QueryxContext(ctx, nselectInvoiceLinesSQL, invoiceHeader.InvoiceHeaderD.Id, "active")
	if err != nil {
		is.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	for rows.Next() {

		invoiceLineTmp := invoicestruct.InvoiceLine{}
		err = rows.StructScan(&invoiceLineTmp)
		if err != nil {
			is.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
			return nil, err
		}

		invoiceLine, err := is.getInvoiceLineStruct(ctx, &getRequest, invoiceLineTmp)
		if err != nil {
			is.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
			return nil, err
		}

		invoiceLines = append(invoiceLines, invoiceLine)
	}
	invoiceLinesResponse := invoiceproto.GetInvoiceLinesResponse{}
	invoiceLinesResponse.InvoiceLines = invoiceLines
	return &invoiceLinesResponse, nil
}

// getInvoiceLineStruct - Get invoice note Line
func (is *InvoiceService) getInvoiceLineStruct(ctx context.Context, in *commonproto.GetRequest, invoiceLineTmp invoicestruct.InvoiceLine) (*invoiceproto.InvoiceLine, error) {
	invoiceLineT := new(invoiceproto.InvoiceLineT)
	invoiceLineT.TaxPointDate = common.TimeToTimestamp(invoiceLineTmp.InvoiceLineT.TaxPointDate)
	invoiceLineT.InvoicePeriodStartDate = common.TimeToTimestamp(invoiceLineTmp.InvoiceLineT.InvoicePeriodStartDate)
	invoiceLineT.InvoicePeriodEndDate = common.TimeToTimestamp(invoiceLineTmp.InvoiceLineT.InvoicePeriodEndDate)
	invoiceLineT.PriceValidityPeriodStartDate = common.TimeToTimestamp(invoiceLineTmp.InvoiceLineT.PriceValidityPeriodStartDate)
	invoiceLineT.PriceValidityPeriodEndDate = common.TimeToTimestamp(invoiceLineTmp.InvoiceLineT.PriceValidityPeriodEndDate)

	uuid4Str, err := common.UUIDBytesToStr(invoiceLineTmp.InvoiceLineD.Uuid4)
	if err != nil {
		is.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	invoiceLineTmp.InvoiceLineD.IdS = uuid4Str

	crUpdTime := new(commonproto.CrUpdTime)
	crUpdTime.CreatedAt = common.TimeToTimestamp(invoiceLineTmp.CrUpdTime.CreatedAt)
	crUpdTime.UpdatedAt = common.TimeToTimestamp(invoiceLineTmp.CrUpdTime.UpdatedAt)

	invoiceLine := invoiceproto.InvoiceLine{InvoiceLineD: invoiceLineTmp.InvoiceLineD, InvoiceLineT: invoiceLineT, CrUpdUser: invoiceLineTmp.CrUpdUser, CrUpdTime: crUpdTime}

	return &invoiceLine, nil
}
