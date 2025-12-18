package generator

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

// MockProcessor is a mock implementation of the Processor interface for testing
type MockProcessor struct {
	name        string
	available   bool
	resizeErr   error
	icoErr      error
	resizeCalls []resizeCall
	icoCalls    []icoCall
}

type resizeCall struct {
	inputPath  string
	outputPath string
	size       int
}

type icoCall struct {
	inputPaths []string
	outputPath string
}

func (m *MockProcessor) Name() string {
	return m.name
}

func (m *MockProcessor) IsAvailable() bool {
	return m.available
}

func (m *MockProcessor) Resize(inputPath, outputPath string, size int) error {
	m.resizeCalls = append(m.resizeCalls, resizeCall{inputPath, outputPath, size})
	if m.resizeErr != nil {
		return m.resizeErr
	}
	return os.WriteFile(outputPath, []byte("mock png data"), 0644)
}

func (m *MockProcessor) ConvertToICO(inputPaths []string, outputPath string) error {
	m.icoCalls = append(m.icoCalls, icoCall{inputPaths, outputPath})
	if m.icoErr != nil {
		return m.icoErr
	}
	return os.WriteFile(outputPath, []byte("mock ico data"), 0644)
}

func createTestEnv(t *testing.T) (tmpDir, sourcePath, outputDir string, cleanup func()) {
	t.Helper()
	tmpDir, err := os.MkdirTemp("", "favicongen-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	sourcePath = filepath.Join(tmpDir, "source.png")
	if err := os.WriteFile(sourcePath, []byte("mock source"), 0644); err != nil {
		os.RemoveAll(tmpDir)
		t.Fatalf("failed to create source file: %v", err)
	}

	outputDir = filepath.Join(tmpDir, "output")
	cleanup = func() { os.RemoveAll(tmpDir) }
	return
}

func newMockProcessor() *MockProcessor {
	return &MockProcessor{name: "mock", available: true}
}

func TestFaviconGeneratorGenerate(t *testing.T) {
	tmpDir, sourcePath, outputDir, cleanup := createTestEnv(t)
	defer cleanup()

	mockProc := newMockProcessor()
	gen := &FaviconGenerator{
		Processor:  mockProc,
		SourcePath: sourcePath,
		OutputDir:  outputDir,
		Sizes:      []int{16, 32, 48},
	}

	result, err := gen.Generate()
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	t.Run("creates output directory", func(t *testing.T) {
		if _, err := os.Stat(outputDir); os.IsNotExist(err) {
			t.Error("output directory was not created")
		}
	})

	t.Run("generates correct number of files", func(t *testing.T) {
		if len(result.GeneratedFiles) != 3 {
			t.Errorf("generated %d files, want 3", len(result.GeneratedFiles))
		}
	})

	t.Run("calls resize for each size", func(t *testing.T) {
		if len(mockProc.resizeCalls) != 3 {
			t.Errorf("resize called %d times, want 3", len(mockProc.resizeCalls))
		}
	})

	t.Run("creates nested output directory", func(t *testing.T) {
		nestedDir := filepath.Join(tmpDir, "deep", "nested", "output")
		gen2 := &FaviconGenerator{
			Processor:  newMockProcessor(),
			SourcePath: sourcePath,
			OutputDir:  nestedDir,
			Sizes:      []int{16},
		}
		if _, err := gen2.Generate(); err != nil {
			t.Fatalf("Generate() error = %v", err)
		}
		if _, err := os.Stat(nestedDir); os.IsNotExist(err) {
			t.Error("nested output directory was not created")
		}
	})
}

func TestFaviconGeneratorGenerateError(t *testing.T) {
	_, sourcePath, outputDir, cleanup := createTestEnv(t)
	defer cleanup()

	mockProc := &MockProcessor{
		name:      "mock",
		available: true,
		resizeErr: errors.New("mock resize error"),
	}

	gen := &FaviconGenerator{
		Processor:  mockProc,
		SourcePath: sourcePath,
		OutputDir:  outputDir,
		Sizes:      []int{16, 32},
	}

	if _, err := gen.Generate(); err == nil {
		t.Error("expected error from Generate()")
	}
}

func TestFaviconGeneratorGenerateICO(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "favicongen-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	mockProc := newMockProcessor()
	gen := &FaviconGenerator{
		Processor: mockProc,
		OutputDir: tmpDir,
		Sizes:     []int{16, 32, 48},
	}

	pngPaths := []string{
		filepath.Join(tmpDir, "favicon-16x16.png"),
		filepath.Join(tmpDir, "favicon-32x32.png"),
		filepath.Join(tmpDir, "favicon-48x48.png"),
	}

	icoPath, err := gen.GenerateICO(pngPaths)
	if err != nil {
		t.Fatalf("GenerateICO() error = %v", err)
	}

	expectedPath := filepath.Join(tmpDir, "favicon.ico")
	if icoPath != expectedPath {
		t.Errorf("GenerateICO() path = %q, want %q", icoPath, expectedPath)
	}

	if len(mockProc.icoCalls) != 1 {
		t.Fatalf("ICO called %d times, want 1", len(mockProc.icoCalls))
	}

	call := mockProc.icoCalls[0]
	if len(call.inputPaths) != 3 {
		t.Errorf("ICO inputPaths count = %d, want 3", len(call.inputPaths))
	}
}

func TestFaviconGeneratorGenerateICOError(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "favicongen-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	mockProc := &MockProcessor{
		name:      "mock",
		available: true,
		icoErr:    errors.New("mock ico error"),
	}

	gen := &FaviconGenerator{
		Processor: mockProc,
		OutputDir: tmpDir,
		Sizes:     []int{16, 32},
	}

	if _, err := gen.GenerateICO([]string{"a.png", "b.png"}); err == nil {
		t.Error("expected error from GenerateICO()")
	}
}

func TestGenerateResultFields(t *testing.T) {
	result := &GenerateResult{
		GeneratedFiles: []string{"a.png", "b.png"},
		ICOPath:        "favicon.ico",
	}

	if len(result.GeneratedFiles) != 2 {
		t.Errorf("GeneratedFiles count = %d, want 2", len(result.GeneratedFiles))
	}

	if result.ICOPath != "favicon.ico" {
		t.Errorf("ICOPath = %q, want %q", result.ICOPath, "favicon.ico")
	}
}
