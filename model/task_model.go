package model

import (
	"errors"
	"fmt"
	"os"
)

type TaskModel struct {
	Id             string
	NumThreads     int    // 线程数量
	NumProcessors  int    // 处理器数量
	GraphicAdapter string // 显卡
	PathModel      string // 模型
	PathAudioFile  string // 音频文件
	Language       string // 发音语种
}

func (customTM *TaskModel) String() string {
	return fmt.Sprintf(
		"编号: %s, 音视频文件: %s, 线程数: %d, 处理器数: %d, 显卡: %s, 模型: %s, 语言: %s",
		customTM.Id, customTM.PathAudioFile, customTM.NumThreads, customTM.NumProcessors,
		customTM.GraphicAdapter, customTM.PathModel, customTM.Language,
	)
}

func (customTM *TaskModel) Validate() error {
	chains := []func() error{
		customTM.validateChainPathModel, customTM.validateChainPathAudioFile,
	}
	for _, chain := range chains {
		if err := chain(); err != nil {
			return err
		}
	}
	return nil
}

func (customTM *TaskModel) validateChainPathModel() error {
	if customTM.PathModel == "" {
		return errors.New("未设置模型")
	}
	if _, err := os.Stat(customTM.PathModel); err != nil {
		return err
	}
	return nil
}

func (customTM *TaskModel) validateChainPathAudioFile() error {
	if customTM.PathAudioFile == "" {
		return errors.New("未设置音频文件")
	}
	if _, err := os.Stat(customTM.PathAudioFile); err != nil {
		return err
	}
	return nil
}
