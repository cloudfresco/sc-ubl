package invoicecontrollers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/cloudfresco/sc-ubl/internal/common"
	"github.com/cloudfresco/sc-ubl/internal/config"
	commonproto "github.com/cloudfresco/sc-ubl/internal/protogen/common/v1"
	invoiceproto "github.com/cloudfresco/sc-ubl/internal/protogen/invoice/v1"
	partyproto "github.com/cloudfresco/sc-ubl/internal/protogen/party/v1"
	invoiceworkflows "github.com/cloudfresco/sc-ubl/internal/workflows/invoiceworkflows"
	"github.com/pborman/uuid"
	"go.uber.org/cadence/client"
	"go.uber.org/zap"
)

// DebitNoteHeaderController - Create DebitNoteHeader Controller
type DebitNoteHeaderController struct {
	log                          *zap.Logger
	UserServiceClient            partyproto.UserServiceClient
	DebitNoteHeaderServiceClient invoiceproto.DebitNoteHeaderServiceClient
	wfHelper                     common.WfHelper
	workflowClient               client.Client
	ServerOpt                    *config.ServerOptions
}

// NewDebitNoteHeaderController - Create DebitNoteHeader Handler
func NewDebitNoteHeaderController(log *zap.Logger, userServiceClient partyproto.UserServiceClient, debitNoteHeaderServiceClient invoiceproto.DebitNoteHeaderServiceClient, wfHelper common.WfHelper, workflowClient client.Client, serverOpt *config.ServerOptions) *DebitNoteHeaderController {
	return &DebitNoteHeaderController{
		log:                          log,
		UserServiceClient:            userServiceClient,
		DebitNoteHeaderServiceClient: debitNoteHeaderServiceClient,
		wfHelper:                     wfHelper,
		workflowClient:               workflowClient,
		ServerOpt:                    serverOpt,
	}
}

// CreateDebitNoteHeader - Create Debit Order Header
func (dc *DebitNoteHeaderController) CreateDebitNoteHeader(w http.ResponseWriter, r *http.Request) {
	ctx, user, token, err := common.GetContextAuthUser(w, r, []string{"debitnote:cud"}, dc.ServerOpt.Auth0Audience, dc.ServerOpt.Auth0Domain, dc.UserServiceClient)
	if err != nil {
		common.RenderErrorJSON(w, "1001", err.Error(), 401, user.RequestId)
		return
	}

	workflowOptions := client.StartWorkflowOptions{
		ID:                              "ubl_" + uuid.New(),
		TaskList:                        invoiceworkflows.ApplicationName,
		ExecutionStartToCloseTimeout:    time.Minute,
		DecisionTaskStartToCloseTimeout: time.Minute,
	}
	form := invoiceproto.CreateDebitNoteHeaderRequest{}
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
	result := wHelper.StartWorkflow(workflowOptions, invoiceworkflows.CreateDebitNoteHeaderWorkflow, &form, token, user, dc.log)
	workflowClient := dc.workflowClient
	workflowRun := workflowClient.GetWorkflow(ctx, result.ID, result.RunID)
	var debitNoteHeader invoiceproto.CreateDebitNoteHeaderResponse
	err = workflowRun.Get(ctx, &debitNoteHeader)

	if err != nil {
		dc.log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		common.RenderErrorJSON(w, "4002", err.Error(), 402, user.RequestId)
		return
	}

	common.RenderJSON(w, debitNoteHeader)
}

// Index - list DebitNotes
func (dc *DebitNoteHeaderController) Index(w http.ResponseWriter, r *http.Request) {
	ctx, user, _, err := common.GetContextAuthUser(w, r, []string{"debitnote:read"}, dc.ServerOpt.Auth0Audience, dc.ServerOpt.Auth0Domain, dc.UserServiceClient)
	if err != nil {
		common.RenderErrorJSON(w, "1001", err.Error(), 401, user.RequestId)
		return
	}

	cursor := r.URL.Query().Get("cursor")
	limit := r.URL.Query().Get("limit")

	resp, err := dc.DebitNoteHeaderServiceClient.GetDebitNoteHeaders(ctx, &invoiceproto.GetDebitNoteHeadersRequest{Limit: limit, NextCursor: cursor, UserEmail: user.Email, RequestId: user.RequestId})
	if err != nil {
		dc.log.Error("Error",
			zap.String("user", user.Email),
			zap.String("reqid", user.RequestId),

			zap.Error(err))
		common.RenderErrorJSON(w, "1301", err.Error(), 402, user.RequestId)
		return
	}

	common.RenderJSON(w, resp)
}

// Show - Show DebitNote
func (dc *DebitNoteHeaderController) Show(w http.ResponseWriter, r *http.Request) {
	ctx, user, _, err := common.GetContextAuthUser(w, r, []string{"debitnote:read"}, dc.ServerOpt.Auth0Audience, dc.ServerOpt.Auth0Domain, dc.UserServiceClient)
	if err != nil {
		common.RenderErrorJSON(w, "1001", err.Error(), 401, user.RequestId)
		return
	}

	id := r.PathValue("id")

	resp, err := dc.DebitNoteHeaderServiceClient.GetDebitNoteHeader(ctx, &invoiceproto.GetDebitNoteHeaderRequest{GetRequest: &commonproto.GetRequest{Id: id, UserEmail: user.Email, RequestId: user.RequestId}})
	if err != nil {
		dc.log.Error("Error",
			zap.String("reqid", user.RequestId),

			zap.Error(err))
		common.RenderErrorJSON(w, "1103", err.Error(), 402, user.RequestId)
		return
	}
	common.RenderJSON(w, resp)
}

// UpdateDebitNoteHeader - Update DebitNoteHeader
func (dc *DebitNoteHeaderController) UpdateDebitNoteHeader(w http.ResponseWriter, r *http.Request) {
	ctx, user, token, err := common.GetContextAuthUser(w, r, []string{"debitnote:cud"}, dc.ServerOpt.Auth0Audience, dc.ServerOpt.Auth0Domain, dc.UserServiceClient)
	if err != nil {
		common.RenderErrorJSON(w, "1001", err.Error(), 401, user.RequestId)
		return
	}

	id := r.PathValue("id")

	workflowOptions := client.StartWorkflowOptions{
		ID:                              "ubl_" + uuid.New(),
		TaskList:                        invoiceworkflows.ApplicationName,
		ExecutionStartToCloseTimeout:    time.Minute,
		DecisionTaskStartToCloseTimeout: time.Minute,
	}

	form := invoiceproto.UpdateDebitNoteHeaderRequest{}
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
	result := wHelper.StartWorkflow(workflowOptions, invoiceworkflows.UpdateDebitNoteHeaderWorkflow, &form, token, user, dc.log)
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

// GetDebitNoteLines - Get DebitNote Lines
func (dc *DebitNoteHeaderController) GetDebitNoteLines(w http.ResponseWriter, r *http.Request) {
	ctx, user, _, err := common.GetContextAuthUser(w, r, []string{"debitnote:read"}, dc.ServerOpt.Auth0Audience, dc.ServerOpt.Auth0Domain, dc.UserServiceClient)
	if err != nil {
		common.RenderErrorJSON(w, "1001", err.Error(), 401, user.RequestId)
		return
	}

	id := r.PathValue("id")

	resp, err := dc.DebitNoteHeaderServiceClient.GetDebitNoteLines(ctx, &invoiceproto.GetDebitNoteLinesRequest{GetRequest: &commonproto.GetRequest{Id: id, UserEmail: user.Email, RequestId: user.RequestId}})
	if err != nil {
		dc.log.Error("Error",
			zap.String("reqid", user.RequestId),

			zap.Error(err))
		common.RenderErrorJSON(w, "1103", err.Error(), 402, user.RequestId)
		return
	}
	common.RenderJSON(w, resp)
}
