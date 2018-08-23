package main

import (
	"cleaner"
	"time"
	"log"
	"flag"
)

const (
	defaultConfigurationPath = "etc/configuration.json"
	defaultRetentionDataPath = "etc/retention_data.json"
)

var configPath = flag.String("config", defaultConfigurationPath, "Path to configuration file")
var retentionDataPath = flag.String("retention_data", defaultRetentionDataPath,
	"Path to client specific retention_data")

func main() {
	flag.Parse()

	config, err := cleaner.MakeConfigurationFromFile(*configPath)
	if err != nil {
		panic(err)
	}

	customerInfo, err := cleaner.NewFileCustomerInfoFromFile(*retentionDataPath, config.DefaultRetentionDays)
	if err != nil {
		panic(err)
	}

	clr := cleaner.NewCleaner(time.Now(), config.RootDir, customerInfo)
	log.Printf("Cleaning starts (root_dir: %s)", config.RootDir)
	clr.Work()
	log.Printf("Cleaning done")
}
