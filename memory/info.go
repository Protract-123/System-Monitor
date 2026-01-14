package memory

import (
	"System_Monitor/utils"
)

type Info struct {
	TotalMemory  utils.ValueUnitPair[uint]    `yaml:"total_memory"`
	UsableMemory utils.ValueUnitPair[float32] `yaml:"usable_memory"`

	FreeMemory utils.ValueUnitPair[float32] `yaml:"free_memory"`
	UsedMemory utils.ValueUnitPair[float32] `yaml:"used_memory"`

	SwapUsed  utils.ValueUnitPair[float32] `yaml:"swap_used"`
	SwapFree  utils.ValueUnitPair[float32] `yaml:"swap_free"`
	SwapTotal utils.ValueUnitPair[uint]    `yaml:"swap_total"`

	PlatformInfo PlatformInfo `yaml:"platform_info,omitempty"`
}

type PlatformInfo interface {
	isPlatformInfo()
}
