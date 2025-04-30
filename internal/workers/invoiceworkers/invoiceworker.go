package invoiceworkers

import (
	"os"

	"github.com/cloudfresco/sc-ubl/internal/config"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/cadence/worker"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/cloudfresco/sc-ubl/internal/common"
	invoiceproto "github.com/cloudfresco/sc-ubl/internal/protogen/invoice/v1"

	invoiceworkflows "github.com/cloudfresco/sc-ubl/internal/workflows/invoiceworkflows"
)

// This needs to be done as part of a bootstrap step when the process starts.
// The workers are supposed to be long running.
func startWorkers(h *common.WfHelper) {
	// Configure worker options.
	workerOptions := worker.Options{
		MetricsScope: h.WorkerMetricScope,
		Logger:       h.Logger,
	}
	h.StartWorkers(h.Config.DomainName, invoiceworkflows.ApplicationName, workerOptions)
}

func StartInvoiceWorker(log *zap.Logger, isTest bool, pwd string, grpcServerOpt *config.GrpcServerOptions, configFilePath string) {
	var h common.WfHelper
	h.SetupServiceConfig(configFilePath)

	creds, err := common.GetClientCred(log, isTest, pwd, grpcServerOpt)
	if err != nil {
		log.Error("Error", zap.Error(err))
		os.Exit(1)
	}

	invoiceconn, err := grpc.NewClient(grpcServerOpt.GrpcInvoiceServerPort, grpc.WithTransportCredentials(creds), grpc.WithStatsHandler(otelgrpc.NewClientHandler()))
	if err != nil {
		log.Error("Error",
			zap.Error(err))
		os.Exit(1)
	}
	invoiceServiceClient := invoiceproto.NewInvoiceServiceClient(invoiceconn)
	invoiceActivities := &invoiceworkflows.InvoiceActivities{InvoiceServiceClient: invoiceServiceClient}

	creditNoteHeaderServiceClient := invoiceproto.NewCreditNoteHeaderServiceClient(invoiceconn)
	creditNoteHeaderActivities := &invoiceworkflows.CreditNoteHeaderActivities{CreditNoteHeaderServiceClient: creditNoteHeaderServiceClient}

	debitNoteHeaderServiceClient := invoiceproto.NewDebitNoteHeaderServiceClient(invoiceconn)
	debitNoteHeaderActivities := &invoiceworkflows.DebitNoteHeaderActivities{DebitNoteHeaderServiceClient: debitNoteHeaderServiceClient}

	h.RegisterWorkflow(invoiceworkflows.CreateCreditNoteHeaderWorkflow)
	h.RegisterWorkflow(invoiceworkflows.UpdateCreditNoteHeaderWorkflow)
	h.RegisterWorkflow(invoiceworkflows.CreateDebitNoteHeaderWorkflow)
	h.RegisterWorkflow(invoiceworkflows.UpdateDebitNoteHeaderWorkflow)
	h.RegisterWorkflow(invoiceworkflows.CreateInvoiceWorkflow)
	h.RegisterWorkflow(invoiceworkflows.UpdateInvoiceWorkflow)
	h.RegisterActivity(creditNoteHeaderActivities)
	h.RegisterActivity(debitNoteHeaderActivities)
	h.RegisterActivity(invoiceActivities)

	startWorkers(&h)

	// The workers are supposed to be long running process that should not exit.
	// Use select{} to block indefinitely for samples, you can quit by CMD+C.
	select {}
}
