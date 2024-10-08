package runtime

import (
	"context"

	"github.com/unix-world/wails-app/internal/frontend"
)

type Screen = frontend.Screen

// ScreenGetAll returns all screens
func ScreenGetAll(ctx context.Context) ([]Screen, error) {
	appFrontend := getFrontend(ctx)
	return appFrontend.ScreenGetAll()
}
