//go:build darwin

package cpu

import (
	"System_Monitor/utils"
	"fmt"
	"log"
	"strings"

	"github.com/mappu/miqt/qt6"
	"golang.org/x/sys/unix"
)

// cacheLevels are the cache tiers macOS exposes via hw.perflevelN.<name>cachesize
var cacheLevels = [...]string{"L1i", "L1d", "L2"}

// maxCoreTypes caps the number of performance levels enumerated
const maxCoreTypes = 16

func GenerateUI() *qt6.QLayout {
	cpuInfo := fetchInfo()

	cpuLayout := qt6.NewQGridLayout(nil)
	cpuLayout.SetContentsMargins(0, 0, 0, 0)
	cpuLayout.SetVerticalSpacing(10)
	cpuLayout.SetHorizontalSpacing(10)

	imageContainer := createCPUImage(cpuInfo.IconPath)
	cpuInfoContainer := createCPUInfoContainer(cpuInfo)
	coreInfoGrid := createCoreInfoGrid(cpuInfo)

	cpuLayout.AddWidget4(imageContainer, 0, 0, qt6.AlignCenter)
	cpuLayout.AddWidget4(cpuInfoContainer, 0, 1, qt6.AlignVCenter|qt6.AlignTop)
	cpuLayout.AddWidget5(coreInfoGrid, 1, 0, 1, 2, qt6.AlignCenter)

	return cpuLayout.QLayout
}

func fetchInfo() info {
	cpuInfo := info{}

	model, err := unix.Sysctl("machdep.cpu.brand_string")
	if err != nil {
		log.Printf("cpu: sysctl %q failed: %v", "machdep.cpu.brand_string", err)
	}
	cpuInfo.Model = model

	cores, err := unix.SysctlUint32("hw.physicalcpu_max")
	if err != nil {
		log.Printf("cpu: sysctl %q failed: %v", "hw.physicalcpu_max", err)
	}
	cpuInfo.Cores = uint(cores)

	threads, err := unix.SysctlUint32("hw.logicalcpu_max")
	if err != nil {
		log.Printf("cpu: sysctl %q failed: %v", "hw.logicalcpu_max", err)
	}
	cpuInfo.Threads = uint(threads)

	coreTypeCount, err := unix.SysctlUint32("hw.nperflevels")
	if err != nil {
		log.Printf("cpu: sysctl %q failed: %v", "hw.nperflevels", err)
	}

	cpuInfo.Codename = lookupBaseModel(cpuInfo.Model, cpuToCodename)
	cpuInfo.IconPath = lookupBaseModel(cpuInfo.Model, cpuToFilePath)

	if coreTypeCount > maxCoreTypes {
		log.Printf("cpu: hw.nperflevels reported %d core types, capping at %d", coreTypeCount, maxCoreTypes)
		coreTypeCount = maxCoreTypes
	}

	for i := 0; i < int(coreTypeCount); i++ {
		coreTypeInfo := coreInfo{}
		coreTypeKey := fmt.Sprintf("hw.perflevel%d", i)

		name, err := unix.Sysctl(fmt.Sprintf("%s.name", coreTypeKey))
		if err != nil {
			log.Printf("cpu: sysctl %q failed: %v", coreTypeKey+".name", err)
		}
		coreTypeInfo.Name = name

		coreTypeAmount, err := unix.SysctlUint32(fmt.Sprintf("%s.physicalcpu_max", coreTypeKey))
		if err != nil {
			log.Printf("cpu: sysctl %q failed: %v", coreTypeKey+".physicalcpu_max", err)
		}
		coreTypeInfo.CoreCount = uint(coreTypeAmount)

		coreTypeThreadAmount, err := unix.SysctlUint32(fmt.Sprintf("%s.logicalcpu_max", coreTypeKey))
		if err != nil {
			log.Printf("cpu: sysctl %q failed: %v", coreTypeKey+".logicalcpu_max", err)
		}
		coreTypeInfo.ThreadCount = uint(coreTypeThreadAmount)

		for _, cacheName := range cacheLevels {
			cacheKey := fmt.Sprintf("%s.%scachesize", coreTypeKey, strings.ToLower(cacheName))

			cacheSize, err := unix.SysctlUint32(cacheKey)
			if err != nil {
				log.Printf("cpu: sysctl %q failed: %v", cacheKey, err)
			}

			cacheEntry := cacheInfo{Name: cacheName}
			cacheEntry.Amount, cacheEntry.Unit = utils.ConvertFromBytesParts(uint(cacheSize))

			coreTypeInfo.CacheLevelInfos = append(coreTypeInfo.CacheLevelInfos, cacheEntry)
		}

		cpuInfo.CoreTypeInfos = append(cpuInfo.CoreTypeInfos, coreTypeInfo)
	}

	return cpuInfo
}
