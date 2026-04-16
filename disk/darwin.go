package disk

import (
	"System_Monitor/utils"
	"strings"
	"syscall"

	"github.com/mappu/miqt/qt6"
)

func GenerateUI() *qt6.QLayout {
	diskInfo := fetchInfo()

	diskLayout := qt6.NewQHBoxLayout2()
	diskLayout.SetContentsMargins(0, 0, 0, 0)

	diskInfoContainer := createDiskInfoContainer(diskInfo)

	diskLayout.AddWidget(diskInfoContainer)
	return diskLayout.QLayout
}

func fetchInfo() info {
	var stat syscall.Statfs_t
	err := syscall.Statfs("/", &stat)
	if err != nil {
		return info{}
	}

	totalBytes := stat.Blocks * uint64(stat.Bsize)
	freeBytes := stat.Bavail * uint64(stat.Bsize)
	usedBytes := totalBytes - freeBytes

	return info{
		DiskMountPoint:   "/",
		FilesystemFormat: strings.TrimRight(int8ToStr(stat.Fstypename[:]), "\x00"),
		Size:             utils.ConvertFromBytes(float32(totalBytes)),
		AvailableSpace:   utils.ConvertFromBytes(float32(freeBytes)),
		UsedSpace:        utils.ConvertFromBytes(float32(usedBytes)),
		UsedPercent:      int((100 * float32(usedBytes)) / float32(totalBytes)),
	}
}

func int8ToStr(arr []int8) string {
	b := make([]byte, len(arr))
	for i, v := range arr {
		if v == 0 {
			return string(b[:i])
		}
		b[i] = byte(v)
	}
	return string(b)
}
