package as

// #include <stdlib.h>
import "C"

import (
	"unsafe"
)

// bytesToOffsetSize returns a pointer and size pair for the given byte slice in a
// format compatible with WebAssembly numeric types.
func bytesToOffsetSize(data []byte) (uint32, uint32) {
	ptr := unsafe.Pointer(unsafe.SliceData(data))
	offset := uint32(uintptr(ptr))
	size := uint32(len(data))
	return offset, size
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
