package cleaner

import (
	"bytes"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"
)

func datePartToString(datePart int) string {
	return strconv.FormatInt(int64(datePart), 10)
}

func TestCleaner_Clean(t *testing.T) {
	rawJson := bytes.NewBuffer([]byte(`{"1289":90,"3574":60}`))

	customerInfo, err := NewFileCustomerInfoFromReader(rawJson, 30)
	if err != nil {
		t.Fatalf("Couldn't decode json (err: %v)", err)
	}

	currentDateStub := time.Date(2018, 8, 22, 0, 0, 0, 0, time.UTC)
	cutoff_1289 := currentDateStub.AddDate(0, 0, -90)
	cutoff_3574 := currentDateStub.AddDate(0, 0, -60)
	cutoff_default := currentDateStub.AddDate(0, 0, -30)

	rootDir := filepath.Join(os.TempDir(), "data")
	defer os.RemoveAll(rootDir)

	delay, _ := time.ParseDuration("")
	cleaner := NewCleaner(currentDateStub, rootDir, customerInfo, delay)

	cases := []struct {
		name            string
		client          string
		path            string
		shouldBeDeleted bool
	}{
		{
			name:            "Current year is not removed",
			client:          "1234",
			path:            "2018",
			shouldBeDeleted: false,
		},
		{
			name:            "Current month is not removed",
			client:          "1234",
			path:            "2018/08",
			shouldBeDeleted: false,
		},
		{
			name:            "Current day is not removed",
			client:          "1234",
			path:            "2018/08/22",
			shouldBeDeleted: false,
		},
		{
			name:            "Default retention year is not deleted",
			client:          "1234",
			path:            datePartToString(cutoff_default.Year()),
			shouldBeDeleted: false,
		},
		{
			name:   "Default cutoff month is not deleted",
			client: "1234",
			path: filepath.Join(
				datePartToString(cutoff_default.Year()),
				datePartToString(int(cutoff_default.Month())),
			),
			shouldBeDeleted: false,
		},
		{
			name:   "Default cutoff day is not deleted",
			client: "1234",
			path: filepath.Join(
				datePartToString(cutoff_default.Year()),
				datePartToString(int(cutoff_default.Month())),
				datePartToString(cutoff_default.Day()),
			),
			shouldBeDeleted: false,
		},
		{
			name:   "Day before default cutoff is deleted",
			client: "1234",
			path: filepath.Join(
				datePartToString(cutoff_default.Year()),
				datePartToString(int(cutoff_default.Month())),
				datePartToString(cutoff_default.Day()-1),
			),
			shouldBeDeleted: true,
		},
		{
			name:   "Month before default cutoff is deleted",
			client: "1234",
			path: filepath.Join(
				datePartToString(cutoff_default.Year()),
				datePartToString(int(cutoff_default.Month())-1),
				datePartToString(cutoff_default.Day()),
			),
			shouldBeDeleted: true,
		},
		{
			name:   "Year before default cutoff is deleted",
			client: "1234",
			path: filepath.Join(
				datePartToString(cutoff_default.Year()-1),
				datePartToString(int(cutoff_default.Month())),
				datePartToString(cutoff_default.Day()),
			),
			shouldBeDeleted: true,
		},
		{
			name:   "Client 1289 cutoff date is not deleted",
			client: "1289",
			path: filepath.Join(
				datePartToString(cutoff_1289.Year()),
				datePartToString(int(cutoff_1289.Month())),
				datePartToString(cutoff_1289.Day()),
			),
			shouldBeDeleted: false,
		},
		{
			name:   "Client 3574 cutoff date is not deleted",
			client: "3574",
			path: filepath.Join(
				datePartToString(cutoff_3574.Year()),
				datePartToString(int(cutoff_3574.Month())),
				datePartToString(cutoff_3574.Day()),
			),
			shouldBeDeleted: false,
		},
		{
			name:   "Day before cutoff date for Client 1289 is deleted",
			client: "1289",
			path: filepath.Join(
				datePartToString(cutoff_1289.Year()),
				datePartToString(int(cutoff_1289.Month())),
				datePartToString(cutoff_1289.Day()-1),
			),
			shouldBeDeleted: true,
		},
		{
			name:   "Day before cutoff date for Client 3574 is deleted",
			client: "3574",
			path: filepath.Join(
				datePartToString(cutoff_3574.Year()),
				datePartToString(int(cutoff_3574.Month())),
				datePartToString(cutoff_3574.Day()-1),
			),
			shouldBeDeleted: true,
		},
		{
			name:   "Future date is not deleted",
			client: "1234",
			path: filepath.Join(
				datePartToString(cutoff_default.Year()),
				datePartToString(int(cutoff_default.Month())),
				datePartToString(cutoff_default.Day()+1),
			),
			shouldBeDeleted: false,
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			fullPath := filepath.Join(rootDir, testCase.client, "1234", testCase.path)
			os.MkdirAll(fullPath, os.ModePerm)

			cleaner.Work()

			_, err := os.Stat(fullPath)
			if os.IsNotExist(err) && !testCase.shouldBeDeleted {
				t.Errorf("File got deleted while it shouldn't")
			} else if err == nil && testCase.shouldBeDeleted {
				t.Errorf("File is present on disk but should get deleted")
			} else if err != nil && !os.IsNotExist(err) {
				t.Fatalf("Unexpected error encountered (err: %v)", err)
			}
		})
	}
}
