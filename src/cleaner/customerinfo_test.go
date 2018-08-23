package cleaner

import (
	"bytes"
	"testing"
)

func TestNewFileCustomerInfoFromReader(t *testing.T) {
	rawJson := bytes.NewBuffer([]byte(`{"1289":90,"3574":60}`))

	customerInfo, err := NewFileCustomerInfoFromReader(rawJson, 30)
	if err != nil {
		t.Fatalf("Couldn't decode json (err: %v", err)
	}

	asserts := []struct {
		key           string
		expectedValue int
	}{
		{key: "1289", expectedValue: 90},
		{key: "3574", expectedValue: 60},
		{key: "7312", expectedValue: 30},
	}

	for _, assert := range asserts {
		value := customerInfo.GetRetentionDays(assert.key)
		if value != assert.expectedValue {
			t.Errorf("Test returned wrong value (expected: %d, got: %d", assert.expectedValue, value)
		}
	}
}
