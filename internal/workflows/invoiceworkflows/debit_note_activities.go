package invoiceworkflows

import (
	"context"

	invoiceproto "github.com/cloudfresco/sc-ubl/internal/protogen/invoice/v1"
	partyproto "github.com/cloudfresco/sc-ubl/internal/protogen/party/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
)

type DebitNoteHeaderActivities struct {
	DebitNoteHeaderServiceClient invoiceproto.DebitNoteHeaderServiceClient
}

// CreateDebitNoteHeaderActivity - Create DebitNoteHeader activity
func (da *DebitNoteHeaderActivities) CreateDebitNoteHeaderActivity(ctx context.Context, form *invoiceproto.CreateDebitNoteHeaderRequest, tokenString string, user *partyproto.GetAuthUserDetailsResponse, log *zap.Logger) (*invoiceproto.CreateDebitNoteHeaderResponse, error) {
	debitNoteHeaderServiceClient := da.DebitNoteHeaderServiceClient
	md := metadata.Pairs("authorization", "Bearer "+tokenString)
	ctxNew := metadata.NewOutgoingContext(ctx, md)
	debitNoteHeader, err := debitNoteHeaderServiceClient.CreateDebitNoteHeader(ctxNew, form)
	if err != nil {
		log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		return nil, err
	}
	return debitNoteHeader, nil
}

// UpdateDebitNoteHeaderActivity - update DebitNoteHeader activity
func (da *DebitNoteHeaderActivities) UpdateDebitNoteHeaderActivity(ctx context.Context, form *invoiceproto.UpdateDebitNoteHeaderRequest, tokenString string, user *partyproto.GetAuthUserDetailsResponse, log *zap.Logger) (string, error) {
	debitNoteHeaderServiceClient := da.DebitNoteHeaderServiceClient
	md := metadata.Pairs("authorization", "Bearer "+tokenString)
	ctxNew := metadata.NewOutgoingContext(ctx, md)
	_, err := debitNoteHeaderServiceClient.UpdateDebitNoteHeader(ctxNew, form)
	if err != nil {
		log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		return "", err
	}
	return "Updated Successfully", nil
}
