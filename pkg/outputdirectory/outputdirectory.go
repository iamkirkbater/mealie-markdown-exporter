package outputdirectory

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

const DefaultOutputDir = "mealie-markdown-export"

// Resolve ensures the output directory is valid and exists. If the directory
// does not exist and it is the default, it will be created. If it is a
// user-provided directory that does not exist, an error is returned.
func Resolve(fs afero.Fs, outputDir string) error {
	if outputDir == "" {
		return fmt.Errorf("output-dir is required. This defaults to %q, but it has been overwritten with an empty value", DefaultOutputDir)
	}

	isDefault := outputDir == DefaultOutputDir

	exists, err := afero.DirExists(fs, outputDir)
	if err != nil {
		return fmt.Errorf("failed to ensure the output directory %q exists: %w", outputDir, err)
	}

	if !exists && !isDefault {
		log.Errorf("When an output directory is provided, it must already exist. The directory %q does not exist.", outputDir)
		return fmt.Errorf("output directory %q does not exist", outputDir)
	}

	if isDefault {
		log.Info("Using default output directory: ", outputDir)
	} else {
		log.Info("Using provided output directory: ", outputDir)
	}

	if !exists {
		log.Warn("Default output directory '", outputDir, "' does not exist. It will be created.")
		if err := fs.MkdirAll(outputDir, 0755); err != nil {
			return fmt.Errorf("failed to create output directory %q: %w", outputDir, err)
		}
	}

	return nil
}
