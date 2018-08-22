package main

import (
	"cleaner"
	"time"
	"log"
)

const (
	defaultConfigurationPath = "etc/configuration.json"
	defaultRetentionDataPath = "etc/retention_data.json"
)

func main() {
	config, err := cleaner.MakeConfigurationFromFile(defaultConfigurationPath)
	if err != nil {
		panic(err)
	}

	customerInfo, err := cleaner.NewFileCustomerInfoFromFile(defaultRetentionDataPath, config.DefaultRetentionDays)
	if err != nil {
		panic(err)
	}

	clr := cleaner.NewCleaner(time.Now(), config.RootDir, customerInfo)
	log.Printf("Cleaning starts (root_dir: %s", config.RootDir)
	clr.Work()
	log.Printf("Cleaning done")
}
