package v1

import (
	commonproto "github.com/cloudfresco/sc-ubl/internal/protogen/common/v1"
	partyproto "github.com/cloudfresco/sc-ubl/internal/protogen/party/v1"
	commonstruct "github.com/cloudfresco/sc-ubl/internal/servicestructs/common/v1"
)

// FinancialInstitution - struct FinancialInstitution
type FinancialInstitution struct {
	*partyproto.FinancialInstitutionD
	*commonproto.CrUpdUser
	*commonstruct.CrUpdTime
}

// FinancialInstitutionBranch - struct FinancialInstitutionBranch
type FinancialInstitutionBranch struct {
	*partyproto.FinancialInstitutionBranchD
	*commonproto.CrUpdUser
	*commonstruct.CrUpdTime
}
