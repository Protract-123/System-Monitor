package utils

import "log"

var units = [...]string{"KiB", "MiB", "GiB", "TiB", "PiB", "EiB", "ZiB", "YiB"}

func ConvertFromBytes[T Number](bytes T) ValueUnitPair[T] {
	value := bytes
	unit := "Bytes"

	i := 0
	for value >= 1024 {
		value = value / 1024
		unit = units[i]
		i += 1
	}

	return ValueUnitPair[T]{
		Value: value,
		Unit:  unit,
	}
}

func ConvertFromBytesParts[T Number](bytes T) (T, string) {
	pair := ConvertFromBytes(bytes)
	return pair.Value, pair.Unit
}

func ConvertFromBytesToUnit[T Number](bytes T, unit string) ValueUnitPair[T] {
	if unit == "Bytes" {
		return ValueUnitPair[T]{Value: bytes, Unit: "Bytes"}
	}

	value := bytes
	currUnit := "Bytes"

	i := 0
	for currUnit != unit {
		if i >= len(units) {
			log.Printf("utils: unknown unit %q, returning raw bytes", unit)
			return ValueUnitPair[T]{Value: bytes, Unit: "Bytes"}
		}
		value = value / 1024
		currUnit = units[i]
		i += 1
	}

	return ValueUnitPair[T]{
		Value: value,
		Unit:  unit,
	}
}

func ConvertFromBytesToUnitParts[T Number](bytes T, unit string) (T, string) {
	pair := ConvertFromBytesToUnit(bytes, unit)
	return pair.Value, pair.Unit
}
