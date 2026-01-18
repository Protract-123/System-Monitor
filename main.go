package main

import (
	"System_Monitor/cpu"
	"System_Monitor/debug"
	"System_Monitor/memory"
	"System_Monitor/ui"
	"fmt"
	"os"

	"github.com/mappu/miqt/qt6"
	"go.yaml.in/yaml/v4"
)

func main() {
	yamlData, _ := yaml.Marshal(memory.FetchInfo())
	println(string(yamlData))

	qApp := qt6.NewQApplication(os.Args)

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
	//debug.AddDebugBorder(rootContainer.QWidget, "red", 1)

	window.SetCentralWidget(rootContainer)
	window.SetContentsMargins(10, 10, 10, 10)
	//window.Resize(400, 300)
	window.Show()

	window.SetFixedSize(window.Size())

	fmt.Println(qApp.ObjectName())
	debug.DumpQObjectTree(window.QObject, 0)

	qt6.QApplication_Exec()
}
