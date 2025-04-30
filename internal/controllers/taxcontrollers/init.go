package taxcontrollers

import (
	"context"
	"net/http"
	"os"
	"path/filepath"

	"github.com/cloudfresco/sc-ubl/internal/common"
	"github.com/cloudfresco/sc-ubl/internal/config"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"

	partyproto "github.com/cloudfresco/sc-ubl/internal/protogen/party/v1"
	taxproto "github.com/cloudfresco/sc-ubl/internal/protogen/tax/v1"

	"github.com/throttled/throttled/v2/store/goredisstore"
	"go.uber.org/cadence/client"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	h              common.WfHelper
	workflowClient client.Client
)

// Init the tax controllers
func Init(log *zap.Logger, mux *http.ServeMux, store *goredisstore.GoRedisStore, serverOpt *config.ServerOptions, grpcServerOpt *config.GrpcServerOptions, uptraceOpt *config.UptraceOptions, configFilePath string) error {
	pwd, _ := os.Getwd()
	keyPath := pwd + filepath.FromSlash(grpcServerOpt.GrpcCaCertPath)

	err := initSetup(log, mux, keyPath, configFilePath, serverOpt, grpcServerOpt)
	if err != nil {
		log.Error("Error", zap.Int("msgnum", 110), zap.Error(err))
		return err
	}

	return nil
}

// InitTest the tax controllers
func InitTest(log *zap.Logger, mux *http.ServeMux, store *goredisstore.GoRedisStore, serverOpt *config.ServerOptions, grpcServerOpt *config.GrpcServerOptions, uptraceOpt *config.UptraceOptions, configFilePath string) error {
	pwd, _ := os.Getwd()
	keyPath := filepath.Join(pwd, filepath.FromSlash("/../../../")+filepath.FromSlash(grpcServerOpt.GrpcCaCertPath))

	err := initSetup(log, mux, keyPath, configFilePath, serverOpt, grpcServerOpt)
	if err != nil {
		log.Error("Error", zap.Int("msgnum", 110), zap.Error(err))
		return err
	}

	return nil
}

func initSetup(log *zap.Logger, mux *http.ServeMux, keyPath string, configFilePath string, serverOpt *config.ServerOptions, grpcServerOpt *config.GrpcServerOptions) error {
	creds, err := credentials.NewClientTLSFromFile(keyPath, "localhost")
	if err != nil {
		log.Error("Error", zap.Int("msgnum", 110), zap.Error(err))
	}

	tp, err := config.InitTracerProvider()
	if err != nil {
		log.Error("Error", zap.Int("msgnum", 9108), zap.Error(err))
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Error("Error", zap.Int("msgnum", 9108), zap.Error(err))
		}
	}()

	h.SetupServiceConfig(configFilePath)
	workflowClient, err = h.Builder.BuildCadenceClient()
	if err != nil {
		panic(err)
	}

	userconn, err := grpc.NewClient(grpcServerOpt.GrpcUserServerPort, grpc.WithTransportCredentials(creds), grpc.WithStatsHandler(otelgrpc.NewClientHandler()))
	if err != nil {
		log.Error("Error", zap.Int("msgnum", 113), zap.Error(err))
		return err
	}

	taxconn, err := grpc.NewClient(grpcServerOpt.GrpcTaxServerPort, grpc.WithTransportCredentials(creds), grpc.WithStatsHandler(otelgrpc.NewClientHandler()))
	if err != nil {
		log.Error("Error", zap.Error(err))
		return err
	}

	u := partyproto.NewUserServiceClient(userconn)
	t := taxproto.NewTaxServiceClient(taxconn)

	initPurchaseOrders(mux, serverOpt, log, u, t, h, workflowClient)

	return nil
}

func initPurchaseOrders(mux *http.ServeMux, serverOpt *config.ServerOptions, log *zap.Logger, u partyproto.UserServiceClient, t taxproto.TaxServiceClient, wfHelper common.WfHelper, workflowClient client.Client) {
	tc := NewTaxController(log, u, t, h, workflowClient, serverOpt)

	mux.Handle("GET /v2.3/tax-schemes", http.HandlerFunc(tc.Index))
	mux.Handle("GET /v2.3/tax-schemes/{id}", http.HandlerFunc(tc.Show))

	mux.Handle("POST /v2.3/tax-schemes", http.HandlerFunc(tc.CreateTaxScheme))
	mux.Handle("POST /v2.3/tax-schemes/{id}/add-tax-scheme-jurisdiction", http.HandlerFunc(tc.CreateTaxSchemeJurisdiction))
	mux.Handle("POST /v2.3/tax-schemes/{id}/add-tax-category", http.HandlerFunc(tc.CreateTaxCategory))
	mux.Handle("POST /v2.3/tax-schemes/tax-categories/{id}/add-tax-total", http.HandlerFunc(tc.CreateTaxTotal))
	mux.Handle("POST /v2.3/tax-schemes/tax-totals/{id}/add-tax-subtotal", http.HandlerFunc(tc.CreateTaxSubTotal))

	mux.Handle("PUT /v2.3/tax-schemes/{id}", http.HandlerFunc(tc.UpdateTaxScheme))
	mux.Handle("PUT /v2.3/tax-schemes/tax-scheme-jurisdictions/{id}", http.HandlerFunc(tc.UpdateTaxSchemeJurisdiction))
	mux.Handle("PUT /v2.3/tax-schemes/tax-categories/{id}", http.HandlerFunc(tc.UpdateTaxCategory))
	mux.Handle("PUT /v2.3/tax-schemes/tax-totals/{id}", http.HandlerFunc(tc.UpdateTaxTotal))
	mux.Handle("PUT /v2.3/tax-schemes/tax-subtotals/{id}", http.HandlerFunc(tc.UpdateTaxSubTotal))
}
