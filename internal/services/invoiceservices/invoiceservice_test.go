package invoiceservices

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/cloudfresco/sc-ubl/internal/common"
	"github.com/cloudfresco/sc-ubl/internal/config"
	partyproto "github.com/cloudfresco/sc-ubl/internal/protogen/party/v1"
	"github.com/cloudfresco/sc-ubl/internal/services/partyservices"
	"github.com/cloudfresco/sc-ubl/test"
	_ "github.com/go-sql-driver/mysql"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

var (
	dbService         *common.DBService
	redisService      *common.RedisService
	userServiceClient partyproto.UserServiceClient
	mailerService     common.MailerIntf
	jwtOpt            *config.JWTOptions
	userTestOpt       *config.UserTestOptions
	redisOpt          *config.RedisOptions
	mailerOpt         *config.MailerOptions
	serverOpt         *config.ServerOptions
	grpcServerOpt     *config.GrpcServerOptions
	oauthOpt          *config.OauthOptions
	userOpt           *config.UserOptions
	uptraceOpt        *config.UptraceOptions
	log               *zap.Logger
	logInvoice        *zap.Logger
	logUser           *zap.Logger
	Layout            string
)

func TestMain(m *testing.M) {
	v, err := config.GetViper()
	if err != nil {
		log.Error("Error", zap.Error(err))
		os.Exit(1)
	}

	logOpt, err := config.GetLogConfig(v)
	if err != nil {
		log.Error("Error", zap.Error(err))
		os.Exit(1)
	}

	log = config.SetUpLogging(logOpt.Path)
	logUser = config.SetUpLogging(logOpt.UserPath)
	logInvoice = config.SetUpLogging(logOpt.InvoicePath)
	Layout = "2006-01-02T15:04:05Z"

	dbOpt, err := config.GetDbConfig(log, v, true, "SC_UBL_DB", "SC_UBL_DBHOST", "SC_UBL_DBPORT", "SC_UBL_DBUSER_TEST", "SC_UBL_DBPASS_TEST", "SC_UBL_DBNAME_TEST", "SC_UBL_DBSQL_MYSQL_TEST", "SC_UBL_DBSQL_MYSQL_SCHEMA", "SC_UBL_DBSQL_MYSQL_TRUNCATE", "SC_UBL_DBSQL_PGSQL_TEST", "SC_UBL_DBSQL_PGSQL_SCHEMA", "SC_UBL_DBSQL_PGSQL_TRUNCATE")
	if err != nil {
		log.Error("Error", zap.Error(err))
		return
	}

	jwtOpt, err = config.GetJWTConfig(log, v, true, "SC_UBL_JWT_KEY_TEST", "SC_UBL_JWT_DURATION_TEST")
	if err != nil {
		log.Error("Error", zap.Error(err))
		return
	}

	userTestOpt, err = config.GetUserTestConfig(log, v)
	if err != nil {
		log.Error("Error", zap.Error(err))
		return
	}

	redisOpt, mailerOpt, serverOpt, grpcServerOpt, oauthOpt, userOpt, uptraceOpt = config.GetConfigOpt(log, v)

	dbService, redisService, _ = common.GetServices(log, true, dbOpt, redisOpt, jwtOpt, mailerOpt)

	mailerService, err = test.CreateMailerServiceTest(log)
	if err != nil {
		log.Error("Error", zap.Error(err))
	}

	pwd, _ := os.Getwd()
	go partyservices.StartUserServer(logUser, true, pwd, dbOpt, redisOpt, mailerOpt, serverOpt, grpcServerOpt, jwtOpt, oauthOpt, userOpt, uptraceOpt, dbService, redisService, mailerService)
	go StartInvoiceServer(logInvoice, true, pwd, dbOpt, redisOpt, mailerOpt, grpcServerOpt, jwtOpt, oauthOpt, userOpt, uptraceOpt, dbService, redisService, mailerService)

	keyPath := filepath.Join(pwd, filepath.FromSlash("/../../../")+filepath.FromSlash(grpcServerOpt.GrpcCaCertPath))
	creds, err := credentials.NewClientTLSFromFile(keyPath, "localhost")
	if err != nil {
		log.Error("Error", zap.Error(err))
	}

	userconn, err := grpc.NewClient(grpcServerOpt.GrpcUserServerPort, grpc.WithTransportCredentials(creds), grpc.WithStatsHandler(otelgrpc.NewClientHandler()))
	if err != nil {
		log.Error("Error", zap.Int("msgnum", 8301), zap.Error(err))
		return
	}

	userServiceClient = partyproto.NewUserServiceClient(userconn)

	os.Exit(m.Run())
}

func LoginUser() context.Context {
	md := metadata.Pairs("authorization", "Bearer "+userTestOpt.Tokenstring)
	ctxNew := metadata.NewIncomingContext(context.Background(), md)
	return ctxNew
}
