package gui_windows

import (
	"context"
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/lxn/win"
	"sync"
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

	ptrMainWindow *walk.MainWindow

	ptrFlyPanel      *walk.GroupBox
	ptrTerminalPanel *walk.GroupBox
	ptrTaskPanel     *walk.GroupBox
}

func (customGW *GuiWindows) Run(ctx context.Context, fnCancel context.CancelFunc) {
	customGW.ctx = ctx
	customGW.fnCancel = fnCancel
	customGW.ptrMainWindow.Run()
}

func (customGW *GuiWindows) Close() {

}

func (customGW *GuiWindows) Init() error {
	if err := customGW.initMainWindow(); err != nil {
		return err
	}

	win.SetWindowLong(customGW.ptrMainWindow.Handle(), win.GWL_STYLE,
		win.GetWindowLong(customGW.ptrMainWindow.Handle(), win.GWL_STYLE) & ^win.WS_MAXIMIZEBOX & ^win.WS_THICKFRAME)
	fmt.Println(customGW.ptrMainWindow.Size())
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
					GroupBox{
						AssignTo: &customGW.ptrFlyPanel,
						Layout:   VBox{},
						Title:    "投递面板",
						MinSize:  Size{Width: sidePanelWidth},
						MaxSize:  Size{Width: sidePanelWidth},
						Children: []Widget{

							Composite{
								Layout: HBox{MarginsZero: true},
								Children: []Widget{
									TextLabel{Text: "目标音频"},
									PushButton{Text: "点击选择文件"},
									HSpacer{},
								},
							},
							Composite{
								Layout: HBox{MarginsZero: true},
								Children: []Widget{
									TextLabel{
										Font:    Font{PointSize: 8},
										Text:    "E:\\下载\\086.Camera Component UI  Game Engine series.mp4",
										Enabled: false,
									},
									HSpacer{},
								},
							},
							Composite{
								Layout: HBox{MarginsZero: true},
								Children: []Widget{
									TextLabel{Text: "语言模型"},
									PushButton{Text: "点击选择文件"},
									HSpacer{},
								},
							},
							Composite{
								Layout: HBox{MarginsZero: true},
								Children: []Widget{
									TextLabel{
										Font:    Font{PointSize: 8},
										Text:    "E:\\下载\\ggml-medium.en.bin",
										Enabled: false,
									},
									HSpacer{},
								},
							},
							Composite{
								Layout: HBox{MarginsZero: true},
								Children: []Widget{
									TextLabel{Text: "指定显卡"},
									ComboBox{Model: repository.ApiCfg().ActionGetGraphics(), CurrentIndex: 0},
									HSpacer{},
								},
							},
							Composite{
								Layout: HBox{MarginsZero: true},
								Children: []Widget{
									TextLabel{Text: "音频语种"},
									ComboBox{Model: []string{"en"}, CurrentIndex: 0},
									HSpacer{},
								},
							},
							Composite{
								Layout: HBox{MarginsZero: true},
								Children: []Widget{
									TextLabel{Text: "线程数量"},
									LineEdit{Text: fmt.Sprintf("%d", repository.ApiCfg().ActionGetWhisperCfg().DefaultThread)},
									TextLabel{Text: "处理器数"},
									LineEdit{Text: fmt.Sprintf("%d", repository.ApiCfg().ActionGetWhisperCfg().DefaultProcessors)},
									HSpacer{},
								},
							},

							Composite{
								Layout: HBox{MarginsZero: true},
								Children: []Widget{
									PushButton{Text: "提交", Enabled: false},
									PushButton{Text: "保存为默认配置", Enabled: false},
									HSpacer{},
								},
							},

							VSpacer{}, // 类似竖直弹簧
						},
					},
					GroupBox{
						AssignTo: &customGW.ptrTerminalPanel,
						MinSize:  Size{Width: terminalPanelWidth},
						MaxSize:  Size{Width: terminalPanelWidth},
						Layout:   VBox{},
						Children: []Widget{
							TextLabel{Text: "欢迎使用 Titmouse 任务平台, 基于 Whisper.exe 处理音频转文本"},
							TextLabel{Text: "测试"},
							HSpacer{},
							VSpacer{},
						},
					},
					GroupBox{
						AssignTo: &customGW.ptrTaskPanel,
						Layout:   VBox{},
						Title:    "任务面板",
						MinSize:  Size{Width: sidePanelWidth},
						MaxSize:  Size{Width: sidePanelWidth},
						Children: []Widget{
							TextLabel{Text: "测试"},
							VSpacer{},
						},
					},
				},
			},
		},
	}.Create()
}

func (customGW *GuiWindows) eventExit() {

}
