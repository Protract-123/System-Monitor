package main

import (
	"System_Monitor/cpu"
	"System_Monitor/debug"
	"fmt"
	"os"

	"github.com/mappu/miqt/qt6"
)

func main() {
	qApp := qt6.NewQApplication(os.Args)

	window := qt6.NewQMainWindow(nil)
	window.SetWindowTitle("MIQT Qt6 App")

	rootContainer := qt6.NewQWidget(nil)
	cpu.GenerateUI(rootContainer)

	window.SetCentralWidget(rootContainer)
	window.Resize(400, 300)
	window.Show()

	fmt.Println(qApp.ObjectName())
	debug.DumpQObjectTree(window.QObject, 0)

	qt6.QApplication_Exec()
}
