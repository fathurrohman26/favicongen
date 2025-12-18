# favicongen

> A lightweight, cross-platform CLI tool for generating complete favicon sets from a single source image.

**favicongen** streamlines favicon creation by automatically generating multiple sizes, multi-resolution `favicon.ico` files, and `manifest.webmanifest` files for progressive web applications. It supports both SVG and PNG source images and leverages ImageMagick or libvips for high-quality image processing.

## Features

- **Fast & Lightweight** - Optimized for speed with minimal dependencies
- **Multiple Formats** - Supports SVG and PNG source images
- **Smart Processing** - Automatically detects and uses ImageMagick or libvips
- **Complete Output** - Generates all standard favicon sizes from a single source
- **Web App Ready** - Creates multi-resolution `favicon.ico` and `manifest.webmanifest` files
- **HTML Integration** - Automatically generates ready-to-use HTML `<link>` tags
- **Highly Customizable** - Configure sizes, output directory, and manifest properties
- **Cross-Platform** - Works seamlessly on Windows, macOS, and Linux
- **Build Integration** - CLI interface perfect for CI/CD pipelines and build scripts

## Installation

### Using Go

Install directly using the Go toolchain:

```bash
go install github.com/fathurrohman26/favicongen/cmd/favicongen@latest
```

### Build from Source

Clone and build manually:

```bash
git clone https://github.com/fathurrohman26/favicongen.git
cd favicongen
make build
```

### Prerequisites

Ensure you have one of the following image processing libraries installed:

- **ImageMagick**

You can install ImageMagick via package managers:

- On macOS (using Homebrew):

  ```bash
  brew install imagemagick
  ```

- For other platforms, refer to the [ImageMagick installation guide](https://imagemagick.org/script/download.php).

- **libvips**

Install libvips via package managers:

- On macOS (using Homebrew):

  ```bash
  brew install vips
  ```

- For other platforms, refer to the [libvips installation guide](https://libvips.github.io/libvips/install.html).

## Usage

### Basic Usage

```bash
# Full command with flags
favicongen --source path/to/logo.svg --output ./public/favicons

# Shorthand (positional arguments)
favicongen path/to/logo.png ./output
```

### Quick Commands

```bash
# Display help
favicongen help

# Show version
favicongen version
```

### Command-Line Options

#### General Options

| Option | Description | Default |
|--------|-------------|---------|
| `--source` | Path to the source image file (SVG or PNG). | N/A |
| `--output` | Path to the output directory where favicon files will be saved. | `./favicons` |
| `--sizes` | Comma-separated list of sizes to generate (includes 180 for Apple Touch Icon). | `16,32,48,64,128,180,256,512` |
| `--backend` | Image processing backend to use (`imagemagick` or `vips`). If not specified, favicongen will auto-detect. | N/A |
| `--html-tags` | Generate HTML tags for the favicons. | True |
| `--manifest` | Generate a `manifest.webmanifest` file. | False |
| `--ico` | Generate a multi-resolution `favicon.ico` file. | True |

#### Manifest Configuration

| Option | Description | Default |
|--------|-------------|---------|
| `--app-name` | Name of the application. | N/A |
| `--app-short-name` | Short name of the application. | N/A |
| `--app-description` | Description of the application. | N/A |
| `--app-start-url` | Start URL of the application. | `/` |
| `--app-display` | Display mode of the application (e.g., `standalone`, `fullscreen`). | `standalone` |
| `--app-orientation` | Default orientation of the application (e.g., `portrait`, `landscape`). | `any` |
| `--app-scope` | Scope of the application. | `/` |
| `--app-theme-color` | Theme color of the application. | `#ffffff` |
| `--app-background-color` | Background color of the application. | `#ffffff` |
| `--app-categories` | Comma-separated list of categories for the application. | N/A |
| `--app-icon` | Path to the application icon file (should be one of the generated favicons). | N/A |

### Examples

#### Specify Image Processing Backend

```bash
# Use ImageMagick
favicongen --source logo.svg --output ./public/favicons --sizes 16,32,48,64 --backend imagemagick

# Use libvips
favicongen --source logo.png --output ./public/favicons --sizes 16,32,48,64 --backend vips
```

#### Custom Output Options

```bash
# Generate only favicon images (skip HTML tags file)
favicongen --source logo.svg --output ./public/favicons --html-tags=false

# Generate favicons with manifest file
favicongen --source logo.svg --output ./public/favicons --manifest --app-name "My App"
```

#### Generate HTML Tags from Existing Favicons

```bash
# Generate complete manifest and HTML tags (outputs to stdout)
favicongen --generate-html-tags \
  --output ./public/favicons \
  --sizes 16,32,48,64 \
  --manifest \
  --app-name "My App" \
  --app-short-name "App" \
  --app-description "My Application" \
  --app-theme-color "#ffffff" \
  --app-background-color "#ffffff" \
  --app-start-url "/" \
  --app-display "standalone" \
  --app-orientation "portrait" \
  --app-scope "/" \
  --app-icon "favicon-512x512.png" \
  --app-categories "utilities,productivity"
```

## Contributing

Contributions are welcome! Here's how you can help:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

Please ensure your code follows the existing style and includes appropriate tests.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Support

If you encounter any issues or have questions:

- Open an issue on [GitHub](https://github.com/fathurrohman26/favicongen/issues)
- Check existing issues for solutions
