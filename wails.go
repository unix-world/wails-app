// Package wails is the main package of the Wails project.
// It is used by client applications.
package wails

import (
	_ "github.com/unix-world/wails-app/internal/goversion" // Add Compile-Time version check for minimum go version
	"github.com/unix-world/wails-app/pkg/application"
	"github.com/unix-world/wails-app/pkg/options"
)

// Run creates an application based on the given config and executes it
func Run(options *options.App) error {
	mainApp := application.NewWithOptions(options)
	return mainApp.Run()
}
