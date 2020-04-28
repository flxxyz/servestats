package utils

import (
	"fmt"
	"strconv"
	"unsafe"
)

func Int2Str(value int) string {
	return strconv.Itoa(value)
}

func Int64toStr(value int64) string {
	return strconv.FormatInt(value, 10)
}

func Uint64toStr(value uint64) string {
	return strconv.FormatUint(value, 10)
}

func Str2Int(value string) (val int) {
	val, _ = strconv.Atoi(value)
	return
}

func Bytes2Str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func Decimal(value float64, bit interface{}) (val float64) {
	format := "%f"
	switch bit.(type) {
	case int:
		format = "%." + Int2Str(bit.(int)) + "f"
	case string:
		format = "%." + bit.(string) + "f"
	}
	val, _ = strconv.ParseFloat(fmt.Sprintf(format, value), 64)
	return
}
