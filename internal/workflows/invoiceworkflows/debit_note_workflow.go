package invoiceworkflows

import (
	"time"

	invoiceproto "github.com/cloudfresco/sc-ubl/internal/protogen/invoice/v1"
	partyproto "github.com/cloudfresco/sc-ubl/internal/protogen/party/v1"

	"go.uber.org/cadence/workflow"
	"go.uber.org/zap"
)

// CreateDebitNoteHeaderWorkflow - Create DebitNoteHeader workflow
func CreateDebitNoteHeaderWorkflow(ctx workflow.Context, form *invoiceproto.CreateDebitNoteHeaderRequest, tokenString string, user *partyproto.GetAuthUserDetailsResponse, log *zap.Logger) (*invoiceproto.CreateDebitNoteHeaderResponse, error) {
	ao := workflow.ActivityOptions{
		ScheduleToStartTimeout: time.Minute,
		StartToCloseTimeout:    time.Minute,
		HeartbeatTimeout:       time.Second * 20,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	logger := workflow.GetLogger(ctx)
	var da *DebitNoteHeaderActivities
	var debitNoteHeader invoiceproto.CreateDebitNoteHeaderResponse
	err := workflow.ExecuteActivity(ctx, da.CreateDebitNoteHeaderActivity, form, tokenString, user, log).Get(ctx, &debitNoteHeader)
	if err != nil {
		logger.Error("Failed to CreateDebitNoteHeaderWorkflow", zap.Error(err))
		log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		return nil, err
	}
	return &debitNoteHeader, nil
}

// UpdateDebitNoteHeaderWorkflow - update DebitNoteHeader workflow
func UpdateDebitNoteHeaderWorkflow(ctx workflow.Context, form *invoiceproto.UpdateDebitNoteHeaderRequest, tokenString string, user *partyproto.GetAuthUserDetailsResponse, log *zap.Logger) (string, error) {
	ao := workflow.ActivityOptions{
		ScheduleToStartTimeout: time.Minute,
		StartToCloseTimeout:    time.Minute,
		HeartbeatTimeout:       time.Second * 20,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	logger := workflow.GetLogger(ctx)
	var da *DebitNoteHeaderActivities
	var resp string
	err := workflow.ExecuteActivity(ctx, da.UpdateDebitNoteHeaderActivity, form, tokenString, user, log).Get(ctx, &resp)
	if err != nil {
		logger.Error("Failed to UpdateDebitNoteHeaderWorkflow", zap.Error(err))
		log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		return "", err
	}
	return resp, nil
}
