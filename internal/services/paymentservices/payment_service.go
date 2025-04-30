package paymentservices

import (
	"context"
	"net"
	"os"
	"time"

	"github.com/cloudfresco/sc-ubl/internal/common"
	"github.com/cloudfresco/sc-ubl/internal/config"
	commonproto "github.com/cloudfresco/sc-ubl/internal/protogen/common/v1"
	partyproto "github.com/cloudfresco/sc-ubl/internal/protogen/party/v1"
	paymentproto "github.com/cloudfresco/sc-ubl/internal/protogen/payment/v1"
	partyservice "github.com/cloudfresco/sc-ubl/internal/services/partyservices"
	commonstruct "github.com/cloudfresco/sc-ubl/internal/servicestructs/common/v1"
	paymentstruct "github.com/cloudfresco/sc-ubl/internal/servicestructs/payment/v1"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// PaymentService - For accessing Payment services
type PaymentService struct {
	log               *zap.Logger
	DBService         *common.DBService
	RedisService      *common.RedisService
	UserServiceClient partyproto.UserServiceClient
	paymentproto.UnimplementedPaymentServiceServer
}

// NewPaymentService - Create Payment service
func NewPaymentService(log *zap.Logger, dbOpt *common.DBService, redisOpt *common.RedisService, userServiceClient partyproto.UserServiceClient) *PaymentService {
	return &PaymentService{
		log:               log,
		DBService:         dbOpt,
		RedisService:      redisOpt,
		UserServiceClient: userServiceClient,
	}
}

// StartPaymentServer - Start Payment server
func StartPaymentServer(log *zap.Logger, isTest bool, pwd string, dbOpt *config.DBOptions, redisOpt *config.RedisOptions, mailerOpt *config.MailerOptions, grpcServerOpt *config.GrpcServerOptions, jwtOpt *config.JWTOptions, oauthOpt *config.OauthOptions, userOpt *config.UserOptions, uptraceOpt *config.UptraceOptions, dbService *common.DBService, redisService *common.RedisService, mailerService common.MailerIntf) {
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
	paymentService := NewPaymentService(log, dbService, redisService, uc)

	lis, err := net.Listen("tcp", grpcServerOpt.GrpcPaymentServerPort)
	if err != nil {
		log.Error("Error", zap.Error(err))
		os.Exit(1)
	}

	srv := grpc.NewServer(srvOpts...)
	paymentproto.RegisterPaymentServiceServer(srv, paymentService)

	if err := srv.Serve(lis); err != nil {
		log.Error("Error", zap.Error(err))
		os.Exit(1)
	}
}

const insertPaymentSQL = `insert into  payments
	  (
    uuid4,
    p_id,
    paid_amount,
    instruction_id,
    payment_mean_id,
    received_date,
    paid_date,
    status_code,
    created_by_user_id,
    updated_by_user_id,
    created_at,
    updated_at)
  values (:uuid4,
:p_id,
:paid_amount,
:instruction_id,
:payment_mean_id,
:received_date,
:paid_date,
:status_code,
:created_by_user_id,
:updated_by_user_id,
:created_at,
:updated_at);`

const insertPaymentMandateSQL = `insert into payment_mandates
	  (
    uuid4,
    pmd_id,
    mandate_type_code,
    maximum_payment_instructions_numeric,
    maximum_paid_amount,
    signature_id,
    payer_party_id,
    payer_financial_account_id,
    clause,
    validity_period_start_date,
    validity_period_end_date,
    payment_reversal_period_start_date,
    payment_reversal_period_end_date,
    status_code,
    created_by_user_id,
    updated_by_user_id,
    created_at,
    updated_at)
  values (:uuid4,
:pmd_id,
:mandate_type_code,
:maximum_payment_instructions_numeric,
:maximum_paid_amount,
:signature_id,
:payer_party_id,
:payer_financial_account_id,
:clause,
:validity_period_start_date,
:validity_period_end_date,
:payment_reversal_period_start_date,
:payment_reversal_period_end_date,
:status_code,
:created_by_user_id,
:updated_by_user_id,
:created_at,
:updated_at);`

const insertPaymentMeanSQL = `insert into payment_means
	  (
    uuid4,
    pm_id,
    payment_means_code,
    payment_channel_code,
    instruction_id,
    instruction_note,
    credit_account_id,
    payment_term_id,
    payment_mandate_id,
    trade_financing_id,
    payer_financial_account_id,
    payee_financial_account_id,
    payment_due_date,
    status_code,
    created_by_user_id,
    updated_by_user_id,
    created_at,
    updated_at)
  values (:uuid4,
:pm_id,
:payment_means_code,
:payment_channel_code,
:instruction_id,
:instruction_note,
:credit_account_id,
:payment_term_id,
:payment_mandate_id,
:trade_financing_id,
:payer_financial_account_id,
:payee_financial_account_id,
:payment_due_date,
:status_code,
:created_by_user_id,
:updated_by_user_id,
:created_at,
:updated_at);`

const insertPaymentTermSQL = `insert into payment_terms
	  (
      uuid4,
      pt_id,
      prepaid_payment_reference_id,
      note,
      reference_event_code,
      settlement_discount_percent,
      penalty_surcharge_percent,
      payment_percent,
      amount,
      settlement_discount_amount,
      penalty_amount,
      payment_terms_details_uri,
      exchange_rate,
      payment_due_date,
      installment_due_date,
      settlement_period_start_date,
      settlement_period_end_date,
      penalty_period_start_date,
      penalty_period_end_date,
      validity_period_start_date,
      validity_period_end_date,
      status_code,
      created_by_user_id,
      updated_by_user_id,
      created_at,
      updated_at)
  values (:uuid4,
:pt_id,
:prepaid_payment_reference_id,
:note,
:reference_event_code,
:settlement_discount_percent,
:penalty_surcharge_percent,
:payment_percent,
:amount,
:settlement_discount_amount,
:penalty_amount,
:payment_terms_details_uri,
:exchange_rate,
:payment_due_date,
:installment_due_date,
:settlement_period_start_date,
:settlement_period_end_date,
:penalty_period_start_date,
:penalty_period_end_date,
:validity_period_start_date,
:validity_period_end_date,
:status_code,
:created_by_user_id,
:updated_by_user_id,
:created_at,
:updated_at);`

const insertPaymentMandateClauseSQL = `insert into payment_mandate_clauses
	  (
    pm_cl_id,
    payment_mandate_id)
  values (:pm_cl_id,
:payment_mandate_id);`

const insertPaymentMandateClauseContentSQL = `insert into payment_mandate_clause_contents
	  (
    content,
    payment_mandate_clause_id,
    status_code,
    created_by_user_id,
    updated_by_user_id,
    created_at,
    updated_at)
  values (:content,
:payment_mandate_clause_id,
:status_code,
:created_by_user_id,
:updated_by_user_id,
:created_at,
:updated_at);`

// CreatePayment - Create Payment
func (ps *PaymentService) CreatePayment(ctx context.Context, in *paymentproto.CreatePaymentRequest) (*paymentproto.CreatePaymentResponse, error) {
	user, err := partyservice.GetUserWithNewContext(ctx, in.UserId, in.UserEmail, in.RequestId, ps.UserServiceClient)
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	ttime := common.GetTimeDetails()
	tn := common.TimeToTimestamp(ttime)

	receivedDate, err := time.Parse(common.Layout, in.ReceivedDate)
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	paidDate, err := time.Parse(common.Layout, in.PaidDate)
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	paymentD := paymentproto.PaymentD{}
	paymentD.Uuid4, err = common.GetUUIDBytes()
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	paymentD.PId = in.PId
	paymentD.PaidAmount = in.PaidAmount
	paymentD.InstructionId = in.InstructionId
	paymentD.PaymentMeanId = in.PaymentMeanId

	paymentT := paymentproto.PaymentT{}
	paymentT.ReceivedDate = common.TimeToTimestamp(receivedDate.UTC().Truncate(time.Second))
	paymentT.PaidDate = common.TimeToTimestamp(paidDate.UTC().Truncate(time.Second))

	crUpdUser := commonproto.CrUpdUser{}
	crUpdUser.StatusCode = "active"
	crUpdUser.CreatedByUserId = user.Id
	crUpdUser.UpdatedByUserId = user.Id

	crUpdTime := commonproto.CrUpdTime{}
	crUpdTime.CreatedAt = tn
	crUpdTime.UpdatedAt = tn

	payment := paymentproto.Payment{PaymentD: &paymentD, PaymentT: &paymentT, CrUpdUser: &crUpdUser, CrUpdTime: &crUpdTime}

	err = ps.insertPayment(ctx, insertPaymentSQL, &payment, in.GetUserEmail(), in.GetRequestId())
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	paymentResponse := paymentproto.CreatePaymentResponse{}
	paymentResponse.Payment = &payment
	return &paymentResponse, nil
}

// insertPayment - Insert payment details into database
func (ps *PaymentService) insertPayment(ctx context.Context, insertPaymentSQL string, payment *paymentproto.Payment, userEmail string, requestID string) error {
	paymentTmp, err := ps.crPaymentStruct(ctx, payment, userEmail, requestID)
	if err != nil {
		ps.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
		return err
	}

	err = ps.DBService.InsUpd(ctx, userEmail, requestID, func(tx *sqlx.Tx) error {
		res, err := tx.NamedExecContext(ctx, insertPaymentSQL, paymentTmp)
		if err != nil {
			ps.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}

		uID, err := res.LastInsertId()
		if err != nil {
			ps.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}
		payment.PaymentD.Id = uint32(uID)
		uuid4Str, err := common.UUIDBytesToStr(payment.PaymentD.Uuid4)
		if err != nil {
			ps.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}
		payment.PaymentD.IdS = uuid4Str
		return nil
	})

	if err != nil {
		ps.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
		return err
	}
	return nil
}

// crPaymentStruct - process Payment details
func (ps *PaymentService) crPaymentStruct(ctx context.Context, payment *paymentproto.Payment, userEmail string, requestID string) (*paymentstruct.Payment, error) {
	paymentT := new(paymentstruct.PaymentT)
	paymentT.ReceivedDate = common.TimestampToTime(payment.PaymentT.ReceivedDate)
	paymentT.PaidDate = common.TimestampToTime(payment.PaymentT.PaidDate)

	crUpdTime := new(commonstruct.CrUpdTime)
	crUpdTime.CreatedAt = common.TimestampToTime(payment.CrUpdTime.CreatedAt)
	crUpdTime.UpdatedAt = common.TimestampToTime(payment.CrUpdTime.UpdatedAt)

	paymentTmp := paymentstruct.Payment{PaymentD: payment.PaymentD, PaymentT: paymentT, CrUpdUser: payment.CrUpdUser, CrUpdTime: crUpdTime}
	return &paymentTmp, nil
}

// CreatePaymentTerm - Create PaymentTerm
func (ps *PaymentService) CreatePaymentTerm(ctx context.Context, in *paymentproto.CreatePaymentTermRequest) (*paymentproto.CreatePaymentTermResponse, error) {
	user, err := partyservice.GetUserWithNewContext(ctx, in.UserId, in.UserEmail, in.RequestId, ps.UserServiceClient)
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	ttime := common.GetTimeDetails()
	tn := common.TimeToTimestamp(ttime)

	paymentDueDate, err := time.Parse(common.Layout, in.PaymentDueDate)
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	installmentDueDate, err := time.Parse(common.Layout, in.InstallmentDueDate)
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	settlementPeriodStartDate, err := time.Parse(common.Layout, in.SettlementPeriodStartDate)
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	settlementPeriodEndDate, err := time.Parse(common.Layout, in.SettlementPeriodEndDate)
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	penaltyPeriodStartDate, err := time.Parse(common.Layout, in.PenaltyPeriodStartDate)
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	penaltyPeriodEndDate, err := time.Parse(common.Layout, in.PenaltyPeriodEndDate)
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	validityPeriodStartDate, err := time.Parse(common.Layout, in.ValidityPeriodStartDate)
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	validityPeriodEndDate, err := time.Parse(common.Layout, in.ValidityPeriodEndDate)
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	paymentTermD := paymentproto.PaymentTermD{}
	paymentTermD.Uuid4, err = common.GetUUIDBytes()
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	paymentTermD.PtId = in.PtId
	paymentTermD.PrepaidPaymentReferenceId = in.PrepaidPaymentReferenceId
	paymentTermD.Note = in.Note
	paymentTermD.ReferenceEventCode = in.ReferenceEventCode
	paymentTermD.SettlementDiscountPercent = in.SettlementDiscountPercent
	paymentTermD.PenaltySurchargePercent = in.PenaltySurchargePercent
	paymentTermD.PaymentPercent = in.PaymentPercent
	paymentTermD.Amount = in.Amount
	paymentTermD.SettlementDiscountAmount = in.SettlementDiscountAmount
	paymentTermD.PenaltyAmount = in.PenaltyAmount
	paymentTermD.PaymentTermsDetailsURI = in.PaymentTermsDetailsURI
	paymentTermD.ExchangeRate = in.ExchangeRate

	paymentTermT := paymentproto.PaymentTermT{}
	paymentTermT.PaymentDueDate = common.TimeToTimestamp(paymentDueDate.UTC().Truncate(time.Second))
	paymentTermT.InstallmentDueDate = common.TimeToTimestamp(installmentDueDate.UTC().Truncate(time.Second))
	paymentTermT.SettlementPeriodStartDate = common.TimeToTimestamp(settlementPeriodStartDate.UTC().Truncate(time.Second))
	paymentTermT.SettlementPeriodEndDate = common.TimeToTimestamp(settlementPeriodEndDate.UTC().Truncate(time.Second))
	paymentTermT.PenaltyPeriodStartDate = common.TimeToTimestamp(penaltyPeriodStartDate.UTC().Truncate(time.Second))
	paymentTermT.PenaltyPeriodEndDate = common.TimeToTimestamp(penaltyPeriodEndDate.UTC().Truncate(time.Second))
	paymentTermT.ValidityPeriodStartDate = common.TimeToTimestamp(validityPeriodStartDate.UTC().Truncate(time.Second))
	paymentTermT.ValidityPeriodEndDate = common.TimeToTimestamp(validityPeriodEndDate.UTC().Truncate(time.Second))

	crUpdUser := commonproto.CrUpdUser{}
	crUpdUser.StatusCode = "active"
	crUpdUser.CreatedByUserId = user.Id
	crUpdUser.UpdatedByUserId = user.Id

	crUpdTime := commonproto.CrUpdTime{}
	crUpdTime.CreatedAt = tn
	crUpdTime.UpdatedAt = tn

	paymentTerm := paymentproto.PaymentTerm{PaymentTermD: &paymentTermD, PaymentTermT: &paymentTermT, CrUpdUser: &crUpdUser, CrUpdTime: &crUpdTime}

	err = ps.insertPaymentTerm(ctx, insertPaymentTermSQL, &paymentTerm, in.GetUserEmail(), in.GetRequestId())
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))

		return nil, err
	}

	paymentTermResponse := paymentproto.CreatePaymentTermResponse{}
	paymentTermResponse.PaymentTerm = &paymentTerm
	return &paymentTermResponse, nil
}

// insertPaymentTerm - Insert PaymentTerm details into database
func (ps *PaymentService) insertPaymentTerm(ctx context.Context, insertPaymentTermSQL string, paymentTerm *paymentproto.PaymentTerm, userEmail string, requestID string) error {
	paymentTermTmp, err := ps.crPaymentTermStruct(ctx, paymentTerm, userEmail, requestID)
	if err != nil {
		ps.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
		return err
	}

	err = ps.DBService.InsUpd(ctx, userEmail, requestID, func(tx *sqlx.Tx) error {
		res, err := tx.NamedExecContext(ctx, insertPaymentTermSQL, paymentTermTmp)
		if err != nil {
			ps.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}

		uID, err := res.LastInsertId()
		if err != nil {
			ps.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}
		paymentTerm.PaymentTermD.Id = uint32(uID)
		uuid4Str, err := common.UUIDBytesToStr(paymentTerm.PaymentTermD.Uuid4)
		if err != nil {
			ps.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}
		paymentTerm.PaymentTermD.IdS = uuid4Str

		return nil
	})

	if err != nil {
		ps.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
		return err
	}
	return nil
}

// crPaymentTermStruct - process PaymentTerm details
func (ps *PaymentService) crPaymentTermStruct(ctx context.Context, paymentTerm *paymentproto.PaymentTerm, userEmail string, requestID string) (*paymentstruct.PaymentTerm, error) {
	paymentTermT := new(paymentstruct.PaymentTermT)
	paymentTermT.PaymentDueDate = common.TimestampToTime(paymentTerm.PaymentTermT.PaymentDueDate)
	paymentTermT.InstallmentDueDate = common.TimestampToTime(paymentTerm.PaymentTermT.InstallmentDueDate)
	paymentTermT.SettlementPeriodStartDate = common.TimestampToTime(paymentTerm.PaymentTermT.SettlementPeriodStartDate)
	paymentTermT.SettlementPeriodEndDate = common.TimestampToTime(paymentTerm.PaymentTermT.SettlementPeriodEndDate)
	paymentTermT.PenaltyPeriodStartDate = common.TimestampToTime(paymentTerm.PaymentTermT.PenaltyPeriodStartDate)
	paymentTermT.PenaltyPeriodEndDate = common.TimestampToTime(paymentTerm.PaymentTermT.PenaltyPeriodEndDate)
	paymentTermT.ValidityPeriodStartDate = common.TimestampToTime(paymentTerm.PaymentTermT.ValidityPeriodStartDate)
	paymentTermT.ValidityPeriodEndDate = common.TimestampToTime(paymentTerm.PaymentTermT.ValidityPeriodEndDate)

	crUpdTime := new(commonstruct.CrUpdTime)
	crUpdTime.CreatedAt = common.TimestampToTime(paymentTerm.CrUpdTime.CreatedAt)
	crUpdTime.UpdatedAt = common.TimestampToTime(paymentTerm.CrUpdTime.UpdatedAt)

	paymentTermTmp := paymentstruct.PaymentTerm{PaymentTermD: paymentTerm.PaymentTermD, PaymentTermT: paymentTermT, CrUpdUser: paymentTerm.CrUpdUser, CrUpdTime: crUpdTime}

	return &paymentTermTmp, nil
}

// CreatePaymentMandate - Create PaymentMandate
func (ps *PaymentService) CreatePaymentMandate(ctx context.Context, in *paymentproto.CreatePaymentMandateRequest) (*paymentproto.CreatePaymentMandateResponse, error) {
	user, err := partyservice.GetUserWithNewContext(ctx, in.UserId, in.UserEmail, in.RequestId, ps.UserServiceClient)
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	ttime := common.GetTimeDetails()
	tn := common.TimeToTimestamp(ttime)

	validityPeriodStartDate, err := time.Parse(common.Layout, in.ValidityPeriodStartDate)
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	validityPeriodEndDate, err := time.Parse(common.Layout, in.ValidityPeriodEndDate)
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	paymentReversalPeriodStartDate, err := time.Parse(common.Layout, in.PaymentReversalPeriodStartDate)
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	paymentReversalPeriodEndDate, err := time.Parse(common.Layout, in.PaymentReversalPeriodEndDate)
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	paymentMandateD := paymentproto.PaymentMandateD{}
	paymentMandateD.Uuid4, err = common.GetUUIDBytes()
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	paymentMandateD.PmdId = in.PmdId
	paymentMandateD.MandateTypeCode = in.MandateTypeCode
	paymentMandateD.MaximumPaymentInstructionsNumeric = in.MaximumPaymentInstructionsNumeric
	paymentMandateD.MaximumPaidAmount = in.MaximumPaidAmount
	paymentMandateD.SignatureId = in.SignatureId
	paymentMandateD.PayerPartyId = in.PayerPartyId
	paymentMandateD.PayerFinancialAccountId = in.PayerPartyId
	paymentMandateD.Clause = in.Clause

	paymentMandateT := paymentproto.PaymentMandateT{}
	paymentMandateT.ValidityPeriodStartDate = common.TimeToTimestamp(validityPeriodStartDate.UTC().Truncate(time.Second))
	paymentMandateT.ValidityPeriodEndDate = common.TimeToTimestamp(validityPeriodEndDate.UTC().Truncate(time.Second))
	paymentMandateT.PaymentReversalPeriodStartDate = common.TimeToTimestamp(paymentReversalPeriodStartDate.UTC().Truncate(time.Second))
	paymentMandateT.PaymentReversalPeriodEndDate = common.TimeToTimestamp(paymentReversalPeriodEndDate.UTC().Truncate(time.Second))

	crUpdUser := commonproto.CrUpdUser{}
	crUpdUser.StatusCode = "active"
	crUpdUser.CreatedByUserId = user.Id
	crUpdUser.UpdatedByUserId = user.Id

	crUpdTime := commonproto.CrUpdTime{}
	crUpdTime.CreatedAt = tn
	crUpdTime.UpdatedAt = tn

	paymentMandate := paymentproto.PaymentMandate{PaymentMandateD: &paymentMandateD, PaymentMandateT: &paymentMandateT, CrUpdUser: &crUpdUser, CrUpdTime: &crUpdTime}

	err = ps.insertPaymentMandate(ctx, insertPaymentMandateSQL, &paymentMandate, in.GetUserEmail(), in.GetRequestId())
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))

		return nil, err
	}
	paymentMandateResponse := paymentproto.CreatePaymentMandateResponse{}
	paymentMandateResponse.PaymentMandate = &paymentMandate
	return &paymentMandateResponse, nil
}

// insertPaymentMandate - Insert payment details into database
func (ps *PaymentService) insertPaymentMandate(ctx context.Context, insertPaymentMandateSQL string, paymentMandate *paymentproto.PaymentMandate, userEmail string, requestID string) error {
	paymentMandateTmp, err := ps.crPaymentMandateStruct(ctx, paymentMandate, userEmail, requestID)
	if err != nil {
		ps.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
		return err
	}

	err = ps.DBService.InsUpd(ctx, userEmail, requestID, func(tx *sqlx.Tx) error {
		res, err := tx.NamedExecContext(ctx, insertPaymentMandateSQL, paymentMandateTmp)
		if err != nil {
			ps.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}

		uID, err := res.LastInsertId()
		if err != nil {
			ps.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}
		paymentMandate.PaymentMandateD.Id = uint32(uID)
		uuid4Str, err := common.UUIDBytesToStr(paymentMandate.PaymentMandateD.Uuid4)
		if err != nil {
			ps.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}
		paymentMandate.PaymentMandateD.IdS = uuid4Str
		return nil
	})

	if err != nil {
		ps.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
		return err
	}
	return nil
}

// crPaymentMandateStruct - process PaymentMandate details
func (ps *PaymentService) crPaymentMandateStruct(ctx context.Context, paymentMandate *paymentproto.PaymentMandate, userEmail string, requestID string) (*paymentstruct.PaymentMandate, error) {
	paymentMandateT := new(paymentstruct.PaymentMandateT)
	paymentMandateT.ValidityPeriodStartDate = common.TimestampToTime(paymentMandate.PaymentMandateT.ValidityPeriodStartDate)
	paymentMandateT.ValidityPeriodEndDate = common.TimestampToTime(paymentMandate.PaymentMandateT.ValidityPeriodEndDate)
	paymentMandateT.PaymentReversalPeriodStartDate = common.TimestampToTime(paymentMandate.PaymentMandateT.PaymentReversalPeriodStartDate)
	paymentMandateT.PaymentReversalPeriodEndDate = common.TimestampToTime(paymentMandate.PaymentMandateT.PaymentReversalPeriodEndDate)

	crUpdTime := new(commonstruct.CrUpdTime)
	crUpdTime.CreatedAt = common.TimestampToTime(paymentMandate.CrUpdTime.CreatedAt)
	crUpdTime.UpdatedAt = common.TimestampToTime(paymentMandate.CrUpdTime.UpdatedAt)

	paymentMandateTmp := paymentstruct.PaymentMandate{PaymentMandateD: paymentMandate.PaymentMandateD, PaymentMandateT: paymentMandateT, CrUpdUser: paymentMandate.CrUpdUser, CrUpdTime: crUpdTime}

	return &paymentMandateTmp, nil
}

// CreatePaymentMean - Create PaymentMean
func (ps *PaymentService) CreatePaymentMean(ctx context.Context, in *paymentproto.CreatePaymentMeanRequest) (*paymentproto.CreatePaymentMeanResponse, error) {
	user, err := partyservice.GetUserWithNewContext(ctx, in.UserId, in.UserEmail, in.RequestId, ps.UserServiceClient)
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	ttime := common.GetTimeDetails()
	tn := common.TimeToTimestamp(ttime)

	paymentDueDate, err := time.Parse(common.Layout, in.PaymentDueDate)
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	paymentMeanD := paymentproto.PaymentMeanD{}
	paymentMeanD.Uuid4, err = common.GetUUIDBytes()
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	paymentMeanD.PmId = in.PmId
	paymentMeanD.PaymentMeansCode = in.PaymentMeansCode
	paymentMeanD.PaymentChannelCode = in.PaymentChannelCode
	paymentMeanD.InstructionId = in.InstructionId
	paymentMeanD.InstructionNote = in.InstructionNote
	paymentMeanD.CreditAccountId = in.CreditAccountId
	paymentMeanD.PaymentTermId = in.PaymentTermId
	paymentMeanD.PaymentMandateId = in.PaymentMandateId
	paymentMeanD.TradeFinancingId = in.TradeFinancingId
	paymentMeanD.PayerFinancialAccountId = in.PayerFinancialAccountId
	paymentMeanD.PayeeFinancialAccountId = in.PayeeFinancialAccountId

	paymentMeanT := paymentproto.PaymentMeanT{}
	paymentMeanT.PaymentDueDate = common.TimeToTimestamp(paymentDueDate.UTC().Truncate(time.Second))

	crUpdUser := commonproto.CrUpdUser{}
	crUpdUser.StatusCode = "active"
	crUpdUser.CreatedByUserId = user.Id
	crUpdUser.UpdatedByUserId = user.Id

	crUpdTime := commonproto.CrUpdTime{}
	crUpdTime.CreatedAt = tn
	crUpdTime.UpdatedAt = tn

	paymentMean := paymentproto.PaymentMean{PaymentMeanD: &paymentMeanD, PaymentMeanT: &paymentMeanT, CrUpdUser: &crUpdUser, CrUpdTime: &crUpdTime}

	err = ps.insertPaymentMean(ctx, insertPaymentMeanSQL, &paymentMean, in.GetUserEmail(), in.GetRequestId())
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))

		return nil, err
	}

	paymentMeanResponse := paymentproto.CreatePaymentMeanResponse{}
	paymentMeanResponse.PaymentMean = &paymentMean
	return &paymentMeanResponse, nil
}

// insertPaymentMean - Insert payment details into database
func (ps *PaymentService) insertPaymentMean(ctx context.Context, insertPaymentMeanSQL string, paymentMean *paymentproto.PaymentMean, userEmail string, requestID string) error {
	paymentMeanTmp, err := ps.crPaymentMeanStruct(ctx, paymentMean, userEmail, requestID)
	if err != nil {
		ps.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
		return err
	}

	err = ps.DBService.InsUpd(ctx, userEmail, requestID, func(tx *sqlx.Tx) error {
		res, err := tx.NamedExecContext(ctx, insertPaymentMeanSQL, paymentMeanTmp)
		if err != nil {
			ps.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}

		uID, err := res.LastInsertId()
		if err != nil {
			ps.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}
		paymentMean.PaymentMeanD.Id = uint32(uID)
		uuid4Str, err := common.UUIDBytesToStr(paymentMean.PaymentMeanD.Uuid4)
		if err != nil {
			ps.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}
		paymentMean.PaymentMeanD.IdS = uuid4Str
		return nil
	})

	if err != nil {
		ps.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
		return err
	}
	return nil
}

// crPaymentMeanStruct - process PaymentMean details
func (ps *PaymentService) crPaymentMeanStruct(ctx context.Context, paymentMean *paymentproto.PaymentMean, userEmail string, requestID string) (*paymentstruct.PaymentMean, error) {
	paymentMeanT := new(paymentstruct.PaymentMeanT)
	paymentMeanT.PaymentDueDate = common.TimestampToTime(paymentMean.PaymentMeanT.PaymentDueDate)

	crUpdTime := new(commonstruct.CrUpdTime)
	crUpdTime.CreatedAt = common.TimestampToTime(paymentMean.CrUpdTime.CreatedAt)
	crUpdTime.UpdatedAt = common.TimestampToTime(paymentMean.CrUpdTime.UpdatedAt)

	paymentMeanTmp := paymentstruct.PaymentMean{PaymentMeanD: paymentMean.PaymentMeanD, PaymentMeanT: paymentMeanT, CrUpdUser: paymentMean.CrUpdUser, CrUpdTime: crUpdTime}

	return &paymentMeanTmp, nil
}

// CreatePaymentMandateClause - Create PaymentMandateClause
func (ps *PaymentService) CreatePaymentMandateClause(ctx context.Context, in *paymentproto.CreatePaymentMandateClauseRequest) (*paymentproto.CreatePaymentMandateClauseResponse, error) {
	paymentMandateClause := paymentproto.PaymentMandateClause{}
	paymentMandateClause.PmClId = in.PmClId
	paymentMandateClause.PaymentMandateId = in.PaymentMandateId

	err := ps.insertPaymentMandateClause(ctx, insertPaymentMandateClauseSQL, &paymentMandateClause, in.GetUserEmail(), in.GetRequestId())
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	paymentMandateClauseResponse := paymentproto.CreatePaymentMandateClauseResponse{}
	paymentMandateClauseResponse.PaymentMandateClause = &paymentMandateClause
	return &paymentMandateClauseResponse, nil
}

// insertPaymentMandateClause - Insert payment details into database
func (ps *PaymentService) insertPaymentMandateClause(ctx context.Context, insertPaymentMandateClauseSQL string, paymentMandateClause *paymentproto.PaymentMandateClause, userEmail string, requestID string) error {
	err := ps.DBService.InsUpd(ctx, userEmail, requestID, func(tx *sqlx.Tx) error {
		res, err := tx.NamedExecContext(ctx, insertPaymentMandateClauseSQL, paymentMandateClause)
		if err != nil {
			ps.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}

		uID, err := res.LastInsertId()
		if err != nil {
			ps.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}
		paymentMandateClause.Id = uint32(uID)
		return nil
	})
	if err != nil {
		ps.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
		return err
	}
	return nil
}

// CreatePaymentMandateClauseContent - Create PaymentMandateClauseContent
func (ps *PaymentService) CreatePaymentMandateClauseContent(ctx context.Context, in *paymentproto.CreatePaymentMandateClauseContentRequest) (*paymentproto.CreatePaymentMandateClauseContentResponse, error) {
	user, err := partyservice.GetUserWithNewContext(ctx, in.UserId, in.UserEmail, in.RequestId, ps.UserServiceClient)
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	ttime := common.GetTimeDetails()
	tn := common.TimeToTimestamp(ttime)

	paymentMandateClauseContentD := paymentproto.PaymentMandateClauseContentD{}
	paymentMandateClauseContentD.Content = in.Content
	paymentMandateClauseContentD.PaymentMandateClauseId = in.PaymentMandateClauseId

	crUpdUser := commonproto.CrUpdUser{}
	crUpdUser.StatusCode = "active"
	crUpdUser.CreatedByUserId = user.Id
	crUpdUser.UpdatedByUserId = user.Id

	crUpdTime := commonproto.CrUpdTime{}
	crUpdTime.CreatedAt = tn
	crUpdTime.UpdatedAt = tn

	paymentMandateClauseContent := paymentproto.PaymentMandateClauseContent{PaymentMandateClauseContentD: &paymentMandateClauseContentD, CrUpdUser: &crUpdUser, CrUpdTime: &crUpdTime}

	err = ps.insertPaymentMandateClauseContent(ctx, insertPaymentMandateClauseContentSQL, &paymentMandateClauseContent, in.GetUserEmail(), in.GetRequestId())
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	paymentMandateClauseContentResponse := paymentproto.CreatePaymentMandateClauseContentResponse{}
	paymentMandateClauseContentResponse.PaymentMandateClauseContent = &paymentMandateClauseContent
	return &paymentMandateClauseContentResponse, nil
}

// insertPaymentMandateClauseContent - Insert payment details into database
func (ps *PaymentService) insertPaymentMandateClauseContent(ctx context.Context, insertPaymentMandateClauseContentSQL string, paymentMandateClauseContent *paymentproto.PaymentMandateClauseContent, userEmail string, requestID string) error {
	paymentMandateClauseContentTmp, err := ps.crPaymentMandateClauseContentStruct(ctx, paymentMandateClauseContent, userEmail, requestID)
	if err != nil {
		ps.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
		return err
	}

	err = ps.DBService.InsUpd(ctx, userEmail, requestID, func(tx *sqlx.Tx) error {
		res, err := tx.NamedExecContext(ctx, insertPaymentMandateClauseContentSQL, paymentMandateClauseContentTmp)
		if err != nil {
			ps.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}

		uID, err := res.LastInsertId()
		if err != nil {
			ps.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}
		paymentMandateClauseContent.PaymentMandateClauseContentD.Id = uint32(uID)
		return nil
	})

	if err != nil {
		ps.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
		return err
	}
	return nil
}

// crPaymentMandateClauseContentStruct - process PaymentMandateClauseContent details
func (ps *PaymentService) crPaymentMandateClauseContentStruct(ctx context.Context, paymentMandateClauseContent *paymentproto.PaymentMandateClauseContent, userEmail string, requestID string) (*paymentstruct.PaymentMandateClauseContent, error) {
	crUpdTime := new(commonstruct.CrUpdTime)
	crUpdTime.CreatedAt = common.TimestampToTime(paymentMandateClauseContent.CrUpdTime.CreatedAt)
	crUpdTime.UpdatedAt = common.TimestampToTime(paymentMandateClauseContent.CrUpdTime.UpdatedAt)

	paymentMandateClauseContentTmp := paymentstruct.PaymentMandateClauseContent{PaymentMandateClauseContentD: paymentMandateClauseContent.PaymentMandateClauseContentD, CrUpdUser: paymentMandateClauseContent.CrUpdUser, CrUpdTime: crUpdTime}

	return &paymentMandateClauseContentTmp, nil
}
