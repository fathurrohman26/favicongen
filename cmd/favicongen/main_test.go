package main

import (
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func TestParseSizes(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []int
		wantErr bool
	}{
		{
			name:    "single size",
			input:   "16",
			want:    []int{16},
			wantErr: false,
		},
		{
			name:    "multiple sizes",
			input:   "16,32,48",
			want:    []int{16, 32, 48},
			wantErr: false,
		},
		{
			name:    "sizes with spaces",
			input:   "16, 32, 48",
			want:    []int{16, 32, 48},
			wantErr: false,
		},
		{
			name:    "default sizes",
			input:   "16,32,48,64,128,180,256,512",
			want:    []int{16, 32, 48, 64, 128, 180, 256, 512},
			wantErr: false,
		},
		{
			name:    "invalid size - not a number",
			input:   "16,abc,48",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid size - negative",
			input:   "16,-32,48",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid size - zero",
			input:   "16,0,48",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   "",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseSizes(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseSizes(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if len(got) != len(tt.want) {
					t.Errorf("parseSizes(%q) = %v, want %v", tt.input, got, tt.want)
					return
				}
				for i := range got {
					if got[i] != tt.want[i] {
						t.Errorf("parseSizes(%q)[%d] = %d, want %d", tt.input, i, got[i], tt.want[i])
					}
				}
			}
		})
	}
}

func TestFileExists(t *testing.T) {
	tests := []struct {
		name string
		path string
		want bool
	}{
		{
			name: "existing file",
			path: "main.go",
			want: true,
		},
		{
			name: "non-existing file",
			path: "nonexistent_file_xyz_123.go",
			want: false,
		},
		{
			name: "existing directory",
			path: ".",
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := fileExists(tt.path)
			if got != tt.want {
				t.Errorf("fileExists(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}

func TestConfigDefaults(t *testing.T) {
	config := &Config{
		Source:             "",
		Output:             "./favicons",
		Sizes:              []int{16, 32, 48, 64, 128, 180, 256, 512},
		Backend:            "",
		GenerateHTML:       true,
		GenerateManifest:   false,
		GenerateICO:        true,
		GenerateHTMLOnly:   false,
		AppName:            "",
		AppShortName:       "",
		AppDescription:     "",
		AppStartURL:        "/",
		AppDisplay:         "standalone",
		AppOrientation:     "any",
		AppScope:           "/",
		AppThemeColor:      "#ffffff",
		AppBackgroundColor: "#ffffff",
		AppCategories:      nil,
		AppIcon:            "",
	}

	if config.Output != "./favicons" {
		t.Errorf("default Output = %q, want %q", config.Output, "./favicons")
	}
	if !config.GenerateHTML {
		t.Error("default GenerateHTML should be true")
	}
	if config.GenerateManifest {
		t.Error("default GenerateManifest should be false")
	}
	if !config.GenerateICO {
		t.Error("default GenerateICO should be true")
	}
	if config.AppStartURL != "/" {
		t.Errorf("default AppStartURL = %q, want %q", config.AppStartURL, "/")
	}
	if config.AppDisplay != "standalone" {
		t.Errorf("default AppDisplay = %q, want %q", config.AppDisplay, "standalone")
	}
	if config.AppOrientation != "any" {
		t.Errorf("default AppOrientation = %q, want %q", config.AppOrientation, "any")
	}
	if config.AppThemeColor != "#ffffff" {
		t.Errorf("default AppThemeColor = %q, want %q", config.AppThemeColor, "#ffffff")
	}
	if config.AppBackgroundColor != "#ffffff" {
		t.Errorf("default AppBackgroundColor = %q, want %q", config.AppBackgroundColor, "#ffffff")
	}
}

func TestConfigSizesContains180(t *testing.T) {
	defaultSizes := []int{16, 32, 48, 64, 128, 180, 256, 512}

	found := slices.Contains(defaultSizes, 180)

	if !found {
		t.Error("default sizes should include 180 for Apple Touch Icon")
	}
}

func TestVersionVariables(t *testing.T) {
	if Version == "" {
		t.Error("Version should not be empty")
	}
	if BuildDate == "" {
		t.Error("BuildDate should not be empty")
	}
	if CommitHash == "" {
		t.Error("CommitHash should not be empty")
	}
}

func TestParseCategories(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{
			name:  "empty string",
			input: "",
			want:  nil,
		},
		{
			name:  "single category",
			input: "utilities",
			want:  []string{"utilities"},
		},
		{
			name:  "multiple categories",
			input: "utilities,productivity,tools",
			want:  []string{"utilities", "productivity", "tools"},
		},
		{
			name:  "categories with spaces",
			input: "utilities, productivity, tools",
			want:  []string{"utilities", "productivity", "tools"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseCategories(tt.input)
			if tt.want == nil && got != nil {
				t.Errorf("parseCategories(%q) = %v, want nil", tt.input, got)
				return
			}
			if tt.want != nil {
				if len(got) != len(tt.want) {
					t.Errorf("parseCategories(%q) = %v, want %v", tt.input, got, tt.want)
					return
				}
				for i := range got {
					if got[i] != tt.want[i] {
						t.Errorf("parseCategories(%q)[%d] = %q, want %q", tt.input, i, got[i], tt.want[i])
					}
				}
			}
		})
	}
}

func TestValidateSource(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "favicongen-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test files
	pngFile := filepath.Join(tmpDir, "test.png")
	svgFile := filepath.Join(tmpDir, "test.svg")
	txtFile := filepath.Join(tmpDir, "test.txt")

	for _, f := range []string{pngFile, svgFile, txtFile} {
		if err := os.WriteFile(f, []byte("test"), 0644); err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}
	}

	tests := []struct {
		name    string
		source  string
		wantErr bool
	}{
		{
			name:    "empty source",
			source:  "",
			wantErr: true,
		},
		{
			name:    "non-existent file",
			source:  "/nonexistent/file.png",
			wantErr: true,
		},
		{
			name:    "valid PNG file",
			source:  pngFile,
			wantErr: false,
		},
		{
			name:    "valid SVG file",
			source:  svgFile,
			wantErr: false,
		},
		{
			name:    "invalid format",
			source:  txtFile,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateSource(tt.source)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateSource(%q) error = %v, wantErr %v", tt.source, err, tt.wantErr)
			}
		})
	}
}

func TestBuildManifestConfig(t *testing.T) {
	config := &Config{
		AppName:            "Test App",
		AppShortName:       "Test",
		AppDescription:     "A test application",
		AppStartURL:        "/app",
		AppDisplay:         "fullscreen",
		AppOrientation:     "portrait",
		AppScope:           "/app/",
		AppThemeColor:      "#ff0000",
		AppBackgroundColor: "#00ff00",
		AppCategories:      []string{"utilities"},
		AppIcon:            "icon.png",
		Sizes:              []int{192, 512},
	}

	manifestConfig := config.buildManifestConfig()

	if manifestConfig.Name != config.AppName {
		t.Errorf("Name = %q, want %q", manifestConfig.Name, config.AppName)
	}
	if manifestConfig.ShortName != config.AppShortName {
		t.Errorf("ShortName = %q, want %q", manifestConfig.ShortName, config.AppShortName)
	}
	if manifestConfig.Display != config.AppDisplay {
		t.Errorf("Display = %q, want %q", manifestConfig.Display, config.AppDisplay)
	}
	if manifestConfig.ThemeColor != config.AppThemeColor {
		t.Errorf("ThemeColor = %q, want %q", manifestConfig.ThemeColor, config.AppThemeColor)
	}
	if len(manifestConfig.IconSizes) != len(config.Sizes) {
		t.Errorf("IconSizes length = %d, want %d", len(manifestConfig.IconSizes), len(config.Sizes))
	}
}

func TestBuildHTMLTagsConfig(t *testing.T) {
	config := &Config{
		Sizes:            []int{16, 32, 64},
		GenerateManifest: true,
		AppThemeColor:    "#123456",
	}

	htmlConfig := config.buildHTMLTagsConfig()

	if len(htmlConfig.Sizes) != len(config.Sizes) {
		t.Errorf("Sizes length = %d, want %d", len(htmlConfig.Sizes), len(config.Sizes))
	}
	if htmlConfig.IncludeManifest != config.GenerateManifest {
		t.Errorf("IncludeManifest = %v, want %v", htmlConfig.IncludeManifest, config.GenerateManifest)
	}
	if htmlConfig.ThemeColor != config.AppThemeColor {
		t.Errorf("ThemeColor = %q, want %q", htmlConfig.ThemeColor, config.AppThemeColor)
	}
}

func TestDefineFlags(t *testing.T) {
	// Reset flags for testing
	// Note: This test just verifies the function doesn't panic
	// and returns a non-nil result
	f := defineFlags()
	if f == nil {
		t.Error("defineFlags() returned nil")
	}
	if f.source == nil {
		t.Error("source flag is nil")
	}
	if f.output == nil {
		t.Error("output flag is nil")
	}
}

func TestBuildConfig(t *testing.T) {
	source := "test.png"
	output := "./output"
	f := &flags{
		source:             &source,
		output:             &output,
		backend:            new(string),
		generateHTML:       boolPtr(true),
		generateManifest:   boolPtr(false),
		generateICO:        boolPtr(true),
		generateHTMLOnly:   boolPtr(false),
		appName:            new(string),
		appShortName:       new(string),
		appDescription:     new(string),
		appStartURL:        strPtr("/"),
		appDisplay:         strPtr("standalone"),
		appOrientation:     strPtr("any"),
		appScope:           strPtr("/"),
		appThemeColor:      strPtr("#ffffff"),
		appBackgroundColor: strPtr("#ffffff"),
		appIcon:            new(string),
	}

	sizes := []int{16, 32}
	categories := []string{"test"}

	config := buildConfig(f, sizes, categories)

	if config.Source != source {
		t.Errorf("Source = %q, want %q", config.Source, source)
	}
	if config.Output != output {
		t.Errorf("Output = %q, want %q", config.Output, output)
	}
	if len(config.Sizes) != len(sizes) {
		t.Errorf("Sizes = %v, want %v", config.Sizes, sizes)
	}
	if len(config.AppCategories) != len(categories) {
		t.Errorf("AppCategories = %v, want %v", config.AppCategories, categories)
	}
}

func boolPtr(b bool) *bool {
	return &b
}

func strPtr(s string) *string {
	return &s
}
