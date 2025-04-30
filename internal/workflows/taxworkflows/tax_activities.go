package taxworkflows

import (
	"context"

	partyproto "github.com/cloudfresco/sc-ubl/internal/protogen/party/v1"
	taxproto "github.com/cloudfresco/sc-ubl/internal/protogen/tax/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
)

type TaxActivities struct {
	TaxServiceClient taxproto.TaxServiceClient
}

// CreateTaxCategoryActivity - Create TaxCategory activity
func (ta *TaxActivities) CreateTaxCategoryActivity(ctx context.Context, form *taxproto.CreateTaxCategoryRequest, tokenString string, user *partyproto.GetAuthUserDetailsResponse, log *zap.Logger) (*taxproto.CreateTaxCategoryResponse, error) {
	taxServiceClient := ta.TaxServiceClient
	md := metadata.Pairs("authorization", "Bearer "+tokenString)
	ctxNew := metadata.NewOutgoingContext(ctx, md)
	taxCategory, err := taxServiceClient.CreateTaxCategory(ctxNew, form)
	if err != nil {
		log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		return nil, err
	}
	return taxCategory, nil
}

// UpdateTaxCategoryActivity - update TaxCategory activity
func (ta *TaxActivities) UpdateTaxCategoryActivity(ctx context.Context, form *taxproto.UpdateTaxCategoryRequest, tokenString string, user *partyproto.GetAuthUserDetailsResponse, log *zap.Logger) (string, error) {
	taxServiceClient := ta.TaxServiceClient
	md := metadata.Pairs("authorization", "Bearer "+tokenString)
	ctxNew := metadata.NewOutgoingContext(ctx, md)
	_, err := taxServiceClient.UpdateTaxCategory(ctxNew, form)
	if err != nil {
		log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		return "", err
	}
	return "Updated Successfully", nil
}

// CreateTaxSchemeActivity - Create TaxScheme activity
func (ta *TaxActivities) CreateTaxSchemeActivity(ctx context.Context, form *taxproto.CreateTaxSchemeRequest, tokenString string, user *partyproto.GetAuthUserDetailsResponse, log *zap.Logger) (*taxproto.CreateTaxSchemeResponse, error) {
	taxServiceClient := ta.TaxServiceClient
	md := metadata.Pairs("authorization", "Bearer "+tokenString)
	ctxNew := metadata.NewOutgoingContext(ctx, md)
	taxScheme, err := taxServiceClient.CreateTaxScheme(ctxNew, form)
	if err != nil {
		log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		return nil, err
	}
	return taxScheme, nil
}

// UpdateTaxSchemeActivity - update TaxScheme activity
func (ta *TaxActivities) UpdateTaxSchemeActivity(ctx context.Context, form *taxproto.UpdateTaxSchemeRequest, tokenString string, user *partyproto.GetAuthUserDetailsResponse, log *zap.Logger) (string, error) {
	taxServiceClient := ta.TaxServiceClient
	md := metadata.Pairs("authorization", "Bearer "+tokenString)
	ctxNew := metadata.NewOutgoingContext(ctx, md)
	_, err := taxServiceClient.UpdateTaxScheme(ctxNew, form)
	if err != nil {
		log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		return "", err
	}
	return "Updated Successfully", nil
}

// CreateTaxSchemeJurisdictionActivity - Create TaxSchemeJurisdiction activity
func (ta *TaxActivities) CreateTaxSchemeJurisdictionActivity(ctx context.Context, form *taxproto.CreateTaxSchemeJurisdictionRequest, tokenString string, user *partyproto.GetAuthUserDetailsResponse, log *zap.Logger) (*taxproto.CreateTaxSchemeJurisdictionResponse, error) {
	taxServiceClient := ta.TaxServiceClient
	md := metadata.Pairs("authorization", "Bearer "+tokenString)
	ctxNew := metadata.NewOutgoingContext(ctx, md)
	taxSchemeJurisdiction, err := taxServiceClient.CreateTaxSchemeJurisdiction(ctxNew, form)
	if err != nil {
		log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		return nil, err
	}
	return taxSchemeJurisdiction, nil
}

// UpdateTaxSchemeJurisdictionActivity - update TaxSchemeJurisdiction activity
func (ta *TaxActivities) UpdateTaxSchemeJurisdictionActivity(ctx context.Context, form *taxproto.UpdateTaxSchemeJurisdictionRequest, tokenString string, user *partyproto.GetAuthUserDetailsResponse, log *zap.Logger) (string, error) {
	taxServiceClient := ta.TaxServiceClient
	md := metadata.Pairs("authorization", "Bearer "+tokenString)
	ctxNew := metadata.NewOutgoingContext(ctx, md)
	_, err := taxServiceClient.UpdateTaxSchemeJurisdiction(ctxNew, form)
	if err != nil {
		log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		return "", err
	}
	return "Updated Successfully", nil
}

// CreateTaxTotalActivity - Create TaxTotal activity
func (ta *TaxActivities) CreateTaxTotalActivity(ctx context.Context, form *taxproto.CreateTaxTotalRequest, tokenString string, user *partyproto.GetAuthUserDetailsResponse, log *zap.Logger) (*taxproto.CreateTaxTotalResponse, error) {
	taxServiceClient := ta.TaxServiceClient
	md := metadata.Pairs("authorization", "Bearer "+tokenString)
	ctxNew := metadata.NewOutgoingContext(ctx, md)
	taxTotal, err := taxServiceClient.CreateTaxTotal(ctxNew, form)
	if err != nil {
		log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		return nil, err
	}
	return taxTotal, nil
}

// UpdateTaxTotalActivity - update TaxTotal activity
func (ta *TaxActivities) UpdateTaxTotalActivity(ctx context.Context, form *taxproto.UpdateTaxTotalRequest, tokenString string, user *partyproto.GetAuthUserDetailsResponse, log *zap.Logger) (string, error) {
	taxServiceClient := ta.TaxServiceClient
	md := metadata.Pairs("authorization", "Bearer "+tokenString)
	ctxNew := metadata.NewOutgoingContext(ctx, md)
	_, err := taxServiceClient.UpdateTaxTotal(ctxNew, form)
	if err != nil {
		log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		return "", err
	}
	return "Updated Successfully", nil
}

// CreateTaxSubTotalActivity - Create TaxSubTotal activity
func (ta *TaxActivities) CreateTaxSubTotalActivity(ctx context.Context, form *taxproto.CreateTaxSubTotalRequest, tokenString string, user *partyproto.GetAuthUserDetailsResponse, log *zap.Logger) (*taxproto.CreateTaxSubTotalResponse, error) {
	taxServiceClient := ta.TaxServiceClient
	md := metadata.Pairs("authorization", "Bearer "+tokenString)
	ctxNew := metadata.NewOutgoingContext(ctx, md)
	taxSubTotal, err := taxServiceClient.CreateTaxSubTotal(ctxNew, form)
	if err != nil {
		log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		return nil, err
	}
	return taxSubTotal, nil
}

// UpdateTaxSubTotalActivity - update TaxSubTotal activity
func (ta *TaxActivities) UpdateTaxSubTotalActivity(ctx context.Context, form *taxproto.UpdateTaxSubTotalRequest, tokenString string, user *partyproto.GetAuthUserDetailsResponse, log *zap.Logger) (string, error) {
	taxServiceClient := ta.TaxServiceClient
	md := metadata.Pairs("authorization", "Bearer "+tokenString)
	ctxNew := metadata.NewOutgoingContext(ctx, md)
	_, err := taxServiceClient.UpdateTaxSubTotal(ctxNew, form)
	if err != nil {
		log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		return "", err
	}
	return "Updated Successfully", nil
}
