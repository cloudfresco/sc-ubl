package logisticsworkflows

import (
	"context"

	logisticsproto "github.com/cloudfresco/sc-ubl/internal/protogen/logistics/v1"
	partyproto "github.com/cloudfresco/sc-ubl/internal/protogen/party/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
)

type DespatchHeaderActivities struct {
	DespatchHeaderServiceClient logisticsproto.DespatchServiceClient
}

// CreateDespatchHeaderActivity - Create DespatchHeader activity
func (dh *DespatchHeaderActivities) CreateDespatchHeaderActivity(ctx context.Context, form *logisticsproto.CreateDespatchHeaderRequest, tokenString string, user *partyproto.GetAuthUserDetailsResponse, log *zap.Logger) (*logisticsproto.CreateDespatchHeaderResponse, error) {
	despatchHeaderServiceClient := dh.DespatchHeaderServiceClient
	md := metadata.Pairs("authorization", "Bearer "+tokenString)
	ctxNew := metadata.NewOutgoingContext(ctx, md)
	despatchHeader, err := despatchHeaderServiceClient.CreateDespatchHeader(ctxNew, form)
	if err != nil {
		log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		return nil, err
	}
	return despatchHeader, nil
}

// UpdateDespatchHeaderActivity - update DespatchHeader activity
func (dh *DespatchHeaderActivities) UpdateDespatchHeaderActivity(ctx context.Context, form *logisticsproto.UpdateDespatchHeaderRequest, tokenString string, user *partyproto.GetAuthUserDetailsResponse, log *zap.Logger) (string, error) {
	despatchHeaderServiceClient := dh.DespatchHeaderServiceClient
	md := metadata.Pairs("authorization", "Bearer "+tokenString)
	ctxNew := metadata.NewOutgoingContext(ctx, md)
	_, err := despatchHeaderServiceClient.UpdateDespatchHeader(ctxNew, form)
	if err != nil {
		log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		return "", err
	}
	return "Updated Successfully", nil
}
