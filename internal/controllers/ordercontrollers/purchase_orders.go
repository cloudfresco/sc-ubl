package ordercontrollers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/cloudfresco/sc-ubl/internal/common"
	"github.com/cloudfresco/sc-ubl/internal/config"
	commonproto "github.com/cloudfresco/sc-ubl/internal/protogen/common/v1"
	orderproto "github.com/cloudfresco/sc-ubl/internal/protogen/order/v1"
	partyproto "github.com/cloudfresco/sc-ubl/internal/protogen/party/v1"
	"github.com/cloudfresco/sc-ubl/internal/workflows/orderworkflows"
	"github.com/pborman/uuid"
	"go.uber.org/cadence/client"
	"go.uber.org/zap"
)

// PurchaseOrderHeaderController - Create PurchaseOrderHeader Controller
type PurchaseOrderHeaderController struct {
	log                              *zap.Logger
	UserServiceClient                partyproto.UserServiceClient
	PurchaseOrderHeaderServiceClient orderproto.PurchaseOrderHeaderServiceClient
	wfHelper                         common.WfHelper
	workflowClient                   client.Client
	ServerOpt                        *config.ServerOptions
}

// NewPurchaseOrderHeaderController - Create PurchaseOrderHeader Handler
func NewPurchaseOrderHeaderController(log *zap.Logger, userServiceClient partyproto.UserServiceClient, purchaseOrderHeaderServiceClient orderproto.PurchaseOrderHeaderServiceClient, wfHelper common.WfHelper, workflowClient client.Client, serverOpt *config.ServerOptions) *PurchaseOrderHeaderController {
	return &PurchaseOrderHeaderController{
		log:                              log,
		UserServiceClient:                userServiceClient,
		PurchaseOrderHeaderServiceClient: purchaseOrderHeaderServiceClient,
		wfHelper:                         wfHelper,
		workflowClient:                   workflowClient,
		ServerOpt:                        serverOpt,
	}
}

// CreatePurchaseOrderHeader - Create Purchase Order Header
func (pc *PurchaseOrderHeaderController) CreatePurchaseOrderHeader(w http.ResponseWriter, r *http.Request) {
	ctx, user, token, err := common.GetContextAuthUser(w, r, []string{"po:cud"}, pc.ServerOpt.Auth0Audience, pc.ServerOpt.Auth0Domain, pc.UserServiceClient)
	if err != nil {
		common.RenderErrorJSON(w, "1001", err.Error(), 401, user.RequestId)
		return
	}

	workflowOptions := client.StartWorkflowOptions{
		ID:                              "ubl_" + uuid.New(),
		TaskList:                        orderworkflows.ApplicationName,
		ExecutionStartToCloseTimeout:    time.Minute,
		DecisionTaskStartToCloseTimeout: time.Minute,
	}

	form := orderproto.CreatePurchaseOrderHeaderRequest{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&form)
	if err != nil {
		pc.log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		common.RenderErrorJSON(w, "4002", err.Error(), 402, user.RequestId)
		return
	}
	form.UserId = user.UserId
	form.UserEmail = user.Email
	form.RequestId = user.RequestId

	wHelper := pc.wfHelper
	result := wHelper.StartWorkflow(workflowOptions, orderworkflows.CreatePurchaseOrderHeaderWorkflow, &form, token, user, pc.log)
	workflowClient := pc.workflowClient
	workflowRun := workflowClient.GetWorkflow(ctx, result.ID, result.RunID)
	var purchaseOrderHeader orderproto.CreatePurchaseOrderHeaderResponse
	err = workflowRun.Get(ctx, &purchaseOrderHeader)
	if err != nil {
		pc.log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		common.RenderErrorJSON(w, "4002", err.Error(), 402, user.RequestId)
		return
	}
	common.RenderJSON(w, purchaseOrderHeader)
}

// Index - list PurchaseOrderHeaders
func (pc *PurchaseOrderHeaderController) Index(w http.ResponseWriter, r *http.Request) {
	ctx, user, _, err := common.GetContextAuthUser(w, r, []string{"po:read"}, pc.ServerOpt.Auth0Audience, pc.ServerOpt.Auth0Domain, pc.UserServiceClient)
	if err != nil {
		common.RenderErrorJSON(w, "1001", err.Error(), 401, user.RequestId)
		return
	}

	cursor := r.URL.Query().Get("cursor")
	limit := r.URL.Query().Get("limit")

	purchaseOrderHeaders, err := pc.PurchaseOrderHeaderServiceClient.GetPurchaseOrderHeaders(ctx, &orderproto.GetPurchaseOrderHeadersRequest{Limit: limit, NextCursor: cursor, UserEmail: user.Email, RequestId: user.RequestId})
	if err != nil {
		pc.log.Error("Error",
			zap.String("user", user.Email),
			zap.String("reqid", user.RequestId),
			zap.Error(err))
		common.RenderErrorJSON(w, "1301", err.Error(), 402, user.RequestId)
		return
	}

	common.RenderJSON(w, purchaseOrderHeaders)
}

// Show - Show PurchaseOrder
func (pc *PurchaseOrderHeaderController) Show(w http.ResponseWriter, r *http.Request) {
	ctx, user, _, err := common.GetContextAuthUser(w, r, []string{"po:read"}, pc.ServerOpt.Auth0Audience, pc.ServerOpt.Auth0Domain, pc.UserServiceClient)
	if err != nil {
		common.RenderErrorJSON(w, "1001", err.Error(), 401, user.RequestId)
		return
	}

	id := r.PathValue("id")

	purchaseOrderHeader, err := pc.PurchaseOrderHeaderServiceClient.GetPurchaseOrderHeader(ctx, &orderproto.GetPurchaseOrderHeaderRequest{GetRequest: &commonproto.GetRequest{Id: id, UserEmail: user.Email, RequestId: user.RequestId}})
	if err != nil {
		pc.log.Error("Error",
			zap.String("reqid", user.RequestId),
			zap.Error(err))
		common.RenderErrorJSON(w, "1103", err.Error(), 402, user.RequestId)
		return
	}
	common.RenderJSON(w, purchaseOrderHeader)
}

// GetPurchaseOrderLines - Get PurchaseOrder Lines
func (pc *PurchaseOrderHeaderController) GetPurchaseOrderLines(w http.ResponseWriter, r *http.Request) {
	ctx, user, _, err := common.GetContextAuthUser(w, r, []string{"po:read"}, pc.ServerOpt.Auth0Audience, pc.ServerOpt.Auth0Domain, pc.UserServiceClient)
	if err != nil {
		common.RenderErrorJSON(w, "1001", err.Error(), 401, user.RequestId)
		return
	}

	id := r.PathValue("id")

	purchaseOrderLines, err := pc.PurchaseOrderHeaderServiceClient.GetPurchaseOrderLines(ctx, &orderproto.GetPurchaseOrderLinesRequest{GetRequest: &commonproto.GetRequest{Id: id, UserEmail: user.Email, RequestId: user.RequestId}})
	if err != nil {
		pc.log.Error("Error",
			zap.String("reqid", user.RequestId),
			zap.Error(err))
		common.RenderErrorJSON(w, "1103", err.Error(), 402, user.RequestId)
		return
	}
	common.RenderJSON(w, purchaseOrderLines)
}

// UpdatePurchaseOrderHeader - Update PurchaseOrderHeader
func (pc *PurchaseOrderHeaderController) UpdatePurchaseOrderHeader(w http.ResponseWriter, r *http.Request) {
	ctx, user, token, err := common.GetContextAuthUser(w, r, []string{"po:cud"}, pc.ServerOpt.Auth0Audience, pc.ServerOpt.Auth0Domain, pc.UserServiceClient)
	if err != nil {
		common.RenderErrorJSON(w, "1001", err.Error(), 401, user.RequestId)
		return
	}

	id := r.PathValue("id")

	workflowOptions := client.StartWorkflowOptions{
		ID:                              "ubl_" + uuid.New(),
		TaskList:                        orderworkflows.ApplicationName,
		ExecutionStartToCloseTimeout:    time.Minute,
		DecisionTaskStartToCloseTimeout: time.Minute,
	}

	form := orderproto.UpdatePurchaseOrderHeaderRequest{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&form)
	if err != nil {
		pc.log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		common.RenderErrorJSON(w, "4009", err.Error(), 402, user.RequestId)
		return
	}
	form.Id = id
	form.UserId = user.UserId
	form.UserEmail = user.Email
	form.RequestId = user.RequestId

	wHelper := pc.wfHelper
	result := wHelper.StartWorkflow(workflowOptions, orderworkflows.UpdatePurchaseOrderHeaderWorkflow, &form, token, user, pc.log)
	workflowClient := pc.workflowClient
	workflowRun := workflowClient.GetWorkflow(ctx, result.ID, result.RunID)
	var response string
	err = workflowRun.Get(ctx, &response)
	if err != nil {
		pc.log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		common.RenderErrorJSON(w, "4009", err.Error(), 402, user.RequestId)
		return
	}
	common.RenderJSON(w, response)
}
