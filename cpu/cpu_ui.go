package cpu

import (
	"System_Monitor/ui"
	"fmt"

	"github.com/mappu/miqt/qt6"
	"go.yaml.in/yaml/v4"
)

func GenerateUI(parent *qt6.QWidget) {
	info := FetchInfoOSX()

	// Print cpu info as yaml
	yamlData, _ := yaml.Marshal(info)
	fmt.Println(string(yamlData))

	// hLayout is the root layout for CPU UI
	hLayout := qt6.NewQHBoxLayout(nil)
	hLayout.SetContentsMargins(10, 10, 10, 10)
	hLayout.SetSpacing(10)

	// Create CPU Image
	imageContainer := CreateCPUImage()
	cpuInfoContainer := CreateCPUInfoContainer(info)
	coreTypeContainer := CreateCoreTypeContainer(info)

	hLayout.AddStretchWithStretch(1)
	hLayout.AddWidget(imageContainer)
	hLayout.AddWidget(cpuInfoContainer)
	hLayout.AddWidget(coreTypeContainer)
	hLayout.AddStretchWithStretch(1)

	parent.SetLayout(hLayout.QLayout)
}

func CreateCPUImage() *qt6.QWidget {
	imageObject := qt6.NewQLabel(nil)
	pixmap := qt6.NewQPixmap7("./cpu/icons/AppleChipM1.png", "", qt6.AutoColor)
	imageObject.SetPixmap(pixmap)

	imageObject.SetAlignment(qt6.AlignCenter)
	imageObject.SetScaledContents(true)
	imageObject.SetMaximumSize2(256, 256)

	imageLayout := qt6.NewQVBoxLayout(nil)
	imageLayout.SetSpacing(0)
	imageLayout.AddWidget(imageObject.QWidget)

	imageContainer := ui.NewBorderContainer(
		nil,
		imageLayout.QLayout,
		2,
		0,
		qt6.NewQColor11(0, 180, 255, 200),
	)
	imageContainer.SetSizePolicy2(qt6.QSizePolicy__Maximum, qt6.QSizePolicy__Maximum)

	return imageContainer.QWidget
}

func CreateCPUInfoContainer(info Info) *qt6.QWidget {
	layout := qt6.NewQVBoxLayout(nil)

	layout.SetContentsMargins(0, 0, 0, 0)
	layout.SetSpacing(2)

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

	coreCountLabel := qt6.NewQLabel5(fmt.Sprintf("Threads/Cores: %d/%d", info.Threads, info.Cores), nil)
	coreCountLabel.SetSizePolicy2(qt6.QSizePolicy__Maximum, qt6.QSizePolicy__Maximum)

	layout.AddStretchWithStretch(1)
	layout.AddWidget3(title.QWidget, 0, qt6.AlignHCenter)
	layout.AddWidget3(modelLabel.QWidget, 0, qt6.AlignHCenter)
	layout.AddWidget3(coreCountLabel.QWidget, 0, qt6.AlignHCenter)
	layout.AddStretchWithStretch(1)

	container := ui.NewBorderContainer(
		nil,
		layout.QLayout,
		2,
		5,
		qt6.NewQColor11(255, 255, 255, 200),
	)
	container.SetSizePolicy2(qt6.QSizePolicy__Maximum, qt6.QSizePolicy__Maximum)

	return container.QWidget
}

func CreateCoreTypeContainer(info Info) *qt6.QWidget {
	layout := qt6.NewQGridLayout(nil)
	layout.SetContentsMargins(5, 5, 5, 5)
	layout.SetSpacing(8)

	boldFont := qt6.NewQFont()
	boldFont.SetWeight(qt6.QFont__DemiBold)
	boldFont.SetPointSize(14)

	for index, coreInfo := range info.CoreTypeInfos {
		colNumber := index

		coreNameLabel := qt6.NewQLabel3(fmt.Sprint(coreInfo.Name))
		coreNameLabel.SetFont(boldFont)
		coreNameLabel.SetSizePolicy2(qt6.QSizePolicy__Maximum, qt6.QSizePolicy__Maximum)
		layout.AddWidget4(coreNameLabel.QWidget, 1, colNumber, qt6.AlignCenter)

		coreCountLabel := qt6.NewQLabel3(fmt.Sprintf("%d Cores", coreInfo.CoreCount))
		coreCountLabel.SetSizePolicy2(qt6.QSizePolicy__Maximum, qt6.QSizePolicy__Maximum)
		layout.AddWidget4(coreCountLabel.QWidget, 2, colNumber, qt6.AlignCenter)
		fmt.Println(coreCountLabel.Font().PointSize())

		coreThreadCountLabel := qt6.NewQLabel3(fmt.Sprintf("%d Threads", coreInfo.ThreadCount))
		coreThreadCountLabel.SetSizePolicy2(qt6.QSizePolicy__Maximum, qt6.QSizePolicy__Maximum)
		layout.AddWidget4(coreThreadCountLabel.QWidget, 3, colNumber, qt6.AlignCenter)

		rowNum := 5
		for i, cacheInfo := range coreInfo.CacheLevelInfos {
			cacheLevelLabel := qt6.NewQLabel3(fmt.Sprintf("%d %s %s", cacheInfo.Amount, cacheInfo.Unit, cacheInfo.Name))
			cacheLevelLabel.SetSizePolicy2(qt6.QSizePolicy__Maximum, qt6.QSizePolicy__Maximum)
			layout.AddWidget4(cacheLevelLabel.QWidget, rowNum+i, colNumber, qt6.AlignCenter)
		}
	}

	cacheIndicatorLabel := qt6.NewQLabel3("Cache Breakdown\n(per core)")
	cacheIndicatorLabel.SetAlignment(qt6.AlignCenter)
	cacheIndicatorLabel.SetContentsMargins(0, 5, 0, 5)
	layout.AddWidget5(cacheIndicatorLabel.QWidget, 4, 0, 1, 2, qt6.AlignCenter)

	layout.SetRowStretch(0, 1)
	layout.SetRowStretch(8, 1)

	container := qt6.NewQWidget(nil)
	container.SetContentsMargins(10, 0, 10, 0)
	container.SetLayout(layout.QLayout)

	return container
}
