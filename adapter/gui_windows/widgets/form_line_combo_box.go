package widgets

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func FormLineComboBox(handle **walk.ComboBox, label string, boxModel interface{}, defaultIdx int) Widget {
	return Composite{Layout: HBox{MarginsZero: true}, Children: []Widget{
		TextLabel{Text: label},
		ComboBox{AssignTo: handle, Model: boxModel, CurrentIndex: defaultIdx},
	}}
}
