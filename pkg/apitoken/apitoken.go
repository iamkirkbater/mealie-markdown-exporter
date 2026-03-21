package apitoken

import (
	"fmt"
	"strings"

	"github.com/spf13/afero"
)

// Resolve resolves an API token string. If the token has a "file://" prefix,
// the token is read from the specified file path on the given filesystem.
// Otherwise it is returned as-is.
func Resolve(fs afero.Fs, token string) (string, error) {
	if filePath, ok := strings.CutPrefix(token, "file://"); ok {
		data, err := afero.ReadFile(fs, filePath)
		if err != nil {
			return "", fmt.Errorf("failed to read api-token from file %q: %w", filePath, err)
		}
		return strings.TrimSpace(string(data)), nil
	}
	return token, nil
}
