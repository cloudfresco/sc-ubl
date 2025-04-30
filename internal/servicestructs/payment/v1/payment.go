package v1

import (
	"time"

	commonproto "github.com/cloudfresco/sc-ubl/internal/protogen/common/v1"
	paymentproto "github.com/cloudfresco/sc-ubl/internal/protogen/payment/v1"
	commonstruct "github.com/cloudfresco/sc-ubl/internal/servicestructs/common/v1"
)

// Payment - struct Payment
type Payment struct {
	*paymentproto.PaymentD
	*PaymentT
	*commonproto.CrUpdUser
	*commonstruct.CrUpdTime
}

// PaymentT - struct PaymentT
type PaymentT struct {
	ReceivedDate time.Time `protobuf:"bytes,1,opt,name=received_date,json=receivedDate,proto3" json:"received_date,omitempty"`
	PaidDate     time.Time `protobuf:"bytes,2,opt,name=paid_date,json=paidDate,proto3" json:"paid_date,omitempty"`
}

// PaymentMandate - struct PaymentMandate
type PaymentMandate struct {
	*paymentproto.PaymentMandateD
	*PaymentMandateT
	*commonproto.CrUpdUser
	*commonstruct.CrUpdTime
}

// PaymentMandateT - struct PaymentMandateT
type PaymentMandateT struct {
	ValidityPeriodStartDate        time.Time `protobuf:"bytes,1,opt,name=validity_period_start_date,json=validityPeriodStartDate,proto3" json:"validity_period_start_date,omitempty"`
	ValidityPeriodEndDate          time.Time `protobuf:"bytes,2,opt,name=validity_period_end_date,json=validityPeriodEndDate,proto3" json:"validity_period_end_date,omitempty"`
	PaymentReversalPeriodStartDate time.Time `protobuf:"bytes,3,opt,name=payment_reversal_period_start_date,json=paymentReversalPeriodStartDate,proto3" json:"payment_reversal_period_start_date,omitempty"`
	PaymentReversalPeriodEndDate   time.Time `protobuf:"bytes,4,opt,name=payment_reversal_period_end_date,json=paymentReversalPeriodEndDate,proto3" json:"payment_reversal_period_end_date,omitempty"`
}

// PaymentMean - struct PaymentMean
type PaymentMean struct {
	*paymentproto.PaymentMeanD
	*PaymentMeanT
	*commonproto.CrUpdUser
	*commonstruct.CrUpdTime
}

// PaymentMeanT - struct PaymentMeanT
type PaymentMeanT struct {
	PaymentDueDate time.Time `protobuf:"bytes,1,opt,name=payment_due_date,json=paymentDueDate,proto3" json:"payment_due_date,omitempty"`
}

// PaymentTerm - struct PaymentTerm
type PaymentTerm struct {
	*paymentproto.PaymentTermD
	*PaymentTermT
	*commonproto.CrUpdUser
	*commonstruct.CrUpdTime
}

// PaymentTermT - struct PaymentTermT
type PaymentTermT struct {
	PaymentDueDate            time.Time `protobuf:"bytes,1,opt,name=payment_due_date,json=paymentDueDate,proto3" json:"payment_due_date,omitempty"`
	InstallmentDueDate        time.Time `protobuf:"bytes,2,opt,name=installment_due_date,json=installmentDueDate,proto3" json:"installment_due_date,omitempty"`
	SettlementPeriodStartDate time.Time `protobuf:"bytes,3,opt,name=settlement_period_start_date,json=settlementPeriodStartDate,proto3" json:"settlement_period_start_date,omitempty"`
	SettlementPeriodEndDate   time.Time `protobuf:"bytes,4,opt,name=settlement_period_end_date,json=settlementPeriodEndDate,proto3" json:"settlement_period_end_date,omitempty"`
	PenaltyPeriodStartDate    time.Time `protobuf:"bytes,5,opt,name=penalty_period_start_date,json=penaltyPeriodStartDate,proto3" json:"penalty_period_start_date,omitempty"`
	PenaltyPeriodEndDate      time.Time `protobuf:"bytes,6,opt,name=penalty_period_end_date,json=penaltyPeriodEndDate,proto3" json:"penalty_period_end_date,omitempty"`
	ValidityPeriodStartDate   time.Time `protobuf:"bytes,7,opt,name=validity_period_start_date,json=validityPeriodStartDate,proto3" json:"validity_period_start_date,omitempty"`
	ValidityPeriodEndDate     time.Time `protobuf:"bytes,8,opt,name=validity_period_end_date,json=validityPeriodEndDate,proto3" json:"validity_period_end_date,omitempty"`
}

// PaymentMandateClauseContent - struct PaymentMandateClauseContent
type PaymentMandateClauseContent struct {
	*paymentproto.PaymentMandateClauseContentD
	*commonproto.CrUpdUser
	*commonstruct.CrUpdTime
}
