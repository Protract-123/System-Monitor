package main

import (
	"System_Monitor/cpu"
	"System_Monitor/disk"
	"System_Monitor/memory"
	"System_Monitor/ui"
	"os"

	"github.com/mappu/miqt/qt6"
)

func main() {
	qt6.NewQApplication(os.Args)
	ui.LoadFonts()

	window := qt6.NewQMainWindow(nil)
	window.SetWindowTitle("System Monitor")

	rootContainer := qt6.NewQWidget2()
	rootLayout := qt6.NewQHBoxLayout2()

	cpuContainer := qt6.NewQWidget2()
	cpuLayout := cpu.GenerateUI()
	cpuContainer.SetLayout(cpuLayout)

	memoryContainer := qt6.NewQWidget2()
	memoryLayout := memory.GenerateUI()
	memoryContainer.SetLayout(memoryLayout)

	diskContainer := qt6.NewQWidget2()
	diskLayout := disk.GenerateUI()
	diskContainer.SetLayout(diskLayout)

	var divider = ui.NewDivider(1, ui.Vertical)
	var divider2 = ui.NewDivider(1, ui.Vertical)
	rootLayout.SetSpacing(10)

	rootLayout.AddStretchWithStretch(1)
	rootLayout.AddWidget(cpuContainer)
	rootLayout.AddWidget(divider.QWidget)
	rootLayout.AddWidget(memoryContainer)
	rootLayout.AddWidget(divider2.QWidget)
	rootLayout.AddWidget(diskContainer)
	rootLayout.AddStretchWithStretch(1)

	rootContainer.SetLayout(rootLayout.QLayout)
	rootContainer.SetSizePolicy2(qt6.QSizePolicy__Maximum, qt6.QSizePolicy__Maximum)

	window.SetCentralWidget(rootContainer)
	window.SetContentsMargins(10, 10, 10, 10)
	window.Show()

	window.SetFixedSize(window.Size())

	qt6.QApplication_Exec()
}
