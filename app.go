package main

import (
	"context"
	"os"
	"titmouse/adapter/gui_windows"
	"titmouse/cfg"
	"titmouse/cron"
	"titmouse/lib/log"
	"titmouse/lib/processor/whisper"
)

func main() {
	ctx := context.Background()
	if err := install(); err != nil {
		panic(err)
	}

	if err := cfg.Api().Load(); err != nil {
		panic(err)
	}
	log.Api().Init(nil)

	if err := whisper.Api().Init(cfg.Api().WhisperCfg); err != nil {
		panic(err)
	}

	if err := gui_windows.ApiGuiWindows().Init(); err != nil {
		panic(err)
	}

	go func() { cron.ApiWhisperCron().Run(ctx, nil) }()

	gui_windows.ApiGuiWindows().Run(ctx, nil)
}

func install() error {
	if _, err := os.Stat(cfg.Api().GetFilename()); err != nil {
		return cfg.Api().Store()
	}
	return nil
}
