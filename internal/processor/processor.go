package processor

import (
	"fmt"
	"os/exec"
)

// Processor defines the interface for image processing backends
type Processor interface {
	// Name returns the name of the processor
	Name() string

	// IsAvailable checks if the processor is available on the system
	IsAvailable() bool

	// Resize resizes an image to the specified dimensions
	Resize(inputPath, outputPath string, size int) error

	// ConvertToICO converts multiple PNGs to a single ICO file
	ConvertToICO(inputPaths []string, outputPath string) error
}

// DetectAvailableProcessor detects which image processor is available
func DetectAvailableProcessor(preferred string) (Processor, error) {
	processors := []Processor{
		&ImageMagickProcessor{},
		&VipsProcessor{},
	}

	// If a preferred processor is specified, try it first
	if preferred != "" {
		for _, p := range processors {
			if p.Name() == preferred && p.IsAvailable() {
				return p, nil
			}
		}
		return nil, fmt.Errorf("preferred processor %s is not available", preferred)
	}

	// Otherwise, return the first available processor
	for _, p := range processors {
		if p.IsAvailable() {
			return p, nil
		}
	}

	return nil, fmt.Errorf("no image processor available (install ImageMagick or libvips)")
}

// commandExists checks if a command is available in PATH
func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}
