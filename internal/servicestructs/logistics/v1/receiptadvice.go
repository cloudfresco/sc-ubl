package v1

import (
	"time"

	commonproto "github.com/cloudfresco/sc-ubl/internal/protogen/common/v1"
	logisticsproto "github.com/cloudfresco/sc-ubl/internal/protogen/logistics/v1"
	commonstruct "github.com/cloudfresco/sc-ubl/internal/servicestructs/common/v1"
)

// ReceiptAdviceHeader - struct ReceiptAdviceHeader
type ReceiptAdviceHeader struct {
	*logisticsproto.ReceiptAdviceHeaderD
	*ReceiptAdviceHeaderT
	*commonproto.CrUpdUser
	*commonstruct.CrUpdTime
}

// ReceiptAdviceHeaderT - struct ReceiptAdviceHeaderT
type ReceiptAdviceHeaderT struct {
	IssueDate time.Time `protobuf:"bytes,1,opt,name=issue_date,json=issueDate,proto3" json:"issue_date,omitempty"`
}

// ReceiptAdviceLine - struct ReceiptAdviceLine
type ReceiptAdviceLine struct {
	*logisticsproto.ReceiptAdviceLineD
	*ReceiptAdviceLineT
	*commonproto.CrUpdUser
	*commonstruct.CrUpdTime
}

// ReceiptAdviceLineT - struct ReceiptAdviceLineT
type ReceiptAdviceLineT struct {
	ReceivedDate time.Time `protobuf:"bytes,1,opt,name=received_date,json=receivedDate,proto3" json:"received_date,omitempty"`
}
