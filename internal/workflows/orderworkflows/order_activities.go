package orderworkflows

import (
	"context"

	orderproto "github.com/cloudfresco/sc-ubl/internal/protogen/order/v1"
	partyproto "github.com/cloudfresco/sc-ubl/internal/protogen/party/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
)

type PurchaseOrderHeaderActivities struct {
	PurchaseOrderHeaderServiceClient orderproto.PurchaseOrderHeaderServiceClient
}

// CreatePurchaseOrderHeaderActivity - Create PurchaseOrderHeader activity
func (pc *PurchaseOrderHeaderActivities) CreatePurchaseOrderHeaderActivity(ctx context.Context, form *orderproto.CreatePurchaseOrderHeaderRequest, tokenString string, user *partyproto.GetAuthUserDetailsResponse, log *zap.Logger) (*orderproto.CreatePurchaseOrderHeaderResponse, error) {
	purchaseOrderHeaderServiceClient := pc.PurchaseOrderHeaderServiceClient
	md := metadata.Pairs("authorization", "Bearer "+tokenString)
	ctxNew := metadata.NewOutgoingContext(ctx, md)
	purchaseOrderHeader, err := purchaseOrderHeaderServiceClient.CreatePurchaseOrderHeader(ctxNew, form)
	if err != nil {
		log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		return nil, err
	}
	return purchaseOrderHeader, nil
}

// UpdatePurchaseOrderHeaderActivity - update PurchaseOrderHeader activity
func (pc *PurchaseOrderHeaderActivities) UpdatePurchaseOrderHeaderActivity(ctx context.Context, form *orderproto.UpdatePurchaseOrderHeaderRequest, tokenString string, user *partyproto.GetAuthUserDetailsResponse, log *zap.Logger) (string, error) {
	purchaseOrderHeaderServiceClient := pc.PurchaseOrderHeaderServiceClient
	md := metadata.Pairs("authorization", "Bearer "+tokenString)
	ctxNew := metadata.NewOutgoingContext(ctx, md)
	_, err := purchaseOrderHeaderServiceClient.UpdatePurchaseOrderHeader(ctxNew, form)
	if err != nil {
		log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		return "", err
	}
	return "Updated Successfully", nil
}
