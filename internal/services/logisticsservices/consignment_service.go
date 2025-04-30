package logisticsservices

import (
	"context"

	"github.com/cloudfresco/sc-ubl/internal/common"
	commonproto "github.com/cloudfresco/sc-ubl/internal/protogen/common/v1"
	logisticsproto "github.com/cloudfresco/sc-ubl/internal/protogen/logistics/v1"
	partyproto "github.com/cloudfresco/sc-ubl/internal/protogen/party/v1"
	partyservice "github.com/cloudfresco/sc-ubl/internal/services/partyservices"
	commonstruct "github.com/cloudfresco/sc-ubl/internal/servicestructs/common/v1"
	logisticsstruct "github.com/cloudfresco/sc-ubl/internal/servicestructs/logistics/v1"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// ConsignmentService - For accessing Consignment services
type ConsignmentService struct {
	log               *zap.Logger
	DBService         *common.DBService
	RedisService      *common.RedisService
	UserServiceClient partyproto.UserServiceClient
	logisticsproto.UnimplementedConsignmentServiceServer
}

// NewConsignmentService - Create Consignment service
func NewConsignmentService(log *zap.Logger, dbOpt *common.DBService, redisOpt *common.RedisService, userServiceClient partyproto.UserServiceClient) *ConsignmentService {
	return &ConsignmentService{
		log:               log,
		DBService:         dbOpt,
		RedisService:      redisOpt,
		UserServiceClient: userServiceClient,
	}
}

const insertConsignmentSQL = `insert into consignments
	  ( 
uuid4,
cons_id,
carrier_assigned_id,
consignee_assigned_id,
consignor_assigned_id,
freight_forwarder_assigned_id,
broker_assigned_id,
contracted_carrier_assigned_id,
performing_carrier_assigned_id,
summary_description,
total_invoice_amount,
declared_customs_value_amount,
tariff_description,
tariff_code,
insurance_premium_amount,
gross_weight_measure,
net_weight_measure,
net_net_weight_measure,
chargeable_weight_measure,
gross_volume_measure,
net_volume_measure,
loading_length_measure,
remarks,
hazardous_risk_indicator,
animal_food_indicator,
human_food_indicator,
livestock_indicator,
bulk_cargo_indicator,
containerized_indicator,
general_cargo_indicator,
special_security_indicator,
third_party_payer_indicator,
carrier_service_instructions,
customs_clearance_service_instructions,
forwarder_service_instructions,
special_service_instructions,
sequence_id,
shipping_priority_level_code,
handling_code,
handling_instructions,
information,
total_goods_item_quantity,
total_transport_handling_unit_quantity,
insurance_value_amount,
declared_for_carriage_value_amount,
declared_statistics_value_amount,
free_on_board_value_amount,
special_instructions,
split_consignment_indicator,
delivery_instructions,
consignment_quantity,
consolidatable_indicator,
haulage_instructions,
loading_sequence_id,
child_consignment_quantity,
total_packages_quantity,
consignee_party_id,
exporter_party_id,
consignor_party_id,
importer_party_id,
carrier_party_id,
freight_forwarder_party_id,
notify_party_id,
original_despatch_party_id,
final_delivery_party_id,
performing_carrier_party_id,
substitute_carrier_party_id,
logistics_operator_party_id,
transport_advisor_party_id,
hazardous_item_notification_party_id,
insurance_party_id,
mortgage_holder_party_id,
bill_of_lading_holder_party_id,
original_departure_country_id_code,
original_departure_country_name,
final_destination_country_id_code,
final_destination_country_name,
transit_country_id_code,
transit_country_name,
delivery_terms_id,
payment_terms_id,
collect_payment_terms_id,
disbursement_payment_terms_id,
prepaid_payment_terms_id,
first_arrival_port_address_id,
last_exit_port_location_address_id,
status_code,
created_by_user_id,
updated_by_user_id,
created_at,
updated_at)
  values (:uuid4,
:cons_id,
:carrier_assigned_id,
:consignee_assigned_id,
:consignor_assigned_id,
:freight_forwarder_assigned_id,
:broker_assigned_id,
:contracted_carrier_assigned_id,
:performing_carrier_assigned_id,
:summary_description,
:total_invoice_amount,
:declared_customs_value_amount,
:tariff_description,
:tariff_code,
:insurance_premium_amount,
:gross_weight_measure,
:net_weight_measure,
:net_net_weight_measure,
:chargeable_weight_measure,
:gross_volume_measure,
:net_volume_measure,
:loading_length_measure,
:remarks,
:hazardous_risk_indicator,
:animal_food_indicator,
:human_food_indicator,
:livestock_indicator,
:bulk_cargo_indicator,
:containerized_indicator,
:general_cargo_indicator,
:special_security_indicator,
:third_party_payer_indicator,
:carrier_service_instructions,
:customs_clearance_service_instructions,
:forwarder_service_instructions,
:special_service_instructions,
:sequence_id,
:shipping_priority_level_code,
:handling_code,
:handling_instructions,
:information,
:total_goods_item_quantity,
:total_transport_handling_unit_quantity,
:insurance_value_amount,
:declared_for_carriage_value_amount,
:declared_statistics_value_amount,
:free_on_board_value_amount,
:special_instructions,
:split_consignment_indicator,
:delivery_instructions,
:consignment_quantity,
:consolidatable_indicator,
:haulage_instructions,
:loading_sequence_id,
:child_consignment_quantity,
:total_packages_quantity,
:consignee_party_id,
:exporter_party_id,
:consignor_party_id,
:importer_party_id,
:carrier_party_id,
:freight_forwarder_party_id,
:notify_party_id,
:original_despatch_party_id,
:final_delivery_party_id,
:performing_carrier_party_id,
:substitute_carrier_party_id,
:logistics_operator_party_id,
:transport_advisor_party_id,
:hazardous_item_notification_party_id,
:insurance_party_id,
:mortgage_holder_party_id,
:bill_of_lading_holder_party_id,
:original_departure_country_id_code,
:original_departure_country_name,
:final_destination_country_id_code,
:final_destination_country_name,
:transit_country_id_code,
:transit_country_name,
:delivery_terms_id,
:payment_terms_id,
:collect_payment_terms_id,
:disbursement_payment_terms_id,
:prepaid_payment_terms_id,
:first_arrival_port_address_id,
:last_exit_port_location_address_id,
:status_code,
:created_by_user_id,
:updated_by_user_id,
:created_at,
:updated_at);`

const selectConsignmentsSQL = `select 
id,
uuid4,
cons_id,
carrier_assigned_id,
consignee_assigned_id,
consignor_assigned_id,
freight_forwarder_assigned_id,
broker_assigned_id,
contracted_carrier_assigned_id,
performing_carrier_assigned_id,
summary_description,
total_invoice_amount,
declared_customs_value_amount,
tariff_description,
tariff_code,
insurance_premium_amount,
gross_weight_measure,
net_weight_measure,
net_net_weight_measure,
chargeable_weight_measure,
gross_volume_measure,
net_volume_measure,
loading_length_measure,
remarks,
hazardous_risk_indicator,
animal_food_indicator,
human_food_indicator,
livestock_indicator,
bulk_cargo_indicator,
containerized_indicator,
general_cargo_indicator,
special_security_indicator,
third_party_payer_indicator,
carrier_service_instructions,
customs_clearance_service_instructions,
forwarder_service_instructions,
special_service_instructions,
sequence_id,
shipping_priority_level_code,
handling_code,
handling_instructions,
information,
total_goods_item_quantity,
total_transport_handling_unit_quantity,
insurance_value_amount,
declared_for_carriage_value_amount,
declared_statistics_value_amount,
free_on_board_value_amount,
special_instructions,
split_consignment_indicator,
delivery_instructions,
consignment_quantity,
consolidatable_indicator,
haulage_instructions,
loading_sequence_id,
child_consignment_quantity,
total_packages_quantity,
consignee_party_id,
exporter_party_id,
consignor_party_id,
importer_party_id,
carrier_party_id,
freight_forwarder_party_id,
notify_party_id,
original_despatch_party_id,
final_delivery_party_id,
performing_carrier_party_id,
substitute_carrier_party_id,
logistics_operator_party_id,
transport_advisor_party_id,
hazardous_item_notification_party_id,
insurance_party_id,
mortgage_holder_party_id,
bill_of_lading_holder_party_id,
original_departure_country_id_code,
original_departure_country_name,
final_destination_country_id_code,
final_destination_country_name,
transit_country_id_code,
transit_country_name,
delivery_terms_id,
payment_terms_id,
collect_payment_terms_id,
disbursement_payment_terms_id,
prepaid_payment_terms_id,
first_arrival_port_address_id,
last_exit_port_location_address_id,
status_code,
created_by_user_id,
updated_by_user_id,
created_at,
updated_at from consignments`

// CreateConsignment - Create Consignment
func (cs *ConsignmentService) CreateConsignment(ctx context.Context, in *logisticsproto.CreateConsignmentRequest) (*logisticsproto.CreateConsignmentResponse, error) {
	user, err := partyservice.GetUserWithNewContext(ctx, in.UserId, in.UserEmail, in.RequestId, cs.UserServiceClient)
	if err != nil {
		cs.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	ttime := common.GetTimeDetails()
	tn := common.TimeToTimestamp(ttime)

	consignmentD := logisticsproto.ConsignmentD{}
	consignmentD.Uuid4, err = common.GetUUIDBytes()
	if err != nil {
		cs.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	consignmentD.ConsId = in.ConsId
	consignmentD.CarrierAssignedId = in.CarrierAssignedId
	consignmentD.ConsigneeAssignedId = in.ConsigneeAssignedId
	consignmentD.ConsignorAssignedId = in.ConsignorAssignedId
	consignmentD.FreightForwarderAssignedId = in.FreightForwarderAssignedId
	consignmentD.BrokerAssignedId = in.BrokerAssignedId
	consignmentD.ContractedCarrierAssignedId = in.ContractedCarrierAssignedId
	consignmentD.PerformingCarrierAssignedId = in.PerformingCarrierAssignedId
	consignmentD.SummaryDescription = in.SummaryDescription
	consignmentD.TotalInvoiceAmount = in.TotalInvoiceAmount
	consignmentD.DeclaredCustomsValueAmount = in.DeclaredCustomsValueAmount
	consignmentD.TariffDescription = in.TariffDescription
	consignmentD.TariffCode = in.TariffCode
	consignmentD.InsurancePremiumAmount = in.InsurancePremiumAmount
	consignmentD.GrossWeightMeasure = in.GrossWeightMeasure
	consignmentD.NetWeightMeasure = in.NetWeightMeasure
	consignmentD.NetNetWeightMeasure = in.NetNetWeightMeasure
	consignmentD.ChargeableWeightMeasure = in.ChargeableWeightMeasure
	consignmentD.GrossVolumeMeasure = in.GrossVolumeMeasure
	consignmentD.NetVolumeMeasure = in.NetVolumeMeasure
	consignmentD.LoadingLengthMeasure = in.LoadingLengthMeasure
	consignmentD.Remarks = in.Remarks
	consignmentD.HazardousRiskIndicator = in.HazardousRiskIndicator
	consignmentD.AnimalFoodIndicator = in.AnimalFoodIndicator
	consignmentD.HumanFoodIndicator = in.HumanFoodIndicator
	consignmentD.LivestockIndicator = in.LivestockIndicator
	consignmentD.BulkCargoIndicator = in.BulkCargoIndicator
	consignmentD.ContainerizedIndicator = in.ContainerizedIndicator
	consignmentD.GeneralCargoIndicator = in.GeneralCargoIndicator
	consignmentD.SpecialSecurityIndicator = in.SpecialSecurityIndicator
	consignmentD.ThirdPartyPayerIndicator = in.ThirdPartyPayerIndicator
	consignmentD.CarrierServiceInstructions = in.CarrierServiceInstructions
	consignmentD.CustomsClearanceServiceInstructions = in.CustomsClearanceServiceInstructions
	consignmentD.ForwarderServiceInstructions = in.ForwarderServiceInstructions
	consignmentD.SpecialServiceInstructions = in.SpecialServiceInstructions
	consignmentD.SequenceId = in.SequenceId
	consignmentD.ShippingPriorityLevelCode = in.ShippingPriorityLevelCode
	consignmentD.HandlingCode = in.HandlingCode
	consignmentD.HandlingInstructions = in.HandlingInstructions
	consignmentD.Information = in.Information
	consignmentD.TotalGoodsItemQuantity = in.TotalGoodsItemQuantity
	consignmentD.TotalTransportHandlingUnitQuantity = in.TotalTransportHandlingUnitQuantity
	consignmentD.InsuranceValueAmount = in.InsuranceValueAmount
	consignmentD.DeclaredForCarriageValueAmount = in.DeclaredForCarriageValueAmount
	consignmentD.DeclaredStatisticsValueAmount = in.DeclaredStatisticsValueAmount
	consignmentD.FreeOnBoardValueAmount = in.FreeOnBoardValueAmount
	consignmentD.SpecialInstructions = in.SpecialInstructions
	consignmentD.SplitConsignmentIndicator = in.SplitConsignmentIndicator
	consignmentD.DeliveryInstructions = in.DeliveryInstructions
	consignmentD.ConsignmentQuantity = in.ConsignmentQuantity
	consignmentD.ConsolidatableIndicator = in.ConsolidatableIndicator
	consignmentD.HaulageInstructions = in.HaulageInstructions
	consignmentD.LoadingSequenceId = in.LoadingSequenceId
	consignmentD.ChildConsignmentQuantity = in.ChildConsignmentQuantity
	consignmentD.TotalPackagesQuantity = in.TotalPackagesQuantity
	consignmentD.ConsigneePartyId = in.ConsigneePartyId
	consignmentD.ExporterPartyId = in.ExporterPartyId
	consignmentD.ConsignorPartyId = in.ConsignorPartyId
	consignmentD.ImporterPartyId = in.ImporterPartyId
	consignmentD.CarrierPartyId = in.CarrierPartyId
	consignmentD.FreightForwarderPartyId = in.FreightForwarderPartyId
	consignmentD.NotifyPartyId = in.NotifyPartyId
	consignmentD.OriginalDespatchPartyId = in.OriginalDespatchPartyId
	consignmentD.FinalDeliveryPartyId = in.FinalDeliveryPartyId
	consignmentD.PerformingCarrierPartyId = in.PerformingCarrierPartyId
	consignmentD.SubstituteCarrierPartyId = in.SubstituteCarrierPartyId
	consignmentD.LogisticsOperatorPartyId = in.LogisticsOperatorPartyId
	consignmentD.TransportAdvisorPartyId = in.TransportAdvisorPartyId
	consignmentD.HazardousItemNotificationPartyId = in.HazardousItemNotificationPartyId
	consignmentD.InsurancePartyId = in.InsurancePartyId
	consignmentD.MortgageHolderPartyId = in.MortgageHolderPartyId
	consignmentD.BillOfLadingHolderPartyId = in.BillOfLadingHolderPartyId
	consignmentD.OriginalDepartureCountryIdCode = in.OriginalDepartureCountryIdCode
	consignmentD.OriginalDepartureCountryName = in.OriginalDepartureCountryName
	consignmentD.FinalDestinationCountryIdCode = in.FinalDestinationCountryIdCode
	consignmentD.FinalDestinationCountryName = in.FinalDestinationCountryName
	consignmentD.TransitCountryIdCode = in.TransitCountryIdCode
	consignmentD.TransitCountryName = in.TransitCountryName
	consignmentD.DeliveryTermsId = in.DeliveryTermsId
	consignmentD.PaymentTermsId = in.PaymentTermsId
	consignmentD.CollectPaymentTermsId = in.CollectPaymentTermsId
	consignmentD.DisbursementPaymentTermsId = in.DisbursementPaymentTermsId
	consignmentD.PrepaidPaymentTermsId = in.PrepaidPaymentTermsId
	consignmentD.FirstArrivalPortAddressId = in.FirstArrivalPortAddressId
	consignmentD.LastExitPortLocationAddressId = in.LastExitPortLocationAddressId

	crUpdUser := commonproto.CrUpdUser{}
	crUpdUser.StatusCode = "active"
	crUpdUser.CreatedByUserId = user.Id
	crUpdUser.UpdatedByUserId = user.Id

	crUpdTime := commonproto.CrUpdTime{}
	crUpdTime.CreatedAt = tn
	crUpdTime.UpdatedAt = tn

	consignment := logisticsproto.Consignment{ConsignmentD: &consignmentD, CrUpdUser: &crUpdUser, CrUpdTime: &crUpdTime}

	err = cs.insertConsignment(ctx, insertConsignmentSQL, &consignment, in.GetUserEmail(), in.GetRequestId())
	if err != nil {
		cs.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	consignmentResponse := logisticsproto.CreateConsignmentResponse{}
	consignmentResponse.Consignment = &consignment
	return &consignmentResponse, nil
}

// insertDespatch - Insert Consignment
func (cs *ConsignmentService) insertConsignment(ctx context.Context, insertConsignmentSQL string, consignment *logisticsproto.Consignment, userEmail string, requestID string) error {
	consignmentTmp, err := cs.crConsignmentStruct(ctx, consignment, userEmail, requestID)
	if err != nil {
		cs.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
		return err
	}

	err = cs.DBService.InsUpd(ctx, userEmail, requestID, func(tx *sqlx.Tx) error {
		res, err := tx.NamedExecContext(ctx, insertConsignmentSQL, consignmentTmp)
		if err != nil {
			cs.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}

		uID, err := res.LastInsertId()
		if err != nil {
			cs.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}
		consignment.ConsignmentD.Id = uint32(uID)
		uuid4Str, err := common.UUIDBytesToStr(consignment.ConsignmentD.Uuid4)
		if err != nil {
			cs.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}
		consignment.ConsignmentD.IdS = uuid4Str
		return nil
	})

	if err != nil {
		cs.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
		return err
	}
	return nil
}

// crConsignmentStruct - process Consignment details
func (cs *ConsignmentService) crConsignmentStruct(ctx context.Context, consignment *logisticsproto.Consignment, userEmail string, requestID string) (*logisticsstruct.Consignment, error) {
	crUpdTime := new(commonstruct.CrUpdTime)
	crUpdTime.CreatedAt = common.TimestampToTime(consignment.CrUpdTime.CreatedAt)
	crUpdTime.UpdatedAt = common.TimestampToTime(consignment.CrUpdTime.UpdatedAt)

	consignmentTmp := logisticsstruct.Consignment{ConsignmentD: consignment.ConsignmentD, CrUpdUser: consignment.CrUpdUser, CrUpdTime: crUpdTime}

	return &consignmentTmp, nil
}

// GetConsignments - Get Consignments
func (cs *ConsignmentService) GetConsignments(ctx context.Context, in *logisticsproto.GetConsignmentsRequest) (*logisticsproto.GetConsignmentsResponse, error) {
	limit := in.GetLimit()
	nextCursor := in.GetNextCursor()
	if limit == "" {
		limit = cs.DBService.LimitSQLRows
	}
	query := "levelp = ? and status_code = ?"
	if nextCursor == "" {
		query = query + " order by id desc " + " limit " + limit + ";"
	} else {
		nextCursor = common.DecodeCursor(nextCursor)
		query = query + " " + "and" + " " + "id <= " + nextCursor + " order by id desc " + " limit " + limit + ";"
	}

	consignments := []*logisticsproto.Consignment{}

	nselectConsignmentsSQL := selectConsignmentsSQL + ` where ` + query

	rows, err := cs.DBService.DB.QueryxContext(ctx, nselectConsignmentsSQL, 0, "active")
	if err != nil {
		cs.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	for rows.Next() {

		consignmentTmp := logisticsstruct.Consignment{}
		err = rows.StructScan(&consignmentTmp)
		if err != nil {
			cs.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
			return nil, err
		}

		getRequest := commonproto.GetRequest{}
		getRequest.UserEmail = in.UserEmail
		getRequest.RequestId = in.RequestId
		consignment, err := cs.getConsignmentStruct(ctx, &getRequest, consignmentTmp)
		if err != nil {
			cs.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
			return nil, err
		}

		consignments = append(consignments, consignment)

	}

	consignmentsResponse := logisticsproto.GetConsignmentsResponse{}
	if len(consignments) != 0 {
		next := consignments[len(consignments)-1].ConsignmentD.Id
		next--
		nextc := common.EncodeCursor(next)
		consignmentsResponse = logisticsproto.GetConsignmentsResponse{Consignments: consignments, NextCursor: nextc}
	} else {
		consignmentsResponse = logisticsproto.GetConsignmentsResponse{Consignments: consignments, NextCursor: "0"}
	}
	return &consignmentsResponse, nil
}

// GetConsignment - Get Consignment
func (cs *ConsignmentService) GetConsignment(ctx context.Context, inReq *logisticsproto.GetConsignmentRequest) (*logisticsproto.GetConsignmentResponse, error) {
	in := inReq.GetRequest
	uuid4byte, err := common.UUIDStrToBytes(in.Id)
	if err != nil {
		cs.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	nselectConsignmentsSQL := selectConsignmentsSQL + ` where uuid4 = ? and status_code = ?;`
	row := cs.DBService.DB.QueryRowxContext(ctx, nselectConsignmentsSQL, uuid4byte, "active")
	consignmentTmp := logisticsstruct.Consignment{}
	err = row.StructScan(&consignmentTmp)
	if err != nil {
		cs.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	consignment, err := cs.getConsignmentStruct(ctx, in, consignmentTmp)
	if err != nil {
		cs.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	consignmentResponse := logisticsproto.GetConsignmentResponse{}
	consignmentResponse.Consignment = consignment
	return &consignmentResponse, nil
}

// GetConsignmentByPk - Get Consignment By Primary key(Id)
func (cs *ConsignmentService) GetConsignmentByPk(ctx context.Context, inReq *logisticsproto.GetConsignmentByPkRequest) (*logisticsproto.GetConsignmentByPkResponse, error) {
	in := inReq.GetByIdRequest
	nselectConsignmentsSQL := selectConsignmentsSQL + ` where id = ? and status_code = ?;`
	row := cs.DBService.DB.QueryRowxContext(ctx, nselectConsignmentsSQL, in.Id, "active")
	consignmentTmp := logisticsstruct.Consignment{}
	err := row.StructScan(&consignmentTmp)
	if err != nil {
		cs.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	getRequest := commonproto.GetRequest{}
	getRequest.UserEmail = in.UserEmail
	getRequest.RequestId = in.RequestId
	consignment, err := cs.getConsignmentStruct(ctx, &getRequest, consignmentTmp)
	if err != nil {
		cs.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	consignmentResponse := logisticsproto.GetConsignmentByPkResponse{}
	consignmentResponse.Consignment = consignment
	return &consignmentResponse, nil
}

// getConsignmentStruct - Get Consignment header
func (cs *ConsignmentService) getConsignmentStruct(ctx context.Context, in *commonproto.GetRequest, consignmentTmp logisticsstruct.Consignment) (*logisticsproto.Consignment, error) {
	uuid4Str, err := common.UUIDBytesToStr(consignmentTmp.ConsignmentD.Uuid4)
	if err != nil {
		cs.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	consignmentTmp.ConsignmentD.IdS = uuid4Str

	crUpdTime := new(commonproto.CrUpdTime)
	crUpdTime.CreatedAt = common.TimeToTimestamp(consignmentTmp.CrUpdTime.CreatedAt)
	crUpdTime.UpdatedAt = common.TimeToTimestamp(consignmentTmp.CrUpdTime.UpdatedAt)

	consignment := logisticsproto.Consignment{ConsignmentD: consignmentTmp.ConsignmentD, CrUpdUser: consignmentTmp.CrUpdUser, CrUpdTime: crUpdTime}

	return &consignment, nil
}
