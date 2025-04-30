package v1

import (
	"time"

	commonproto "github.com/cloudfresco/sc-ubl/internal/protogen/common/v1"
)

// CrUpdTime - struct CrUpdTime
type CrUpdTime struct {
	CreatedAt time.Time `protobuf:"bytes,1,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt time.Time `protobuf:"bytes,2,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
}

// Location - struct Location
type Location struct {
	*commonproto.LocationD
	*LocationT
}

type LocationT struct {
	ValidityPeriodStartDate time.Time `protobuf:"bytes,1,opt,name=validity_period_start_date,json=validityPeriodStartDate,proto3" json:"validity_period_start_date,omitempty"`
	ValidityPeriodEndDate   time.Time `protobuf:"bytes,2,opt,name=validity_period_end_date,json=validityPeriodEndDate,proto3" json:"validity_period_end_date,omitempty"`
}

// PartyLegalEntity - struct PartyLegalEntity
type PartyLegalEntity struct {
	*commonproto.PartyLegalEntityD
	*PartyLegalEntityT
}

// PartyLegalEntityT - struct PartyLegalEntityT
type PartyLegalEntityT struct {
	RegistrationDate           time.Time `protobuf:"bytes,1,opt,name=registration_date,json=registrationDate,proto3" json:"registration_date,omitempty"`
	RegistrationExpirationDate time.Time `protobuf:"bytes,2,opt,name=registration_expiration_date,json=registrationExpirationDate,proto3" json:"registration_expiration_date,omitempty"`
}
