package disk

import (
	"System_Monitor/utils"
	"log"
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
		log.Printf("disk: statfs(\"/\") failed: %v", err)
		return info{
			DiskMountPoint:   "/",
			FilesystemFormat: "Unknown",
		}
	}

	totalBytes := stat.Blocks * uint64(stat.Bsize)
	freeBytes := stat.Bavail * uint64(stat.Bsize)

	var usedBytes uint64
	if freeBytes < totalBytes {
		usedBytes = totalBytes - freeBytes
	}

	usedPercent := 0
	if totalBytes > 0 {
		usedPercent = int((100 * float32(usedBytes)) / float32(totalBytes))
	}

	return info{
		DiskMountPoint:   "/",
		FilesystemFormat: strings.TrimRight(int8ToStr(stat.Fstypename[:]), "\x00"),
		Size:             utils.ConvertFromBytes(float32(totalBytes)),
		AvailableSpace:   utils.ConvertFromBytes(float32(freeBytes)),
		UsedSpace:        utils.ConvertFromBytes(float32(usedBytes)),
		UsedPercent:      usedPercent,
	}
}

// int8ToStr converts a C-style int8 byte array to a Go string, stopping at the
// first null terminator. If no null is present, the whole array is returned.
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
