package memory

import (
	"fmt"
	"time"

	"github.com/mappu/miqt/qt-restricted-extras/charts6"
	"github.com/mappu/miqt/qt6"
	"github.com/mappu/miqt/qt6/mainthread"
)

func GenerateUI() *qt6.QLayout {
	info := FetchInfo()

	memoryLayout := qt6.NewQHBoxLayout2()
	memoryLayout.SetContentsMargins(10, 10, 10, 10)

	memoryInfoContainer, memoryInfoUpdateFunc := CreateMemoryInfoContainer(info)
	memoryChart, memoryChartUpdateFunc := CreateMemoryGraphContainer(info)

	memoryLayout.AddWidget(memoryInfoContainer)
	memoryLayout.AddWidget(memoryChart)

	go func() {
		secondTicker := time.NewTicker(time.Second)
		minuteTicker := time.NewTicker(time.Minute)

		for {
			select {
			case <-secondTicker.C:
				if err := info.UpdateInfo(&info); err != nil {
					return
				}
				memoryInfoUpdateFunc(&info)

			case <-minuteTicker.C:
				memoryChartUpdateFunc(&info)
			}
		}
	}()

	return memoryLayout.QLayout
}

func CreateMemoryInfoContainer(info Info) (*qt6.QWidget, func(info2 *Info)) {
	boldFont := qt6.NewQFont()
	boldFont.SetWeight(qt6.QFont__DemiBold)
	boldFont.SetPointSize(13)

	titleFont := qt6.NewQFont()
	titleFont.SetWeight(qt6.QFont__Black)
	titleFont.SetPointSize(16)

	// Root layout
	rootLayout := qt6.NewQVBoxLayout(nil)
	rootLayout.SetContentsMargins(5, 15, 0, 15)
	rootLayout.SetSpacing(6)

	// Title
	title := qt6.NewQLabel5("Memory Breakdown", nil)
	title.SetFont(titleFont)
	title.SetAlignment(qt6.AlignTop | qt6.AlignHCenter)
	title.SetContentsMargins(0, 0, 0, 10)

	rootLayout.AddWidget(title.QWidget)

	// Helper to create a row
	addRow := func(labelText, valueText string) *qt6.QLabel {
		row := qt6.NewQHBoxLayout(nil)
		row.SetContentsMargins(0, 0, 0, 0)
		row.SetSpacing(30)

		label := qt6.NewQLabel5(labelText, nil)
		label.SetFont(boldFont)
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

	// Helper to add a horizontal divider
	addDivider := func() {
		divider := qt6.NewQFrame(nil)
		divider.SetFrameShape(qt6.QFrame__HLine)
		divider.SetFrameShadow(qt6.QFrame__Plain)
		divider.SetLineWidth(1)
		rootLayout.AddWidget(divider.QWidget)
	}

	// Rows
	addRow("System Memory", fmt.Sprintf("%d%s", info.TotalMemory.Value, info.TotalMemory.Unit))
	addDivider()
	usedMemoryLabel := addRow("Used Memory", fmt.Sprintf("%.2f%s", info.UsedMemory.Value, info.UsedMemory.Unit))
	freeMemoryLabel := addRow("Free Memory", fmt.Sprintf("%.2f%s", info.FreeMemory.Value, info.FreeMemory.Unit))
	addRow("Usable Memory", fmt.Sprintf("%.2f%s", info.UsableMemory.Value, info.UsableMemory.Unit))

	addDivider()
	swapUsedLabel := addRow("Used Swap", fmt.Sprintf("%.2f%s", info.SwapUsed.Value, info.SwapUsed.Unit))
	swapFreeLabel := addRow("Free Swap", fmt.Sprintf("%.2f%s", info.SwapFree.Value, info.SwapFree.Unit))
	swapTotalLabel := addRow("Total Swap", fmt.Sprintf("%d%s", info.SwapTotal.Value, info.SwapTotal.Unit))

	// Update Func
	updateFunc := func(info *Info) {
		mainthread.Start(func() {
			usedMemoryLabel.SetText(fmt.Sprintf("%.2f%s", info.UsedMemory.Value, info.UsedMemory.Unit))
			freeMemoryLabel.SetText(fmt.Sprintf("%.2f%s", info.FreeMemory.Value, info.FreeMemory.Unit))
			swapUsedLabel.SetText(fmt.Sprintf("%.2f%s", info.SwapUsed.Value, info.SwapUsed.Unit))
			swapFreeLabel.SetText(fmt.Sprintf("%.2f%s", info.SwapFree.Value, info.SwapFree.Unit))
			swapTotalLabel.SetText(fmt.Sprintf("%d%s", info.SwapTotal.Value, info.SwapTotal.Unit))
		})
	}

	container := qt6.NewQWidget(nil)
	container.SetLayout(rootLayout.QLayout)
	container.SetSizePolicy2(qt6.QSizePolicy__Preferred, qt6.QSizePolicy__Preferred)

	return container, updateFunc
}

func CreateMemoryGraphContainer(info Info) (*qt6.QWidget, func(info2 *Info)) {
	memoryChartView := charts6.NewQChartView2()
	memoryChart := charts6.NewQChart()

	memoryUsedValues := [30]float32{}
	memoryUsedValues[29] = info.UsedMemory.Value

	memoryUsedLineSeries := charts6.NewQLineSeries()
	memoryUsedLineSeries.SetColor(qt6.NewQColor3(186, 225, 255))

	memoryUsedScatterSeries := charts6.NewQScatterSeries()
	memoryUsedScatterSeries.OnHovered(func(point *qt6.QPointF, state bool) {
		if !state {
			qt6.QToolTip_HideText()
			return
		}

		text := fmt.Sprintf(
			"Time: %.0f\nMemory: %.2f",
			point.X(),
			point.Y(),
		)

		qt6.QToolTip_ShowText(
			qt6.QCursor_Pos(),
			text,
		)
	})
	memoryUsedScatterSeries.SetMarkerSize(3)
	memoryUsedScatterSeries.SetColor(qt6.NewQColor3(186, 225, 255))

	memoryChartXAxis := charts6.NewQValueAxis()
	memoryChartXAxis.SetMin(0)
	memoryChartXAxis.SetMax(30)
	memoryChartXAxis.SetTitleText("Time")

	memoryChartYAxis := charts6.NewQValueAxis()
	memoryChartYAxis.SetMin(0)
	memoryChartYAxis.SetMax(float64(info.UsableMemory.Value))
	memoryChartYAxis.SetTitleText("Memory Used")

	memoryChart.AddSeries(memoryUsedLineSeries.QAbstractSeries)
	memoryChart.AddSeries(memoryUsedScatterSeries.QAbstractSeries)

	memoryChart.AddAxis(memoryChartXAxis.QAbstractAxis, qt6.AlignBottom)
	memoryChart.AddAxis(memoryChartYAxis.QAbstractAxis, qt6.AlignRight)

	memoryUsedLineSeries.AttachAxis(memoryChartXAxis.QAbstractAxis)
	memoryUsedLineSeries.AttachAxis(memoryChartYAxis.QAbstractAxis)

	memoryUsedScatterSeries.AttachAxis(memoryChartXAxis.QAbstractAxis)
	memoryUsedScatterSeries.AttachAxis(memoryChartYAxis.QAbstractAxis)

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
			copy(memoryUsedValues[0:], memoryUsedValues[1:])
			memoryUsedValues[len(memoryUsedValues)-1] = info2.UsedMemory.Value

			var points [30]qt6.QPointF
			for i := 0; i < len(memoryUsedValues); i++ {
				// Removing non (0,0) causes crashes, for what reason I have no idea
				// Perhaps QT doesn't like non-continuous lines?
				points[i] = *qt6.NewQPointF3(float64(i), float64(memoryUsedValues[i]))
			}

			memoryUsedLineSeries.ReplaceWithPoints(points[:])
			memoryUsedScatterSeries.ReplaceWithPoints(points[:])
		})
	}

	memoryChartView.OnChangeEvent(func(super func(*qt6.QEvent), e *qt6.QEvent) {
		if e.Type() == qt6.QEvent__PaletteChange {
			applyChartPalette(memoryChart, memoryChartXAxis.QAbstractAxis, memoryChartYAxis.QAbstractAxis)
		}
	})

	return memoryChartView.QWidget, updateFunc
}

func applyChartPalette(
	chart *charts6.QChart,
	xAxis *charts6.QAbstractAxis,
	yAxis *charts6.QAbstractAxis,
) {
	p := qt6.QGuiApplication_Palette()

	text := p.ColorWithCr(qt6.QPalette__Text)
	window := p.ColorWithCr(qt6.QPalette__Window)
	base := p.ColorWithCr(qt6.QPalette__Base)
	grid := p.ColorWithCr(qt6.QPalette__Mid)

	// Chart backgrounds
	chart.SetBackgroundBrush(qt6.NewQBrush3(window))
	chart.SetPlotAreaBackgroundVisibleWithVisible(true)
	chart.SetPlotAreaBackgroundBrush(qt6.NewQBrush3(base))

	// Chart title
	chart.SetTitleBrush(qt6.NewQBrush3(text))

	//chart.Legend().SetBrush(qt6.NewQBrush3(window))
	chart.Legend().SetLabelBrush(qt6.NewQBrush3(text))

	// Axes text
	xAxis.SetLabelsBrush(qt6.NewQBrush3(text))
	xAxis.SetTitleBrush(qt6.NewQBrush3(text))
	xAxis.SetGridLineColor(grid)

	yAxis.SetLabelsBrush(qt6.NewQBrush3(text))
	yAxis.SetTitleBrush(qt6.NewQBrush3(text))
	yAxis.SetGridLineColor(grid)
}
