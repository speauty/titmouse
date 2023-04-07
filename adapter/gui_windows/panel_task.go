package gui_windows

import (
	. "github.com/lxn/walk/declarative"
)

func (customGW *GuiWindows) panelTaskBuilder() Widget {
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
