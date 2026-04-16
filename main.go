// System Monitor - a desktop utility exposing CPU, memory, and disk information.
// Copyright (C) 2026 Yuvraj Verma
//
// This program is free software: you can redistribute it and/or modify it under
// the terms of the GNU General Public License as published by the Free Software
// Foundation, either version 3 of the License, or (at your option) any later
// version.
//
// This program is distributed in the hope that it will be useful, but WITHOUT
// ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS
// FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along with
// this program. If not, see <https://www.gnu.org/licenses/>.

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

	cpuMemoryDivider := ui.NewDivider(1, ui.Vertical)
	memoryDiskDivider := ui.NewDivider(1, ui.Vertical)
	rootLayout.SetSpacing(10)

	rootLayout.AddStretchWithStretch(1)
	rootLayout.AddWidget(cpuContainer)
	rootLayout.AddWidget(cpuMemoryDivider.QWidget)
	rootLayout.AddWidget(memoryContainer)
	rootLayout.AddWidget(memoryDiskDivider.QWidget)
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
