package v1

import (
	"time"

	commonproto "github.com/cloudfresco/sc-ubl/internal/protogen/common/v1"
	orderproto "github.com/cloudfresco/sc-ubl/internal/protogen/order/v1"
	commonstruct "github.com/cloudfresco/sc-ubl/internal/servicestructs/common/v1"
)

// PurchaseOrderHeader - struct PurchaseOrderHeader
type PurchaseOrderHeader struct {
	*orderproto.PurchaseOrderHeaderD
	*PurchaseOrderHeaderT
	*commonproto.CrUpdUser
	*commonstruct.CrUpdTime
}

// PurchaseOrderHeaderT - struct PurchaseOrderHeaderT
type PurchaseOrderHeaderT struct {
	IssueDate      time.Time `protobuf:"bytes,1,opt,name=issue_date,json=issueDate,proto3" json:"issue_date,omitempty"`
	ValidityPeriod time.Time `protobuf:"bytes,2,opt,name=validity_period,json=validityPeriod,proto3" json:"validity_period,omitempty"`
	TaxExDate      time.Time `protobuf:"bytes,3,opt,name=tax_ex_date,json=taxExDate,proto3" json:"tax_ex_date,omitempty"`
	PricingExDate  time.Time `protobuf:"bytes,4,opt,name=pricing_ex_date,json=pricingExDate,proto3" json:"pricing_ex_date,omitempty"`
	PaymentExDate  time.Time `protobuf:"bytes,5,opt,name=payment_ex_date,json=paymentExDate,proto3" json:"payment_ex_date,omitempty"`
}

// PurchaseOrderLine - struct PurchaseOrderLine
type PurchaseOrderLine struct {
	*orderproto.PurchaseOrderLineD
	*PurchaseOrderLineT
	*commonproto.CrUpdUser
	*commonstruct.CrUpdTime
}

// PurchaseOrderLineT - struct PurchaseOrderLineT
type PurchaseOrderLineT struct {
	PriceValidityPeriodStartDate time.Time `protobuf:"bytes,1,opt,name=price_validity_period_start_date,json=priceValidityPeriodStartDate,proto3" json:"price_validity_period_start_date,omitempty"`
	PriceValidityPeriodEndDate   time.Time `protobuf:"bytes,2,opt,name=price_validity_period_end_date,json=priceValidityPeriodEndDate,proto3" json:"price_validity_period_end_date,omitempty"`
}
