package cpu

type Info struct {
	Model         string     `yaml:"cpu_model"`
	Cores         uint       `yaml:"cpu_cores"`
	Threads       uint       `yaml:"cpu_threads"`
	Codename      string     `yaml:"cpu_codename"`
	CoreTypeInfos []CoreInfo `yaml:"cpu_core_infos"`
}

type CoreInfo struct {
	Name            string      `yaml:"core_name"`
	CoreCount       uint        `yaml:"core_count"`
	ThreadCount     uint        `yaml:"core_thread_count"`
	CacheLevelInfos []CacheInfo `yaml:"core_cache_infos"`
}

type CacheInfo struct {
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
