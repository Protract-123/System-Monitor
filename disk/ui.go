package disk

import (
	"System_Monitor/ui"
	"fmt"

	"github.com/mappu/miqt/qt6"
)

func createDiskInfoContainer(diskInfo info) *qt6.QWidget {
	rootLayout := qt6.NewQVBoxLayout2()
	rootLayout.SetContentsMargins(0, 0, 0, 0)
	rootLayout.SetSpacing(6)

	title := qt6.NewQLabel5("Disk Info", nil)
	title.SetFont(ui.HeadingFont)
	title.SetAlignment(qt6.AlignCenter)
	title.SetContentsMargins(0, 0, 0, 0)

	rootLayout.AddWidget(title.QWidget)

	addRow := func(labelText, valueText string) *qt6.QLabel {
		row := qt6.NewQHBoxLayout2()
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

	addRow("Disk Mount Point", diskInfo.DiskMountPoint)
	addRow("Filesystem Format", diskInfo.FilesystemFormat)
	addDivider()

	addRow("Total Disk Space", fmt.Sprintf("%.2f%s", diskInfo.Size.Value, diskInfo.Size.Unit))
	addRow("•  Used Disk Space", fmt.Sprintf("%.2f%s", diskInfo.UsedSpace.Value, diskInfo.UsedSpace.Unit))
	addRow("•  Available Disk Space", fmt.Sprintf("%.2f%s", diskInfo.AvailableSpace.Value, diskInfo.AvailableSpace.Unit))

	rootLayout.AddLayout(createDiskUsagePill(diskInfo))

	container := qt6.NewQWidget(nil)
	container.SetLayout(rootLayout.QLayout)
	container.SetSizePolicy2(qt6.QSizePolicy__Preferred, qt6.QSizePolicy__Preferred)

	return container
}

func createDiskUsagePill(diskInfo info) *qt6.QLayout {
	row := qt6.NewQHBoxLayout2()
	row.SetContentsMargins(0, 0, 0, 0)
	row.SetSpacing(5)

	label := qt6.NewQLabel5(fmt.Sprintf("%d%% used", diskInfo.UsedPercent), nil)
	label.SetFont(ui.BoldFont)
	label.SetAlignment(qt6.AlignLeft | qt6.AlignVCenter)
	label.SetSizePolicy2(qt6.QSizePolicy__Preferred, qt6.QSizePolicy__Preferred)

	pill := qt6.NewQProgressBar(nil)
	pill.SetMinimum(0)
	pill.SetMaximum(100)
	pill.SetValue(diskInfo.UsedPercent)
	pill.SetTextVisible(true)
	pill.SetFormat(fmt.Sprintf("%d%%", diskInfo.UsedPercent))
	pill.SetAlignment(qt6.AlignRight | qt6.AlignVCenter)
	pill.SetFixedHeight(16)
	pill.SetSizePolicy2(qt6.QSizePolicy__Expanding, qt6.QSizePolicy__Fixed)

	row.AddWidget(label.QWidget)
	row.AddWidget(pill.QWidget)

	return row.QLayout
}
