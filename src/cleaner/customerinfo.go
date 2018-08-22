package cleaner

import (
	"encoding/json"
	"os"
	"io"
)

type CustomerInfo interface {
	GetRetentionDays(string) int
}

type FileCustomerInfo struct {
	data map[string]int
	defaultRetention int
}

func NewFileCustomerInfoFromFile(jsonPath string, retentionDays int) (*FileCustomerInfo, error) {
	reader, err := os.Open(jsonPath)
	if err != nil {
		return nil, err
	}

	return NewFileCustomerInfoFromReader(reader, retentionDays)
}

func NewFileCustomerInfoFromReader(reader io.Reader, retentionDays int) (*FileCustomerInfo, error) {
	data := map[string]int{}
	decoder := json.NewDecoder(reader)

	err := decoder.Decode(&data)
	if err != nil {
		return nil, err
	}

	return &FileCustomerInfo{data: data, defaultRetention: retentionDays}, nil
}

func (fci *FileCustomerInfo) GetRetentionDays(customerID string) int {
	days, ok := fci.data[customerID]
	if !ok {
		return fci.defaultRetention
	}

	return days
}
