package log

type Cfg struct {
	InfoFile  string
	ErrorFile string
}

func (customCfg *Cfg) Default() *Cfg {
	return &Cfg{
		InfoFile:  "logs/info.log",
		ErrorFile: "logs/err.log",
	}
}
