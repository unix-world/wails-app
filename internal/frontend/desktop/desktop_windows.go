//go:build windows
// +build windows

package desktop

import (
	"context"
	"github.com/unix-world/wails-app/internal/binding"
	"github.com/unix-world/wails-app/internal/frontend"
	"github.com/unix-world/wails-app/internal/frontend/desktop/windows"
	"github.com/unix-world/wails-app/internal/logger"
	"github.com/unix-world/wails-app/pkg/options"
)

func NewFrontend(ctx context.Context, appoptions *options.App, logger *logger.Logger, appBindings *binding.Bindings, dispatcher frontend.Dispatcher) frontend.Frontend {
	return windows.NewFrontend(ctx, appoptions, logger, appBindings, dispatcher)
}
