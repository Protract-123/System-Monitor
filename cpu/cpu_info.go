package cpu

import (
	"fmt"
	"strings"

	"golang.org/x/sys/unix"
)

type Info struct {
	Model         string     `yaml:"cpu_model"`
	Cores         uint8      `yaml:"cpu_cores"`
	Threads       uint8      `yaml:"cpu_threads"`
	Codename      string     `yaml:"cpu_codename"`
	CoreTypeCount uint8      `yaml:"cpu_core_info_count"`
	CoreTypeInfos []CoreInfo `yaml:"cpu_core_infos"`
}

type CoreInfo struct {
	Name            string      `yaml:"core_name"`
	CoreCount       uint8       `yaml:"core_count"`
	ThreadCount     uint8       `yaml:"core_thread_count"`
	CacheLevelCount uint8       `yaml:"core_cache_levels"`
	CacheLevelInfos []CacheInfo `yaml:"core_cache_infos"`
}

type CacheInfo struct {
	Name   string `yaml:"cache_name"`
	Amount uint16 `yaml:"cache_amount"`
	Unit   string `yaml:"cache_unit"`
}

var cpuToCodename = map[string]string{
	// Retrieved from https://en.wikipedia.org/wiki/List_of_Apple_codenames#M_series
	// Named after islands, pretty cool naming scheme tbh
	"apple m1": "Tonga",
	"apple m2": "Staten",
	"apple m3": "Ibiza",
	"apple m4": "Donan",
	"apple m5": "Hidra",
}

func FetchInfoOSX() Info {
	info := Info{}

	info.Model, _ = unix.Sysctl("machdep.cpu.brand_string")

	// hw.physicalcpu is the number of active cores in current power profile
	// hw.physicalcpu_max is the number of cores which the system actually has
	cores, _ := unix.SysctlUint32("hw.physicalcpu_max")
	info.Cores = uint8(cores)

	// hw.logicalcpu is the number of active threads in current power profile
	// hw.logicalcpu_max is the number of threads which the system actually has
	threads, _ := unix.SysctlUint32("hw.logicalcpu_max")
	info.Threads = uint8(threads)

	coreTypeCount, _ := unix.SysctlUint32("hw.nperflevels")
	info.CoreTypeCount = uint8(coreTypeCount)

	info.Codename = cpuToCodename[strings.ToLower(info.Model)]

	for i := 0; i < int(coreTypeCount); i++ {
		coreTypeInfo := CoreInfo{}
		coreTypeString := fmt.Sprintf("hw.perflevel%d", i)

		coreTypeInfo.Name, _ = unix.Sysctl(fmt.Sprintf("%s.name", coreTypeString))

		coreTypeAmount, _ := unix.SysctlUint32(fmt.Sprintf("%s.physicalcpu_max", coreTypeString))
		coreTypeInfo.CoreCount = uint8(coreTypeAmount)

		coreTypeThreadAmount, _ := unix.SysctlUint32(fmt.Sprintf("%s.logicalcpu_max", coreTypeString))
		coreTypeInfo.ThreadCount = uint8(coreTypeThreadAmount)

		cacheTypes := [3]string{"L1i", "L1d", "L2"}
		coreTypeInfo.CacheLevelCount = 3

		for _, cacheType := range cacheTypes {
			cacheKey := fmt.Sprintf("%s.%scachesize", coreTypeString, strings.ToLower(cacheType))

			cacheSize, _ := unix.SysctlUint32(cacheKey)

			cacheInfo := CacheInfo{}
			cacheInfo.Name = cacheType
			cacheInfo.Unit = "Bytes"

			// add hw.perflevel1.cpusperl2 to get cores per l2

			cacheUnits := [3]string{"KiB", "MiB", "GiB"}
			i := 0
			for cacheSize >= 1024 {
				cacheSize = cacheSize / 1024
				cacheInfo.Unit = cacheUnits[i]
				i += 1
			}
			cacheInfo.Amount = uint16(cacheSize)

			coreTypeInfo.CacheLevelInfos = append(coreTypeInfo.CacheLevelInfos, cacheInfo)
		}

		info.CoreTypeInfos = append(info.CoreTypeInfos, coreTypeInfo)
	}

	return info
}
