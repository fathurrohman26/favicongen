package generator

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// ManifestConfig contains configuration for the web app manifest
type ManifestConfig struct {
	Name            string
	ShortName       string
	Description     string
	StartURL        string
	Display         string
	Orientation     string
	Scope           string
	ThemeColor      string
	BackgroundColor string
	Categories      []string
	IconPath        string
	IconSizes       []int
}

// Manifest represents a web app manifest
type Manifest struct {
	Name            string   `json:"name,omitempty"`
	ShortName       string   `json:"short_name,omitempty"`
	Description     string   `json:"description,omitempty"`
	StartURL        string   `json:"start_url"`
	Display         string   `json:"display"`
	Orientation     string   `json:"orientation,omitempty"`
	Scope           string   `json:"scope,omitempty"`
	ThemeColor      string   `json:"theme_color,omitempty"`
	BackgroundColor string   `json:"background_color,omitempty"`
	Categories      []string `json:"categories,omitempty"`
	Icons           []Icon   `json:"icons"`
}

// Icon represents an icon in the manifest
type Icon struct {
	Src     string `json:"src"`
	Sizes   string `json:"sizes"`
	Type    string `json:"type"`
	Purpose string `json:"purpose,omitempty"`
}

// GenerateManifest creates a manifest.webmanifest file
func GenerateManifest(config *ManifestConfig, outputDir string) (string, error) {
	manifest := &Manifest{
		Name:            config.Name,
		ShortName:       config.ShortName,
		Description:     config.Description,
		StartURL:        config.StartURL,
		Display:         config.Display,
		Orientation:     config.Orientation,
		Scope:           config.Scope,
		ThemeColor:      config.ThemeColor,
		BackgroundColor: config.BackgroundColor,
		Categories:      config.Categories,
		Icons:           make([]Icon, 0),
	}

	// Add icons
	for _, size := range config.IconSizes {
		icon := Icon{
			Src:   fmt.Sprintf("favicon-%dx%d.png", size, size),
			Sizes: fmt.Sprintf("%dx%d", size, size),
			Type:  "image/png",
		}

		// Mark larger icons as suitable for any purpose
		if size >= 192 {
			icon.Purpose = "any maskable"
		}

		manifest.Icons = append(manifest.Icons, icon)
	}

	// Marshal to JSON
	data, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal manifest: %w", err)
	}

	// Write to file
	manifestPath := filepath.Join(outputDir, "manifest.webmanifest")
	if err := os.WriteFile(manifestPath, data, 0644); err != nil {
		return "", fmt.Errorf("failed to write manifest file: %w", err)
	}

	return manifestPath, nil
}
