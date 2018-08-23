package cleaner

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type cutoffData struct {
	cutoffDate time.Time
	path       string
}

type Cleaner struct {
	rootDir      string
	currentTime  time.Time
	customerInfo CustomerInfo
	toWipe       chan string
	removalDelay time.Duration
}

const (
	WipeQueueLength      = 2
	HarvesterQueueLength = 1
)

func NewCleaner(currentTime time.Time, rootDir string, customerInfo CustomerInfo, removalDelay time.Duration) *Cleaner {
	return &Cleaner{
		rootDir:      rootDir,
		currentTime:  currentTime,
		customerInfo: customerInfo,
		removalDelay: removalDelay,
	}
}

// TODO: maybe put rootdir here?
func (cl *Cleaner) Work() {
	clientPaths := make(chan string, HarvesterQueueLength)
	clientToYear := make(chan cutoffData, HarvesterQueueLength)
	yearToMonth := make(chan cutoffData, HarvesterQueueLength)
	monthToDay := make(chan cutoffData, HarvesterQueueLength)
	dayToFin := make(chan cutoffData, HarvesterQueueLength)
	cl.toWipe = make(chan string, WipeQueueLength)

	go clientFinder(cl.rootDir, clientPaths)
	go cl.cutoffAdder(clientPaths, clientToYear)

	go cl.harvester(clientToYear, yearToMonth, func(cutoffDate time.Time) uint64 {
		return uint64(cutoffDate.Year())
	})
	go cl.harvester(yearToMonth, monthToDay, func(cutoffDate time.Time) uint64 {
		return uint64(cutoffDate.Month())
	})
	go cl.harvester(monthToDay, dayToFin, func(cutoffDate time.Time) uint64 {
		return uint64(cutoffDate.Day())
	})

	go cl.deadEnd(dayToFin)

	cl.wipeRoutine(cl.toWipe)
}

// clientFinder scans root directory searching for client directories
func clientFinder(rootDir string, clientPaths chan<- string) {
	// Little assumption that there's nothing but directories in this folder
	customerDirs, err := filepath.Glob(rootDir + "/*")
	if err != nil {
		panic(err)
	}

	for _, customerDir := range customerDirs {
		clientPaths <- customerDir
	}

	close(clientPaths)
}

// cutoffAdder calculates cutoff date based on customer name and customerInfo structure
// also
func (cl *Cleaner) cutoffAdder(paths <-chan string, toHarvesters chan<- cutoffData) {
	for path := range paths {
		clientName := filepath.Base(path)

		retentionDays := cl.customerInfo.GetRetentionDays(clientName)
		cutoffDate := cl.currentTime.AddDate(0, 0, -retentionDays)
		deviceDirs, err := filepath.Glob(path + "/*")
		if err != nil {
			log.Printf("Couldn't get devices for client (err: %v, client: %s)", err, clientName)
		}

		for _, deviceDir := range deviceDirs {
			toHarvesters <- cutoffData{path: deviceDir, cutoffDate: cutoffDate}
		}
	}

	close(toHarvesters)
}

// wipeRoutine goroutine function responsible for removal of data,
// rate limiting based on some data should be done here
// Ideas
// 1. Inject configurable/randomized timeout between each removal
// 2. Monitor system load in another goroutine, inject sleep period between files
// Moreover, in case this app finds whole year of data to delete, some heuristic should
// be used to split this task into smaller pieces
func (cl *Cleaner) wipeRoutine(toWipe <-chan string) {
	for path := range toWipe {
		os.RemoveAll(path)
		log.Printf("Wiped (%s)", path)
		<-time.After(cl.removalDelay)
	}
}

type Extractor func(time.Time) uint64

// harvester function implementing interface to filter out folders at given depth
// based on information from cutoff date, which is added as function transforming time.Time into uint64
// Three cases are considered here
// 1. Cutoff part at current depth is equal to folder name: pass to next harvester (increase depth)
// 2. Cutoff part at current depth is bigger than folder name: schedule deletion
// 3. Cutoff part at current depth is less than folder name: ignore, data storage is still required
func (cl *Cleaner) harvester(inputs <-chan cutoffData, toNextHarvester chan<- cutoffData, extractor Extractor) {
	for input := range inputs {
		// yet again, check for directories could be added here
		directories, err := filepath.Glob(input.path + "/*")
		if err != nil {
			log.Printf("Couldn't read directory contents, skipping (err: %v, dir: %s)", err, input.path)
			continue
		}

		for _, directory := range directories {
			datePart := filepath.Base(directory)
			datePartNum, err := strconv.ParseUint(datePart, 10, 64)
			if err != nil {
				log.Printf("Unable to parse integer from datePart (err: %v, part: %s)", err, datePart)
			}

			cutoff := extractor(input.cutoffDate)
			if datePartNum == cutoff {
				toNextHarvester <- cutoffData{cutoffDate: input.cutoffDate, path: directory}
			} else if datePartNum < cutoff {
				cl.toWipe <- directory
			}
		}
	}

	close(toNextHarvester)
}

// deadEnd fetches paths from last harvester and does nothing
// honestly its there just to keep Harvesters identical
func (cl *Cleaner) deadEnd(inputs <-chan cutoffData) {
	for range inputs {
		// maybe get some stats about files which are kept
	}

	close(cl.toWipe)
}
