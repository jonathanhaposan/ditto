package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)


func TestIsEqualJson(t *testing.T) {
	jsonA := `{"payouts":[{"beneficiary_name":"beneficiary error case","beneficiary_account":"123456","beneficiary_bank":"mandiri","amount":"20396","notes":"2002741746489L8E","reference_no":"ID2002741746489L8E"}]}`
	jsonB := `{"payouts":[{"beneficiary_name":"success case","beneficiary_account":"123456","beneficiary_bank":"mandiri"}]}`

	assert.False(t, IsEqualJson([]byte(jsonA), []byte(jsonB)))
}
