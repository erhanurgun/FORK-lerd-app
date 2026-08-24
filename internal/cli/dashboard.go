package cli

import (
	"fmt"

	"github.com/geodro/lerd/internal/dashboard"
	"github.com/geodro/lerd/internal/desktopnotify"
	"github.com/spf13/cobra"
)

// NewDashboardCmd returns the dashboard command.
func NewDashboardCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "dashboard",
		Short: "Open the Lerd dashboard, starting Lerd first if it is not running",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			// Opening the dashboard on a stopped lerd is a request to work on
			// something, so bring the stack up rather than land on a page whose
			// only content is a button to press.
			if !dashboard.Serving() {
				if err := startLerd(nil, nil); err != nil {
					return err
				}
			}
			// Prefer the desktop app when it's the registered lerd:// handler;
			// it focuses the running window rather than opening a new tab.
			if desktopnotify.AppInstalled() {
				if err := desktopnotify.OpenApp(""); err == nil {
					fmt.Println("Opening the Lerd desktop app")
					return nil
				}
			}
			url := dashboard.URL()
			fmt.Printf("Opening %s\n", url)
			return openBrowser(url)
		},
	}
}
