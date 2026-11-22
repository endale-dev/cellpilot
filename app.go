package main

import "context"

type Config struct {
	Name string
}

type App struct {
	ctx  context.Context
	name string
}

func NewApp(cfg Config) *App {
	return &App{name: cfg.Name}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) Name() string {
	return a.name
}
