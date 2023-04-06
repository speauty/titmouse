package model

import "fmt"

type TaskModel struct {
	SeqNo          int    // 任务序号
	NumThreads     int    // 线程数量
	NumProcessors  int    // 处理器数量
	GraphicAdapter string // 显卡
	PathModel      string // 模型
	PathAudioFile  string // 音频文件
}

func (customTM *TaskModel) String() string {
	return fmt.Sprintf(
		"序号: %d, 线程数: %d, 处理器数: %d, 显卡: %s, 模型: %s, 音视频文件: %s",
		customTM.SeqNo, customTM.NumThreads, customTM.NumProcessors,
		customTM.GraphicAdapter, customTM.PathModel, customTM.PathAudioFile,
	)
}
