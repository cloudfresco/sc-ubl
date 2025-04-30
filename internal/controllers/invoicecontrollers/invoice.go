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

// InvoiceHeaderController - Create InvoiceHeader Controller
type InvoiceHeaderController struct {
	log                  *zap.Logger
	UserServiceClient    partyproto.UserServiceClient
	InvoiceServiceClient invoiceproto.InvoiceServiceClient
	wfHelper             common.WfHelper
	workflowClient       client.Client
	ServerOpt            *config.ServerOptions
}

// NewInvoiceHeaderController - Create Invoice Handler
func NewInvoiceHeaderController(log *zap.Logger, userServiceClient partyproto.UserServiceClient, invoiceHeaderServiceClient invoiceproto.InvoiceServiceClient, wfHelper common.WfHelper, workflowClient client.Client, serverOpt *config.ServerOptions) *InvoiceHeaderController {
	return &InvoiceHeaderController{
		log:                  log,
		UserServiceClient:    userServiceClient,
		InvoiceServiceClient: invoiceHeaderServiceClient,
		wfHelper:             wfHelper,
		workflowClient:       workflowClient,
		ServerOpt:            serverOpt,
	}
}

// CreateInvoice - Create Invoice
func (ic *InvoiceHeaderController) CreateInvoice(w http.ResponseWriter, r *http.Request) {
	ctx, user, token, err := common.GetContextAuthUser(w, r, []string{"invoice:cud"}, ic.ServerOpt.Auth0Audience, ic.ServerOpt.Auth0Domain, ic.UserServiceClient)
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

	form := invoiceproto.CreateInvoiceRequest{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&form)
	if err != nil {
		ic.log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		common.RenderErrorJSON(w, "4002", err.Error(), 402, user.RequestId)
		return
	}
	form.UserId = user.UserId
	form.UserEmail = user.Email
	form.RequestId = user.RequestId

	wHelper := ic.wfHelper
	result := wHelper.StartWorkflow(workflowOptions, invoiceworkflows.CreateInvoiceWorkflow, &form, token, user, ic.log)
	workflowClient := ic.workflowClient
	workflowRun := workflowClient.GetWorkflow(ctx, result.ID, result.RunID)
	var invoice invoiceproto.CreateInvoiceResponse
	err = workflowRun.Get(ctx, &invoice)

	if err != nil {
		ic.log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		common.RenderErrorJSON(w, "4002", err.Error(), 402, user.RequestId)
		return
	}

	common.RenderJSON(w, invoice)
}

// Index - list Invoice
func (ic *InvoiceHeaderController) Index(w http.ResponseWriter, r *http.Request) {
	ctx, user, _, err := common.GetContextAuthUser(w, r, []string{"invoice:read"}, ic.ServerOpt.Auth0Audience, ic.ServerOpt.Auth0Domain, ic.UserServiceClient)
	if err != nil {
		common.RenderErrorJSON(w, "1001", err.Error(), 401, user.RequestId)
		return
	}

	cursor := r.URL.Query().Get("cursor")
	limit := r.URL.Query().Get("limit")

	invoice, err := ic.InvoiceServiceClient.GetInvoices(ctx, &invoiceproto.GetInvoicesRequest{Limit: limit, NextCursor: cursor, UserEmail: user.Email, RequestId: user.RequestId})
	if err != nil {
		ic.log.Error("Error",
			zap.String("user", user.Email),
			zap.String("reqid", user.RequestId),
			zap.Error(err))
		common.RenderErrorJSON(w, "1301", err.Error(), 402, user.RequestId)
		return
	}

	common.RenderJSON(w, invoice)
}

// Show - Show Invoice
func (ic *InvoiceHeaderController) Show(w http.ResponseWriter, r *http.Request) {
	ctx, user, _, err := common.GetContextAuthUser(w, r, []string{"invoice:read"}, ic.ServerOpt.Auth0Audience, ic.ServerOpt.Auth0Domain, ic.UserServiceClient)
	if err != nil {
		common.RenderErrorJSON(w, "1001", err.Error(), 401, user.RequestId)
		return
	}
	id := r.PathValue("id")

	invoice, err := ic.InvoiceServiceClient.GetInvoice(ctx, &invoiceproto.GetInvoiceRequest{GetRequest: &commonproto.GetRequest{Id: id, UserEmail: user.Email, RequestId: user.RequestId}})
	if err != nil {
		ic.log.Error("Error",
			zap.String("reqid", user.RequestId),

			zap.Error(err))
		common.RenderErrorJSON(w, "1103", err.Error(), 402, user.RequestId)
		return
	}
	common.RenderJSON(w, invoice)
}

func (ic *InvoiceHeaderController) UpdateInvoice(w http.ResponseWriter, r *http.Request) {
	ctx, user, token, err := common.GetContextAuthUser(w, r, []string{"invoice:cud"}, ic.ServerOpt.Auth0Audience, ic.ServerOpt.Auth0Domain, ic.UserServiceClient)
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

	form := invoiceproto.UpdateInvoiceRequest{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&form)
	if err != nil {
		ic.log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		common.RenderErrorJSON(w, "4009", err.Error(), 402, user.RequestId)
		return
	}
	form.Id = id
	form.UserId = user.UserId
	form.UserEmail = user.Email
	form.RequestId = user.RequestId

	wHelper := ic.wfHelper
	result := wHelper.StartWorkflow(workflowOptions, invoiceworkflows.UpdateInvoiceWorkflow, &form, token, user, ic.log)
	workflowClient := ic.workflowClient
	workflowRun := workflowClient.GetWorkflow(ctx, result.ID, result.RunID)
	var response string
	err = workflowRun.Get(ctx, &response)
	if err != nil {
		ic.log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		common.RenderErrorJSON(w, "4009", err.Error(), 402, user.RequestId)
		return
	}
	common.RenderJSON(w, response)
}

// GetInvoiceLines - Get Invoice Lines
func (ic *InvoiceHeaderController) GetInvoiceLines(w http.ResponseWriter, r *http.Request) {
	ctx, user, _, err := common.GetContextAuthUser(w, r, []string{"invoice:read"}, ic.ServerOpt.Auth0Audience, ic.ServerOpt.Auth0Domain, ic.UserServiceClient)
	if err != nil {
		common.RenderErrorJSON(w, "1001", err.Error(), 401, user.RequestId)
		return
	}

	id := r.PathValue("id")

	invoiceLines, err := ic.InvoiceServiceClient.GetInvoiceLines(ctx, &invoiceproto.GetInvoiceLinesRequest{GetRequest: &commonproto.GetRequest{Id: id, UserEmail: user.Email, RequestId: user.RequestId}})
	if err != nil {
		ic.log.Error("Error",
			zap.String("reqid", user.RequestId),
			zap.Error(err))
		common.RenderErrorJSON(w, "1103", err.Error(), 402, user.RequestId)
		return
	}
	common.RenderJSON(w, invoiceLines)
}
