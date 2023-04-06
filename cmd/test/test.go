package main

import (
	"titmouse/cfg"
	"titmouse/lib/log"
	"titmouse/lib/processor/whisper"
)

func main() {
	if err := cfg.Api().Load(); err != nil {
		panic(err)
	}
	log.Api().Init(nil)

	if err := whisper.Api().Init(cfg.Api().WhisperCfg); err != nil {
		panic(err)
	}
}
