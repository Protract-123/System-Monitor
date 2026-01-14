//go:build darwin

package memory

/*
   #include <mach/mach.h>
*/
import "C"
import (
	"System_Monitor/utils"
	"encoding/binary"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"unsafe"

	"golang.org/x/sys/unix"
)

type DarwinPlatformInfo struct {
	AppMemory        utils.ValueUnitPair[float32] `yaml:"app_memory"`
	CompressedMemory utils.ValueUnitPair[float32] `yaml:"compressed_memory"`
	WiredMemory      utils.ValueUnitPair[float32] `yaml:"wired_memory"`
	CachedMemory     utils.ValueUnitPair[float32] `yaml:"cached_memory"`

	MemoryPressure utils.ValueUnitPair[uint] `yaml:"memory_pressure"`
}

func (DarwinPlatformInfo) isPlatformInfo() {}

func FetchInfo() Info {
	info := Info{}
	platformInfo := DarwinPlatformInfo{}

	totalMemory, _ := unix.SysctlUint64("hw.memsize")
	info.TotalMemory = utils.ConvertFromBytes(uint(totalMemory))

	usableMemory, _ := unix.SysctlUint64("hw.memsize_usable")
	info.UsableMemory = utils.ConvertFromBytes(float32(usableMemory))

	rawSwapUsage, err := unix.SysctlRaw("vm.swapusage")
	if err != nil || len(rawSwapUsage) < 32 {
		return Info{}
	}

	// Bit Layout for vm.swapusage as seen in sysctl.h
	// 0–7   uint64 xsu_total
	// 8–15  uint64 xsu_avail
	// 16–23 uint64 xsu_used
	// 24–27 uint32 xsu_pagesize
	// 28–31 uint32 xsu_encrypted

	totalSwap := binary.LittleEndian.Uint64(rawSwapUsage[0:8])
	swapAvailable := binary.LittleEndian.Uint64(rawSwapUsage[8:16])
	swapUsed := binary.LittleEndian.Uint64(rawSwapUsage[16:24])

	pageSize := binary.LittleEndian.Uint32(rawSwapUsage[24:28])

	info.SwapUsed = utils.ConvertFromBytes(float32(swapUsed))
	info.SwapFree = utils.ConvertFromBytes(float32(swapAvailable))
	info.SwapTotal = utils.ConvertFromBytes(uint(totalSwap))

	err = AddVMStats(&info, &platformInfo, uint(pageSize))
	if err != nil {
		return Info{}
	}

	err = AddSystemPressure(&platformInfo)
	if err != nil {
		return Info{}
	}

	info.PlatformInfo = platformInfo

	return info
}

func AddVMStats(info *Info, darwinInfo *DarwinPlatformInfo, pageSize uint) error {
	var stats C.vm_statistics64_data_t
	var count C.mach_msg_type_number_t = C.HOST_VM_INFO64_COUNT

	host := C.mach_host_self()

	ret := C.host_statistics64(
		host,
		C.HOST_VM_INFO64,
		C.host_info64_t(unsafe.Pointer(&stats)),
		&count,
	)
	if ret != C.KERN_SUCCESS {
		return fmt.Errorf("host_statistics64 failed: %d", ret)
	}

	formatMemoryValue := func(value uint) utils.ValueUnitPair[float32] {
		unit := "Pages"

		if pageSize == 0 {
			return utils.ValueUnitPair[float32]{
				Value: float32(value),
				Unit:  unit,
			}
		}

		return utils.ConvertFromBytes[float32](float32(value * pageSize))
	}

	darwinInfo.AppMemory = formatMemoryValue(uint(stats.internal_page_count) - uint(stats.purgeable_count))
	darwinInfo.CompressedMemory = formatMemoryValue(uint(stats.compressor_page_count))
	darwinInfo.WiredMemory = formatMemoryValue(uint(stats.wire_count))
	darwinInfo.CachedMemory = formatMemoryValue(uint(stats.external_page_count) + uint(stats.speculative_count))

	// freePages + usedPages approximately equals hw.memsize_usable, as it should
	freePages := uint(stats.external_page_count) + uint(stats.free_count)
	usedPages := uint(stats.internal_page_count) - uint(stats.purgeable_count) + uint(stats.compressor_page_count) + uint(stats.wire_count)

	info.FreeMemory = formatMemoryValue(freePages)
	info.UsedMemory = formatMemoryValue(usedPages)

	return nil
}

func AddSystemPressure(darwinInfo *DarwinPlatformInfo) error {
	cmd := exec.Command("memory_pressure")
	out, err := cmd.Output()
	if err != nil {
		return err
	}

	re := regexp.MustCompile(`System-wide memory free percentage:\s+(\d+)%`)
	matches := re.FindSubmatch(out)
	if len(matches) < 2 {
		return fmt.Errorf("percentage not found")
	}

	// This is system-wide memory free percentage, not memory pressure percentage
	percent, err := strconv.Atoi(string(matches[1]))
	if err != nil {
		return err
	}

	darwinInfo.MemoryPressure = utils.ValueUnitPair[uint]{
		Value: 100 - uint(percent),
		Unit:  "Percent",
	}

	/*
		Memory Pressure Boundaries
		System-wide memory free percentage less than 30% = RED
		System-wide memory free percentage less than 40% = YELLOW
		otherwise GREEN
	*/

	return nil
}
