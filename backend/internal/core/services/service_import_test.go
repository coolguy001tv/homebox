package services

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestIsWithin(t *testing.T) {
	base := `C:\test-images`

	tests := []struct {
		name   string
		target string
		want   bool
	}{
		{name: "base itself", target: `C:\test-images`, want: true},
		{name: "base itself lower", target: `c:\test-images`, want: true},
		{name: "file inside", target: `C:\test-images\foo.jpg`, want: true},
		{name: "subdir inside", target: `C:\test-images\subdir\bar.jpg`, want: true},
		{name: "outside", target: `C:\Windows\explorer.exe`, want: false},
		{name: "sibling folder", target: `C:\test-images-other`, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isWithin(tt.target, base)
			if got != tt.want {
				t.Errorf("isWithin(%q, %q) = %v, want %v", tt.target, base, got, tt.want)
			}
		})
	}
}

func TestIsWithin_Subpath(t *testing.T) {
	absBase, _ := filepath.Abs(`C:\test-images`)
	absBase = filepath.Clean(absBase)

	tests := []struct {
		name   string
		path   string
		wantOK bool
	}{
		{name: "root empty", path: "", wantOK: true},
		{name: "root dot", path: ".", wantOK: true},
		{name: "absolute file", path: `C:\test-images\foo.jpg`, wantOK: true},
		{name: "relative subdir", path: `subdir`, wantOK: true},
		{name: "traversal attempt", path: `..\Windows\System32`, wantOK: false},
		{name: "absolute outside", path: `D:\other\file.jpg`, wantOK: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var target string
			if filepath.IsAbs(tt.path) {
				target = filepath.Clean(tt.path)
			} else if tt.path == "." || tt.path == "" {
				target = absBase
			} else {
				target = filepath.Clean(filepath.Join(absBase, tt.path))
			}

			ok := target == absBase || (len(target) > len(absBase)+1 &&
				strings.EqualFold(target[:len(absBase)+1], absBase+string(filepath.Separator)))

			t.Logf("target=%q base=%q ok=%v", target, absBase, ok)

			if ok != tt.wantOK {
				t.Errorf("path=%q: got ok=%v, want %v", tt.path, ok, tt.wantOK)
			}
		})
	}
}
