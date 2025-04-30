package orderservices

import (
	"context"
	"net"
	"os"
	"time"

	"github.com/cloudfresco/sc-ubl/internal/common"
	"github.com/cloudfresco/sc-ubl/internal/config"
	commonproto "github.com/cloudfresco/sc-ubl/internal/protogen/common/v1"
	orderproto "github.com/cloudfresco/sc-ubl/internal/protogen/order/v1"
	partyproto "github.com/cloudfresco/sc-ubl/internal/protogen/party/v1"
	partyservice "github.com/cloudfresco/sc-ubl/internal/services/partyservices"
	commonstruct "github.com/cloudfresco/sc-ubl/internal/servicestructs/common/v1"
	orderstruct "github.com/cloudfresco/sc-ubl/internal/servicestructs/order/v1"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// PurchaseOrderHeaderService - For accessing PurchaseOrder services
type PurchaseOrderHeaderService struct {
	log               *zap.Logger
	DBService         *common.DBService
	RedisService      *common.RedisService
	UserServiceClient partyproto.UserServiceClient
	orderproto.UnimplementedPurchaseOrderHeaderServiceServer
}

// NewPurchaseOrderHeaderService - Create PurchaseOrder service
func NewPurchaseOrderHeaderService(log *zap.Logger, dbOpt *common.DBService, redisOpt *common.RedisService, userServiceClient partyproto.UserServiceClient) *PurchaseOrderHeaderService {
	return &PurchaseOrderHeaderService{
		log:               log,
		DBService:         dbOpt,
		RedisService:      redisOpt,
		UserServiceClient: userServiceClient,
	}
}

// StartOrderServer - Start Order server
func StartOrderServer(log *zap.Logger, isTest bool, pwd string, dbOpt *config.DBOptions, redisOpt *config.RedisOptions, mailerOpt *config.MailerOptions, grpcServerOpt *config.GrpcServerOptions, jwtOpt *config.JWTOptions, oauthOpt *config.OauthOptions, userOpt *config.UserOptions, uptraceOpt *config.UptraceOptions, dbService *common.DBService, redisService *common.RedisService, mailerService common.MailerIntf) {
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
	purchaseOrderService := NewPurchaseOrderHeaderService(log, dbService, redisService, uc)

	lis, err := net.Listen("tcp", grpcServerOpt.GrpcOrderServerPort)
	if err != nil {
		log.Error("Error", zap.Error(err))
		os.Exit(1)
	}

	srv := grpc.NewServer(srvOpts...)
	orderproto.RegisterPurchaseOrderHeaderServiceServer(srv, purchaseOrderService)

	if err := srv.Serve(lis); err != nil {
		log.Error("Error", zap.Error(err))
		os.Exit(1)
	}
}

const insertPurchaseOrderHeaderSQL = `insert into purchase_order_headers
	    (uuid4,
poh_id,
sales_order_id,
order_type_code,
note,
requested_invoice_currency_code,
document_currency_code,
pricing_currency_code,
tax_currency_code,
accounting_cost_code,
accounting_cost,
line_count_numeric,
quotation_id,
order_id,
catalogue_id,
buyer_customer_party_id,
seller_supplier_party_id,
originator_customer_party_id,
freight_forwarder_party_id,
accounting_customer_party_id,
transaction_conditions,
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
destination_country,
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
validity_period,
tax_ex_date,
pricing_ex_date,
payment_ex_date,
status_code,
created_by_user_id,
updated_by_user_id,
created_at,
updated_at)
  values(:uuid4,
:poh_id,
:sales_order_id,
:order_type_code,
:note,
:requested_invoice_currency_code,
:document_currency_code,
:pricing_currency_code,
:tax_currency_code,
:accounting_cost_code,
:accounting_cost,
:line_count_numeric,
:quotation_id,
:order_id,
:catalogue_id,
:buyer_customer_party_id,
:seller_supplier_party_id,
:originator_customer_party_id,
:freight_forwarder_party_id,
:accounting_customer_party_id,
:transaction_conditions,
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
:destination_country,
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
:validity_period,
:tax_ex_date,
:pricing_ex_date,
:payment_ex_date,
:status_code,
:created_by_user_id,
:updated_by_user_id,
:created_at,
:updated_at);`

const selectPurchaseOrderHeadersSQL = `select 
    id,
uuid4,
poh_id,
sales_order_id,
order_type_code,
note,
requested_invoice_currency_code,
document_currency_code,
pricing_currency_code,
tax_currency_code,
accounting_cost_code,
accounting_cost,
line_count_numeric,
quotation_id,
order_id,
catalogue_id,
buyer_customer_party_id,
seller_supplier_party_id,
originator_customer_party_id,
freight_forwarder_party_id,
accounting_customer_party_id,
transaction_conditions,
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
destination_country,
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
validity_period,
tax_ex_date,
pricing_ex_date,
payment_ex_date,
status_code,
created_by_user_id,
updated_by_user_id,
created_at,
updated_at from purchase_order_headers`

// updatePurchaseOrderHeaderSQL - update PurchaseOrderHeaderSQL query
const updatePurchaseOrderHeaderSQL = `update purchase_order_headers set 
  order_type_code= ?,
  note= ?,
  requested_invoice_currency_code= ?,
  document_currency_code= ?,
  pricing_currency_code= ?,
  tax_currency_code= ?,
  accounting_cost_code= ?,
  accounting_cost= ?,
  updated_at = ? where uuid4 = ?;`

// CreatePurchaseOrderHeader - Create PurchaseOrderHeader
func (ps *PurchaseOrderHeaderService) CreatePurchaseOrderHeader(ctx context.Context, in *orderproto.CreatePurchaseOrderHeaderRequest) (*orderproto.CreatePurchaseOrderHeaderResponse, error) {
	user, err := partyservice.GetUserWithNewContext(ctx, in.UserId, in.UserEmail, in.RequestId, ps.UserServiceClient)
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	issueDate, err := time.Parse(common.Layout, in.IssueDate)
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	validityPeriod, err := time.Parse(common.Layout, in.ValidityPeriod)
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	taxExDate, err := time.Parse(common.Layout, in.TaxExDate)
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	pricingExDate, err := time.Parse(common.Layout, in.PricingExDate)
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	paymentExDate, err := time.Parse(common.Layout, in.PaymentExDate)
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	ttime := common.GetTimeDetails()
	tn := common.TimeToTimestamp(ttime)

	purchaseOrderHeaderD := orderproto.PurchaseOrderHeaderD{}
	purchaseOrderHeaderD.Uuid4, err = common.GetUUIDBytes()
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	purchaseOrderHeaderD.PohId = in.PohId
	purchaseOrderHeaderD.SalesOrderId = in.SalesOrderId
	purchaseOrderHeaderD.OrderTypeCode = in.OrderTypeCode
	purchaseOrderHeaderD.Note = in.Note

	purchaseOrderHeaderD.RequestedInvoiceCurrencyCode = in.RequestedInvoiceCurrencyCode
	purchaseOrderHeaderD.DocumentCurrencyCode = in.DocumentCurrencyCode
	purchaseOrderHeaderD.PricingCurrencyCode = in.PricingCurrencyCode
	purchaseOrderHeaderD.TaxCurrencyCode = in.TaxCurrencyCode

	purchaseOrderHeaderD.AccountingCostCode = in.AccountingCostCode
	purchaseOrderHeaderD.AccountingCost = in.AccountingCost
	purchaseOrderHeaderD.LineCountNumeric = in.LineCountNumeric

	purchaseOrderHeaderD.QuotationId = in.QuotationId
	purchaseOrderHeaderD.OrderId = in.OrderId
	purchaseOrderHeaderD.CatalogueId = in.CatalogueId

	purchaseOrderHeaderD.BuyerCustomerPartyId = in.BuyerCustomerPartyId
	purchaseOrderHeaderD.SellerSupplierPartyId = in.SellerSupplierPartyId
	purchaseOrderHeaderD.OriginatorCustomerPartyId = in.OriginatorCustomerPartyId
	purchaseOrderHeaderD.FreightForwarderPartyId = in.FreightForwarderPartyId
	purchaseOrderHeaderD.AccountingCustomerPartyId = in.AccountingCustomerPartyId

	purchaseOrderHeaderD.TransactionConditions = in.TransactionConditions

	purchaseOrderHeaderD.TaxExSourceCurrencyCode = in.TaxExSourceCurrencyCode
	purchaseOrderHeaderD.TaxExSourceCurrencyBaseRate = in.TaxExSourceCurrencyBaseRate
	purchaseOrderHeaderD.TaxExTargetCurrencyCode = in.TaxExTargetCurrencyCode
	purchaseOrderHeaderD.TaxExTargetCurrencyBaseRate = in.TaxExTargetCurrencyBaseRate
	purchaseOrderHeaderD.TaxExExchangeMarketId = in.TaxExExchangeMarketId
	purchaseOrderHeaderD.TaxExCalculationRate = in.TaxExCalculationRate
	purchaseOrderHeaderD.TaxExMathematicOperatorCode = in.TaxExMathematicOperatorCode

	purchaseOrderHeaderD.PricingExSourceCurrencyCode = in.PricingExSourceCurrencyCode
	purchaseOrderHeaderD.PricingExSourceCurrencyBaseRate = in.PricingExSourceCurrencyBaseRate
	purchaseOrderHeaderD.PricingExTargetCurrencyCode = in.PricingExTargetCurrencyCode
	purchaseOrderHeaderD.PricingExTargetCurrencyBaseRate = in.PricingExTargetCurrencyBaseRate
	purchaseOrderHeaderD.PricingExExchangeMarketId = in.PricingExExchangeMarketId
	purchaseOrderHeaderD.PricingExCalculationRate = in.PricingExCalculationRate
	purchaseOrderHeaderD.PricingExMathematicOperatorCode = in.PricingExMathematicOperatorCode

	purchaseOrderHeaderD.PaymentExSourceCurrencyCode = in.PaymentExSourceCurrencyCode
	purchaseOrderHeaderD.PaymentExSourceCurrencyBaseRate = in.PaymentExSourceCurrencyBaseRate
	purchaseOrderHeaderD.PaymentExTargetCurrencyCode = in.PaymentExTargetCurrencyCode
	purchaseOrderHeaderD.PaymentExTargetCurrencyBaseRate = in.PaymentExTargetCurrencyBaseRate
	purchaseOrderHeaderD.PaymentExExchangeMarketId = in.PaymentExExchangeMarketId
	purchaseOrderHeaderD.PaymentExCalculationRate = in.PaymentExCalculationRate
	purchaseOrderHeaderD.PaymentExMathematicOperatorCode = in.PaymentExMathematicOperatorCode

	purchaseOrderHeaderD.LineExtensionAmount = in.LineExtensionAmount
	purchaseOrderHeaderD.TaxExclusiveAmount = in.TaxExclusiveAmount
	purchaseOrderHeaderD.TaxInclusiveAmount = in.TaxInclusiveAmount
	purchaseOrderHeaderD.AllowanceTotalAmount = in.AllowanceTotalAmount
	purchaseOrderHeaderD.ChargeTotalAmount = in.ChargeTotalAmount
	purchaseOrderHeaderD.WithholdingTaxTotalAmount = in.WithholdingTaxTotalAmount
	purchaseOrderHeaderD.PrepaidAmount = in.PrepaidAmount
	purchaseOrderHeaderD.PayableRoundingAmount = in.PayableRoundingAmount
	purchaseOrderHeaderD.PayableAmount = in.PayableAmount
	purchaseOrderHeaderD.PayableAlternativeAmount = in.PayableAlternativeAmount

	purchaseOrderHeaderT := orderproto.PurchaseOrderHeaderT{}
	purchaseOrderHeaderT.IssueDate = common.TimeToTimestamp(issueDate.UTC().Truncate(time.Second))
	purchaseOrderHeaderT.ValidityPeriod = common.TimeToTimestamp(validityPeriod.UTC().Truncate(time.Second))
	purchaseOrderHeaderT.TaxExDate = common.TimeToTimestamp(taxExDate.UTC().Truncate(time.Second))
	purchaseOrderHeaderT.PaymentExDate = common.TimeToTimestamp(paymentExDate.UTC().Truncate(time.Second))
	purchaseOrderHeaderT.PricingExDate = common.TimeToTimestamp(pricingExDate.UTC().Truncate(time.Second))

	crUpdUser := commonproto.CrUpdUser{}
	crUpdUser.StatusCode = "active"
	crUpdUser.CreatedByUserId = user.Id
	crUpdUser.UpdatedByUserId = user.Id

	crUpdTime := commonproto.CrUpdTime{}
	crUpdTime.CreatedAt = tn
	crUpdTime.UpdatedAt = tn

	purchaseOrderHeader := orderproto.PurchaseOrderHeader{PurchaseOrderHeaderD: &purchaseOrderHeaderD, PurchaseOrderHeaderT: &purchaseOrderHeaderT, CrUpdUser: &crUpdUser, CrUpdTime: &crUpdTime}

	purchaseOrderLines := []*orderproto.PurchaseOrderLine{}
	// we will do for loop on lines which is comes from client form
	for _, line := range in.PurchaseOrderLines {
		line.UserId = in.UserId
		line.UserEmail = in.UserEmail
		line.RequestId = in.RequestId
		// we wl call CreatePurchaseOrderLine function which wl populate form values to purchaseorderline struct
		purchaseOrderLine, err := ps.ProcessPurchaseOrderLineRequest(ctx, line)
		if err != nil {
			ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
			return nil, err
		}
		purchaseOrderLines = append(purchaseOrderLines, purchaseOrderLine)
	}

	err = ps.insertPurchaseOrderHeader(ctx, insertPurchaseOrderHeaderSQL, &purchaseOrderHeader, insertPurchaseOrderLineSQL, purchaseOrderLines, in.GetUserEmail(), in.GetRequestId())

	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	purchaseOrderHeaderResponse := orderproto.CreatePurchaseOrderHeaderResponse{}
	purchaseOrderHeaderResponse.PurchaseOrderHeader = &purchaseOrderHeader
	return &purchaseOrderHeaderResponse, nil
}

func (ps *PurchaseOrderHeaderService) insertPurchaseOrderHeader(ctx context.Context, insertPurchaseOrderHeaderSQL string, purchaseOrderHeader *orderproto.PurchaseOrderHeader, insertPurchaseOrderLineSQL string, purchaseOrderLines []*orderproto.PurchaseOrderLine, userEmail string, requestID string) error {
	purchaseOrderHeaderTmp, err := ps.CrPurchaseOrderHeaderStruct(ctx, purchaseOrderHeader, userEmail, requestID)
	if err != nil {
		ps.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
		return err
	}
	err = ps.DBService.InsUpd(ctx, userEmail, requestID, func(tx *sqlx.Tx) error {
		// header creation
		res, err := tx.NamedExecContext(ctx, insertPurchaseOrderHeaderSQL, purchaseOrderHeaderTmp)
		if err != nil {
			ps.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}

		uID, err := res.LastInsertId()
		if err != nil {
			ps.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}
		purchaseOrderHeader.PurchaseOrderHeaderD.Id = uint32(uID)
		uuid4Str, err := common.UUIDBytesToStr(purchaseOrderHeader.PurchaseOrderHeaderD.Uuid4)
		if err != nil {
			ps.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}
		purchaseOrderHeader.PurchaseOrderHeaderD.IdS = uuid4Str

		for _, purchaseOrderLine := range purchaseOrderLines {
			purchaseOrderLine.PurchaseOrderLineD.PurchaseOrderHeaderId = purchaseOrderHeader.PurchaseOrderHeaderD.Id
			purchaseOrderLineTmp, err := ps.CrPurchaseOrderLineStruct(ctx, purchaseOrderLine, userEmail, requestID)
			if err != nil {
				ps.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
				return err
			}
			_, err = tx.NamedExecContext(ctx, insertPurchaseOrderLineSQL, purchaseOrderLineTmp)
			if err != nil {
				ps.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
				return err
			}
		}
		return nil
	})
	if err != nil {
		ps.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
		return err
	}
	return nil
}

// CrPurchaseOrderHeaderStruct - process PurchaseOrderHeader details
func (ps *PurchaseOrderHeaderService) CrPurchaseOrderHeaderStruct(ctx context.Context, purchaseOrderHeader *orderproto.PurchaseOrderHeader, userEmail string, requestID string) (*orderstruct.PurchaseOrderHeader, error) {
	purchaseOrderHeaderT := new(orderstruct.PurchaseOrderHeaderT)
	purchaseOrderHeaderT.IssueDate = common.TimestampToTime(purchaseOrderHeader.PurchaseOrderHeaderT.IssueDate)
	purchaseOrderHeaderT.ValidityPeriod = common.TimestampToTime(purchaseOrderHeader.PurchaseOrderHeaderT.ValidityPeriod)
	purchaseOrderHeaderT.TaxExDate = common.TimestampToTime(purchaseOrderHeader.PurchaseOrderHeaderT.TaxExDate)
	purchaseOrderHeaderT.PaymentExDate = common.TimestampToTime(purchaseOrderHeader.PurchaseOrderHeaderT.PaymentExDate)
	purchaseOrderHeaderT.PricingExDate = common.TimestampToTime(purchaseOrderHeader.PurchaseOrderHeaderT.PricingExDate)

	crUpdTime := new(commonstruct.CrUpdTime)
	crUpdTime.CreatedAt = common.TimestampToTime(purchaseOrderHeader.CrUpdTime.CreatedAt)
	crUpdTime.UpdatedAt = common.TimestampToTime(purchaseOrderHeader.CrUpdTime.UpdatedAt)

	purchaseOrderHeaderTmp := orderstruct.PurchaseOrderHeader{PurchaseOrderHeaderD: purchaseOrderHeader.PurchaseOrderHeaderD, PurchaseOrderHeaderT: purchaseOrderHeaderT, CrUpdUser: purchaseOrderHeader.CrUpdUser, CrUpdTime: crUpdTime}

	return &purchaseOrderHeaderTmp, nil
}

// GetPurchaseOrderHeaders - Get PurchaseOrderHeaders
func (ps *PurchaseOrderHeaderService) GetPurchaseOrderHeaders(ctx context.Context, in *orderproto.GetPurchaseOrderHeadersRequest) (*orderproto.GetPurchaseOrderHeadersResponse, error) {
	limit := in.GetLimit()
	nextCursor := in.GetNextCursor()
	if limit == "" {
		limit = ps.DBService.LimitSQLRows
	}
	query := "status_code = ?"
	if nextCursor == "" {
		query = query + " order by id desc " + " limit " + limit + ";"
	} else {
		nextCursor = common.DecodeCursor(nextCursor)
		query = query + " " + "and" + " " + "id <= " + nextCursor + " order by id desc " + " limit " + limit + ";"
	}

	purchaseOrderHeaders := []*orderproto.PurchaseOrderHeader{}

	nselectPurchaseOrderHeadersSQL := selectPurchaseOrderHeadersSQL + ` where ` + query

	rows, err := ps.DBService.DB.QueryxContext(ctx, nselectPurchaseOrderHeadersSQL, "active")
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	for rows.Next() {

		purchaseOrderHeaderTmp := orderstruct.PurchaseOrderHeader{}
		err = rows.StructScan(&purchaseOrderHeaderTmp)
		if err != nil {
			ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
			return nil, err
		}

		getRequest := commonproto.GetRequest{}
		getRequest.UserEmail = in.UserEmail
		getRequest.RequestId = in.RequestId
		purchaseOrderHeader, err := ps.getPurchaseOrderHeaderStruct(ctx, &getRequest, &purchaseOrderHeaderTmp)
		if err != nil {
			ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
			return nil, err
		}
		purchaseOrderHeaders = append(purchaseOrderHeaders, purchaseOrderHeader)

	}

	purchaseOrderHeadersResponse := orderproto.GetPurchaseOrderHeadersResponse{}
	if len(purchaseOrderHeaders) != 0 {
		next := purchaseOrderHeaders[len(purchaseOrderHeaders)-1].PurchaseOrderHeaderD.Id
		next--
		nextc := common.EncodeCursor(next)
		purchaseOrderHeadersResponse = orderproto.GetPurchaseOrderHeadersResponse{PurchaseOrderHeaders: purchaseOrderHeaders, NextCursor: nextc}
	} else {
		purchaseOrderHeadersResponse = orderproto.GetPurchaseOrderHeadersResponse{PurchaseOrderHeaders: purchaseOrderHeaders, NextCursor: "0"}
	}
	return &purchaseOrderHeadersResponse, nil
}

// GetPurchaseOrderHeader - Get PurchaseOrderHeader
func (ps *PurchaseOrderHeaderService) GetPurchaseOrderHeader(ctx context.Context, inReq *orderproto.GetPurchaseOrderHeaderRequest) (*orderproto.GetPurchaseOrderHeaderResponse, error) {
	in := inReq.GetRequest
	uuid4byte, err := common.UUIDStrToBytes(in.Id)
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	nselectPurchaseOrderHeadersSQL := selectPurchaseOrderHeadersSQL + ` where uuid4 = ? and status_code = ?;`
	row := ps.DBService.DB.QueryRowxContext(ctx, nselectPurchaseOrderHeadersSQL, uuid4byte, "active")
	purchaseOrderHeaderTmp := orderstruct.PurchaseOrderHeader{}
	err = row.StructScan(&purchaseOrderHeaderTmp)
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	purchaseOrderHeader, err := ps.getPurchaseOrderHeaderStruct(ctx, in, &purchaseOrderHeaderTmp)
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	purchaseOrderHeaderResponse := orderproto.GetPurchaseOrderHeaderResponse{}
	purchaseOrderHeaderResponse.PurchaseOrderHeader = purchaseOrderHeader
	return &purchaseOrderHeaderResponse, nil
}

// GetPurchaseOrderHeaderByPk - Get PurchaseOrderHeader By Primary key(Id)
func (ps *PurchaseOrderHeaderService) GetPurchaseOrderHeaderByPk(ctx context.Context, inReq *orderproto.GetPurchaseOrderHeaderByPkRequest) (*orderproto.GetPurchaseOrderHeaderByPkResponse, error) {
	in := inReq.GetByIdRequest
	nselectPurchaseOrderHeadersSQL := selectPurchaseOrderHeadersSQL + ` where id = ? and status_code = ?;`
	row := ps.DBService.DB.QueryRowxContext(ctx, nselectPurchaseOrderHeadersSQL, in.Id, "active")
	purchaseOrderHeaderTmp := orderstruct.PurchaseOrderHeader{}
	err := row.StructScan(&purchaseOrderHeaderTmp)
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	getRequest := commonproto.GetRequest{}
	getRequest.UserEmail = in.UserEmail
	getRequest.RequestId = in.RequestId
	purchaseOrderHeader, err := ps.getPurchaseOrderHeaderStruct(ctx, &getRequest, &purchaseOrderHeaderTmp)
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	purchaseOrderHeaderResponse := orderproto.GetPurchaseOrderHeaderByPkResponse{}
	purchaseOrderHeaderResponse.PurchaseOrderHeader = purchaseOrderHeader
	return &purchaseOrderHeaderResponse, nil
}

// getPurchaseOrderHeaderStruct - Get PurchaseOrderHeader
func (ps *PurchaseOrderHeaderService) getPurchaseOrderHeaderStruct(ctx context.Context, in *commonproto.GetRequest, purchaseOrderHeaderTmp *orderstruct.PurchaseOrderHeader) (*orderproto.PurchaseOrderHeader, error) {
	purchaseOrderHeaderT := new(orderproto.PurchaseOrderHeaderT)
	purchaseOrderHeaderT.IssueDate = common.TimeToTimestamp(purchaseOrderHeaderTmp.PurchaseOrderHeaderT.IssueDate)
	purchaseOrderHeaderT.ValidityPeriod = common.TimeToTimestamp(purchaseOrderHeaderTmp.PurchaseOrderHeaderT.ValidityPeriod)
	purchaseOrderHeaderT.TaxExDate = common.TimeToTimestamp(purchaseOrderHeaderTmp.PurchaseOrderHeaderT.TaxExDate)
	purchaseOrderHeaderT.PricingExDate = common.TimeToTimestamp(purchaseOrderHeaderTmp.PurchaseOrderHeaderT.PricingExDate)
	purchaseOrderHeaderT.PaymentExDate = common.TimeToTimestamp(purchaseOrderHeaderTmp.PurchaseOrderHeaderT.PaymentExDate)

	crUpdTime := new(commonproto.CrUpdTime)
	crUpdTime.CreatedAt = common.TimeToTimestamp(purchaseOrderHeaderTmp.CrUpdTime.CreatedAt)
	crUpdTime.UpdatedAt = common.TimeToTimestamp(purchaseOrderHeaderTmp.CrUpdTime.UpdatedAt)

	uuid4Str, err := common.UUIDBytesToStr(purchaseOrderHeaderTmp.PurchaseOrderHeaderD.Uuid4)
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	purchaseOrderHeaderTmp.PurchaseOrderHeaderD.IdS = uuid4Str

	purchaseOrderHeader := orderproto.PurchaseOrderHeader{PurchaseOrderHeaderD: purchaseOrderHeaderTmp.PurchaseOrderHeaderD, PurchaseOrderHeaderT: purchaseOrderHeaderT, CrUpdUser: purchaseOrderHeaderTmp.CrUpdUser, CrUpdTime: crUpdTime}

	return &purchaseOrderHeader, nil
}

// UpdatePurchaseOrderHeader - Update PurchaseOrderHeader
func (ps *PurchaseOrderHeaderService) UpdatePurchaseOrderHeader(ctx context.Context, in *orderproto.UpdatePurchaseOrderHeaderRequest) (*orderproto.UpdatePurchaseOrderHeaderResponse, error) {
	db := ps.DBService.DB
	tn := common.GetTimeDetails()

	uuid4byte, err := common.UUIDStrToBytes(in.Id)
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	stmt, err := db.PreparexContext(ctx, updatePurchaseOrderHeaderSQL)
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	err = ps.DBService.InsUpd(ctx, in.GetUserEmail(), in.GetRequestId(), func(tx *sqlx.Tx) error {
		_, err = tx.StmtxContext(ctx, stmt).ExecContext(ctx,
			in.OrderTypeCode,
			in.Note,
			in.RequestedInvoiceCurrencyCode,
			in.DocumentCurrencyCode,
			in.PricingCurrencyCode,
			in.TaxCurrencyCode,
			in.AccountingCostCode,
			in.AccountingCost,
			tn,
			uuid4byte)
		if err != nil {
			ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
			err1 := stmt.Close()
			if err1 != nil {
				ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err1))
				return err1
			}
			return err
		}
		return nil
	})

	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	err = stmt.Close()
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	return &orderproto.UpdatePurchaseOrderHeaderResponse{}, nil
}
