package invoiceworkflows

import (
	"time"

	invoiceproto "github.com/cloudfresco/sc-ubl/internal/protogen/invoice/v1"
	partyproto "github.com/cloudfresco/sc-ubl/internal/protogen/party/v1"

	"go.uber.org/cadence/workflow"
	"go.uber.org/zap"
)

const (
	// ApplicationName is the task list
	ApplicationName = "ubl"
)

// CreateCreditNoteHeaderWorkflow - Create CreditNoteHeader workflow
func CreateCreditNoteHeaderWorkflow(ctx workflow.Context, form *invoiceproto.CreateCreditNoteHeaderRequest, tokenString string, user *partyproto.GetAuthUserDetailsResponse, log *zap.Logger) (*invoiceproto.CreateCreditNoteHeaderResponse, error) {
	ao := workflow.ActivityOptions{
		ScheduleToStartTimeout: time.Minute,
		StartToCloseTimeout:    time.Minute,
		HeartbeatTimeout:       time.Second * 20,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	logger := workflow.GetLogger(ctx)
	var ca *CreditNoteHeaderActivities
	var creditNoteHeader invoiceproto.CreateCreditNoteHeaderResponse
	err := workflow.ExecuteActivity(ctx, ca.CreateCreditNoteHeaderActivity, form, tokenString, user, log).Get(ctx, &creditNoteHeader)
	if err != nil {
		logger.Error("Failed to CreateCreditNoteHeaderWorkflow", zap.Error(err))
		log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		return nil, err
	}
	return &creditNoteHeader, nil
}

// UpdateCreditNoteHeaderWorkflow - update CreditNoteHeader workflow
func UpdateCreditNoteHeaderWorkflow(ctx workflow.Context, form *invoiceproto.UpdateCreditNoteHeaderRequest, tokenString string, user *partyproto.GetAuthUserDetailsResponse, log *zap.Logger) (string, error) {
	ao := workflow.ActivityOptions{
		ScheduleToStartTimeout: time.Minute,
		StartToCloseTimeout:    time.Minute,
		HeartbeatTimeout:       time.Second * 20,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	logger := workflow.GetLogger(ctx)
	var ca *CreditNoteHeaderActivities
	var resp string
	err := workflow.ExecuteActivity(ctx, ca.UpdateCreditNoteHeaderActivity, form, tokenString, user, log).Get(ctx, &resp)
	if err != nil {
		logger.Error("Failed to UpdateCreditNoteHeaderWorkflow", zap.Error(err))
		log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		return "", err
	}
	return resp, nil
}
