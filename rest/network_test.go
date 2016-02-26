package rest

import (
	"testing"

	"github.com/pilatuz/go-devicehive"
	"github.com/stretchr/testify/assert"
)

// Test GetNetworkList method
func TestNetworkList(t *testing.T) {
	service := testNewRest(t)
	if service == nil {
		return // nothing to test
	}

	networks, err := service.GetNetworkList(0, 0, testWaitTimeout)
	assert.NoError(t, err, "Failed to get list of networks")
	assert.NotEmpty(t, networks, "No any network avaialble")
	// t.Logf("networks: %s", networks)
}

// Test GetNetwork method
func TestNetworkGet(t *testing.T) {
	service := testNewRest(t)
	if service == nil {
		return // nothing to test
	}

	networks, err := service.GetNetworkList(0, 0, testWaitTimeout)
	assert.NoError(t, err, "Failed to get list of networks")
	assert.NotEmpty(t, networks, "No any network avaialble")
	// t.Logf("networks: %s", networks)

	for _, a := range networks {
		b, err := service.GetNetwork(a.ID, testWaitTimeout)
		assert.NoError(t, err, "Failed to get network")
		assert.NotNil(t, b, "No network avaialble")
		assert.JSONEq(t, toJsonStr(a), toJsonStr(b), "Networks are not the same")
	}
}

// Test UpdateNetwork method
func TestNetworkUpdate(t *testing.T) {
	return // IGNORED, DOESN'T WORK with playground

	service := testNewRest(t)
	if service == nil {
		return // nothing to test
	}

	networks, err := service.GetNetworkList(0, 0, testWaitTimeout)
	assert.NoError(t, err, "Failed to get list of networks")
	assert.NotEmpty(t, networks, "No any network avaialble")
	// t.Logf("networks: %s", networks)

	for _, a := range networks {
		a.Description += "-updated"
		err := service.UpdateNetwork(a, testWaitTimeout)
		assert.NoError(t, err, "Failed to update network")
	}
}

// Test InsertNetwork DeleteNetwork methods
func TestNetworkInsertAndDelete(t *testing.T) {
	return // IGNORED, DOESN'T WORK with playground

	service := testNewRest(t)
	if service == nil {
		return // nothing to test
	}

	network := devicehive.NewNetwork("test-net", "no-secure-key")
	err := service.InsertNetwork(network, testWaitTimeout)
	assert.NoError(t, err, "Failed to insert network")
	assert.NotEmpty(t, network.ID, "No network identifier provided")

	err = service.DeleteNetwork(network, testWaitTimeout)
	assert.NoError(t, err, "Failed to delete network")
}
