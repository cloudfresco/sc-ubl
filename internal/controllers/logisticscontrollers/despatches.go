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

// DespatchHeaderController - Create DespatchHeader Controller
type DespatchHeaderController struct {
	log                   *zap.Logger
	UserServiceClient     partyproto.UserServiceClient
	DespatchServiceClient logisticsproto.DespatchServiceClient
	wfHelper              common.WfHelper
	workflowClient        client.Client
	ServerOpt             *config.ServerOptions
}

// NewDespatchHeaderController - Create DespatchHeader Handler
func NewDespatchHeaderController(log *zap.Logger, userServiceClient partyproto.UserServiceClient, despatchServiceClient logisticsproto.DespatchServiceClient, wfHelper common.WfHelper, workflowClient client.Client, serverOpt *config.ServerOptions) *DespatchHeaderController {
	return &DespatchHeaderController{
		log:                   log,
		UserServiceClient:     userServiceClient,
		DespatchServiceClient: despatchServiceClient,
		wfHelper:              wfHelper,
		workflowClient:        workflowClient,
		ServerOpt:             serverOpt,
	}
}

// CreateDespatchHeader - Create Despatch Header
func (dc *DespatchHeaderController) CreateDespatchHeader(w http.ResponseWriter, r *http.Request) {
	ctx, user, token, err := common.GetContextAuthUser(w, r, []string{"despatch:cud"}, dc.ServerOpt.Auth0Audience, dc.ServerOpt.Auth0Domain, dc.UserServiceClient)
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

	form := logisticsproto.CreateDespatchHeaderRequest{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&form)
	if err != nil {
		dc.log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		common.RenderErrorJSON(w, "4002", err.Error(), 402, user.RequestId)
		return
	}
	form.UserId = user.UserId
	form.UserEmail = user.Email
	form.RequestId = user.RequestId

	wHelper := dc.wfHelper
	result := wHelper.StartWorkflow(workflowOptions, logisticsworkflows.CreateDespatchHeaderWorkflow, &form, token, user, dc.log)
	workflowClient := dc.workflowClient
	workflowRun := workflowClient.GetWorkflow(ctx, result.ID, result.RunID)
	var despatchHeader logisticsproto.CreateDespatchHeaderResponse
	err = workflowRun.Get(ctx, &despatchHeader)

	if err != nil {
		dc.log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		common.RenderErrorJSON(w, "4002", err.Error(), 402, user.RequestId)
		return
	}

	common.RenderJSON(w, despatchHeader)
}

// Index - list Despatch Headers
func (dc *DespatchHeaderController) Index(w http.ResponseWriter, r *http.Request) {
	ctx, user, _, err := common.GetContextAuthUser(w, r, []string{"despatch:read"}, dc.ServerOpt.Auth0Audience, dc.ServerOpt.Auth0Domain, dc.UserServiceClient)
	if err != nil {
		common.RenderErrorJSON(w, "1001", err.Error(), 401, user.RequestId)
		return
	}

	cursor := r.URL.Query().Get("cursor")
	limit := r.URL.Query().Get("limit")

	despatchHeader, err := dc.DespatchServiceClient.GetDespatchHeaders(ctx, &logisticsproto.GetDespatchHeadersRequest{Limit: limit, NextCursor: cursor, UserEmail: user.Email, RequestId: user.RequestId})
	if err != nil {
		dc.log.Error("Error",
			zap.String("user", user.Email),
			zap.String("reqid", user.RequestId),
			zap.Error(err))
		common.RenderErrorJSON(w, "1301", err.Error(), 402, user.RequestId)
		return
	}

	common.RenderJSON(w, despatchHeader)
}

// Show - Show Despatch
func (dc *DespatchHeaderController) Show(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	ctx, user, _, err := common.GetContextAuthUser(w, r, []string{"despatch:read"}, dc.ServerOpt.Auth0Audience, dc.ServerOpt.Auth0Domain, dc.UserServiceClient)
	if err != nil {
		common.RenderErrorJSON(w, "1001", err.Error(), 401, user.RequestId)
		return
	}
	despatchHeader, err := dc.DespatchServiceClient.GetDespatchHeader(ctx, &logisticsproto.GetDespatchHeaderRequest{GetRequest: &commonproto.GetRequest{Id: id, UserEmail: user.Email, RequestId: user.RequestId}})
	if err != nil {
		dc.log.Error("Error",
			zap.String("reqid", user.RequestId),
			zap.Error(err))
		common.RenderErrorJSON(w, "1103", err.Error(), 402, user.RequestId)
		return
	}
	common.RenderJSON(w, despatchHeader)
}

// GetDespatchLines - Get Despatch Lines
func (dc *DespatchHeaderController) GetDespatchLines(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	ctx, user, _, err := common.GetContextAuthUser(w, r, []string{"despatch:read"}, dc.ServerOpt.Auth0Audience, dc.ServerOpt.Auth0Domain, dc.UserServiceClient)
	if err != nil {
		common.RenderErrorJSON(w, "1001", err.Error(), 401, user.RequestId)
		return
	}
	despatchLines, err := dc.DespatchServiceClient.GetDespatchLines(ctx, &logisticsproto.GetDespatchLinesRequest{GetRequest: &commonproto.GetRequest{Id: id, UserEmail: user.Email, RequestId: user.RequestId}})
	if err != nil {
		dc.log.Error("Error",
			zap.String("reqid", user.RequestId),

			zap.Error(err))
		common.RenderErrorJSON(w, "1103", err.Error(), 402, user.RequestId)
		return
	}
	common.RenderJSON(w, despatchLines)
}

// UpdateDespatchHeader - Update DespatchHeader
func (dc *DespatchHeaderController) UpdateDespatchHeader(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	ctx, user, token, err := common.GetContextAuthUser(w, r, []string{"despatch:cud"}, dc.ServerOpt.Auth0Audience, dc.ServerOpt.Auth0Domain, dc.UserServiceClient)
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

	form := logisticsproto.UpdateDespatchHeaderRequest{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&form)
	if err != nil {
		dc.log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		common.RenderErrorJSON(w, "4009", err.Error(), 402, user.RequestId)
		return
	}
	form.Id = id
	form.UserId = user.UserId
	form.UserEmail = user.Email
	form.RequestId = user.RequestId

	wHelper := dc.wfHelper
	result := wHelper.StartWorkflow(workflowOptions, logisticsworkflows.UpdateDespatchHeaderWorkflow, &form, token, user, dc.log)
	workflowClient := dc.workflowClient
	workflowRun := workflowClient.GetWorkflow(ctx, result.ID, result.RunID)
	var response string
	err = workflowRun.Get(ctx, &response)
	if err != nil {
		dc.log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		common.RenderErrorJSON(w, "4009", err.Error(), 402, user.RequestId)
		return
	}
	common.RenderJSON(w, response)
}
