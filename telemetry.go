package acctelemetry

import (
	"fmt"
	"unsafe"
)

const STATIC_FILE_MMAP = "Local\\acpmf_static"
const PHYSICS_FILE_MMAP = "Local\\acpmf_physics"
const GRAPHIS_FILE_MMAP = "Local\\acpmf_graphics"

type AccTelemetry struct {
	staticData   *accDataHolder[AccStatic]
	physicsData  *accDataHolder[AccPhysics]
	graphicsData *accDataHolder[AccGraphic]
}

type accDataHolder[T AccGraphic | AccPhysics | AccStatic] struct {
	mmap *mmap
	data *T
}

func (d *accDataHolder[T]) Close() error {
	if d.mmap != nil {
		d.mmap.Close()
		d.mmap = nil
	}
	d.data = nil
	return nil
}

func (t *AccTelemetry) Connect() error {
	var accStatic AccStatic
	staticMMap, err := mapFile(STATIC_FILE_MMAP, unsafe.Sizeof(accStatic))
	if err != nil {
		return fmt.Errorf("Failed to create mapping to ACC static file: %w", err)
	}
	t.staticData = &accDataHolder[AccStatic]{
		mmap: staticMMap,
		data: (*AccStatic)(staticMMap.pointer()),
	}

	var accPhysics AccPhysics
	physicsMMap, err := mapFile(PHYSICS_FILE_MMAP, unsafe.Sizeof(accPhysics))
	if err != nil {
		return fmt.Errorf("Failed to create mapping to ACC physics file: %w", err)
	}
	t.physicsData = &accDataHolder[AccPhysics]{
		mmap: physicsMMap,
		data: (*AccPhysics)(physicsMMap.pointer()),
	}

	var accGraphic AccGraphic
	graphicsMMap, err := mapFile(GRAPHIS_FILE_MMAP, unsafe.Sizeof(accGraphic))
	if err != nil {
		return fmt.Errorf("Failed to create mapping to ACC physics file: %w", err)
	}
	t.graphicsData = &accDataHolder[AccGraphic]{
		mmap: graphicsMMap,
		data: (*AccGraphic)(graphicsMMap.pointer()),
	}

	return nil
}

func New() *AccTelemetry {
	return &AccTelemetry{}
}

// this returns direct pointer to the memory so underlying struct will change over time
func (t *AccTelemetry) GraphicsPointer() *AccGraphic {
	if t.graphicsData != nil {
		return t.graphicsData.data
	}
	return nil
}

// this returns direct pointer to the memory so underlying struct will change over time
func (t *AccTelemetry) StaticPointer() *AccStatic {
	if t.staticData != nil {
		return t.staticData.data
	}
	return nil
}

// this returns direct pointer to the memory so underlying struct will change over time
func (t *AccTelemetry) PhysicsPointer() *AccPhysics {
	if t.physicsData != nil {
		return t.physicsData.data
	}
	return nil
}

func (t *AccTelemetry) Close() error {
	t.graphicsData.Close()
	t.staticData.Close()
	t.physicsData.Close()
	return nil
}
