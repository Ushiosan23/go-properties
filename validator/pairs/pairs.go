package pairs

import (
	"errors"
	"strings"
)

// ############################################################
// Functions
// ############################################################

// ValidateKey Checks if data is a valid key name.
// Returns an error if data is not valid
func ValidateKey(data string) error {
	// Clean data
	data = strings.TrimSpace(data)
	// Check data size
	if len(data) == 0 {
		return errors.New("invalid key name. the key cannot be empty")
	}
	// Not error
	return nil
}
