package gui_windows

import (
	"context"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/lxn/win"
	"sync"
	"titmouse/cron"
	"titmouse/repository"
)

var (
	apiGuiWindows  *GuiWindows
	onceGuiWindows sync.Once
)

const (
	title              = "titmouse@speauty"
	icon               = "favicon.ico"
	width              = 1200
	height             = 675
	minWidth           = 800
	minHeight          = 450
	sidePanelWidth     = 250
	terminalPanelWidth = 650
)

func ApiGuiWindows() *GuiWindows {
	onceGuiWindows.Do(func() {
		apiGuiWindows = new(GuiWindows)
	})
	return apiGuiWindows
}

type GuiWindows struct {
	ctx      context.Context
	fnCancel context.CancelFunc

	graphics  []string
	languages []string

	ptrMainWindow *walk.MainWindow

	ptrFlyPanel      *walk.GroupBox
	ptrTerminalPanel *walk.GroupBox
	ptrTaskPanel     *walk.GroupBox

	ptrBtnAudioFile  *walk.PushButton
	ptrEchoAudioFile *walk.TextLabel

	ptrBtnModelFile  *walk.PushButton
	ptrEchoModelFile *walk.TextLabel

	ptrGraphic *walk.ComboBox
	ptrLang    *walk.ComboBox

	ptrNumThreads    *walk.LineEdit
	ptrNumProcessors *walk.LineEdit

	ptrBtnCmdFile  *walk.PushButton
	ptrEchoCmdFile *walk.TextLabel

	chanMsg   chan string
	ptrOutput *walk.ListBox
	outputMsg []string

	ptrTaskLisInMemory *walk.ListBox

	ptrTimeRun *walk.TextLabel
	ptrCntTask *walk.TextLabel
}

func (customGW *GuiWindows) Run(ctx context.Context, fnCancel context.CancelFunc) {
	customGW.ctx = ctx
	customGW.fnCancel = fnCancel
	customGW.ptrMainWindow.Run()
}

func (customGW *GuiWindows) Close() {

}

func (customGW *GuiWindows) Init() error {
	customGW.graphics = repository.ApiCfg().ActionGetGraphics()
	customGW.languages = []string{"en"}
	customGW.chanMsg = make(chan string, 20)
	cron.ApiWhisperCron().SetChanMsgRedirect(customGW.chanMsg)

	go func() {
		for true {
			select {
			case msg := <-customGW.chanMsg:
				if msg != "" {
					customGW.outputMsg = append([]string{msg}, customGW.outputMsg...)
					if customGW.ptrOutput != nil {
						_ = customGW.ptrOutput.SetModel(customGW.outputMsg)
						customGW.ptrOutput.RequestLayout()
					}
				}
			}
		}
	}()

	customGW.ptrMainWindow = new(walk.MainWindow)
	if err := customGW.initMainWindow(); err != nil {
		return err
	}

	win.SetWindowLong(customGW.ptrMainWindow.Handle(), win.GWL_STYLE,
		win.GetWindowLong(customGW.ptrMainWindow.Handle(), win.GWL_STYLE) & ^win.WS_MAXIMIZEBOX & ^win.WS_THICKFRAME)
	return nil
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

func (customGW *GuiWindows) eventExit() {

}
