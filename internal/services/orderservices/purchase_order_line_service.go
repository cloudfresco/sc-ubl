package orderservices

import (
	"context"
	"time"

	"github.com/cloudfresco/sc-ubl/internal/common"
	commonproto "github.com/cloudfresco/sc-ubl/internal/protogen/common/v1"
	orderproto "github.com/cloudfresco/sc-ubl/internal/protogen/order/v1"
	partyproto "github.com/cloudfresco/sc-ubl/internal/protogen/party/v1"
	commonstruct "github.com/cloudfresco/sc-ubl/internal/servicestructs/common/v1"
	orderstruct "github.com/cloudfresco/sc-ubl/internal/servicestructs/order/v1"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

const insertPurchaseOrderLineSQL = `insert into purchase_order_lines
	  (uuid4,
pol_id,
substitution_status_code,
note,
sales_order_id,
line_status_code,
quantity,
line_extension_amount,
total_tax_amount,
minimum_quantity,
maximum_quantity,
minimum_backorder_quantity,
maximum_backorder_quantity,
inspection_method_code,
partial_delivery_indicator,
back_order_allowed_indicator,
accounting_cost_code,
accounting_cost,
warranty_information,
originator_party_id,
item_id,
price_amount,
price_base_quantity,
price_change_reason,
price_type_code,
price_type,
orderable_unit_factor_rate,
price_list_id,
purchase_order_header_id,
price_validity_period_start_date,
price_validity_period_end_date,
status_code,
created_by_user_id,
updated_by_user_id,
created_at,
updated_at)
  values (:uuid4,
:pol_id,
:substitution_status_code,
:note,
:sales_order_id,
:line_status_code,
:quantity,
:line_extension_amount,
:total_tax_amount,
:minimum_quantity,
:maximum_quantity,
:minimum_backorder_quantity,
:maximum_backorder_quantity,
:inspection_method_code,
:partial_delivery_indicator,
:back_order_allowed_indicator,
:accounting_cost_code,
:accounting_cost,
:warranty_information,
:originator_party_id,
:item_id,
:price_amount,
:price_base_quantity,
:price_change_reason,
:price_type_code,
:price_type,
:orderable_unit_factor_rate,
:price_list_id,
:purchase_order_header_id,
:price_validity_period_start_date,
:price_validity_period_end_date,
:status_code,
:created_by_user_id,
:updated_by_user_id,
:created_at,
:updated_at);`

const selectPurchaseOrderLinesSQL = `select 
  id,
uuid4,
pol_id,
substitution_status_code,
note,
sales_order_id,
line_status_code,
quantity,
line_extension_amount,
total_tax_amount,
minimum_quantity,
maximum_quantity,
minimum_backorder_quantity,
maximum_backorder_quantity,
inspection_method_code,
partial_delivery_indicator,
back_order_allowed_indicator,
accounting_cost_code,
accounting_cost,
warranty_information,
originator_party_id,
item_id,
price_amount,
price_base_quantity,
price_change_reason,
price_type_code,
price_type,
orderable_unit_factor_rate,
price_list_id,
purchase_order_header_id,
price_validity_period_start_date,
price_validity_period_end_date,
status_code,
created_by_user_id,
updated_by_user_id,
created_at,
updated_at from purchase_order_lines`

// CreatePurchaseOrderLine - Create PurchaseOrderLine
func (ps *PurchaseOrderHeaderService) CreatePurchaseOrderLine(ctx context.Context, in *orderproto.CreatePurchaseOrderLineRequest) (*orderproto.CreatePurchaseOrderLineResponse, error) {
	purchaseOrderLine, err := ps.ProcessPurchaseOrderLineRequest(ctx, in)
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	err = ps.insertPurchaseOrderLine(ctx, insertPurchaseOrderLineSQL, purchaseOrderLine, in.GetUserEmail(), in.GetRequestId())
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	purchaseOrderLineResponse := orderproto.CreatePurchaseOrderLineResponse{}
	purchaseOrderLineResponse.PurchaseOrderLine = purchaseOrderLine
	return &purchaseOrderLineResponse, nil
}

// ProcessPurchaseOrderLineRequest - ProcessPurchaseOrderLineRequest
func (ps *PurchaseOrderHeaderService) ProcessPurchaseOrderLineRequest(ctx context.Context, in *orderproto.CreatePurchaseOrderLineRequest) (*orderproto.PurchaseOrderLine, error) {
	getRequest := commonproto.GetRequest{}
	getRequest.Id = in.UserId
	getRequest.UserEmail = in.UserEmail
	getRequest.RequestId = in.RequestId
	ctxNew, err := common.CreateCtxJWT(ctx)
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	form := partyproto.GetUserRequest{}
	form.GetRequest = &getRequest
	userResponse, err := ps.UserServiceClient.GetUser(ctxNew, &form)
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	user := userResponse.User

	priceValidityPeriodStartDate, err := time.Parse(common.Layout, in.PriceValidityPeriodStartDate)
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	priceValidityPeriodEndDate, err := time.Parse(common.Layout, in.PriceValidityPeriodEndDate)
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	ttime := common.GetTimeDetails()
	tn := common.TimeToTimestamp(ttime)

	purchaseOrderLineD := orderproto.PurchaseOrderLineD{}
	purchaseOrderLineD.Uuid4, err = common.GetUUIDBytes()
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	purchaseOrderLineD.PolId = in.PolId
	purchaseOrderLineD.SubstitutionStatusCode = in.SubstitutionStatusCode
	purchaseOrderLineD.Note = in.Note
	purchaseOrderLineD.SalesOrderId = in.SalesOrderId
	purchaseOrderLineD.LineStatusCode = in.LineStatusCode
	purchaseOrderLineD.Quantity = in.Quantity
	purchaseOrderLineD.LineExtensionAmount = in.LineExtensionAmount
	purchaseOrderLineD.TotalTaxAmount = in.TotalTaxAmount
	purchaseOrderLineD.MinimumQuantity = in.MinimumQuantity
	purchaseOrderLineD.MaximumQuantity = in.MaximumQuantity
	purchaseOrderLineD.MinimumBackorderQuantity = in.MinimumBackorderQuantity
	purchaseOrderLineD.MaximumBackorderQuantity = in.MaximumBackorderQuantity
	purchaseOrderLineD.InspectionMethodCode = in.InspectionMethodCode
	purchaseOrderLineD.PartialDeliveryIndicator = in.PartialDeliveryIndicator
	purchaseOrderLineD.BackOrderAllowedIndicator = in.BackOrderAllowedIndicator
	purchaseOrderLineD.AccountingCostCode = in.AccountingCostCode
	purchaseOrderLineD.AccountingCost = in.AccountingCost
	purchaseOrderLineD.WarrantyInformation = in.WarrantyInformation
	purchaseOrderLineD.OriginatorPartyId = in.OriginatorPartyId
	purchaseOrderLineD.ItemId = in.ItemId
	purchaseOrderLineD.PriceAmount = in.PriceAmount
	purchaseOrderLineD.PriceBaseQuantity = in.PriceBaseQuantity
	purchaseOrderLineD.PriceChangeReason = in.PriceChangeReason
	purchaseOrderLineD.PriceTypeCode = in.PriceTypeCode
	purchaseOrderLineD.PriceType = in.PriceType
	purchaseOrderLineD.OrderableUnitFactorRate = in.OrderableUnitFactorRate
	purchaseOrderLineD.PriceListId = in.PriceListId
	purchaseOrderLineD.PurchaseOrderHeaderId = in.PurchaseOrderHeaderId

	purchaseOrderLineT := orderproto.PurchaseOrderLineT{}
	purchaseOrderLineT.PriceValidityPeriodStartDate = common.TimeToTimestamp(priceValidityPeriodStartDate.UTC().Truncate(time.Second))
	purchaseOrderLineT.PriceValidityPeriodEndDate = common.TimeToTimestamp(priceValidityPeriodEndDate.UTC().Truncate(time.Second))

	crUpdUser := commonproto.CrUpdUser{}
	crUpdUser.StatusCode = "active"
	crUpdUser.CreatedByUserId = user.Id
	crUpdUser.UpdatedByUserId = user.Id

	crUpdTime := commonproto.CrUpdTime{}
	crUpdTime.CreatedAt = tn
	crUpdTime.UpdatedAt = tn

	purchaseOrderLine := orderproto.PurchaseOrderLine{PurchaseOrderLineD: &purchaseOrderLineD, PurchaseOrderLineT: &purchaseOrderLineT, CrUpdUser: &crUpdUser, CrUpdTime: &crUpdTime}

	return &purchaseOrderLine, nil
}

// insertPurchaseOrderLine - Insert PurchaseOrderLine details into database
func (ps *PurchaseOrderHeaderService) insertPurchaseOrderLine(ctx context.Context, insertPurchaseOrderLineSQL string, purchaseOrderLine *orderproto.PurchaseOrderLine, userEmail string, requestID string) error {
	purchaseOrderLineTmp, err := ps.CrPurchaseOrderLineStruct(ctx, purchaseOrderLine, userEmail, requestID)
	if err != nil {
		ps.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
		return err
	}
	err = ps.DBService.InsUpd(ctx, userEmail, requestID, func(tx *sqlx.Tx) error {
		res, err := tx.NamedExecContext(ctx, insertPurchaseOrderLineSQL, purchaseOrderLineTmp)
		if err != nil {
			ps.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}

		uID, err := res.LastInsertId()
		if err != nil {
			ps.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}
		purchaseOrderLine.PurchaseOrderLineD.Id = uint32(uID)
		uuid4Str, err := common.UUIDBytesToStr(purchaseOrderLine.PurchaseOrderLineD.Uuid4)
		if err != nil {
			ps.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}
		purchaseOrderLine.PurchaseOrderLineD.IdS = uuid4Str
		return nil
	})

	if err != nil {
		ps.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
		return err
	}
	return nil
}

// CrPurchaseOrderLineStruct - process PurchaseOrderLine details
func (ps *PurchaseOrderHeaderService) CrPurchaseOrderLineStruct(ctx context.Context, purchaseOrderLine *orderproto.PurchaseOrderLine, userEmail string, requestID string) (*orderstruct.PurchaseOrderLine, error) {
	purchaseOrderLineT := new(orderstruct.PurchaseOrderLineT)
	purchaseOrderLineT.PriceValidityPeriodStartDate = common.TimestampToTime(purchaseOrderLine.PurchaseOrderLineT.PriceValidityPeriodStartDate)
	purchaseOrderLineT.PriceValidityPeriodEndDate = common.TimestampToTime(purchaseOrderLine.PurchaseOrderLineT.PriceValidityPeriodEndDate)

	crUpdTime := new(commonstruct.CrUpdTime)
	crUpdTime.CreatedAt = common.TimestampToTime(purchaseOrderLine.CrUpdTime.CreatedAt)
	crUpdTime.UpdatedAt = common.TimestampToTime(purchaseOrderLine.CrUpdTime.UpdatedAt)

	purchaseOrderLineTmp := orderstruct.PurchaseOrderLine{PurchaseOrderLineD: purchaseOrderLine.PurchaseOrderLineD, PurchaseOrderLineT: purchaseOrderLineT, CrUpdUser: purchaseOrderLine.CrUpdUser, CrUpdTime: crUpdTime}

	return &purchaseOrderLineTmp, nil
}

// GetPurchaseOrderLines - GetPurchaseOrderLines
func (ps *PurchaseOrderHeaderService) GetPurchaseOrderLines(ctx context.Context, inReq *orderproto.GetPurchaseOrderLinesRequest) (*orderproto.GetPurchaseOrderLinesResponse, error) {
	in := inReq.GetRequest
	getRequest := commonproto.GetRequest{}
	getRequest.Id = in.Id
	getRequest.UserEmail = in.UserEmail
	getRequest.RequestId = in.RequestId
	form := orderproto.GetPurchaseOrderHeaderRequest{}
	form.GetRequest = &getRequest

	purchaseOrderHeaderResponse, err := ps.GetPurchaseOrderHeader(ctx, &form)
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	purchaseOrderHeader := purchaseOrderHeaderResponse.PurchaseOrderHeader

	purchaseOrderLines := []*orderproto.PurchaseOrderLine{}

	nselectPurchaseOrderLinesSQL := selectPurchaseOrderLinesSQL + ` where purchase_order_header_id = ? and status_code = ?;`
	rows, err := ps.DBService.DB.QueryxContext(ctx, nselectPurchaseOrderLinesSQL, purchaseOrderHeader.PurchaseOrderHeaderD.Id, "active")
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	for rows.Next() {

		purchaseOrderLineTmp := orderstruct.PurchaseOrderLine{}
		err = rows.StructScan(&purchaseOrderLineTmp)
		if err != nil {
			ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
			return nil, err
		}
		purchaseOrderLine, err := ps.getPurchaseOrderLineStruct(ctx, &getRequest, purchaseOrderLineTmp)
		if err != nil {
			ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
			return nil, err
		}

		purchaseOrderLines = append(purchaseOrderLines, purchaseOrderLine)
	}

	purchaseOrderLinesResponse := orderproto.GetPurchaseOrderLinesResponse{}
	purchaseOrderLinesResponse.PurchaseOrderLines = purchaseOrderLines
	return &purchaseOrderLinesResponse, nil
}

// getPurchaseOrderLineStruct - Get PurchaseOrderLine
func (ps *PurchaseOrderHeaderService) getPurchaseOrderLineStruct(ctx context.Context, in *commonproto.GetRequest, purchaseOrderLineTmp orderstruct.PurchaseOrderLine) (*orderproto.PurchaseOrderLine, error) {
	purchaseOrderLineT := new(orderproto.PurchaseOrderLineT)
	purchaseOrderLineT.PriceValidityPeriodStartDate = common.TimeToTimestamp(purchaseOrderLineTmp.PurchaseOrderLineT.PriceValidityPeriodStartDate)
	purchaseOrderLineT.PriceValidityPeriodEndDate = common.TimeToTimestamp(purchaseOrderLineTmp.PurchaseOrderLineT.PriceValidityPeriodEndDate)

	crUpdTime := new(commonproto.CrUpdTime)
	crUpdTime.CreatedAt = common.TimeToTimestamp(purchaseOrderLineTmp.CrUpdTime.CreatedAt)
	crUpdTime.UpdatedAt = common.TimeToTimestamp(purchaseOrderLineTmp.CrUpdTime.UpdatedAt)

	uuid4Str, err := common.UUIDBytesToStr(purchaseOrderLineTmp.PurchaseOrderLineD.Uuid4)
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	purchaseOrderLineTmp.PurchaseOrderLineD.IdS = uuid4Str

	purchaseOrderLine := orderproto.PurchaseOrderLine{PurchaseOrderLineD: purchaseOrderLineTmp.PurchaseOrderLineD, PurchaseOrderLineT: purchaseOrderLineT, CrUpdUser: purchaseOrderLineTmp.CrUpdUser, CrUpdTime: crUpdTime}

	return &purchaseOrderLine, nil
}
