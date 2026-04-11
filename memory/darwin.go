//go:build darwin

package memory

/*
   #include <mach/mach.h>
*/
import "C"
import (
	"System_Monitor/ui"
	"System_Monitor/utils"
	"encoding/binary"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"time"
	"unsafe"

	"github.com/mappu/miqt/qt-restricted-extras/charts6"
	"github.com/mappu/miqt/qt6"
	"github.com/mappu/miqt/qt6/mainthread"
	"golang.org/x/sys/unix"
)

func GenerateUI() *qt6.QLayout {
	memoryInfo := fetchInfo()

	memoryLayout := qt6.NewQHBoxLayout2()
	memoryLayout.SetContentsMargins(10, 10, 10, 10)

	memoryInfoContainer, memoryInfoUpdateFunc := createMemoryInfoContainer(memoryInfo.info)

	//memoryChart, memoryChartUpdateFunc := CreateMemoryGraphContainer(info)
	memoryChart, memoryChartUpdateFunc := createMemoryAreaGraph(memoryInfo)

	memoryLayout.AddWidget(memoryInfoContainer)
	memoryLayout.AddWidget(memoryChart)

	go func() {
		secondTicker := time.NewTicker(time.Second)
		minuteTicker := time.NewTicker(time.Minute)

		for {
			select {
			case <-secondTicker.C:
				if err := updateInfo(&memoryInfo); err != nil {
					return
				}
				memoryInfoUpdateFunc(&memoryInfo.info)

			case <-minuteTicker.C:
				memoryChartUpdateFunc(&memoryInfo)
			}
		}
	}()

	return memoryLayout.QLayout
}

func createMemoryAreaGraph(memoryInfo darwinInfo) (*qt6.QWidget, func(memoryInfo *darwinInfo)) {
	memoryChartView := charts6.NewQChartView2()
	memoryChart := charts6.NewQChart()

	compressedMemoryValues := [30]float32{} // 2nd Level
	appMemoryValues := [30]float32{}        // 3rd Level
	wiredMemoryValues := [30]float32{}      // 4th Level
	// Bottom level is X Axis

	timestamps := [30]string{}

	generateLevel := func(color *qt6.QColor, dataArray *[30]float32, hoverTitle string) (*charts6.QLineSeries, *charts6.QScatterSeries) {
		level := charts6.NewQLineSeries()
		level.SetColor(color)

		levelMarkers := charts6.NewQScatterSeries()
		levelMarkers.OnHovered(func(point *qt6.QPointF, state bool) {
			if !state {
				qt6.QToolTip_HideText()
				return
			}

			text := fmt.Sprintf(
				"Time: %s\n%s: %.2f",
				timestamps[int(point.X())],
				hoverTitle,
				dataArray[int(point.X())],
			)

			qt6.QToolTip_ShowText(
				qt6.QCursor_Pos(),
				text,
			)
		})
		levelMarkers.SetMarkerSize(4)
		levelMarkers.SetBorderColor(qt6.NewQColor3(0, 0, 0))
		levelMarkers.SetColor(color)

		return level, levelMarkers

	}

	levelTwo, levelTwoMarkers := generateLevel(qt6.NewQColor3(186, 225, 255), &compressedMemoryValues, "Compressed Memory")
	levelThree, levelThreeMarkers := generateLevel(qt6.NewQColor3(186, 225, 255), &appMemoryValues, "App Memory")
	levelFour, levelFourMarkers := generateLevel(qt6.NewQColor3(186, 225, 255), &wiredMemoryValues, "Wired Memory")

	areaTwo := charts6.NewQAreaSeries4(levelTwo, levelThree)
	areaThree := charts6.NewQAreaSeries4(levelThree, levelFour)
	areaFour := charts6.NewQAreaSeries2(levelFour)

	memoryChartXAxis := charts6.NewQValueAxis()
	memoryChartXAxis.SetMin(0)
	memoryChartXAxis.SetMax(30)
	memoryChartXAxis.SetTitleText("Time")

	memoryChartYAxis := charts6.NewQValueAxis()
	memoryChartYAxis.SetMin(0)
	memoryChartYAxis.SetMax(float64(memoryInfo.UsableMemory.Value))
	memoryChartYAxis.SetTitleText("Memory Used")
	//memoryChartYAxis.SetLabelFormat(fmt.Sprintf("%%.2f %s", info.TotalMemory.Unit))

	addSeries := func(series *charts6.QAbstractSeries) {
		memoryChart.AddSeries(series)

		series.AttachAxis(memoryChartXAxis.QAbstractAxis)
		series.AttachAxis(memoryChartYAxis.QAbstractAxis)
	}

	memoryChart.AddAxis(memoryChartXAxis.QAbstractAxis, qt6.AlignBottom)
	memoryChart.AddAxis(memoryChartYAxis.QAbstractAxis, qt6.AlignRight)

	addSeries(areaTwo.QAbstractSeries)
	addSeries(areaThree.QAbstractSeries)
	addSeries(areaFour.QAbstractSeries)

	addSeries(levelTwoMarkers.QAbstractSeries)
	addSeries(levelThreeMarkers.QAbstractSeries)
	addSeries(levelFourMarkers.QAbstractSeries)

	memoryChart.SetMargins(qt6.NewQMargins2(0, 0, 0, 0))

	memoryChartXAxis.Hide()
	memoryChart.Legend().Hide()
	memoryChartYAxis.SetTitleVisibleWithVisible(false)

	memoryChartView.SetChart(memoryChart)
	memoryChartView.SetRenderHint2(qt6.QPainter__Antialiasing, true)
	memoryChartView.SetFixedSize(qt6.NewQSize2(250, 200))

	ui.ApplyChartPalette(memoryChart, memoryChartXAxis.QAbstractAxis, memoryChartYAxis.QAbstractAxis)

	updateFunc := func(memoryInfo *darwinInfo) {
		mainthread.Start(func() {
			copy(timestamps[0:], timestamps[1:])
			timestamps[len(timestamps)-1] = time.Now().Format("15:04")

			updateArray := func(array *[30]float32, newValue float32) {
				copy(array[0:], array[1:])
				array[len(array)-1] = newValue
			}

			updateData := func(sum *[30]float32, array *[30]float32, series1 *charts6.QLineSeries, series2 *charts6.QScatterSeries) {
				var points [30]qt6.QPointF
				for i := 0; i < len(sum); i++ {
					// Removing non (0,0) causes crashes, for what reason I have no idea
					// Perhaps QT doesn't like non-continuous lines?
					sum[i] += array[i]
					points[i] = *qt6.NewQPointF3(float64(i), float64(sum[i]))
				}

				series1.ReplaceWithPoints(points[:])
				series2.ReplaceWithPoints(points[:])
			}

			updateArray(&compressedMemoryValues, memoryInfo.CompressedMemory.Value)
			updateArray(&wiredMemoryValues, memoryInfo.WiredMemory.Value)
			updateArray(&appMemoryValues, memoryInfo.AppMemory.Value)

			rollingSum := [30]float32{}

			updateData(&rollingSum, &wiredMemoryValues, levelFour, levelFourMarkers)
			updateData(&rollingSum, &appMemoryValues, levelThree, levelThreeMarkers)
			updateData(&rollingSum, &compressedMemoryValues, levelTwo, levelTwoMarkers)
		})
	}

	memoryChartView.OnChangeEvent(func(super func(*qt6.QEvent), e *qt6.QEvent) {
		if e.Type() == qt6.QEvent__PaletteChange {
			ui.ApplyChartPalette(memoryChart, memoryChartXAxis.QAbstractAxis, memoryChartYAxis.QAbstractAxis)
		}
	})

	updateFunc(&memoryInfo)

	return memoryChartView.QWidget, updateFunc
}

type platformInfo struct {
	AppMemory        utils.ValueUnitPair[float32] `yaml:"app_memory"`
	CompressedMemory utils.ValueUnitPair[float32] `yaml:"compressed_memory"`
	WiredMemory      utils.ValueUnitPair[float32] `yaml:"wired_memory"`
	CachedMemory     utils.ValueUnitPair[float32] `yaml:"cached_memory"`
	MemoryPressure   utils.ValueUnitPair[uint]    `yaml:"memory_pressure"`
}

type darwinInfo struct {
	info
	platformInfo
}

func fetchInfo() darwinInfo {
	memoryInfo := darwinInfo{}

	totalMemory, _ := unix.SysctlUint64("hw.memsize")
	memoryInfo.TotalMemory = utils.ConvertFromBytes(uint(totalMemory))

	usableMemory, _ := unix.SysctlUint64("hw.memsize_usable")
	memoryInfo.UsableMemory = utils.ConvertFromBytes(float32(usableMemory))

	err := addSwapStats(&memoryInfo)
	if err != nil {
		return darwinInfo{}
	}

	err = addVMStats(&memoryInfo)
	if err != nil {
		return darwinInfo{}
	}

	err = addSystemPressure(&memoryInfo)
	if err != nil {
		return darwinInfo{}
	}

	return memoryInfo
}

func updateInfo(info *darwinInfo) error {

	err := addSwapStats(info)
	if err != nil {
		return err
	}

	err = addVMStats(info)
	if err != nil {
		return err
	}

	err = addSystemPressure(info)
	if err != nil {
		return err
	}

	return nil
}

func addVMStats(info *darwinInfo) error {
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
		return fmt.Errorf("host_statistics64 failed: %d", int(ret))
	}

	formatMemoryValue := func(value uint) utils.ValueUnitPair[float32] {
		unit := "Pages"

		if info.PageSize.Value == 0 {
			return utils.ValueUnitPair[float32]{
				Value: float32(value),
				Unit:  unit,
			}
		}

		return utils.ConvertFromBytes[float32](float32(value * info.PageSize.Value))
	}

	info.AppMemory = formatMemoryValue(uint(stats.internal_page_count) - uint(stats.purgeable_count))
	info.CompressedMemory = formatMemoryValue(uint(stats.compressor_page_count))
	info.WiredMemory = formatMemoryValue(uint(stats.wire_count))
	info.CachedMemory = formatMemoryValue(uint(stats.external_page_count) + uint(stats.speculative_count))

	// freePages + usedPages approximately equals hw.memsize_usable, as it should
	freePages := uint(stats.external_page_count) + uint(stats.free_count)
	usedPages := uint(stats.internal_page_count) - uint(stats.purgeable_count) + uint(stats.compressor_page_count) + uint(stats.wire_count)

	info.FreeMemory = formatMemoryValue(freePages)
	info.UsedMemory = formatMemoryValue(usedPages)

	return nil
}

func addSwapStats(info *darwinInfo) error {
	rawSwapUsage, err := unix.SysctlRaw("vm.swapusage")
	if err != nil || len(rawSwapUsage) < 32 {
		return err
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

	info.PageSize = utils.ValueUnitPair[uint]{
		Value: uint(pageSize),
		Unit:  "Bytes",
	}

	info.SwapUsed = utils.ConvertFromBytes(float32(swapUsed))
	info.SwapFree = utils.ConvertFromBytes(float32(swapAvailable))
	info.SwapTotal = utils.ConvertFromBytes(uint(totalSwap))

	return nil
}

func addSystemPressure(darwinInfo *darwinInfo) error {
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
