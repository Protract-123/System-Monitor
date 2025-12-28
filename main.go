package main

import (
	"fmt"
	"os"

	"github.com/mappu/miqt/qt6"
	"golang.org/x/sys/unix"

	"System_Monitor/ui"
)

type CpuInfo struct {
	Model   string `yaml:"cpu_model"`
	Cores   uint8  `yaml:"cpu_cores"`
	Threads uint8  `yaml:"cpu_threads"`
}

func main() {
	// Getting CPU Info on macOS
	info := CpuInfo{}
	info.Model, _ = unix.Sysctl("machdep.cpu.brand_string")

	// hw.physicalcpu is the number of active cores in current power profile
	// hw.physicalcpu_max is the number of cores which the system actually has
	cores, _ := unix.SysctlUint32("hw.physicalcpu_max")
	info.Cores = uint8(cores)

	// hw.logicalcpu is the number of active threads in current power profile
	// hw.logicalcpu_max is the number of threads which the system actually has
	threads, _ := unix.SysctlUint32("hw.logicalcpu_max")
	info.Threads = uint8(threads)

	fmt.Println("CPU Info:", info)

	qApp := qt6.NewQApplication(os.Args)

	window := qt6.NewQMainWindow(nil)
	window.SetWindowTitle("MIQT Qt6 App")

	rootContainer := qt6.NewQWidget(nil)

	hLayout := qt6.NewQHBoxLayout(nil)
	hLayout.SetContentsMargins(10, 0, 10, 0)
	hLayout.SetSpacing(10)

	image := ui.NewGlowImageWidget(rootContainer, "image.png", 10)
	image.SetSizePolicy2(qt6.QSizePolicy__Maximum, qt6.QSizePolicy__Maximum)

	labelContainer := qt6.NewQWidget(nil)
	//AddDebugBorder(labelContainer, "red", 1)
	labelLayout := qt6.NewQVBoxLayout(nil)
	labelLayout.SetContentsMargins(20, 0, 20, 0)
	labelLayout.SetSpacing(2)

	boldFont := qt6.NewQFont()
	boldFont.SetWeight(qt6.QFont__Black)
	boldFont.SetPointSize(16)

	label1 := qt6.NewQLabel5("CPU Information", nil)
	label1.SetSizePolicy2(qt6.QSizePolicy__Maximum, qt6.QSizePolicy__Maximum)
	label1.SetFont(boldFont)
	label1.SetContentsMargins(0, 0, 0, 10)

	label2 := qt6.NewQLabel5(fmt.Sprintf("CPU Model: %s", info.Model), nil)
	label2.SetSizePolicy2(qt6.QSizePolicy__Maximum, qt6.QSizePolicy__Maximum)

	label3 := qt6.NewQLabel5(fmt.Sprintf("Threads/Cores: %d/%d", info.Threads, info.Cores), nil)
	label3.SetSizePolicy2(qt6.QSizePolicy__Maximum, qt6.QSizePolicy__Maximum)

	labelLayout.AddStretchWithStretch(1)
	labelLayout.AddWidget3(label1.QWidget, 0, qt6.AlignHCenter)
	labelLayout.AddWidget3(label2.QWidget, 0, qt6.AlignHCenter)
	labelLayout.AddWidget3(label3.QWidget, 0, qt6.AlignHCenter)
	labelLayout.AddStretchWithStretch(1)

	labelContainer.SetLayout(labelLayout.QLayout)

	hLayout.AddStretchWithStretch(1)
	hLayout.AddWidget(image.QWidget)
	hLayout.AddWidget(labelContainer)
	hLayout.AddStretchWithStretch(1)

	rootContainer.SetLayout(hLayout.QLayout)

	window.SetCentralWidget(rootContainer)
	window.Resize(400, 300)
	window.Show()

	fmt.Println(qApp.ObjectName())
	DumpQObjectTree(window.QObject, 0)

	qt6.QApplication_Exec()
}
