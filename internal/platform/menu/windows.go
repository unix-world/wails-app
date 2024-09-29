//go:build windows

package menu

import "github.com/unix-world/wails-app/internal/platform/win32"

type Menu struct {
	menu win32.HMENU
}
