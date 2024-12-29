// internal/utils/utils.go
package utils

import (
	"os"
	"path/filepath"
	"strings"
)

// ExpandPath expands the tilde (~) in a path to the user's home directory
func ExpandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return path
		}
		return filepath.Join(home, path[2:])
	}
	return path
}
