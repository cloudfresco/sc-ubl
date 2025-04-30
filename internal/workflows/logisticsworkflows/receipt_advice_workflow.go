package logisticsworkflows

import (
	"time"

	logisticsproto "github.com/cloudfresco/sc-ubl/internal/protogen/logistics/v1"
	partyproto "github.com/cloudfresco/sc-ubl/internal/protogen/party/v1"

	"go.uber.org/cadence/workflow"
	"go.uber.org/zap"
)

const (
	// ApplicationName is the task list
	ApplicationName = "ubl"
)

// CreateReceiptAdviceHeaderWorkflow - Create ReceiptAdviceHeader workflow
func CreateReceiptAdviceHeaderWorkflow(ctx workflow.Context, form *logisticsproto.CreateReceiptAdviceHeaderRequest, tokenString string, user *partyproto.GetAuthUserDetailsResponse, log *zap.Logger) (*logisticsproto.CreateReceiptAdviceHeaderResponse, error) {
	ao := workflow.ActivityOptions{
		ScheduleToStartTimeout: time.Minute,
		StartToCloseTimeout:    time.Minute,
		HeartbeatTimeout:       time.Second * 20,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	logger := workflow.GetLogger(ctx)
	var raha *ReceiptAdviceHeaderActivities
	var receiptAdviceHeader logisticsproto.CreateReceiptAdviceHeaderResponse
	err := workflow.ExecuteActivity(ctx, raha.CreateReceiptAdviceHeaderActivity, form, tokenString, user, log).Get(ctx, &receiptAdviceHeader)
	if err != nil {
		logger.Error("Failed to CreateReceiptAdviceHeaderWorkflow", zap.Error(err))
		log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		return nil, err
	}
	return &receiptAdviceHeader, nil
}

// UpdateReceiptAdviceHeaderWorkflow - update ReceiptAdviceHeader workflow
func UpdateReceiptAdviceHeaderWorkflow(ctx workflow.Context, form *logisticsproto.UpdateReceiptAdviceHeaderRequest, tokenString string, user *partyproto.GetAuthUserDetailsResponse, log *zap.Logger) (string, error) {
	ao := workflow.ActivityOptions{
		ScheduleToStartTimeout: time.Minute,
		StartToCloseTimeout:    time.Minute,
		HeartbeatTimeout:       time.Second * 20,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	logger := workflow.GetLogger(ctx)
	var rah *ReceiptAdviceHeaderActivities
	var resp string
	err := workflow.ExecuteActivity(ctx, rah.UpdateReceiptAdviceHeaderActivity, form, tokenString, user, log).Get(ctx, &resp)
	if err != nil {
		logger.Error("Failed to UpdateReceiptAdviceHeaderWorkflow", zap.Error(err))
		log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		return "", err
	}
	return resp, nil
}
