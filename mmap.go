package acctelemetry

import (
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

type mmap struct {
	kernel32 *syscall.LazyDLL
	hMap     uintptr
	addr     uintptr
}

func mapFile(fileName string, sharedMemorySize uintptr) (*mmap, error) {
	mmap := &mmap{kernel32: syscall.NewLazyDLL("kernel32.dll")}

	// Open shared memory with read access
	openFileMapping := mmap.kernel32.NewProc("OpenFileMappingW")
	filePointer, err := syscall.UTF16PtrFromString(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to get mmap file pointer '%s': %w", fileName, err)
	}
	hMap, _, _ := openFileMapping.Call(
		windows.FILE_MAP_READ,
		0, uintptr(unsafe.Pointer(filePointer)),
	)

	if hMap == 0 {
		return nil, fmt.Errorf("Failed to open shared memory file '%s'", fileName)
	}
	mmap.hMap = hMap

	// Map view of the file
	mapViewOfFile := mmap.kernel32.NewProc("MapViewOfFile")
	addr, _, _ := mapViewOfFile.Call(
		hMap,
		windows.FILE_MAP_READ,
		0, 0, sharedMemorySize,
	)

	if addr == 0 {
		return nil, fmt.Errorf("Failed to map shared memory.")
	}
	mmap.addr = addr

	return mmap, nil
}

func (m *mmap) pointer() unsafe.Pointer {
	return unsafe.Pointer(m.addr)
}

func (m *mmap) Close() error {
	unmapViewOfFile := m.kernel32.NewProc("UnmapViewOfFile")
	syscall.CloseHandle(syscall.Handle(m.hMap))
	unmapViewOfFile.Call(m.addr)
	return nil
}
