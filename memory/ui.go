package memory

import (
	"System_Monitor/ui"
	"fmt"
	"time"

	"github.com/mappu/miqt/qt-restricted-extras/charts6"
	"github.com/mappu/miqt/qt6"
	"github.com/mappu/miqt/qt6/mainthread"
)

func createMemoryInfoContainer(memoryInfo info) (*qt6.QWidget, func(memoryInfo *info)) {
	rootLayout := qt6.NewQVBoxLayout(nil)
	rootLayout.SetContentsMargins(0, 0, 0, 0)
	rootLayout.SetSpacing(6)

	title := qt6.NewQLabel5("Memory Breakdown", nil)
	title.SetFont(ui.HeadingFont)
	title.SetAlignment(qt6.AlignCenter)
	title.SetContentsMargins(0, 0, 0, 0)

	rootLayout.AddWidget(title.QWidget)

	addRow := func(labelText, valueText string) *qt6.QLabel {
		row := qt6.NewQHBoxLayout(nil)
		row.SetContentsMargins(0, 0, 0, 0)
		row.SetSpacing(30)

		label := qt6.NewQLabel5(labelText, nil)
		label.SetFont(ui.BoldFont)
		label.SetAlignment(qt6.AlignLeft | qt6.AlignVCenter)

		value := qt6.NewQLabel5(valueText, nil)
		value.SetAlignment(qt6.AlignRight | qt6.AlignVCenter)

		label.SetSizePolicy2(qt6.QSizePolicy__Expanding, qt6.QSizePolicy__Preferred)
		value.SetSizePolicy2(qt6.QSizePolicy__Expanding, qt6.QSizePolicy__Preferred)

		row.AddWidget(label.QWidget)
		row.AddWidget(value.QWidget)

		rootLayout.AddLayout(row.QLayout)

		return value
	}

	addDivider := func() {
		divider := ui.NewDivider(1, ui.Horizontal)
		rootLayout.AddWidget(divider.QWidget)
	}

	addRow("System Memory", fmt.Sprintf("%d%s", memoryInfo.TotalMemory.Value, memoryInfo.TotalMemory.Unit))
	addDivider()
	addRow("Usable Memory", fmt.Sprintf("%.2f%s", memoryInfo.UsableMemory.Value, memoryInfo.UsableMemory.Unit))

	usedMemoryLabel := addRow("•  Used Memory", fmt.Sprintf("%.2f%s", memoryInfo.UsedMemory.Value, memoryInfo.UsedMemory.Unit))
	freeMemoryLabel := addRow("•  Free Memory", fmt.Sprintf("%.2f%s", memoryInfo.FreeMemory.Value, memoryInfo.FreeMemory.Unit))

	addDivider()
	swapTotalLabel := addRow("Total Swap", fmt.Sprintf("%d%s", memoryInfo.SwapTotal.Value, memoryInfo.SwapTotal.Unit))
	swapUsedLabel := addRow("•  Used Swap", fmt.Sprintf("%.2f%s", memoryInfo.SwapUsed.Value, memoryInfo.SwapUsed.Unit))
	swapFreeLabel := addRow("•  Free Swap", fmt.Sprintf("%.2f%s", memoryInfo.SwapFree.Value, memoryInfo.SwapFree.Unit))

	updateFunc := func(memoryInfo *info) {
		mainthread.Start(func() {
			usedMemoryLabel.SetText(fmt.Sprintf("%.2f%s", memoryInfo.UsedMemory.Value, memoryInfo.UsedMemory.Unit))
			freeMemoryLabel.SetText(fmt.Sprintf("%.2f%s", memoryInfo.FreeMemory.Value, memoryInfo.FreeMemory.Unit))
			swapUsedLabel.SetText(fmt.Sprintf("%.2f%s", memoryInfo.SwapUsed.Value, memoryInfo.SwapUsed.Unit))
			swapFreeLabel.SetText(fmt.Sprintf("%.2f%s", memoryInfo.SwapFree.Value, memoryInfo.SwapFree.Unit))
			swapTotalLabel.SetText(fmt.Sprintf("%d%s", memoryInfo.SwapTotal.Value, memoryInfo.SwapTotal.Unit))
		})
	}

	container := qt6.NewQWidget(nil)
	container.SetLayout(rootLayout.QLayout)
	container.SetSizePolicy2(qt6.QSizePolicy__Preferred, qt6.QSizePolicy__Preferred)

	return container, updateFunc
}

func createMemoryGraphContainer(memoryInfo info) (*qt6.QWidget, func(memoryInfo *info)) {
	memoryChartView := charts6.NewQChartView2()
	memoryChart := charts6.NewQChart()

	memoryUsedValues := [30]float32{}
	timestamps := [30]string{}

	memoryUsedLine := charts6.NewQLineSeries()
	memoryUsedLine.SetColor(qt6.NewQColor3(186, 225, 255))

	memoryUsedMarkers := charts6.NewQScatterSeries()
	memoryUsedMarkers.OnHovered(func(point *qt6.QPointF, state bool) {
		if !state {
			qt6.QToolTip_HideText()
			return
		}

		text := fmt.Sprintf(
			"Time: %s\nMemory Used: %.2f",
			timestamps[int(point.X())],
			point.Y(),
		)

		qt6.QToolTip_ShowText(
			qt6.QCursor_Pos(),
			text,
		)
	})
	memoryUsedMarkers.SetMarkerSize(4)
	memoryUsedMarkers.SetBorderColor(qt6.NewQColor3(0, 0, 0))
	memoryUsedMarkers.SetColor(qt6.NewQColor3(186, 225, 255))

	memoryChartXAxis := charts6.NewQValueAxis()
	memoryChartXAxis.SetMin(0)
	memoryChartXAxis.SetMax(30)
	memoryChartXAxis.SetTitleText("Time")

	memoryChartYAxis := charts6.NewQValueAxis()
	memoryChartYAxis.SetMin(0)
	memoryChartYAxis.SetMax(float64(memoryInfo.UsableMemory.Value))
	memoryChartYAxis.SetTitleText("Memory Used")
	//memoryChartYAxis.SetLabelFormat(fmt.Sprintf("%%.2f %s", info.TotalMemory.Unit))

	memoryChart.AddSeries(memoryUsedLine.QAbstractSeries)
	memoryChart.AddSeries(memoryUsedMarkers.QAbstractSeries)

	memoryChart.AddAxis(memoryChartXAxis.QAbstractAxis, qt6.AlignBottom)
	memoryChart.AddAxis(memoryChartYAxis.QAbstractAxis, qt6.AlignRight)

	memoryUsedLine.AttachAxis(memoryChartXAxis.QAbstractAxis)
	memoryUsedLine.AttachAxis(memoryChartYAxis.QAbstractAxis)

	memoryUsedMarkers.AttachAxis(memoryChartXAxis.QAbstractAxis)
	memoryUsedMarkers.AttachAxis(memoryChartYAxis.QAbstractAxis)

	memoryChart.SetMargins(qt6.NewQMargins2(0, 0, 0, 0))

	memoryChartXAxis.Hide()
	memoryChart.Legend().Hide()
	memoryChartYAxis.SetTitleVisibleWithVisible(false)

	memoryChartView.SetChart(memoryChart)
	memoryChartView.SetRenderHint2(qt6.QPainter__Antialiasing, true)
	memoryChartView.SetFixedSize(qt6.NewQSize2(250, 200))

	ui.ApplyChartPalette(memoryChart, memoryChartXAxis.QAbstractAxis, memoryChartYAxis.QAbstractAxis)

	updateFunc := func(memoryInfo *info) {
		mainthread.Start(func() {
			copy(memoryUsedValues[0:], memoryUsedValues[1:])
			memoryUsedValues[len(memoryUsedValues)-1] = memoryInfo.UsedMemory.Value

			copy(timestamps[0:], timestamps[1:])
			timestamps[len(timestamps)-1] = time.Now().Format("15:04")

			var points [30]qt6.QPointF
			for i := 0; i < len(memoryUsedValues); i++ {
				// Removing non (0,0) causes crashes, for what reason I have no idea
				// Perhaps QT doesn't like non-continuous lines?
				points[i] = *qt6.NewQPointF3(float64(i), float64(memoryUsedValues[i]))
			}

			memoryUsedLine.ReplaceWithPoints(points[:])
			memoryUsedMarkers.ReplaceWithPoints(points[:])
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
