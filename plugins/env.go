package plugins

import (
	"os"
	"regexp"
	"strings"
)

var (
	environmentRegex = regexp.MustCompile(`\${([a-zA-Z0-9_]+)}`)
)

// ############################################################
// Plugins
// ############################################################

func EnvironmentDetector(data string) string {
	return environmentRegex.ReplaceAllStringFunc(data, func(match string) string {
		// Clean key
		key := strings.TrimSpace(match[2 : len(match)-1])
		// Check environment variables
		for _, v := range os.Environ() {
			// Split elements
			env := strings.SplitN(v, "=", 2)
			envKey := strings.TrimSpace(env[0])
			envValue := strings.TrimSpace(env[1])
			// Compare elements
			if key == envKey {
				return envValue
			}
		}
		return match
	})
}
