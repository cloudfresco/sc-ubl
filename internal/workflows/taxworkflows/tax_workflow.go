package taxworkflows

import (
	"time"

	partyproto "github.com/cloudfresco/sc-ubl/internal/protogen/party/v1"
	taxproto "github.com/cloudfresco/sc-ubl/internal/protogen/tax/v1"

	"go.uber.org/cadence/workflow"
	"go.uber.org/zap"
)

const (
	// ApplicationName is the task list
	ApplicationName = "ubl"
)

// CreateTaxCategoryWorkflow - Create TaxCategory workflow
func CreateTaxCategoryWorkflow(ctx workflow.Context, form *taxproto.CreateTaxCategoryRequest, tokenString string, user *partyproto.GetAuthUserDetailsResponse, log *zap.Logger) (*taxproto.CreateTaxCategoryResponse, error) {
	ao := workflow.ActivityOptions{
		ScheduleToStartTimeout: time.Minute,
		StartToCloseTimeout:    time.Minute,
		HeartbeatTimeout:       time.Second * 20,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	logger := workflow.GetLogger(ctx)
	var ta *TaxActivities
	var taxCategory taxproto.CreateTaxCategoryResponse
	err := workflow.ExecuteActivity(ctx, ta.CreateTaxCategoryActivity, form, tokenString, user, log).Get(ctx, &taxCategory)
	if err != nil {
		logger.Error("Failed to CreateTaxCategoryWorkflow", zap.Error(err))
		log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		return nil, err
	}
	return &taxCategory, nil
}

// UpdateTaxCategoryWorkflow - update TaxCategory workflow
func UpdateTaxCategoryWorkflow(ctx workflow.Context, form *taxproto.UpdateTaxCategoryRequest, tokenString string, user *partyproto.GetAuthUserDetailsResponse, log *zap.Logger) (string, error) {
	ao := workflow.ActivityOptions{
		ScheduleToStartTimeout: time.Minute,
		StartToCloseTimeout:    time.Minute,
		HeartbeatTimeout:       time.Second * 20,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	logger := workflow.GetLogger(ctx)
	var ta *TaxActivities
	var resp string
	err := workflow.ExecuteActivity(ctx, ta.UpdateTaxCategoryActivity, form, tokenString, user, log).Get(ctx, &resp)
	if err != nil {
		logger.Error("Failed to UpdateTaxCategoryWorkflow", zap.Error(err))
		log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		return "", err
	}
	return resp, nil
}

// CreateTaxSchemeWorkflow - Create TaxScheme workflow
func CreateTaxSchemeWorkflow(ctx workflow.Context, form *taxproto.CreateTaxSchemeRequest, tokenString string, user *partyproto.GetAuthUserDetailsResponse, log *zap.Logger) (*taxproto.CreateTaxSchemeResponse, error) {
	ao := workflow.ActivityOptions{
		ScheduleToStartTimeout: time.Minute,
		StartToCloseTimeout:    time.Minute,
		HeartbeatTimeout:       time.Second * 20,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	logger := workflow.GetLogger(ctx)
	var ta *TaxActivities
	var taxScheme taxproto.CreateTaxSchemeResponse
	err := workflow.ExecuteActivity(ctx, ta.CreateTaxSchemeActivity, form, tokenString, user, log).Get(ctx, &taxScheme)
	if err != nil {
		logger.Error("Failed to CreateTaxSchemeWorkflow", zap.Error(err))
		log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		return nil, err
	}
	return &taxScheme, nil
}

// UpdateTaxSchemeWorkflow - update TaxScheme workflow
func UpdateTaxSchemeWorkflow(ctx workflow.Context, form *taxproto.UpdateTaxSchemeRequest, tokenString string, user *partyproto.GetAuthUserDetailsResponse, log *zap.Logger) (string, error) {
	ao := workflow.ActivityOptions{
		ScheduleToStartTimeout: time.Minute,
		StartToCloseTimeout:    time.Minute,
		HeartbeatTimeout:       time.Second * 20,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	logger := workflow.GetLogger(ctx)
	var ta *TaxActivities
	var resp string
	err := workflow.ExecuteActivity(ctx, ta.UpdateTaxSchemeActivity, form, tokenString, user, log).Get(ctx, &resp)
	if err != nil {
		logger.Error("Failed to UpdateTaxSchemeWorkflow", zap.Error(err))
		log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		return "", err
	}
	return resp, nil
}

// CreateTaxSchemeJurisdictionWorkflow - Create TaxSchemeJurisdiction workflow
func CreateTaxSchemeJurisdictionWorkflow(ctx workflow.Context, form *taxproto.CreateTaxSchemeJurisdictionRequest, tokenString string, user *partyproto.GetAuthUserDetailsResponse, log *zap.Logger) (*taxproto.CreateTaxSchemeJurisdictionResponse, error) {
	ao := workflow.ActivityOptions{
		ScheduleToStartTimeout: time.Minute,
		StartToCloseTimeout:    time.Minute,
		HeartbeatTimeout:       time.Second * 20,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	logger := workflow.GetLogger(ctx)
	var ta *TaxActivities
	var taxSchemeJurisdiction taxproto.CreateTaxSchemeJurisdictionResponse
	err := workflow.ExecuteActivity(ctx, ta.CreateTaxSchemeJurisdictionActivity, form, tokenString, user, log).Get(ctx, &taxSchemeJurisdiction)
	if err != nil {
		logger.Error("Failed to CreateTaxSchemeJurisdictionWorkflow", zap.Error(err))
		log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		return nil, err
	}
	return &taxSchemeJurisdiction, nil
}

// UpdateTaxSchemeJurisdictionWorkflow - update TaxSchemeJurisdiction workflow
func UpdateTaxSchemeJurisdictionWorkflow(ctx workflow.Context, form *taxproto.UpdateTaxSchemeJurisdictionRequest, tokenString string, user *partyproto.GetAuthUserDetailsResponse, log *zap.Logger) (string, error) {
	ao := workflow.ActivityOptions{
		ScheduleToStartTimeout: time.Minute,
		StartToCloseTimeout:    time.Minute,
		HeartbeatTimeout:       time.Second * 20,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	logger := workflow.GetLogger(ctx)
	var ta *TaxActivities
	var resp string
	err := workflow.ExecuteActivity(ctx, ta.UpdateTaxSchemeJurisdictionActivity, form, tokenString, user, log).Get(ctx, &resp)
	if err != nil {
		logger.Error("Failed to UpdateTaxSchemeJurisdictionWorkflow", zap.Error(err))
		log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		return "", err
	}
	return resp, nil
}

// CreateTaxTotalWorkflow - Create TaxTotal workflow
func CreateTaxTotalWorkflow(ctx workflow.Context, form *taxproto.CreateTaxTotalRequest, tokenString string, user *partyproto.GetAuthUserDetailsResponse, log *zap.Logger) (*taxproto.CreateTaxTotalResponse, error) {
	ao := workflow.ActivityOptions{
		ScheduleToStartTimeout: time.Minute,
		StartToCloseTimeout:    time.Minute,
		HeartbeatTimeout:       time.Second * 20,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	logger := workflow.GetLogger(ctx)
	var ta *TaxActivities
	var taxTotal taxproto.CreateTaxTotalResponse
	err := workflow.ExecuteActivity(ctx, ta.CreateTaxTotalActivity, form, tokenString, user, log).Get(ctx, &taxTotal)
	if err != nil {
		logger.Error("Failed to CreateTaxTotalWorkflow", zap.Error(err))
		log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		return nil, err
	}
	return &taxTotal, nil
}

// UpdateTaxTotalWorkflow - update TaxTotal workflow
func UpdateTaxTotalWorkflow(ctx workflow.Context, form *taxproto.UpdateTaxTotalRequest, tokenString string, user *partyproto.GetAuthUserDetailsResponse, log *zap.Logger) (string, error) {
	ao := workflow.ActivityOptions{
		ScheduleToStartTimeout: time.Minute,
		StartToCloseTimeout:    time.Minute,
		HeartbeatTimeout:       time.Second * 20,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	logger := workflow.GetLogger(ctx)
	var ta *TaxActivities
	var resp string
	err := workflow.ExecuteActivity(ctx, ta.UpdateTaxTotalActivity, form, tokenString, user, log).Get(ctx, &resp)
	if err != nil {
		logger.Error("Failed to UpdateTaxTotalWorkflow", zap.Error(err))
		log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		return "", err
	}
	return resp, nil
}

// CreateTaxSubTotalWorkflow - Create TaxSubTotal workflow
func CreateTaxSubTotalWorkflow(ctx workflow.Context, form *taxproto.CreateTaxSubTotalRequest, tokenString string, user *partyproto.GetAuthUserDetailsResponse, log *zap.Logger) (*taxproto.CreateTaxSubTotalResponse, error) {
	ao := workflow.ActivityOptions{
		ScheduleToStartTimeout: time.Minute,
		StartToCloseTimeout:    time.Minute,
		HeartbeatTimeout:       time.Second * 20,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	logger := workflow.GetLogger(ctx)
	var ta *TaxActivities
	var taxSubTotal taxproto.CreateTaxSubTotalResponse
	err := workflow.ExecuteActivity(ctx, ta.CreateTaxSubTotalActivity, form, tokenString, user, log).Get(ctx, &taxSubTotal)
	if err != nil {
		logger.Error("Failed to CreateTaxSubTotalWorkflow", zap.Error(err))
		log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		return nil, err
	}
	return &taxSubTotal, nil
}

// UpdateTaxSubTotalWorkflow - update TaxSubTotal workflow
func UpdateTaxSubTotalWorkflow(ctx workflow.Context, form *taxproto.UpdateTaxSubTotalRequest, tokenString string, user *partyproto.GetAuthUserDetailsResponse, log *zap.Logger) (string, error) {
	ao := workflow.ActivityOptions{
		ScheduleToStartTimeout: time.Minute,
		StartToCloseTimeout:    time.Minute,
		HeartbeatTimeout:       time.Second * 20,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	logger := workflow.GetLogger(ctx)
	var ta *TaxActivities
	var resp string
	err := workflow.ExecuteActivity(ctx, ta.UpdateTaxSubTotalActivity, form, tokenString, user, log).Get(ctx, &resp)
	if err != nil {
		logger.Error("Failed to UpdateTaxSubTotalWorkflow", zap.Error(err))
		log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		return "", err
	}
	return resp, nil
}
