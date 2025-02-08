package acctelemetry

import (
	"fmt"
	"unsafe"
)

const STATIC_FILE_MMAP = "Local\\acpmf_static"
const PHYSICS_FILE_MMAP = "Local\\acpmf_physics"
const GRAPHIS_FILE_MMAP = "Local\\acpmf_graphics"

type accTelemetry struct {
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

func (t *accTelemetry) Connect() error {
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

func AccTelemetry() *accTelemetry {
	return &accTelemetry{}
}

// this returns direct pointer to the memory so underlying struct will change over time
func (t *accTelemetry) GraphicsPointer() *AccGraphic {
	return t.graphicsData.data
}

// this returns direct pointer to the memory so underlying struct will change over time
func (t *accTelemetry) StaticPointer() *AccStatic {
	return t.staticData.data
}

// this returns direct pointer to the memory so underlying struct will change over time
func (t *accTelemetry) PhysicsPointer() *AccPhysics {
	return t.physicsData.data
}

func (t *accTelemetry) Close() error {
	t.graphicsData.Close()
	t.staticData.Close()
	t.physicsData.Close()
	return nil
}
