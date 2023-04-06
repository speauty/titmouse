package cfg

import (
	"sync"
	"titmouse/lib/processor/whisper"
)

var (
	apiCfg  *Cfg
	onceCfg sync.Once
)

func Api() *Cfg {
	onceCfg.Do(func() {
		apiCfg = new(Cfg)
		apiCfg.WhisperCfg = new(whisper.Cfg).Default()
	})
	return apiCfg
}

type Cfg struct {
	FlagInstalled bool
	WhisperCfg    *whisper.Cfg
}
