package as

// the flags for cgo below are a workaround for fixing some includes that
// seem to not be sent properly with nix. Hash for the nix store acquired
// by running `which tinygo`

// #cgo CFLAGS: -I /nix/store/73zik9gpf3n9vr699nyn77abk21v32kk-tinygo-0.31.2/share/tinygo/lib/wasi-libc/sysroot/include/
// #include <stdlib.h>
import "C"

import (
	"unsafe"
)

// bytesToPtr returns a pointer and size pair for the given byte slice in a
// format compatible with WebAssembly numeric types.
func bytesToPtr(data []byte) (uint32, uint32) {
	ptr := unsafe.Pointer(unsafe.SliceData(data))
	offset := uint32(uintptr(ptr))
	length := uint32(len(data))
	return offset, length
}

// bytesFromPtr takes a packed offset and size pointer in a format compatible
// with WebAssembly numeric types and returns the corresponding byte slice.
func bytesFromPtr(inPtr uint64) []byte {
	offset := uint32(inPtr >> 32)
	size := uint32(inPtr)
	return unsafe.Slice((*byte)(unsafe.Pointer(uintptr(offset))), size)
}

func Bytes(inPtr uint64) []byte {
	return bytesFromPtr(inPtr)
}

// leakToSharedMem is used to persist data in shared memory beyond the execution
// of this program. We expect the data will be cleaned up by the host.
func leakToSharedMem(v []byte) uint64 {
	size := C.ulong(len(v))
	ptr := unsafe.Pointer(C.malloc(size))
	copy(unsafe.Slice((*byte)(ptr), size), v)
	return (uint64(uintptr(ptr)) << uint64(32)) | uint64(size)
}

func ShareWithHost(v []byte) uint64 {
	return leakToSharedMem(v)
}
