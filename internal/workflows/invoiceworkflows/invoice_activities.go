package invoiceworkflows

import (
	"context"

	invoiceproto "github.com/cloudfresco/sc-ubl/internal/protogen/invoice/v1"
	partyproto "github.com/cloudfresco/sc-ubl/internal/protogen/party/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
)

type InvoiceActivities struct {
	InvoiceServiceClient invoiceproto.InvoiceServiceClient
}

// CreateInvoiceActivity - Create Invoice activity
func (ia *InvoiceActivities) CreateInvoiceActivity(ctx context.Context, form *invoiceproto.CreateInvoiceRequest, tokenString string, user *partyproto.GetAuthUserDetailsResponse, log *zap.Logger) (*invoiceproto.CreateInvoiceResponse, error) {
	invoiceServiceClient := ia.InvoiceServiceClient
	md := metadata.Pairs("authorization", "Bearer "+tokenString)
	ctxNew := metadata.NewOutgoingContext(ctx, md)
	invoice, err := invoiceServiceClient.CreateInvoice(ctxNew, form)
	if err != nil {
		log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		return nil, err
	}
	return invoice, nil
}

// UpdateInvoiceActivity - update Invoice activity
func (ia *InvoiceActivities) UpdateInvoiceActivity(ctx context.Context, form *invoiceproto.UpdateInvoiceRequest, tokenString string, user *partyproto.GetAuthUserDetailsResponse, log *zap.Logger) (string, error) {
	invoiceServiceClient := ia.InvoiceServiceClient
	md := metadata.Pairs("authorization", "Bearer "+tokenString)
	ctxNew := metadata.NewOutgoingContext(ctx, md)
	_, err := invoiceServiceClient.UpdateInvoice(ctxNew, form)
	if err != nil {
		log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		return "", err
	}
	return "Updated Successfully", nil
}
