//go:build linux

package desktopapp

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/geodro/lerd/internal/config"
)

//go:embed assets/lerd-mark.png
var markPNG []byte

// entryName is the desktop file's basename. It doubles as the icon's, so a
// launcher that resolves Icon= through the theme finds the same artwork.
const entryName = "lerd"

// dataHome is the XDG base for the desktop entry and the icon, held in a var so
// a test can point it somewhere harmless.
var dataHome = func() string {
	if v := os.Getenv("XDG_DATA_HOME"); v != "" {
		return v
	}
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return ""
	}
	return filepath.Join(home, ".local", "share")
}

// Path is the desktop entry, empty when there is nowhere to write it.
func Path() string {
	dir := dataHome()
	if dir == "" {
		return ""
	}
	return filepath.Join(dir, "applications", entryName+".desktop")
}

// iconPath is where the entry's Icon= points. An absolute path rather than a
// theme name: a themed icon has to land in the right size directory to be found
// at all, and the mark ships at one size.
func iconPath() string {
	return filepath.Join(dataHome(), entryName, entryName+".png")
}

// Install writes the desktop entry and its icon, and returns the entry's path.
// Rewritten on every run so an upgraded lerd repoints it at the binary that is
// live now.
//
// The entry runs `lerd dashboard --splash`, which is the same command the tray
// and the CLI use: it starts the environment when it is down, opens the desktop
// app when that is installed and the browser otherwise. --splash is what makes
// it usable from an icon, where there is no terminal to show progress in.
func Install() (string, error) {
	entry := Path()
	if entry == "" {
		return "", fmt.Errorf("locating the home directory")
	}
	icon := iconPath()
	if err := os.MkdirAll(filepath.Dir(icon), 0755); err != nil {
		return "", fmt.Errorf("creating %s: %w", filepath.Dir(icon), err)
	}
	if err := os.WriteFile(icon, markPNG, 0644); err != nil {
		return "", fmt.Errorf("writing the icon: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(entry), 0755); err != nil {
		return "", fmt.Errorf("creating %s: %w", filepath.Dir(entry), err)
	}
	if err := os.WriteFile(entry, []byte(desktopEntry(config.LerdBinary(), icon)), 0644); err != nil {
		return "", fmt.Errorf("writing the desktop entry: %w", err)
	}
	refreshDesktopDatabase(filepath.Dir(entry))
	return entry, nil
}

// desktopEntry renders the .desktop file. Exec is the resolved binary rather
// than a bare name: a launcher starts it with a minimal PATH that need not carry
// ~/.local/bin.
func desktopEntry(bin, icon string) string {
	return "[Desktop Entry]\n" +
		"Type=Application\n" +
		"Name=" + Name + "\n" +
		"GenericName=Local PHP development environment\n" +
		"Comment=Start Lerd and open its dashboard\n" +
		"Exec=" + bin + " dashboard --splash\n" +
		"Icon=" + icon + "\n" +
		"Terminal=false\n" +
		"Categories=Development;\n" +
		"Keywords=php;laravel;podman;lerd;\n" +
		"StartupNotify=true\n"
}

// refreshDesktopDatabase asks the desktop to notice the new entry now rather
// than at the next login. Best effort: the entry is readable either way.
func refreshDesktopDatabase(dir string) {
	bin, err := exec.LookPath("update-desktop-database")
	if err != nil {
		return
	}
	_ = exec.Command(bin, dir).Run()
}

// Remove deletes the desktop entry and its icon. Missing is success.
func Remove() error {
	entry := Path()
	if entry == "" {
		return nil
	}
	if err := os.Remove(entry); err != nil && !os.IsNotExist(err) {
		return err
	}
	_ = os.RemoveAll(filepath.Dir(iconPath()))
	refreshDesktopDatabase(filepath.Dir(entry))
	return nil
}
