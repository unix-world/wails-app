//go:build !dev && !production && !bindings && (linux || darwin)

package app

import (
	"fmt"

	"github.com/unix-world/wails-app/pkg/options"
)

func (a *App) Run() error {
	return nil
}

// CreateApp creates the app!
func CreateApp(_ *options.App) (*App, error) {
	return nil, fmt.Errorf(`Wails applications will not build without the correct build tags.`)
}
