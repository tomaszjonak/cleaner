package cleaner

import (
	"testing"
	"time"
	"bytes"
)

func TestCleaner_Clean(t *testing.T) {
	currentDateStub := time.Date(2018,8,22,0,0,0,0, time.UTC)

	rootDir := "data"

	rawJson := bytes.NewBuffer([]byte(`{"1289":90,"3574":60}`))

	customerInfo, err := NewFileCustomerInfoFromReader(rawJson, 30)
	if err != nil {
		t.Fatalf("Couldn't decode json (err: %v)", err)
	}

	cleaner := NewCleaner(currentDateStub, rootDir, customerInfo)
	cleaner.Work()
}
