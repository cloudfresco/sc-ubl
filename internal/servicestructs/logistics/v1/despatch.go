package v1

import (
	"time"

	commonproto "github.com/cloudfresco/sc-ubl/internal/protogen/common/v1"
	logisticsproto "github.com/cloudfresco/sc-ubl/internal/protogen/logistics/v1"
	commonstruct "github.com/cloudfresco/sc-ubl/internal/servicestructs/common/v1"
)

// DespatchHeader - struct DespatchHeader
type DespatchHeader struct {
	*logisticsproto.DespatchHeaderD
	*DespatchHeaderT
	*commonproto.CrUpdUser
	*commonstruct.CrUpdTime
}

// DespatchHeaderT - struct DespatchHeaderT
type DespatchHeaderT struct {
	IssueDate time.Time `protobuf:"bytes,1,opt,name=issue_date,json=issueDate,proto3" json:"issue_date,omitempty"`
}

// DespatchLine - struct DespatchLine
type DespatchLine struct {
	*logisticsproto.DespatchLineD
	*commonproto.CrUpdUser
	*commonstruct.CrUpdTime
}
