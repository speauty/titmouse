package cron

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-module/carbon"
	"runtime"
	"sync"
	"sync/atomic"
	"titmouse/lib/log"
	"titmouse/lib/processor/whisper"
	"titmouse/model"
)

var (
	apiWhisperCron  *WhisperCron
	onceWhisperCron sync.Once
)

func ApiWhisperCron() *WhisperCron {
	onceWhisperCron.Do(func() {
		apiWhisperCron = new(WhisperCron)
		apiWhisperCron.init()
	})
	return apiWhisperCron
}

type WhisperData struct {
	id   string           // 数据中的id
	data *model.TaskModel // 数据(已校验过的, 否则会出现异常)
}

func (customWD *WhisperData) toWhisperTransformParams() *whisper.TransformParams {
	return &whisper.TransformParams{
		NumThreads:     customWD.data.NumThreads,
		NumProcessors:  customWD.data.NumProcessors,
		GraphicAdapter: customWD.data.GraphicAdapter,
		PathModel:      customWD.data.PathModel,
		PathAudioFile:  customWD.data.PathAudioFile,
		Language:       customWD.data.Language,
	}
}

type WhisperCron struct {
	ctx              context.Context
	fnCancel         context.CancelFunc
	listInMemory     sync.Map
	cntListInMemory  *atomic.Int32
	chanTransform    chan *WhisperData
	maxChanTransform int
	chanMsg          chan string
	chanMsgRedirect  chan string
}

func (customWC *WhisperCron) SetChanMsgRedirect(chanMsgRedirect chan string) {
	customWC.chanMsgRedirect = chanMsgRedirect
}

func (customWC *WhisperCron) GetCntListInMemory() int32 {
	return customWC.cntListInMemory.Load()
}

func (customWC *WhisperCron) Push(task *model.TaskModel) error {
	if _, flagIsExisted := customWC.listInMemory.Load(task.Id); flagIsExisted {
		return errors.New("当前任务已存在, 请勿反复推送")
	}
	currentData := &WhisperData{id: task.Id, data: task}
	customWC.listInMemory.Store(task.Id, currentData)
	customWC.cntListInMemory.Add(1)
	customWC.chanTransform <- currentData
	return nil
}

func (customWC *WhisperCron) Run(ctx context.Context, fnCancel context.CancelFunc) {
	customWC.ctx = ctx
	customWC.fnCancel = fnCancel

	customWC.jobTransform()
	customWC.jobMsgRedirect()
}

func (customWC *WhisperCron) Close() {
	fmt.Println("WhisperCron exit")
}

func (customWC *WhisperCron) jobTransform() {
	for i := 0; i < customWC.maxChanTransform; i++ {
		go func(ctx context.Context, currentChanIdx int, chanTransform chan *WhisperData, chanMsg chan string) {
			coroutineName := fmt.Sprintf("Whisper转换任务@%d", currentChanIdx)
			currentLog := customWC.log().Scope(coroutineName)
			apiWhisper := whisper.Api()

			for true {
				select {
				case <-ctx.Done():
					currentLog.Warnf("%s被迫结束, 错误: 主程序退出", coroutineName)
					runtime.Goexit()
				case currentTransform, isOpen := <-chanTransform:
					timeStarted := carbon.Now()
					if isOpen == false && currentTransform == nil {
						currentLog.Warnf("%s被迫结束, 错误: 通道关闭", coroutineName)
						runtime.Goexit()
					}
					customWC.cntListInMemory.Add(-1)
					chanMsg <- fmt.Sprintf("正在处理, 当前数据[%s]", currentTransform.data)
					if err := apiWhisper.Transform(currentTransform.toWhisperTransformParams()); err != nil {
						currentLog.Errorf("数据[%s], %s", currentTransform.data.String(), err)
						chanMsg <- fmt.Sprintf("数据[%s], %s", currentTransform.data.String(), err)
						continue
					}
					currentLog.Infof("转换成功, 当前数据[%s], 耗时: %d", currentTransform.data, carbon.Now().DiffAbsInSeconds(timeStarted))
					chanMsg <- fmt.Sprintf("转换成功, 当前数据[%s], 耗时: %d", currentTransform.data, carbon.Now().DiffAbsInSeconds(timeStarted))
					customWC.listInMemory.Delete(currentTransform.id)
				}
			}
		}(customWC.ctx, i, customWC.chanTransform, customWC.chanMsg)
	}
}

func (customWC *WhisperCron) jobMsgRedirect() {
	go func(ctx context.Context, chanMsg, chanMsgRedirect chan string) {
		coroutineName := "消息重定向"
		currentLog := customWC.log().Scope(coroutineName)
		for true {
			select {
			case <-ctx.Done():
				currentLog.Warnf("%s被迫结束, 错误: 主程序退出", coroutineName)
				runtime.Goexit()
			case currentMsg, isOpen := <-chanMsg:
				if isOpen == false && currentMsg == "" {
					currentLog.Warnf("%s被迫结束, 错误: 通道关闭", coroutineName)
					runtime.Goexit()
				}
				fmt.Println(currentMsg)
				if chanMsgRedirect == nil {
					continue
				}
				chanMsgRedirect <- currentMsg
			}
		}
	}(customWC.ctx, customWC.chanMsg, customWC.chanMsgRedirect)
}

func (customWC *WhisperCron) init() {
	customWC.maxChanTransform = 10
	customWC.chanTransform = make(chan *WhisperData, customWC.maxChanTransform)
	customWC.chanMsg = make(chan string, 50)
	customWC.cntListInMemory = new(atomic.Int32)
}

func (customWC *WhisperCron) log() *log.Log {
	return log.Api()
}
