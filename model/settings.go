package model

import "titmouse/lib/storage/single_file"

type Settings struct {
	Name string
}

func (customS *Settings) GetFilename() string {
	return "settings.tm"
}

func (customS *Settings) Store() error {
	return single_file.Store(customS)
}

func (customS *Settings) Load() error {
	return single_file.Load(customS)
}
