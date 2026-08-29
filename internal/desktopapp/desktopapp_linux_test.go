//go:build linux

package desktopapp

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInstall_writesAClickableEntryAndItsIcon(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_DATA_HOME", dir)

	entry, err := Install()
	if err != nil {
		t.Fatal(err)
	}
	if want := filepath.Join(dir, "applications", "lerd.desktop"); entry != want {
		t.Errorf("entry = %q, want %q", entry, want)
	}
	body, err := os.ReadFile(entry)
	if err != nil {
		t.Fatal(err)
	}
	text := string(body)
	// Clicked from a launcher there is no terminal, so the entry has to say so
	// and has to ask for the splash, or the click is a minute of nothing.
	for _, want := range []string{
		"Type=Application", "Name=Lerd", "Terminal=false",
		"dashboard --splash", "Categories=Development;",
	} {
		if !strings.Contains(text, want) {
			t.Errorf("desktop entry is missing %q:\n%s", want, text)
		}
	}
	icon := filepath.Join(dir, "lerd", "lerd.png")
	if !strings.Contains(text, "Icon="+icon) {
		t.Errorf("entry does not point at the installed icon:\n%s", text)
	}
	if info, err := os.Stat(icon); err != nil || info.Size() == 0 {
		t.Errorf("icon not written: %v", err)
	}
}

// Exec names the resolved binary: a launcher starts it with a minimal PATH that
// need not carry ~/.local/bin.
func TestDesktopEntry_execIsAnAbsolutePath(t *testing.T) {
	body := desktopEntry("/home/someone/.local/bin/lerd", "/tmp/lerd.png")
	if !strings.Contains(body, "Exec=/home/someone/.local/bin/lerd dashboard --splash") {
		t.Errorf("Exec line is not the resolved binary:\n%s", body)
	}
}

func TestRemove_takesTheEntryAndTheIcon(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_DATA_HOME", dir)
	if _, err := Install(); err != nil {
		t.Fatal(err)
	}
	if err := Remove(); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(Path()); !os.IsNotExist(err) {
		t.Errorf("desktop entry survived removal: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "lerd", "lerd.png")); !os.IsNotExist(err) {
		t.Errorf("icon survived removal: %v", err)
	}
	// Removing twice is what an uninstall on a machine that never installed
	// the entry does, and it must not fail.
	if err := Remove(); err != nil {
		t.Errorf("second Remove = %v, want nil", err)
	}
}
