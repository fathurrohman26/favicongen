package processor

import (
	"testing"
)

func TestCommandExists(t *testing.T) {
	tests := []struct {
		name     string
		cmd      string
		wantExist bool
	}{
		{
			name:      "existing command (go)",
			cmd:       "go",
			wantExist: true,
		},
		{
			name:      "non-existing command",
			cmd:       "nonexistent_command_xyz_123",
			wantExist: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := commandExists(tt.cmd)
			if got != tt.wantExist {
				t.Errorf("commandExists(%q) = %v, want %v", tt.cmd, got, tt.wantExist)
			}
		})
	}
}

func TestImageMagickProcessorName(t *testing.T) {
	p := &ImageMagickProcessor{}
	if got := p.Name(); got != "imagemagick" {
		t.Errorf("ImageMagickProcessor.Name() = %q, want %q", got, "imagemagick")
	}
}

func TestVipsProcessorName(t *testing.T) {
	p := &VipsProcessor{}
	if got := p.Name(); got != "vips" {
		t.Errorf("VipsProcessor.Name() = %q, want %q", got, "vips")
	}
}

func TestDetectAvailableProcessor(t *testing.T) {
	t.Run("no preferred processor", func(t *testing.T) {
		proc, err := DetectAvailableProcessor("")
		// Either we get a processor or an error about no processor available
		if err != nil {
			if err.Error() != "no image processor available (install ImageMagick or libvips)" {
				t.Errorf("unexpected error: %v", err)
			}
			return
		}
		if proc == nil {
			t.Error("expected non-nil processor")
		}
	})

	t.Run("non-existent preferred processor", func(t *testing.T) {
		_, err := DetectAvailableProcessor("nonexistent")
		if err == nil {
			t.Error("expected error for non-existent processor")
		}
	})
}

func TestImageMagickProcessorIsAvailable(t *testing.T) {
	p := &ImageMagickProcessor{}
	// Just verify it doesn't panic - actual availability depends on system
	_ = p.IsAvailable()
}

func TestVipsProcessorIsAvailable(t *testing.T) {
	p := &VipsProcessor{}
	// Just verify it doesn't panic - actual availability depends on system
	_ = p.IsAvailable()
}

func TestImageMagickGetConvertCommand(t *testing.T) {
	p := &ImageMagickProcessor{}
	cmdName, baseArgs := p.getConvertCommand()

	// Should return either "magick" with ["convert"] or "convert" with []
	if cmdName == "magick" {
		if len(baseArgs) != 1 || baseArgs[0] != "convert" {
			t.Errorf("expected baseArgs=[convert] for magick, got %v", baseArgs)
		}
	} else if cmdName == "convert" {
		if len(baseArgs) != 0 {
			t.Errorf("expected empty baseArgs for convert, got %v", baseArgs)
		}
	} else {
		t.Errorf("unexpected command name: %s", cmdName)
	}
}
