package build

import (
	"github.com/unix-world/wails-app/internal/project"
	"github.com/unix-world/wails-app/pkg/clilogger"
)

// Builder defines a builder that can build Wails applications
type Builder interface {
	SetProjectData(projectData *project.Project)
	BuildFrontend(logger *clilogger.CLILogger) error
	CompileProject(options *Options) error
	OutputFilename(options *Options) string
	CleanUp()
}
