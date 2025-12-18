package processor

import (
	"fmt"
	"os/exec"
)

// ImageMagickProcessor implements image processing using ImageMagick
type ImageMagickProcessor struct{}

func (p *ImageMagickProcessor) Name() string {
	return "imagemagick"
}

func (p *ImageMagickProcessor) IsAvailable() bool {
	// ImageMagick 7+ uses 'magick', older versions use 'convert'
	return commandExists("magick") || commandExists("convert")
}

// getConvertCommand returns the appropriate command for ImageMagick convert operations
func (p *ImageMagickProcessor) getConvertCommand() (string, []string) {
	// ImageMagick 7+ uses 'magick convert' or just 'magick'
	if commandExists("magick") {
		return "magick", []string{"convert"}
	}
	// ImageMagick 6.x uses 'convert' directly
	return "convert", []string{}
}

func (p *ImageMagickProcessor) Resize(inputPath, outputPath string, size int) error {
	cmdName, baseArgs := p.getConvertCommand()
	args := append(baseArgs,
		inputPath,
		"-resize", fmt.Sprintf("%dx%d", size, size),
		"-background", "none",
		"-gravity", "center",
		"-extent", fmt.Sprintf("%dx%d", size, size),
		outputPath,
	)

	cmd := exec.Command(cmdName, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("imagemagick resize failed: %w, output: %s", err, string(output))
	}

	return nil
}

func (p *ImageMagickProcessor) ConvertToICO(inputPaths []string, outputPath string) error {
	cmdName, baseArgs := p.getConvertCommand()
	args := append(baseArgs, inputPaths...)
	args = append(args, outputPath)

	cmd := exec.Command(cmdName, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("imagemagick ico conversion failed: %w, output: %s", err, string(output))
	}

	return nil
}
