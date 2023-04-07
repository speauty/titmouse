package widgets

import (
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func FormLineFileDialog(win walk.Form, btnHandle **walk.PushButton, echoHandle **walk.TextLabel, label, btnLabel string, isFile bool, fileFilter string) Widget {
	return Composite{Layout: HBox{MarginsZero: true}, Children: []Widget{
		TextLabel{Text: label},
		PushButton{AssignTo: btnHandle, Text: btnLabel, OnClicked: func() {
			dlg := new(walk.FileDialog)
			if fileFilter != "" {
				dlg.Filter = fileFilter
			}
			var isSelected bool
			var err error
			if isFile {
				isSelected, err = dlg.ShowOpen(win)
			} else {
				isSelected, err = dlg.ShowBrowseFolder(win)
			}
			if err != nil {
				walk.MsgBox(win, "错误", fmt.Sprintf("选择异常, 错误: %s", err), walk.MsgBoxOK)
				return
			}
			if isSelected && echoHandle != nil {
				_ = (*echoHandle).SetText(dlg.FilePath)
			}
		}},
	}}
}
