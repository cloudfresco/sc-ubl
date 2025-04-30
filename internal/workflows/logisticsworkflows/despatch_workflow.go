package logisticsworkflows

import (
	"time"

	logisticsproto "github.com/cloudfresco/sc-ubl/internal/protogen/logistics/v1"
	partyproto "github.com/cloudfresco/sc-ubl/internal/protogen/party/v1"

	"go.uber.org/cadence/workflow"
	"go.uber.org/zap"
)

// CreateDespatchHeaderWorkflow - Create DespatchHeader workflow
func CreateDespatchHeaderWorkflow(ctx workflow.Context, form *logisticsproto.CreateDespatchHeaderRequest, tokenString string, user *partyproto.GetAuthUserDetailsResponse, log *zap.Logger) (*logisticsproto.CreateDespatchHeaderResponse, error) {
	ao := workflow.ActivityOptions{
		ScheduleToStartTimeout: time.Minute,
		StartToCloseTimeout:    time.Minute,
		HeartbeatTimeout:       time.Second * 20,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	logger := workflow.GetLogger(ctx)
	var dha *DespatchHeaderActivities
	var despatchHeader logisticsproto.CreateDespatchHeaderResponse
	err := workflow.ExecuteActivity(ctx, dha.CreateDespatchHeaderActivity, form, tokenString, user, log).Get(ctx, &despatchHeader)
	if err != nil {
		logger.Error("Failed to CreateDespatchHeaderWorkflow", zap.Error(err))
		log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		return nil, err
	}
	return &despatchHeader, nil
}

// UpdateDespatchHeaderWorkflow - update DespatchHeader workflow
func UpdateDespatchHeaderWorkflow(ctx workflow.Context, form *logisticsproto.UpdateDespatchHeaderRequest, tokenString string, user *partyproto.GetAuthUserDetailsResponse, log *zap.Logger) (string, error) {
	ao := workflow.ActivityOptions{
		ScheduleToStartTimeout: time.Minute,
		StartToCloseTimeout:    time.Minute,
		HeartbeatTimeout:       time.Second * 20,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	logger := workflow.GetLogger(ctx)
	var dh *DespatchHeaderActivities
	var resp string
	err := workflow.ExecuteActivity(ctx, dh.UpdateDespatchHeaderActivity, form, tokenString, user, log).Get(ctx, &resp)
	if err != nil {
		logger.Error("Failed to UpdateDespatchHeaderWorkflow", zap.Error(err))
		log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		return "", err
	}
	return resp, nil
}
