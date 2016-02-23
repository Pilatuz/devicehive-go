package devicehive

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test Equipment JSON marshaling
func TestEquipmentJson(t *testing.T) {
	equipment := NewEquipment("eqp-name", "eqp-code", "eqp-type")
	assert.JSONEq(t, toJsonStr(t, equipment), `{"name":"eqp-name","code":"eqp-code","type":"eqp-type"}`)

	equipment.Data = "custom data"
	equipment.Id = 100
	assert.JSONEq(t, toJsonStr(t, equipment), `{"id":100,"name":"eqp-name","code":"eqp-code","type":"eqp-type","data":"custom data"}`)
}
