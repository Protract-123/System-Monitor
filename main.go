package main

import (
	"System_Monitor/cpu"
	"System_Monitor/debug"
	"System_Monitor/ui"
	"fmt"
	"os"

	"github.com/mappu/miqt/qt6"
)

func main() {
	qApp := qt6.NewQApplication(os.Args)

	window := qt6.NewQMainWindow(nil)
	window.SetWindowTitle("MIQT Qt6 App")

	rootContainer := ui.NewBorderContainer(nil, 2, 8, qt6.NewQColor11(255, 255, 255, 255))
	cpuLayout := cpu.GenerateUI()
	rootContainer.SetLayout(cpuLayout)
	rootContainer.SetSizePolicy2(qt6.QSizePolicy__Maximum, qt6.QSizePolicy__Maximum)
	//debug.AddDebugBorder(rootContainer.QWidget, "red", 1)

	window.SetCentralWidget(rootContainer.QWidget)
	window.SetContentsMargins(10, 10, 10, 10)
	//window.Resize(400, 300)
	window.Show()

	fmt.Println(qApp.ObjectName())
	debug.DumpQObjectTree(window.QObject, 0)

	qt6.QApplication_Exec()
}
