package main

import (
	"github.com/pterm/pterm"
	"github.com/unix-world/wails-app/cmd/wails/flags"
	"github.com/unix-world/wails-app/cmd/wails/internal"
	"github.com/unix-world/wails-app/internal/colour"
	"github.com/unix-world/wails-app/internal/github"
)

func showReleaseNotes(f *flags.ShowReleaseNotes) error {
	if f.NoColour {
		pterm.DisableColor()
		colour.ColourEnabled = false
	}

	version := internal.Version
	if f.Version != "" {
		version = f.Version
	}

	app.PrintBanner()
	releaseNotes := github.GetReleaseNotes(version, f.NoColour)
	pterm.Println(releaseNotes)

	return nil
}
