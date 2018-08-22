package cleaner

import (
	"time"
	"io/ioutil"
	"path"
	"fmt"
	"os"
	"log"
	"strconv"
)

// TODO: move to configuration file
const (
	ConfigFilePath    = "etc/cleanup.json"
	RootDataDirectory = "data"
)

type OldCleaner struct {
	rootDir       string
	currentTime   time.Time
	customerInfo CustomerInfo
}

func NewOldCleaner(currentTime time.Time, rootDir string, customerInfo CustomerInfo) *OldCleaner {
	return &OldCleaner{
		rootDir:     rootDir,
		currentTime: currentTime,
		customerInfo: customerInfo,
	}
}

type cutoffTuple struct {
	cutoffDate time.Time
	paths []string
}

type cutoffCustomerTuple struct {
	cutoffDate time.Time
	path string
}

/* Step by step guide
1. Get all customer directories based on root directory injected
2. Create list of tuples denoting customer and his retention time, based on customer info data structure
3. Branching time, if year is older than current we don't need to further compare any dates in months/days (goes recursive_
4. Rinse and repeat for each resolution step
5. os.RemoveAll() for each path on final list

Questions:
- Do we need to get resolution going as far as to seconds? Sort out with "feature owner"
- In case of whole month qualifying for deletion: I'm not confident enough to say whether we can just RemoveAll there
  or some rate limiting should be injected
- This "script" may monitor system load and inject wait periods between removals (those will be most io intensive)
  should it get implemented?
 */
func (cl *OldCleaner) Clean() {
	//	//cutoffDate := calculateCutoffDate(cl.currentTime, cl.defaultCutoff)
	//	//cutoffYear, cutoffMonth, cutoffDay := cutoffDate.Date()
	//
	customerInfos, err := locateCustomers(cl.rootDir)
	if err != nil {
		panic(err)
	}

	customers := cl.addCutoffsToCustomers(customerInfos)
	devices := findDevicesForCustomers(customers)
	yearData := findMatchingDescendants(devices, func(datePart int64, cutoff time.Time) bool {
		return datePart <= int64(cutoff.Year())
	})
	monthData := findMatchingDescendants(yearData, func(datePart int64, cutoff time.Time) bool {
		return datePart <= int64(cutoff.Month())
	})
	fmt.Println(monthData)
	dayData := findMatchingDescendants(yearData, func(datePart int64, cutoff time.Time) bool {
		return datePart <= int64(cutoff.Day())
	})
	fmt.Println(dayData)
}

func locateCustomers(rootDir string) ([]os.FileInfo, error) {
	customerIDs, err := ioutil.ReadDir(rootDir)
	if err != nil {
		return nil, err
	}

	return customerIDs, nil
}

func (cl *OldCleaner) addCutoffsToCustomers(customerInfos []os.FileInfo) []cutoffCustomerTuple {
	var customers []cutoffCustomerTuple
	for _, customerDir := range customerInfos {
		customerID := customerDir.Name()
		cutoffDate := cl.calculateCutoffDate(customerID)

		customers = append(customers, cutoffCustomerTuple{
			cutoffDate: cutoffDate,
			path:    path.Join(cl.rootDir, customerID),
		})
	}

	return customers
}

func findDevicesForCustomers(customers []cutoffCustomerTuple) []cutoffTuple {
	var deviceDirectories []cutoffTuple
	for _, customer := range customers {
		customerPath := customer.path
		devices, err := ioutil.ReadDir(customerPath)
		if err != nil {
			log.Printf("Couldn't read customer directory, skipping (dir: %s, err: %v", customerPath, err)
			continue
		}

		var deviceDirs []string
		for _, deviceInfo := range devices {
			deviceDirs = append(deviceDirs, path.Join(customerPath, deviceInfo.Name()))
		}

		deviceDirectories = append(deviceDirectories, cutoffTuple{
			cutoffDate: customer.cutoffDate,
			paths: deviceDirs,
		})
	}

	return deviceDirectories
}

type CutoffPredicate func(date int64, cutoff time.Time) bool

func findMatchingDescendants(parentDirectories []cutoffTuple, predicate CutoffPredicate) []cutoffTuple {
	var results []cutoffTuple
	for _, parentDirectories := range parentDirectories {
		var paths []string
		for _, parentDirectory := range parentDirectories.paths {
			childDirectories, err := ioutil.ReadDir(parentDirectory)
			if err != nil {
				fmt.Printf("Couldn't read directory, skipping (dir: %s, err: %v", parentDirectory, err)
			}

			for _, childDirInfo := range childDirectories {
				childDirName := childDirInfo.Name()
				year, err := strconv.ParseInt(childDirName, 10, 64)
				if err != nil {
					log.Printf("Unable to parse folder name, skipping (name: %s, err: %v)", childDirName, err)
					continue
				}

				if predicate(year, parentDirectories.cutoffDate) {
					paths = append(paths, path.Join(parentDirectory, childDirName))
				}
			}
		}

		results = append(results, cutoffTuple{
			cutoffDate: parentDirectories.cutoffDate,
			paths:      paths,
		})
	}

	return results
}

func (cl *OldCleaner) calculateCutoffDate(customerID string) time.Time {
	retentionDays := cl.customerInfo.GetRetentionDays(customerID)

	return cl.currentTime.AddDate(0, 0, -retentionDays)
}