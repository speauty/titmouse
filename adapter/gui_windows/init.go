package gui_windows

import (
	. "github.com/lxn/walk/declarative"
	"titmouse/cron"
	"titmouse/lib/processor/whisper"
	"titmouse/repository"
)

func (customGW *GuiWindows) initProps() {

	customGW.graphics = repository.ApiCfg().ActionGetGraphics()
	customGW.languages = whisper.LanguagesSupported

	customGW.chanMsg = make(chan string, 20)
	cron.ApiWhisperCron().SetChanMsgRedirect(customGW.chanMsg)
}

func (customGW *GuiWindows) initMainWindow() error {
	return MainWindow{
		AssignTo: &customGW.ptrMainWindow, Title: title, Icon: icon,
		MinSize: Size{Width: minWidth, Height: minHeight},
		MaxSize: Size{Width: width, Height: height},
		Size:    Size{Width: width, Height: height},
		Layout:  VBox{MarginsZero: true},
		Children: []Widget{
			Composite{
				Layout: HBox{Alignment: AlignHCenterVNear},
				Children: []Widget{
					customGW.panelFlyBuilder(),
					customGW.panelTerminalBuilder(),
					customGW.panelTaskBuilder(),
				},
			},
		},
	}.Create()
}
