package orderworkers

import (
	"os"

	"github.com/cloudfresco/sc-ubl/internal/config"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/cadence/worker"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/cloudfresco/sc-ubl/internal/common"
	orderproto "github.com/cloudfresco/sc-ubl/internal/protogen/order/v1"

	orderworkflows "github.com/cloudfresco/sc-ubl/internal/workflows/orderworkflows"
)

// This needs to be done as part of a bootstrap step when the process starts.
// The workers are supposed to be long running.
func startWorkers(h *common.WfHelper) {
	// Configure worker options.
	workerOptions := worker.Options{
		MetricsScope: h.WorkerMetricScope,
		Logger:       h.Logger,
	}
	h.StartWorkers(h.Config.DomainName, orderworkflows.ApplicationName, workerOptions)
}

func StartOrderWorker(log *zap.Logger, isTest bool, pwd string, grpcServerOpt *config.GrpcServerOptions, configFilePath string) {
	var h common.WfHelper
	h.SetupServiceConfig(configFilePath)

	creds, err := common.GetClientCred(log, isTest, pwd, grpcServerOpt)
	if err != nil {
		log.Error("Error", zap.Error(err))
		os.Exit(1)
	}

	orderconn, err := grpc.NewClient(grpcServerOpt.GrpcOrderServerPort, grpc.WithTransportCredentials(creds), grpc.WithStatsHandler(otelgrpc.NewClientHandler()))
	if err != nil {
		log.Error("Error",
			zap.Error(err))
		os.Exit(1)
	}
	purchaseOrderHeaderServiceClient := orderproto.NewPurchaseOrderHeaderServiceClient(orderconn)
	purchaseOrderHeaderActivities := &orderworkflows.PurchaseOrderHeaderActivities{PurchaseOrderHeaderServiceClient: purchaseOrderHeaderServiceClient}

	h.RegisterWorkflow(orderworkflows.CreatePurchaseOrderHeaderWorkflow)
	h.RegisterWorkflow(orderworkflows.UpdatePurchaseOrderHeaderWorkflow)
	h.RegisterActivity(purchaseOrderHeaderActivities)

	startWorkers(&h)

	// The workers are supposed to be long running process that should not exit.
	// Use select{} to block indefinitely for samples, you can quit by CMD+C.
	select {}
}
