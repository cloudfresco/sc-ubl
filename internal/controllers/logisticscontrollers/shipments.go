package logisticscontrollers

import (
	"encoding/json"
	"time"

	// "encoding/json"
	"net/http"

	"github.com/cloudfresco/sc-ubl/internal/common"
	"github.com/cloudfresco/sc-ubl/internal/config"
	logisticsworkflows "github.com/cloudfresco/sc-ubl/internal/workflows/logisticsworkflows"
	"github.com/pborman/uuid"

	// commonproto "github.com/cloudfresco/sc-ubl/internal/protogen/common/v1"
	logisticsproto "github.com/cloudfresco/sc-ubl/internal/protogen/logistics/v1"
	partyproto "github.com/cloudfresco/sc-ubl/internal/protogen/party/v1"
	"go.uber.org/cadence/client"
	"go.uber.org/zap"
)

// ShipmentController - Create Shipment Controller
type ShipmentController struct {
	log                   *zap.Logger
	UserServiceClient     partyproto.UserServiceClient
	ShipmentServiceClient logisticsproto.ShipmentServiceClient
	wfHelper              common.WfHelper
	workflowClient        client.Client
	ServerOpt             *config.ServerOptions
}

// NewShipmentController - Create Shipment Handler
func NewShipmentController(log *zap.Logger, userServiceClient partyproto.UserServiceClient, shipmentServiceClient logisticsproto.ShipmentServiceClient, wfHelper common.WfHelper, workflowClient client.Client, serverOpt *config.ServerOptions) *ShipmentController {
	return &ShipmentController{
		log:                   log,
		UserServiceClient:     userServiceClient,
		ShipmentServiceClient: shipmentServiceClient,
		wfHelper:              wfHelper,
		workflowClient:        workflowClient,
		ServerOpt:             serverOpt,
	}
}

func (sc *ShipmentController) CreateShipment(w http.ResponseWriter, r *http.Request) {
	ctx, user, token, err := common.GetContextAuthUser(w, r, []string{"shipment:cud"}, sc.ServerOpt.Auth0Audience, sc.ServerOpt.Auth0Domain, sc.UserServiceClient)
	if err != nil {
		common.RenderErrorJSON(w, "1001", err.Error(), 401, user.RequestId)
		return
	}

	workflowOptions := client.StartWorkflowOptions{
		ID:                              "ubl_" + uuid.New(),
		TaskList:                        logisticsworkflows.ApplicationName,
		ExecutionStartToCloseTimeout:    time.Minute,
		DecisionTaskStartToCloseTimeout: time.Minute,
	}

	form := logisticsproto.CreateShipmentRequest{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&form)
	if err != nil {
		sc.log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		common.RenderErrorJSON(w, "4002", err.Error(), 402, user.RequestId)
		return
	}
	form.UserId = user.UserId
	form.UserEmail = user.Email
	form.RequestId = user.RequestId

	wHelper := sc.wfHelper
	result := wHelper.StartWorkflow(workflowOptions, logisticsworkflows.CreateShipmentWorkflow, &form, token, user, sc.log)
	workflowClient := sc.workflowClient
	workflowRun := workflowClient.GetWorkflow(ctx, result.ID, result.RunID)
	var shipment logisticsproto.CreateShipmentResponse
	err = workflowRun.Get(ctx, &shipment)

	if err != nil {
		sc.log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		common.RenderErrorJSON(w, "4002", err.Error(), 402, user.RequestId)
		return
	}

	common.RenderJSON(w, shipment)
}
