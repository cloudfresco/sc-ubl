package v1

import (
	commonproto "github.com/cloudfresco/sc-ubl/internal/protogen/common/v1"
	partyproto "github.com/cloudfresco/sc-ubl/internal/protogen/party/v1"
	commonstruct "github.com/cloudfresco/sc-ubl/internal/servicestructs/common/v1"
)

// Party - struct Party
type Party struct {
	*partyproto.PartyD
	*commonproto.PartyLegalEntityD
	*commonstruct.PartyLegalEntityT
	*commonproto.CrUpdUser
	*commonstruct.CrUpdTime
}

// PartyChd - struct PartyChd
type PartyChd struct {
	*partyproto.PartyChdD
	*commonproto.CrUpdUser
	*commonstruct.CrUpdTime
}

// PartyContact - struct PartyContact
type PartyContact struct {
	*partyproto.PartyContactD
	*commonproto.PartyInfo
	*commonproto.CrUpdUser
	*commonstruct.CrUpdTime
}

// PartySocialProfile - struct PartySocialProfile
type PartySocialProfile struct {
	*partyproto.PartySocialProfileD
	*commonproto.PartyInfo
	*commonproto.CrUpdUser
	*commonstruct.CrUpdTime
}

// PartyCorporateJurisdiction - struct PartyCorporateJurisdiction
type PartyCorporateJurisdiction struct {
	*partyproto.PartyCorporateJurisdictionD
	*commonproto.PartyInfo
	*commonproto.CrUpdUser
	*commonstruct.CrUpdTime
}
