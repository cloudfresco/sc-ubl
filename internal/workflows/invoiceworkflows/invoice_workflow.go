package invoiceworkflows

import (
	"time"

	invoiceproto "github.com/cloudfresco/sc-ubl/internal/protogen/invoice/v1"
	partyproto "github.com/cloudfresco/sc-ubl/internal/protogen/party/v1"

	"go.uber.org/cadence/workflow"
	"go.uber.org/zap"
)

// CreateInvoiceWorkflow - Create Invoice workflow
func CreateInvoiceWorkflow(ctx workflow.Context, form *invoiceproto.CreateInvoiceRequest, tokenString string, user *partyproto.GetAuthUserDetailsResponse, log *zap.Logger) (*invoiceproto.CreateInvoiceResponse, error) {
	ao := workflow.ActivityOptions{
		ScheduleToStartTimeout: time.Minute,
		StartToCloseTimeout:    time.Minute,
		HeartbeatTimeout:       time.Second * 20,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	logger := workflow.GetLogger(ctx)
	var t *InvoiceActivities
	var invoice invoiceproto.CreateInvoiceResponse
	err := workflow.ExecuteActivity(ctx, t.CreateInvoiceActivity, form, tokenString, user, log).Get(ctx, &invoice)
	if err != nil {
		logger.Error("Failed to CreateInvoiceWorkflow", zap.Error(err))
		log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		return nil, err
	}
	return &invoice, nil
}

// UpdateInvoiceWorkflow - update Invoice workflow
func UpdateInvoiceWorkflow(ctx workflow.Context, form *invoiceproto.UpdateInvoiceRequest, tokenString string, user *partyproto.GetAuthUserDetailsResponse, log *zap.Logger) (string, error) {
	ao := workflow.ActivityOptions{
		ScheduleToStartTimeout: time.Minute,
		StartToCloseTimeout:    time.Minute,
		HeartbeatTimeout:       time.Second * 20,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	logger := workflow.GetLogger(ctx)
	var ia *InvoiceActivities
	var resp string
	err := workflow.ExecuteActivity(ctx, ia.UpdateInvoiceActivity, form, tokenString, user, log).Get(ctx, &resp)
	if err != nil {
		logger.Error("Failed to UpdateInvoiceWorkflow", zap.Error(err))
		log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		return "", err
	}
	return resp, nil
}
