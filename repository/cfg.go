package repository

import (
	"sync"
	"titmouse/cfg"
	"titmouse/lib/processor/whisper"
)

var (
	apiCfgRepository  *CfgRepository
	onceCfgRepository sync.Once
)

func ApiCfg() *CfgRepository {
	onceCfgRepository.Do(func() {
		apiCfgRepository = new(CfgRepository)
		apiCfgRepository.whisperClient = whisper.Api()
	})
	return apiCfgRepository
}

// CfgRepository 关于配置的相关操作仓库
type CfgRepository struct {
	*Repository
	whisperClient *whisper.Whisper
}

// ActionGetGraphics 获取显卡集合
func (customCR CfgRepository) ActionGetGraphics() []string {
	_ = customCR.whisperClient.LoadGraphics()
	return customCR.whisperClient.Graphics()
}

// ActionGetWhisperCfg 获取whisper配置
func (customCR CfgRepository) ActionGetWhisperCfg() *whisper.Cfg {
	return customCR.whisperClient.Cfg()
}

// ActionSetWhisperCfg 设置whisper配置
func (customCR CfgRepository) ActionSetWhisperCfg(whisperCfg *whisper.Cfg) error {
	defer func() {
		_ = customCR.ActionSync()
	}()
	if whisperCfg == nil {
		return nil
	}
	return customCR.whisperClient.SetCfg(whisperCfg)
}

func (customCR CfgRepository) ActionGetFlagCacheWhisperData() bool {
	return cfg.Api().FlagCacheWhisperData
}

func (customCR CfgRepository) ActionSetFlagCacheWhisperData(flag bool) {
	defer func() {
		_ = customCR.ActionSync()
	}()
	cfg.Api().FlagCacheWhisperData = flag
}

// ActionSync 同步配置
func (customCR CfgRepository) ActionSync() error {
	return cfg.Api().Store()
}
