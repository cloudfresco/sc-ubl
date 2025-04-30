package v1

import (
	"time"

	commonproto "github.com/cloudfresco/sc-ubl/internal/protogen/common/v1"
	itemproto "github.com/cloudfresco/sc-ubl/internal/protogen/item/v1"
	commonstruct "github.com/cloudfresco/sc-ubl/internal/servicestructs/common/v1"
)

// Item - struct Item
type Item struct {
	*itemproto.ItemD
	*commonproto.CrUpdUser
	*commonstruct.CrUpdTime
}

// ItemProperty - struct ItemProperty
type ItemProperty struct {
	*itemproto.ItemPropertyD
	*ItemPropertyT
	*commonproto.CrUpdUser
	*commonstruct.CrUpdTime
}

// ItemPropertyT - struct ItemPropertyT
type ItemPropertyT struct {
	UsabilityPeriodStartDate time.Time `protobuf:"bytes,1,opt,name=usability_period_start_date,json=usabilityPeriodStartDate,proto3" json:"usability_period_start_date,omitempty"`
	UsabilityPeriodEndDate   time.Time `protobuf:"bytes,2,opt,name=usability_period_end_date,json=usabilityPeriodEndDate,proto3" json:"usability_period_end_date,omitempty"`
}

// ItemInstance - struct ItemInstance
type ItemInstance struct {
	*itemproto.ItemInstanceD
	*ItemInstanceT
	*commonproto.CrUpdUser
	*commonstruct.CrUpdTime
}

// ItemInstanceT - struct ItemInstanceT
type ItemInstanceT struct {
	ManufactureDate time.Time `protobuf:"bytes,1,opt,name=manufacture_date,json=manufactureDate,proto3" json:"manufacture_date,omitempty"`
	BestBeforeDate  time.Time `protobuf:"bytes,2,opt,name=best_before_date,json=bestBeforeDate,proto3" json:"best_before_date,omitempty"`
	LotExpiryDate   time.Time `protobuf:"bytes,3,opt,name=lot_expiry_date,json=lotExpiryDate,proto3" json:"lot_expiry_date,omitempty"`
}
