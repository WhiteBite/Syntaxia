package main

import (
	"context"
	"embed"
	"log"
	"os"
	"syntaxia/cmd/app"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed assets/ignore.glob
var embeddedIgnoreGlob string

const defaultCustomPromptRulesContent = "no additional rules"

// getStartupPath returns path from command line arguments if provided
func getStartupPath() string {
	if len(os.Args) > 1 {
		path := os.Args[1]
		// Verify path exists and is a directory
		if info, err := os.Stat(path); err == nil && info.IsDir() {
			return path
		}
	}
	return ""
}

func main() {
	appInstance := &App{}
	startupPath := getStartupPath()

	err := wails.Run(&options.App{
		Title:     "Syntaxia",
		Width:     1600,
		Height:    900,
		MinWidth:  1280,
		MinHeight: 720,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup: func(ctx context.Context) {
			container, err := app.NewContainer(ctx, embeddedIgnoreGlob, defaultCustomPromptRulesContent)
			if err != nil {
				log.Fatalf("Failed to create DI container: %v", err)
			}
			appInstance.startup(ctx, container)
			appInstance.startupPath = startupPath
		},
		OnDomReady: appInstance.domReady,
		OnShutdown: appInstance.shutdown,
		Bind: []interface{}{
			appInstance,
		},
	})

	if err != nil {
		log.Fatal(err)
	}
}
