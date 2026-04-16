package disk

import "System_Monitor/utils"

type info struct {
	DiskMountPoint   string
	FilesystemFormat string
	Size             utils.ValueUnitPair[float32]
	AvailableSpace   utils.ValueUnitPair[float32]
	UsedSpace        utils.ValueUnitPair[float32]
	UsedPercent      int
}
