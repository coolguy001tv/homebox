package services

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestIsWithin(t *testing.T) {
	// 用 os.TempDir() 构造双平台绝对路径：Windows 为 C:\...\Temp\TestImages，
	// Linux 为 /tmp/TestImages。isWithin 是纯字符串逻辑，路径无需真实存在。
	base := filepath.Join(os.TempDir(), "TestImages")

	tests := []struct {
		name   string
		target string
		want   bool
	}{
		{name: "base itself", target: base, want: true},
		{name: "base itself different case", target: strings.ToLower(base), want: true},
		{name: "file inside", target: filepath.Join(base, "foo.jpg"), want: true},
		{name: "subdir inside", target: filepath.Join(base, "subdir", "bar.jpg"), want: true},
		{name: "outside", target: filepath.Join(os.TempDir(), "other", "file.jpg"), want: false},
		{name: "sibling folder", target: base + "-other", want: false},
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
	absBase, _ := filepath.Abs(filepath.Join(os.TempDir(), "TestImages"))
	absBase = filepath.Clean(absBase)

	tests := []struct {
		name   string
		path   string
		wantOK bool
	}{
		{name: "root empty", path: "", wantOK: true},
		{name: "root dot", path: ".", wantOK: true},
		{name: "absolute file", path: filepath.Join(absBase, "foo.jpg"), wantOK: true},
		{name: "relative subdir", path: "subdir", wantOK: true},
		// 前向斜杠穿越路径在双平台都能解析：Windows 把 / 视为分隔符，
		// Linux 亦然，Clean 后落到 base 之外。
		{name: "traversal attempt", path: "../Windows/System32", wantOK: false},
		{name: "absolute outside", path: filepath.Join(os.TempDir(), "othere", "file.jpg"), wantOK: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var target string
			switch {
			case filepath.IsAbs(tt.path):
				target = filepath.Clean(tt.path)
			case tt.path == "." || tt.path == "":
				target = absBase
			default:
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
