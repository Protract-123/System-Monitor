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

	cpuContainer := qt6.NewQWidget2()
	cpuContainer.SetContentsMargins(0, 0, 0, 0)
	cpuLayout := cpu.GenerateUI()
	cpuContainer.SetLayout(cpuLayout)

	//utils.AddDebugBorder(cpuContainer, "red", 1)

	memoryContainer := qt6.NewQWidget2()
	memoryLayout := memory.GenerateUI()
	memoryContainer.SetLayout(memoryLayout)

	//utils.AddDebugBorder(memoryContainer, "red", 1)

	var divider = ui.NewDivider(1, ui.Vertical)
	rootLayout.SetSpacing(10)

	rootLayout.AddStretchWithStretch(1)
	rootLayout.AddWidget(cpuContainer)
	rootLayout.AddWidget(divider.QWidget)
	rootLayout.AddWidget(memoryContainer)
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
