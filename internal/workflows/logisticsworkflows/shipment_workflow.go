package logisticsworkflows

import (
	"time"

	logisticsproto "github.com/cloudfresco/sc-ubl/internal/protogen/logistics/v1"
	partyproto "github.com/cloudfresco/sc-ubl/internal/protogen/party/v1"

	"go.uber.org/cadence/workflow"
	"go.uber.org/zap"
)

// CreateShipmentWorkflow - Create Shipment workflow
func CreateShipmentWorkflow(ctx workflow.Context, form *logisticsproto.CreateShipmentRequest, tokenString string, user *partyproto.GetAuthUserDetailsResponse, log *zap.Logger) (*logisticsproto.CreateShipmentResponse, error) {
	ao := workflow.ActivityOptions{
		ScheduleToStartTimeout: time.Minute,
		StartToCloseTimeout:    time.Minute,
		HeartbeatTimeout:       time.Second * 20,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	logger := workflow.GetLogger(ctx)
	var sa *ShipmentActivities
	var shipment logisticsproto.CreateShipmentResponse
	err := workflow.ExecuteActivity(ctx, sa.CreateShipmentActivity, form, tokenString, user, log).Get(ctx, &shipment)
	if err != nil {
		logger.Error("Failed to CreateShipmentWorkflow", zap.Error(err))
		log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		return nil, err
	}
	return &shipment, nil
}
