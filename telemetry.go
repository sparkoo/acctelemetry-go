package acctelemetry

import (
	"fmt"
	"unsafe"
)

const PHYSICS_FILE_MMAP = "Local\\acpmf_physics"
const STATIC_FILE_MMAP = "Local\\acpmf_static"
const GRAPHIS_FILE_MMAP = "Local\\acpmf_graphics"

type accTelemetry struct {
	graphicsData *accDataHolder[AccGraphic]
	staticData   *accDataHolder[AccStatic]
	physicsData  *accDataHolder[AccPhysics]
}

type accDataHolder[T AccGraphic | AccPhysics | AccStatic] struct {
	mmap *mmap
	data *T
}

func (d *accDataHolder[T]) Close() error {
	if d.mmap != nil {
		d.Close()
	}
	d.data = nil
	return nil
}

func AccTelemetry() (*accTelemetry, error) {
	var accGraphic AccGraphic
	graphicsMMap, err := mapFile(GRAPHIS_FILE_MMAP, unsafe.Sizeof(accGraphic))
	if err != nil {
		return nil, fmt.Errorf("Failed to create mapping to ACC physics file: %w", err)
	}

	var accStatic AccStatic
	staticMMap, err := mapFile(STATIC_FILE_MMAP, unsafe.Sizeof(accStatic))
	if err != nil {
		return nil, fmt.Errorf("Failed to create mapping to ACC static file: %w", err)
	}

	var accPhysics AccPhysics
	physicsMMap, err := mapFile(PHYSICS_FILE_MMAP, unsafe.Sizeof(accPhysics))
	if err != nil {
		return nil, fmt.Errorf("Failed to create mapping to ACC physics file: %w", err)
	}

	return &accTelemetry{
		graphicsData: &accDataHolder[AccGraphic]{
			mmap: graphicsMMap,
			data: (*AccGraphic)(graphicsMMap.pointer()),
		},

		staticData: &accDataHolder[AccStatic]{
			mmap: staticMMap,
			data: (*AccStatic)(staticMMap.pointer()),
		},

		physicsData: &accDataHolder[AccPhysics]{
			mmap: physicsMMap,
			data: (*AccPhysics)(physicsMMap.pointer()),
		},
	}, nil
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
