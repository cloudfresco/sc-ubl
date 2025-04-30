package taxcontrollers

import (
	"encoding/json"
	"net/http"
	"time"

	//"github.com/bufbuild/protovalidate-go"
	"github.com/cloudfresco/sc-ubl/internal/common"
	"github.com/cloudfresco/sc-ubl/internal/config"
	commonproto "github.com/cloudfresco/sc-ubl/internal/protogen/common/v1"
	partyproto "github.com/cloudfresco/sc-ubl/internal/protogen/party/v1"
	taxproto "github.com/cloudfresco/sc-ubl/internal/protogen/tax/v1"
	"github.com/cloudfresco/sc-ubl/internal/workflows/taxworkflows"
	"github.com/pborman/uuid"
	"go.uber.org/cadence/client"
	"go.uber.org/zap"
)

// TaxController - Create Tax Controller
type TaxController struct {
	log               *zap.Logger
	UserServiceClient partyproto.UserServiceClient
	TaxServiceClient  taxproto.TaxServiceClient
	wfHelper          common.WfHelper
	workflowClient    client.Client
	ServerOpt         *config.ServerOptions
}

// NewTaxController - Create Tax Handler
func NewTaxController(log *zap.Logger, userServiceClient partyproto.UserServiceClient, taxServiceClient taxproto.TaxServiceClient, wfHelper common.WfHelper, workflowClient client.Client, serverOpt *config.ServerOptions) *TaxController {
	return &TaxController{
		log:               log,
		UserServiceClient: userServiceClient,
		TaxServiceClient:  taxServiceClient,
		wfHelper:          wfHelper,
		workflowClient:    workflowClient,
		ServerOpt:         serverOpt,
	}
}

// Index - list TaxScheme
func (tc *TaxController) Index(w http.ResponseWriter, r *http.Request) {
	ctx, user, _, err := common.GetContextAuthUser(w, r, []string{"tax:read"}, tc.ServerOpt.Auth0Audience, tc.ServerOpt.Auth0Domain, tc.UserServiceClient)
	if err != nil {
		common.RenderErrorJSON(w, "1001", err.Error(), 401, user.RequestId)
		return
	}

	cursor := r.URL.Query().Get("cursor")
	limit := r.URL.Query().Get("limit")

	taxSchemes, err := tc.TaxServiceClient.GetTaxSchemes(ctx, &taxproto.GetTaxSchemesRequest{Limit: limit, NextCursor: cursor, UserEmail: user.Email, RequestId: user.RequestId})
	if err != nil {
		tc.log.Error("Error",
			zap.String("user", user.Email),
			zap.String("reqid", user.RequestId),
			zap.Error(err))
		common.RenderErrorJSON(w, "1301", err.Error(), 402, user.RequestId)
		return
	}

	common.RenderJSON(w, taxSchemes)
}

// Show - Show TaxScheme
func (tc *TaxController) Show(w http.ResponseWriter, r *http.Request) {
	ctx, user, _, err := common.GetContextAuthUser(w, r, []string{"tax:read"}, tc.ServerOpt.Auth0Audience, tc.ServerOpt.Auth0Domain, tc.UserServiceClient)
	if err != nil {
		common.RenderErrorJSON(w, "1001", err.Error(), 401, user.RequestId)
		return
	}

	id := r.PathValue("id")
	taxScheme, err := tc.TaxServiceClient.GetTaxScheme(ctx, &taxproto.GetTaxSchemeRequest{GetRequest: &commonproto.GetRequest{Id: id, UserEmail: user.Email, RequestId: user.RequestId}})
	if err != nil {
		tc.log.Error("Error",
			zap.String("reqid", user.RequestId),
			zap.Error(err))
		common.RenderErrorJSON(w, "1103", err.Error(), 402, user.RequestId)
		return
	}
	common.RenderJSON(w, taxScheme)
}

// CreateTaxScheme - Create TaxScheme
func (tc *TaxController) CreateTaxScheme(w http.ResponseWriter, r *http.Request) {
	ctx, user, token, err := common.GetContextAuthUser(w, r, []string{"tax:cud"}, tc.ServerOpt.Auth0Audience, tc.ServerOpt.Auth0Domain, tc.UserServiceClient)
	if err != nil {
		common.RenderErrorJSON(w, "1001", err.Error(), 401, user.RequestId)
		return
	}

	workflowOptions := client.StartWorkflowOptions{
		ID:                              "ubl_" + uuid.New(),
		TaskList:                        taxworkflows.ApplicationName,
		ExecutionStartToCloseTimeout:    time.Minute,
		DecisionTaskStartToCloseTimeout: time.Minute,
	}

	form := taxproto.CreateTaxSchemeRequest{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&form)
	if err != nil {
		tc.log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		common.RenderErrorJSON(w, "4002", err.Error(), 402, user.RequestId)
		return
	}
	form.UserId = user.UserId
	form.UserEmail = user.Email
	form.RequestId = user.RequestId

	wHelper := tc.wfHelper
	result := wHelper.StartWorkflow(workflowOptions, taxworkflows.CreateTaxSchemeWorkflow, &form, token, user, tc.log)
	workflowClient := tc.workflowClient
	workflowRun := workflowClient.GetWorkflow(ctx, result.ID, result.RunID)
	var taxScheme taxproto.CreateTaxSchemeResponse
	err = workflowRun.Get(ctx, &taxScheme)
	if err != nil {
		tc.log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		common.RenderErrorJSON(w, "4002", err.Error(), 402, user.RequestId)
		return
	}
	common.RenderJSON(w, taxScheme)
}

// UpdateTaxScheme - Update TaxScheme
func (tc *TaxController) UpdateTaxScheme(w http.ResponseWriter, r *http.Request) {
	ctx, user, token, err := common.GetContextAuthUser(w, r, []string{"tax:cud"}, tc.ServerOpt.Auth0Audience, tc.ServerOpt.Auth0Domain, tc.UserServiceClient)
	if err != nil {
		common.RenderErrorJSON(w, "1001", err.Error(), 401, user.RequestId)
		return
	}

	id := r.PathValue("id")

	workflowOptions := client.StartWorkflowOptions{
		ID:                              "ubl_" + uuid.New(),
		TaskList:                        taxworkflows.ApplicationName,
		ExecutionStartToCloseTimeout:    time.Minute,
		DecisionTaskStartToCloseTimeout: time.Minute,
	}

	form := taxproto.UpdateTaxSchemeRequest{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&form)
	if err != nil {
		tc.log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		common.RenderErrorJSON(w, "4009", err.Error(), 402, user.RequestId)
		return
	}
	form.Id = id
	form.UserId = user.UserId
	form.UserEmail = user.Email
	form.RequestId = user.RequestId

	wHelper := tc.wfHelper
	result := wHelper.StartWorkflow(workflowOptions, taxworkflows.UpdateTaxSchemeWorkflow, &form, token, user, tc.log)
	workflowClient := tc.workflowClient
	workflowRun := workflowClient.GetWorkflow(ctx, result.ID, result.RunID)
	var response string
	err = workflowRun.Get(ctx, &response)
	if err != nil {
		tc.log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		common.RenderErrorJSON(w, "4009", err.Error(), 402, user.RequestId)
		return
	}
	common.RenderJSON(w, response)
}

// CreateTaxSchemeJurisdiction - Create TaxSchemeJurisdiction
func (tc *TaxController) CreateTaxSchemeJurisdiction(w http.ResponseWriter, r *http.Request) {
	ctx, user, token, err := common.GetContextAuthUser(w, r, []string{"tax:cud"}, tc.ServerOpt.Auth0Audience, tc.ServerOpt.Auth0Domain, tc.UserServiceClient)
	if err != nil {
		common.RenderErrorJSON(w, "1001", err.Error(), 401, user.RequestId)
		return
	}

	id := r.PathValue("id")

	workflowOptions := client.StartWorkflowOptions{
		ID:                              "ubl_" + uuid.New(),
		TaskList:                        taxworkflows.ApplicationName,
		ExecutionStartToCloseTimeout:    time.Minute,
		DecisionTaskStartToCloseTimeout: time.Minute,
	}

	form := taxproto.CreateTaxSchemeJurisdictionRequest{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&form)
	if err != nil {
		tc.log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		common.RenderErrorJSON(w, "4002", err.Error(), 402, user.RequestId)
		return
	}

	form.TaxSchemeId = id
	form.UserId = user.UserId
	form.UserEmail = user.Email
	form.RequestId = user.RequestId

	wHelper := tc.wfHelper
	result := wHelper.StartWorkflow(workflowOptions, taxworkflows.CreateTaxSchemeJurisdictionWorkflow, &form, token, user, tc.log)
	workflowClient := tc.workflowClient
	workflowRun := workflowClient.GetWorkflow(ctx, result.ID, result.RunID)
	var taxSchemeJurisdiction taxproto.CreateTaxSchemeJurisdictionResponse
	err = workflowRun.Get(ctx, &taxSchemeJurisdiction)
	if err != nil {
		tc.log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		common.RenderErrorJSON(w, "4002", err.Error(), 402, user.RequestId)
		return
	}
	common.RenderJSON(w, taxSchemeJurisdiction)
}

// UpdateTaxSchemeJurisdiction - Update TaxSchemeJurisdiction
func (tc *TaxController) UpdateTaxSchemeJurisdiction(w http.ResponseWriter, r *http.Request) {
	ctx, user, token, err := common.GetContextAuthUser(w, r, []string{"tax:cud"}, tc.ServerOpt.Auth0Audience, tc.ServerOpt.Auth0Domain, tc.UserServiceClient)
	if err != nil {
		common.RenderErrorJSON(w, "1001", err.Error(), 401, user.RequestId)
		return
	}

	id := r.PathValue("id")

	workflowOptions := client.StartWorkflowOptions{
		ID:                              "ubl_" + uuid.New(),
		TaskList:                        taxworkflows.ApplicationName,
		ExecutionStartToCloseTimeout:    time.Minute,
		DecisionTaskStartToCloseTimeout: time.Minute,
	}

	form := taxproto.UpdateTaxSchemeJurisdictionRequest{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&form)
	if err != nil {
		tc.log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		common.RenderErrorJSON(w, "4009", err.Error(), 402, user.RequestId)
		return
	}
	form.Id = id
	form.UserId = user.UserId
	form.UserEmail = user.Email
	form.RequestId = user.RequestId

	wHelper := tc.wfHelper
	result := wHelper.StartWorkflow(workflowOptions, taxworkflows.UpdateTaxSchemeJurisdictionWorkflow, &form, token, user, tc.log)
	workflowClient := tc.workflowClient
	workflowRun := workflowClient.GetWorkflow(ctx, result.ID, result.RunID)
	var response string
	err = workflowRun.Get(ctx, &response)
	if err != nil {
		tc.log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		common.RenderErrorJSON(w, "4009", err.Error(), 402, user.RequestId)
		return
	}
	common.RenderJSON(w, response)
}

// CreateTaxCategory - Create TaxCategory
func (tc *TaxController) CreateTaxCategory(w http.ResponseWriter, r *http.Request) {
	ctx, user, token, err := common.GetContextAuthUser(w, r, []string{"tax:cud"}, tc.ServerOpt.Auth0Audience, tc.ServerOpt.Auth0Domain, tc.UserServiceClient)
	if err != nil {
		common.RenderErrorJSON(w, "1001", err.Error(), 401, user.RequestId)
		return
	}

	id := r.PathValue("id")

	workflowOptions := client.StartWorkflowOptions{
		ID:                              "ubl_" + uuid.New(),
		TaskList:                        taxworkflows.ApplicationName,
		ExecutionStartToCloseTimeout:    time.Minute,
		DecisionTaskStartToCloseTimeout: time.Minute,
	}

	form := taxproto.CreateTaxCategoryRequest{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&form)
	if err != nil {
		tc.log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		common.RenderErrorJSON(w, "4002", err.Error(), 402, user.RequestId)
		return
	}

	form.TaxSchemeId = id
	form.UserId = user.UserId
	form.UserEmail = user.Email
	form.RequestId = user.RequestId

	wHelper := tc.wfHelper
	result := wHelper.StartWorkflow(workflowOptions, taxworkflows.CreateTaxCategoryWorkflow, &form, token, user, tc.log)
	workflowClient := tc.workflowClient
	workflowRun := workflowClient.GetWorkflow(ctx, result.ID, result.RunID)
	var taxCategory taxproto.CreateTaxCategoryResponse
	err = workflowRun.Get(ctx, &taxCategory)
	if err != nil {
		tc.log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		common.RenderErrorJSON(w, "4002", err.Error(), 402, user.RequestId)
		return
	}
	common.RenderJSON(w, taxCategory)
}

// UpdateTaxCategory - Update TaxCategory
func (tc *TaxController) UpdateTaxCategory(w http.ResponseWriter, r *http.Request) {
	ctx, user, token, err := common.GetContextAuthUser(w, r, []string{"tax:cud"}, tc.ServerOpt.Auth0Audience, tc.ServerOpt.Auth0Domain, tc.UserServiceClient)
	if err != nil {
		common.RenderErrorJSON(w, "1001", err.Error(), 401, user.RequestId)
		return
	}

	id := r.PathValue("id")

	workflowOptions := client.StartWorkflowOptions{
		ID:                              "ubl_" + uuid.New(),
		TaskList:                        taxworkflows.ApplicationName,
		ExecutionStartToCloseTimeout:    time.Minute,
		DecisionTaskStartToCloseTimeout: time.Minute,
	}

	form := taxproto.UpdateTaxCategoryRequest{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&form)
	if err != nil {
		tc.log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		common.RenderErrorJSON(w, "4009", err.Error(), 402, user.RequestId)
		return
	}
	form.Id = id
	form.UserId = user.UserId
	form.UserEmail = user.Email
	form.RequestId = user.RequestId

	wHelper := tc.wfHelper
	result := wHelper.StartWorkflow(workflowOptions, taxworkflows.UpdateTaxCategoryWorkflow, &form, token, user, tc.log)
	workflowClient := tc.workflowClient
	workflowRun := workflowClient.GetWorkflow(ctx, result.ID, result.RunID)
	var response string
	err = workflowRun.Get(ctx, &response)
	if err != nil {
		tc.log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		common.RenderErrorJSON(w, "4009", err.Error(), 402, user.RequestId)
		return
	}
	common.RenderJSON(w, response)
}

// CreateTaxTotal - Create TaxTotal
func (tc *TaxController) CreateTaxTotal(w http.ResponseWriter, r *http.Request) {
	ctx, user, token, err := common.GetContextAuthUser(w, r, []string{"tax:cud"}, tc.ServerOpt.Auth0Audience, tc.ServerOpt.Auth0Domain, tc.UserServiceClient)
	if err != nil {
		common.RenderErrorJSON(w, "1001", err.Error(), 401, user.RequestId)
		return
	}

	id := r.PathValue("id")

	workflowOptions := client.StartWorkflowOptions{
		ID:                              "ubl_" + uuid.New(),
		TaskList:                        taxworkflows.ApplicationName,
		ExecutionStartToCloseTimeout:    time.Minute,
		DecisionTaskStartToCloseTimeout: time.Minute,
	}

	form := taxproto.CreateTaxTotalRequest{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&form)
	if err != nil {
		tc.log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		common.RenderErrorJSON(w, "4002", err.Error(), 402, user.RequestId)
		return
	}

	form.TaxCategoryId = id
	form.UserId = user.UserId
	form.UserEmail = user.Email
	form.RequestId = user.RequestId

	wHelper := tc.wfHelper
	result := wHelper.StartWorkflow(workflowOptions, taxworkflows.CreateTaxTotalWorkflow, &form, token, user, tc.log)

	workflowClient := tc.workflowClient
	workflowRun := workflowClient.GetWorkflow(ctx, result.ID, result.RunID)
	var taxTotal taxproto.CreateTaxTotalResponse
	err = workflowRun.Get(ctx, &taxTotal)
	if err != nil {
		tc.log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		common.RenderErrorJSON(w, "4002", err.Error(), 402, user.RequestId)
		return
	}
	common.RenderJSON(w, taxTotal)
}

// UpdateTaxTotal - Update TaxTotal
func (tc *TaxController) UpdateTaxTotal(w http.ResponseWriter, r *http.Request) {
	ctx, user, token, err := common.GetContextAuthUser(w, r, []string{"tax:cud"}, tc.ServerOpt.Auth0Audience, tc.ServerOpt.Auth0Domain, tc.UserServiceClient)
	if err != nil {
		common.RenderErrorJSON(w, "1001", err.Error(), 401, user.RequestId)
		return
	}

	id := r.PathValue("id")

	workflowOptions := client.StartWorkflowOptions{
		ID:                              "ubl_" + uuid.New(),
		TaskList:                        taxworkflows.ApplicationName,
		ExecutionStartToCloseTimeout:    time.Minute,
		DecisionTaskStartToCloseTimeout: time.Minute,
	}

	form := taxproto.UpdateTaxTotalRequest{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&form)
	if err != nil {
		tc.log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		common.RenderErrorJSON(w, "4009", err.Error(), 402, user.RequestId)
		return
	}
	form.Id = id
	form.UserId = user.UserId
	form.UserEmail = user.Email
	form.RequestId = user.RequestId

	wHelper := tc.wfHelper
	result := wHelper.StartWorkflow(workflowOptions, taxworkflows.UpdateTaxTotalWorkflow, &form, token, user, tc.log)
	workflowClient := tc.workflowClient
	workflowRun := workflowClient.GetWorkflow(ctx, result.ID, result.RunID)
	var response string
	err = workflowRun.Get(ctx, &response)
	if err != nil {
		tc.log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		common.RenderErrorJSON(w, "4009", err.Error(), 402, user.RequestId)
		return
	}
	common.RenderJSON(w, response)
}

// CreateTaxSubTotal - Create TaxSubTotal
func (tc *TaxController) CreateTaxSubTotal(w http.ResponseWriter, r *http.Request) {
	ctx, user, token, err := common.GetContextAuthUser(w, r, []string{"tax:cud"}, tc.ServerOpt.Auth0Audience, tc.ServerOpt.Auth0Domain, tc.UserServiceClient)
	if err != nil {
		common.RenderErrorJSON(w, "1001", err.Error(), 401, user.RequestId)
		return
	}

	id := r.PathValue("id")

	workflowOptions := client.StartWorkflowOptions{
		ID:                              "ubl_" + uuid.New(),
		TaskList:                        taxworkflows.ApplicationName,
		ExecutionStartToCloseTimeout:    time.Minute,
		DecisionTaskStartToCloseTimeout: time.Minute,
	}

	form := taxproto.CreateTaxSubTotalRequest{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&form)
	if err != nil {
		tc.log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		common.RenderErrorJSON(w, "4002", err.Error(), 402, user.RequestId)
		return
	}

	form.TaxTotalId = id
	form.UserId = user.UserId
	form.UserEmail = user.Email
	form.RequestId = user.RequestId

	wHelper := tc.wfHelper
	result := wHelper.StartWorkflow(workflowOptions, taxworkflows.CreateTaxSubTotalWorkflow, &form, token, user, tc.log)
	workflowClient := tc.workflowClient
	workflowRun := workflowClient.GetWorkflow(ctx, result.ID, result.RunID)
	var taxSubTotal taxproto.CreateTaxSubTotalResponse
	err = workflowRun.Get(ctx, &taxSubTotal)
	if err != nil {
		tc.log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		common.RenderErrorJSON(w, "4002", err.Error(), 402, user.RequestId)
		return
	}
	common.RenderJSON(w, taxSubTotal)
}

// UpdateTaxSubTotal - Update TaxSubTotal
func (tc *TaxController) UpdateTaxSubTotal(w http.ResponseWriter, r *http.Request) {
	ctx, user, token, err := common.GetContextAuthUser(w, r, []string{"tax:cud"}, tc.ServerOpt.Auth0Audience, tc.ServerOpt.Auth0Domain, tc.UserServiceClient)
	if err != nil {
		common.RenderErrorJSON(w, "1001", err.Error(), 401, user.RequestId)
		return
	}

	id := r.PathValue("id")

	workflowOptions := client.StartWorkflowOptions{
		ID:                              "ubl_" + uuid.New(),
		TaskList:                        taxworkflows.ApplicationName,
		ExecutionStartToCloseTimeout:    time.Minute,
		DecisionTaskStartToCloseTimeout: time.Minute,
	}

	form := taxproto.UpdateTaxSubTotalRequest{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&form)
	if err != nil {
		tc.log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		common.RenderErrorJSON(w, "4009", err.Error(), 402, user.RequestId)
		return
	}
	form.Id = id
	form.UserId = user.UserId
	form.UserEmail = user.Email
	form.RequestId = user.RequestId

	wHelper := tc.wfHelper
	result := wHelper.StartWorkflow(workflowOptions, taxworkflows.UpdateTaxSubTotalWorkflow, &form, token, user, tc.log)
	workflowClient := tc.workflowClient
	workflowRun := workflowClient.GetWorkflow(ctx, result.ID, result.RunID)
	var response string
	err = workflowRun.Get(ctx, &response)
	if err != nil {
		tc.log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		common.RenderErrorJSON(w, "4009", err.Error(), 402, user.RequestId)
		return
	}
	common.RenderJSON(w, response)
}
