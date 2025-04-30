package logisticscontrollers

import (
	"context"
	"net/http"
	"os"
	"path/filepath"

	"github.com/cloudfresco/sc-ubl/internal/common"
	"github.com/cloudfresco/sc-ubl/internal/config"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"

	logisticsproto "github.com/cloudfresco/sc-ubl/internal/protogen/logistics/v1"
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

// Init the invoice controllers
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

// InitTest the invoice controllers
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

	logisticsconn, err := grpc.NewClient(grpcServerOpt.GrpcLogisticsServerPort, grpc.WithTransportCredentials(creds), grpc.WithStatsHandler(otelgrpc.NewClientHandler()))
	if err != nil {
		log.Error("Error", zap.Error(err))
		return err
	}

	u := partyproto.NewUserServiceClient(userconn)
	c := logisticsproto.NewConsignmentServiceClient(logisticsconn)
	r := logisticsproto.NewReceiptAdviceHeaderServiceClient(logisticsconn)
	d := logisticsproto.NewDespatchServiceClient(logisticsconn)
	s := logisticsproto.NewShipmentServiceClient(logisticsconn)

	initConsignments(mux, serverOpt, log, u, c, h, workflowClient)
	initReceiptAdviceHeaders(mux, serverOpt, log, u, r, h, workflowClient)
	initDespatches(mux, serverOpt, log, u, d, h, workflowClient)
	initShipments(mux, serverOpt, log, u, s, h, workflowClient)

	return nil
}

func initConsignments(mux *http.ServeMux, serverOpt *config.ServerOptions, log *zap.Logger, u partyproto.UserServiceClient, c logisticsproto.ConsignmentServiceClient, wfHelper common.WfHelper, workflowClient client.Client) {
	cc := NewConsignmentController(log, u, c, h, workflowClient, serverOpt)

	mux.Handle("GET /v2.3/consignments", http.HandlerFunc(cc.Index))
	mux.Handle("GET /v2.3/consignments/{id}", http.HandlerFunc(cc.Show))

	mux.Handle("POST /v2.3/consignments", http.HandlerFunc(cc.CreateConsignment))
}

func initReceiptAdviceHeaders(mux *http.ServeMux, serverOpt *config.ServerOptions, log *zap.Logger, u partyproto.UserServiceClient, r logisticsproto.ReceiptAdviceHeaderServiceClient, wfHelper common.WfHelper, workflowClient client.Client) {
	rc := NewReceiptAdviceHeaderController(log, u, r, h, workflowClient, serverOpt)

	mux.Handle("GET /v2.3/receipt-advices", http.HandlerFunc(rc.Index))
	mux.Handle("GET /v2.3/receipt-advices/{id}", http.HandlerFunc(rc.Show))
	mux.Handle("GET /v2.3/receipt-advices/{id}/lines", http.HandlerFunc(rc.GetReceiptAdviceLines))

	mux.Handle("POST /v2.3/receipt-advices", http.HandlerFunc(rc.CreateReceiptAdviceHeader))

	mux.Handle("PUT /v2.3/receipt-advices/{id}", http.HandlerFunc(rc.UpdateReceiptAdviceHeader))
}

func initDespatches(mux *http.ServeMux, serverOpt *config.ServerOptions, log *zap.Logger, u partyproto.UserServiceClient, d logisticsproto.DespatchServiceClient, wfHelper common.WfHelper, workflowClient client.Client) {
	dc := NewDespatchHeaderController(log, u, d, h, workflowClient, serverOpt)

	mux.Handle("GET /v2.3/despatches", http.HandlerFunc(dc.Index))
	mux.Handle("GET /v2.3/despatches/{id}", http.HandlerFunc(dc.Show))
	mux.Handle("GET /v2.3/despatches/{id}/lines", http.HandlerFunc(dc.GetDespatchLines))

	mux.Handle("POST /v2.3/despatches", http.HandlerFunc(dc.CreateDespatchHeader))

	mux.Handle("PUT /v2.3/despatches/{id}", http.HandlerFunc(dc.UpdateDespatchHeader))
}

func initShipments(mux *http.ServeMux, serverOpt *config.ServerOptions, log *zap.Logger, u partyproto.UserServiceClient, s logisticsproto.ShipmentServiceClient, wfHelper common.WfHelper, workflowClient client.Client) {
	sc := NewShipmentController(log, u, s, h, workflowClient, serverOpt)

	mux.Handle("POST /v2.3/shipments", http.HandlerFunc(sc.CreateShipment))
}
