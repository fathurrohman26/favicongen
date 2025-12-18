package processor

import (
	"fmt"
	"os/exec"
)

// VipsProcessor implements image processing using libvips
type VipsProcessor struct{}

func (p *VipsProcessor) Name() string {
	return "vips"
}

func (p *VipsProcessor) IsAvailable() bool {
	return commandExists("vips")
}

func (p *VipsProcessor) Resize(inputPath, outputPath string, size int) error {
	cmd := exec.Command("vips",
		"thumbnail",
		inputPath,
		outputPath,
		fmt.Sprintf("%d", size),
		"--size", "down",
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("vips resize failed: %w, output: %s", err, string(output))
	}

	return nil
}

func (p *VipsProcessor) ConvertToICO(inputPaths []string, outputPath string) error {
	// Vips doesn't support ICO creation directly, fall back to ImageMagick
	if !commandExists("convert") {
		return fmt.Errorf("ICO creation requires ImageMagick (convert command)")
	}

	magick := &ImageMagickProcessor{}
	return magick.ConvertToICO(inputPaths, outputPath)
}
