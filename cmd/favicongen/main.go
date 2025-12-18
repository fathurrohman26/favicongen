package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/fathurrohman26/favicongen/internal/generator"
	"github.com/fathurrohman26/favicongen/internal/processor"
)

var (
	Version    = "dev"
	BuildDate  = "unknown"
	CommitHash = "unknown"
)

type Config struct {
	Source             string
	Output             string
	Sizes              []int
	Backend            string
	GenerateHTML       bool
	GenerateManifest   bool
	GenerateICO        bool
	GenerateHTMLOnly   bool
	AppName            string
	AppShortName       string
	AppDescription     string
	AppStartURL        string
	AppDisplay         string
	AppOrientation     string
	AppScope           string
	AppThemeColor      string
	AppBackgroundColor string
	AppCategories      []string
	AppIcon            string
}

type flags struct {
	source             *string
	output             *string
	sizesStr           *string
	backend            *string
	generateHTML       *bool
	generateManifest   *bool
	generateICO        *bool
	generateHTMLOnly   *bool
	appName            *string
	appShortName       *string
	appDescription     *string
	appStartURL        *string
	appDisplay         *string
	appOrientation     *string
	appScope           *string
	appThemeColor      *string
	appBackgroundColor *string
	appCategories      *string
	appIcon            *string
	showVersion        *bool
	showHelp           *bool
}

func defineFlags() *flags {
	return &flags{
		source:             flag.String("source", "", "Path to source image (SVG or PNG)"),
		output:             flag.String("output", "./favicons", "Output directory for generated files"),
		sizesStr:           flag.String("sizes", "16,32,48,64,128,180,256,512", "Comma-separated list of sizes"),
		backend:            flag.String("backend", "", "Image processor backend (imagemagick or vips)"),
		generateHTML:       flag.Bool("html-tags", true, "Generate HTML link tags"),
		generateManifest:   flag.Bool("manifest", false, "Generate manifest.webmanifest file"),
		generateICO:        flag.Bool("ico", true, "Generate favicon.ico file"),
		generateHTMLOnly:   flag.Bool("generate-html-tags", false, "Only generate HTML tags from existing favicons"),
		appName:            flag.String("app-name", "", "Application name for manifest"),
		appShortName:       flag.String("app-short-name", "", "Short application name for manifest"),
		appDescription:     flag.String("app-description", "", "Application description for manifest"),
		appStartURL:        flag.String("app-start-url", "/", "Start URL for manifest"),
		appDisplay:         flag.String("app-display", "standalone", "Display mode for manifest"),
		appOrientation:     flag.String("app-orientation", "any", "Orientation for manifest"),
		appScope:           flag.String("app-scope", "/", "Scope for manifest"),
		appThemeColor:      flag.String("app-theme-color", "#ffffff", "Theme color for manifest"),
		appBackgroundColor: flag.String("app-background-color", "#ffffff", "Background color for manifest"),
		appCategories:      flag.String("app-categories", "", "Comma-separated categories for manifest"),
		appIcon:            flag.String("app-icon", "", "Icon path for manifest"),
		showVersion:        flag.Bool("version", false, "Show version information"),
		showHelp:           flag.Bool("help", false, "Show help information"),
	}
}

func shouldShowVersion(f *flags) bool {
	return *f.showVersion || (len(os.Args) > 1 && os.Args[1] == "version")
}

func shouldShowHelp(f *flags) bool {
	return *f.showHelp || len(os.Args) == 1 || (len(os.Args) > 1 && os.Args[1] == "help")
}

func printVersion() {
	fmt.Printf("favicongen version %s\n", Version)
	fmt.Printf("Build date: %s\n", BuildDate)
	fmt.Printf("Commit: %s\n", CommitHash)
}

func parsePositionalArgs(f *flags) {
	args := flag.Args()
	if len(args) >= 1 && *f.source == "" {
		*f.source = args[0]
	}
	if len(args) >= 2 && *f.output == "./favicons" {
		*f.output = args[1]
	}
}

func parseCategories(categoriesStr string) []string {
	if categoriesStr == "" {
		return nil
	}
	categories := strings.Split(categoriesStr, ",")
	for i := range categories {
		categories[i] = strings.TrimSpace(categories[i])
	}
	return categories
}

func buildConfig(f *flags, sizes []int, categories []string) *Config {
	return &Config{
		Source:             *f.source,
		Output:             *f.output,
		Sizes:              sizes,
		Backend:            *f.backend,
		GenerateHTML:       *f.generateHTML,
		GenerateManifest:   *f.generateManifest,
		GenerateICO:        *f.generateICO,
		GenerateHTMLOnly:   *f.generateHTMLOnly,
		AppName:            *f.appName,
		AppShortName:       *f.appShortName,
		AppDescription:     *f.appDescription,
		AppStartURL:        *f.appStartURL,
		AppDisplay:         *f.appDisplay,
		AppOrientation:     *f.appOrientation,
		AppScope:           *f.appScope,
		AppThemeColor:      *f.appThemeColor,
		AppBackgroundColor: *f.appBackgroundColor,
		AppCategories:      categories,
		AppIcon:            *f.appIcon,
	}
}

func main() {
	f := defineFlags()
	flag.Parse()

	if shouldShowVersion(f) {
		printVersion()
		return
	}

	if shouldShowHelp(f) {
		showUsage()
		return
	}

	parsePositionalArgs(f)

	sizes, err := parseSizes(*f.sizesStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: invalid sizes format: %v\n", err)
		os.Exit(1)
	}

	categories := parseCategories(*f.appCategories)
	config := buildConfig(f, sizes, categories)

	if err := run(config); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func (c *Config) buildManifestConfig() *generator.ManifestConfig {
	return &generator.ManifestConfig{
		Name:            c.AppName,
		ShortName:       c.AppShortName,
		Description:     c.AppDescription,
		StartURL:        c.AppStartURL,
		Display:         c.AppDisplay,
		Orientation:     c.AppOrientation,
		Scope:           c.AppScope,
		ThemeColor:      c.AppThemeColor,
		BackgroundColor: c.AppBackgroundColor,
		Categories:      c.AppCategories,
		IconPath:        c.AppIcon,
		IconSizes:       c.Sizes,
	}
}

func (c *Config) buildHTMLTagsConfig() *generator.HTMLTagsConfig {
	return &generator.HTMLTagsConfig{
		Sizes:           c.Sizes,
		IncludeManifest: c.GenerateManifest,
		ThemeColor:      c.AppThemeColor,
	}
}

func runHTMLOnlyMode(config *Config) error {
	htmlTags := generator.GenerateHTMLTags(config.buildHTMLTagsConfig())
	fmt.Println(htmlTags)

	if config.GenerateManifest {
		manifestPath, err := generator.GenerateManifest(config.buildManifestConfig(), config.Output)
		if err != nil {
			return fmt.Errorf("failed to generate manifest: %w", err)
		}
		fmt.Printf("\n✓ Generated manifest: %s\n", manifestPath)
	}

	return nil
}

func validateSource(source string) error {
	if source == "" {
		return fmt.Errorf("source image is required (use --source or provide as first argument)")
	}

	if _, err := os.Stat(source); os.IsNotExist(err) {
		return fmt.Errorf("source file does not exist: %s", source)
	}

	ext := strings.ToLower(filepath.Ext(source))
	if ext != ".svg" && ext != ".png" {
		return fmt.Errorf("source must be SVG or PNG format, got: %s", ext)
	}

	return nil
}

func generateICOFile(gen *generator.FaviconGenerator, config *Config) {
	var icoSizes []string
	for _, size := range []int{16, 32, 48} {
		path := filepath.Join(config.Output, fmt.Sprintf("favicon-%dx%d.png", size, size))
		if fileExists(path) {
			icoSizes = append(icoSizes, path)
		}
	}

	if len(icoSizes) == 0 {
		return
	}

	icoPath, err := gen.GenerateICO(icoSizes)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to generate ICO file: %v\n", err)
		return
	}
	fmt.Printf("✓ Generated favicon.ico: %s\n", icoPath)
}

func generateHTMLTagsFile(config *Config) error {
	htmlTags := generator.GenerateHTMLTags(config.buildHTMLTagsConfig())

	htmlPath := filepath.Join(config.Output, "favicon-tags.html")
	if err := os.WriteFile(htmlPath, []byte(htmlTags), 0644); err != nil {
		return fmt.Errorf("failed to write HTML tags: %w", err)
	}

	fmt.Printf("✓ Generated HTML tags: %s\n", htmlPath)
	fmt.Println("\nHTML tags to include in your <head>:")
	fmt.Println(strings.Repeat("-", 50))
	fmt.Println(htmlTags)
	fmt.Println(strings.Repeat("-", 50))

	return nil
}

func run(config *Config) error {
	if config.GenerateHTMLOnly {
		return runHTMLOnlyMode(config)
	}

	if err := validateSource(config.Source); err != nil {
		return err
	}

	proc, err := processor.DetectAvailableProcessor(config.Backend)
	if err != nil {
		return fmt.Errorf("failed to initialize image processor: %w", err)
	}

	fmt.Printf("Using %s for image processing\n", proc.Name())
	fmt.Printf("Source: %s\n", config.Source)
	fmt.Printf("Output: %s\n", config.Output)
	fmt.Printf("Sizes: %v\n", config.Sizes)

	gen := &generator.FaviconGenerator{
		Processor:  proc,
		SourcePath: config.Source,
		OutputDir:  config.Output,
		Sizes:      config.Sizes,
	}

	result, err := gen.Generate()
	if err != nil {
		return fmt.Errorf("failed to generate favicons: %w", err)
	}

	fmt.Printf("\n✓ Generated %d favicon files\n", len(result.GeneratedFiles))

	if config.GenerateICO {
		generateICOFile(gen, config)
	}

	if config.GenerateManifest {
		manifestPath, err := generator.GenerateManifest(config.buildManifestConfig(), config.Output)
		if err != nil {
			return fmt.Errorf("failed to generate manifest: %w", err)
		}
		fmt.Printf("✓ Generated manifest: %s\n", manifestPath)
	}

	if config.GenerateHTML {
		if err := generateHTMLTagsFile(config); err != nil {
			return err
		}
	}

	return nil
}

func parseSizes(sizesStr string) ([]int, error) {
	parts := strings.Split(sizesStr, ",")
	sizes := make([]int, 0, len(parts))

	for _, part := range parts {
		part = strings.TrimSpace(part)
		size, err := strconv.Atoi(part)
		if err != nil {
			return nil, fmt.Errorf("invalid size: %s", part)
		}
		if size <= 0 {
			return nil, fmt.Errorf("size must be positive: %d", size)
		}
		sizes = append(sizes, size)
	}

	return sizes, nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func showUsage() {
	fmt.Printf("favicongen v%s - Generate favicon files from a single image\n\n", Version)
	fmt.Println("Usage:")
	fmt.Println("  favicongen --source <image> --output <dir> [options]")
	fmt.Println("  favicongen <image> <dir>  (shorthand)")
	fmt.Println("  favicongen version        (show version)")
	fmt.Println("  favicongen help           (show this help)")
	fmt.Println()
	fmt.Println("General Options:")
	fmt.Println("  --source <path>          Source image file (SVG or PNG)")
	fmt.Println("  --output <dir>           Output directory (default: ./favicons)")
	fmt.Println("  --sizes <sizes>          Comma-separated sizes (default: 16,32,48,64,128,180,256,512)")
	fmt.Println("  --backend <name>         Image processor: imagemagick or vips (auto-detect if not specified)")
	fmt.Println("  --html-tags              Generate HTML tags file (default: true)")
	fmt.Println("  --manifest               Generate manifest.webmanifest (default: false)")
	fmt.Println("  --ico                    Generate favicon.ico (default: true)")
	fmt.Println("  --generate-html-tags     Only generate HTML tags from existing favicons")
	fmt.Println()
	fmt.Println("Manifest Options:")
	fmt.Println("  --app-name <name>               Application name")
	fmt.Println("  --app-short-name <name>         Short application name")
	fmt.Println("  --app-description <text>        Application description")
	fmt.Println("  --app-start-url <url>           Start URL (default: /)")
	fmt.Println("  --app-display <mode>            Display mode (default: standalone)")
	fmt.Println("  --app-orientation <mode>        Orientation (default: any)")
	fmt.Println("  --app-scope <path>              Scope (default: /)")
	fmt.Println("  --app-theme-color <color>       Theme color (default: #ffffff)")
	fmt.Println("  --app-background-color <color>  Background color (default: #ffffff)")
	fmt.Println("  --app-categories <list>         Comma-separated categories")
	fmt.Println("  --app-icon <path>               Icon path")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  favicongen logo.svg ./public/favicons")
	fmt.Println("  favicongen --source logo.png --output ./dist --sizes 16,32,64")
	fmt.Println("  favicongen --source logo.svg --manifest --app-name \"My App\"")
	fmt.Println("  favicongen --generate-html-tags --output ./public --sizes 16,32,64")
}
