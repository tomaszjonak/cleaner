package main

import (
	"cleaner"
	"flag"
	"log"
	"time"
)

const (
	defaultConfigurationPath = "etc/configuration.json"
	defaultRetentionDataPath = "etc/retention_data.json"
)

var configPath = flag.String("config", defaultConfigurationPath, "Path to configuration file")
var retentionDataPath = flag.String("retention_data", defaultRetentionDataPath,
	"Path to client specific retention_data")
var removalDelay = flag.String("delay", "0s", "Sets up wait period between consecutive removals")

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

	delay, err := time.ParseDuration(*removalDelay)
	if err != nil {
		panic(err)
	}

	currentTime := time.Now()
	clr := cleaner.NewCleaner(currentTime, config.RootDir, customerInfo, delay)
	log.Printf("Cleaning starts (root_dir: %s)", config.RootDir)
	log.Printf("Default cutoff date: %s",
		currentTime.AddDate(0,0,-config.DefaultRetentionDays).Format(time.RFC3339)[:10])

	clr.Work()

	log.Printf("Cleaning done")
}
