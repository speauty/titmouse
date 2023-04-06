package cfg

import "titmouse/lib/storage/single_file"

func (customCfg *Cfg) GetFilename() string {
	return "cfg.data"
}

func (customCfg *Cfg) Store() error {
	return single_file.Store(customCfg)
}

func (customCfg *Cfg) Load() error {
	return single_file.Load(customCfg)
}
