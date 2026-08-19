package factorio

import (
	"path/filepath"
	"testing"
)

func TestCompareVersions(t *testing.T) {
	tests := []struct {
		v1       string
		v2       string
		expected int
	}{
		{"1.0.0", "1.0.0", 0},
		{"1.1.107", "2.0.28", -1},
		{"2.0.28", "1.1.107", 1},
		{"2.0.9", "2.0.28", -1},
		{"2.0.28", "2.0.9", 1},
		{"2.0.28", "2.0.28", 0},
		{"v2.0.28", "2.0.28", 0},
		{"2.0.28.1", "2.0.28", 1},
		{"2.0.28", "2.0.28.1", -1},
		{"", "1.0.0", -1},
	}

	for _, tc := range tests {
		t.Run(tc.v1+"_vs_"+tc.v2, func(t *testing.T) {
			result := CompareVersions(tc.v1, tc.v2)
			if result != tc.expected {
				t.Errorf("CompareVersions(%q, %q) = %d; expected %d", tc.v1, tc.v2, result, tc.expected)
			}
		})
	}
}

func TestVersionFileLifecycle(t *testing.T) {
	tempDir := t.TempDir()
	versionFile := filepath.Join(tempDir, "version.txt")

	// Missing file should return empty string and no error
	ver, err := ReadDiskVersion(versionFile)
	if err != nil {
		t.Fatalf("expected no error for missing file, got: %v", err)
	}
	if ver != "" {
		t.Fatalf("expected empty string for missing file, got: %q", ver)
	}

	// Update version file
	if err := UpdateVersionFile(versionFile, "2.0.28"); err != nil {
		t.Fatalf("failed to update version file: %v", err)
	}

	// Read version file
	ver, err = ReadDiskVersion(versionFile)
	if err != nil {
		t.Fatalf("failed to read version file: %v", err)
	}
	if ver != "2.0.28" {
		t.Fatalf("expected version %q, got: %q", "2.0.28", ver)
	}

	// Update to newer version
	if err := UpdateVersionFile(versionFile, "2.0.30\n"); err != nil {
		t.Fatalf("failed to update version file: %v", err)
	}

	ver, err = ReadDiskVersion(versionFile)
	if err != nil {
		t.Fatalf("failed to read version file: %v", err)
	}
	if ver != "2.0.30" {
		t.Fatalf("expected version %q, got: %q", "2.0.30", ver)
	}
}
