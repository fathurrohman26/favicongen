package generator

import (
	"fmt"
	"strings"
)

// HTMLTagsConfig contains configuration for HTML tag generation
type HTMLTagsConfig struct {
	Sizes           []int
	IncludeManifest bool
	ThemeColor      string
}

// GenerateHTMLTags creates HTML link tags for favicons
func GenerateHTMLTags(config *HTMLTagsConfig) string {
	var tags []string

	// Add favicon.ico link (default browser favicon)
	tags = append(tags, `<link rel="icon" href="/favicon.ico" sizes="any">`)

	// Add PNG favicons for each size
	for _, size := range config.Sizes {
		tag := fmt.Sprintf(`<link rel="icon" type="image/png" sizes="%dx%d" href="/favicon-%dx%d.png">`,
			size, size, size, size)
		tags = append(tags, tag)
	}

	// Add Apple Touch Icon (typically 180x180)
	for _, size := range config.Sizes {
		if size == 180 || size >= 180 {
			tag := fmt.Sprintf(`<link rel="apple-touch-icon" sizes="%dx%d" href="/favicon-%dx%d.png">`,
				size, size, size, size)
			tags = append(tags, tag)
			break
		}
	}

	// Add manifest link
	if config.IncludeManifest {
		tags = append(tags, `<link rel="manifest" href="/manifest.webmanifest">`)
	}

	// Add theme color meta tag
	if config.ThemeColor != "" {
		tags = append(tags, fmt.Sprintf(`<meta name="theme-color" content="%s">`, config.ThemeColor))
	}

	return strings.Join(tags, "\n")
}
