package invoiceservices

import (
	"context"
	"time"

	"github.com/cloudfresco/sc-ubl/internal/common"
	commonproto "github.com/cloudfresco/sc-ubl/internal/protogen/common/v1"
	invoiceproto "github.com/cloudfresco/sc-ubl/internal/protogen/invoice/v1"
	partyproto "github.com/cloudfresco/sc-ubl/internal/protogen/party/v1"
	partyservice "github.com/cloudfresco/sc-ubl/internal/services/partyservices"
	commonstruct "github.com/cloudfresco/sc-ubl/internal/servicestructs/common/v1"
	invoicestruct "github.com/cloudfresco/sc-ubl/internal/servicestructs/invoice/v1"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// DebitNoteHeaderService - For accessing DebitNoteHeader services
type DebitNoteHeaderService struct {
	log               *zap.Logger
	DBService         *common.DBService
	RedisService      *common.RedisService
	UserServiceClient partyproto.UserServiceClient
	invoiceproto.UnimplementedDebitNoteHeaderServiceServer
}

// NewDebitNoteHeaderService - Create DebitNoteHeader service
func NewDebitNoteHeaderService(log *zap.Logger, dbOpt *common.DBService, redisOpt *common.RedisService, userServiceClient partyproto.UserServiceClient) *DebitNoteHeaderService {
	return &DebitNoteHeaderService{
		log:               log,
		DBService:         dbOpt,
		RedisService:      redisOpt,
		UserServiceClient: userServiceClient,
	}
}

const insertDebitNoteHeaderSQL = `insert into debit_note_headers
	  (uuid4,
dnh_id,
note,
document_currency_code,
tax_currency_code,
pricing_currency_code,
payment_currency_code,
payment_alt_currency_code,
accounting_cost_code,
accounting_cost,
line_count_numeric,
discrepancy_response,
order_id,
billing_id,
despatch_id,
receipt_id,
statement_id,
contract_id,
accounting_supplier_party_id,
accounting_customer_party_id,
payee_party_id,
buyer_customer_party_id,
seller_supplier_party_id,
tax_representative_party_id,
tax_ex_source_currency_code,
tax_ex_source_currency_base_rate,
tax_ex_target_currency_code,
tax_ex_target_currency_base_rate,
tax_ex_exchange_market_id,
tax_ex_calculation_rate,
tax_ex_mathematic_operator_code,
pricing_ex_source_currency_code,
pricing_ex_source_currency_base_rate,
pricing_ex_target_currency_code,
pricing_ex_target_currency_base_rate,
pricing_ex_exchange_market_id,
pricing_ex_calculation_rate,
pricing_ex_mathematic_operator_code,
payment_ex_source_currency_code,
payment_ex_source_currency_base_rate,
payment_ex_target_currency_code,
payment_ex_target_currency_base_rate,
payment_ex_exchange_market_id,
payment_ex_calculation_rate,
payment_ex_mathematic_operator_code,
payment_alt_ex_source_currency_code,
payment_alt_ex_source_currency_base_rate,
payment_alt_ex_target_currency_code,
payment_alt_ex_target_currency_base_rate,
payment_alt_ex_exchange_market_id,
payment_alt_ex_calculation_rate,
payment_alt_ex_mathematic_operator_code,
line_extension_amount,
tax_exclusive_amount,
tax_inclusive_amount,
allowance_total_amount,
charge_total_amount,
withholding_tax_total_amount,
prepaid_amount,
payable_rounding_amount,
payable_amount,
payable_alternative_amount,
issue_date,
tax_point_date,
invoice_period_start_date,
invoice_period_end_date,
tax_ex_date,
pricing_ex_date,
payment_ex_date,
payment_alt_ex_date,
status_code,
created_by_user_id,
updated_by_user_id,
created_at,
updated_at)
        values(:uuid4,
:dnh_id,
:note,
:document_currency_code,
:tax_currency_code,
:pricing_currency_code,
:payment_currency_code,
:payment_alt_currency_code,
:accounting_cost_code,
:accounting_cost,
:line_count_numeric,
:discrepancy_response,
:order_id,
:billing_id,
:despatch_id,
:receipt_id,
:statement_id,
:contract_id,
:accounting_supplier_party_id,
:accounting_customer_party_id,
:payee_party_id,
:buyer_customer_party_id,
:seller_supplier_party_id,
:tax_representative_party_id,
:tax_ex_source_currency_code,
:tax_ex_source_currency_base_rate,
:tax_ex_target_currency_code,
:tax_ex_target_currency_base_rate,
:tax_ex_exchange_market_id,
:tax_ex_calculation_rate,
:tax_ex_mathematic_operator_code,
:pricing_ex_source_currency_code,
:pricing_ex_source_currency_base_rate,
:pricing_ex_target_currency_code,
:pricing_ex_target_currency_base_rate,
:pricing_ex_exchange_market_id,
:pricing_ex_calculation_rate,
:pricing_ex_mathematic_operator_code,
:payment_ex_source_currency_code,
:payment_ex_source_currency_base_rate,
:payment_ex_target_currency_code,
:payment_ex_target_currency_base_rate,
:payment_ex_exchange_market_id,
:payment_ex_calculation_rate,
:payment_ex_mathematic_operator_code,
:payment_alt_ex_source_currency_code,
:payment_alt_ex_source_currency_base_rate,
:payment_alt_ex_target_currency_code,
:payment_alt_ex_target_currency_base_rate,
:payment_alt_ex_exchange_market_id,
:payment_alt_ex_calculation_rate,
:payment_alt_ex_mathematic_operator_code,
:line_extension_amount,
:tax_exclusive_amount,
:tax_inclusive_amount,
:allowance_total_amount,
:charge_total_amount,
:withholding_tax_total_amount,
:prepaid_amount,
:payable_rounding_amount,
:payable_amount,
:payable_alternative_amount,
:issue_date,
:tax_point_date,
:invoice_period_start_date,
:invoice_period_end_date,
:tax_ex_date,
:pricing_ex_date,
:payment_ex_date,
:payment_alt_ex_date,
:status_code,
:created_by_user_id,
:updated_by_user_id,
:created_at,
:updated_at);`

const selectDebitNoteHeadersSQL = `select 
id,
uuid4,
dnh_id,
note,
document_currency_code,
tax_currency_code,
pricing_currency_code,
payment_currency_code,
payment_alt_currency_code,
accounting_cost_code,
accounting_cost,
line_count_numeric,
discrepancy_response,
order_id,
billing_id,
despatch_id,
receipt_id,
statement_id,
contract_id,
accounting_supplier_party_id,
accounting_customer_party_id,
payee_party_id,
buyer_customer_party_id,
seller_supplier_party_id,
tax_representative_party_id,
tax_ex_source_currency_code,
tax_ex_source_currency_base_rate,
tax_ex_target_currency_code,
tax_ex_target_currency_base_rate,
tax_ex_exchange_market_id,
tax_ex_calculation_rate,
tax_ex_mathematic_operator_code,
pricing_ex_source_currency_code,
pricing_ex_source_currency_base_rate,
pricing_ex_target_currency_code,
pricing_ex_target_currency_base_rate,
pricing_ex_exchange_market_id,
pricing_ex_calculation_rate,
pricing_ex_mathematic_operator_code,
payment_ex_source_currency_code,
payment_ex_source_currency_base_rate,
payment_ex_target_currency_code,
payment_ex_target_currency_base_rate,
payment_ex_exchange_market_id,
payment_ex_calculation_rate,
payment_ex_mathematic_operator_code,
payment_alt_ex_source_currency_code,
payment_alt_ex_source_currency_base_rate,
payment_alt_ex_target_currency_code,
payment_alt_ex_target_currency_base_rate,
payment_alt_ex_exchange_market_id,
payment_alt_ex_calculation_rate,
payment_alt_ex_mathematic_operator_code,
line_extension_amount,
tax_exclusive_amount,
tax_inclusive_amount,
allowance_total_amount,
charge_total_amount,
withholding_tax_total_amount,
prepaid_amount,
payable_rounding_amount,
payable_amount,
payable_alternative_amount,
issue_date,
tax_point_date,
invoice_period_start_date,
invoice_period_end_date,
tax_ex_date,
pricing_ex_date,
payment_ex_date,
payment_alt_ex_date,
status_code,
created_by_user_id,
updated_by_user_id,
created_at,
updated_at from debit_note_headers`

// updateDebitNoteHeaderSQL - update DebitNoteHeaderSQL query
const updateDebitNoteHeaderSQL = `update debit_note_headers set 
note = ?, 
document_currency_code = ?, 
accounting_cost = ?, 
charge_total_amount = ?, 
prepaid_amount = ?, 
payable_rounding_amount = ?,
payable_amount = ?, updated_at = ? where uuid4 = ?;`

// CreateDebitNoteHeader - Create DebitNoteHeader
func (ds *DebitNoteHeaderService) CreateDebitNoteHeader(ctx context.Context, in *invoiceproto.CreateDebitNoteHeaderRequest) (*invoiceproto.CreateDebitNoteHeaderResponse, error) {
	user, err := partyservice.GetUserWithNewContext(ctx, in.UserId, in.UserEmail, in.RequestId, ds.UserServiceClient)
	if err != nil {
		ds.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	issueDate, err := time.Parse(common.Layout, in.IssueDate)
	if err != nil {
		ds.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	taxPointDate, err := time.Parse(common.Layout, in.TaxPointDate)
	if err != nil {
		ds.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	invoicePeriodStartDate, err := time.Parse(common.Layout, in.InvoicePeriodStartDate)
	if err != nil {
		ds.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	invoicePeriodEndDate, err := time.Parse(common.Layout, in.InvoicePeriodEndDate)
	if err != nil {
		ds.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	taxExDate, err := time.Parse(common.Layout, in.TaxExDate)
	if err != nil {
		ds.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	pricingExDate, err := time.Parse(common.Layout, in.PricingExDate)
	if err != nil {
		ds.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	paymentExDate, err := time.Parse(common.Layout, in.PaymentExDate)
	if err != nil {
		ds.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	paymentAltExDate, err := time.Parse(common.Layout, in.PaymentAltExDate)
	if err != nil {
		ds.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	ttime := common.GetTimeDetails()
	tn := common.TimeToTimestamp(ttime)

	debitNoteHeaderD := invoiceproto.DebitNoteHeaderD{}
	debitNoteHeaderD.Uuid4, err = common.GetUUIDBytes()
	if err != nil {
		ds.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	debitNoteHeaderD.DnhId = in.DnhId
	debitNoteHeaderD.Note = in.Note
	debitNoteHeaderD.DocumentCurrencyCode = in.DocumentCurrencyCode
	debitNoteHeaderD.TaxCurrencyCode = in.TaxCurrencyCode
	debitNoteHeaderD.PricingCurrencyCode = in.PricingCurrencyCode
	debitNoteHeaderD.PaymentCurrencyCode = in.PaymentCurrencyCode
	debitNoteHeaderD.PaymentAltCurrencyCode = in.PaymentAltCurrencyCode
	debitNoteHeaderD.AccountingCostCode = in.AccountingCostCode
	debitNoteHeaderD.AccountingCost = in.AccountingCost
	debitNoteHeaderD.LineCountNumeric = in.LineCountNumeric
	debitNoteHeaderD.DiscrepancyResponse = in.DiscrepancyResponse
	debitNoteHeaderD.OrderId = in.OrderId
	debitNoteHeaderD.BillingId = in.BillingId
	debitNoteHeaderD.DespatchId = in.DespatchId
	debitNoteHeaderD.ReceiptId = in.ReceiptId
	debitNoteHeaderD.StatementId = in.StatementId
	debitNoteHeaderD.ContractId = in.ContractId
	debitNoteHeaderD.AccountingSupplierPartyId = in.AccountingSupplierPartyId
	debitNoteHeaderD.AccountingCustomerPartyId = in.AccountingCustomerPartyId
	debitNoteHeaderD.PayeePartyId = in.PayeePartyId
	debitNoteHeaderD.BuyerCustomerPartyId = in.BuyerCustomerPartyId
	debitNoteHeaderD.SellerSupplierPartyId = in.SellerSupplierPartyId
	debitNoteHeaderD.TaxRepresentativePartyId = in.TaxRepresentativePartyId
	debitNoteHeaderD.TaxExSourceCurrencyCode = in.TaxExSourceCurrencyCode
	debitNoteHeaderD.TaxExSourceCurrencyBaseRate = in.TaxExSourceCurrencyBaseRate
	debitNoteHeaderD.TaxExTargetCurrencyCode = in.TaxExTargetCurrencyCode
	debitNoteHeaderD.TaxExTargetCurrencyBaseRate = in.TaxExTargetCurrencyBaseRate
	debitNoteHeaderD.TaxExExchangeMarketId = in.TaxExExchangeMarketId
	debitNoteHeaderD.TaxExCalculationRate = in.TaxExCalculationRate
	debitNoteHeaderD.TaxExMathematicOperatorCode = in.TaxExMathematicOperatorCode
	debitNoteHeaderD.PricingExSourceCurrencyCode = in.PricingExSourceCurrencyCode
	debitNoteHeaderD.PricingExSourceCurrencyBaseRate = in.PricingExSourceCurrencyBaseRate
	debitNoteHeaderD.PricingExTargetCurrencyCode = in.PricingExTargetCurrencyCode
	debitNoteHeaderD.PricingExTargetCurrencyBaseRate = in.PricingExTargetCurrencyBaseRate
	debitNoteHeaderD.PricingExExchangeMarketId = in.PricingExExchangeMarketId
	debitNoteHeaderD.PricingExCalculationRate = in.PricingExCalculationRate
	debitNoteHeaderD.PricingExMathematicOperatorCode = in.PricingExMathematicOperatorCode
	debitNoteHeaderD.PaymentExSourceCurrencyCode = in.PaymentExSourceCurrencyCode
	debitNoteHeaderD.PaymentExSourceCurrencyBaseRate = in.PaymentExSourceCurrencyBaseRate
	debitNoteHeaderD.PaymentExTargetCurrencyCode = in.PaymentExTargetCurrencyCode
	debitNoteHeaderD.PaymentExTargetCurrencyBaseRate = in.PaymentExTargetCurrencyBaseRate
	debitNoteHeaderD.PaymentExExchangeMarketId = in.PaymentExExchangeMarketId
	debitNoteHeaderD.PaymentExCalculationRate = in.PaymentExCalculationRate
	debitNoteHeaderD.PaymentExMathematicOperatorCode = in.PaymentExMathematicOperatorCode
	debitNoteHeaderD.PaymentAltExSourceCurrencyCode = in.PaymentAltExSourceCurrencyCode
	debitNoteHeaderD.PaymentAltExSourceCurrencyBaseRate = in.PaymentAltExSourceCurrencyCode
	debitNoteHeaderD.PaymentAltExTargetCurrencyCode = in.PaymentAltExTargetCurrencyCode
	debitNoteHeaderD.PaymentAltExTargetCurrencyBaseRate = in.PaymentAltExTargetCurrencyBaseRate
	debitNoteHeaderD.PaymentAltExExchangeMarketId = in.PaymentAltExExchangeMarketId
	debitNoteHeaderD.PaymentAltExCalculationRate = in.PaymentAltExCalculationRate
	debitNoteHeaderD.PaymentAltExMathematicOperatorCode = in.PaymentAltExMathematicOperatorCode
	debitNoteHeaderD.LineExtensionAmount = in.LineExtensionAmount
	debitNoteHeaderD.TaxExclusiveAmount = in.TaxExclusiveAmount
	debitNoteHeaderD.TaxInclusiveAmount = in.TaxInclusiveAmount
	debitNoteHeaderD.AllowanceTotalAmount = in.AllowanceTotalAmount
	debitNoteHeaderD.ChargeTotalAmount = in.ChargeTotalAmount
	debitNoteHeaderD.WithholdingTaxTotalAmount = in.WithholdingTaxTotalAmount
	debitNoteHeaderD.PrepaidAmount = in.PrepaidAmount
	debitNoteHeaderD.PayableRoundingAmount = in.PayableRoundingAmount
	debitNoteHeaderD.PayableAmount = in.PayableAmount
	debitNoteHeaderD.PayableAlternativeAmount = in.PayableAlternativeAmount

	debitNoteHeaderT := invoiceproto.DebitNoteHeaderT{}
	debitNoteHeaderT.IssueDate = common.TimeToTimestamp(issueDate.UTC().Truncate(time.Second))
	debitNoteHeaderT.TaxPointDate = common.TimeToTimestamp(taxPointDate.UTC().Truncate(time.Second))
	debitNoteHeaderT.InvoicePeriodStartDate = common.TimeToTimestamp(invoicePeriodStartDate.UTC().Truncate(time.Second))
	debitNoteHeaderT.InvoicePeriodEndDate = common.TimeToTimestamp(invoicePeriodEndDate.UTC().Truncate(time.Second))
	debitNoteHeaderT.TaxExDate = common.TimeToTimestamp(taxExDate.UTC().Truncate(time.Second))
	debitNoteHeaderT.PricingExDate = common.TimeToTimestamp(pricingExDate.UTC().Truncate(time.Second))
	debitNoteHeaderT.PaymentExDate = common.TimeToTimestamp(paymentExDate.UTC().Truncate(time.Second))
	debitNoteHeaderT.PaymentAltExDate = common.TimeToTimestamp(paymentAltExDate.UTC().Truncate(time.Second))

	crUpdUser := commonproto.CrUpdUser{}
	crUpdUser.StatusCode = "active"
	crUpdUser.CreatedByUserId = user.Id
	crUpdUser.UpdatedByUserId = user.Id

	crUpdTime := commonproto.CrUpdTime{}
	crUpdTime.CreatedAt = tn
	crUpdTime.UpdatedAt = tn

	debitNoteHeader := invoiceproto.DebitNoteHeader{DebitNoteHeaderD: &debitNoteHeaderD, DebitNoteHeaderT: &debitNoteHeaderT, CrUpdUser: &crUpdUser, CrUpdTime: &crUpdTime}

	debitNoteLines := []*invoiceproto.DebitNoteLine{}
	// we will do for loop on lines which is comes from client form
	for _, line := range in.DebitNoteLines {
		line.UserId = in.UserId
		line.UserEmail = in.UserEmail
		line.RequestId = in.RequestId
		debitNoteLine, err := ds.ProcessDebitNoteLineRequest(ctx, line)
		if err != nil {
			ds.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
			return nil, err
		}

		debitNoteLines = append(debitNoteLines, debitNoteLine)
	}

	err = ds.insertDebitNoteHeader(ctx, insertDebitNoteHeaderSQL, &debitNoteHeader, insertDebitNoteLineSQL, debitNoteLines, in.GetUserEmail(), in.GetRequestId())
	if err != nil {
		ds.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	debitNoteHeaderResponse := invoiceproto.CreateDebitNoteHeaderResponse{}
	debitNoteHeaderResponse.DebitNoteHeader = &debitNoteHeader
	return &debitNoteHeaderResponse, nil
}

func (ds *DebitNoteHeaderService) insertDebitNoteHeader(ctx context.Context, insertDebitNoteHeaderSQL string, debitNoteHeader *invoiceproto.DebitNoteHeader, insertDebitNoteLineSQL string, debitNoteLines []*invoiceproto.DebitNoteLine, userEmail string, requestID string) error {
	debitNoteHeaderTmp, err := ds.crDebitNoteHeaderStruct(ctx, debitNoteHeader, userEmail, requestID)
	if err != nil {
		ds.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
		return err
	}
	err = ds.DBService.InsUpd(ctx, userEmail, requestID, func(tx *sqlx.Tx) error {
		res, err := tx.NamedExecContext(ctx, insertDebitNoteHeaderSQL, debitNoteHeaderTmp)
		if err != nil {
			ds.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}

		uID, err := res.LastInsertId()
		if err != nil {
			ds.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}
		debitNoteHeader.DebitNoteHeaderD.Id = uint32(uID)
		uuid4Str, err := common.UUIDBytesToStr(debitNoteHeader.DebitNoteHeaderD.Uuid4)
		if err != nil {
			ds.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}
		debitNoteHeader.DebitNoteHeaderD.IdS = uuid4Str

		for _, debitNoteLine := range debitNoteLines {
			debitNoteLine.DebitNoteLineD.DebitNoteHeaderId = debitNoteHeader.DebitNoteHeaderD.Id
			debitNoteLineTmp, err := ds.crDebitNoteLineStruct(ctx, debitNoteLine, userEmail, requestID)
			if err != nil {
				ds.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
				return err
			}
			_, err = tx.NamedExecContext(ctx, insertDebitNoteLineSQL, debitNoteLineTmp)
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

// crDebitNoteHeaderStruct - process DebitNoteHeader details
func (ds *DebitNoteHeaderService) crDebitNoteHeaderStruct(ctx context.Context, debitNoteHeader *invoiceproto.DebitNoteHeader, userEmail string, requestID string) (*invoicestruct.DebitNoteHeader, error) {
	crUpdTime := new(commonstruct.CrUpdTime)
	crUpdTime.CreatedAt = common.TimestampToTime(debitNoteHeader.CrUpdTime.CreatedAt)
	crUpdTime.UpdatedAt = common.TimestampToTime(debitNoteHeader.CrUpdTime.UpdatedAt)

	debitNoteHeaderT := new(invoicestruct.DebitNoteHeaderT)
	debitNoteHeaderT.IssueDate = common.TimestampToTime(debitNoteHeader.DebitNoteHeaderT.IssueDate)
	debitNoteHeaderT.TaxPointDate = common.TimestampToTime(debitNoteHeader.DebitNoteHeaderT.TaxPointDate)
	debitNoteHeaderT.InvoicePeriodStartDate = common.TimestampToTime(debitNoteHeader.DebitNoteHeaderT.InvoicePeriodStartDate)
	debitNoteHeaderT.InvoicePeriodEndDate = common.TimestampToTime(debitNoteHeader.DebitNoteHeaderT.InvoicePeriodEndDate)
	debitNoteHeaderT.TaxExDate = common.TimestampToTime(debitNoteHeader.DebitNoteHeaderT.TaxExDate)
	debitNoteHeaderT.PricingExDate = common.TimestampToTime(debitNoteHeader.DebitNoteHeaderT.PricingExDate)
	debitNoteHeaderT.PaymentExDate = common.TimestampToTime(debitNoteHeader.DebitNoteHeaderT.PaymentExDate)
	debitNoteHeaderT.PaymentAltExDate = common.TimestampToTime(debitNoteHeader.DebitNoteHeaderT.PaymentAltExDate)

	debitNoteHeaderTmp := invoicestruct.DebitNoteHeader{DebitNoteHeaderD: debitNoteHeader.DebitNoteHeaderD, DebitNoteHeaderT: debitNoteHeaderT, CrUpdUser: debitNoteHeader.CrUpdUser, CrUpdTime: crUpdTime}
	return &debitNoteHeaderTmp, nil
}

// GetDebitNoteHeaders - Get DebitNoteHeaders
func (ds *DebitNoteHeaderService) GetDebitNoteHeaders(ctx context.Context, in *invoiceproto.GetDebitNoteHeadersRequest) (*invoiceproto.GetDebitNoteHeadersResponse, error) {
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

	debitNoteHeaders := []*invoiceproto.DebitNoteHeader{}

	nselectDebitNoteHeadersSQL := selectDebitNoteHeadersSQL + ` where ` + query

	rows, err := ds.DBService.DB.QueryxContext(ctx, nselectDebitNoteHeadersSQL, "active")
	if err != nil {
		ds.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	for rows.Next() {

		debitNoteHeaderTmp := invoicestruct.DebitNoteHeader{}
		err = rows.StructScan(&debitNoteHeaderTmp)
		if err != nil {
			ds.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
			return nil, err
		}

		getRequest := commonproto.GetRequest{}
		getRequest.UserEmail = in.UserEmail
		getRequest.RequestId = in.RequestId
		debitNoteHeader, err := ds.getDebitNoteHeaderStruct(ctx, &getRequest, debitNoteHeaderTmp)
		if err != nil {
			ds.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
			return nil, err
		}
		debitNoteHeaders = append(debitNoteHeaders, debitNoteHeader)

	}

	debitNoteHeadersResponse := invoiceproto.GetDebitNoteHeadersResponse{}
	if len(debitNoteHeaders) != 0 {
		next := debitNoteHeaders[len(debitNoteHeaders)-1].DebitNoteHeaderD.Id
		next--
		nextc := common.EncodeCursor(next)
		debitNoteHeadersResponse = invoiceproto.GetDebitNoteHeadersResponse{DebitNoteHeaders: debitNoteHeaders, NextCursor: nextc}
	} else {
		debitNoteHeadersResponse = invoiceproto.GetDebitNoteHeadersResponse{DebitNoteHeaders: debitNoteHeaders, NextCursor: "0"}
	}
	return &debitNoteHeadersResponse, nil
}

// GetDebitNoteHeader - Get DebitNoteHeader
func (ds *DebitNoteHeaderService) GetDebitNoteHeader(ctx context.Context, inReq *invoiceproto.GetDebitNoteHeaderRequest) (*invoiceproto.GetDebitNoteHeaderResponse, error) {
	in := inReq.GetRequest

	uuid4byte, err := common.UUIDStrToBytes(in.Id)
	if err != nil {
		ds.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	nselectDebitNoteHeadersSQL := selectDebitNoteHeadersSQL + ` where uuid4 = ? and status_code = ?;`
	row := ds.DBService.DB.QueryRowxContext(ctx, nselectDebitNoteHeadersSQL, uuid4byte, "active")
	debitNoteHeaderTmp := invoicestruct.DebitNoteHeader{}
	err = row.StructScan(&debitNoteHeaderTmp)
	if err != nil {
		ds.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	debitNoteHeader, err := ds.getDebitNoteHeaderStruct(ctx, in, debitNoteHeaderTmp)
	if err != nil {
		ds.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	debitNoteHeaderResponse := invoiceproto.GetDebitNoteHeaderResponse{}
	debitNoteHeaderResponse.DebitNoteHeader = debitNoteHeader
	return &debitNoteHeaderResponse, nil
}

// GetDebitNoteHeaderByPk - Get DebitNoteHeader By Primary key(Id)
func (ds *DebitNoteHeaderService) GetDebitNoteHeaderByPk(ctx context.Context, inReq *invoiceproto.GetDebitNoteHeaderByPkRequest) (*invoiceproto.GetDebitNoteHeaderByPkResponse, error) {
	in := inReq.GetByIdRequest

	nselectDebitNoteHeadersSQL := selectDebitNoteHeadersSQL + ` where id = ? and status_code = ?;`
	row := ds.DBService.DB.QueryRowxContext(ctx, nselectDebitNoteHeadersSQL, in.Id, "active")
	debitNoteHeaderTmp := invoicestruct.DebitNoteHeader{}
	err := row.StructScan(&debitNoteHeaderTmp)
	if err != nil {
		ds.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	getRequest := commonproto.GetRequest{}
	getRequest.UserEmail = in.UserEmail
	getRequest.RequestId = in.RequestId
	debitNoteHeader, err := ds.getDebitNoteHeaderStruct(ctx, &getRequest, debitNoteHeaderTmp)
	if err != nil {
		ds.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	debitNoteHeaderResponse := invoiceproto.GetDebitNoteHeaderByPkResponse{}
	debitNoteHeaderResponse.DebitNoteHeader = debitNoteHeader
	return &debitNoteHeaderResponse, nil
}

// getDebitNoteHeaderStruct - Get debit note header
func (ds *DebitNoteHeaderService) getDebitNoteHeaderStruct(ctx context.Context, in *commonproto.GetRequest, debitNoteHeaderTmp invoicestruct.DebitNoteHeader) (*invoiceproto.DebitNoteHeader, error) {
	debitNoteHeaderT := new(invoiceproto.DebitNoteHeaderT)
	debitNoteHeaderT.IssueDate = common.TimeToTimestamp(debitNoteHeaderTmp.DebitNoteHeaderT.IssueDate)
	debitNoteHeaderT.TaxPointDate = common.TimeToTimestamp(debitNoteHeaderTmp.DebitNoteHeaderT.TaxPointDate)
	debitNoteHeaderT.TaxExDate = common.TimeToTimestamp(debitNoteHeaderTmp.DebitNoteHeaderT.TaxExDate)
	debitNoteHeaderT.InvoicePeriodStartDate = common.TimeToTimestamp(debitNoteHeaderTmp.DebitNoteHeaderT.InvoicePeriodStartDate)
	debitNoteHeaderT.InvoicePeriodEndDate = common.TimeToTimestamp(debitNoteHeaderTmp.DebitNoteHeaderT.InvoicePeriodEndDate)
	debitNoteHeaderT.PricingExDate = common.TimeToTimestamp(debitNoteHeaderTmp.DebitNoteHeaderT.PricingExDate)
	debitNoteHeaderT.PaymentExDate = common.TimeToTimestamp(debitNoteHeaderTmp.DebitNoteHeaderT.PaymentExDate)
	debitNoteHeaderT.PaymentAltExDate = common.TimeToTimestamp(debitNoteHeaderTmp.DebitNoteHeaderT.PaymentAltExDate)

	uuid4Str, err := common.UUIDBytesToStr(debitNoteHeaderTmp.DebitNoteHeaderD.Uuid4)
	if err != nil {
		ds.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	debitNoteHeaderTmp.DebitNoteHeaderD.IdS = uuid4Str

	crUpdTime := new(commonproto.CrUpdTime)
	crUpdTime.CreatedAt = common.TimeToTimestamp(debitNoteHeaderTmp.CrUpdTime.CreatedAt)
	crUpdTime.UpdatedAt = common.TimeToTimestamp(debitNoteHeaderTmp.CrUpdTime.UpdatedAt)

	debitNoteHeader := invoiceproto.DebitNoteHeader{DebitNoteHeaderD: debitNoteHeaderTmp.DebitNoteHeaderD, DebitNoteHeaderT: debitNoteHeaderT, CrUpdUser: debitNoteHeaderTmp.CrUpdUser, CrUpdTime: crUpdTime}
	return &debitNoteHeader, nil
}

// UpdateDebitNoteHeader - Update DebitNoteHeader
func (ds *DebitNoteHeaderService) UpdateDebitNoteHeader(ctx context.Context, in *invoiceproto.UpdateDebitNoteHeaderRequest) (*invoiceproto.UpdateDebitNoteHeaderResponse, error) {
	db := ds.DBService.DB
	tn := common.GetTimeDetails()

	uuid4byte, err := common.UUIDStrToBytes(in.Id)
	if err != nil {
		ds.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	stmt, err := db.PreparexContext(ctx, updateDebitNoteHeaderSQL)
	if err != nil {
		ds.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	err = ds.DBService.InsUpd(ctx, in.GetUserEmail(), in.GetRequestId(), func(tx *sqlx.Tx) error {
		_, err = tx.StmtxContext(ctx, stmt).ExecContext(ctx,
			in.Note,
			in.DocumentCurrencyCode,
			in.AccountingCost,
			in.ChargeTotalAmount,
			in.PrepaidAmount,
			in.PayableRoundingAmount,
			in.PayableAmount,
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

	return &invoiceproto.UpdateDebitNoteHeaderResponse{}, nil
}
