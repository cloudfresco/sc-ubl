package logisticscontrollers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/cloudfresco/sc-ubl/internal/common"
	"github.com/cloudfresco/sc-ubl/internal/config"
	commonproto "github.com/cloudfresco/sc-ubl/internal/protogen/common/v1"
	logisticsproto "github.com/cloudfresco/sc-ubl/internal/protogen/logistics/v1"
	partyproto "github.com/cloudfresco/sc-ubl/internal/protogen/party/v1"
	logisticsworkflows "github.com/cloudfresco/sc-ubl/internal/workflows/logisticsworkflows"
	"github.com/pborman/uuid"
	"go.uber.org/cadence/client"
	"go.uber.org/zap"
)

// ConsignmentController - Create Consignment Controller
type ConsignmentController struct {
	log                      *zap.Logger
	UserServiceClient        partyproto.UserServiceClient
	ConsignmentServiceClient logisticsproto.ConsignmentServiceClient
	wfHelper                 common.WfHelper
	workflowClient           client.Client
	ServerOpt                *config.ServerOptions
}

// NewConsignmentController - Create Consignment Handler
func NewConsignmentController(log *zap.Logger, userServiceClient partyproto.UserServiceClient, consignmentServiceClient logisticsproto.ConsignmentServiceClient, wfHelper common.WfHelper, workflowClient client.Client, serverOpt *config.ServerOptions) *ConsignmentController {
	return &ConsignmentController{
		log:                      log,
		UserServiceClient:        userServiceClient,
		ConsignmentServiceClient: consignmentServiceClient,
		wfHelper:                 wfHelper,
		workflowClient:           workflowClient,
		ServerOpt:                serverOpt,
	}
}

// Index - list Consignment
func (cc *ConsignmentController) Index(w http.ResponseWriter, r *http.Request) {
	ctx, user, _, err := common.GetContextAuthUser(w, r, []string{"cons:read"}, cc.ServerOpt.Auth0Audience, cc.ServerOpt.Auth0Domain, cc.UserServiceClient)
	if err != nil {
		common.RenderErrorJSON(w, "1001", err.Error(), 401, user.RequestId)
		return
	}

	cursor := r.URL.Query().Get("cursor")
	limit := r.URL.Query().Get("limit")

	consignment, err := cc.ConsignmentServiceClient.GetConsignments(ctx, &logisticsproto.GetConsignmentsRequest{Limit: limit, NextCursor: cursor, UserEmail: user.Email, RequestId: user.RequestId})
	if err != nil {
		cc.log.Error("Error",
			zap.String("user", user.Email),
			zap.String("reqid", user.RequestId),
			zap.Error(err))
		common.RenderErrorJSON(w, "1301", err.Error(), 402, user.RequestId)
		return
	}

	common.RenderJSON(w, consignment)
}

// Show - Show Consignment
func (cc *ConsignmentController) Show(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	ctx, user, _, err := common.GetContextAuthUser(w, r, []string{"cons:read"}, cc.ServerOpt.Auth0Audience, cc.ServerOpt.Auth0Domain, cc.UserServiceClient)
	if err != nil {
		common.RenderErrorJSON(w, "1001", err.Error(), 401, user.RequestId)
		return
	}
	consignment, err := cc.ConsignmentServiceClient.GetConsignment(ctx, &logisticsproto.GetConsignmentRequest{GetRequest: &commonproto.GetRequest{Id: id, UserEmail: user.Email, RequestId: user.RequestId}})
	if err != nil {
		cc.log.Error("Error",
			zap.String("reqid", user.RequestId),
			zap.Error(err))
		common.RenderErrorJSON(w, "1103", err.Error(), 402, user.RequestId)
		return
	}
	common.RenderJSON(w, consignment)
}

// CreateConsignment - Create Consignment
func (cc *ConsignmentController) CreateConsignment(w http.ResponseWriter, r *http.Request) {
	ctx, user, token, err := common.GetContextAuthUser(w, r, []string{"cons:cud"}, cc.ServerOpt.Auth0Audience, cc.ServerOpt.Auth0Domain, cc.UserServiceClient)
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

	form := logisticsproto.CreateConsignmentRequest{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&form)
	if err != nil {
		cc.log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		common.RenderErrorJSON(w, "4002", err.Error(), 402, user.RequestId)
		return
	}
	form.UserId = user.UserId
	form.UserEmail = user.Email
	form.RequestId = user.RequestId

	wHelper := cc.wfHelper
	result := wHelper.StartWorkflow(workflowOptions, logisticsworkflows.CreateConsignmentWorkflow, &form, token, user, cc.log)
	workflowClient := cc.workflowClient
	workflowRun := workflowClient.GetWorkflow(ctx, result.ID, result.RunID)
	var consignment logisticsproto.CreateConsignmentResponse
	err = workflowRun.Get(ctx, &consignment)

	if err != nil {
		cc.log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		common.RenderErrorJSON(w, "4002", err.Error(), 402, user.RequestId)
		return
	}

	common.RenderJSON(w, consignment)
}
