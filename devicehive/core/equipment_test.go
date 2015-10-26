package core

import (
	"encoding/json"
	"testing"
)

// Test Equipment JSON marshaling
func TestJsonEquipment(t *testing.T) {
	check := func(equipment Equipment, expectedJson string) {
		jsonBytes, err := json.Marshal(equipment)
		if err != nil {
			t.Errorf("Equipment: Cannot convert %+v to JSON (error: %s)", equipment, err)
		}
		if expectedJson != string(jsonBytes) {
			t.Errorf("Equipment: Cannot convert %+v to JSON\n\t   found: %s,\n\texpected:%s",
				equipment, string(jsonBytes), expectedJson)
		}

		t.Logf("Equipment: %+v converted to %s", equipment, string(jsonBytes))
	}

	equipment := Equipment{
		Name: "eqp-name",
		Code: "eqp-code",
		Type: "eqp-type"}

	check(equipment, `{"name":"eqp-name","code":"eqp-code","type":"eqp-type"}`)

	equipment.Data = "custom data"
	equipment.Id = 100
	check(equipment, `{"id":100,"name":"eqp-name","code":"eqp-code","type":"eqp-type","data":"custom data"}`)
}
