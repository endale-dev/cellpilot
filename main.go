package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	return wails.Run(applicationOptions(compose()))
}

func compose() *App {
	return NewApp(Config{Name: "CellPilot"})
}

func applicationOptions(app *App) *options.App {
	return &options.App{
		Title:     app.Name(),
		Width:     1280,
		Height:    800,
		MinWidth:  960,
		MinHeight: 640,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 244, G: 245, B: 247, A: 255},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	}
}
