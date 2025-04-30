package common

import (
	"context"

	commonproto "github.com/cloudfresco/sc-ubl/internal/protogen/common/v1"
	commonstruct "github.com/cloudfresco/sc-ubl/internal/servicestructs/common/v1"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

const insertLocationSQL = `insert into locations
	  ( 
    uuid4,
    loc_id,
    description,
    conditions,
    country_subentity,
    country_subentity_code,
    location_type_code,
    information_uri,
    loc_name,
    location_coord_lat,
    location_coord_lon,
    altitude_measure,
    address_id,
    validity_period_start_date,
    validity_period_end_date)
  values (:uuid4,
    :loc_id,
    :description,
    :conditions,
    :country_subentity,
    :country_subentity_code,
    :location_type_code,
    :information_uri,
    :loc_name,
    :location_coord_lat,
    :location_coord_lon,
    :altitude_measure,
    :address_id,
    :validity_period_start_date,
    :validity_period_end_date);`

// CreateLocation - For Creating Location
func CreateLocation(ctx context.Context, in *commonproto.Location, userEmail string, requestID string) (*commonproto.Location, error) {
	var err error
	locationD := commonproto.LocationD{}
	locationD.Uuid4, err = GetUUIDBytes()
	if err != nil {
		log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
		return nil, err
	}

	locationD.LocId = in.LocationD.LocId
	locationD.Description = in.LocationD.Description
	locationD.Conditions = in.LocationD.Conditions
	locationD.CountrySubentity = in.LocationD.CountrySubentity
	locationD.CountrySubentityCode = in.LocationD.CountrySubentityCode
	locationD.LocationTypeCode = in.LocationD.LocationTypeCode
	locationD.InformationURI = in.LocationD.InformationURI
	locationD.LocName = in.LocationD.LocName
	locationD.LocationCoordLat = in.LocationD.LocationCoordLat
	locationD.LocationCoordLon = in.LocationD.LocationCoordLon
	locationD.AltitudeMeasure = in.LocationD.AltitudeMeasure
	locationD.AddressId = in.LocationD.AddressId

	locationT := commonproto.LocationT{}
	locationT.ValidityPeriodStartDate = in.LocationT.ValidityPeriodStartDate
	locationT.ValidityPeriodEndDate = in.LocationT.ValidityPeriodEndDate

	location := commonproto.Location{LocationD: &locationD, LocationT: &locationT}

	return &location, nil
}

// InsertLocation - For Inserting Location
func InsertLocation(ctx context.Context, tx *sqlx.Tx, location *commonproto.Location, userEmail string, requestID string) (*commonproto.Location, error) {
	locationTmp, err := crLocationStruct(ctx, location, userEmail, requestID)
	if err != nil {
		log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
		return nil, err
	}

	res, err := tx.NamedExecContext(ctx, insertLocationSQL, locationTmp)
	if err != nil {
		log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
		return nil, err
	}

	uID, err := res.LastInsertId()
	if err != nil {
		log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
		return nil, err
	}
	location.LocationD.Id = uint32(uID)
	uuid4Str, err := UUIDBytesToStr(location.LocationD.Uuid4)
	if err != nil {
		log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
		return nil, err
	}
	location.LocationD.IdS = uuid4Str
	return location, nil
}

// crLocationStruct - process Location details
func crLocationStruct(ctx context.Context, location *commonproto.Location, userEmail string, requestID string) (*commonstruct.Location, error) {
	locationT := new(commonstruct.LocationT)
	locationT.ValidityPeriodStartDate = TimestampToTime(location.LocationT.ValidityPeriodStartDate)
	locationT.ValidityPeriodEndDate = TimestampToTime(location.LocationT.ValidityPeriodEndDate)

	locationTmp := commonstruct.Location{LocationD: location.LocationD, LocationT: locationT}

	return &locationTmp, nil
}
