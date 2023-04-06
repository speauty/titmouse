package whisper

import (
	"errors"
	"os"
)

type Cfg struct {
	CmdPath           string
	DefaultGraphic    string
	DefaultThread     uint
	DefaultProcessors uint
	DefaultLanguage   string
	DefaultModelPath  string
}

func (customCfg *Cfg) Default() *Cfg {
	return &Cfg{
		DefaultThread:     1,
		DefaultProcessors: 2,
		DefaultLanguage:   "en",
	}
}

func (customCfg *Cfg) Validate() error {
	chains := []func() error{
		customCfg.validateChainCmd, customCfg.validateChainModel,
	}
	for _, chain := range chains {
		if err := chain(); err != nil {
			return err
		}
	}
	return nil
}

func (customCfg *Cfg) validateChainCmd() error {
	if customCfg.CmdPath == "" {
		return errors.New("the path of whisper.exe not configured")
	}
	if _, err := os.Stat(customCfg.CmdPath); err != nil {
		return err
	}
	return nil
}

func (customCfg *Cfg) validateChainModel() error {
	if customCfg.DefaultModelPath != "" {
		if _, err := os.Stat(customCfg.DefaultModelPath); err != nil {
			return err
		}
	}

	return nil
}
