package main

import (
	"System_Monitor/cpu"
	"System_Monitor/memory"
	"System_Monitor/ui"
	"System_Monitor/utils"
	"fmt"
	"os"

	"github.com/mappu/miqt/qt6"
)

func main() {
	qApp := qt6.NewQApplication(os.Args)
	ui.LoadFonts()

	window := qt6.NewQMainWindow(nil)
	window.SetWindowTitle("MIQT Qt6 App")

	rootContainer := qt6.NewQWidget2()
	rootLayout := qt6.NewQHBoxLayout2()

	cpuContainer := ui.NewBorderContainer(nil, 2, 8, qt6.NewQColor11(255, 255, 255, 255))
	cpuLayout := cpu.GenerateUI()
	cpuContainer.SetLayout(cpuLayout)

	memoryContainer := ui.NewBorderContainer(nil, 2, 8, qt6.NewQColor11(255, 255, 255, 255))
	memoryLayout := memory.GenerateUI()
	memoryContainer.SetLayout(memoryLayout)

	rootLayout.AddStretchWithStretch(1)
	rootLayout.AddWidget(cpuContainer.QWidget)
	rootLayout.AddWidget(memoryContainer.QWidget)
	rootLayout.AddStretchWithStretch(1)

	rootContainer.SetLayout(rootLayout.QLayout)
	rootContainer.SetSizePolicy2(qt6.QSizePolicy__Maximum, qt6.QSizePolicy__Maximum)

	window.SetCentralWidget(rootContainer)
	window.SetContentsMargins(10, 10, 10, 10)
	window.Show()

	window.SetFixedSize(window.Size())

	fmt.Println(qApp.ObjectName())
	utils.DumpQObjectTree(window.QObject, 0)

	qt6.QApplication_Exec()
}
