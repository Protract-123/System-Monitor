package debug

import (
	"fmt"

	"github.com/mappu/miqt/qt6"
)

func DumpQObjectTree(obj *qt6.QObject, indent int) {
	prefix := ""
	for i := 0; i < indent; i++ {
		prefix += "  "
	}

	name := obj.ObjectName()
	if name == "" {
		name = "<unnamed>"
	}

	fmt.Printf("%s%s (%s)\n",
		prefix,
		obj.MetaObject().ClassName(),
		name,
	)

	for _, child := range obj.Children() {
		DumpQObjectTree(child, indent+1)
	}
}

func AddDebugBorder(widget *qt6.QWidget, color string, width int) {
	styleString := fmt.Sprintf("QWidget {border: %dpx solid %s;}", width, color)
	widget.SetStyleSheet(styleString)
}

func AddRootDebugBorder(widget *qt6.QApplication, color string, width int) {
	styleString := fmt.Sprintf("QWidget {border: %dpx solid %s;}", width, color)
	widget.SetStyleSheet(styleString)
}
