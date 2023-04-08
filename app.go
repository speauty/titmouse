package main

import (
	"context"
	"fmt"
	"os"
	"titmouse/adapter/gui_windows"
	"titmouse/cfg"
	"titmouse/cron"
	"titmouse/lib/log"
	"titmouse/lib/processor/nohup"
	"titmouse/lib/processor/whisper"
)

func main() {
	ctx, _ := context.WithCancel(context.Background())
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

	nohup.NewResident(ctx, cron.ApiWhisperCron(), gui_windows.ApiGuiWindows())
}

func install() error {
	fd, err := os.Stat("logs")
	if err != nil || (fd != nil && !fd.IsDir()) {
		if err = os.Mkdir("logs", os.ModePerm); err != nil {
			return fmt.Errorf("创建日志目录失败(如果存在相应logs文件, 请手动处理), 错误: %s", err)
		}
	}

	if _, err = os.Stat(cfg.Api().GetFilename()); err != nil {
		return cfg.Api().Store()
	}
	return nil
}
