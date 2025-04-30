package invoiceservices

import (
	"context"
	"net"
	"os"
	"time"

	"github.com/cloudfresco/sc-ubl/internal/common"
	"github.com/cloudfresco/sc-ubl/internal/config"
	commonproto "github.com/cloudfresco/sc-ubl/internal/protogen/common/v1"
	invoiceproto "github.com/cloudfresco/sc-ubl/internal/protogen/invoice/v1"
	partyproto "github.com/cloudfresco/sc-ubl/internal/protogen/party/v1"
	partyservice "github.com/cloudfresco/sc-ubl/internal/services/partyservices"
	commonstruct "github.com/cloudfresco/sc-ubl/internal/servicestructs/common/v1"
	invoicestruct "github.com/cloudfresco/sc-ubl/internal/servicestructs/invoice/v1"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// InvoiceService - For accessing Invoice services
type InvoiceService struct {
	log               *zap.Logger
	DBService         *common.DBService
	RedisService      *common.RedisService
	UserServiceClient partyproto.UserServiceClient
	invoiceproto.UnimplementedInvoiceServiceServer
}

// NewInvoiceService - Create Invoice service
func NewInvoiceService(log *zap.Logger, dbOpt *common.DBService, redisOpt *common.RedisService, userServiceClient partyproto.UserServiceClient) *InvoiceService {
	return &InvoiceService{
		log:               log,
		DBService:         dbOpt,
		RedisService:      redisOpt,
		UserServiceClient: userServiceClient,
	}
}

// StartInvoiceServer - Start Invoice server
func StartInvoiceServer(log *zap.Logger, isTest bool, pwd string, dbOpt *config.DBOptions, redisOpt *config.RedisOptions, mailerOpt *config.MailerOptions, grpcServerOpt *config.GrpcServerOptions, jwtOpt *config.JWTOptions, oauthOpt *config.OauthOptions, userOpt *config.UserOptions, uptraceOpt *config.UptraceOptions, dbService *common.DBService, redisService *common.RedisService, mailerService common.MailerIntf) {
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
	creditNoteHeaderService := NewCreditNoteHeaderService(log, dbService, redisService, uc)
	debitNoteService := NewDebitNoteHeaderService(log, dbService, redisService, uc)
	invoiceService := NewInvoiceService(log, dbService, redisService, uc)

	lis, err := net.Listen("tcp", grpcServerOpt.GrpcInvoiceServerPort)
	if err != nil {
		log.Error("Error", zap.Error(err))
		os.Exit(1)
	}

	srv := grpc.NewServer(srvOpts...)
	invoiceproto.RegisterCreditNoteHeaderServiceServer(srv, creditNoteHeaderService)
	invoiceproto.RegisterDebitNoteHeaderServiceServer(srv, debitNoteService)
	invoiceproto.RegisterInvoiceServiceServer(srv, invoiceService)

	if err := srv.Serve(lis); err != nil {
		log.Error("Error", zap.Error(err))
		os.Exit(1)
	}
}

const insertInvoiceHeaderSQL = `insert into invoice_headers
	    (uuid4,
ih_id,
invoice_type_code,
note,
document_currency_code,
tax_currency_code,
pricing_currency_code,
payment_currency_code,
payment_alt_currency_code,
accounting_cost_code,
accounting_cost,
line_count_numeric,
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
:ih_id,
:invoice_type_code,
:note,
:document_currency_code,
:tax_currency_code,
:pricing_currency_code,
:payment_currency_code,
:payment_alt_currency_code,
:accounting_cost_code,
:accounting_cost,
:line_count_numeric,
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

const selectInvoiceHeadersSQL = `select 
id,
uuid4,
ih_id,
invoice_type_code,
note,
document_currency_code,
tax_currency_code,
pricing_currency_code,
payment_currency_code,
payment_alt_currency_code,
accounting_cost_code,
accounting_cost,
line_count_numeric,
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
updated_at from invoice_headers`

// updateInvoiceSQL - update InvoiceSQL query
const updateInvoiceSQL = `update invoice_headers set 
note = ?, 
invoice_type_code = ?, 
charge_total_amount = ?, 
prepaid_amount = ?, 
payable_rounding_amount = ?,
payable_amount = ?, updated_at = ? where uuid4 = ?;`

// CreateInvoice - Create Invoice
func (is *InvoiceService) CreateInvoice(ctx context.Context, in *invoiceproto.CreateInvoiceRequest) (*invoiceproto.CreateInvoiceResponse, error) {
	user, err := partyservice.GetUserWithNewContext(ctx, in.UserId, in.UserEmail, in.RequestId, is.UserServiceClient)
	if err != nil {
		is.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	issueDate, err := time.Parse(common.Layout, in.IssueDate)
	if err != nil {
		is.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	dueDate, err := time.Parse(common.Layout, in.DueDate)
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

	taxPointDate, err := time.Parse(common.Layout, in.TaxPointDate)
	if err != nil {
		is.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	taxExDate, err := time.Parse(common.Layout, in.TaxExDate)
	if err != nil {
		is.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	pricingExDate, err := time.Parse(common.Layout, in.PricingExDate)
	if err != nil {
		is.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	paymentExDate, err := time.Parse(common.Layout, in.PaymentExDate)
	if err != nil {
		is.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	paymentAltExDate, err := time.Parse(common.Layout, in.PaymentAltExDate)
	if err != nil {
		is.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	ttime := common.GetTimeDetails()
	tn := common.TimeToTimestamp(ttime)

	invoiceHeaderD := invoiceproto.InvoiceHeaderD{}
	invoiceHeaderD.Uuid4, err = common.GetUUIDBytes()
	if err != nil {
		is.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	invoiceHeaderD.IhId = in.IhId
	invoiceHeaderD.InvoiceTypeCode = in.InvoiceTypeCode
	invoiceHeaderD.Note = in.Note
	invoiceHeaderD.DocumentCurrencyCode = in.DocumentCurrencyCode
	invoiceHeaderD.TaxCurrencyCode = in.TaxCurrencyCode
	invoiceHeaderD.PricingCurrencyCode = in.PricingCurrencyCode
	invoiceHeaderD.PaymentCurrencyCode = in.PaymentCurrencyCode
	invoiceHeaderD.PaymentAltCurrencyCode = in.PaymentAltCurrencyCode
	invoiceHeaderD.AccountingCostCode = in.AccountingCostCode
	invoiceHeaderD.AccountingCost = in.AccountingCost
	invoiceHeaderD.LineCountNumeric = in.LineCountNumeric
	invoiceHeaderD.OrderId = in.OrderId
	invoiceHeaderD.BillingId = in.BillingId
	invoiceHeaderD.DespatchId = in.DespatchId
	invoiceHeaderD.ReceiptId = in.ReceiptId
	invoiceHeaderD.StatementId = in.StatementId
	invoiceHeaderD.ContractId = in.ContractId
	invoiceHeaderD.AccountingSupplierPartyId = in.AccountingSupplierPartyId
	invoiceHeaderD.AccountingCustomerPartyId = in.AccountingCustomerPartyId
	invoiceHeaderD.PayeePartyId = in.PayeePartyId
	invoiceHeaderD.BuyerCustomerPartyId = in.BuyerCustomerPartyId
	invoiceHeaderD.SellerSupplierPartyId = in.SellerSupplierPartyId
	invoiceHeaderD.TaxRepresentativePartyId = in.TaxRepresentativePartyId
	invoiceHeaderD.TaxExSourceCurrencyCode = in.TaxExSourceCurrencyCode
	invoiceHeaderD.TaxExSourceCurrencyBaseRate = in.TaxExSourceCurrencyBaseRate
	invoiceHeaderD.TaxExTargetCurrencyCode = in.TaxExTargetCurrencyCode
	invoiceHeaderD.TaxExTargetCurrencyBaseRate = in.TaxExTargetCurrencyBaseRate
	invoiceHeaderD.TaxExExchangeMarketId = in.TaxExExchangeMarketId
	invoiceHeaderD.TaxExCalculationRate = in.TaxExCalculationRate
	invoiceHeaderD.TaxExMathematicOperatorCode = in.TaxExMathematicOperatorCode

	invoiceHeaderD.PricingExSourceCurrencyCode = in.PricingExSourceCurrencyCode
	invoiceHeaderD.PricingExSourceCurrencyBaseRate = in.PricingExSourceCurrencyBaseRate
	invoiceHeaderD.PricingExTargetCurrencyCode = in.PricingExTargetCurrencyCode
	invoiceHeaderD.PricingExTargetCurrencyBaseRate = in.PricingExTargetCurrencyBaseRate
	invoiceHeaderD.PricingExExchangeMarketId = in.PricingExExchangeMarketId
	invoiceHeaderD.PricingExCalculationRate = in.PricingExCalculationRate
	invoiceHeaderD.PricingExMathematicOperatorCode = in.PricingExMathematicOperatorCode

	invoiceHeaderD.PaymentExSourceCurrencyCode = in.PaymentExSourceCurrencyCode
	invoiceHeaderD.PaymentExSourceCurrencyBaseRate = in.PaymentExSourceCurrencyBaseRate
	invoiceHeaderD.PaymentExTargetCurrencyCode = in.PaymentExTargetCurrencyCode
	invoiceHeaderD.PaymentExTargetCurrencyBaseRate = in.PaymentExTargetCurrencyBaseRate
	invoiceHeaderD.PaymentExExchangeMarketId = in.PaymentExExchangeMarketId
	invoiceHeaderD.PaymentExCalculationRate = in.PaymentExCalculationRate
	invoiceHeaderD.PaymentExMathematicOperatorCode = in.PaymentExMathematicOperatorCode

	invoiceHeaderD.PaymentAltExSourceCurrencyCode = in.PaymentAltExSourceCurrencyCode
	invoiceHeaderD.PaymentAltExSourceCurrencyBaseRate = in.PaymentAltExSourceCurrencyCode
	invoiceHeaderD.PaymentAltExTargetCurrencyCode = in.PaymentAltExTargetCurrencyCode
	invoiceHeaderD.PaymentAltExTargetCurrencyBaseRate = in.PaymentAltExTargetCurrencyBaseRate
	invoiceHeaderD.PaymentAltExExchangeMarketId = in.PaymentAltExExchangeMarketId
	invoiceHeaderD.PaymentAltExCalculationRate = in.PaymentAltExCalculationRate
	invoiceHeaderD.PaymentAltExMathematicOperatorCode = in.PaymentAltExMathematicOperatorCode

	invoiceHeaderD.LineExtensionAmount = in.LineExtensionAmount
	invoiceHeaderD.TaxExclusiveAmount = in.TaxExclusiveAmount
	invoiceHeaderD.TaxInclusiveAmount = in.TaxInclusiveAmount
	invoiceHeaderD.AllowanceTotalAmount = in.AllowanceTotalAmount
	invoiceHeaderD.ChargeTotalAmount = in.ChargeTotalAmount
	invoiceHeaderD.WithholdingTaxTotalAmount = in.WithholdingTaxTotalAmount
	invoiceHeaderD.PrepaidAmount = in.PrepaidAmount
	invoiceHeaderD.PayableRoundingAmount = in.PayableRoundingAmount
	invoiceHeaderD.PayableAmount = in.PayableAmount
	invoiceHeaderD.PayableAlternativeAmount = in.PayableAlternativeAmount

	invoiceHeaderT := invoiceproto.InvoiceHeaderT{}
	invoiceHeaderT.IssueDate = common.TimeToTimestamp(issueDate.UTC().Truncate(time.Second))
	invoiceHeaderT.DueDate = common.TimeToTimestamp(dueDate.UTC().Truncate(time.Second))
	invoiceHeaderT.TaxPointDate = common.TimeToTimestamp(taxPointDate.UTC().Truncate(time.Second))
	invoiceHeaderT.InvoicePeriodStartDate = common.TimeToTimestamp(invoicePeriodStartDate.UTC().Truncate(time.Second))
	invoiceHeaderT.InvoicePeriodEndDate = common.TimeToTimestamp(invoicePeriodEndDate.UTC().Truncate(time.Second))
	invoiceHeaderT.TaxExDate = common.TimeToTimestamp(taxExDate.UTC().Truncate(time.Second))
	invoiceHeaderT.PricingExDate = common.TimeToTimestamp(pricingExDate.UTC().Truncate(time.Second))
	invoiceHeaderT.PaymentExDate = common.TimeToTimestamp(paymentExDate.UTC().Truncate(time.Second))
	invoiceHeaderT.PaymentAltExDate = common.TimeToTimestamp(paymentAltExDate.UTC().Truncate(time.Second))

	crUpdUser := commonproto.CrUpdUser{}
	crUpdUser.StatusCode = "active"
	crUpdUser.CreatedByUserId = user.Id
	crUpdUser.UpdatedByUserId = user.Id

	crUpdTime := commonproto.CrUpdTime{}
	crUpdTime.CreatedAt = tn
	crUpdTime.UpdatedAt = tn

	invoiceHeader := invoiceproto.InvoiceHeader{InvoiceHeaderD: &invoiceHeaderD, InvoiceHeaderT: &invoiceHeaderT, CrUpdUser: &crUpdUser, CrUpdTime: &crUpdTime}

	invoiceLines := []*invoiceproto.InvoiceLine{}
	// we will do for loop on lines which is comes from client form
	for _, line := range in.InvoiceLines {
		line.UserId = in.UserId
		line.UserEmail = in.UserEmail
		line.RequestId = in.RequestId
		invoiceLine, err := is.ProcessInvoiceLineRequest(ctx, line)
		if err != nil {
			is.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
			return nil, err
		}

		invoiceLines = append(invoiceLines, invoiceLine)
	}

	err = is.insertInvoiceHeader(ctx, insertInvoiceHeaderSQL, &invoiceHeader, insertInvoiceLineSQL, invoiceLines, in.GetUserEmail(), in.GetRequestId())
	if err != nil {
		is.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	invoiceHeaderResponse := invoiceproto.CreateInvoiceResponse{}
	invoiceHeaderResponse.InvoiceHeader = &invoiceHeader
	return &invoiceHeaderResponse, nil
}

// insertInvoiceHeader - insertInvoiceHeader
func (is *InvoiceService) insertInvoiceHeader(ctx context.Context, insertInvoiceHeaderSQL string, invoiceHeader *invoiceproto.InvoiceHeader, insertInvoiceLineSQL string, invoiceLines []*invoiceproto.InvoiceLine, userEmail string, requestID string) error {
	invoiceHeaderTmp, err := is.crInvoiceHeaderStruct(ctx, invoiceHeader, userEmail, requestID)
	if err != nil {
		is.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
		return err
	}
	err = is.DBService.InsUpd(ctx, userEmail, requestID, func(tx *sqlx.Tx) error {
		res, err := tx.NamedExecContext(ctx, insertInvoiceHeaderSQL, invoiceHeaderTmp)
		if err != nil {
			is.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}

		uID, err := res.LastInsertId()
		if err != nil {
			is.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}
		invoiceHeader.InvoiceHeaderD.Id = uint32(uID)
		uuid4Str, err := common.UUIDBytesToStr(invoiceHeader.InvoiceHeaderD.Uuid4)
		if err != nil {
			is.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}
		invoiceHeader.InvoiceHeaderD.IdS = uuid4Str

		for _, invoiceLine := range invoiceLines {
			invoiceLine.InvoiceLineD.InvoiceHeaderId = invoiceHeader.InvoiceHeaderD.Id
			invoiceLineTmp, err := is.crInvoiceLineStruct(ctx, invoiceLine, userEmail, requestID)
			if err != nil {
				is.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
				return err
			}
			_, err = tx.NamedExecContext(ctx, insertInvoiceLineSQL, invoiceLineTmp)
			if err != nil {
				is.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
				return err
			}
		}

		return nil
	})

	if err != nil {
		is.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
		return err
	}
	return nil
}

// crInvoiceHeaderStruct - process InvoiceHeader details
func (is *InvoiceService) crInvoiceHeaderStruct(ctx context.Context, invoiceHeader *invoiceproto.InvoiceHeader, userEmail string, requestID string) (*invoicestruct.InvoiceHeader, error) {
	crUpdTime := new(commonstruct.CrUpdTime)
	crUpdTime.CreatedAt = common.TimestampToTime(invoiceHeader.CrUpdTime.CreatedAt)
	crUpdTime.UpdatedAt = common.TimestampToTime(invoiceHeader.CrUpdTime.UpdatedAt)

	invoiceHeaderT := new(invoicestruct.InvoiceHeaderT)
	invoiceHeaderT.IssueDate = common.TimestampToTime(invoiceHeader.InvoiceHeaderT.IssueDate)
	invoiceHeaderT.DueDate = common.TimestampToTime(invoiceHeader.InvoiceHeaderT.DueDate)
	invoiceHeaderT.TaxPointDate = common.TimestampToTime(invoiceHeader.InvoiceHeaderT.TaxPointDate)
	invoiceHeaderT.InvoicePeriodStartDate = common.TimestampToTime(invoiceHeader.InvoiceHeaderT.InvoicePeriodStartDate)
	invoiceHeaderT.InvoicePeriodEndDate = common.TimestampToTime(invoiceHeader.InvoiceHeaderT.InvoicePeriodEndDate)
	invoiceHeaderT.TaxExDate = common.TimestampToTime(invoiceHeader.InvoiceHeaderT.TaxExDate)
	invoiceHeaderT.PricingExDate = common.TimestampToTime(invoiceHeader.InvoiceHeaderT.PricingExDate)
	invoiceHeaderT.PaymentExDate = common.TimestampToTime(invoiceHeader.InvoiceHeaderT.PaymentExDate)
	invoiceHeaderT.PaymentAltExDate = common.TimestampToTime(invoiceHeader.InvoiceHeaderT.PaymentAltExDate)

	invoiceHeaderTmp := invoicestruct.InvoiceHeader{InvoiceHeaderD: invoiceHeader.InvoiceHeaderD, InvoiceHeaderT: invoiceHeaderT, CrUpdUser: invoiceHeader.CrUpdUser, CrUpdTime: crUpdTime}

	return &invoiceHeaderTmp, nil
}

// GetInvoices - Get Invoices
func (is *InvoiceService) GetInvoices(ctx context.Context, in *invoiceproto.GetInvoicesRequest) (*invoiceproto.GetInvoicesResponse, error) {
	limit := in.GetLimit()
	nextCursor := in.GetNextCursor()
	if limit == "" {
		limit = is.DBService.LimitSQLRows
	}
	query := "status_code = ?"
	if nextCursor == "" {
		query = query + " order by id desc " + " limit " + limit + ";"
	} else {
		nextCursor = common.DecodeCursor(nextCursor)
		query = query + " " + "and" + " " + "id <= " + nextCursor + " order by id desc " + " limit " + limit + ";"
	}

	invoiceHeaders := []*invoiceproto.InvoiceHeader{}

	nselectInvoiceHeadersSQL := selectInvoiceHeadersSQL + ` where ` + query

	rows, err := is.DBService.DB.QueryxContext(ctx, nselectInvoiceHeadersSQL, "active")
	if err != nil {
		is.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	for rows.Next() {

		invoiceHeaderTmp := invoicestruct.InvoiceHeader{}
		err = rows.StructScan(&invoiceHeaderTmp)
		if err != nil {
			is.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
			return nil, err
		}

		getRequest := commonproto.GetRequest{}
		getRequest.UserEmail = in.UserEmail
		getRequest.RequestId = in.RequestId
		invoiceHeader, err := is.getInvoiceHeaderStruct(ctx, &getRequest, invoiceHeaderTmp)
		if err != nil {
			is.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
			return nil, err
		}
		invoiceHeaders = append(invoiceHeaders, invoiceHeader)

	}

	invoiceHeadersResponse := invoiceproto.GetInvoicesResponse{}
	if len(invoiceHeaders) != 0 {
		next := invoiceHeaders[len(invoiceHeaders)-1].InvoiceHeaderD.Id
		next--
		nextc := common.EncodeCursor(next)
		invoiceHeadersResponse = invoiceproto.GetInvoicesResponse{InvoiceHeaders: invoiceHeaders, NextCursor: nextc}
	} else {
		invoiceHeadersResponse = invoiceproto.GetInvoicesResponse{InvoiceHeaders: invoiceHeaders, NextCursor: "0"}
	}
	return &invoiceHeadersResponse, nil
}

// GetInvoice - Get Invoice
func (is *InvoiceService) GetInvoice(ctx context.Context, inReq *invoiceproto.GetInvoiceRequest) (*invoiceproto.GetInvoiceResponse, error) {
	in := inReq.GetRequest
	uuid4byte, err := common.UUIDStrToBytes(in.Id)
	if err != nil {
		is.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	nselectInvoiceHeadersSQL := selectInvoiceHeadersSQL + ` where uuid4 = ? and status_code = ?;`
	row := is.DBService.DB.QueryRowxContext(ctx, nselectInvoiceHeadersSQL, uuid4byte, "active")
	invoiceHeaderTmp := invoicestruct.InvoiceHeader{}
	err = row.StructScan(&invoiceHeaderTmp)
	if err != nil {
		is.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	invoiceHeader, err := is.getInvoiceHeaderStruct(ctx, in, invoiceHeaderTmp)
	if err != nil {
		is.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	invoiceHeaderResponse := invoiceproto.GetInvoiceResponse{}
	invoiceHeaderResponse.InvoiceHeader = invoiceHeader
	return &invoiceHeaderResponse, nil
}

// GetInvoiceByPk - Get Invoice By Primary key(Id)
func (is *InvoiceService) GetInvoiceByPk(ctx context.Context, inReq *invoiceproto.GetInvoiceByPkRequest) (*invoiceproto.GetInvoiceByPkResponse, error) {
	in := inReq.GetByIdRequest
	nselectInvoiceHeadersSQL := selectInvoiceHeadersSQL + ` where id = ? and status_code = ?;`
	row := is.DBService.DB.QueryRowxContext(ctx, nselectInvoiceHeadersSQL, in.Id, "active")
	invoiceHeaderTmp := invoicestruct.InvoiceHeader{}
	err := row.StructScan(&invoiceHeaderTmp)
	if err != nil {
		is.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	getRequest := commonproto.GetRequest{}
	getRequest.UserEmail = in.UserEmail
	getRequest.RequestId = in.RequestId
	invoiceHeader, err := is.getInvoiceHeaderStruct(ctx, &getRequest, invoiceHeaderTmp)
	if err != nil {
		is.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	invoiceHeaderResponse := invoiceproto.GetInvoiceByPkResponse{}
	invoiceHeaderResponse.InvoiceHeader = invoiceHeader
	return &invoiceHeaderResponse, nil
}

// getInvoiceHeaderStruct - Get invoice
func (is *InvoiceService) getInvoiceHeaderStruct(ctx context.Context, in *commonproto.GetRequest, invoiceHeaderTmp invoicestruct.InvoiceHeader) (*invoiceproto.InvoiceHeader, error) {
	invoiceHeaderT := new(invoiceproto.InvoiceHeaderT)
	invoiceHeaderT.IssueDate = common.TimeToTimestamp(invoiceHeaderTmp.InvoiceHeaderT.IssueDate)
	invoiceHeaderT.DueDate = common.TimeToTimestamp(invoiceHeaderTmp.InvoiceHeaderT.DueDate)
	invoiceHeaderT.TaxPointDate = common.TimeToTimestamp(invoiceHeaderTmp.InvoiceHeaderT.TaxPointDate)
	invoiceHeaderT.TaxExDate = common.TimeToTimestamp(invoiceHeaderTmp.InvoiceHeaderT.TaxExDate)
	invoiceHeaderT.InvoicePeriodStartDate = common.TimeToTimestamp(invoiceHeaderTmp.InvoiceHeaderT.InvoicePeriodStartDate)
	invoiceHeaderT.InvoicePeriodEndDate = common.TimeToTimestamp(invoiceHeaderTmp.InvoiceHeaderT.InvoicePeriodEndDate)
	invoiceHeaderT.PricingExDate = common.TimeToTimestamp(invoiceHeaderTmp.InvoiceHeaderT.PricingExDate)
	invoiceHeaderT.PaymentExDate = common.TimeToTimestamp(invoiceHeaderTmp.InvoiceHeaderT.PaymentExDate)
	invoiceHeaderT.PaymentAltExDate = common.TimeToTimestamp(invoiceHeaderTmp.InvoiceHeaderT.PaymentAltExDate)

	uuid4Str, err := common.UUIDBytesToStr(invoiceHeaderTmp.InvoiceHeaderD.Uuid4)
	if err != nil {
		is.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	invoiceHeaderTmp.InvoiceHeaderD.IdS = uuid4Str

	crUpdTime := new(commonproto.CrUpdTime)
	crUpdTime.CreatedAt = common.TimeToTimestamp(invoiceHeaderTmp.CrUpdTime.CreatedAt)
	crUpdTime.UpdatedAt = common.TimeToTimestamp(invoiceHeaderTmp.CrUpdTime.UpdatedAt)

	invoiceHeader := invoiceproto.InvoiceHeader{InvoiceHeaderD: invoiceHeaderTmp.InvoiceHeaderD, InvoiceHeaderT: invoiceHeaderT, CrUpdUser: invoiceHeaderTmp.CrUpdUser, CrUpdTime: crUpdTime}

	return &invoiceHeader, nil
}

// UpdateInvoice - Update Invoice
func (is *InvoiceService) UpdateInvoice(ctx context.Context, in *invoiceproto.UpdateInvoiceRequest) (*invoiceproto.UpdateInvoiceResponse, error) {
	db := is.DBService.DB
	tn := common.GetTimeDetails()

	uuid4byte, err := common.UUIDStrToBytes(in.Id)
	if err != nil {
		is.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	stmt, err := db.PreparexContext(ctx, updateInvoiceSQL)
	if err != nil {
		is.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	err = is.DBService.InsUpd(ctx, in.GetUserEmail(), in.GetRequestId(), func(tx *sqlx.Tx) error {
		_, err = tx.StmtxContext(ctx, stmt).ExecContext(ctx,
			in.Note,
			in.InvoiceTypeCode,
			in.ChargeTotalAmount,
			in.PrepaidAmount,
			in.PayableRoundingAmount,
			in.PayableAmount,
			tn,
			uuid4byte)
		if err != nil {
			is.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
			err1 := stmt.Close()
			if err1 != nil {
				is.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err1))
				return err1
			}
			return err
		}
		return nil
	})

	if err != nil {
		is.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	err = stmt.Close()
	if err != nil {
		is.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	return &invoiceproto.UpdateInvoiceResponse{}, nil
}
