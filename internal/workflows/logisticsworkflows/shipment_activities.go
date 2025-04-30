package logisticsworkflows

import (
	"context"

	logisticsproto "github.com/cloudfresco/sc-ubl/internal/protogen/logistics/v1"
	partyproto "github.com/cloudfresco/sc-ubl/internal/protogen/party/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
)

type ShipmentActivities struct {
	ShipmentServiceClient logisticsproto.ShipmentServiceClient
}

// CreateShipmentActivity - Create Shipment activity
func (sa *ShipmentActivities) CreateShipmentActivity(ctx context.Context, form *logisticsproto.CreateShipmentRequest, tokenString string, user *partyproto.GetAuthUserDetailsResponse, log *zap.Logger) (*logisticsproto.CreateShipmentResponse, error) {
	shipmentServiceClient := sa.ShipmentServiceClient
	md := metadata.Pairs("authorization", "Bearer "+tokenString)
	ctxNew := metadata.NewOutgoingContext(ctx, md)
	shipment, err := shipmentServiceClient.CreateShipment(ctxNew, form)
	if err != nil {
		log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		return nil, err
	}
	return shipment, nil
}
