package cpu

import (
	"embed"
	"regexp"
	"strings"
)

type info struct {
	Model         string     `yaml:"cpu_model"`
	Cores         uint       `yaml:"cpu_cores"`
	Threads       uint       `yaml:"cpu_threads"`
	Codename      string     `yaml:"cpu_codename"`
	CoreTypeInfos []coreInfo `yaml:"cpu_core_infos"`
	IconPath      string     `yaml:"icon_path"`
}

type coreInfo struct {
	Name            string      `yaml:"core_name"`
	CoreCount       uint        `yaml:"core_count"`
	ThreadCount     uint        `yaml:"core_thread_count"`
	CacheLevelInfos []cacheInfo `yaml:"core_cache_infos"`
}

type cacheInfo struct {
	Name   string `yaml:"cache_name"`
	Amount uint   `yaml:"cache_amount"`
	Unit   string `yaml:"cache_unit"`
}

var cpuToCodename = map[string]string{
	// Retrieved from https://en.wikipedia.org/wiki/List_of_Apple_codenames#M_series
	// Named after islands, pretty cool naming scheme tbh
	"apple m1": "Tonga",
	"apple m2": "Staten",
	"apple m3": "Ibiza",
	"apple m4": "Donan",
	"apple m5": "Hidra",
}

var cpuToFilePath = map[string]string{
	"apple m1": "icons/AppleM1.png",
	"apple m2": "icons/AppleM2.png",
	"apple m3": "icons/AppleM3.png",
	"apple m4": "icons/AppleM4.png",
	"apple m5": "icons/AppleM5.png",
}

var baseModelRegex = regexp.MustCompile(`(?i)apple m\d+`)

func lookupBaseModel(model string, m map[string]string) string {
	base := baseModelRegex.FindString(model)
	return m[strings.ToLower(base)]
}

//go:embed icons
var icons embed.FS
