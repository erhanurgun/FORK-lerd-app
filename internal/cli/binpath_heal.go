package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/geodro/lerd/internal/config"
	"github.com/geodro/lerd/internal/services"
	lerdSystemd "github.com/geodro/lerd/internal/systemd"
)

// daemonUnits are the lerd services whose ExecStart names the lerd binary, so
// they are the ones a moved binary takes down. lerd-autostart is macOS-only and
// reads as absent elsewhere.
var daemonUnits = []string{"lerd-ui", "lerd-watcher", "lerd-tray", "lerd-autostart"}

// healLerdBinaryMove repairs what a lerd binary that moved leaves behind, and
// reports what it rewrote. A package manager can replace lerd underneath a
// working install: `brew upgrade lerd` retires the previous version's keg, and
// every unit and shim written against it then runs a path that is gone, which
// surfaces as daemons failing with a bare exit status 203 and `php` reporting
// no such file (#1432).
func healLerdBinaryMove() (units, shims []string) {
	return healDaemonUnits(), healShimBinaryPaths(config.LerdBinary())
}

// repairSummary is the line the start prints once a move has been repaired, so
// a rewrite of the user's own files is never silent.
func repairSummary(units, shims []string) string {
	var parts []string
	if len(units) > 0 {
		parts = append(parts, "services "+strings.Join(units, ", "))
	}
	if len(shims) > 0 {
		parts = append(parts, "shims "+strings.Join(shims, ", "))
	}
	return "Repointed " + strings.Join(parts, " and ") + " at " + config.LerdBinary()
}

// healDaemonUnits rewrites daemon units whose ExecStart names a binary that is
// no longer on disk, so they run the lerd that is installed now. A unit that
// still resolves is left alone, which keeps a start from a checkout build from
// quietly taking the login daemons with it.
func healDaemonUnits() []string {
	var healed []string
	for _, name := range daemonUnits {
		installed := services.InstalledUnitBinary(name)
		if installed == "" {
			continue
		}
		if _, err := os.Stat(installed); err == nil {
			continue
		}
		content, err := lerdSystemd.GetUnit(name)
		if err != nil {
			continue
		}
		// Nothing to gain when the rewrite would name the same missing file,
		// which is what a host without the tray helper installed looks like.
		if services.UnitExecBinary(content) == installed {
			continue
		}
		if err := writeUserServiceWithReload(name, content); err != nil {
			fmt.Printf("  WARN: repointing the %s unit at %s: %v\n", name, config.LerdBinary(), err)
			continue
		}
		healed = append(healed, name)
	}
	return healed
}

// healShimBinaryPaths repoints shims in lerd's bin dir that run a lerd binary
// which is no longer on disk. Only a dead path is rewritten: a shim that still
// resolves is left exactly as installed, so starting a build from a checkout
// never takes the user's shims with it.
func healShimBinaryPaths(lerdBin string) []string {
	if info, err := os.Stat(lerdBin); err != nil || info.IsDir() {
		return nil
	}
	binDir := config.BinDir()
	entries, err := os.ReadDir(binDir)
	if err != nil {
		return nil
	}
	var healed []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		path := filepath.Join(binDir, entry.Name())
		data, err := os.ReadFile(path)
		if err != nil || !strings.HasPrefix(string(data), "#!") {
			continue
		}
		repaired, changed := healedShim(string(data), lerdBin)
		if !changed {
			continue
		}
		if err := os.WriteFile(path, []byte(repaired), 0755); err != nil {
			fmt.Printf("  WARN: repointing the %s shim at %s: %v\n", entry.Name(), lerdBin, err)
			continue
		}
		healed = append(healed, entry.Name())
	}
	return healed
}

// healedShim swaps every absolute path to a lerd binary that is not there any
// more for lerdBin. Anything else the shim runs (composer.phar, a version
// manager, the tool itself) is left alone.
func healedShim(content, lerdBin string) (string, bool) {
	changed := false
	for _, token := range shimTokens(content) {
		if token == lerdBin || !filepath.IsAbs(token) || filepath.Base(token) != "lerd" {
			continue
		}
		if _, err := os.Stat(token); err == nil {
			continue
		}
		content = strings.ReplaceAll(content, token, lerdBin)
		changed = true
	}
	return content, changed
}

// shimTokens splits shim source into the words a shell would see, so a path is
// recognised whether it was written bare or quoted.
func shimTokens(content string) []string {
	return strings.FieldsFunc(content, func(r rune) bool {
		return unicode.IsSpace(r) || r == '"' || r == '\''
	})
}
