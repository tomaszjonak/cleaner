package cleaner

import (
	"testing"
	"time"
	"fmt"
	"os"
	"bytes"
)

func TestCleaner_Clean(t *testing.T) {
	currentDateStub := time.Date(2018,8,22,0,0,0,0, time.UTC)

	rootDir := "data"
	fmt.Println(os.Getwd())

	rawJson := bytes.NewBuffer([]byte(`{"1289":90,"3574":60}`))

	customerInfo, err := NewFileCustomerInfoFromReader(rawJson)
	if err != nil {
		t.Fatalf("Couldn't decode json (err: %v", err)
	}

	cleaner := NewCleaner(currentDateStub, rootDir, customerInfo)
	cleaner.Work()
}