package taxworkers

import (
	"os"

	"github.com/cloudfresco/sc-ubl/internal/config"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/cadence/worker"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/cloudfresco/sc-ubl/internal/common"
	taxproto "github.com/cloudfresco/sc-ubl/internal/protogen/tax/v1"

	taxworkflows "github.com/cloudfresco/sc-ubl/internal/workflows/taxworkflows"
)

// This needs to be done as part of a bootstrap step when the process starts.
// The workers are supposed to be long running.
func startWorkers(h *common.WfHelper) {
	// Configure worker options.
	workerOptions := worker.Options{
		MetricsScope: h.WorkerMetricScope,
		Logger:       h.Logger,
	}
	h.StartWorkers(h.Config.DomainName, taxworkflows.ApplicationName, workerOptions)
}

func StartTaxWorker(log *zap.Logger, isTest bool, pwd string, grpcServerOpt *config.GrpcServerOptions, configFilePath string) {
	var h common.WfHelper
	h.SetupServiceConfig(configFilePath)

	creds, err := common.GetClientCred(log, isTest, pwd, grpcServerOpt)
	if err != nil {
		log.Error("Error", zap.Error(err))
		os.Exit(1)
	}

	taxconn, err := grpc.NewClient(grpcServerOpt.GrpcTaxServerPort, grpc.WithTransportCredentials(creds), grpc.WithStatsHandler(otelgrpc.NewClientHandler()))
	if err != nil {
		log.Error("Error",
			zap.Error(err))
		os.Exit(1)
	}
	taxServiceClient := taxproto.NewTaxServiceClient(taxconn)
	taxActivities := &taxworkflows.TaxActivities{TaxServiceClient: taxServiceClient}

	h.RegisterWorkflow(taxworkflows.CreateTaxSchemeWorkflow)
	h.RegisterWorkflow(taxworkflows.UpdateTaxSchemeWorkflow)
	h.RegisterWorkflow(taxworkflows.CreateTaxCategoryWorkflow)
	h.RegisterWorkflow(taxworkflows.UpdateTaxCategoryWorkflow)
	h.RegisterWorkflow(taxworkflows.CreateTaxSchemeJurisdictionWorkflow)
	h.RegisterWorkflow(taxworkflows.UpdateTaxSchemeJurisdictionWorkflow)
	h.RegisterWorkflow(taxworkflows.CreateTaxTotalWorkflow)
	h.RegisterWorkflow(taxworkflows.UpdateTaxTotalWorkflow)
	h.RegisterWorkflow(taxworkflows.CreateTaxSubTotalWorkflow)
	h.RegisterWorkflow(taxworkflows.UpdateTaxSubTotalWorkflow)
	h.RegisterActivity(taxActivities)

	startWorkers(&h)

	// The workers are supposed to be long running process that should not exit.
	// Use select{} to block indefinitely for samples, you can quit by CMD+C.
	select {}
}
