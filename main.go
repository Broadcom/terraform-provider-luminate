package main

import (
	"github.com/Broadcom/terraform-provider-luminate/provider"
	"github.com/hashicorp/terraform/plugin"
	//"os"
	//log "github.com/sirupsen/logrus"
)

var RateLimitSleepDuration = 5

func main() {
	/*
	log.SetLevel(log.DebugLevel)
	/*
	logFile := "/tmp/luminate-terraform-provider.log"
	f, err := os.OpenFile(logFile, os.O_WRONLY | os.O_CREATE, 0755)
	if err != nil {
		log.Fatalf("Error opening file %s: %s", logFile, err)
	}
	log.SetOutput(f)
	log.Infof("Logger initialized, log file: %s", logFile)
	*/
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: provider.Provider,
	})
}
