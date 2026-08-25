package ui

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// terminalScriptUTI is the file type of a .command script. The app LaunchServices
// records as its handler is what macOS means by the default terminal, and it is
// what "Open in terminal" should reach for before any emulator that merely
// happens to be on PATH.
const terminalScriptUTI = "com.apple.terminal.shell-script"

// launchServicesPlist is where the user's own handler choices live. The system
// defaults are not written here, so an absent file or a missing entry both mean
// the user never changed their terminal.
const launchServicesPlist = "Library/Preferences/com.apple.LaunchServices/com.apple.launchservices.secure.plist"

// macDefaultTerminalBundle reads the bundle id LaunchServices records for
// terminal scripts out of the plist's JSON form. It returns "" when the user
// never picked a terminal, which leaves the caller on its ordinary fallbacks.
func macDefaultTerminalBundle(plistJSON []byte) string {
	var doc struct {
		Handlers []struct {
			ContentType string `json:"LSHandlerContentType"`
			RoleAll     string `json:"LSHandlerRoleAll"`
			RoleShell   string `json:"LSHandlerRoleShell"`
		} `json:"LSHandlers"`
	}
	if err := json.Unmarshal(plistJSON, &doc); err != nil {
		return ""
	}
	for _, h := range doc.Handlers {
		if h.ContentType != terminalScriptUTI {
			continue
		}
		// A terminal claims the shell role, but the picker in Finder writes the
		// catch-all one, so take whichever the entry carries.
		if h.RoleAll != "" {
			return h.RoleAll
		}
		return h.RoleShell
	}
	return ""
}

// macDefaultTerminal returns the user's chosen default terminal, or "" when
// there is none to honour. plutil is what turns the binary plist into something
// readable without a cgo dependency on LaunchServices itself.
func macDefaultTerminal() string {
	if runtime.GOOS != "darwin" {
		return ""
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	out, err := exec.Command("plutil", "-convert", "json", "-o", "-", filepath.Join(home, launchServicesPlist)).Output()
	if err != nil {
		return ""
	}
	return macDefaultTerminalBundle(out)
}
