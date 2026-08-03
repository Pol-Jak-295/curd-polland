package internal

import (
	"path/filepath"
	"strconv"
	"strings"
)

var curdVersion = "dev"

// SetCurdVersion records the running application version for storage migrations.
func SetCurdVersion(version string) {
	version = strings.TrimSpace(version)
	if version != "" {
		curdVersion = version
	}
}

// CurdVersion returns the running application version.
func CurdVersion() string {
	if curdVersion == "" {
		return "dev"
	}
	return curdVersion
}

func storageVersionFilePath(storagePath string) string {
	return filepath.Join(strings.TrimSpace(storagePath), "curd_version")
}

// parseVersionParts parses "2.0.3" / "v2.0.3-rc1" into [major, minor, patch].
// Non-numeric suffixes are ignored. Unparseable versions return nil.
func parseVersionParts(version string) []int {
	version = strings.TrimSpace(version)
	version = strings.TrimPrefix(version, "v")
	version = strings.TrimPrefix(version, "V")
	if version == "" || strings.EqualFold(version, "dev") {
		return nil
	}
	// Drop pre-release / build metadata: 2.0.3-rc1+meta → 2.0.3
	if i := strings.IndexAny(version, "-+"); i >= 0 {
		version = version[:i]
	}
	parts := strings.Split(version, ".")
	if len(parts) == 0 {
		return nil
	}
	out := make([]int, 0, 3)
	for i := 0; i < 3; i++ {
		if i >= len(parts) {
			out = append(out, 0)
			continue
		}
		n, err := strconv.Atoi(parts[i])
		if err != nil {
			return nil
		}
		out = append(out, n)
	}
	return out
}

// compareVersions returns -1 if a<b, 0 if a==b, 1 if a>b.
// Unparseable versions compare as equal to each other and less than parseable ones
// only when empty; "dev" is treated as greater than any release (no auto-inject from dev churn).
func compareVersions(a, b string) int {
	a = strings.TrimSpace(a)
	b = strings.TrimSpace(b)
	if a == b {
		return 0
	}
	ap := parseVersionParts(a)
	bp := parseVersionParts(b)
	if ap == nil && bp == nil {
		if a < b {
			return -1
		}
		if a > b {
			return 1
		}
		return 0
	}
	if ap == nil {
		// empty/unknown stored version → treat as older than any real release
		if a == "" {
			return -1
		}
		// "dev" or garbage: treat as equal to itself only; vs release → greater (skip inject churn)
		return 1
	}
	if bp == nil {
		if b == "" {
			return 1
		}
		return -1
	}
	for i := 0; i < 3; i++ {
		if ap[i] < bp[i] {
			return -1
		}
		if ap[i] > bp[i] {
			return 1
		}
	}
	return 0
}

// versionLess reports whether a is strictly older than b.
func versionLess(a, b string) bool { return compareVersions(a, b) < 0 }

// versionLessOrEqual reports whether a is older than or equal to b.
func versionLessOrEqual(a, b string) bool { return compareVersions(a, b) <= 0 }
