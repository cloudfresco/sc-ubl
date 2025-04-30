package invoicecontrollers

import (
	"context"
	"net/http"
	"os"
	"path/filepath"

	"github.com/cloudfresco/sc-ubl/internal/common"
	"github.com/cloudfresco/sc-ubl/internal/config"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"

	invoiceproto "github.com/cloudfresco/sc-ubl/internal/protogen/invoice/v1"
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

	invoiceconn, err := grpc.NewClient(grpcServerOpt.GrpcInvoiceServerPort, grpc.WithTransportCredentials(creds), grpc.WithStatsHandler(otelgrpc.NewClientHandler()))
	if err != nil {
		log.Error("Error", zap.Error(err))
		return err
	}

	u := partyproto.NewUserServiceClient(userconn)
	c := invoiceproto.NewCreditNoteHeaderServiceClient(invoiceconn)
	d := invoiceproto.NewDebitNoteHeaderServiceClient(invoiceconn)
	i := invoiceproto.NewInvoiceServiceClient(invoiceconn)

	initCreditNoteHeaders(mux, serverOpt, log, u, c, h, workflowClient)
	initDebitNoteHeaders(mux, serverOpt, log, u, d, h, workflowClient)
	initInvoiceHeaders(mux, serverOpt, log, u, i, h, workflowClient)

	return nil
}

func initCreditNoteHeaders(mux *http.ServeMux, serverOpt *config.ServerOptions, log *zap.Logger, u partyproto.UserServiceClient, c invoiceproto.CreditNoteHeaderServiceClient, wfHelper common.WfHelper, workflowClient client.Client) {
	cc := NewCreditNoteHeaderController(log, u, c, h, workflowClient, serverOpt)

	mux.Handle("GET /v2.3/credit-notes", http.HandlerFunc(cc.Index))
	mux.Handle("GET /v2.3/credit-notes/{id}", http.HandlerFunc(cc.Show))
	mux.Handle("GET /v2.3/credit-notes/{id}/lines", http.HandlerFunc(cc.GetCreditNoteLines))

	mux.Handle("POST /v2.3/credit-notes", http.HandlerFunc(cc.CreateCreditNoteHeader))

	mux.Handle("PUT /v2.3/credit-notes/{id}", http.HandlerFunc(cc.UpdateCreditNoteHeader))
}

func initDebitNoteHeaders(mux *http.ServeMux, serverOpt *config.ServerOptions, log *zap.Logger, u partyproto.UserServiceClient, d invoiceproto.DebitNoteHeaderServiceClient, wfHelper common.WfHelper, workflowClient client.Client) {
	dc := NewDebitNoteHeaderController(log, u, d, h, workflowClient, serverOpt)

	mux.Handle("GET /v2.3/debit-notes", http.HandlerFunc(dc.Index))
	mux.Handle("GET /v2.3/debit-notes/{id}", http.HandlerFunc(dc.Show))
	mux.Handle("GET /v2.3/debit-notes/{id}/lines", http.HandlerFunc(dc.GetDebitNoteLines))

	mux.Handle("POST /v2.3/debit-notes", http.HandlerFunc(dc.CreateDebitNoteHeader))

	mux.Handle("PUT /v2.3/debit-notes/{id}", http.HandlerFunc(dc.UpdateDebitNoteHeader))
}

func initInvoiceHeaders(mux *http.ServeMux, serverOpt *config.ServerOptions, log *zap.Logger, u partyproto.UserServiceClient, i invoiceproto.InvoiceServiceClient, wfHelper common.WfHelper, workflowClient client.Client) {
	ic := NewInvoiceHeaderController(log, u, i, h, workflowClient, serverOpt)

	mux.Handle("GET /v2.3/invoices", http.HandlerFunc(ic.Index))
	mux.Handle("GET /v2.3/invoices/{id}", http.HandlerFunc(ic.Show))
	mux.Handle("GET /v2.3/invoices/{id}/lines", http.HandlerFunc(ic.GetInvoiceLines))

	mux.Handle("POST /v2.3/invoices", http.HandlerFunc(ic.CreateInvoice))

	mux.Handle("PUT /v2.3/invoices/{id}", http.HandlerFunc(ic.UpdateInvoice))
}
