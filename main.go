package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/go-toast/toast"
	"github.com/spf13/cobra"
)

// Overridden via ldflags
var (
	version   = "99.0.1-devbuild"
	commit    = "unknown"
	date      = "unknown"
	goversion = "unknown"
)

func main() {
	var showHelp bool
	var showVersion bool
	var icon string
	var onClick string
	var category string
	var appID string

	rootCmd := &cobra.Command{
		Use:   "wsl-notify-send",
		Short: "wsl-notify-send - a WSL integration for notify-send",
		Long:  "wsl-notify-send provides a Windows.exe that accepts parameters similar to the Linux notify-send utility to aid interop. For more customisability, see the toast CLI at https://github.com/go-toast/toast",
		Run: func(cmd *cobra.Command, args []string) {
			if showVersion {
				fmt.Printf("wsl-notify-send version %s\nBuilt %s (commit %s)\n%s\n\n", version, date, commit, goversion)
				return
			}
			if showHelp || len(args) == 0 {
				_ = cmd.Usage()
				return
			}

			random := strconv.Itoa(rand.Intn(100) + 1)

			if len(icon) > 0 && (strings.HasPrefix(icon, "http://") || strings.HasPrefix(icon, "https://")) {
				tmpFolder := os.TempDir()

				err := DownloadFile(icon, filepath.Join(tmpFolder, "wsl-notify-send-icon-tmp"+random+".png"))
				if err != nil {
					log.Fatalln(err)
				}
				tmpFile.Close()
				err = DownloadFile(icon, tmpFile.Name())
				if err != nil {
			if err != nil {
				log.Fatalln(err)
			} else {
				} else {
					// had to comment this out because the toast wasn't getting invoked before the file was removed
					// defer os.Remove("wsl-notify-send-icon-tmp"+random+".png")
					icon = filepath.Join(tmpFolder, "wsl-notify-send-icon-tmp"+random+".png")
				}
			}

			notification := &toast.Notification{
				AppID:               appID,
				Title:               category,
				Message:             args[0],
				ActivationArguments: onClick,
				Icon:                icon,
			}

			if err := notification.Push(); err != nil {
				log.Fatalln(err)
			}
		},
	}
	// Standard flags
	rootCmd.Flags().BoolVarP(&showHelp, "help", "?", false, "Show a help message")
	rootCmd.Flags().StringVarP(&icon, "icon", "i", "", "An icon filename to display (stock icons are not currently supported)")
	rootCmd.Flags().StringVarP(&category, "category", "c", "wsl-notify-send", "Specifies the notification category")
	rootCmd.Flags().StringVarP(&onClick, "onClick", "l", "", "Link which should open when clicking the notification")
	// Standard flags that are ignored
	rootCmd.Flags().IntP("expire-time", "t", -1, "[Ignored in wsl-notiy-send]") // TODO - extend go-toast to support https://docs.microsoft.com/en-us/uwp/api/windows.ui.notifications.toastnotification.expirationtime?view=winrt-19041
	rootCmd.Flags().StringArrayP("hint", "h", []string{}, "Ignored in wsl-notify-send")
	rootCmd.Flags().StringArrayP("urgency", "u", []string{}, "Ignored in wsl-notify-send")
	// Custom flags
	rootCmd.Flags().StringVar(&appID, "appId", "wsl-notify-send", "[non-standard] Specifies the app ID")
	rootCmd.Flags().BoolVar(&showVersion, "version", false, "Show version information")
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error executing command: %v", err)
	}
}

func DownloadFile(url string, filepath string) error {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Ensure we received a successful response before writing to file
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download file from %s: %s", url, resp.Status)
	}

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

// TODO - explore mapping icons: https://wiki.ubuntu.com/NotificationDevelopmentGuidelines#How_do_I_get_these_slick_icons
//      https://docs.microsoft.com/en-us/uwp/api/windows.ui.notifications.toastnotification?view=winrt-19041
//      https://docs.microsoft.com/en-us/uwp/schemas/tiles/toastschema/schema-root
