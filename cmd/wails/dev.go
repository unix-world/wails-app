package main

import (
	"os"

	"github.com/pterm/pterm"
	"github.com/unix-world/wails-app/cmd/wails/flags"
	"github.com/unix-world/wails-app/cmd/wails/internal/dev"
	"github.com/unix-world/wails-app/internal/colour"
	"github.com/unix-world/wails-app/pkg/clilogger"
)

func devApplication(f *flags.Dev) error {
	if f.NoColour {
		pterm.DisableColor()
		colour.ColourEnabled = false
	}

	quiet := f.Verbosity == flags.Quiet

	// Create logger
	logger := clilogger.New(os.Stdout)
	logger.Mute(quiet)

	if quiet {
		pterm.DisableOutput()
	} else {
		app.PrintBanner()
	}

	err := f.Process()
	if err != nil {
		return err
	}

	return dev.Application(f, logger)
}
