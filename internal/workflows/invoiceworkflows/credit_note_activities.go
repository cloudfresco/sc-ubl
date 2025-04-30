package invoiceworkflows

import (
	"context"

	invoiceproto "github.com/cloudfresco/sc-ubl/internal/protogen/invoice/v1"
	partyproto "github.com/cloudfresco/sc-ubl/internal/protogen/party/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
)

type CreditNoteHeaderActivities struct {
	CreditNoteHeaderServiceClient invoiceproto.CreditNoteHeaderServiceClient
}

// CreateCreditNoteHeaderActivity - Create CreditNoteHeader activity
func (ca *CreditNoteHeaderActivities) CreateCreditNoteHeaderActivity(ctx context.Context, form *invoiceproto.CreateCreditNoteHeaderRequest, tokenString string, user *partyproto.GetAuthUserDetailsResponse, log *zap.Logger) (*invoiceproto.CreateCreditNoteHeaderResponse, error) {
	creditNoteHeaderServiceClient := ca.CreditNoteHeaderServiceClient
	md := metadata.Pairs("authorization", "Bearer "+tokenString)
	ctxNew := metadata.NewOutgoingContext(ctx, md)
	creditNoteHeader, err := creditNoteHeaderServiceClient.CreateCreditNoteHeader(ctxNew, form)
	if err != nil {
		log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		return nil, err
	}
	return creditNoteHeader, nil
}

// UpdateCreditNoteHeaderActivity - update CreditNoteHeader activity
func (ca *CreditNoteHeaderActivities) UpdateCreditNoteHeaderActivity(ctx context.Context, form *invoiceproto.UpdateCreditNoteHeaderRequest, tokenString string, user *partyproto.GetAuthUserDetailsResponse, log *zap.Logger) (string, error) {
	creditNoteHeaderServiceClient := ca.CreditNoteHeaderServiceClient
	md := metadata.Pairs("authorization", "Bearer "+tokenString)
	ctxNew := metadata.NewOutgoingContext(ctx, md)
	_, err := creditNoteHeaderServiceClient.UpdateCreditNoteHeader(ctxNew, form)
	if err != nil {
		log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		return "", err
	}
	return "Updated Successfully", nil
}
