package ui

import (
	"github.com/mappu/miqt/qt-restricted-extras/charts6"
	"github.com/mappu/miqt/qt6"
)

func ApplyChartPalette(
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

	chart.Legend().SetBrush(qt6.NewQBrush3(window))
	chart.Legend().SetLabelBrush(qt6.NewQBrush3(text))

	// Axes text
	xAxis.SetLabelsBrush(qt6.NewQBrush3(text))
	xAxis.SetTitleBrush(qt6.NewQBrush3(text))
	xAxis.SetGridLineColor(grid)

	yAxis.SetLabelsBrush(qt6.NewQBrush3(text))
	yAxis.SetTitleBrush(qt6.NewQBrush3(text))
	yAxis.SetGridLineColor(grid)
}
