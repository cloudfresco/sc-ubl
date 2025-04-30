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

// ReceiptAdviceHeaderController - Create ReceiptAdviceHeader Controller
type ReceiptAdviceHeaderController struct {
	log                              *zap.Logger
	UserServiceClient                partyproto.UserServiceClient
	ReceiptAdviceHeaderServiceClient logisticsproto.ReceiptAdviceHeaderServiceClient
	wfHelper                         common.WfHelper
	workflowClient                   client.Client
	ServerOpt                        *config.ServerOptions
}

// NewReceiptAdviceHeaderController - Create ReceiptAdviceHeader Handler
func NewReceiptAdviceHeaderController(log *zap.Logger, userServiceClient partyproto.UserServiceClient, receiptAdviceHeaderServiceClient logisticsproto.ReceiptAdviceHeaderServiceClient, wfHelper common.WfHelper, workflowClient client.Client, serverOpt *config.ServerOptions) *ReceiptAdviceHeaderController {
	return &ReceiptAdviceHeaderController{
		log:                              log,
		UserServiceClient:                userServiceClient,
		ReceiptAdviceHeaderServiceClient: receiptAdviceHeaderServiceClient,
		wfHelper:                         wfHelper,
		workflowClient:                   workflowClient,
		ServerOpt:                        serverOpt,
	}
}

// CreateReceiptAdviceHeader - Create Receipt Advice Header
func (rc *ReceiptAdviceHeaderController) CreateReceiptAdviceHeader(w http.ResponseWriter, r *http.Request) {
	ctx, user, token, err := common.GetContextAuthUser(w, r, []string{"rcptadv:cud"}, rc.ServerOpt.Auth0Audience, rc.ServerOpt.Auth0Domain, rc.UserServiceClient)
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

	form := logisticsproto.CreateReceiptAdviceHeaderRequest{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&form)
	if err != nil {
		rc.log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		common.RenderErrorJSON(w, "4002", err.Error(), 402, user.RequestId)
		return
	}
	form.UserId = user.UserId
	form.UserEmail = user.Email
	form.RequestId = user.RequestId

	wHelper := rc.wfHelper
	result := wHelper.StartWorkflow(workflowOptions, logisticsworkflows.CreateReceiptAdviceHeaderWorkflow, &form, token, user, rc.log)
	workflowClient := rc.workflowClient
	workflowRun := workflowClient.GetWorkflow(ctx, result.ID, result.RunID)
	var receiptAdviceHeader logisticsproto.CreateReceiptAdviceHeaderResponse
	err = workflowRun.Get(ctx, &receiptAdviceHeader)

	if err != nil {
		rc.log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		common.RenderErrorJSON(w, "4002", err.Error(), 402, user.RequestId)
		return
	}

	common.RenderJSON(w, receiptAdviceHeader)
}

// Index - list Receipt Advice Headers
func (rc *ReceiptAdviceHeaderController) Index(w http.ResponseWriter, r *http.Request) {
	ctx, user, _, err := common.GetContextAuthUser(w, r, []string{"rcptadv:read"}, rc.ServerOpt.Auth0Audience, rc.ServerOpt.Auth0Domain, rc.UserServiceClient)
	if err != nil {
		common.RenderErrorJSON(w, "1001", err.Error(), 401, user.RequestId)
		return
	}

	cursor := r.URL.Query().Get("cursor")
	limit := r.URL.Query().Get("limit")

	receiptAdviceHeader, err := rc.ReceiptAdviceHeaderServiceClient.GetReceiptAdviceHeaders(ctx, &logisticsproto.GetReceiptAdviceHeadersRequest{Limit: limit, NextCursor: cursor, UserEmail: user.Email, RequestId: user.RequestId})
	if err != nil {
		rc.log.Error("Error",
			zap.String("user", user.Email),
			zap.String("reqid", user.RequestId),
			zap.Error(err))
		common.RenderErrorJSON(w, "1301", err.Error(), 402, user.RequestId)
		return
	}

	common.RenderJSON(w, receiptAdviceHeader)
}

// Show - Show ReceiptAdvice
func (rc *ReceiptAdviceHeaderController) Show(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	ctx, user, _, err := common.GetContextAuthUser(w, r, []string{"rcptadv:read"}, rc.ServerOpt.Auth0Audience, rc.ServerOpt.Auth0Domain, rc.UserServiceClient)
	if err != nil {
		common.RenderErrorJSON(w, "1001", err.Error(), 401, user.RequestId)
		return
	}
	receiptAdviceHeader, err := rc.ReceiptAdviceHeaderServiceClient.GetReceiptAdviceHeader(ctx, &logisticsproto.GetReceiptAdviceHeaderRequest{GetRequest: &commonproto.GetRequest{Id: id, UserEmail: user.Email, RequestId: user.RequestId}})
	if err != nil {
		rc.log.Error("Error",
			zap.String("reqid", user.RequestId),
			zap.Error(err))
		common.RenderErrorJSON(w, "1103", err.Error(), 402, user.RequestId)
		return
	}
	common.RenderJSON(w, receiptAdviceHeader)
}

// GetReceiptAdviceLines - Get ReceiptAdvice Lines
func (rc *ReceiptAdviceHeaderController) GetReceiptAdviceLines(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	ctx, user, _, err := common.GetContextAuthUser(w, r, []string{"rcptadv:read"}, rc.ServerOpt.Auth0Audience, rc.ServerOpt.Auth0Domain, rc.UserServiceClient)
	if err != nil {
		common.RenderErrorJSON(w, "1001", err.Error(), 401, user.RequestId)
		return
	}

	receiptAdviceLines, err := rc.ReceiptAdviceHeaderServiceClient.GetReceiptAdviceLines(ctx, &logisticsproto.GetReceiptAdviceLinesRequest{GetRequest: &commonproto.GetRequest{Id: id, UserEmail: user.Email, RequestId: user.RequestId}})
	if err != nil {
		rc.log.Error("Error",
			zap.String("reqid", user.RequestId),
			zap.Error(err))
		common.RenderErrorJSON(w, "1103", err.Error(), 402, user.RequestId)
		return
	}
	common.RenderJSON(w, receiptAdviceLines)
}

// UpdateReceiptAdviceHeader - Update ReceiptAdviceHeader
func (rc *ReceiptAdviceHeaderController) UpdateReceiptAdviceHeader(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	ctx, user, token, err := common.GetContextAuthUser(w, r, []string{"rcptadv:cud"}, rc.ServerOpt.Auth0Audience, rc.ServerOpt.Auth0Domain, rc.UserServiceClient)
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

	form := logisticsproto.UpdateReceiptAdviceHeaderRequest{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&form)
	if err != nil {
		rc.log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		common.RenderErrorJSON(w, "4009", err.Error(), 402, user.RequestId)
		return
	}
	form.Id = id
	form.UserId = user.UserId
	form.UserEmail = user.Email
	form.RequestId = user.RequestId

	wHelper := rc.wfHelper
	result := wHelper.StartWorkflow(workflowOptions, logisticsworkflows.UpdateReceiptAdviceHeaderWorkflow, &form, token, user, rc.log)
	workflowClient := rc.workflowClient
	workflowRun := workflowClient.GetWorkflow(ctx, result.ID, result.RunID)
	var response string
	err = workflowRun.Get(ctx, &response)
	if err != nil {
		rc.log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		common.RenderErrorJSON(w, "4009", err.Error(), 402, user.RequestId)
		return
	}
	common.RenderJSON(w, response)
}
