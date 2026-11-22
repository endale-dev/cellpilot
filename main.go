package main

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"log"

	"cellpilot/internal/repository/sqlite"

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

func run() (err error) {
	ctx := context.Background()
	db, err := openDatabase(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := db.Close(); cerr != nil && err == nil {
			err = fmt.Errorf("close database: %w", cerr)
		}
	}()

	if err = sqlite.Migrate(ctx, db); err != nil {
		return err
	}

	return wails.Run(applicationOptions(compose()))
}

func openDatabase(ctx context.Context) (*sql.DB, error) {
	path, err := sqlite.DefaultPath()
	if err != nil {
		return nil, err
	}
	return sqlite.Open(ctx, path)
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
