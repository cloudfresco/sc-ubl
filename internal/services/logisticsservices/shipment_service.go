package logisticsservices

import (
	"context"
	"net"
	"os"
	"time"

	"github.com/cloudfresco/sc-ubl/internal/common"
	"github.com/cloudfresco/sc-ubl/internal/config"
	commonproto "github.com/cloudfresco/sc-ubl/internal/protogen/common/v1"
	logisticsproto "github.com/cloudfresco/sc-ubl/internal/protogen/logistics/v1"
	partyproto "github.com/cloudfresco/sc-ubl/internal/protogen/party/v1"
	partyservice "github.com/cloudfresco/sc-ubl/internal/services/partyservices"
	commonstruct "github.com/cloudfresco/sc-ubl/internal/servicestructs/common/v1"
	logisticsstruct "github.com/cloudfresco/sc-ubl/internal/servicestructs/logistics/v1"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// ShipmentService - For accessing Shipment services
type ShipmentService struct {
	log               *zap.Logger
	DBService         *common.DBService
	RedisService      *common.RedisService
	UserServiceClient partyproto.UserServiceClient
	logisticsproto.UnimplementedShipmentServiceServer
}

// NewShipmentService - Create Shipment service
func NewShipmentService(log *zap.Logger, dbOpt *common.DBService, redisOpt *common.RedisService, userServiceClient partyproto.UserServiceClient) *ShipmentService {
	return &ShipmentService{
		log:               log,
		DBService:         dbOpt,
		RedisService:      redisOpt,
		UserServiceClient: userServiceClient,
	}
}

// StartLogisticsServer - Start Logistics server
func StartLogisticsServer(log *zap.Logger, isTest bool, pwd string, dbOpt *config.DBOptions, redisOpt *config.RedisOptions, mailerOpt *config.MailerOptions, grpcServerOpt *config.GrpcServerOptions, jwtOpt *config.JWTOptions, oauthOpt *config.OauthOptions, userOpt *config.UserOptions, uptraceOpt *config.UptraceOptions, dbService *common.DBService, redisService *common.RedisService, mailerService common.MailerIntf) {
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
	shipmentService := NewShipmentService(log, dbService, redisService, uc)
	despatchService := NewDespatchService(log, dbService, redisService, uc)
	consignmentService := NewConsignmentService(log, dbService, redisService, uc)
	receiptAdviceHeaderService := NewReceiptAdviceHeaderService(log, dbService, redisService, uc)

	lis, err := net.Listen("tcp", grpcServerOpt.GrpcLogisticsServerPort)
	if err != nil {
		log.Error("Error", zap.Error(err))
		os.Exit(1)
	}

	srv := grpc.NewServer(srvOpts...)
	logisticsproto.RegisterShipmentServiceServer(srv, shipmentService)
	logisticsproto.RegisterDespatchServiceServer(srv, despatchService)
	logisticsproto.RegisterConsignmentServiceServer(srv, consignmentService)
	logisticsproto.RegisterReceiptAdviceHeaderServiceServer(srv, receiptAdviceHeaderService)

	if err := srv.Serve(lis); err != nil {
		log.Error("Error", zap.Error(err))
		os.Exit(1)
	}
}

const insertAllowanceChargeSQL = `insert into allowance_charges
	  ( 
    ac_id,
    charge_indicator,
    allowance_charge_reason_code,
    allowance_charge_reason,
    multiplier_factor_numeric,
    prepaid_indicator,
    sequence_numeric,
    amount,
    base_amount,
    per_unit_amount,
    tax_category_id,
    tax_total_id)
  values (:ac_id,
:charge_indicator,
:allowance_charge_reason_code,
:allowance_charge_reason,
:multiplier_factor_numeric,
:prepaid_indicator,
:sequence_numeric,
:amount,
:base_amount,
:per_unit_amount,
:tax_category_id,
:tax_total_id);`

const insertDeliverySQL = `insert into deliveries
	  ( 
    uuid4,
    del_id,
    quantity,
    minimum_quantity,
    maximum_quantity,
    release_id,
    tracking_id,
    minimum_batch_quantity,
    maximum_batch_quantity,
    consumer_unit_quantity,
    hazardous_risk_indicator,
    delivery_address_id,
    delivery_location_id,
    alternative_delivery_location_id,
    carrier_party_id,
    delivery_party_id,
    notify_party_id,
    despatch_id,
    shipment_id,
    actual_delivery_date,
    latest_delivery_date,
    requested_delivery_period_start_date,
    requested_delivery_period_end_date,
    promised_delivery_period_start_date,
    promised_delivery_period_end_date,
    estimated_delivery_period_start_date,
    estimated_delivery_period_end_date,
    status_code,
    created_by_user_id,
    updated_by_user_id,
    created_at,
    updated_at)
  values (:uuid4,
:del_id,
:quantity,
:minimum_quantity,
:maximum_quantity,
:release_id,
:tracking_id,
:minimum_batch_quantity,
:maximum_batch_quantity,
:consumer_unit_quantity,
:hazardous_risk_indicator,
:delivery_address_id,
:delivery_location_id,
:alternative_delivery_location_id,
:carrier_party_id,
:delivery_party_id,
:notify_party_id,
:despatch_id,
:shipment_id,
:actual_delivery_date,
:latest_delivery_date,
:requested_delivery_period_start_date,
:requested_delivery_period_end_date,
:promised_delivery_period_start_date,
:promised_delivery_period_end_date,
:estimated_delivery_period_start_date,
:estimated_delivery_period_end_date,
:status_code,
:created_by_user_id,
:updated_by_user_id,
:created_at,
:updated_at);`

const insertDeliveryTermSQL = `insert into delivery_terms
	  ( 
    del_term_id,
    special_terms,
    loss_risk_responsibility_code,
    loss_risk,
    amount,
    delivery_location_id,
    del_term_allowance_charge_id,
    delivery_id,
    status_code,
    created_by_user_id,
    updated_by_user_id,
    created_at,
    updated_at)
  values (:del_term_id,
:special_terms,
:loss_risk_responsibility_code,
:loss_risk,
:amount,
:delivery_location_id,
:del_term_allowance_charge_id,
:delivery_id,
:status_code,
:created_by_user_id,
:updated_by_user_id,
:created_at,
:updated_at);`

const insertDespatchSQL = `insert into despatches
	  ( 
    desp_id,
    release_id,
    instructions,
    despatch_address_id,
    despatch_location_id,
    despatch_party_contact,
    despatch_party_id,
    carrier_party_id,
    notify_party_id,
    requested_despatch_date,
    estimated_despatch_date,
    actual_despatch_date,
    guaranteed_despatch_date,
    estimated_despatch_period_start_date,
    estimated_despatch_period_end_date,
    requested_despatch_period_start_date,
    requested_despatch_period_end_date)
  values (:desp_id,
:release_id,
:instructions,
:despatch_address_id,
:despatch_location_id,
:despatch_party_contact,
:despatch_party_id,
:carrier_party_id,
:notify_party_id,
:requested_despatch_date,
:estimated_despatch_date,
:actual_despatch_date,
:guaranteed_despatch_date,
:estimated_despatch_period_start_date,
:estimated_despatch_period_end_date,
:requested_despatch_period_start_date,
:requested_despatch_period_end_date);`

const insertShipmentSQL = `insert into shipments
	  (
    uuid4,
    sh_id,
    shipping_priority_level_code,
    handling_code,
    handling_instructions,
    information,
    gross_weight_measure,
    net_weight_measure,
    net_net_weight_measure,
    gross_volume_measure,
    net_volume_measure,
    total_goods_item_quantity,
    total_transport_handling_unit_quantity,
    insurance_value_amount,
    declared_customs_value_amount,
    declared_for_carriage_value_amount,
    declared_statistics_value_amount,
    free_on_board_value_amount,
    special_instructions,
    delivery_instructions,
    split_consignment_indicator,
    consignment_quantity,
    return_address_id,
    origin_address_id,
    first_arrival_port_location_id,
    last_exit_port_location_id,
    export_country_id_code,
    export_country_name,
    status_code,
    created_by_user_id,
    updated_by_user_id,
    created_at,
    updated_at)
  values (:uuid4,
:sh_id,
:shipping_priority_level_code,
:handling_code,
:handling_instructions,
:information,
:gross_weight_measure,
:net_weight_measure,
:net_net_weight_measure,
:gross_volume_measure,
:net_volume_measure,
:total_goods_item_quantity,
:total_transport_handling_unit_quantity,
:insurance_value_amount,
:declared_customs_value_amount,
:declared_for_carriage_value_amount,
:declared_statistics_value_amount,
:free_on_board_value_amount,
:special_instructions,
:delivery_instructions,
:split_consignment_indicator,
:consignment_quantity,
:return_address_id,
:origin_address_id,
:first_arrival_port_location_id,
:last_exit_port_location_id,
:export_country_id_code,
:export_country_name,
:status_code,
:created_by_user_id,
:updated_by_user_id,
:created_at,
:updated_at);`

const insertBillOfLadingSQL = `insert into bill_of_ladings
	  (
    uuid4,
    bill_of_lading_id,
    carrier_assigned_id,
    name1,
    description,
    note,
    document_status_code,
    shipping_order_id,
    to_order_indicator,
    ad_valorem_indicator,
    declared_carriage_value_amount,
    declared_carriage_value_amount_currency_code,
    other_instruction,
    consignor_party_id,
    carrier_party_id,
    freight_forwarder_party_id,
    shipment_id,
    issue_date,
    status_code,
    created_by_user_id,
    updated_by_user_id,
    created_at,
    updated_at)
  values (:uuid4,
:bill_of_lading_id,
:carrier_assigned_id,
:name1,
:description,
:note,
:document_status_code,
:shipping_order_id,
:to_order_indicator,
:ad_valorem_indicator,
:declared_carriage_value_amount,
:declared_carriage_value_amount_currency_code,
:other_instruction,
:consignor_party_id,
:carrier_party_id,
:freight_forwarder_party_id,
:shipment_id,
:issue_date,
:status_code,
:created_by_user_id,
:updated_by_user_id,
:created_at,
:updated_at);`

const insertWaybillSQL = `insert into waybills
	  (
    uuid4,
    waybill_id,
    carrier_assigned_id,
    name1,
    description,
    note,
    shipping_order_id,
    ad_valorem_indicator,
    declared_carriage_value_amount,
    declared_carriage_value_amount_currency_code,
    other_instruction,
    consignor_party_id,
    carrier_party_id,
    freight_forwarder_party_id,
    shipment_id,
    issue_date,
    status_code,
    created_by_user_id,
    updated_by_user_id,
    created_at,
    updated_at)
  values (:uuid4,
:waybill_id,
:carrier_assigned_id,
:name1,
:description,
:note,
:shipping_order_id,
:ad_valorem_indicator,
:declared_carriage_value_amount,
:declared_carriage_value_amount_currency_code,
:other_instruction,
:consignor_party_id,
:carrier_party_id,
:freight_forwarder_party_id,
:shipment_id,
:issue_date,
:status_code,
:created_by_user_id,
:updated_by_user_id,
:created_at,
:updated_at);`

// CreateShipment - Create Shipment
func (ss *ShipmentService) CreateShipment(ctx context.Context, in *logisticsproto.CreateShipmentRequest) (*logisticsproto.CreateShipmentResponse, error) {
	user, err := partyservice.GetUserWithNewContext(ctx, in.UserId, in.UserEmail, in.RequestId, ss.UserServiceClient)
	if err != nil {
		ss.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	ttime := common.GetTimeDetails()
	tn := common.TimeToTimestamp(ttime)

	shipmentD := logisticsproto.ShipmentD{}
	shipmentD.Uuid4, err = common.GetUUIDBytes()
	if err != nil {
		ss.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	shipmentD.ShId = in.ShId
	shipmentD.ShippingPriorityLevelCode = in.ShippingPriorityLevelCode
	shipmentD.HandlingCode = in.HandlingCode
	shipmentD.HandlingInstructions = in.HandlingInstructions
	shipmentD.Information = in.Information
	shipmentD.GrossWeightMeasure = in.GrossWeightMeasure
	shipmentD.NetWeightMeasure = in.NetWeightMeasure
	shipmentD.NetNetWeightMeasure = in.NetNetWeightMeasure
	shipmentD.GrossVolumeMeasure = in.GrossVolumeMeasure
	shipmentD.NetVolumeMeasure = in.NetVolumeMeasure
	shipmentD.TotalGoodsItemQuantity = in.TotalGoodsItemQuantity
	shipmentD.TotalTransportHandlingUnitQuantity = in.TotalTransportHandlingUnitQuantity
	shipmentD.InsuranceValueAmount = in.InsuranceValueAmount
	shipmentD.DeclaredCustomsValueAmount = in.DeclaredCustomsValueAmount
	shipmentD.DeclaredForCarriageValueAmount = in.DeclaredForCarriageValueAmount
	shipmentD.DeclaredStatisticsValueAmount = in.DeclaredStatisticsValueAmount
	shipmentD.FreeOnBoardValueAmount = in.FreeOnBoardValueAmount
	shipmentD.SpecialInstructions = in.SpecialInstructions
	shipmentD.DeliveryInstructions = in.DeliveryInstructions
	shipmentD.SplitConsignmentIndicator = in.SplitConsignmentIndicator
	shipmentD.ConsignmentQuantity = in.ConsignmentQuantity
	shipmentD.ExportCountryIdCode = in.ExportCountryIdCode
	shipmentD.ExportCountryName = in.ExportCountryName

	crUpdUser := commonproto.CrUpdUser{}
	crUpdUser.StatusCode = "active"
	crUpdUser.CreatedByUserId = user.Id
	crUpdUser.UpdatedByUserId = user.Id

	crUpdTime := commonproto.CrUpdTime{}
	crUpdTime.CreatedAt = tn
	crUpdTime.UpdatedAt = tn

	returnAddr, err := common.CreateAddress(ctx, in.ReturnAddress, in.GetUserEmail(), in.GetRequestId())
	if err != nil {
		ss.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	originAddr, err := common.CreateAddress(ctx, in.OriginalAddress, in.GetUserEmail(), in.GetRequestId())
	if err != nil {
		ss.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	arrivalLocation, err := common.CreateLocation(ctx, in.FirstArrivalPortLocation, in.GetUserEmail(), in.GetRequestId())
	if err != nil {
		ss.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	exitLocation, err := common.CreateLocation(ctx, in.LastExitPortLocation, in.GetUserEmail(), in.GetRequestId())
	if err != nil {
		ss.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	shipment := logisticsproto.Shipment{ShipmentD: &shipmentD, CrUpdUser: &crUpdUser, CrUpdTime: &crUpdTime}

	err = ss.insertShipment(ctx, insertShipmentSQL, &shipment, returnAddr, originAddr, arrivalLocation, exitLocation, in.GetUserEmail(), in.GetRequestId())
	if err != nil {
		ss.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	shipmentResponse := logisticsproto.CreateShipmentResponse{}
	shipmentResponse.Shipment = &shipment
	return &shipmentResponse, nil
}

// insertShipment - Insert shipment details into database
func (ss *ShipmentService) insertShipment(ctx context.Context, insertShipmentSQL string, shipment *logisticsproto.Shipment, raddr *commonproto.Address, oaddr *commonproto.Address, al *commonproto.Location, el *commonproto.Location, userEmail string, requestID string) error {
	shipmentTmp, err := ss.crShipmentStruct(ctx, shipment, userEmail, requestID)
	if err != nil {
		ss.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
		return err
	}

	err = ss.DBService.InsUpd(ctx, userEmail, requestID, func(tx *sqlx.Tx) error {
		returnAddr, err := common.InsertAddress(ctx, tx, raddr, userEmail, requestID)
		if err != nil {
			ss.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}
		originAddr, err := common.InsertAddress(ctx, tx, oaddr, userEmail, requestID)
		if err != nil {
			ss.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}
		arrivalLocation, err := common.InsertLocation(ctx, tx, al, userEmail, requestID)
		if err != nil {
			ss.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}
		exitLocation, err := common.InsertLocation(ctx, tx, el, userEmail, requestID)
		if err != nil {
			ss.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}

		shipment.ShipmentD.ReturnAddressId = returnAddr.Id
		shipment.ShipmentD.OriginAddressId = originAddr.Id
		shipment.ShipmentD.FirstArrivalPortLocationId = arrivalLocation.LocationD.Id
		shipment.ShipmentD.LastExitPortLocationId = exitLocation.LocationD.Id
		shipmentTmp.ShipmentD.ReturnAddressId = returnAddr.Id
		shipmentTmp.ShipmentD.OriginAddressId = originAddr.Id
		shipmentTmp.ShipmentD.FirstArrivalPortLocationId = arrivalLocation.LocationD.Id
		shipmentTmp.ShipmentD.LastExitPortLocationId = exitLocation.LocationD.Id

		res, err := tx.NamedExecContext(ctx, insertShipmentSQL, shipmentTmp)
		if err != nil {
			ss.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}

		uID, err := res.LastInsertId()
		if err != nil {
			ss.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}
		shipment.ShipmentD.Id = uint32(uID)
		uuid4Str, err := common.UUIDBytesToStr(shipment.ShipmentD.Uuid4)
		if err != nil {
			ss.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}
		shipment.ShipmentD.IdS = uuid4Str
		return nil
	})

	if err != nil {
		ss.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
		return err
	}
	return nil
}

// crShipmentStruct - process Shipment details
func (ss *ShipmentService) crShipmentStruct(ctx context.Context, shipment *logisticsproto.Shipment, userEmail string, requestID string) (*logisticsstruct.Shipment, error) {
	crUpdTime := new(commonstruct.CrUpdTime)
	crUpdTime.CreatedAt = common.TimestampToTime(shipment.CrUpdTime.CreatedAt)
	crUpdTime.UpdatedAt = common.TimestampToTime(shipment.CrUpdTime.UpdatedAt)

	shipmentTmp := logisticsstruct.Shipment{ShipmentD: shipment.ShipmentD, CrUpdUser: shipment.CrUpdUser, CrUpdTime: crUpdTime}

	return &shipmentTmp, nil
}

// CreateAllowanceCharge - CreateAllowanceCharge
func (ss *ShipmentService) CreateAllowanceCharge(ctx context.Context, in *logisticsproto.CreateAllowanceChargeRequest) (*logisticsproto.CreateAllowanceChargeResponse, error) {
	allowanceCharge := logisticsproto.AllowanceCharge{}
	allowanceCharge.AcId = in.AcId
	allowanceCharge.ChargeIndicator = in.ChargeIndicator
	allowanceCharge.AllowanceChargeReasonCode = in.AllowanceChargeReasonCode
	allowanceCharge.AllowanceChargeReason = in.AllowanceChargeReason
	allowanceCharge.MultiplierFactorNumeric = in.MultiplierFactorNumeric
	allowanceCharge.PrepaidIndicator = in.PrepaidIndicator
	allowanceCharge.SequenceNumeric = in.SequenceNumeric
	allowanceCharge.Amount = in.Amount
	allowanceCharge.BaseAmount = in.BaseAmount
	allowanceCharge.PerUnitAmount = in.PerUnitAmount
	allowanceCharge.TaxCategoryId = in.TaxCategoryId
	allowanceCharge.TaxTotalId = in.TaxTotalId

	err := ss.insertAllowanceCharge(ctx, insertAllowanceChargeSQL, &allowanceCharge, in.GetUserEmail(), in.GetRequestId())
	if err != nil {
		ss.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	allowanceChargeResponse := logisticsproto.CreateAllowanceChargeResponse{}
	allowanceChargeResponse.AllowanceCharge = &allowanceCharge
	return &allowanceChargeResponse, nil
}

// insertAllowanceCharge - InsertAllowanceCharge details into database
func (ss *ShipmentService) insertAllowanceCharge(ctx context.Context, insertAllowanceChargeSQL string, allowanceCharge *logisticsproto.AllowanceCharge, userEmail string, requestID string) error {
	err := ss.DBService.InsUpd(ctx, userEmail, requestID, func(tx *sqlx.Tx) error {
		res, err := tx.NamedExecContext(ctx, insertAllowanceChargeSQL, allowanceCharge)
		if err != nil {
			ss.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}

		uID, err := res.LastInsertId()
		if err != nil {
			ss.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}
		allowanceCharge.Id = uint32(uID)
		return nil
	})
	if err != nil {
		ss.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
		return err
	}
	return nil
}

// CreateDespatch - CreateDespatch
func (ss *ShipmentService) CreateDespatch(ctx context.Context, in *logisticsproto.CreateDespatchRequest) (*logisticsproto.CreateDespatchResponse, error) {
	requestedDespatchDate, err := time.Parse(common.Layout, in.RequestedDespatchDate)
	if err != nil {
		ss.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	estimatedDespatchDate, err := time.Parse(common.Layout, in.EstimatedDespatchDate)
	if err != nil {
		ss.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	actualDespatchDate, err := time.Parse(common.Layout, in.ActualDespatchDate)
	if err != nil {
		ss.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	guaranteedDespatchDate, err := time.Parse(common.Layout, in.GuaranteedDespatchDate)
	if err != nil {
		ss.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	estimatedDespatchPeriodStartDate, err := time.Parse(common.Layout, in.EstimatedDespatchPeriodStartDate)
	if err != nil {
		ss.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	estimatedDespatchPeriodEndDate, err := time.Parse(common.Layout, in.EstimatedDespatchPeriodEndDate)
	if err != nil {
		ss.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	requestedDespatchPeriodStartDate, err := time.Parse(common.Layout, in.RequestedDespatchPeriodStartDate)
	if err != nil {
		ss.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	requestedDespatchPeriodEndDate, err := time.Parse(common.Layout, in.RequestedDespatchPeriodEndDate)
	if err != nil {
		ss.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	despatchAddr, err := common.CreateAddress(ctx, in.DespatchAddress, in.GetUserEmail(), in.GetRequestId())
	if err != nil {
		ss.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	despatchLocation, err := common.CreateLocation(ctx, in.DespatchLocation, in.GetUserEmail(), in.GetRequestId())
	if err != nil {
		ss.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	despatchD := logisticsproto.DespatchD{}
	despatchD.DespId = in.DespId
	despatchD.ReleaseId = in.ReleaseId
	despatchD.Instructions = in.Instructions
	despatchD.DespatchPartyContact = in.DespatchPartyContact
	despatchD.DespatchPartyId = in.DespatchPartyId
	despatchD.CarrierPartyId = in.CarrierPartyId
	despatchD.NotifyPartyId = in.NotifyPartyId

	despatchT := logisticsproto.DespatchT{}
	despatchT.RequestedDespatchDate = common.TimeToTimestamp(requestedDespatchDate.UTC().Truncate(time.Second))
	despatchT.EstimatedDespatchDate = common.TimeToTimestamp(estimatedDespatchDate.UTC().Truncate(time.Second))
	despatchT.ActualDespatchDate = common.TimeToTimestamp(actualDespatchDate.UTC().Truncate(time.Second))
	despatchT.GuaranteedDespatchDate = common.TimeToTimestamp(guaranteedDespatchDate.UTC().Truncate(time.Second))
	despatchT.EstimatedDespatchPeriodStartDate = common.TimeToTimestamp(estimatedDespatchPeriodStartDate.UTC().Truncate(time.Second))
	despatchT.EstimatedDespatchPeriodEndDate = common.TimeToTimestamp(estimatedDespatchPeriodEndDate.UTC().Truncate(time.Second))
	despatchT.RequestedDespatchPeriodStartDate = common.TimeToTimestamp(requestedDespatchPeriodStartDate.UTC().Truncate(time.Second))
	despatchT.RequestedDespatchPeriodEndDate = common.TimeToTimestamp(requestedDespatchPeriodEndDate.UTC().Truncate(time.Second))

	despatch := logisticsproto.Despatch{DespatchD: &despatchD, DespatchT: &despatchT}

	err = ss.insertDespatch(ctx, insertDespatchSQL, &despatch, despatchAddr, despatchLocation, in.GetUserEmail(), in.GetRequestId())
	if err != nil {
		ss.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))

		return nil, err
	}
	despatchResponse := logisticsproto.CreateDespatchResponse{}
	despatchResponse.Despatch = &despatch
	return &despatchResponse, nil
}

// insertDespatch - InsertDespatch details into database
func (ss *ShipmentService) insertDespatch(ctx context.Context, insertDespatchSQL string, despatch *logisticsproto.Despatch, daddr *commonproto.Address, dl *commonproto.Location, userEmail string, requestID string) error {
	despatchTmp, err := ss.crDespatchStruct(ctx, despatch, userEmail, requestID)
	if err != nil {
		ss.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
		return err
	}
	err = ss.DBService.InsUpd(ctx, userEmail, requestID, func(tx *sqlx.Tx) error {
		despatchAddress, err := common.InsertAddress(ctx, tx, daddr, userEmail, requestID)
		if err != nil {
			ss.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}
		despatchLocation, err := common.InsertLocation(ctx, tx, dl, userEmail, requestID)
		if err != nil {
			ss.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}

		despatch.DespatchD.DespatchAddressId = despatchAddress.Id
		despatch.DespatchD.DespatchLocationId = despatchLocation.LocationD.Id
		despatchTmp.DespatchD.DespatchAddressId = despatchAddress.Id
		despatchTmp.DespatchD.DespatchLocationId = despatchLocation.LocationD.Id

		res, err := tx.NamedExecContext(ctx, insertDespatchSQL, despatchTmp)
		if err != nil {
			ss.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}

		uID, err := res.LastInsertId()
		if err != nil {
			ss.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}
		despatch.DespatchD.Id = uint32(uID)
		return nil
	})

	if err != nil {
		ss.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
		return err
	}
	return nil
}

// crDespatchStruct - process Despatch details
func (ss *ShipmentService) crDespatchStruct(ctx context.Context, despatch *logisticsproto.Despatch, userEmail string, requestID string) (*logisticsstruct.Despatch, error) {
	despatchT := new(logisticsstruct.DespatchT)
	despatchT.RequestedDespatchDate = common.TimestampToTime(despatch.DespatchT.RequestedDespatchDate)
	despatchT.EstimatedDespatchDate = common.TimestampToTime(despatch.DespatchT.EstimatedDespatchDate)
	despatchT.ActualDespatchDate = common.TimestampToTime(despatch.DespatchT.ActualDespatchDate)
	despatchT.GuaranteedDespatchDate = common.TimestampToTime(despatch.DespatchT.GuaranteedDespatchDate)
	despatchT.EstimatedDespatchPeriodStartDate = common.TimestampToTime(despatch.DespatchT.EstimatedDespatchPeriodStartDate)
	despatchT.EstimatedDespatchPeriodEndDate = common.TimestampToTime(despatch.DespatchT.EstimatedDespatchPeriodEndDate)
	despatchT.RequestedDespatchPeriodStartDate = common.TimestampToTime(despatch.DespatchT.RequestedDespatchPeriodStartDate)
	despatchT.RequestedDespatchPeriodEndDate = common.TimestampToTime(despatch.DespatchT.RequestedDespatchPeriodEndDate)

	despatchTmp := logisticsstruct.Despatch{DespatchD: despatch.DespatchD, DespatchT: despatchT}

	return &despatchTmp, nil
}

// CreateDelivery - Create Delivery
func (ss *ShipmentService) CreateDelivery(ctx context.Context, in *logisticsproto.CreateDeliveryRequest) (*logisticsproto.CreateDeliveryResponse, error) {
	user, err := partyservice.GetUserWithNewContext(ctx, in.UserId, in.UserEmail, in.RequestId, ss.UserServiceClient)
	if err != nil {
		ss.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	ttime := common.GetTimeDetails()
	tn := common.TimeToTimestamp(ttime)

	actualDeliveryDate, err := time.Parse(common.Layout, in.ActualDeliveryDate)
	if err != nil {
		ss.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	latestDeliveryDate, err := time.Parse(common.Layout, in.LatestDeliveryDate)
	if err != nil {
		ss.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	requestedDeliveryPeriodStartDate, err := time.Parse(common.Layout, in.RequestedDeliveryPeriodStartDate)
	if err != nil {
		ss.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	requestedDeliveryPeriodEndDate, err := time.Parse(common.Layout, in.RequestedDeliveryPeriodEndDate)
	if err != nil {
		ss.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	promisedDeliveryPeriodStartDate, err := time.Parse(common.Layout, in.PromisedDeliveryPeriodStartDate)
	if err != nil {
		ss.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	promisedDeliveryPeriodEndDate, err := time.Parse(common.Layout, in.PromisedDeliveryPeriodEndDate)
	if err != nil {
		ss.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	estimatedDeliveryPeriodStartDate, err := time.Parse(common.Layout, in.EstimatedDeliveryPeriodStartDate)
	if err != nil {
		ss.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	estimatedDeliveryPeriodEndDate, err := time.Parse(common.Layout, in.EstimatedDeliveryPeriodEndDate)
	if err != nil {
		ss.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	deliveryD := logisticsproto.DeliveryD{}
	deliveryD.Uuid4, err = common.GetUUIDBytes()
	if err != nil {
		ss.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	deliveryD.DelId = in.DelId
	deliveryD.Quantity = in.Quantity
	deliveryD.MinimumQuantity = in.MinimumQuantity
	deliveryD.MaximumQuantity = in.MaximumQuantity
	deliveryD.ReleaseId = in.ReleaseId
	deliveryD.TrackingId = in.TrackingId
	deliveryD.MinimumBatchQuantity = in.MinimumBatchQuantity
	deliveryD.MaximumBatchQuantity = in.MaximumBatchQuantity
	deliveryD.ConsumerUnitQuantity = in.ConsumerUnitQuantity
	deliveryD.HazardousRiskIndicator = in.HazardousRiskIndicator
	deliveryD.CarrierPartyId = in.CarrierPartyId
	deliveryD.DeliveryPartyId = in.DeliveryPartyId
	deliveryD.NotifyPartyId = in.NotifyPartyId
	deliveryD.DespatchId = in.DespatchId
	deliveryD.ShipmentId = in.DespatchId

	deliveryT := logisticsproto.DeliveryT{}
	deliveryT.ActualDeliveryDate = common.TimeToTimestamp(actualDeliveryDate.UTC().Truncate(time.Second))
	deliveryT.LatestDeliveryDate = common.TimeToTimestamp(latestDeliveryDate.UTC().Truncate(time.Second))
	deliveryT.RequestedDeliveryPeriodStartDate = common.TimeToTimestamp(requestedDeliveryPeriodStartDate.UTC().Truncate(time.Second))
	deliveryT.RequestedDeliveryPeriodEndDate = common.TimeToTimestamp(requestedDeliveryPeriodEndDate.UTC().Truncate(time.Second))
	deliveryT.PromisedDeliveryPeriodStartDate = common.TimeToTimestamp(promisedDeliveryPeriodStartDate.UTC().Truncate(time.Second))
	deliveryT.PromisedDeliveryPeriodEndDate = common.TimeToTimestamp(promisedDeliveryPeriodEndDate.UTC().Truncate(time.Second))
	deliveryT.EstimatedDeliveryPeriodStartDate = common.TimeToTimestamp(estimatedDeliveryPeriodStartDate.UTC().Truncate(time.Second))
	deliveryT.EstimatedDeliveryPeriodEndDate = common.TimeToTimestamp(estimatedDeliveryPeriodEndDate.UTC().Truncate(time.Second))

	crUpdUser := commonproto.CrUpdUser{}
	crUpdUser.StatusCode = "active"
	crUpdUser.CreatedByUserId = user.Id
	crUpdUser.UpdatedByUserId = user.Id

	crUpdTime := commonproto.CrUpdTime{}
	crUpdTime.CreatedAt = tn
	crUpdTime.UpdatedAt = tn

	deliveryAddress, err := common.CreateAddress(ctx, in.DeliveryAddress, in.GetUserEmail(), in.GetRequestId())
	if err != nil {
		ss.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	deliveryLocation, err := common.CreateLocation(ctx, in.DeliveryLocation, in.GetUserEmail(), in.GetRequestId())
	if err != nil {
		ss.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	alternativeDeliveryLocation, err := common.CreateLocation(ctx, in.AlternativeDeliveryLocation, in.GetUserEmail(), in.GetRequestId())
	if err != nil {
		ss.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	delivery := logisticsproto.Delivery{DeliveryD: &deliveryD, DeliveryT: &deliveryT, CrUpdUser: &crUpdUser, CrUpdTime: &crUpdTime}

	err = ss.insertDelivery(ctx, insertDeliverySQL, &delivery, deliveryAddress, deliveryLocation, alternativeDeliveryLocation, in.GetUserEmail(), in.GetRequestId())
	if err != nil {
		ss.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	deliveryResponse := logisticsproto.CreateDeliveryResponse{}
	deliveryResponse.Delivery = &delivery
	return &deliveryResponse, nil
}

// insertDelivery - Insert Delivery details into database
func (ss *ShipmentService) insertDelivery(ctx context.Context, insertDeliverySQL string, delivery *logisticsproto.Delivery, daddr *commonproto.Address, dl *commonproto.Location, adl *commonproto.Location, userEmail string, requestID string) error {
	deliveryTmp, err := ss.crDeliveryStruct(ctx, delivery, userEmail, requestID)
	if err != nil {
		ss.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
		return err
	}

	err = ss.DBService.InsUpd(ctx, userEmail, requestID, func(tx *sqlx.Tx) error {
		deliveryAddress, err := common.InsertAddress(ctx, tx, daddr, userEmail, requestID)
		if err != nil {
			ss.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}
		deliveryLocation, err := common.InsertLocation(ctx, tx, dl, userEmail, requestID)
		if err != nil {
			ss.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}
		alternativeDeliveryLocation, err := common.InsertLocation(ctx, tx, adl, userEmail, requestID)
		if err != nil {
			ss.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}

		delivery.DeliveryD.DeliveryAddressId = deliveryAddress.Id
		delivery.DeliveryD.DeliveryLocationId = deliveryLocation.LocationD.Id
		delivery.DeliveryD.AlternativeDeliveryLocationId = alternativeDeliveryLocation.LocationD.Id
		deliveryTmp.DeliveryD.DeliveryAddressId = deliveryAddress.Id
		deliveryTmp.DeliveryD.DeliveryLocationId = deliveryLocation.LocationD.Id
		deliveryTmp.DeliveryD.AlternativeDeliveryLocationId = alternativeDeliveryLocation.LocationD.Id

		res, err := tx.NamedExecContext(ctx, insertDeliverySQL, deliveryTmp)
		if err != nil {
			ss.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}

		uID, err := res.LastInsertId()
		if err != nil {
			ss.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}
		delivery.DeliveryD.Id = uint32(uID)
		uuid4Str, err := common.UUIDBytesToStr(delivery.DeliveryD.Uuid4)
		if err != nil {
			ss.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}
		delivery.DeliveryD.IdS = uuid4Str
		return nil
	})

	if err != nil {
		ss.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
		return err
	}
	return nil
}

// crDeliveryStruct - process Delivery details
func (ss *ShipmentService) crDeliveryStruct(ctx context.Context, delivery *logisticsproto.Delivery, userEmail string, requestID string) (*logisticsstruct.Delivery, error) {
	deliveryT := new(logisticsstruct.DeliveryT)
	deliveryT.ActualDeliveryDate = common.TimestampToTime(delivery.DeliveryT.ActualDeliveryDate)
	deliveryT.LatestDeliveryDate = common.TimestampToTime(delivery.DeliveryT.LatestDeliveryDate)
	deliveryT.RequestedDeliveryPeriodStartDate = common.TimestampToTime(delivery.DeliveryT.RequestedDeliveryPeriodStartDate)
	deliveryT.RequestedDeliveryPeriodEndDate = common.TimestampToTime(delivery.DeliveryT.RequestedDeliveryPeriodEndDate)
	deliveryT.PromisedDeliveryPeriodStartDate = common.TimestampToTime(delivery.DeliveryT.PromisedDeliveryPeriodStartDate)
	deliveryT.PromisedDeliveryPeriodEndDate = common.TimestampToTime(delivery.DeliveryT.PromisedDeliveryPeriodEndDate)
	deliveryT.EstimatedDeliveryPeriodStartDate = common.TimestampToTime(delivery.DeliveryT.EstimatedDeliveryPeriodStartDate)
	deliveryT.EstimatedDeliveryPeriodEndDate = common.TimestampToTime(delivery.DeliveryT.EstimatedDeliveryPeriodEndDate)

	crUpdTime := new(commonstruct.CrUpdTime)
	crUpdTime.CreatedAt = common.TimestampToTime(delivery.CrUpdTime.CreatedAt)
	crUpdTime.UpdatedAt = common.TimestampToTime(delivery.CrUpdTime.UpdatedAt)

	deliveryTmp := logisticsstruct.Delivery{DeliveryD: delivery.DeliveryD, DeliveryT: deliveryT, CrUpdUser: delivery.CrUpdUser, CrUpdTime: crUpdTime}

	return &deliveryTmp, nil
}

// CreateDeliveryTerm - CreateDeliveryTerm
func (ss *ShipmentService) CreateDeliveryTerm(ctx context.Context, in *logisticsproto.CreateDeliveryTermRequest) (*logisticsproto.CreateDeliveryTermResponse, error) {
	user, err := partyservice.GetUserWithNewContext(ctx, in.UserId, in.UserEmail, in.RequestId, ss.UserServiceClient)
	if err != nil {
		ss.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	ttime := common.GetTimeDetails()
	tn := common.TimeToTimestamp(ttime)

	deliveryLocation, err := common.CreateLocation(ctx, in.DeliveryLocation, in.GetUserEmail(), in.GetRequestId())
	if err != nil {
		ss.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	deliveryTermD := logisticsproto.DeliveryTermD{}
	deliveryTermD.Uuid4, err = common.GetUUIDBytes()
	if err != nil {
		ss.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	deliveryTermD.DelTermId = in.DelTermId
	deliveryTermD.SpecialTerms = in.SpecialTerms
	deliveryTermD.LossRiskResponsibilityCode = in.LossRiskResponsibilityCode
	deliveryTermD.LossRisk = in.LossRisk
	deliveryTermD.Amount = in.Amount
	deliveryTermD.DelTermAllowanceChargeId = in.DelTermAllowanceChargeId
	deliveryTermD.DeliveryId = in.DeliveryId

	crUpdUser := commonproto.CrUpdUser{}
	crUpdUser.StatusCode = "active"
	crUpdUser.CreatedByUserId = user.Id
	crUpdUser.UpdatedByUserId = user.Id

	crUpdTime := commonproto.CrUpdTime{}
	crUpdTime.CreatedAt = tn
	crUpdTime.UpdatedAt = tn

	deliveryTerm := logisticsproto.DeliveryTerm{DeliveryTermD: &deliveryTermD, CrUpdUser: &crUpdUser, CrUpdTime: &crUpdTime}

	err = ss.insertDeliveryTerm(ctx, insertDeliveryTermSQL, &deliveryTerm, deliveryLocation, in.GetUserEmail(), in.GetRequestId())
	if err != nil {
		ss.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))

		return nil, err
	}

	deliveryTermResponse := logisticsproto.CreateDeliveryTermResponse{}
	deliveryTermResponse.DeliveryTerm = &deliveryTerm
	return &deliveryTermResponse, nil
}

// insertDeliveryTerm - InsertDeliveryTerm details into database
func (ss *ShipmentService) insertDeliveryTerm(ctx context.Context, insertDeliveryTermSQL string, deliveryTerm *logisticsproto.DeliveryTerm, dl *commonproto.Location, userEmail string, requestID string) error {
	deliveryTermTmp, err := ss.crDeliveryTermStruct(ctx, deliveryTerm, userEmail, requestID)
	if err != nil {
		ss.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
		return err
	}
	err = ss.DBService.InsUpd(ctx, userEmail, requestID, func(tx *sqlx.Tx) error {
		deliveryLocation, err := common.InsertLocation(ctx, tx, dl, userEmail, requestID)
		if err != nil {
			ss.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}
		deliveryTerm.DeliveryTermD.DeliveryLocationId = deliveryLocation.LocationD.Id
		deliveryTermTmp.DeliveryTermD.DeliveryLocationId = deliveryLocation.LocationD.Id
		res, err := tx.NamedExecContext(ctx, insertDeliveryTermSQL, deliveryTermTmp)
		if err != nil {
			ss.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}

		uID, err := res.LastInsertId()
		if err != nil {
			ss.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}
		deliveryTerm.DeliveryTermD.Id = uint32(uID)
		uuid4Str, err := common.UUIDBytesToStr(deliveryTerm.DeliveryTermD.Uuid4)
		if err != nil {
			ss.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}
		deliveryTerm.DeliveryTermD.IdS = uuid4Str

		return nil
	})
	if err != nil {
		ss.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
		return err
	}
	return nil
}

// crDeliveryTermStruct - process DeliveryTerm details
func (ss *ShipmentService) crDeliveryTermStruct(ctx context.Context, deliveryTerm *logisticsproto.DeliveryTerm, userEmail string, requestID string) (*logisticsstruct.DeliveryTerm, error) {
	crUpdTime := new(commonstruct.CrUpdTime)
	crUpdTime.CreatedAt = common.TimestampToTime(deliveryTerm.CrUpdTime.CreatedAt)
	crUpdTime.UpdatedAt = common.TimestampToTime(deliveryTerm.CrUpdTime.UpdatedAt)

	deliveryTermTmp := logisticsstruct.DeliveryTerm{DeliveryTermD: deliveryTerm.DeliveryTermD, CrUpdUser: deliveryTerm.CrUpdUser, CrUpdTime: crUpdTime}

	return &deliveryTermTmp, nil
}

// CreateBillOfLading - Create BillOfLading
func (ss *ShipmentService) CreateBillOfLading(ctx context.Context, in *logisticsproto.CreateBillOfLadingRequest) (*logisticsproto.CreateBillOfLadingResponse, error) {
	getRequest := commonproto.GetRequest{}
	getRequest.Id = in.UserId
	getRequest.UserEmail = in.UserEmail
	getRequest.RequestId = in.RequestId
	ctxNew, err := common.CreateCtxJWT(ctx)
	if err != nil {
		ss.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	form := partyproto.GetUserRequest{}
	form.GetRequest = &getRequest
	userResponse, err := ss.UserServiceClient.GetUser(ctxNew, &form)
	if err != nil {
		ss.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	user := userResponse.User
	ttime := common.GetTimeDetails()
	tn := common.TimeToTimestamp(ttime)

	issueDate, err := time.Parse(common.Layout, in.IssueDate)
	if err != nil {
		ss.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	billOfLadingD := logisticsproto.BillOfLadingD{}
	billOfLadingD.Uuid4, err = common.GetUUIDBytes()
	if err != nil {
		ss.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	billOfLadingD.BillOfLadingId = in.BillOfLadingId
	billOfLadingD.CarrierAssignedId = in.CarrierAssignedId
	billOfLadingD.Name1 = in.Name1
	billOfLadingD.Description = in.Description
	billOfLadingD.Note = in.Note
	billOfLadingD.DocumentStatusCode = in.DocumentStatusCode
	billOfLadingD.ShippingOrderId = in.ShippingOrderId
	billOfLadingD.ToOrderIndicator = in.ToOrderIndicator
	billOfLadingD.AdValoremIndicator = in.AdValoremIndicator
	billOfLadingD.DeclaredCarriageValueAmount = in.DeclaredCarriageValueAmount
	billOfLadingD.DeclaredCarriageValueAmountCurrencyCode = in.DeclaredCarriageValueAmountCurrencyCode
	billOfLadingD.OtherInstruction = in.OtherInstruction
	billOfLadingD.ConsignorPartyId = in.ConsignorPartyId
	billOfLadingD.CarrierPartyId = in.CarrierPartyId
	billOfLadingD.FreightForwarderPartyId = in.FreightForwarderPartyId
	billOfLadingD.ShipmentId = in.ShipmentId

	billOfLadingT := logisticsproto.BillOfLadingT{}
	billOfLadingT.IssueDate = common.TimeToTimestamp(issueDate.UTC().Truncate(time.Second))

	crUpdUser := commonproto.CrUpdUser{}
	crUpdUser.StatusCode = "active"
	crUpdUser.CreatedByUserId = user.Id
	crUpdUser.UpdatedByUserId = user.Id

	crUpdTime := commonproto.CrUpdTime{}
	crUpdTime.CreatedAt = tn
	crUpdTime.UpdatedAt = tn

	billOfLading := logisticsproto.BillOfLading{BillOfLadingD: &billOfLadingD, BillOfLadingT: &billOfLadingT, CrUpdUser: &crUpdUser, CrUpdTime: &crUpdTime}

	err = ss.insertBillOfLading(ctx, insertBillOfLadingSQL, &billOfLading, in.GetUserEmail(), in.GetRequestId())
	if err != nil {
		ss.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	billOfLadingResponse := logisticsproto.CreateBillOfLadingResponse{}
	billOfLadingResponse.BillOfLading = &billOfLading
	return &billOfLadingResponse, nil
}

// insertBillOfLading - Insert BillOfLading details into database
func (ss *ShipmentService) insertBillOfLading(ctx context.Context, insertBillOfLadingSQL string, billOfLading *logisticsproto.BillOfLading, userEmail string, requestID string) error {
	billOfLadingTmp, err := ss.crBillOfLadingStruct(ctx, billOfLading, userEmail, requestID)
	if err != nil {
		ss.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
		return err
	}

	err = ss.DBService.InsUpd(ctx, userEmail, requestID, func(tx *sqlx.Tx) error {
		res, err := tx.NamedExecContext(ctx, insertBillOfLadingSQL, billOfLadingTmp)
		if err != nil {
			ss.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}

		uID, err := res.LastInsertId()
		if err != nil {
			ss.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}
		billOfLading.BillOfLadingD.Id = uint32(uID)
		uuid4Str, err := common.UUIDBytesToStr(billOfLading.BillOfLadingD.Uuid4)
		if err != nil {
			ss.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}
		billOfLading.BillOfLadingD.IdS = uuid4Str
		return nil
	})

	if err != nil {
		ss.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
		return err
	}
	return nil
}

// crBillOfLadingStruct - process BillOfLading details
func (ss *ShipmentService) crBillOfLadingStruct(ctx context.Context, billOfLading *logisticsproto.BillOfLading, userEmail string, requestID string) (*logisticsstruct.BillOfLading, error) {
	billOfLadingT := new(logisticsstruct.BillOfLadingT)
	billOfLadingT.IssueDate = common.TimestampToTime(billOfLading.BillOfLadingT.IssueDate)

	crUpdTime := new(commonstruct.CrUpdTime)
	crUpdTime.CreatedAt = common.TimestampToTime(billOfLading.CrUpdTime.CreatedAt)
	crUpdTime.UpdatedAt = common.TimestampToTime(billOfLading.CrUpdTime.UpdatedAt)

	billOfLadingTmp := logisticsstruct.BillOfLading{BillOfLadingD: billOfLading.BillOfLadingD, BillOfLadingT: billOfLadingT, CrUpdUser: billOfLading.CrUpdUser, CrUpdTime: crUpdTime}

	return &billOfLadingTmp, nil
}

// CreateWaybill - Create Waybill
func (ss *ShipmentService) CreateWaybill(ctx context.Context, in *logisticsproto.CreateWaybillRequest) (*logisticsproto.CreateWaybillResponse, error) {
	user, err := partyservice.GetUserWithNewContext(ctx, in.UserId, in.UserEmail, in.RequestId, ss.UserServiceClient)
	if err != nil {
		ss.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	ttime := common.GetTimeDetails()
	tn := common.TimeToTimestamp(ttime)

	issueDate, err := time.Parse(common.Layout, in.IssueDate)
	if err != nil {
		ss.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	waybillD := logisticsproto.WaybillD{}
	waybillD.Uuid4, err = common.GetUUIDBytes()
	if err != nil {
		ss.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	waybillD.WaybillId = in.WaybillId
	waybillD.CarrierAssignedId = in.CarrierAssignedId
	waybillD.Name1 = in.Name1
	waybillD.Description = in.Description
	waybillD.Note = in.Note
	waybillD.ShippingOrderId = in.ShippingOrderId
	waybillD.AdValoremIndicator = in.AdValoremIndicator
	waybillD.DeclaredCarriageValueAmount = in.DeclaredCarriageValueAmount
	waybillD.DeclaredCarriageValueAmountCurrencyCode = in.DeclaredCarriageValueAmountCurrencyCode
	waybillD.OtherInstruction = in.OtherInstruction
	waybillD.ConsignorPartyId = in.ConsignorPartyId
	waybillD.CarrierPartyId = in.CarrierPartyId
	waybillD.FreightForwarderPartyId = in.FreightForwarderPartyId
	waybillD.ShipmentId = in.ShipmentId

	waybillT := logisticsproto.WaybillT{}
	waybillT.IssueDate = common.TimeToTimestamp(issueDate.UTC().Truncate(time.Second))

	crUpdUser := commonproto.CrUpdUser{}
	crUpdUser.StatusCode = "active"
	crUpdUser.CreatedByUserId = user.Id
	crUpdUser.UpdatedByUserId = user.Id

	crUpdTime := commonproto.CrUpdTime{}
	crUpdTime.CreatedAt = tn
	crUpdTime.UpdatedAt = tn

	waybill := logisticsproto.Waybill{WaybillD: &waybillD, WaybillT: &waybillT, CrUpdUser: &crUpdUser, CrUpdTime: &crUpdTime}

	err = ss.insertWaybill(ctx, insertWaybillSQL, &waybill, in.GetUserEmail(), in.GetRequestId())
	if err != nil {
		ss.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	waybillResponse := logisticsproto.CreateWaybillResponse{}
	waybillResponse.Waybill = &waybill
	return &waybillResponse, nil
}

// insertWaybill - Insert Waybill details into database
func (ss *ShipmentService) insertWaybill(ctx context.Context, insertWaybillSQL string, waybill *logisticsproto.Waybill, userEmail string, requestID string) error {
	waybillTmp, err := ss.crWaybillStruct(ctx, waybill, userEmail, requestID)
	if err != nil {
		ss.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
		return err
	}

	err = ss.DBService.InsUpd(ctx, userEmail, requestID, func(tx *sqlx.Tx) error {
		res, err := tx.NamedExecContext(ctx, insertWaybillSQL, waybillTmp)
		if err != nil {
			ss.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}

		uID, err := res.LastInsertId()
		if err != nil {
			ss.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}
		waybill.WaybillD.Id = uint32(uID)
		uuid4Str, err := common.UUIDBytesToStr(waybill.WaybillD.Uuid4)
		if err != nil {
			ss.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}
		waybill.WaybillD.IdS = uuid4Str
		return nil
	})

	if err != nil {
		ss.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
		return err
	}
	return nil
}

// crWaybillStruct - process Waybill details
func (ss *ShipmentService) crWaybillStruct(ctx context.Context, waybill *logisticsproto.Waybill, userEmail string, requestID string) (*logisticsstruct.Waybill, error) {
	waybillT := new(logisticsstruct.WaybillT)
	waybillT.IssueDate = common.TimestampToTime(waybill.WaybillT.IssueDate)

	crUpdTime := new(commonstruct.CrUpdTime)
	crUpdTime.CreatedAt = common.TimestampToTime(waybill.CrUpdTime.CreatedAt)
	crUpdTime.UpdatedAt = common.TimestampToTime(waybill.CrUpdTime.UpdatedAt)

	waybillTmp := logisticsstruct.Waybill{WaybillD: waybill.WaybillD, WaybillT: waybillT, CrUpdUser: waybill.CrUpdUser, CrUpdTime: crUpdTime}

	return &waybillTmp, nil
}
