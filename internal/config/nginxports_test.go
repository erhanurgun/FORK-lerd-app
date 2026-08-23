package config

import "testing"

func TestNginxPorts(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	t.Run("an install with no config gets the defaults", func(t *testing.T) {
		http, https := NginxPorts()
		if http != 80 || https != 443 {
			t.Errorf("got %d/%d, want 80/443", http, https)
		}
	})

	t.Run("configured ports win", func(t *testing.T) {
		cfg := defaultConfig()
		cfg.Nginx.HTTPPort = 10080
		cfg.Nginx.HTTPSPort = 10443
		if err := SaveGlobal(cfg); err != nil {
			t.Fatalf("saving config: %v", err)
		}
		http, https := NginxPorts()
		if http != 10080 || https != 10443 {
			t.Errorf("got %d/%d, want 10080/10443", http, https)
		}
	})
}
