package v1

import (
	commonproto "github.com/cloudfresco/sc-ubl/internal/protogen/common/v1"
	taxproto "github.com/cloudfresco/sc-ubl/internal/protogen/tax/v1"
	commonstruct "github.com/cloudfresco/sc-ubl/internal/servicestructs/common/v1"
)

// TaxScheme - struct TaxScheme
type TaxScheme struct {
	*taxproto.TaxSchemeD
	*commonproto.CrUpdUser
	*commonstruct.CrUpdTime
}

// TaxCategory - struct TaxCategory
type TaxCategory struct {
	*taxproto.TaxCategoryD
	*commonproto.CrUpdUser
	*commonstruct.CrUpdTime
}

// TaxSubTotal - struct TaxSubTotal
type TaxSubTotal struct {
	*taxproto.TaxSubTotalD
	*commonproto.CrUpdUser
	*commonstruct.CrUpdTime
}

// TaxTotal - struct TaxTotal
type TaxTotal struct {
	*taxproto.TaxTotalD
	*commonproto.CrUpdUser
	*commonstruct.CrUpdTime
}

// TaxSchemeJurisdiction - struct TaxSchemeJurisdiction
type TaxSchemeJurisdiction struct {
	*taxproto.TaxSchemeJurisdictionD
	*commonproto.TaxSchemeInfo
	*commonproto.CrUpdUser
	*commonstruct.CrUpdTime
}
