package gui_windows

import (
	"fmt"
	"time"
	"titmouse/cron"
)

func (customGW *GuiWindows) syncChanMsg() {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				customGW.log().Scope(customGW.name()).Errorf("协程-同步输出信息异常退出, 错误: %s", err)
			}
		}()
		for true {
			select {
			case msg := <-customGW.chanMsg:
				if msg != "" {
					customGW.outputMsg = append([]string{msg}, customGW.outputMsg...)
					if customGW.ptrOutput != nil {
						customGW.ptrOutput.Synchronize(func() {
							_ = customGW.ptrOutput.SetModel(customGW.outputMsg)
						})
					}
				}
			}
		}
	}()
}

func (customGW *GuiWindows) syncTask() {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				customGW.log().Scope(customGW.name()).Errorf("协程-同步任务信息异常退出, 错误: %s", err)
			}
		}()
		apiCron := cron.ApiWhisperCron()
		timeTicker := time.NewTicker(time.Second)
		for true {
			select {
			case <-timeTicker.C:
				if customGW.ptrTimeRun != nil {
					customGW.ptrTimeRun.Synchronize(func() {
						_ = customGW.ptrTimeRun.SetText(fmt.Sprintf("运行时长(s) %d", apiCron.GetTimeRun()))
					})
				}
				if customGW.ptrCntTask != nil {
					customGW.ptrCntTask.Synchronize(func() {
						_ = customGW.ptrCntTask.SetText(fmt.Sprintf(
							"任务总数 %d    任务排队 %d",
							apiCron.GetCntTask(), apiCron.GetCntWait(),
						))
					})
				}
				if customGW.ptrTaskLisInMemory != nil {
					customGW.ptrTaskLisInMemory.Synchronize(func() {
						_ = customGW.ptrTaskLisInMemory.SetModel(apiCron.GetListInMemory())
					})
				}
			}
		}
	}()
}
