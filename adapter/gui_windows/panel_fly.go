package gui_windows

import (
	"errors"
	"github.com/golang-module/carbon"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"strconv"
	"titmouse/adapter/gui_windows/util"
	"titmouse/adapter/gui_windows/widgets"
	"titmouse/cron"
	"titmouse/lib/processor/whisper"
	"titmouse/model"
	"titmouse/repository"
)

func (customGW *GuiWindows) panelFlyBuilder() Widget {
	defaultCfg := repository.ApiCfg().ActionGetWhisperCfg()
	defaultGraphicIdx := -1
	if len(customGW.graphics) > 0 {
		defaultGraphicIdx = 0
	}
	for idx, graphic := range customGW.graphics {
		if graphic == defaultCfg.DefaultGraphic {
			defaultGraphicIdx = idx
			break
		}
	}

	return GroupBox{
		AssignTo: &customGW.ptrFlyPanel,
		Layout:   VBox{},
		Title:    "投递面板",
		MinSize:  Size{Width: sidePanelWidth},
		MaxSize:  Size{Width: sidePanelWidth},
		Children: []Widget{
			Composite{
				Layout: HBox{MarginsZero: true},
				Children: []Widget{
					widgets.FormLineFileDialog(
						customGW.ptrMainWindow, &customGW.ptrBtnAudioFile, &customGW.ptrEchoAudioFile,
						"目标音频", "点击选择文件", true,
						"*.*|*.mp4|*.avi|*.mp3",
					),
					HSpacer{},
				},
			},
			Composite{
				Layout: HBox{MarginsZero: true},
				Children: []Widget{
					TextLabel{AssignTo: &customGW.ptrEchoAudioFile, Font: Font{PointSize: 8}, Enabled: false},
					HSpacer{},
				},
			},
			Composite{
				Layout: HBox{MarginsZero: true},
				Children: []Widget{
					widgets.FormLineFileDialog(
						customGW.ptrMainWindow, &customGW.ptrBtnModelFile, &customGW.ptrEchoModelFile,
						"语言模型", "点击选择文件", true,
						"*.*|*.bin",
					),
					HSpacer{},
				},
			},
			Composite{
				Layout: HBox{MarginsZero: true},
				Children: []Widget{
					TextLabel{
						AssignTo: &customGW.ptrEchoModelFile, Font: Font{PointSize: 8}, Enabled: false,
						Text: defaultCfg.DefaultModelPath,
					},
					HSpacer{},
				},
			},
			Composite{
				Layout: HBox{MarginsZero: true},
				Children: []Widget{
					widgets.FormLineComboBox(&customGW.ptrGraphic, "指定显卡", customGW.graphics, defaultGraphicIdx),
					HSpacer{},
				},
			},
			Composite{
				Layout: HBox{MarginsZero: true},
				Children: []Widget{
					widgets.FormLineComboBox(&customGW.ptrLang, "音频语种", customGW.languages, 0),
					HSpacer{},
				},
			},
			Composite{
				Layout: HBox{MarginsZero: true},
				Children: []Widget{
					widgets.FormLineInput(&customGW.ptrNumThreads, "线程数量", util.Int2Str(int(defaultCfg.DefaultThread))),
					widgets.FormLineInput(&customGW.ptrNumProcessors, "处理器数", util.Int2Str(int(defaultCfg.DefaultProcessors))),
					HSpacer{},
				},
			},
			Composite{
				Layout: HBox{MarginsZero: true},
				Children: []Widget{
					PushButton{Text: "提交", OnClicked: customGW.eventSubmit},
					PushButton{Text: "保存为默认配置", OnClicked: customGW.eventSave2DefaultConfig},
					HSpacer{},
				},
			},
			VSeparator{MinSize: Size{Height: 10}, MaxSize: Size{Height: 10}},
			Composite{
				Layout: HBox{MarginsZero: true},
				Children: []Widget{
					widgets.FormLineFileDialog(
						customGW.ptrMainWindow, &customGW.ptrBtnCmdFile, &customGW.ptrEchoCmdFile,
						"Whisper文件", "点击选择文件", true,
						"*.*|*.exe",
					),
					HSpacer{},
				},
			},
			Composite{
				Layout: HBox{MarginsZero: true},
				Children: []Widget{
					TextLabel{AssignTo: &customGW.ptrEchoCmdFile, Font: Font{PointSize: 8}, Enabled: false, Text: defaultCfg.CmdPath},
					HSpacer{},
				},
			},
			Composite{
				Layout: HBox{MarginsZero: true},
				Children: []Widget{
					PushButton{Text: "更新", OnClicked: customGW.eventUpdateCfg},
					HSpacer{},
				},
			},
			/*Composite{
				Layout: HBox{MarginsZero: true},
				Children: []Widget{
					CheckBox{AssignTo: &customGW.ptrFlagCacheWhisperData, Text: "是否缓存未处理完的任务[临时]", Checked: repository.ApiCfg().ActionGetFlagCacheWhisperData(), OnCheckedChanged: func() {
						repository.ApiCfg().ActionSetFlagCacheWhisperData(customGW.ptrFlagCacheWhisperData.Checked())
					}},
					HSpacer{},
				},
			},*/
			VSeparator{MinSize: Size{Height: 10}, MaxSize: Size{Height: 10}},
			Composite{
				Layout: HBox{MarginsZero: true},
				Children: []Widget{
					TextLabel{Text: "开发人员"},
					TextLabel{Text: "speauty"},
					HSpacer{},
				},
			},
			Composite{
				Layout: HBox{MarginsZero: true},
				Children: []Widget{
					TextLabel{Text: "联系邮箱"},
					TextLabel{Text: "speauty@163.com"},
					HSpacer{},
				},
			},
			Composite{
				Layout: HBox{MarginsZero: true},
				Children: []Widget{
					TextLabel{Text: "项目地址"},
					TextLabel{Text: "https://github.com/speauty/titmouse"},
					HSpacer{},
				},
			},
			VSpacer{},
		},
	}
}

func (customGW *GuiWindows) eventSubmit() {
	audioFile, modelFile, graphic, language, numThreads, numProcessors := customGW.formFormat()

	validateChains := []func() error{
		customGW.validateChainAudioFile, customGW.validateChainModelFile,
	}
	for _, chain := range validateChains {
		if err := chain(); err != nil {
			walk.MsgBox(customGW.ptrMainWindow, "警告", err.Error(), walk.MsgBoxOK)
			return
		}
	}
	modelTask := &model.TaskModel{
		Id:             carbon.Now().Layout(carbon.ShortDateTimeLayout),
		NumThreads:     numThreads,
		NumProcessors:  numProcessors,
		GraphicAdapter: graphic,
		PathModel:      modelFile,
		PathAudioFile:  audioFile,
		Language:       language,
	}
	if err := cron.ApiWhisperCron().Push(modelTask); err != nil {
		walk.MsgBox(customGW.ptrMainWindow, "警告", err.Error(), walk.MsgBoxOK)
		return
	}
	_ = customGW.ptrEchoAudioFile.SetText("")
	return
}

func (customGW *GuiWindows) eventSave2DefaultConfig() {
	_, modelFile, graphic, language, numThreads, numProcessors := customGW.formFormat()
	if err := customGW.validateChainModelFile(); err != nil {
		walk.MsgBox(customGW.ptrMainWindow, "警告", err.Error(), walk.MsgBoxOK)
		return
	}

	defaultCfg := repository.ApiCfg().ActionGetWhisperCfg()
	defaultCfg.DefaultGraphic = graphic
	defaultCfg.DefaultModelPath = modelFile
	defaultCfg.DefaultLanguage = language
	defaultCfg.DefaultThread = uint(numThreads)
	defaultCfg.DefaultProcessors = uint(numProcessors)
	if err := repository.ApiCfg().ActionSetWhisperCfg(defaultCfg); err != nil {
		walk.MsgBox(customGW.ptrMainWindow, "警告", err.Error(), walk.MsgBoxOK)
		return
	}
	walk.MsgBox(customGW.ptrMainWindow, "提示", "已保存为默认设置", walk.MsgBoxOK)
	return
}
func (customGW *GuiWindows) eventUpdateCfg() {
	cmdFile := customGW.ptrEchoCmdFile.Text()
	if cmdFile == "" {
		walk.MsgBox(customGW.ptrMainWindow, "警告", "请选择Whisper文件", walk.MsgBoxOK)
		return
	}

	defaultCfg := repository.ApiCfg().ActionGetWhisperCfg()
	if defaultCfg.CmdPath == cmdFile {
		walk.MsgBox(customGW.ptrMainWindow, "警告", "暂无配置更新", walk.MsgBoxOK)
		return
	}
	defaultCfg.CmdPath = cmdFile

	if err := repository.ApiCfg().ActionSetWhisperCfg(defaultCfg); err != nil {
		walk.MsgBox(customGW.ptrMainWindow, "警告", err.Error(), walk.MsgBoxOK)
		return
	}
	if err := whisper.Api().SetCfg(defaultCfg); err != nil {
		walk.MsgBox(customGW.ptrMainWindow, "警告", err.Error(), walk.MsgBoxOK)
		return
	}
	walk.MsgBox(customGW.ptrMainWindow, "提示", "当前配置已经更新, 建议重启应用", walk.MsgBoxOK)
	return
}

func (customGW *GuiWindows) formFormat() (string, string, string, string, int, int) {
	audioFile := customGW.ptrEchoAudioFile.Text()
	modelFile := customGW.ptrEchoModelFile.Text()
	graphic := ""
	if len(customGW.graphics) > 0 {
		graphic = customGW.graphics[customGW.ptrGraphic.CurrentIndex()]
	}
	language := customGW.languages[customGW.ptrLang.CurrentIndex()]
	numThreads, _ := strconv.Atoi(customGW.ptrNumThreads.Text())
	numProcessors, _ := strconv.Atoi(customGW.ptrNumProcessors.Text())
	return audioFile, modelFile, graphic, language, numThreads, numProcessors
}

func (customGW *GuiWindows) validateChainAudioFile() error {
	if customGW.ptrEchoAudioFile.Text() == "" {
		return errors.New("未检测到音视频文件")
	}
	return nil
}

func (customGW *GuiWindows) validateChainModelFile() error {
	if customGW.ptrEchoModelFile.Text() == "" {
		return errors.New("未检测到语言模型文件")
	}
	return nil
}
