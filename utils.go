package devicehive

import (
	"github.com/mitchellh/mapstructure"
)

const (
	// DateTimeLayout is layout used for timestamps.
	DateTimeLayout = "2006-01-02T15:04:05.999"
)

// Assign fields from JSON map.
// This method is used to assign already parsed JSON data.
func FromJSON(result interface{}, data interface{}) error {
	config := new(mapstructure.DecoderConfig)
	config.WeaklyTypedInput = true
	config.Result = result
	config.TagName = "json"

	dec, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}
	return dec.Decode(data)
}
