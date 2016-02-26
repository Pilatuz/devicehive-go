package devicehive

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TODO: Network to string
// TODO: Network from map

// Test Network JSON marshaling
func TestJsonNetwork(t *testing.T) {
	network := NewNetwork("net-name", "net-key")
	network.Description = "custom description"
	assert.JSONEq(t, toJsonStr(t, network), `{"name":"net-name","key":"net-key","description":"custom description"}`)

	network.Description = ""
	network.ID = 100
	assert.JSONEq(t, toJsonStr(t, network), `{"id":100,"name":"net-name","key":"net-key"}`)
}
