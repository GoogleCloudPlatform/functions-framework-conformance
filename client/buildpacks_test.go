package main

import (
	"fmt"
	"net"
	"testing"
	"time"
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

func TestWaitForServer(t *testing.T) {
	ln, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	addr := ln.Addr().String()

	if err := waitForServer(addr, time.Second); err != nil {
		t.Errorf("waitForServer(%q) with a listening server: got %v, want nil", addr, err)
	}

	ln.Close()
	if err := waitForServer(addr, 500*time.Millisecond); err == nil {
		t.Errorf("waitForServer(%q) with no server: got nil, want error", addr)
	}
}

func TestWaitForServerStartsLate(t *testing.T) {
	ln, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	addr := ln.Addr().String()
	ln.Close()

	go func() {
		time.Sleep(500 * time.Millisecond)
		late, err := net.Listen("tcp", addr)
		if err != nil {
			return
		}
		t.Cleanup(func() { late.Close() })
	}()

	if err := waitForServer(addr, 5*time.Second); err != nil {
		t.Errorf("waitForServer(%q) with a server that starts late: got %v, want nil", addr, err)
	}
}
