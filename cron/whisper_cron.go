package cron

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-module/carbon"
	"path/filepath"
	"runtime"
	"sort"
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
	id                 string           // 数据中的id
	data               *model.TaskModel // 数据(已校验过的, 否则会出现异常)
	flagIsTransforming bool             // 是否正在转换
	pushedAt           carbon.Carbon
	transformedAt      carbon.Carbon
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
	ctx                context.Context
	fnCancel           context.CancelFunc
	listInMemory       sync.Map
	cntList            *atomic.Int32
	cntWait            *atomic.Int32
	audioFilesInMemory sync.Map
	chanTransform      chan *WhisperData
	maxChanTransform   int
	chanMsg            chan string
	chanMsgRedirect    chan string
	maxCoroutine       int
	timeStarted        carbon.Carbon
}

func (customWC *WhisperCron) SetChanMsgRedirect(chanMsgRedirect chan string) {
	customWC.chanMsgRedirect = chanMsgRedirect
}

func (customWC *WhisperCron) GetCntWait() int {
	return int(customWC.cntWait.Load())
}

func (customWC *WhisperCron) GetMaxWait() int {
	return customWC.maxChanTransform - customWC.maxCoroutine
}

func (customWC *WhisperCron) GetCntTask() int {
	return int(customWC.cntList.Load())
}

func (customWC *WhisperCron) GetTimeRun() int {
	return int(carbon.Now().DiffAbsInSeconds(customWC.timeStarted))
}

func (customWC *WhisperCron) GetListInMemory() []string {
	var results []string
	customWC.listInMemory.Range(func(idx, value any) bool {
		tmpData := value.(*WhisperData)
		tmpStr := fmt.Sprintf("排队中(%ds)", carbon.Now().DiffAbsInSeconds(tmpData.pushedAt))
		if tmpData.flagIsTransforming {
			tmpStr = fmt.Sprintf("转换中(%ds)", carbon.Now().DiffAbsInSeconds(tmpData.transformedAt))
		}

		results = append(results, fmt.Sprintf(
			"%s, 状态: %s",
			filepath.Base(tmpData.data.PathAudioFile), tmpStr,
		))
		return true
	})
	sort.Strings(results)
	return results
}

func (customWC *WhisperCron) Push(task *model.TaskModel) error {
	if _, flagIsExisted := customWC.audioFilesInMemory.Load(task.PathAudioFile); flagIsExisted {
		return errors.New("当前任务已存在, 请勿反复推送")
	}
	if _, flagIsExisted := customWC.listInMemory.Load(task.Id); flagIsExisted {
		return errors.New("当前任务已存在, 请勿反复推送")
	}
	numWait := int(customWC.cntWait.Load())
	if customWC.maxChanTransform <= numWait+customWC.maxCoroutine {
		return errors.New(fmt.Sprintf("当前任务过多(数量: %d), 请稍后推送", numWait))
	}
	currentData := &WhisperData{id: task.Id, data: task, flagIsTransforming: false, pushedAt: carbon.Now()}
	customWC.listInMemory.Store(task.Id, currentData)
	customWC.cntWait.Add(1)
	customWC.cntList.Add(1)
	customWC.audioFilesInMemory.Store(task.PathAudioFile, 1)
	customWC.chanMsg <- fmt.Sprintf("已接收, 当前数据[%s]", task)
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
	for i := 0; i < customWC.maxCoroutine; i++ {
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
					customWC.cntWait.Add(-1)
					chanMsg <- fmt.Sprintf("正在处理, 当前数据[%s]", currentTransform.data)

					currentTransform.flagIsTransforming = true
					currentTransform.transformedAt = carbon.Now()
					customWC.listInMemory.LoadOrStore(currentTransform.id, currentTransform)

					if err := apiWhisper.Transform(currentTransform.toWhisperTransformParams()); err != nil {
						currentLog.Errorf("%s, 数据[%s]", currentTransform.data.String(), err)
						chanMsg <- fmt.Sprintf("%s, 数据[%s]", err, currentTransform.data.String())
						customWC.listInMemory.Delete(currentTransform.id)
						customWC.cntList.Add(-1)
						customWC.audioFilesInMemory.Delete(currentTransform.data.PathAudioFile)
						continue
					}

					currentLog.Infof("转换成功, 当前数据[%s], 耗时: %d", currentTransform.data, carbon.Now().DiffAbsInSeconds(timeStarted))
					chanMsg <- fmt.Sprintf("转换成功, 耗时: %d, 当前数据[%s]", carbon.Now().DiffAbsInSeconds(timeStarted), currentTransform.data)

					customWC.listInMemory.Delete(currentTransform.id)
					customWC.cntList.Add(-1)
					customWC.audioFilesInMemory.Delete(currentTransform.data.PathAudioFile)
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
				if chanMsgRedirect == nil {
					continue
				}
				chanMsgRedirect <- fmt.Sprintf("[%s]%s", carbon.Now(), currentMsg)
			}
		}
	}(customWC.ctx, customWC.chanMsg, customWC.chanMsgRedirect)
}

func (customWC *WhisperCron) init() {
	customWC.maxCoroutine = 1
	customWC.maxChanTransform = 50
	customWC.chanTransform = make(chan *WhisperData, customWC.maxChanTransform)
	customWC.chanMsg = make(chan string, 50)
	customWC.cntList = new(atomic.Int32)
	customWC.cntWait = new(atomic.Int32)
	customWC.timeStarted = carbon.Now()
}

func (customWC *WhisperCron) log() *log.Log {
	return log.Api()
}
