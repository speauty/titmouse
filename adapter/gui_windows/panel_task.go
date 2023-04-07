package gui_windows

import (
	"fmt"
	. "github.com/lxn/walk/declarative"
	"time"
	"titmouse/cron"
)

func (customGW *GuiWindows) panelTaskBuilder() Widget {
	customGW.autoUpdater()
	return GroupBox{
		AssignTo: &customGW.ptrTaskPanel,
		Layout:   VBox{},
		Title:    "任务面板",
		MinSize:  Size{Width: sidePanelWidth},
		MaxSize:  Size{Width: sidePanelWidth},
		Children: []Widget{
			TextLabel{AssignTo: &customGW.ptrTimeRun, Text: "运行时长(s) 0"},
			TextLabel{AssignTo: &customGW.ptrCntTask, Text: "任务总数 0    任务排队 0"},
			VSeparator{},
			ListBox{
				AssignTo: &customGW.ptrTaskLisInMemory,
				MinSize:  Size{Height: height - 139},
				MaxSize:  Size{Height: height - 139},
			},
			VSpacer{},
		},
	}
}

func (customGW *GuiWindows) autoUpdater() {
	go func() {
		apiCron := cron.ApiWhisperCron()
		timeTicker := time.NewTicker(time.Second)
		for true {
			select {
			case <-timeTicker.C:
				if customGW.ptrTimeRun != nil {
					_ = customGW.ptrTimeRun.SetText(fmt.Sprintf("运行时长(s) %d", apiCron.GetTimeRun()))
				}
				if customGW.ptrCntTask != nil {
					customGW.ptrCntTask.Synchronize(func() {
						_ = customGW.ptrCntTask.SetText(fmt.Sprintf(
							"任务总数 %d    任务排队 %d",
							apiCron.GetCntTask(),
							apiCron.GetCntWait(),
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
