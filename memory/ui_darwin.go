//go:build darwin

package memory

import (
	"fmt"
	"time"

	"github.com/mappu/miqt/qt-restricted-extras/charts6"
	"github.com/mappu/miqt/qt6"
	"github.com/mappu/miqt/qt6/mainthread"
)

func CreateMemoryAreaGraph(info Info) (*qt6.QWidget, func(info2 *Info)) {
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
	memoryChartYAxis.SetMax(float64(info.UsableMemory.Value))
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

	applyChartPalette(memoryChart, memoryChartXAxis.QAbstractAxis, memoryChartYAxis.QAbstractAxis)

	updateFunc := func(info2 *Info) {
		mainthread.Start(func() {
			platformInfo, ok := info2.PlatformInfo.(DarwinPlatformInfo)
			if !ok {
				return
			}

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

			updateArray(&compressedMemoryValues, platformInfo.CompressedMemory.Value)
			updateArray(&wiredMemoryValues, platformInfo.WiredMemory.Value)
			updateArray(&appMemoryValues, platformInfo.AppMemory.Value)

			rollingSum := [30]float32{}

			updateData(&rollingSum, &wiredMemoryValues, levelFour, levelFourMarkers)
			updateData(&rollingSum, &appMemoryValues, levelThree, levelThreeMarkers)
			updateData(&rollingSum, &compressedMemoryValues, levelTwo, levelTwoMarkers)
		})
	}

	memoryChartView.OnChangeEvent(func(super func(*qt6.QEvent), e *qt6.QEvent) {
		if e.Type() == qt6.QEvent__PaletteChange {
			applyChartPalette(memoryChart, memoryChartXAxis.QAbstractAxis, memoryChartYAxis.QAbstractAxis)
		}
	})

	updateFunc(&info)

	return memoryChartView.QWidget, updateFunc
}
