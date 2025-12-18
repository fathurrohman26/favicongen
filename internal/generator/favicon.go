package generator

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fathurrohman26/favicongen/internal/processor"
)

// FaviconGenerator handles favicon generation
type FaviconGenerator struct {
	Processor  processor.Processor
	SourcePath string
	OutputDir  string
	Sizes      []int
}

// GenerateResult contains the results of favicon generation
type GenerateResult struct {
	GeneratedFiles []string
	ICOPath        string
}

// Generate creates all favicon files
func (g *FaviconGenerator) Generate() (*GenerateResult, error) {
	// Create output directory
	if err := os.MkdirAll(g.OutputDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create output directory: %w", err)
	}

	result := &GenerateResult{
		GeneratedFiles: make([]string, 0, len(g.Sizes)),
	}

	// Generate each size
	for _, size := range g.Sizes {
		outputPath := filepath.Join(g.OutputDir, fmt.Sprintf("favicon-%dx%d.png", size, size))

		if err := g.Processor.Resize(g.SourcePath, outputPath, size); err != nil {
			return nil, fmt.Errorf("failed to generate %dx%d favicon: %w", size, size, err)
		}

		result.GeneratedFiles = append(result.GeneratedFiles, outputPath)
	}

	return result, nil
}

// GenerateICO creates a multi-resolution ICO file
func (g *FaviconGenerator) GenerateICO(pngPaths []string) (string, error) {
	icoPath := filepath.Join(g.OutputDir, "favicon.ico")

	if err := g.Processor.ConvertToICO(pngPaths, icoPath); err != nil {
		return "", fmt.Errorf("failed to generate ICO file: %w", err)
	}

	return icoPath, nil
}
