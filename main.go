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

	qt6.NewQApplication(os.Args)

	window := qt6.NewQMainWindow(nil)
	window.SetWindowTitle("MIQT Qt6 App")

	widget := ui.NewGlowImageWidget(nil, "image.png")

	window.SetCentralWidget(widget.QWidget)
	window.Resize(400, 300)
	window.Show()

	qt6.QApplication_Exec()
}
