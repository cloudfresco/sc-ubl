package v1

import (
	commonproto "github.com/cloudfresco/sc-ubl/internal/protogen/common/v1"
	logisticsproto "github.com/cloudfresco/sc-ubl/internal/protogen/logistics/v1"
	commonstruct "github.com/cloudfresco/sc-ubl/internal/servicestructs/common/v1"
)

// Consignment - struct Consignment
type Consignment struct {
	*logisticsproto.ConsignmentD
	*commonproto.CrUpdUser
	*commonstruct.CrUpdTime
}
