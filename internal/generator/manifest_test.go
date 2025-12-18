package generator

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func createManifestTestDir(t *testing.T) (string, func()) {
	t.Helper()
	tmpDir, err := os.MkdirTemp("", "favicongen-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	return tmpDir, func() { os.RemoveAll(tmpDir) }
}

func readManifest(t *testing.T, path string) map[string]interface{} {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read manifest: %v", err)
	}
	var manifest map[string]interface{}
	if err := json.Unmarshal(data, &manifest); err != nil {
		t.Fatalf("failed to parse manifest: %v", err)
	}
	return manifest
}

func TestGenerateManifestBasic(t *testing.T) {
	tmpDir, cleanup := createManifestTestDir(t)
	defer cleanup()

	config := &ManifestConfig{
		Name:            "Test App",
		ShortName:       "Test",
		Description:     "A test application",
		StartURL:        "/",
		Display:         "standalone",
		Orientation:     "any",
		Scope:           "/",
		ThemeColor:      "#ffffff",
		BackgroundColor: "#000000",
		Categories:      []string{"utilities", "productivity"},
		IconSizes:       []int{192, 512},
	}

	manifestPath, err := GenerateManifest(config, tmpDir)
	if err != nil {
		t.Fatalf("GenerateManifest() error = %v", err)
	}

	expectedPath := filepath.Join(tmpDir, "manifest.webmanifest")
	if manifestPath != expectedPath {
		t.Errorf("manifest path = %q, want %q", manifestPath, expectedPath)
	}

	manifest := readManifest(t, manifestPath)

	if manifest["name"] != "Test App" {
		t.Errorf("name = %v, want %q", manifest["name"], "Test App")
	}
	if manifest["display"] != "standalone" {
		t.Errorf("display = %v, want %q", manifest["display"], "standalone")
	}
}

func TestGenerateManifestMinimal(t *testing.T) {
	tmpDir, cleanup := createManifestTestDir(t)
	defer cleanup()

	config := &ManifestConfig{
		StartURL:  "/",
		Display:   "standalone",
		IconSizes: []int{192},
	}

	manifestPath, err := GenerateManifest(config, tmpDir)
	if err != nil {
		t.Fatalf("GenerateManifest() error = %v", err)
	}

	manifest := readManifest(t, manifestPath)

	if manifest["start_url"] != "/" {
		t.Errorf("start_url = %v, want %q", manifest["start_url"], "/")
	}
	if manifest["display"] != "standalone" {
		t.Errorf("display = %v, want %q", manifest["display"], "standalone")
	}
}

func TestGenerateManifestIcons(t *testing.T) {
	tmpDir, cleanup := createManifestTestDir(t)
	defer cleanup()

	config := &ManifestConfig{
		StartURL:  "/",
		Display:   "standalone",
		IconSizes: []int{48, 192, 512},
	}

	manifestPath, err := GenerateManifest(config, tmpDir)
	if err != nil {
		t.Fatalf("GenerateManifest() error = %v", err)
	}

	data, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatalf("failed to read manifest: %v", err)
	}

	var manifest Manifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		t.Fatalf("failed to parse manifest: %v", err)
	}

	if len(manifest.Icons) != 3 {
		t.Fatalf("got %d icons, want 3", len(manifest.Icons))
	}

	// Check first icon (48px, no purpose)
	if manifest.Icons[0].Src != "favicon-48x48.png" {
		t.Errorf("icon[0].Src = %q, want %q", manifest.Icons[0].Src, "favicon-48x48.png")
	}
	if manifest.Icons[0].Purpose != "" {
		t.Errorf("icon[0].Purpose = %q, want empty", manifest.Icons[0].Purpose)
	}

	// Check large icon (192px, has purpose)
	if manifest.Icons[1].Purpose != "any maskable" {
		t.Errorf("icon[1].Purpose = %q, want %q", manifest.Icons[1].Purpose, "any maskable")
	}
}

func TestGenerateManifestInvalidDir(t *testing.T) {
	config := &ManifestConfig{
		StartURL:  "/",
		Display:   "standalone",
		IconSizes: []int{192},
	}

	_, err := GenerateManifest(config, "/nonexistent/path/that/cannot/exist")
	if err == nil {
		t.Error("expected error for invalid directory")
	}
}

func TestManifestCategories(t *testing.T) {
	tmpDir, cleanup := createManifestTestDir(t)
	defer cleanup()

	config := &ManifestConfig{
		StartURL:   "/",
		Display:    "standalone",
		Categories: []string{"games", "entertainment", "social"},
		IconSizes:  []int{192},
	}

	manifestPath, err := GenerateManifest(config, tmpDir)
	if err != nil {
		t.Fatalf("GenerateManifest() error = %v", err)
	}

	manifest := readManifest(t, manifestPath)

	categories, ok := manifest["categories"].([]interface{})
	if !ok {
		t.Fatal("manifest missing categories array")
	}

	if len(categories) != 3 {
		t.Errorf("categories count = %d, want 3", len(categories))
	}
}
