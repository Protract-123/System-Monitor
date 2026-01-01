package cpu

import (
	"System_Monitor/ui"
	"fmt"

	"github.com/mappu/miqt/qt6"
	"go.yaml.in/yaml/v4"
)

func GenerateUI() *qt6.QLayout {
	info := FetchInfoOSX()

	// Print cpu info as yaml
	yamlData, _ := yaml.Marshal(info)
	fmt.Println(string(yamlData))

	// cpuLayout is the root layout for CPU UI
	cpuLayout := qt6.NewQGridLayout(nil)
	cpuLayout.SetContentsMargins(10, 10, 10, 10)
	cpuLayout.SetVerticalSpacing(10)
	cpuLayout.SetHorizontalSpacing(5)

	// Create UI Components
	imageContainer := CreateCPUImage()
	cpuInfoContainer := CreateCPUInfoContainer(info)
	coreInfoGrid := CreateCoreInfoGrid(info)
	cacheInfoGrid := CreateCacheInfoGrid(info)

	//cpuLayout.SetRowStretch(0, 1)
	//cpuLayout.SetColumnStretch(0, 1)
	cpuLayout.AddWidget4(imageContainer, 1, 1, qt6.AlignCenter)
	cpuLayout.AddWidget4(cpuInfoContainer, 1, 2, qt6.AlignCenter)
	cpuLayout.AddWidget4(coreInfoGrid, 2, 1, qt6.AlignCenter)
	cpuLayout.AddWidget4(cacheInfoGrid, 2, 2, qt6.AlignCenter)
	//cpuLayout.SetColumnStretch(3, 1)
	//cpuLayout.SetRowStretch(3, 1)

	return cpuLayout.QLayout
}

func CreateCPUImage() *qt6.QWidget {
	imageObject := qt6.NewQLabel(nil)
	pixmap := qt6.NewQPixmap7("./cpu/icons/AppleChipM1.png", "", qt6.AutoColor)
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

func CreateCPUInfoContainer(info Info) *qt6.QWidget {
	layout := qt6.NewQVBoxLayout(nil)

	layout.SetSpacing(5)
	layout.SetContentsMargins(0, 0, 0, 0)

	// Create bold font for styled text
	boldFont := qt6.NewQFont()
	boldFont.SetWeight(qt6.QFont__Black)
	boldFont.SetPointSize(16)

	title := qt6.NewQLabel5("CPU Information", nil)
	title.SetSizePolicy2(qt6.QSizePolicy__Maximum, qt6.QSizePolicy__Maximum)
	title.SetFont(boldFont)
	title.SetContentsMargins(0, 0, 0, 10)

	modelLabel := qt6.NewQLabel5(fmt.Sprintf("CPU Model: %s", info.Model), nil)
	modelLabel.SetSizePolicy2(qt6.QSizePolicy__Maximum, qt6.QSizePolicy__Maximum)

	coreCountLabel := qt6.NewQLabel5(fmt.Sprintf("Core Count: %d", info.Cores), nil)
	coreCountLabel.SetSizePolicy2(qt6.QSizePolicy__Maximum, qt6.QSizePolicy__Maximum)

	threadCountLabel := qt6.NewQLabel5(fmt.Sprintf("Thread Count: %d", info.Threads), nil)
	threadCountLabel.SetSizePolicy2(qt6.QSizePolicy__Maximum, qt6.QSizePolicy__Maximum)

	codenameLabel := qt6.NewQLabel5(fmt.Sprintf("Codename: %s", info.Codename), nil)
	codenameLabel.SetSizePolicy2(qt6.QSizePolicy__Maximum, qt6.QSizePolicy__Maximum)

	layout.AddStretchWithStretch(1)
	layout.AddWidget3(title.QWidget, 0, qt6.AlignHCenter)
	layout.AddWidget3(modelLabel.QWidget, 0, qt6.AlignHCenter)
	layout.AddWidget3(codenameLabel.QWidget, 0, qt6.AlignHCenter)
	layout.AddWidget3(coreCountLabel.QWidget, 0, qt6.AlignHCenter)
	layout.AddWidget3(threadCountLabel.QWidget, 0, qt6.AlignHCenter)
	layout.AddStretchWithStretch(1)

	container := qt6.NewQWidget(nil)
	container.SetLayout(layout.QLayout)

	container.SetSizePolicy2(qt6.QSizePolicy__Maximum, qt6.QSizePolicy__Maximum)

	return container
}

func CreateCoreInfoGrid(info Info) *qt6.QWidget {
	layout := qt6.NewQGridLayout(nil)
	layout.SetContentsMargins(0, 0, 0, 0)
	layout.SetSpacing(0)

	boldFont := qt6.NewQFont()
	boldFont.SetWeight(qt6.QFont__DemiBold)
	boldFont.SetPointSize(14)

	metrics := qt6.NewQFontMetrics(boldFont)
	maxWidth := 0

	for _, coreInfo := range info.CoreTypeInfos {
		w := metrics.HorizontalAdvance(coreInfo.Name) // + margins
		if w > maxWidth {
			maxWidth = w
		}
	}

	for index, coreInfo := range info.CoreTypeInfos {
		colNumber := index

		layout.SetColumnStretch(colNumber, 1)

		coreNameLabel := qt6.NewQLabel3(fmt.Sprint(coreInfo.Name))
		coreNameLabel.SetFont(boldFont)
		coreNameLabel.SetAlignment(qt6.AlignCenter)
		coreNameLabel.SetMinimumWidth(maxWidth)
		coreNameLabel.SetSizePolicy2(qt6.QSizePolicy__Expanding, qt6.QSizePolicy__Preferred)
		layout.AddWidget2(coreNameLabel.QWidget, 1, colNumber)

		coreCountLabel := qt6.NewQLabel3(fmt.Sprintf("%d Cores", coreInfo.CoreCount))
		coreCountLabel.SetSizePolicy2(qt6.QSizePolicy__Maximum, qt6.QSizePolicy__Maximum)
		layout.AddWidget4(coreCountLabel.QWidget, 2, colNumber, qt6.AlignCenter)

		coreThreadCountLabel := qt6.NewQLabel3(fmt.Sprintf("%d Threads", coreInfo.ThreadCount))
		coreThreadCountLabel.SetSizePolicy2(qt6.QSizePolicy__Maximum, qt6.QSizePolicy__Maximum)
		layout.AddWidget4(coreThreadCountLabel.QWidget, 3, colNumber, qt6.AlignCenter)
	}

	layout.SetRowStretch(0, 1)
	layout.SetRowStretch(5, 1)

	container := qt6.NewQWidget(nil)
	container.SetContentsMargins(0, 0, 0, 0)
	container.SetLayout(layout.QLayout)

	return container
}

func CreateCacheInfoGrid(info Info) *qt6.QWidget {
	layout := qt6.NewQGridLayout(nil)
	layout.SetContentsMargins(0, 0, 0, 0)
	layout.SetSpacing(0)

	boldFont := qt6.NewQFont()
	boldFont.SetWeight(qt6.QFont__DemiBold)
	boldFont.SetPointSize(14)

	for index, coreInfo := range info.CoreTypeInfos {
		colNumber := index

		rowNum := 2
		for i, cacheInfo := range coreInfo.CacheLevelInfos {
			cacheLevelLabel := qt6.NewQLabel3(fmt.Sprintf("%d %s %s", cacheInfo.Amount, cacheInfo.Unit, cacheInfo.Name))
			cacheLevelLabel.SetSizePolicy2(qt6.QSizePolicy__Maximum, qt6.QSizePolicy__Maximum)
			layout.AddWidget4(cacheLevelLabel.QWidget, rowNum+i, colNumber, qt6.AlignCenter)
		}
	}

	cacheIndicatorLabel := qt6.NewQLabel3("Cache per core")
	cacheIndicatorLabel.SetAlignment(qt6.AlignCenter)
	cacheIndicatorLabel.SetFont(boldFont)
	layout.AddWidget5(cacheIndicatorLabel.QWidget, 1, 0, 1, 2, qt6.AlignCenter)

	layout.SetRowStretch(0, 1)
	layout.SetRowStretch(5, 1)
	layout.SetHorizontalSpacing(5)
	layout.SetVerticalSpacing(2)

	container := qt6.NewQWidget(nil)
	container.SetContentsMargins(0, 0, 0, 0)
	container.SetLayout(layout.QLayout)

	return container
}
