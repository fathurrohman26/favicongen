package generator

import (
	"strings"
	"testing"
)

func TestGenerateHTMLTags(t *testing.T) {
	tests := []struct {
		name           string
		config         *HTMLTagsConfig
		wantContains   []string
		wantNotContain []string
	}{
		{
			name: "basic sizes",
			config: &HTMLTagsConfig{
				Sizes:           []int{16, 32},
				IncludeManifest: false,
				ThemeColor:      "",
			},
			wantContains: []string{
				`<link rel="icon" href="/favicon.ico" sizes="any">`,
				`<link rel="icon" type="image/png" sizes="16x16" href="/favicon-16x16.png">`,
				`<link rel="icon" type="image/png" sizes="32x32" href="/favicon-32x32.png">`,
			},
			wantNotContain: []string{
				`manifest`,
				`theme-color`,
			},
		},
		{
			name: "with manifest",
			config: &HTMLTagsConfig{
				Sizes:           []int{16, 32},
				IncludeManifest: true,
				ThemeColor:      "",
			},
			wantContains: []string{
				`<link rel="manifest" href="/manifest.webmanifest">`,
			},
		},
		{
			name: "with theme color",
			config: &HTMLTagsConfig{
				Sizes:           []int{16},
				IncludeManifest: false,
				ThemeColor:      "#ff0000",
			},
			wantContains: []string{
				`<meta name="theme-color" content="#ff0000">`,
			},
		},
		{
			name: "with apple touch icon size 180",
			config: &HTMLTagsConfig{
				Sizes:           []int{16, 32, 180},
				IncludeManifest: false,
				ThemeColor:      "",
			},
			wantContains: []string{
				`<link rel="apple-touch-icon" sizes="180x180" href="/favicon-180x180.png">`,
			},
		},
		{
			name: "with apple touch icon fallback (size >= 180)",
			config: &HTMLTagsConfig{
				Sizes:           []int{16, 32, 192},
				IncludeManifest: false,
				ThemeColor:      "",
			},
			wantContains: []string{
				`<link rel="apple-touch-icon" sizes="192x192" href="/favicon-192x192.png">`,
			},
		},
		{
			name: "no apple touch icon when sizes < 180",
			config: &HTMLTagsConfig{
				Sizes:           []int{16, 32, 64},
				IncludeManifest: false,
				ThemeColor:      "",
			},
			wantNotContain: []string{
				`apple-touch-icon`,
			},
		},
		{
			name: "full configuration",
			config: &HTMLTagsConfig{
				Sizes:           []int{16, 32, 64, 180, 512},
				IncludeManifest: true,
				ThemeColor:      "#ffffff",
			},
			wantContains: []string{
				`<link rel="icon" href="/favicon.ico" sizes="any">`,
				`<link rel="icon" type="image/png" sizes="16x16" href="/favicon-16x16.png">`,
				`<link rel="icon" type="image/png" sizes="32x32" href="/favicon-32x32.png">`,
				`<link rel="icon" type="image/png" sizes="64x64" href="/favicon-64x64.png">`,
				`<link rel="icon" type="image/png" sizes="180x180" href="/favicon-180x180.png">`,
				`<link rel="icon" type="image/png" sizes="512x512" href="/favicon-512x512.png">`,
				`<link rel="apple-touch-icon" sizes="180x180" href="/favicon-180x180.png">`,
				`<link rel="manifest" href="/manifest.webmanifest">`,
				`<meta name="theme-color" content="#ffffff">`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateHTMLTags(tt.config)

			for _, want := range tt.wantContains {
				if !strings.Contains(got, want) {
					t.Errorf("GenerateHTMLTags() missing expected content:\nwant: %s\ngot:\n%s", want, got)
				}
			}

			for _, notWant := range tt.wantNotContain {
				if strings.Contains(got, notWant) {
					t.Errorf("GenerateHTMLTags() contains unexpected content:\nnot want: %s\ngot:\n%s", notWant, got)
				}
			}
		})
	}
}

func TestGenerateHTMLTagsEmptySizes(t *testing.T) {
	config := &HTMLTagsConfig{
		Sizes:           []int{},
		IncludeManifest: false,
		ThemeColor:      "",
	}

	got := GenerateHTMLTags(config)

	// Should still have favicon.ico link
	if !strings.Contains(got, `<link rel="icon" href="/favicon.ico" sizes="any">`) {
		t.Error("expected favicon.ico link even with empty sizes")
	}
}
