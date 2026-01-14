//go:build darwin

package cpu

import (
	"System_Monitor/utils"
	"fmt"
	"strings"

	"golang.org/x/sys/unix"
)

func FetchInfo() Info {
	info := Info{}

	info.Model, _ = unix.Sysctl("machdep.cpu.brand_string")

	// hw.physicalcpu is the number of active cores in current power profile
	// hw.physicalcpu_max is the number of cores which the system actually has
	cores, _ := unix.SysctlUint32("hw.physicalcpu_max")
	info.Cores = uint(cores)

	// hw.logicalcpu is the number of active threads in current power profile
	// hw.logicalcpu_max is the number of threads which the system actually has
	threads, _ := unix.SysctlUint32("hw.logicalcpu_max")
	info.Threads = uint(threads)

	coreTypeCount, _ := unix.SysctlUint32("hw.nperflevels")

	info.Codename = cpuToCodename[strings.ToLower(info.Model)]

	for i := 0; i < int(coreTypeCount); i++ {
		coreTypeInfo := CoreInfo{}
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

			cacheInfo := CacheInfo{}
			cacheInfo.Name = cacheType

			var convertedCacheSize uint

			convertedCacheSize, cacheInfo.Unit = utils.ConvertFromBytesParts(uint(cacheSize))
			cacheInfo.Amount = convertedCacheSize

			coreTypeInfo.CacheLevelInfos = append(coreTypeInfo.CacheLevelInfos, cacheInfo)
		}

		info.CoreTypeInfos = append(info.CoreTypeInfos, coreTypeInfo)
	}

	return info
}
