//go:build windows && wv2runtime.browser
// +build windows,wv2runtime.browser

package wv2installer

import (
	"fmt"
	"github.com/unix-world/wails-app/internal/webview2runtime"
	"github.com/unix-world/wails-app/pkg/options/windows"
)

func doInstallationStrategy(installStatus installationStatus, messages *windows.Messages) error {
	confirmed, err := webview2runtime.Confirm(messages.DownloadPage+MinimumRuntimeVersion, messages.MissingRequirements)
	if err != nil {
		return err
	}
	if confirmed {
		err = webview2runtime.OpenInstallerDownloadWebpage()
		if err != nil {
			return err
		}
	}

	return fmt.Errorf(messages.FailedToInstall)
}
