package main

import (
	"context"
	"fmt"
	"titmouse/cfg"
	"titmouse/cron"
	"titmouse/lib/log"
	"titmouse/lib/processor/nohup"
	"titmouse/lib/processor/whisper"
	"titmouse/model"
)

func main() {
	ctx := context.Background()
	if err := cfg.Api().Load(); err != nil {
		panic(err)
	}
	log.Api().Init(nil)

	if err := whisper.Api().Init(cfg.Api().WhisperCfg); err != nil {
		panic(err)
	}

	cron.ApiWhisperCron().Push(&model.TaskModel{
		Id:            "",
		NumThreads:    8,
		NumProcessors: 4,
		PathModel:     "E:\\下载\\whisper\\ggml-tiny.bin",
		PathAudioFile: "E:\\工作空间\\搬运翻译\\ChernoOpenGL\\26.创建纹理测试.mp4",
	})
	fmt.Println(cron.ApiWhisperCron().GetCntWait())
	nohup.NewResident(ctx, cron.ApiWhisperCron())
}
