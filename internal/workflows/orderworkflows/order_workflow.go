package orderworkflows

import (
	"time"

	orderproto "github.com/cloudfresco/sc-ubl/internal/protogen/order/v1"
	partyproto "github.com/cloudfresco/sc-ubl/internal/protogen/party/v1"

	"go.uber.org/cadence/workflow"
	"go.uber.org/zap"
)

const (
	// ApplicationName is the task list
	ApplicationName = "ubl"
)

// CreatePurchaseOrderHeaderWorkflow - Create PurchaseOrderHeader workflow
func CreatePurchaseOrderHeaderWorkflow(ctx workflow.Context, form *orderproto.CreatePurchaseOrderHeaderRequest, tokenString string, user *partyproto.GetAuthUserDetailsResponse, log *zap.Logger) (*orderproto.CreatePurchaseOrderHeaderResponse, error) {
	ao := workflow.ActivityOptions{
		ScheduleToStartTimeout: time.Minute,
		StartToCloseTimeout:    time.Minute,
		HeartbeatTimeout:       time.Second * 20,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	logger := workflow.GetLogger(ctx)
	var pc *PurchaseOrderHeaderActivities
	var purchaseOrderHeader orderproto.CreatePurchaseOrderHeaderResponse
	err := workflow.ExecuteActivity(ctx, pc.CreatePurchaseOrderHeaderActivity, form, tokenString, user, log).Get(ctx, &purchaseOrderHeader)
	if err != nil {
		logger.Error("Failed to CreatePurchaseOrderHeaderWorkflow", zap.Error(err))
		log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		return nil, err
	}
	return &purchaseOrderHeader, nil
}

// UpdatePurchaseOrderHeaderWorkflow - update PurchaseOrderHeader workflow
func UpdatePurchaseOrderHeaderWorkflow(ctx workflow.Context, form *orderproto.UpdatePurchaseOrderHeaderRequest, tokenString string, user *partyproto.GetAuthUserDetailsResponse, log *zap.Logger) (string, error) {
	ao := workflow.ActivityOptions{
		ScheduleToStartTimeout: time.Minute,
		StartToCloseTimeout:    time.Minute,
		HeartbeatTimeout:       time.Second * 20,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	logger := workflow.GetLogger(ctx)
	var pc *PurchaseOrderHeaderActivities
	var resp string
	err := workflow.ExecuteActivity(ctx, pc.UpdatePurchaseOrderHeaderActivity, form, tokenString, user, log).Get(ctx, &resp)
	if err != nil {
		logger.Error("Failed to UpdatePurchaseOrderHeaderWorkflow", zap.Error(err))
		log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		return "", err
	}
	return resp, nil
}
