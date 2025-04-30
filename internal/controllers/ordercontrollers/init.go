package ordercontrollers

import (
	"context"
	"net/http"
	"os"
	"path/filepath"

	"github.com/cloudfresco/sc-ubl/internal/common"
	"github.com/cloudfresco/sc-ubl/internal/config"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"

	orderproto "github.com/cloudfresco/sc-ubl/internal/protogen/order/v1"
	partyproto "github.com/cloudfresco/sc-ubl/internal/protogen/party/v1"

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

// Init the order controllers
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

// InitTest the order controllers
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

	orderconn, err := grpc.NewClient(grpcServerOpt.GrpcOrderServerPort, grpc.WithTransportCredentials(creds), grpc.WithStatsHandler(otelgrpc.NewClientHandler()))
	if err != nil {
		log.Error("Error", zap.Error(err))
		return err
	}

	u := partyproto.NewUserServiceClient(userconn)
	p := orderproto.NewPurchaseOrderHeaderServiceClient(orderconn)

	initPurchaseOrders(mux, serverOpt, log, u, p, h, workflowClient)

	return nil
}

func initPurchaseOrders(mux *http.ServeMux, serverOpt *config.ServerOptions, log *zap.Logger, u partyproto.UserServiceClient, p orderproto.PurchaseOrderHeaderServiceClient, wfHelper common.WfHelper, workflowClient client.Client) {
	po := NewPurchaseOrderHeaderController(log, u, p, h, workflowClient, serverOpt)

	mux.Handle("GET /v2.3/purchase-orders", http.HandlerFunc(po.Index))
	mux.Handle("GET /v2.3/purchase-orders/{id}", http.HandlerFunc(po.Show))
	mux.Handle("GET /v2.3/parties/{id}/lines", http.HandlerFunc(po.GetPurchaseOrderLines))

	mux.Handle("POST /v2.3/purchase-orders", http.HandlerFunc(po.CreatePurchaseOrderHeader))

	mux.Handle("PUT /v2.3/purchase-orders/{id}", http.HandlerFunc(po.UpdatePurchaseOrderHeader))
}
