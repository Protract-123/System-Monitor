//go:build darwin

package cpu

import (
	"System_Monitor/utils"
	"fmt"
	"strings"

	"github.com/mappu/miqt/qt6"
	"golang.org/x/sys/unix"
)

func GenerateUI() *qt6.QLayout {
	cpuInfo := fetchInfo()

	cpuLayout := qt6.NewQGridLayout(nil)
	cpuLayout.SetContentsMargins(0, 0, 0, 0)
	cpuLayout.SetVerticalSpacing(10)
	cpuLayout.SetHorizontalSpacing(10)

	imageContainer := createCPUImage()
	cpuInfoContainer := createCPUInfoContainer(cpuInfo)
	coreInfoGrid := createCoreInfoGrid(cpuInfo)

	cpuLayout.AddWidget4(imageContainer, 0, 0, qt6.AlignCenter)
	cpuLayout.AddWidget4(cpuInfoContainer, 0, 1, qt6.AlignVCenter|qt6.AlignTop)
	cpuLayout.AddWidget5(coreInfoGrid, 1, 0, 1, 2, qt6.AlignCenter)

	return cpuLayout.QLayout
}

func fetchInfo() info {
	cpuInfo := info{}

	cpuInfo.Model, _ = unix.Sysctl("machdep.cpu.brand_string")

	cores, _ := unix.SysctlUint32("hw.physicalcpu_max")
	cpuInfo.Cores = uint(cores)

	threads, _ := unix.SysctlUint32("hw.logicalcpu_max")
	cpuInfo.Threads = uint(threads)

	coreTypeCount, _ := unix.SysctlUint32("hw.nperflevels")

	cpuInfo.Codename = cpuToCodename[strings.ToLower(cpuInfo.Model)]

	for i := 0; i < int(coreTypeCount); i++ {
		coreTypeInfo := coreInfo{}
		coreTypeString := fmt.Sprintf("hw.perflevel%d", i)

		coreTypeInfo.Name, _ = unix.Sysctl(fmt.Sprintf("%s.name", coreTypeString))

		coreTypeAmount, _ := unix.SysctlUint32(fmt.Sprintf("%s.physicalcpu_max", coreTypeString))
		coreTypeInfo.CoreCount = uint(coreTypeAmount)

		coreTypeThreadAmount, _ := unix.SysctlUint32(fmt.Sprintf("%s.logicalcpu_max", coreTypeString))
		coreTypeInfo.ThreadCount = uint(coreTypeThreadAmount)

		cacheTypes := [3]string{"L1i", "L1d", "L2"}

		for _, cacheType := range cacheTypes {
			cacheKey := fmt.Sprintf("%s.%scachesize", coreTypeString, strings.ToLower(cacheType))

			cacheSize, _ := unix.SysctlUint32(cacheKey)

			cacheLevelInfo := cacheInfo{}
			cacheLevelInfo.Name = cacheType

			var convertedCacheSize uint

			convertedCacheSize, cacheLevelInfo.Unit = utils.ConvertFromBytesParts(uint(cacheSize))
			cacheLevelInfo.Amount = convertedCacheSize

			coreTypeInfo.CacheLevelInfos = append(coreTypeInfo.CacheLevelInfos, cacheLevelInfo)
		}

		cpuInfo.CoreTypeInfos = append(cpuInfo.CoreTypeInfos, coreTypeInfo)
	}

	return cpuInfo
}
