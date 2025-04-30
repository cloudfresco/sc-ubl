package main

import (
	"os"

	"github.com/cloudfresco/sc-ubl/internal/config"
	"github.com/cloudfresco/sc-ubl/internal/workers/logisticsworkers"
	"go.uber.org/zap"
)

func main() {
	v, err := config.GetViper()
	if err != nil {
		os.Exit(1)
	}

	configFilePath := v.GetString("SC_UBL_WORKFLOW_CONFIG_FILE_PATH")

	logOpt, err := config.GetLogConfig(v)
	if err != nil {
		os.Exit(1)
	}

	log := config.SetUpLogging(logOpt.Path)

	_, _, _, grpcServerOpt, _, _, _ := config.GetConfigOpt(log, v)
	if err != nil {
		log.Error("Error", zap.Error(err))
		os.Exit(1)
	}

	pwd, _ := os.Getwd()

	logisticsworkers.StartLogisticsWorker(log, false, pwd, grpcServerOpt, configFilePath)
}
