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

// CreditNoteHeaderService - For accessing CreditNoteHeader services
type CreditNoteHeaderService struct {
	log               *zap.Logger
	DBService         *common.DBService
	RedisService      *common.RedisService
	UserServiceClient partyproto.UserServiceClient
	invoiceproto.UnimplementedCreditNoteHeaderServiceServer
}

// NewCreditNoteHeaderService - Create CreditNoteHeader service
func NewCreditNoteHeaderService(log *zap.Logger, dbOpt *common.DBService, redisOpt *common.RedisService, userServiceClient partyproto.UserServiceClient) *CreditNoteHeaderService {
	return &CreditNoteHeaderService{
		log:               log,
		DBService:         dbOpt,
		RedisService:      redisOpt,
		UserServiceClient: userServiceClient,
	}
}

const insertCreditNoteHeaderSQL = `insert into credit_note_headers
	  (uuid4,
cnh_id,
credit_note_type_code,
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
contract_id,
statement_id,
signature,
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
due_date,
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
:cnh_id,
:credit_note_type_code,
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
:contract_id,
:statement_id,
:signature,
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
:due_date,
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

const selectCreditNoteHeadersSQL = `select 
  id,
uuid4,
cnh_id,
credit_note_type_code,
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
contract_id,
statement_id,
signature,
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
due_date,
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
updated_at from credit_note_headers`

// updateCreditNoteHeaderSQL - update CreditNoteHeaderSQL query
const updateCreditNoteHeaderSQL = `update credit_note_headers set 
note = ?, 
tax_currency_code = ?, 
charge_total_amount = ?, 
prepaid_amount = ?, 
payable_rounding_amount = ?,
payable_amount = ?, updated_at = ? where uuid4 = ?;`

// CreateCreditNoteHeader - Create CreditNoteHeader
func (cs *CreditNoteHeaderService) CreateCreditNoteHeader(ctx context.Context, in *invoiceproto.CreateCreditNoteHeaderRequest) (*invoiceproto.CreateCreditNoteHeaderResponse, error) {
	user, err := partyservice.GetUserWithNewContext(ctx, in.UserId, in.UserEmail, in.RequestId, cs.UserServiceClient)
	if err != nil {
		cs.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	issueDate, err := time.Parse(common.Layout, in.IssueDate)
	if err != nil {
		cs.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	dueDate, err := time.Parse(common.Layout, in.DueDate)
	if err != nil {
		cs.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

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

	taxExDate, err := time.Parse(common.Layout, in.TaxExDate)
	if err != nil {
		cs.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	pricingExDate, err := time.Parse(common.Layout, in.PricingExDate)
	if err != nil {
		cs.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	paymentExDate, err := time.Parse(common.Layout, in.PaymentExDate)
	if err != nil {
		cs.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	paymentAltExDate, err := time.Parse(common.Layout, in.PaymentAltExDate)
	if err != nil {
		cs.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	ttime := common.GetTimeDetails()
	tn := common.TimeToTimestamp(ttime)

	creditNoteHeaderD := invoiceproto.CreditNoteHeaderD{}
	creditNoteHeaderD.Uuid4, err = common.GetUUIDBytes()
	if err != nil {
		cs.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	creditNoteHeaderD.CnhId = in.CnhId
	creditNoteHeaderD.CreditNoteTypeCode = in.CreditNoteTypeCode
	creditNoteHeaderD.Note = in.Note
	creditNoteHeaderD.DocumentCurrencyCode = in.DocumentCurrencyCode
	creditNoteHeaderD.TaxCurrencyCode = in.TaxCurrencyCode
	creditNoteHeaderD.PricingCurrencyCode = in.PricingCurrencyCode
	creditNoteHeaderD.PaymentCurrencyCode = in.PaymentCurrencyCode
	creditNoteHeaderD.PaymentAltCurrencyCode = in.PaymentAltCurrencyCode
	creditNoteHeaderD.AccountingCostCode = in.AccountingCostCode
	creditNoteHeaderD.AccountingCost = in.AccountingCost
	creditNoteHeaderD.LineCountNumeric = in.LineCountNumeric
	creditNoteHeaderD.DiscrepancyResponse = in.DiscrepancyResponse
	creditNoteHeaderD.OrderId = in.OrderId
	creditNoteHeaderD.BillingId = in.BillingId
	creditNoteHeaderD.DespatchId = in.DespatchId
	creditNoteHeaderD.ReceiptId = in.ReceiptId
	creditNoteHeaderD.ContractId = in.ContractId
	creditNoteHeaderD.StatementId = in.StatementId
	creditNoteHeaderD.Signature = in.Signature
	creditNoteHeaderD.AccountingSupplierPartyId = in.AccountingSupplierPartyId
	creditNoteHeaderD.AccountingCustomerPartyId = in.AccountingCustomerPartyId
	creditNoteHeaderD.PayeePartyId = in.PayeePartyId
	creditNoteHeaderD.BuyerCustomerPartyId = in.BuyerCustomerPartyId
	creditNoteHeaderD.SellerSupplierPartyId = in.SellerSupplierPartyId
	creditNoteHeaderD.TaxRepresentativePartyId = in.TaxRepresentativePartyId
	creditNoteHeaderD.TaxExSourceCurrencyCode = in.TaxExSourceCurrencyCode
	creditNoteHeaderD.TaxExSourceCurrencyBaseRate = in.TaxExSourceCurrencyBaseRate
	creditNoteHeaderD.TaxExTargetCurrencyCode = in.TaxExTargetCurrencyCode
	creditNoteHeaderD.TaxExTargetCurrencyBaseRate = in.TaxExTargetCurrencyBaseRate
	creditNoteHeaderD.TaxExExchangeMarketId = in.TaxExExchangeMarketId
	creditNoteHeaderD.TaxExCalculationRate = in.TaxExCalculationRate
	creditNoteHeaderD.TaxExMathematicOperatorCode = in.TaxExMathematicOperatorCode
	creditNoteHeaderD.PricingExSourceCurrencyCode = in.PricingExSourceCurrencyCode
	creditNoteHeaderD.PricingExSourceCurrencyBaseRate = in.PricingExSourceCurrencyBaseRate
	creditNoteHeaderD.PricingExTargetCurrencyCode = in.PricingExTargetCurrencyCode
	creditNoteHeaderD.PricingExTargetCurrencyBaseRate = in.PricingExTargetCurrencyBaseRate
	creditNoteHeaderD.PricingExExchangeMarketId = in.PricingExExchangeMarketId
	creditNoteHeaderD.PricingExCalculationRate = in.PricingExCalculationRate
	creditNoteHeaderD.PricingExMathematicOperatorCode = in.PricingExMathematicOperatorCode
	creditNoteHeaderD.PaymentExSourceCurrencyCode = in.PaymentExSourceCurrencyCode
	creditNoteHeaderD.PaymentExSourceCurrencyBaseRate = in.PaymentExSourceCurrencyBaseRate
	creditNoteHeaderD.PaymentExTargetCurrencyCode = in.PaymentExTargetCurrencyCode
	creditNoteHeaderD.PaymentExTargetCurrencyBaseRate = in.PaymentExTargetCurrencyBaseRate
	creditNoteHeaderD.PaymentExExchangeMarketId = in.PaymentExExchangeMarketId
	creditNoteHeaderD.PaymentExCalculationRate = in.PaymentExCalculationRate
	creditNoteHeaderD.PaymentExMathematicOperatorCode = in.PaymentExMathematicOperatorCode
	creditNoteHeaderD.PaymentAltExSourceCurrencyCode = in.PaymentAltExSourceCurrencyCode
	creditNoteHeaderD.PaymentAltExSourceCurrencyBaseRate = in.PaymentAltExSourceCurrencyCode
	creditNoteHeaderD.PaymentAltExTargetCurrencyCode = in.PaymentAltExTargetCurrencyCode
	creditNoteHeaderD.PaymentAltExTargetCurrencyBaseRate = in.PaymentAltExTargetCurrencyBaseRate
	creditNoteHeaderD.PaymentAltExExchangeMarketId = in.PaymentAltExExchangeMarketId
	creditNoteHeaderD.PaymentAltExCalculationRate = in.PaymentAltExCalculationRate
	creditNoteHeaderD.PaymentAltExMathematicOperatorCode = in.PaymentAltExMathematicOperatorCode
	creditNoteHeaderD.LineExtensionAmount = in.LineExtensionAmount
	creditNoteHeaderD.TaxExclusiveAmount = in.TaxExclusiveAmount
	creditNoteHeaderD.TaxInclusiveAmount = in.TaxInclusiveAmount
	creditNoteHeaderD.AllowanceTotalAmount = in.AllowanceTotalAmount
	creditNoteHeaderD.ChargeTotalAmount = in.ChargeTotalAmount
	creditNoteHeaderD.WithholdingTaxTotalAmount = in.WithholdingTaxTotalAmount
	creditNoteHeaderD.PrepaidAmount = in.PrepaidAmount
	creditNoteHeaderD.PayableRoundingAmount = in.PayableRoundingAmount
	creditNoteHeaderD.PayableAmount = in.PayableAmount
	creditNoteHeaderD.PayableAlternativeAmount = in.PayableAlternativeAmount

	creditNoteHeaderT := invoiceproto.CreditNoteHeaderT{}
	creditNoteHeaderT.IssueDate = common.TimeToTimestamp(issueDate.UTC().Truncate(time.Second))
	creditNoteHeaderT.DueDate = common.TimeToTimestamp(dueDate.UTC().Truncate(time.Second))
	creditNoteHeaderT.TaxPointDate = common.TimeToTimestamp(taxPointDate.UTC().Truncate(time.Second))
	creditNoteHeaderT.InvoicePeriodStartDate = common.TimeToTimestamp(invoicePeriodStartDate.UTC().Truncate(time.Second))
	creditNoteHeaderT.InvoicePeriodEndDate = common.TimeToTimestamp(invoicePeriodEndDate.UTC().Truncate(time.Second))
	creditNoteHeaderT.TaxExDate = common.TimeToTimestamp(taxExDate.UTC().Truncate(time.Second))
	creditNoteHeaderT.PricingExDate = common.TimeToTimestamp(pricingExDate.UTC().Truncate(time.Second))
	creditNoteHeaderT.PaymentExDate = common.TimeToTimestamp(paymentExDate.UTC().Truncate(time.Second))
	creditNoteHeaderT.PaymentAltExDate = common.TimeToTimestamp(paymentAltExDate.UTC().Truncate(time.Second))

	crUpdUser := commonproto.CrUpdUser{}
	crUpdUser.StatusCode = "active"
	crUpdUser.CreatedByUserId = user.Id
	crUpdUser.UpdatedByUserId = user.Id

	crUpdTime := commonproto.CrUpdTime{}
	crUpdTime.CreatedAt = tn
	crUpdTime.UpdatedAt = tn

	creditNoteHeader := invoiceproto.CreditNoteHeader{CreditNoteHeaderD: &creditNoteHeaderD, CreditNoteHeaderT: &creditNoteHeaderT, CrUpdUser: &crUpdUser, CrUpdTime: &crUpdTime}

	creditNoteLines := []*invoiceproto.CreditNoteLine{}
	// we will do for loop on lines which is comes from client form
	for _, line := range in.CreditNoteLines {
		line.UserId = in.UserId
		line.UserEmail = in.UserEmail
		line.RequestId = in.RequestId
		creditNoteLine, err := cs.ProcessCreditNoteLineRequest(ctx, line)
		if err != nil {
			cs.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
			return nil, err
		}

		creditNoteLines = append(creditNoteLines, creditNoteLine)
	}

	err = cs.insertCreditNoteHeader(ctx, insertCreditNoteHeaderSQL, &creditNoteHeader, insertCreditNoteLineSQL, creditNoteLines, in.GetUserEmail(), in.GetRequestId())
	if err != nil {
		cs.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	creditNoteHeaderResponse := invoiceproto.CreateCreditNoteHeaderResponse{}
	creditNoteHeaderResponse.CreditNoteHeader = &creditNoteHeader
	return &creditNoteHeaderResponse, nil
}

func (cs *CreditNoteHeaderService) insertCreditNoteHeader(ctx context.Context, insertCreditNoteHeaderSQL string, creditNoteHeader *invoiceproto.CreditNoteHeader, insertCreditNoteLineSQL string, creditNoteLines []*invoiceproto.CreditNoteLine, userEmail string, requestID string) error {
	creditNoteHeaderTmp, err := cs.crCreditNoteHeaderStruct(ctx, creditNoteHeader, userEmail, requestID)
	if err != nil {
		cs.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
		return err
	}
	err = cs.DBService.InsUpd(ctx, userEmail, requestID, func(tx *sqlx.Tx) error {
		res, err := tx.NamedExecContext(ctx, insertCreditNoteHeaderSQL, creditNoteHeaderTmp)
		if err != nil {
			cs.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}

		uID, err := res.LastInsertId()
		if err != nil {
			cs.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}
		creditNoteHeader.CreditNoteHeaderD.Id = uint32(uID)
		uuid4Str, err := common.UUIDBytesToStr(creditNoteHeader.CreditNoteHeaderD.Uuid4)
		if err != nil {
			cs.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}
		creditNoteHeader.CreditNoteHeaderD.IdS = uuid4Str

		for _, creditNoteLine := range creditNoteLines {
			creditNoteLine.CreditNoteLineD.CreditNoteHeaderId = creditNoteHeader.CreditNoteHeaderD.Id
			creditNoteLineTmp, err := cs.crCreditNoteLineStruct(ctx, creditNoteLine, userEmail, requestID)
			if err != nil {
				cs.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
				return err
			}
			_, err = tx.NamedExecContext(ctx, insertCreditNoteLineSQL, creditNoteLineTmp)
			if err != nil {
				cs.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
				return err
			}

		}

		return nil
	})

	if err != nil {
		cs.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
		return err
	}
	return nil
}

// crCreditNoteHeaderStruct - process CreditNoteHeader details
func (cs *CreditNoteHeaderService) crCreditNoteHeaderStruct(ctx context.Context, creditNoteHeader *invoiceproto.CreditNoteHeader, userEmail string, requestID string) (*invoicestruct.CreditNoteHeader, error) {
	crUpdTime := new(commonstruct.CrUpdTime)
	crUpdTime.CreatedAt = common.TimestampToTime(creditNoteHeader.CrUpdTime.CreatedAt)
	crUpdTime.UpdatedAt = common.TimestampToTime(creditNoteHeader.CrUpdTime.UpdatedAt)

	creditNoteHeaderT := new(invoicestruct.CreditNoteHeaderT)
	creditNoteHeaderT.IssueDate = common.TimestampToTime(creditNoteHeader.CreditNoteHeaderT.IssueDate)
	creditNoteHeaderT.DueDate = common.TimestampToTime(creditNoteHeader.CreditNoteHeaderT.DueDate)
	creditNoteHeaderT.TaxPointDate = common.TimestampToTime(creditNoteHeader.CreditNoteHeaderT.TaxPointDate)
	creditNoteHeaderT.InvoicePeriodStartDate = common.TimestampToTime(creditNoteHeader.CreditNoteHeaderT.InvoicePeriodStartDate)
	creditNoteHeaderT.InvoicePeriodEndDate = common.TimestampToTime(creditNoteHeader.CreditNoteHeaderT.InvoicePeriodEndDate)
	creditNoteHeaderT.TaxExDate = common.TimestampToTime(creditNoteHeader.CreditNoteHeaderT.TaxExDate)
	creditNoteHeaderT.PricingExDate = common.TimestampToTime(creditNoteHeader.CreditNoteHeaderT.PricingExDate)
	creditNoteHeaderT.PaymentExDate = common.TimestampToTime(creditNoteHeader.CreditNoteHeaderT.PaymentExDate)
	creditNoteHeaderT.PaymentAltExDate = common.TimestampToTime(creditNoteHeader.CreditNoteHeaderT.PaymentAltExDate)

	creditNoteHeaderTmp := invoicestruct.CreditNoteHeader{CreditNoteHeaderD: creditNoteHeader.CreditNoteHeaderD, CreditNoteHeaderT: creditNoteHeaderT, CrUpdUser: creditNoteHeader.CrUpdUser, CrUpdTime: crUpdTime}

	return &creditNoteHeaderTmp, nil
}

// GetCreditNoteHeaders - Get CreditNoteHeaders
func (cs *CreditNoteHeaderService) GetCreditNoteHeaders(ctx context.Context, in *invoiceproto.GetCreditNoteHeadersRequest) (*invoiceproto.GetCreditNoteHeadersResponse, error) {
	limit := in.GetLimit()
	nextCursor := in.GetNextCursor()
	if limit == "" {
		limit = cs.DBService.LimitSQLRows
	}
	query := "status_code = ?"
	if nextCursor == "" {
		query = query + " order by id desc " + " limit " + limit + ";"
	} else {
		nextCursor = common.DecodeCursor(nextCursor)
		query = query + " " + "and" + " " + "id <= " + nextCursor + " order by id desc " + " limit " + limit + ";"
	}

	creditNoteHeaders := []*invoiceproto.CreditNoteHeader{}

	nselectCreditNoteHeadersSQL := selectCreditNoteHeadersSQL + ` where ` + query

	rows, err := cs.DBService.DB.QueryxContext(ctx, nselectCreditNoteHeadersSQL, "active")
	if err != nil {
		cs.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	for rows.Next() {

		creditNoteHeaderTmp := invoicestruct.CreditNoteHeader{}
		err = rows.StructScan(&creditNoteHeaderTmp)
		if err != nil {
			cs.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
			return nil, err
		}

		getRequest := commonproto.GetRequest{}
		getRequest.UserEmail = in.UserEmail
		getRequest.RequestId = in.RequestId
		creditNoteHeader, err := cs.getCreditNoteHeaderStruct(ctx, &getRequest, creditNoteHeaderTmp)
		if err != nil {
			cs.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
			return nil, err
		}
		creditNoteHeaders = append(creditNoteHeaders, creditNoteHeader)

	}

	creditNoteHeadersResponse := invoiceproto.GetCreditNoteHeadersResponse{}
	if len(creditNoteHeaders) != 0 {
		next := creditNoteHeaders[len(creditNoteHeaders)-1].CreditNoteHeaderD.Id
		next--
		nextc := common.EncodeCursor(next)
		creditNoteHeadersResponse = invoiceproto.GetCreditNoteHeadersResponse{CreditNoteHeaders: creditNoteHeaders, NextCursor: nextc}
	} else {
		creditNoteHeadersResponse = invoiceproto.GetCreditNoteHeadersResponse{CreditNoteHeaders: creditNoteHeaders, NextCursor: "0"}
	}
	return &creditNoteHeadersResponse, nil
}

// GetCreditNoteHeader - Get CreditNoteHeader
func (cs *CreditNoteHeaderService) GetCreditNoteHeader(ctx context.Context, inReq *invoiceproto.GetCreditNoteHeaderRequest) (*invoiceproto.GetCreditNoteHeaderResponse, error) {
	in := inReq.GetRequest
	uuid4byte, err := common.UUIDStrToBytes(in.Id)
	if err != nil {
		cs.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	nselectCreditNoteHeadersSQL := selectCreditNoteHeadersSQL + ` where uuid4 = ? and status_code = ?;`
	row := cs.DBService.DB.QueryRowxContext(ctx, nselectCreditNoteHeadersSQL, uuid4byte, "active")
	creditNoteHeaderTmp := invoicestruct.CreditNoteHeader{}
	err = row.StructScan(&creditNoteHeaderTmp)
	if err != nil {
		cs.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	creditNoteHeader, err := cs.getCreditNoteHeaderStruct(ctx, in, creditNoteHeaderTmp)
	if err != nil {
		cs.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	creditNoteHeaderResponse := invoiceproto.GetCreditNoteHeaderResponse{}
	creditNoteHeaderResponse.CreditNoteHeader = creditNoteHeader
	return &creditNoteHeaderResponse, nil
}

// GetCreditNoteHeaderByPk - Get CreditNoteHeader By Primary key(Id)
func (cs *CreditNoteHeaderService) GetCreditNoteHeaderByPk(ctx context.Context, inReq *invoiceproto.GetCreditNoteHeaderByPkRequest) (*invoiceproto.GetCreditNoteHeaderByPkResponse, error) {
	in := inReq.GetByIdRequest

	nselectCreditNoteHeadersSQL := selectCreditNoteHeadersSQL + ` where id = ? and status_code = ?;`
	row := cs.DBService.DB.QueryRowxContext(ctx, nselectCreditNoteHeadersSQL, in.Id, "active")
	creditNoteHeaderTmp := invoicestruct.CreditNoteHeader{}
	err := row.StructScan(&creditNoteHeaderTmp)
	if err != nil {
		cs.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	getRequest := commonproto.GetRequest{}
	getRequest.UserEmail = in.UserEmail
	getRequest.RequestId = in.RequestId
	creditNoteHeader, err := cs.getCreditNoteHeaderStruct(ctx, &getRequest, creditNoteHeaderTmp)
	if err != nil {
		cs.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	creditNoteHeaderResponse := invoiceproto.GetCreditNoteHeaderByPkResponse{}
	creditNoteHeaderResponse.CreditNoteHeader = creditNoteHeader
	return &creditNoteHeaderResponse, nil
}

// getCreditNoteHeaderStruct - Get credit note header
func (cs *CreditNoteHeaderService) getCreditNoteHeaderStruct(ctx context.Context, in *commonproto.GetRequest, creditNoteHeaderTmp invoicestruct.CreditNoteHeader) (*invoiceproto.CreditNoteHeader, error) {
	creditNoteHeaderT := new(invoiceproto.CreditNoteHeaderT)
	creditNoteHeaderT.IssueDate = common.TimeToTimestamp(creditNoteHeaderTmp.CreditNoteHeaderT.IssueDate)
	creditNoteHeaderT.DueDate = common.TimeToTimestamp(creditNoteHeaderTmp.CreditNoteHeaderT.DueDate)
	creditNoteHeaderT.TaxPointDate = common.TimeToTimestamp(creditNoteHeaderTmp.CreditNoteHeaderT.TaxPointDate)
	creditNoteHeaderT.TaxExDate = common.TimeToTimestamp(creditNoteHeaderTmp.CreditNoteHeaderT.TaxExDate)
	creditNoteHeaderT.InvoicePeriodStartDate = common.TimeToTimestamp(creditNoteHeaderTmp.CreditNoteHeaderT.InvoicePeriodStartDate)
	creditNoteHeaderT.InvoicePeriodEndDate = common.TimeToTimestamp(creditNoteHeaderTmp.CreditNoteHeaderT.InvoicePeriodEndDate)
	creditNoteHeaderT.PricingExDate = common.TimeToTimestamp(creditNoteHeaderTmp.CreditNoteHeaderT.PricingExDate)
	creditNoteHeaderT.PaymentExDate = common.TimeToTimestamp(creditNoteHeaderTmp.CreditNoteHeaderT.PaymentExDate)
	creditNoteHeaderT.PaymentAltExDate = common.TimeToTimestamp(creditNoteHeaderTmp.CreditNoteHeaderT.PaymentAltExDate)

	uuid4Str, err := common.UUIDBytesToStr(creditNoteHeaderTmp.CreditNoteHeaderD.Uuid4)
	if err != nil {
		cs.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	creditNoteHeaderTmp.CreditNoteHeaderD.IdS = uuid4Str

	crUpdTime := new(commonproto.CrUpdTime)
	crUpdTime.CreatedAt = common.TimeToTimestamp(creditNoteHeaderTmp.CrUpdTime.CreatedAt)
	crUpdTime.UpdatedAt = common.TimeToTimestamp(creditNoteHeaderTmp.CrUpdTime.UpdatedAt)

	creditNoteHeader := invoiceproto.CreditNoteHeader{CreditNoteHeaderD: creditNoteHeaderTmp.CreditNoteHeaderD, CreditNoteHeaderT: creditNoteHeaderT, CrUpdUser: creditNoteHeaderTmp.CrUpdUser, CrUpdTime: crUpdTime}

	return &creditNoteHeader, nil
}

// UpdateCreditNoteHeader - Update CreditNoteHeader
func (cs *CreditNoteHeaderService) UpdateCreditNoteHeader(ctx context.Context, in *invoiceproto.UpdateCreditNoteHeaderRequest) (*invoiceproto.UpdateCreditNoteHeaderResponse, error) {
	db := cs.DBService.DB
	tn := common.GetTimeDetails()

	uuid4byte, err := common.UUIDStrToBytes(in.Id)
	if err != nil {
		cs.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	stmt, err := db.PreparexContext(ctx, updateCreditNoteHeaderSQL)
	if err != nil {
		cs.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	err = cs.DBService.InsUpd(ctx, in.GetUserEmail(), in.GetRequestId(), func(tx *sqlx.Tx) error {
		_, err = tx.StmtxContext(ctx, stmt).ExecContext(ctx,
			in.Note,
			in.TaxCurrencyCode,
			in.ChargeTotalAmount,
			in.PrepaidAmount,
			in.PayableRoundingAmount,
			in.PayableAmount,
			tn,
			uuid4byte)
		if err != nil {
			cs.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
			err1 := stmt.Close()
			if err1 != nil {
				cs.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err1))
				return err1
			}
			return err
		}
		return nil
	})

	if err != nil {
		cs.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	err = stmt.Close()
	if err != nil {
		cs.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	return &invoiceproto.UpdateCreditNoteHeaderResponse{}, nil
}
