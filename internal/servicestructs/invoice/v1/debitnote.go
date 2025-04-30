package v1

import (
	"time"

	commonproto "github.com/cloudfresco/sc-ubl/internal/protogen/common/v1"
	invoiceproto "github.com/cloudfresco/sc-ubl/internal/protogen/invoice/v1"
	commonstruct "github.com/cloudfresco/sc-ubl/internal/servicestructs/common/v1"
)

// DebitNoteHeader - struct DebitNoteHeader
type DebitNoteHeader struct {
	*invoiceproto.DebitNoteHeaderD
	*DebitNoteHeaderT
	*commonproto.CrUpdUser
	*commonstruct.CrUpdTime
}

// DebitNoteHeaderT - struct DebitNoteHeaderT
type DebitNoteHeaderT struct {
	IssueDate              time.Time `protobuf:"bytes,1,opt,name=issue_date,json=issueDate,proto3" json:"issue_date,omitempty"`
	TaxPointDate           time.Time `protobuf:"bytes,2,opt,name=tax_point_date,json=taxPointDate,proto3" json:"tax_point_date,omitempty"`
	InvoicePeriodStartDate time.Time `protobuf:"bytes,3,opt,name=invoice_period_start_date,json=invoicePeriodStartDate,proto3" json:"invoice_period_start_date,omitempty"`
	InvoicePeriodEndDate   time.Time `protobuf:"bytes,4,opt,name=invoice_period_end_date,json=invoicePeriodEndDate,proto3" json:"invoice_period_end_date,omitempty"`
	TaxExDate              time.Time `protobuf:"bytes,5,opt,name=tax_ex_date,json=taxExDate,proto3" json:"tax_ex_date,omitempty"`
	PricingExDate          time.Time `protobuf:"bytes,6,opt,name=pricing_ex_date,json=pricingExDate,proto3" json:"pricing_ex_date,omitempty"`
	PaymentExDate          time.Time `protobuf:"bytes,7,opt,name=payment_ex_date,json=paymentExDate,proto3" json:"payment_ex_date,omitempty"`
	PaymentAltExDate       time.Time `protobuf:"bytes,8,opt,name=payment_alt_ex_date,json=paymentAltExDate,proto3" json:"payment_alt_ex_date,omitempty"`
}

// DebitNoteLine - struct DebitNoteLine
type DebitNoteLine struct {
	*invoiceproto.DebitNoteLineD
	*DebitNoteLineT
	*commonproto.CrUpdUser
	*commonstruct.CrUpdTime
}

// DebitNoteLineT - struct DebitNoteLineT
type DebitNoteLineT struct {
	TaxPointDate                 time.Time `protobuf:"bytes,1,opt,name=tax_point_date,json=taxPointDate,proto3" json:"tax_point_date,omitempty"`
	PriceValidityPeriodStartDate time.Time `protobuf:"bytes,2,opt,name=price_validity_period_start_date,json=priceValidityPeriodStartDate,proto3" json:"price_validity_period_start_date,omitempty"`
	PriceValidityPeriodEndDate   time.Time `protobuf:"bytes,3,opt,name=price_validity_period_end_date,json=priceValidityPeriodEndDate,proto3" json:"price_validity_period_end_date,omitempty"`
}
