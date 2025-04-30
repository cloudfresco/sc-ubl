package v1

import (
	"time"

	commonproto "github.com/cloudfresco/sc-ubl/internal/protogen/common/v1"
	logisticsproto "github.com/cloudfresco/sc-ubl/internal/protogen/logistics/v1"
	commonstruct "github.com/cloudfresco/sc-ubl/internal/servicestructs/common/v1"
)

// Delivery - struct Delivery
type Delivery struct {
	*logisticsproto.DeliveryD
	*DeliveryT
	*commonproto.CrUpdUser
	*commonstruct.CrUpdTime
}

type DeliveryT struct {
	ActualDeliveryDate               time.Time `protobuf:"bytes,1,opt,name=actual_delivery_date,json=actualDeliveryDate,proto3" json:"actual_delivery_date,omitempty"`
	LatestDeliveryDate               time.Time `protobuf:"bytes,2,opt,name=latest_delivery_date,json=latestDeliveryDate,proto3" json:"latest_delivery_date,omitempty"`
	RequestedDeliveryPeriodStartDate time.Time `protobuf:"bytes,3,opt,name=requested_delivery_period_start_date,json=requestedDeliveryPeriodStartDate,proto3" json:"requested_delivery_period_start_date,omitempty"`
	RequestedDeliveryPeriodEndDate   time.Time `protobuf:"bytes,4,opt,name=requested_delivery_period_end_date,json=requestedDeliveryPeriodEndDate,proto3" json:"requested_delivery_period_end_date,omitempty"`
	PromisedDeliveryPeriodStartDate  time.Time `protobuf:"bytes,5,opt,name=promised_delivery_period_start_date,json=promisedDeliveryPeriodStartDate,proto3" json:"promised_delivery_period_start_date,omitempty"`
	PromisedDeliveryPeriodEndDate    time.Time `protobuf:"bytes,6,opt,name=promised_delivery_period_end_date,json=promisedDeliveryPeriodEndDate,proto3" json:"promised_delivery_period_end_date,omitempty"`
	EstimatedDeliveryPeriodStartDate time.Time `protobuf:"bytes,7,opt,name=estimated_delivery_period_start_date,json=estimatedDeliveryPeriodStartDate,proto3" json:"estimated_delivery_period_start_date,omitempty"`
	EstimatedDeliveryPeriodEndDate   time.Time `protobuf:"bytes,8,opt,name=estimated_delivery_period_end_date,json=estimatedDeliveryPeriodEndDate,proto3" json:"estimated_delivery_period_end_date,omitempty"`
}

// DeliveryTerm - struct DeliveryTerm
type DeliveryTerm struct {
	*logisticsproto.DeliveryTermD
	*commonproto.CrUpdUser
	*commonstruct.CrUpdTime
}

// Despatch - struct Despatch
type Despatch struct {
	*logisticsproto.DespatchD
	*DespatchT
	*commonproto.CrUpdUser
	*commonstruct.CrUpdTime
}

// DespatchT - struct DespatchT
type DespatchT struct {
	RequestedDespatchDate            time.Time `protobuf:"bytes,1,opt,name=requested_despatch_date,json=requestedDespatchDate,proto3" json:"requested_despatch_date,omitempty"`
	EstimatedDespatchDate            time.Time `protobuf:"bytes,2,opt,name=estimated_despatch_date,json=estimatedDespatchDate,proto3" json:"estimated_despatch_date,omitempty"`
	ActualDespatchDate               time.Time `protobuf:"bytes,3,opt,name=actual_despatch_date,json=actualDespatchDate,proto3" json:"actual_despatch_date,omitempty"`
	GuaranteedDespatchDate           time.Time `protobuf:"bytes,4,opt,name=guaranteed_despatch_date,json=guaranteedDespatchDate,proto3" json:"guaranteed_despatch_date,omitempty"`
	EstimatedDespatchPeriodStartDate time.Time `protobuf:"bytes,5,opt,name=estimated_despatch_period_start_date,json=estimatedDespatchPeriodStartDate,proto3" json:"estimated_despatch_period_start_date,omitempty"`
	EstimatedDespatchPeriodEndDate   time.Time `protobuf:"bytes,6,opt,name=estimated_despatch_period_end_date,json=estimatedDespatchPeriodEndDate,proto3" json:"estimated_despatch_period_end_date,omitempty"`
	RequestedDespatchPeriodStartDate time.Time `protobuf:"bytes,7,opt,name=requested_despatch_period_start_date,json=requestedDespatchPeriodStartDate,proto3" json:"requested_despatch_period_start_date,omitempty"`
	RequestedDespatchPeriodEndDate   time.Time `protobuf:"bytes,8,opt,name=requested_despatch_period_end_date,json=requestedDespatchPeriodEndDate,proto3" json:"requested_despatch_period_end_date,omitempty"`
}

// Shipment - struct Shipment
type Shipment struct {
	*logisticsproto.ShipmentD
	*commonproto.CrUpdUser
	*commonstruct.CrUpdTime
}

// BillOfLading - struct BillOfLading
type BillOfLading struct {
	*logisticsproto.BillOfLadingD
	*BillOfLadingT
	*commonproto.CrUpdUser
	*commonstruct.CrUpdTime
}

// BillOfLadingT - struct BillOfLadingT
type BillOfLadingT struct {
	IssueDate time.Time `protobuf:"bytes,1,opt,name=issue_date,json=issueDate,proto3" json:"issue_date,omitempty"`
}

// Waybill - struct Waybill
type Waybill struct {
	*logisticsproto.WaybillD
	*WaybillT
	*commonproto.CrUpdUser
	*commonstruct.CrUpdTime
}

// WaybillT - struct WaybillT
type WaybillT struct {
	IssueDate time.Time `protobuf:"bytes,1,opt,name=issue_date,json=issueDate,proto3" json:"issue_date,omitempty"`
}
