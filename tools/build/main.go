package main

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"crypto/sha256"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

var (
	Version    = "dev"
	BuildDate  = "unknown"
	CommitHash = "unknown"
)

var (
	outputDir    string
	targetOS     string
	targetArch   string
	buildAll     bool
	ldflags      string
	checksumFile string
	checksums    []checksumEntry
)

type checksumEntry struct {
	filename string
	checksum string
}

func init() {
	// Parse command-line flags
	flag.StringVar(&outputDir, "o", "dist", "Output directory for built binaries")
	flag.StringVar(&targetOS, "os", runtime.GOOS, "Target operating system (linux, darwin, windows)")
	flag.StringVar(&targetArch, "arch", runtime.GOARCH, "Target architecture (amd64, arm64)")
	flag.BoolVar(&buildAll, "all", false, "Build for all supported platforms")

	// Get build metadata from environment or defaults
	Version = getEnv("VERSION", getLastTag(Version))
	BuildDate = getEnv("BUILD_DATE", time.Now().UTC().Format("2006-01-02T15:04:05Z"))
	CommitHash = getEnv("COMMIT_HASH", getCommitHash())
}

func main() {
	flag.Parse()

	// Build ldflags with version information
	ldflags = fmt.Sprintf("-s -w -X main.Version=%s -X main.BuildDate=%s -X main.CommitHash=%s",
		Version, BuildDate, CommitHash)

	checksumFile = filepath.Join(outputDir, "checksums.txt")

	if buildAll {
		buildAllPlatforms()
	} else {
		build(targetOS, targetArch)
	}

	// Write checksums file
	if err := writeChecksums(); err != nil {
		log.Printf("Warning: failed to write checksums: %v", err)
	}
}

func buildAllPlatforms() {
	platforms := []struct {
		os   string
		arch string
	}{
		{"linux", "amd64"},
		{"linux", "arm64"},
		{"darwin", "amd64"},
		{"darwin", "arm64"},
		{"windows", "amd64"},
		{"windows", "arm64"},
	}

	log.Println("Building for all platforms...")
	for _, p := range platforms {
		if err := build(p.os, p.arch); err != nil {
			log.Printf("Failed to build for %s/%s: %v", p.os, p.arch, err)
		}
	}
}

func build(goos, goarch string) error {
	// Create output directory
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Determine binary name
	binaryName := "favicongen"
	if goos == "windows" {
		binaryName += ".exe"
	}

	// Create platform-specific subdirectory
	platformDir := filepath.Join(outputDir, fmt.Sprintf("%s-%s", goos, goarch))
	if err := os.MkdirAll(platformDir, 0755); err != nil {
		return fmt.Errorf("failed to create platform directory: %w", err)
	}

	outputPath := filepath.Join(platformDir, binaryName)

	log.Printf("Building for %s/%s...", goos, goarch)
	log.Printf("  Version:    %s", Version)
	log.Printf("  BuildDate:  %s", BuildDate)
	log.Printf("  CommitHash: %s", CommitHash)
	log.Printf("  Output:     %s", outputPath)

	// Build command
	cmd := exec.Command("go", "build",
		"-ldflags", ldflags,
		"-o", outputPath,
		"./cmd/favicongen",
	)

	// Set environment variables for cross-compilation
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("GOOS=%s", goos),
		fmt.Sprintf("GOARCH=%s", goarch),
		"CGO_ENABLED=0", // Disable CGO for better cross-platform compatibility
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("build failed: %w", err)
	}

	// Get file size
	info, err := os.Stat(outputPath)
	if err == nil {
		log.Printf("  Size:       %.2f MB", float64(info.Size())/(1024*1024))
	}

	log.Printf("✓ Successfully built %s/%s\n", goos, goarch)

	// Create archive
	archivePath, err := createArchive(goos, goarch, platformDir, binaryName)
	if err != nil {
		return fmt.Errorf("failed to create archive: %w", err)
	}

	// Calculate checksum
	checksum, err := calculateSHA256(archivePath)
	if err != nil {
		log.Printf("Warning: failed to calculate checksum for %s: %v", archivePath, err)
	} else {
		checksums = append(checksums, checksumEntry{
			filename: filepath.Base(archivePath),
			checksum: checksum,
		})
		log.Printf("  Checksum:   %s", checksum)
	}

	log.Printf("  Archive:    %s\n", archivePath)

	return nil
}

func getEnv(key, defaultValue string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return defaultValue
}

func getLastTag(fallback string) string {
	cmd := exec.Command("git", "describe", "--tags", "--abbrev=0")
	output, err := cmd.Output()
	if err != nil {
		return fallback
	}
	return strings.TrimSpace(string(output))
}

func getCommitHash() string {
	cmd := exec.Command("git", "rev-parse", "--short", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "unknown"
	}
	return strings.TrimSpace(string(output))
}

func createArchive(goos, goarch, sourceDir, binaryName string) (string, error) {
	// Determine archive name and format
	archiveName := fmt.Sprintf("favicongen_%s_%s_%s", Version, goos, goarch)
	var archivePath string

	if goos == "windows" {
		// Create .zip for Windows
		archivePath = filepath.Join(outputDir, archiveName+".zip")
		return archivePath, createZipArchive(archivePath, sourceDir, binaryName)
	}

	// Create .tar.gz for Unix-like systems
	archivePath = filepath.Join(outputDir, archiveName+".tar.gz")
	return archivePath, createTarGzArchive(archivePath, sourceDir, binaryName)
}

func createTarGzArchive(archivePath, sourceDir, binaryName string) error {
	// Create archive file
	file, err := os.Create(archivePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create gzip writer
	gzipWriter := gzip.NewWriter(file)
	defer gzipWriter.Close()

	// Create tar writer
	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	// Add binary to archive
	binaryPath := filepath.Join(sourceDir, binaryName)
	return addFileToTar(tarWriter, binaryPath, binaryName)
}

func addFileToTar(tarWriter *tar.Writer, filePath, nameInArchive string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return err
	}

	header, err := tar.FileInfoHeader(info, "")
	if err != nil {
		return err
	}
	header.Name = nameInArchive

	if err := tarWriter.WriteHeader(header); err != nil {
		return err
	}

	_, err = io.Copy(tarWriter, file)
	return err
}

func createZipArchive(archivePath, sourceDir, binaryName string) error {
	// Create archive file
	file, err := os.Create(archivePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create zip writer
	zipWriter := zip.NewWriter(file)
	defer zipWriter.Close()

	// Add binary to archive
	binaryPath := filepath.Join(sourceDir, binaryName)
	return addFileToZip(zipWriter, binaryPath, binaryName)
}

func addFileToZip(zipWriter *zip.Writer, filePath, nameInArchive string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}
	header.Name = nameInArchive
	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, file)
	return err
}

func calculateSHA256(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

func writeChecksums() error {
	if len(checksums) == 0 {
		return nil
	}

	file, err := os.Create(checksumFile)
	if err != nil {
		return err
	}
	defer file.Close()

	log.Printf("\nWriting checksums to %s", checksumFile)
	for _, entry := range checksums {
		line := fmt.Sprintf("%s  %s\n", entry.checksum, entry.filename)
		if _, err := file.WriteString(line); err != nil {
			return err
		}
		log.Printf("  %s: %s", entry.filename, entry.checksum)
	}

	return nil
}
