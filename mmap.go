package acctelemetry

import (
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

type mmap struct {
	hMap uintptr
	addr uintptr
}

type win struct {
	openFileMapping *windows.LazyProc
	mapViewOfFile   *windows.LazyProc
	unmapViewOfFile *windows.LazyProc
}

var w = initWindows()

func initWindows() *win {
	var kernel32 = windows.NewLazyDLL("kernel32.dll")
	return &win{
		openFileMapping: kernel32.NewProc("OpenFileMappingW"),
		mapViewOfFile:   kernel32.NewProc("MapViewOfFile"),
		unmapViewOfFile: kernel32.NewProc("UnmapViewOfFile"),
	}
}

func mapFile(fileName string, sharedMemorySize uintptr) (*mmap, error) {
	mmap := &mmap{}

	// Open shared memory with read access
	filePointer, err := syscall.UTF16PtrFromString(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to get mmap file pointer '%s': %w", fileName, err)
	}
	hMap, _, _ := w.openFileMapping.Call(
		windows.FILE_MAP_READ,
		0, uintptr(unsafe.Pointer(filePointer)),
	)

	if hMap == 0 {
		return nil, fmt.Errorf("Failed to open shared memory file '%s'", fileName)
	}
	mmap.hMap = hMap

	// Map view of the file
	addr, _, _ := w.mapViewOfFile.Call(
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
	syscall.CloseHandle(syscall.Handle(m.hMap))
	w.unmapViewOfFile.Call(m.addr)
	return nil
}
