//go:build (linux || darwin) && !bindings

package app

import (
	"github.com/unix-world/wails-app/internal/logger"
	"github.com/unix-world/wails-app/pkg/options"
)

func PreflightChecks(_ *options.App, _ *logger.Logger) error {
	return nil
}
