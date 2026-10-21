package sidebar

import (
	"runtime"
	"testing"

	"github.com/yorukot/superfile/src/pkg/utils"
)

func TestShouldListDiskWithConfig(t *testing.T) {
	if runtime.GOOS == utils.OsWindows {
		// The prefix filtering only applies to POSIX systems; on Windows every
		// drive is always listed.
		t.Skip("prefix filtering is not used on Windows")
	}

	defaultInclude := []string{"/mnt", "/media", "/run/media", "/Volumes"}
	defaultExclude := []string{"/Volumes/.timemachine"}

	testCases := []struct {
		name       string
		mountPoint string
		include    []string
		exclude    []string
		expected   bool
	}{
		{
			name:       "root is always listed",
			mountPoint: "/",
			include:    defaultInclude,
			exclude:    defaultExclude,
			expected:   true,
		},
		{
			name:       "root is listed even with empty include list",
			mountPoint: "/",
			include:    nil,
			exclude:    nil,
			expected:   true,
		},
		{
			name:       "mount under an included prefix is listed",
			mountPoint: "/mnt/almighty",
			include:    defaultInclude,
			exclude:    defaultExclude,
			expected:   true,
		},
		{
			name:       "sshfs mount under /media is listed with defaults",
			mountPoint: "/media/remote",
			include:    defaultInclude,
			exclude:    defaultExclude,
			expected:   true,
		},
		{
			name:       "system mount not under any prefix is hidden",
			mountPoint: "/boot",
			include:    defaultInclude,
			exclude:    defaultExclude,
			expected:   false,
		},
		{
			name:       "proc is hidden with default prefixes",
			mountPoint: "/proc",
			include:    defaultInclude,
			exclude:    defaultExclude,
			expected:   false,
		},
		{
			name:       "excluded prefix wins over included prefix",
			mountPoint: "/Volumes/.timemachine/foo",
			include:    defaultInclude,
			exclude:    defaultExclude,
			expected:   false,
		},
		{
			name:       "empty include list lists every mount",
			mountPoint: "/proc",
			include:    []string{},
			exclude:    nil,
			expected:   true,
		},
		{
			name:       "empty include list still respects the exclude list",
			mountPoint: "/proc",
			include:    []string{},
			exclude:    []string{"/proc"},
			expected:   false,
		},
		{
			name:       "no include and no exclude lists everything",
			mountPoint: "/mnt/whatever",
			include:    nil,
			exclude:    nil,
			expected:   true,
		},
		{
			name:       "custom whitelist excludes non-matching mount",
			mountPoint: "/mnt/data",
			include:    []string{"/media"},
			exclude:    nil,
			expected:   false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := shouldListDiskWithConfig(tc.mountPoint, tc.include, tc.exclude)
			if got != tc.expected {
				t.Errorf("shouldListDiskWithConfig(%q, %v, %v) = %v, want %v",
					tc.mountPoint, tc.include, tc.exclude, got, tc.expected)
			}
		})
	}
}

func TestHasAnyPrefix(t *testing.T) {
	testCases := []struct {
		name     string
		s        string
		prefixes []string
		expected bool
	}{
		{name: "matches first prefix", s: "/mnt/x", prefixes: []string{"/mnt", "/media"}, expected: true},
		{name: "matches later prefix", s: "/media/x", prefixes: []string{"/mnt", "/media"}, expected: true},
		{name: "no match", s: "/boot", prefixes: []string{"/mnt", "/media"}, expected: false},
		{name: "empty prefixes", s: "/mnt/x", prefixes: nil, expected: false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := hasAnyPrefix(tc.s, tc.prefixes); got != tc.expected {
				t.Errorf("hasAnyPrefix(%q, %v) = %v, want %v", tc.s, tc.prefixes, got, tc.expected)
			}
		})
	}
}
