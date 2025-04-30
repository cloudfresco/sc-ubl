package itemservices

import (
	"context"
	"net"
	"os"
	"time"

	"github.com/cloudfresco/sc-ubl/internal/common"
	"github.com/cloudfresco/sc-ubl/internal/config"
	commonproto "github.com/cloudfresco/sc-ubl/internal/protogen/common/v1"
	itemproto "github.com/cloudfresco/sc-ubl/internal/protogen/item/v1"
	partyproto "github.com/cloudfresco/sc-ubl/internal/protogen/party/v1"
	partyservice "github.com/cloudfresco/sc-ubl/internal/services/partyservices"
	commonstruct "github.com/cloudfresco/sc-ubl/internal/servicestructs/common/v1"
	itemstruct "github.com/cloudfresco/sc-ubl/internal/servicestructs/item/v1"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// ItemService - For accessing Item services
type ItemService struct {
	log               *zap.Logger
	DBService         *common.DBService
	RedisService      *common.RedisService
	UserServiceClient partyproto.UserServiceClient
	itemproto.UnimplementedItemServiceServer
}

// NewItemService - Create Item service
func NewItemService(log *zap.Logger, dbOpt *common.DBService, redisOpt *common.RedisService, userServiceClient partyproto.UserServiceClient) *ItemService {
	return &ItemService{
		log:               log,
		DBService:         dbOpt,
		RedisService:      redisOpt,
		UserServiceClient: userServiceClient,
	}
}

// StartItemServer - Start Item server
func StartItemServer(log *zap.Logger, isTest bool, pwd string, dbOpt *config.DBOptions, redisOpt *config.RedisOptions, mailerOpt *config.MailerOptions, grpcServerOpt *config.GrpcServerOptions, jwtOpt *config.JWTOptions, oauthOpt *config.OauthOptions, userOpt *config.UserOptions, uptraceOpt *config.UptraceOptions, dbService *common.DBService, redisService *common.RedisService, mailerService common.MailerIntf) {
	common.SetJWTOpt(jwtOpt)

	creds, err := common.GetSrvCred(log, isTest, pwd, grpcServerOpt)
	if err != nil {
		os.Exit(1)
	}

	userCreds, err := common.GetClientCred(log, isTest, pwd, grpcServerOpt)
	if err != nil {
		os.Exit(1)
	}

	var srvOpts []grpc.ServerOption

	userConn, err := grpc.NewClient(grpcServerOpt.GrpcUserServerPort, grpc.WithTransportCredentials(userCreds), grpc.WithStatsHandler(otelgrpc.NewClientHandler()))
	if err != nil {
		log.Error("Error", zap.Error(err))
		os.Exit(1)
	}

	srvOpts = append(srvOpts, grpc.Creds(creds))

	srvOpts = append(srvOpts, grpc.StatsHandler(otelgrpc.NewServerHandler()))

	uc := partyproto.NewUserServiceClient(userConn)
	itemService := NewItemService(log, dbService, redisService, uc)

	lis, err := net.Listen("tcp", grpcServerOpt.GrpcItemServerPort)
	if err != nil {
		log.Error("Error", zap.Error(err))
		os.Exit(1)
	}

	srv := grpc.NewServer(srvOpts...)
	itemproto.RegisterItemServiceServer(srv, itemService)

	if err := srv.Serve(lis); err != nil {
		log.Error("Error", zap.Error(err))
		os.Exit(1)
	}
}

const insertItemSQL = `insert into items
	  (
uuid4,
description,
pack_quantity,
pack_size_numeric,
catalogue_indicator,
item_name,
hazardous_risk_indicator,
additional_information,
keyword,
brand_name,
model_name,
buyers_item_identification_id,
sellers_item_identification_id,
manufacturers_item_identification_id,
standard_item_identification_id,
catalogue_item_identification_id,
additional_item_identification_id,
origin_country_id_code,
origin_country_name,
manufacturer_party_id,
information_content_provider_party_id,
tax_category_id,
status_code,
created_by_user_id,
updated_by_user_id,
created_at,
updated_at)
  values (:uuid4,
:description,
:pack_quantity,
:pack_size_numeric,
:catalogue_indicator,
:item_name,
:hazardous_risk_indicator,
:additional_information,
:keyword,
:brand_name,
:model_name,
:buyers_item_identification_id,
:sellers_item_identification_id,
:manufacturers_item_identification_id,
:standard_item_identification_id,
:catalogue_item_identification_id,
:additional_item_identification_id,
:origin_country_id_code,
:origin_country_name,
:manufacturer_party_id,
:information_content_provider_party_id,
:tax_category_id,
:status_code,
:created_by_user_id,
:updated_by_user_id,
:created_at,
:updated_at);`

const insertItemPropertySQL = `insert into item_properties
	  (
uuid4,
item_property_id,
item_property_name,
item_property_name_code,
test_method,
value,
value_quantity,
value_qualifier,
importance_code,
list_value,
usability_period_start_date,
usability_period_end_date,
item_property_range_measure,
item_property_range_min_value,
item_property_range_max_value,
item_id,
status_code,
created_by_user_id,
updated_by_user_id,
created_at,
updated_at)
  values (:uuid4,
:item_property_id,
:item_property_name,
:item_property_name_code,
:test_method,
:value,
:value_quantity,
:value_qualifier,
:importance_code,
:list_value,
:usability_period_start_date,
:usability_period_end_date,
:item_property_range_measure,
:item_property_range_min_value,
:item_property_range_max_value,
:item_id,
:status_code,
:created_by_user_id,
:updated_by_user_id,
:created_at,
:updated_at);`

const insertItemPropertyGroupSQL = `insert into item_property_groups
	  (
uuid4,
item_property_group_id,
item_property_group_name,
item_property_group_importance_code,
item_property_id)
  values (:uuid4,
:item_property_group_id,
:item_property_group_name,
:item_property_group_importance_code,
:item_property_id);`

const insertItemDimensionSQL = `insert into item_dimensions
	  (
attribute_id,
measure,
description,
minimum_measure,
maximum_measure,
item_id)
  values (:attribute_id,
:measure,
:description,
:minimum_measure,
:maximum_measure,
:item_id);`

const insertItemCommodityClassificationSQL = `insert into item_commodity_classifications
	  (
nature_code,
cargo_type_code,
commodity_code,
item_classification_code,
item_id)
  values (:nature_code,
:cargo_type_code,
:commodity_code,
:item_classification_code,
:item_id);`

const insertItemCertificateSQL = `insert into item_certificates
	  (
cert_id,
certificate_type_code,
certificate_type,
remarks,
party_id,
item_id)
  values (:cert_id,
:certificate_type_code,
:certificate_type,
:remarks,
:party_id,
:item_id);`

const insertItemInstanceSQL = `insert into item_instances
	  (
uuid4,
product_trace_id,
manufacture_date,
best_before_date,
registration_id,
serial_id,
lot_number_id,
lot_expiry_date,
item_id,
status_code,
created_by_user_id,
updated_by_user_id,
created_at,
updated_at)
  values (:uuid4,
:product_trace_id,
:manufacture_date,
:best_before_date,
:registration_id,
:serial_id,
:lot_number_id,
:lot_expiry_date,
:item_id,
:status_code,
:created_by_user_id,
:updated_by_user_id,
:created_at,
:updated_at);`

// CreateItem - Create Item
func (is *ItemService) CreateItem(ctx context.Context, in *itemproto.CreateItemRequest) (*itemproto.CreateItemResponse, error) {
	user, err := partyservice.GetUserWithNewContext(ctx, in.UserId, in.UserEmail, in.RequestId, is.UserServiceClient)
	if err != nil {
		is.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	ttime := common.GetTimeDetails()
	tn := common.TimeToTimestamp(ttime)

	itemD := itemproto.ItemD{}
	itemD.Uuid4, err = common.GetUUIDBytes()
	if err != nil {
		is.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	itemD.Description = in.Description
	itemD.PackQuantity = in.PackQuantity
	itemD.PackSizeNumeric = in.PackSizeNumeric
	itemD.CatalogueIndicator = in.CatalogueIndicator
	itemD.ItemName = in.ItemName
	itemD.HazardousRiskIndicator = in.HazardousRiskIndicator
	itemD.AdditionalInformation = in.AdditionalInformation
	itemD.Keyword = in.Keyword
	itemD.BrandName = in.BrandName
	itemD.ModelName = in.ModelName
	itemD.BuyersItemIdentificationId = in.BuyersItemIdentificationId
	itemD.SellersItemIdentificationId = in.SellersItemIdentificationId
	itemD.ManufacturersItemIdentificationId = in.ManufacturersItemIdentificationId
	itemD.StandardItemIdentificationId = in.StandardItemIdentificationId
	itemD.CatalogueItemIdentificationId = in.CatalogueItemIdentificationId
	itemD.AdditionalItemIdentificationId = in.AdditionalItemIdentificationId
	itemD.OriginCountryIdCode = in.OriginCountryIdCode
	itemD.OriginCountryName = in.OriginCountryName
	itemD.ManufacturerPartyId = in.ManufacturerPartyId
	itemD.InformationContentProviderPartyId = in.InformationContentProviderPartyId
	itemD.TaxCategoryId = in.TaxCategoryId

	crUpdUser := commonproto.CrUpdUser{}
	crUpdUser.StatusCode = "active"
	crUpdUser.CreatedByUserId = user.Id
	crUpdUser.UpdatedByUserId = user.Id

	crUpdTime := commonproto.CrUpdTime{}
	crUpdTime.CreatedAt = tn
	crUpdTime.UpdatedAt = tn

	item := itemproto.Item{ItemD: &itemD, CrUpdUser: &crUpdUser, CrUpdTime: &crUpdTime}

	err = is.insertItem(ctx, insertItemSQL, &item, in.GetUserEmail(), in.GetRequestId())
	if err != nil {
		is.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))

		return nil, err
	}

	itemResponse := itemproto.CreateItemResponse{}
	itemResponse.Item = &item
	return &itemResponse, nil
}

// insertItem - Insert item details into database
func (is *ItemService) insertItem(ctx context.Context, insertItemSQL string, item *itemproto.Item, userEmail string, requestID string) error {
	itemTmp, err := is.crItemStruct(ctx, item, userEmail, requestID)
	if err != nil {
		is.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
		return err
	}
	err = is.DBService.InsUpd(ctx, userEmail, requestID, func(tx *sqlx.Tx) error {
		res, err := tx.NamedExecContext(ctx, insertItemSQL, itemTmp)
		if err != nil {
			is.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}

		uID, err := res.LastInsertId()
		if err != nil {
			is.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}
		item.ItemD.Id = uint32(uID)
		uuid4Str, err := common.UUIDBytesToStr(item.ItemD.Uuid4)
		if err != nil {
			is.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}
		item.ItemD.IdS = uuid4Str
		return nil
	})

	if err != nil {
		is.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
		return err
	}
	return nil
}

// crItemStruct - process Item details
func (is *ItemService) crItemStruct(ctx context.Context, item *itemproto.Item, userEmail string, requestID string) (*itemstruct.Item, error) {
	crUpdTime := new(commonstruct.CrUpdTime)
	crUpdTime.CreatedAt = common.TimestampToTime(item.CrUpdTime.CreatedAt)
	crUpdTime.UpdatedAt = common.TimestampToTime(item.CrUpdTime.UpdatedAt)

	itemTmp := itemstruct.Item{ItemD: item.ItemD, CrUpdUser: item.CrUpdUser, CrUpdTime: crUpdTime}

	return &itemTmp, nil
}

// CreateItemProperty - Create ItemProperty
func (is *ItemService) CreateItemProperty(ctx context.Context, in *itemproto.CreateItemPropertyRequest) (*itemproto.CreateItemPropertyResponse, error) {
	user, err := partyservice.GetUserWithNewContext(ctx, in.UserId, in.UserEmail, in.RequestId, is.UserServiceClient)
	if err != nil {
		is.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	ttime := common.GetTimeDetails()
	tn := common.TimeToTimestamp(ttime)

	usabilityStartDate, err := time.Parse(common.Layout, in.UsabilityPeriodStartDate)
	if err != nil {
		is.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	usabilityEndDate, err := time.Parse(common.Layout, in.UsabilityPeriodEndDate)
	if err != nil {
		is.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	itemPropertyD := itemproto.ItemPropertyD{}
	itemPropertyD.Uuid4, err = common.GetUUIDBytes()
	if err != nil {
		is.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	itemPropertyD.ItemPropertyId = in.ItemPropertyId
	itemPropertyD.ItemPropertyName = in.ItemPropertyName
	itemPropertyD.ItemPropertyNameCode = in.ItemPropertyNameCode
	itemPropertyD.TestMethod = in.TestMethod
	itemPropertyD.Value = in.Value
	itemPropertyD.ValueQuantity = in.ValueQuantity
	itemPropertyD.ValueQualifier = in.ValueQualifier
	itemPropertyD.ImportanceCode = in.ImportanceCode
	itemPropertyD.ListValue = in.ListValue
	itemPropertyD.ItemPropertyRangeMeasure = in.ItemPropertyRangeMeasure
	itemPropertyD.ItemPropertyRangeMinValue = in.ItemPropertyRangeMinValue
	itemPropertyD.ItemPropertyRangeMaxValue = in.ItemPropertyRangeMaxValue
	itemPropertyD.ItemId = in.ItemId

	itemPropertyT := itemproto.ItemPropertyT{}
	itemPropertyT.UsabilityPeriodStartDate = common.TimeToTimestamp(usabilityEndDate.UTC().Truncate(time.Second))
	itemPropertyT.UsabilityPeriodEndDate = common.TimeToTimestamp(usabilityStartDate.UTC().Truncate(time.Second))

	crUpdUser := commonproto.CrUpdUser{}
	crUpdUser.StatusCode = "active"
	crUpdUser.CreatedByUserId = user.Id
	crUpdUser.UpdatedByUserId = user.Id

	crUpdTime := commonproto.CrUpdTime{}
	crUpdTime.CreatedAt = tn
	crUpdTime.UpdatedAt = tn

	itemProperty := itemproto.ItemProperty{ItemPropertyD: &itemPropertyD, ItemPropertyT: &itemPropertyT, CrUpdUser: &crUpdUser, CrUpdTime: &crUpdTime}

	err = is.insertItemProperty(ctx, insertItemPropertySQL, &itemProperty, in.GetUserEmail(), in.GetRequestId())
	if err != nil {
		is.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))

		return nil, err
	}

	itemPropertyResponse := itemproto.CreateItemPropertyResponse{}
	itemPropertyResponse.ItemProperty = &itemProperty
	return &itemPropertyResponse, nil
}

// insertItemProperty - Insert item property details into database
func (is *ItemService) insertItemProperty(ctx context.Context, insertItemPropertySQL string, itemProperty *itemproto.ItemProperty, userEmail string, requestID string) error {
	itemPropertyTmp, err := is.crItemPropertyStruct(ctx, itemProperty, userEmail, requestID)
	if err != nil {
		is.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
		return err
	}

	err = is.DBService.InsUpd(ctx, userEmail, requestID, func(tx *sqlx.Tx) error {
		res, err := tx.NamedExecContext(ctx, insertItemPropertySQL, itemPropertyTmp)
		if err != nil {
			is.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}

		uID, err := res.LastInsertId()
		if err != nil {
			is.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}
		itemProperty.ItemPropertyD.Id = uint32(uID)
		uuid4Str, err := common.UUIDBytesToStr(itemProperty.ItemPropertyD.Uuid4)
		if err != nil {
			is.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}
		itemProperty.ItemPropertyD.IdS = uuid4Str
		return nil
	})

	if err != nil {
		is.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
		return err
	}
	return nil
}

// crItemPropertyStruct - process ItemProperty details
func (is *ItemService) crItemPropertyStruct(ctx context.Context, itemProperty *itemproto.ItemProperty, userEmail string, requestID string) (*itemstruct.ItemProperty, error) {
	itemPropertyT := new(itemstruct.ItemPropertyT)
	itemPropertyT.UsabilityPeriodStartDate = common.TimestampToTime(itemProperty.ItemPropertyT.UsabilityPeriodStartDate)
	itemPropertyT.UsabilityPeriodEndDate = common.TimestampToTime(itemProperty.ItemPropertyT.UsabilityPeriodEndDate)

	crUpdTime := new(commonstruct.CrUpdTime)
	crUpdTime.CreatedAt = common.TimestampToTime(itemProperty.CrUpdTime.CreatedAt)
	crUpdTime.UpdatedAt = common.TimestampToTime(itemProperty.CrUpdTime.UpdatedAt)

	itemPropertyTmp := itemstruct.ItemProperty{ItemPropertyD: itemProperty.ItemPropertyD, ItemPropertyT: itemPropertyT, CrUpdUser: itemProperty.CrUpdUser, CrUpdTime: crUpdTime}

	return &itemPropertyTmp, nil
}

// CreateItemPropertyGroup - Create ItemPropertyGroup
func (is *ItemService) CreateItemPropertyGroup(ctx context.Context, in *itemproto.CreateItemPropertyGroupRequest) (*itemproto.CreateItemPropertyGroupResponse, error) {
	var err error
	itemPropertyGroup := itemproto.ItemPropertyGroup{}
	itemPropertyGroup.Uuid4, err = common.GetUUIDBytes()
	if err != nil {
		is.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	itemPropertyGroup.ItemPropertyGroupId = in.ItemPropertyGroupId
	itemPropertyGroup.ItemPropertyGroupName = in.ItemPropertyGroupName
	itemPropertyGroup.ItemPropertyGroupImportanceCode = in.ItemPropertyGroupImportanceCode
	itemPropertyGroup.ItemPropertyId = in.ItemPropertyId

	err = is.insertItemPropertyGroup(ctx, insertItemPropertyGroupSQL, &itemPropertyGroup, in.GetUserEmail(), in.GetRequestId())
	if err != nil {
		is.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))

		return nil, err
	}

	itemPropertyGroupResponse := itemproto.CreateItemPropertyGroupResponse{}
	itemPropertyGroupResponse.ItemPropertyGroup = &itemPropertyGroup
	return &itemPropertyGroupResponse, nil
}

// insertItemPropertyGroup - Insert item property details into database
func (is *ItemService) insertItemPropertyGroup(ctx context.Context, insertItemPropertyGroupSQL string, itemPropertyGroup *itemproto.ItemPropertyGroup, userEmail string, requestID string) error {
	err := is.DBService.InsUpd(ctx, userEmail, requestID, func(tx *sqlx.Tx) error {
		res, err := tx.NamedExecContext(ctx, insertItemPropertyGroupSQL, itemPropertyGroup)
		if err != nil {
			is.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}

		uID, err := res.LastInsertId()
		if err != nil {
			is.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}
		itemPropertyGroup.Id = uint32(uID)
		uuid4Str, err := common.UUIDBytesToStr(itemPropertyGroup.Uuid4)
		if err != nil {
			is.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}
		itemPropertyGroup.IdS = uuid4Str
		return nil
	})
	if err != nil {
		is.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
		return err
	}
	return nil
}

// CreateItemDimension - Create ItemDimension
func (is *ItemService) CreateItemDimension(ctx context.Context, in *itemproto.CreateItemDimensionRequest) (*itemproto.CreateItemDimensionResponse, error) {
	itemDimension := itemproto.ItemDimension{}
	itemDimension.AttributeId = in.AttributeId
	itemDimension.Measure = in.Measure
	itemDimension.Description = in.Description
	itemDimension.MinimumMeasure = in.MinimumMeasure
	itemDimension.MaximumMeasure = in.MaximumMeasure
	itemDimension.ItemId = in.ItemId

	err := is.insertItemDimension(ctx, insertItemDimensionSQL, &itemDimension, in.GetUserEmail(), in.GetRequestId())
	if err != nil {
		is.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))

		return nil, err
	}

	itemDimensionResponse := itemproto.CreateItemDimensionResponse{}
	itemDimensionResponse.ItemDimension = &itemDimension
	return &itemDimensionResponse, nil
}

// insertItemDimension - Insert item property details into database
func (is *ItemService) insertItemDimension(ctx context.Context, insertItemDimensionSQL string, itemDimension *itemproto.ItemDimension, userEmail string, requestID string) error {
	err := is.DBService.InsUpd(ctx, userEmail, requestID, func(tx *sqlx.Tx) error {
		res, err := tx.NamedExecContext(ctx, insertItemDimensionSQL, itemDimension)
		if err != nil {
			is.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}

		uID, err := res.LastInsertId()
		if err != nil {
			is.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}
		itemDimension.Id = uint32(uID)
		return nil
	})
	if err != nil {
		is.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
		return err
	}
	return nil
}

// CreateItemCommodityClassification - Create ItemCommodityClassification
func (is *ItemService) CreateItemCommodityClassification(ctx context.Context, in *itemproto.CreateItemCommodityClassificationRequest) (*itemproto.CreateItemCommodityClassificationResponse, error) {
	itemCommodityClassification := itemproto.ItemCommodityClassification{}
	itemCommodityClassification.NatureCode = in.NatureCode
	itemCommodityClassification.CargoTypeCode = in.CargoTypeCode
	itemCommodityClassification.CommodityCode = in.CommodityCode
	itemCommodityClassification.ItemClassificationCode = in.ItemClassificationCode
	itemCommodityClassification.ItemId = in.ItemId

	err := is.insertItemCommodityClassification(ctx, insertItemCommodityClassificationSQL, &itemCommodityClassification, in.GetUserEmail(), in.GetRequestId())
	if err != nil {
		is.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))

		return nil, err
	}

	itemCommodityClassificationResponse := itemproto.CreateItemCommodityClassificationResponse{}
	itemCommodityClassificationResponse.ItemCommodityClassification = &itemCommodityClassification
	return &itemCommodityClassificationResponse, nil
}

// insertItemCommodityClassification - Insert item property details into database
func (is *ItemService) insertItemCommodityClassification(ctx context.Context, insertItemCommodityClassificationSQL string, itemCommodityClassification *itemproto.ItemCommodityClassification, userEmail string, requestID string) error {
	err := is.DBService.InsUpd(ctx, userEmail, requestID, func(tx *sqlx.Tx) error {
		res, err := tx.NamedExecContext(ctx, insertItemCommodityClassificationSQL, itemCommodityClassification)
		if err != nil {
			is.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}

		uID, err := res.LastInsertId()
		if err != nil {
			is.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}
		itemCommodityClassification.Id = uint32(uID)
		return nil
	})
	if err != nil {
		is.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
		return err
	}
	return nil
}

// CreateItemCertificate - Create ItemCertificate
func (is *ItemService) CreateItemCertificate(ctx context.Context, in *itemproto.CreateItemCertificateRequest) (*itemproto.CreateItemCertificateResponse, error) {
	itemCertificate := itemproto.ItemCertificate{}
	itemCertificate.CertId = in.CertId
	itemCertificate.CertificateTypeCode = in.CertificateTypeCode
	itemCertificate.CertificateType = in.CertificateType
	itemCertificate.Remarks = in.Remarks
	itemCertificate.PartyId = in.PartyId
	itemCertificate.ItemId = in.ItemId

	err := is.insertItemCertificate(ctx, insertItemCertificateSQL, &itemCertificate, in.GetUserEmail(), in.GetRequestId())
	if err != nil {
		is.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))

		return nil, err
	}

	itemCertificateResponse := itemproto.CreateItemCertificateResponse{}
	itemCertificateResponse.ItemCertificate = &itemCertificate
	return &itemCertificateResponse, nil
}

// insertItemCertificate - Insert item property details into database
func (is *ItemService) insertItemCertificate(ctx context.Context, insertItemCertificateSQL string, itemCertificate *itemproto.ItemCertificate, userEmail string, requestID string) error {
	err := is.DBService.InsUpd(ctx, userEmail, requestID, func(tx *sqlx.Tx) error {
		res, err := tx.NamedExecContext(ctx, insertItemCertificateSQL, itemCertificate)
		if err != nil {
			is.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}

		uID, err := res.LastInsertId()
		if err != nil {
			is.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}
		itemCertificate.Id = uint32(uID)
		return nil
	})
	if err != nil {
		is.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
		return err
	}
	return nil
}

// CreateItemInstance - Create ItemInstance
func (is *ItemService) CreateItemInstance(ctx context.Context, in *itemproto.CreateItemInstanceRequest) (*itemproto.CreateItemInstanceResponse, error) {
	user, err := partyservice.GetUserWithNewContext(ctx, in.UserId, in.UserEmail, in.RequestId, is.UserServiceClient)
	if err != nil {
		is.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	ttime := common.GetTimeDetails()
	tn := common.TimeToTimestamp(ttime)

	manufactureDate, err := time.Parse(common.Layout, in.ManufactureDate)
	if err != nil {
		is.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	bestBeforeDate, err := time.Parse(common.Layout, in.BestBeforeDate)
	if err != nil {
		is.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	lotExpiryDate, err := time.Parse(common.Layout, in.LotExpiryDate)
	if err != nil {
		is.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	itemInstanceD := itemproto.ItemInstanceD{}
	itemInstanceD.Uuid4, err = common.GetUUIDBytes()
	if err != nil {
		is.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	itemInstanceD.ProductTraceId = in.ProductTraceId
	itemInstanceD.RegistrationId = in.RegistrationId
	itemInstanceD.SerialId = in.SerialId
	itemInstanceD.LotNumberId = in.LotNumberId
	itemInstanceD.ItemId = in.ItemId

	itemInstanceT := itemproto.ItemInstanceT{}
	itemInstanceT.ManufactureDate = common.TimeToTimestamp(manufactureDate.UTC().Truncate(time.Second))
	itemInstanceT.BestBeforeDate = common.TimeToTimestamp(bestBeforeDate.UTC().Truncate(time.Second))
	itemInstanceT.LotExpiryDate = common.TimeToTimestamp(lotExpiryDate.UTC().Truncate(time.Second))

	crUpdUser := commonproto.CrUpdUser{}
	crUpdUser.StatusCode = "active"
	crUpdUser.CreatedByUserId = user.Id
	crUpdUser.UpdatedByUserId = user.Id

	crUpdTime := commonproto.CrUpdTime{}
	crUpdTime.CreatedAt = tn
	crUpdTime.UpdatedAt = tn

	itemInstance := itemproto.ItemInstance{ItemInstanceD: &itemInstanceD, ItemInstanceT: &itemInstanceT, CrUpdUser: &crUpdUser, CrUpdTime: &crUpdTime}

	err = is.insertItemInstance(ctx, insertItemInstanceSQL, &itemInstance, in.GetUserEmail(), in.GetRequestId())
	if err != nil {
		is.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))

		return nil, err
	}

	itemInstanceResponse := itemproto.CreateItemInstanceResponse{}
	itemInstanceResponse.ItemInstance = &itemInstance
	return &itemInstanceResponse, nil
}

// insertItemInstance - Insert item property details into database
func (is *ItemService) insertItemInstance(ctx context.Context, insertItemInstanceSQL string, itemInstance *itemproto.ItemInstance, userEmail string, requestID string) error {
	itemInstanceTmp, err := is.crItemInstanceStruct(ctx, itemInstance, userEmail, requestID)
	if err != nil {
		is.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
		return err
	}
	err = is.DBService.InsUpd(ctx, userEmail, requestID, func(tx *sqlx.Tx) error {
		res, err := tx.NamedExecContext(ctx, insertItemInstanceSQL, itemInstanceTmp)
		if err != nil {
			is.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}

		uID, err := res.LastInsertId()
		if err != nil {
			is.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}
		itemInstance.ItemInstanceD.Id = uint32(uID)
		uuid4Str, err := common.UUIDBytesToStr(itemInstance.ItemInstanceD.Uuid4)
		if err != nil {
			is.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}
		itemInstance.ItemInstanceD.IdS = uuid4Str
		return nil
	})
	if err != nil {
		is.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
		return err
	}
	return nil
}

// crItemInstanceStruct - process ItemInstance details
func (is *ItemService) crItemInstanceStruct(ctx context.Context, itemInstance *itemproto.ItemInstance, userEmail string, requestID string) (*itemstruct.ItemInstance, error) {
	itemInstanceT := new(itemstruct.ItemInstanceT)
	itemInstanceT.ManufactureDate = common.TimestampToTime(itemInstance.ItemInstanceT.ManufactureDate)
	itemInstanceT.BestBeforeDate = common.TimestampToTime(itemInstance.ItemInstanceT.BestBeforeDate)
	itemInstanceT.LotExpiryDate = common.TimestampToTime(itemInstance.ItemInstanceT.LotExpiryDate)

	crUpdTime := new(commonstruct.CrUpdTime)
	crUpdTime.CreatedAt = common.TimestampToTime(itemInstance.CrUpdTime.CreatedAt)
	crUpdTime.UpdatedAt = common.TimestampToTime(itemInstance.CrUpdTime.UpdatedAt)

	itemInstanceTmp := itemstruct.ItemInstance{ItemInstanceD: itemInstance.ItemInstanceD, ItemInstanceT: itemInstanceT, CrUpdUser: itemInstance.CrUpdUser, CrUpdTime: crUpdTime}

	return &itemInstanceTmp, nil
}
