package utils

type Number interface {
	~int | ~uint | ~float32 | ~float64
}

type ValueUnitPair[T Number] struct {
	Value T      `yaml:"value"`
	Unit  string `yaml:"unit"`
}
