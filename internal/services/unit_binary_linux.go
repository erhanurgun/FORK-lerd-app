//go:build linux

package services

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/geodro/lerd/internal/config"
)

// InstalledUnitBinary returns the program the installed unit runs, or "" when
// no unit of that name is on disk. Callers use it to tell a unit that still
// names a binary which is there from one left behind by a binary that moved.
func InstalledUnitBinary(name string) string {
	data, err := os.ReadFile(filepath.Join(config.SystemdUserDir(), name+".service"))
	if err != nil {
		return ""
	}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "ExecStart=") {
			continue
		}
		args := SplitExecStart(strings.TrimPrefix(line, "ExecStart="))
		if len(args) == 0 {
			return ""
		}
		return args[0]
	}
	return ""
}
