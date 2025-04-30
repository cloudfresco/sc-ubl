package logisticsworkflows

import (
	"context"

	logisticsproto "github.com/cloudfresco/sc-ubl/internal/protogen/logistics/v1"
	partyproto "github.com/cloudfresco/sc-ubl/internal/protogen/party/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
)

type ConsignmentActivities struct {
	ConsignmentServiceClient logisticsproto.ConsignmentServiceClient
}

// CreateConsignmentActivity - Create Consignment Activity
func (ca *ConsignmentActivities) CreateConsignmentActivity(ctx context.Context, form *logisticsproto.CreateConsignmentRequest, tokenString string, user *partyproto.GetAuthUserDetailsResponse, log *zap.Logger) (*logisticsproto.CreateConsignmentResponse, error) {
	consignmentServiceClient := ca.ConsignmentServiceClient
	md := metadata.Pairs("authorization", "Bearer "+tokenString)
	ctxNew := metadata.NewOutgoingContext(ctx, md)
	consignment, err := consignmentServiceClient.CreateConsignment(ctxNew, form)
	if err != nil {
		log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		return nil, err
	}
	return consignment, nil
}
