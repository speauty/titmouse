package gui_windows

import . "github.com/lxn/walk/declarative"

func (customGW *GuiWindows) panelTerminalBuilder() Widget {

	return GroupBox{
		AssignTo: &customGW.ptrTerminalPanel,
		MinSize:  Size{Width: terminalPanelWidth},
		MaxSize:  Size{Width: terminalPanelWidth},
		Layout:   VBox{},
		Children: []Widget{
			TextLabel{Text: "欢迎使用 Titmouse 任务平台, 基于 Whisper.exe 处理音频转文本"},
			VSeparator{},
			ListBox{
				AssignTo: &customGW.ptrOutput,
				MinSize:  Size{Height: height - 120},
				MaxSize:  Size{Height: height - 120},
				Model:    customGW.outputMsg,
			},

			VSpacer{},
		},
	}
}
