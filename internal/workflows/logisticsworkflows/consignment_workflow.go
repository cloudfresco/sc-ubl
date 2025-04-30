package logisticsworkflows

import (
	"time"

	logisticsproto "github.com/cloudfresco/sc-ubl/internal/protogen/logistics/v1"
	partyproto "github.com/cloudfresco/sc-ubl/internal/protogen/party/v1"

	"go.uber.org/cadence/workflow"
	"go.uber.org/zap"
)

// CreateConsignmentWorkflow - Create Consignment workflow
func CreateConsignmentWorkflow(ctx workflow.Context, form *logisticsproto.CreateConsignmentRequest, tokenString string, user *partyproto.GetAuthUserDetailsResponse, log *zap.Logger) (*logisticsproto.CreateConsignmentResponse, error) {
	ao := workflow.ActivityOptions{
		ScheduleToStartTimeout: time.Minute,
		StartToCloseTimeout:    time.Minute,
		HeartbeatTimeout:       time.Second * 20,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	logger := workflow.GetLogger(ctx)
	var ca *ConsignmentActivities
	var consignment logisticsproto.CreateConsignmentResponse
	err := workflow.ExecuteActivity(ctx, ca.CreateConsignmentActivity, form, tokenString, user, log).Get(ctx, &consignment)
	if err != nil {
		logger.Error("Failed to CreateConsignmentWorkflow", zap.Error(err))
		log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		return nil, err
	}
	return &consignment, nil
}
