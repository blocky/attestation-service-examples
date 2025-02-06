package as

import (
	"fmt"
	"runtime"
)

// Imported host functions and supporting input/output types

//go:wasmimport env bufferLog
func _hostFuncBufferLog(ptr, size uint32)

//go:wasmimport env consoleLog
func _hostFuncConsoleLog(ptr, size uint32)

func LogToHost(msg string) {
	msgData := []byte(msg)
	inOffset, inLen := bytesToPtr(msgData)
	_hostFuncConsoleLog(inOffset, inLen)
	runtime.KeepAlive(msgData)
}

func Logf(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	msgData := []byte(msg)
	inOffset, inLen := bytesToPtr(msgData)
	_hostFuncConsoleLog(inOffset, inLen)
	runtime.KeepAlive(msgData)
}
