package main

import (
	"os"

	"github.com/cloudfresco/sc-ubl/internal/common"
	"github.com/cloudfresco/sc-ubl/internal/config"
	itemservices "github.com/cloudfresco/sc-ubl/internal/services/itemservices"
	_ "github.com/go-sql-driver/mysql" // mysql
	"go.uber.org/zap"
)

func main() {
	v, err := config.GetViper()
	if err != nil {
		os.Exit(1)
	}

	logOpt, err := config.GetLogConfig(v)
	if err != nil {
		os.Exit(1)
	}

	log := config.SetUpLogging(logOpt.ItemPath)

	dbOpt, err := config.GetDbConfig(log, v, false, "SC_UBL_DB", "SC_UBL_DBHOST", "SC_UBL_DBPORT", "SC_UBL_DBUSER", "SC_UBL_DBPASS", "SC_UBL_DBNAME", "", "", "", "", "", "")
	if err != nil {
		log.Error("Error", zap.Error(err))
		os.Exit(1)
	}

	jwtOpt, err := config.GetJWTConfig(log, v, false, "SC_UBL_JWT_KEY", "SC_UBL_JWT_DURATION")
	if err != nil {
		log.Error("Error", zap.Error(err))
		os.Exit(1)
	}

	redisOpt, mailerOpt, _, grpcServerOpt, oauthOpt, userOpt, uptraceOpt := config.GetConfigOpt(log, v)

	dbService, redisService, mailerService := common.GetServices(log, false, dbOpt, redisOpt, jwtOpt, mailerOpt)

	pwd, _ := os.Getwd()
	itemservices.StartItemServer(log, false, pwd, dbOpt, redisOpt, mailerOpt, grpcServerOpt, jwtOpt, oauthOpt, userOpt, uptraceOpt, dbService, redisService, mailerService)
}
