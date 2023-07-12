package main

import (
	"fmt"
	"testing"
)

func TestBuildpackBuilderImage(t *testing.T) {
	testCases := []struct {
		name           string
		runtime        string
		tag            string
		wantError      bool
		wantBuilderURL string
	}{
		{
			name:           "Extracts go from go119",
			runtime:        "go119",
			tag:            "1.2.3",
			wantBuilderURL: fmt.Sprintf(defaultBuilderURLTemplate, "go", "1.2.3"),
		},
		{
			name:      "Fails with incorrect runtime format",
			runtime:   "11go",
			tag:       "latest",
			wantError: true,
		},
		{
			name:           "Extracts php from php82",
			runtime:        "php82",
			tag:            "latest",
			wantBuilderURL: fmt.Sprintf(defaultBuilderURLTemplate, "php", "latest"),
		},
		{
			name:           "Extracts nodejs from nodejs18",
			runtime:        "nodejs18",
			tag:            "18",
			wantBuilderURL: fmt.Sprintf(defaultBuilderURLTemplate, "nodejs", "18"),
		},
		{
			name:           "Extracts dotnet from dotnet6",
			runtime:        "dotnet6",
			tag:            "123",
			wantBuilderURL: fmt.Sprintf(defaultBuilderURLTemplate, "dotnet", "123"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b := &buildpacksFunctionServer{
				runtime: tc.runtime,
				tag:     tc.tag,
			}

			builderImage, err := b.buildpackBuilderImage()

			if err != nil && !tc.wantError {
				t.Fatalf("Got unexpected error: %v", err)
			}

			if builderImage != tc.wantBuilderURL {
				t.Errorf("buildpackBuilderImage() = %q, want %q", builderImage, tc.wantBuilderURL)
			}
		})
	}
}
