package util

import (
	"hash/crc32"
	"math"
)

func Hash(k interface{}) uint32 {
	switch k.(type) {
	case int:
		return uint32(k.(int))
	case int8:
		return uint32(k.(int8))
	case int16:
		return uint32(k.(int16))
	case int32:
		return uint32(k.(int32))
	case int64:
		return uint32(k.(int64))
	case uint:
		return uint32(k.(uint))
	case uint8:
		return uint32(k.(uint8))
	case uint16:
		return uint32(k.(uint16))
	case uint32:
		return k.(uint32)
	case uint64:
		return uint32(k.(uint64))
	case uintptr:
		return uint32(k.(uintptr))
	case float32:
		return math.Float32bits(k.(float32))
	case float64:
		return uint32(math.Float64bits(k.(float64)))
	case string:
		return crc32.ChecksumIEEE([]byte(k.(string)))
	default:
		panic(k)
	}
}
