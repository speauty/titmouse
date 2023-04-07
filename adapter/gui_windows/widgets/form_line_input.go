package widgets

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func FormLineInput(handle **walk.LineEdit, label, value string) Widget {
	return Composite{
		Layout: HBox{MarginsZero: true},
		Children: []Widget{
			TextLabel{Text: label},
			LineEdit{AssignTo: handle, Text: value},
		},
	}
}
