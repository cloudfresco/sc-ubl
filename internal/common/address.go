package common

import (
	"context"

	commonproto "github.com/cloudfresco/sc-ubl/internal/protogen/common/v1"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

const insertAddressSQL = `insert into addresses
	  ( 
    uuid4,
    addr_list_agency_id,
    addr_list_id,
    addr_list_version_id,
    address_type_code,
    address_format_code,
    postbox,
    floor1,
    room,
    street_name,
    additional_street_name,
    block_name,
    building_name,
    building_number,
    inhouse_mail,
    department,
    mark_attention,
    mark_care,
    plot_identification,
    city_subdivision_name,
    city_name,
    postal_zone,
    country_subentity,
    country_subentity_code,
    region,
    district,
    timezone_offset,
    country_id_code,
    country_name,
    location_coord_lat,
    location_coord_lon,
    note)
  values (:uuid4,
    :addr_list_agency_id,
    :addr_list_id,
    :addr_list_version_id,
    :address_type_code,
    :address_format_code,
    :postbox,
    :floor1,
    :room,
    :street_name,
    :additional_street_name,
    :block_name,
    :building_name,
    :building_number,
    :inhouse_mail,
    :department,
    :mark_attention,
    :mark_care,
    :plot_identification,
    :city_subdivision_name,
    :city_name,
    :postal_zone,
    :country_subentity,
    :country_subentity_code,
    :region,
    :district,
    :timezone_offset,
    :country_id_code,
    :country_name,
    :location_coord_lat,
    :location_coord_lon,
    :note);`

// CreateAddress - For Creating Address
func CreateAddress(ctx context.Context, in *commonproto.Address, userEmail string, requestID string) (*commonproto.Address, error) {
	var err error
	addr := commonproto.Address{}
	addr.Uuid4, err = GetUUIDBytes()
	if err != nil {
		log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
		return nil, err
	}
	addr.AddrListAgencyId = in.AddrListAgencyId
	addr.AddrListId = in.AddrListAgencyId
	addr.AddrListVersionId = in.AddrListAgencyId
	addr.AddressTypeCode = in.AddressTypeCode
	addr.AddressFormatCode = in.AddressFormatCode
	addr.Postbox = in.Postbox
	addr.Floor1 = in.Floor1
	addr.Room = in.Room
	addr.StreetName = in.StreetName
	addr.AdditionalStreetName = in.AdditionalStreetName
	addr.BlockName = in.BlockName
	addr.BuildingName = in.BuildingName
	addr.BuildingNumber = in.BuildingNumber
	addr.InhouseMail = in.InhouseMail
	addr.Department = in.Department
	addr.MarkAttention = in.MarkAttention
	addr.MarkCare = in.MarkCare
	addr.PlotIdentification = in.PlotIdentification
	addr.CitySubdivisionName = in.CitySubdivisionName
	addr.CityName = in.CityName
	addr.PostalZone = in.PostalZone
	addr.CountrySubentity = in.CountrySubentity
	addr.CountrySubentityCode = in.CountrySubentityCode
	addr.Region = in.Region
	addr.District = in.District
	addr.TimezoneOffset = in.TimezoneOffset
	addr.CountryIdCode = in.CountryIdCode
	addr.CountryName = in.CountryName
	addr.LocationCoordLat = in.LocationCoordLat
	addr.LocationCoordLon = in.LocationCoordLon
	addr.Note = in.Note

	return &addr, nil
}

// InsertAddress - For Inserting Address
func InsertAddress(ctx context.Context, tx *sqlx.Tx, addr *commonproto.Address, userEmail string, requestID string) (*commonproto.Address, error) {
	res, err := tx.NamedExecContext(ctx, insertAddressSQL, addr)
	if err != nil {
		log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
		return nil, err
	}

	uID, err := res.LastInsertId()
	if err != nil {
		log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
		return nil, err
	}
	addr.Id = uint32(uID)
	uuid4Str, err := UUIDBytesToStr(addr.Uuid4)
	if err != nil {
		log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
		return nil, err
	}
	addr.IdS = uuid4Str
	return addr, nil
}
