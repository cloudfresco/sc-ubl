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

// CreditNoteHeaderController - Create CreditNoteHeader Controller
type CreditNoteHeaderController struct {
	log                           *zap.Logger
	UserServiceClient             partyproto.UserServiceClient
	CreditNoteHeaderServiceClient invoiceproto.CreditNoteHeaderServiceClient
	wfHelper                      common.WfHelper
	workflowClient                client.Client
	ServerOpt                     *config.ServerOptions
}

// NewCreditNoteHeaderController - Create CreditNoteHeader Handler
func NewCreditNoteHeaderController(log *zap.Logger, userServiceClient partyproto.UserServiceClient, creditNoteHeaderServiceClient invoiceproto.CreditNoteHeaderServiceClient, wfHelper common.WfHelper, workflowClient client.Client, serverOpt *config.ServerOptions) *CreditNoteHeaderController {
	return &CreditNoteHeaderController{
		log:                           log,
		UserServiceClient:             userServiceClient,
		CreditNoteHeaderServiceClient: creditNoteHeaderServiceClient,
		wfHelper:                      wfHelper,
		workflowClient:                workflowClient,
		ServerOpt:                     serverOpt,
	}
}

// CreateCreditNoteHeader - Create Credit Order Header
func (cc *CreditNoteHeaderController) CreateCreditNoteHeader(w http.ResponseWriter, r *http.Request) {
	ctx, user, token, err := common.GetContextAuthUser(w, r, []string{"creditnote:cud"}, cc.ServerOpt.Auth0Audience, cc.ServerOpt.Auth0Domain, cc.UserServiceClient)
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

	form := invoiceproto.CreateCreditNoteHeaderRequest{}
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
	result := wHelper.StartWorkflow(workflowOptions, invoiceworkflows.CreateCreditNoteHeaderWorkflow, &form, token, user, cc.log)
	workflowClient := cc.workflowClient
	workflowRun := workflowClient.GetWorkflow(ctx, result.ID, result.RunID)
	var creditNoteHeader invoiceproto.CreateCreditNoteHeaderResponse
	err = workflowRun.Get(ctx, &creditNoteHeader)

	if err != nil {
		cc.log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		common.RenderErrorJSON(w, "4002", err.Error(), 402, user.RequestId)
		return
	}

	common.RenderJSON(w, creditNoteHeader)
}

// Index - list CrediNotes
func (cc *CreditNoteHeaderController) Index(w http.ResponseWriter, r *http.Request) {
	ctx, user, _, err := common.GetContextAuthUser(w, r, []string{"creditnote:read"}, cc.ServerOpt.Auth0Audience, cc.ServerOpt.Auth0Domain, cc.UserServiceClient)
	if err != nil {
		common.RenderErrorJSON(w, "1001", err.Error(), 401, user.RequestId)
		return
	}

	cursor := r.URL.Query().Get("cursor")
	limit := r.URL.Query().Get("limit")

	creditNoteHeaders, err := cc.CreditNoteHeaderServiceClient.GetCreditNoteHeaders(ctx, &invoiceproto.GetCreditNoteHeadersRequest{Limit: limit, NextCursor: cursor, UserEmail: user.Email, RequestId: user.RequestId})
	if err != nil {
		cc.log.Error("Error",
			zap.String("user", user.Email),
			zap.String("reqid", user.RequestId),
			zap.Error(err))
		common.RenderErrorJSON(w, "1301", err.Error(), 402, user.RequestId)
		return
	}

	common.RenderJSON(w, creditNoteHeaders)
}

// Show - Show CreditNote
func (cc *CreditNoteHeaderController) Show(w http.ResponseWriter, r *http.Request) {
	ctx, user, _, err := common.GetContextAuthUser(w, r, []string{"creditnote:read"}, cc.ServerOpt.Auth0Audience, cc.ServerOpt.Auth0Domain, cc.UserServiceClient)
	if err != nil {
		common.RenderErrorJSON(w, "1001", err.Error(), 401, user.RequestId)
		return
	}

	id := r.PathValue("id")

	creditNoteHeader, err := cc.CreditNoteHeaderServiceClient.GetCreditNoteHeader(ctx, &invoiceproto.GetCreditNoteHeaderRequest{GetRequest: &commonproto.GetRequest{Id: id, UserEmail: user.Email, RequestId: user.RequestId}})
	if err != nil {
		cc.log.Error("Error",
			zap.String("reqid", user.RequestId),
			zap.Error(err))
		common.RenderErrorJSON(w, "1103", err.Error(), 402, user.RequestId)
		return
	}
	common.RenderJSON(w, creditNoteHeader)
}

// UpdateCreditNoteHeader - Update CreditNoteHeader
func (cc *CreditNoteHeaderController) UpdateCreditNoteHeader(w http.ResponseWriter, r *http.Request) {
	ctx, user, token, err := common.GetContextAuthUser(w, r, []string{"creditnote:cud"}, cc.ServerOpt.Auth0Audience, cc.ServerOpt.Auth0Domain, cc.UserServiceClient)
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

	form := invoiceproto.UpdateCreditNoteHeaderRequest{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&form)
	if err != nil {
		cc.log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		common.RenderErrorJSON(w, "4009", err.Error(), 402, user.RequestId)
		return
	}
	form.Id = id
	form.UserId = user.UserId
	form.UserEmail = user.Email
	form.RequestId = user.RequestId

	wHelper := cc.wfHelper
	result := wHelper.StartWorkflow(workflowOptions, invoiceworkflows.UpdateCreditNoteHeaderWorkflow, &form, token, user, cc.log)
	workflowClient := cc.workflowClient
	workflowRun := workflowClient.GetWorkflow(ctx, result.ID, result.RunID)
	var response string
	err = workflowRun.Get(ctx, &response)
	if err != nil {
		cc.log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		common.RenderErrorJSON(w, "4009", err.Error(), 402, user.RequestId)
		return
	}
	common.RenderJSON(w, response)
}

// GetCreditNoteLines - Get GetCreditNoteLines
func (cc *CreditNoteHeaderController) GetCreditNoteLines(w http.ResponseWriter, r *http.Request) {
	ctx, user, _, err := common.GetContextAuthUser(w, r, []string{"creditnote:read"}, cc.ServerOpt.Auth0Audience, cc.ServerOpt.Auth0Domain, cc.UserServiceClient)
	if err != nil {
		common.RenderErrorJSON(w, "1001", err.Error(), 401, user.RequestId)
		return
	}

	id := r.PathValue("id")

	creditNoteLines, err := cc.CreditNoteHeaderServiceClient.GetCreditNoteLines(ctx, &invoiceproto.GetCreditNoteLinesRequest{GetRequest: &commonproto.GetRequest{Id: id, UserEmail: user.Email, RequestId: user.RequestId}})
	if err != nil {
		cc.log.Error("Error",
			zap.String("reqid", user.RequestId),
			zap.Error(err))
		common.RenderErrorJSON(w, "1103", err.Error(), 402, user.RequestId)
		return
	}
	common.RenderJSON(w, creditNoteLines)
}
