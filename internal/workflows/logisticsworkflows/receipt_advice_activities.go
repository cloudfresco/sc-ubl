package logisticsworkflows

import (
	"context"

	logisticsproto "github.com/cloudfresco/sc-ubl/internal/protogen/logistics/v1"
	partyproto "github.com/cloudfresco/sc-ubl/internal/protogen/party/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
)

type ReceiptAdviceHeaderActivities struct {
	ReceiptAdviceHeaderServiceClient logisticsproto.ReceiptAdviceHeaderServiceClient
}

// CreateReceiptAdviceHeaderActivity - Create Receipt Advice Header activity
func (rah *ReceiptAdviceHeaderActivities) CreateReceiptAdviceHeaderActivity(ctx context.Context, form *logisticsproto.CreateReceiptAdviceHeaderRequest, tokenString string, user *partyproto.GetAuthUserDetailsResponse, log *zap.Logger) (*logisticsproto.CreateReceiptAdviceHeaderResponse, error) {
	receiptAdviceHeaderServiceClient := rah.ReceiptAdviceHeaderServiceClient
	md := metadata.Pairs("authorization", "Bearer "+tokenString)
	ctxNew := metadata.NewOutgoingContext(ctx, md)
	receiptAdviceHeader, err := receiptAdviceHeaderServiceClient.CreateReceiptAdviceHeader(ctxNew, form)
	if err != nil {
		log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		return nil, err
	}
	return receiptAdviceHeader, nil
}

// UpdateReceiptAdviceHeaderActivity - update ReceiptAdviceHeader activity
func (rah *ReceiptAdviceHeaderActivities) UpdateReceiptAdviceHeaderActivity(ctx context.Context, form *logisticsproto.UpdateReceiptAdviceHeaderRequest, tokenString string, user *partyproto.GetAuthUserDetailsResponse, log *zap.Logger) (string, error) {
	receiptAdviceHeaderServiceClient := rah.ReceiptAdviceHeaderServiceClient
	md := metadata.Pairs("authorization", "Bearer "+tokenString)
	ctxNew := metadata.NewOutgoingContext(ctx, md)
	_, err := receiptAdviceHeaderServiceClient.UpdateReceiptAdviceHeader(ctxNew, form)
	if err != nil {
		log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		return "", err
	}
	return "Updated Successfully", nil
}
