package whisper

import (
	"bufio"
	"go.uber.org/zap"
	"io"
	"os/exec"
	"strings"
	"sync"
	"titmouse/lib/log"
)

var (
	apiWhisper  *Whisper
	onceWhisper sync.Once
)

func Api() *Whisper {
	onceWhisper.Do(func() {
		apiWhisper = new(Whisper)
	})
	return apiWhisper
}

func New() *Whisper {
	return &Whisper{
		cfg:      Api().cfg,
		graphics: Api().graphics,
	}
}

type Whisper struct {
	cfg      *Cfg
	graphics []string
}

func (customW *Whisper) SetCfg(cfg *Cfg) error {
	customW.cfg = cfg
	return nil
}

func (customW *Whisper) Cfg() *Cfg {
	return customW.cfg
}

func (customW *Whisper) Graphics() []string {
	return customW.graphics
}

func (customW *Whisper) Init(cfg *Cfg) error {
	if customW.cfg != nil {
		customW.log().Warn("重复初始化")
		return nil
	}

	if err := customW.SetCfg(cfg); err != nil {
		return err
	}
	return nil
}

func (customW *Whisper) LoadGraphics() error {
	if err := customW.cfg.Validate(); err != nil {
		customW.log().Errorf("配置验证失败, 错误: %s", err)
		return err
	}

	currentCMD := exec.Command(customW.cfg.CmdPath, "-la")
	currentStdout, err := currentCMD.StdoutPipe()
	if err != nil {
		customW.log().Errorf("绑定标准输出失败, 错误: %s, 指令: %s", err, currentCMD.String())
		return err
	}
	if err := currentCMD.Start(); err != nil {
		customW.log().Errorf("执行执行失败, 错误: %s, 指令: %s", err, currentCMD.String())
		return err
	}

	newReader := bufio.NewReader(currentStdout)
	var graphics []string
	for true {
		line, _, lineErr := newReader.ReadLine()
		if lineErr != nil || lineErr == io.EOF {
			if lineErr != io.EOF {
				customW.log().Errorf("读取输出行失败, 错误: %s", lineErr.Error())
			}
			break
		}
		graphics = append(graphics, strings.ReplaceAll(string(line), "\"", ""))
	}
	_ = currentCMD.Wait()
	if len(graphics) > 1 {
		customW.graphics = graphics[1:]
	} else {
		customW.log().Warnf("未检测到相关显卡, 指令: %s", currentCMD.String())
	}
	return nil
}

func (customW *Whisper) GetName() string {
	return "whisper"
}

func (customW *Whisper) log() *zap.SugaredLogger {
	return log.Api().Scope(customW.GetName())
}
