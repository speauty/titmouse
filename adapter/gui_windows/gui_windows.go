package gui_windows

import (
	"context"
	"github.com/lxn/walk"
	"github.com/lxn/win"
	"sync"
	"titmouse/lib/log"
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
	if err := customGW.ptrMainWindow.Close(); err != nil {
		customGW.log().Scope(customGW.name()).Errorf("主窗口关闭异常, 错误: %s", err)
	}
}

func (customGW *GuiWindows) Init() error {
	customGW.initProps()

	customGW.ptrMainWindow = new(walk.MainWindow)

	if err := customGW.initMainWindow(); err != nil {
		customGW.log().Scope(customGW.name()).Errorf("主窗口初始化异常, 错误: %s", err)
		return err
	}

	customGW.syncChanMsg()
	customGW.syncTask()

	customGW.ptrMainWindow.Closing().Attach(customGW.eventCustomClose)

	customGW.disableResize()

	return nil
}

func (customGW *GuiWindows) eventCustomClose(canceled *bool, reason walk.CloseReason) {
	reason = walk.CloseReasonUser
	*canceled = false
	customGW.fnCancel()
}

func (customGW *GuiWindows) disableResize() {
	win.SetWindowLong(customGW.ptrMainWindow.Handle(), win.GWL_STYLE,
		win.GetWindowLong(customGW.ptrMainWindow.Handle(), win.GWL_STYLE) & ^win.WS_MAXIMIZEBOX & ^win.WS_THICKFRAME)
}

func (customGW *GuiWindows) log() *log.Log {
	return log.Api()
}

func (customGW *GuiWindows) name() string {
	return "GUI窗体"
}
