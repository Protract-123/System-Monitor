package cpu

import (
	"System_Monitor/ui"
	"fmt"

	"github.com/mappu/miqt/qt6"
)

func createCPUImage() *qt6.QWidget {
	imageObject := qt6.NewQLabel(nil)
	imageBytes, _ := icons.ReadFile("icons/AppleM1.png")

	pixmap := qt6.NewQPixmap()
	pixmap.LoadFromData5(imageBytes, "", qt6.AutoColor)
	imageObject.SetPixmap(pixmap)

	imageObject.SetAlignment(qt6.AlignCenter)
	imageObject.SetScaledContents(true)
	imageObject.SetMaximumSize2(160, 160)

	imageLayout := qt6.NewQVBoxLayout(nil)
	imageLayout.SetSpacing(0)
	imageLayout.AddWidget(imageObject.QWidget)

	imageContainer := ui.NewBorderContainer(
		nil,
		2,
		0,
		qt6.NewQColor11(0, 180, 255, 200),
	)
	imageContainer.SetLayout(imageLayout.QLayout)
	imageContainer.SetSizePolicy2(qt6.QSizePolicy__Maximum, qt6.QSizePolicy__Maximum)

	return imageContainer.QWidget
}

func createCPUInfoContainer(info info) *qt6.QWidget {
	grid := qt6.NewQGridLayout(nil)
	grid.SetContentsMargins(0, 0, 0, 0)
	grid.SetHorizontalSpacing(10)
	grid.SetVerticalSpacing(4)

	title := qt6.NewQLabel5("CPU Information", nil)
	title.SetFont(ui.HeadingFont)
	title.SetAlignment(qt6.AlignCenter)
	title.SetContentsMargins(0, 0, 0, 20)

	grid.AddWidget5(title.QWidget, 0, 0, 1, 3, qt6.AlignCenter)

	leftCol := qt6.NewQWidget(nil)
	rightCol := qt6.NewQWidget(nil)

	leftLayout := qt6.NewQVBoxLayout(nil)
	rightLayout := qt6.NewQVBoxLayout(nil)

	leftLayout.SetContentsMargins(0, 0, 0, 0)
	rightLayout.SetContentsMargins(0, 0, 0, 0)

	leftLayout.SetSpacing(4)
	rightLayout.SetSpacing(4)

	leftCol.SetLayout(leftLayout.QLayout)
	rightCol.SetLayout(rightLayout.QLayout)

	leftCol.SetSizePolicy2(qt6.QSizePolicy__Expanding, qt6.QSizePolicy__Preferred)
	rightCol.SetSizePolicy2(qt6.QSizePolicy__Expanding, qt6.QSizePolicy__Preferred)

	addRow := func(labelText string, valueText string) {
		label := qt6.NewQLabel5(labelText, nil)
		label.SetFont(ui.BoldFont)
		label.SetAlignment(qt6.AlignRight | qt6.AlignVCenter)

		value := qt6.NewQLabel5(valueText, nil)
		value.SetAlignment(qt6.AlignLeft | qt6.AlignVCenter)

		leftLayout.AddWidget(label.QWidget)
		rightLayout.AddWidget(value.QWidget)
	}

	addRow("Model", info.Model)
	addRow("Codename", info.Codename)
	addRow("Cores", fmt.Sprintf("%d Cores", info.Cores))
	addRow("Threads", fmt.Sprintf("%d Threads", info.Threads))

	verticalDivider := ui.NewDivider(1, ui.Vertical)
	grid.AddWidget3(verticalDivider.QWidget, 1, 1, 3, 1)

	grid.SetRowMinimumHeight(1, 4)
	grid.SetRowMinimumHeight(3, 4)

	grid.AddWidget2(leftCol, 2, 0)
	grid.AddWidget2(rightCol, 2, 2)

	// Force equal-width columns
	grid.SetColumnStretch(0, 1)
	grid.SetColumnStretch(2, 1)

	container := qt6.NewQWidget(nil)
	container.SetLayout(grid.QLayout)

	container.SetSizePolicy2(qt6.QSizePolicy__Preferred, qt6.QSizePolicy__Preferred)

	return container
}

func createCoreInfoGrid(info info) *qt6.QWidget {
	gridLayout := qt6.NewQGridLayout(nil)
	gridLayout.SetContentsMargins(0, 0, 0, 0)
	gridLayout.SetSpacing(0)

	boldFont := qt6.NewQFont()
	boldFont.SetWeight(qt6.QFont__DemiBold)
	boldFont.SetPointSize(14)

	gridColumnCount := 0
	gridRowCount := 0
	var spacerRowIndexes []int

	for index, coreInfo := range info.CoreTypeInfos {
		currentRow := index * 2

		coreNameLabel := qt6.NewQLabel3(fmt.Sprint(coreInfo.Name))
		coreNameLabel.SetFont(boldFont)
		coreNameLabel.SetAlignment(qt6.AlignCenter)
		coreNameLabel.SetSizePolicy2(qt6.QSizePolicy__Expanding, qt6.QSizePolicy__Preferred)
		gridLayout.AddWidget4(coreNameLabel.QWidget, currentRow, 1, qt6.AlignVCenter|qt6.AlignLeft)

		coreCountLabel := qt6.NewQLabel3(fmt.Sprintf("%dC/%dT", coreInfo.CoreCount, coreInfo.ThreadCount))
		coreCountLabel.SetSizePolicy2(qt6.QSizePolicy__Maximum, qt6.QSizePolicy__Maximum)
		coreCountLabel.SetToolTip("Cores/Threads")
		gridLayout.AddWidget4(coreCountLabel.QWidget, currentRow, 3, qt6.AlignCenter)

		currentColumn := 4

		for _, cacheInfo := range coreInfo.CacheLevelInfos {
			cacheLevelLabel := qt6.NewQLabel3(fmt.Sprintf("%d%s %s", cacheInfo.Amount, cacheInfo.Unit, cacheInfo.Name))
			cacheLevelLabel.SetSizePolicy2(qt6.QSizePolicy__Maximum, qt6.QSizePolicy__Maximum)
			gridLayout.AddWidget4(cacheLevelLabel.QWidget, currentRow, currentColumn, qt6.AlignCenter)

			currentColumn++
		}

		if currentColumn > gridColumnCount {
			gridColumnCount = currentColumn
		}
		if currentRow > gridRowCount {
			gridRowCount = currentRow
		}

		if currentRow != 0 {
			spacerRowIndexes = append(spacerRowIndexes, currentRow-1)
		}
	}

	for _, row := range spacerRowIndexes {
		horizontalDivider := ui.NewDivider(1, ui.Horizontal)
		gridLayout.AddWidget3(horizontalDivider.QWidget, row, 1, 1, gridColumnCount)
	}

	verticalDivider := ui.NewDivider(1, ui.Vertical)
	gridLayout.AddWidget3(verticalDivider.QWidget, 0, 2, gridRowCount+1, 1)

	gridLayout.SetColumnStretch(0, 1)
	gridLayout.SetColumnStretch(4, 1)
	gridLayout.SetHorizontalSpacing(5)
	gridLayout.SetVerticalSpacing(4)

	container := qt6.NewQWidget(nil)
	container.SetContentsMargins(0, 0, 0, 15)
	container.SetLayout(gridLayout.QLayout)

	return container
}
