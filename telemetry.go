package acctelemetry

import (
	"fmt"
	"time"
	"unsafe"

	"github.com/sparkoo/acctelemetry-go/pkg/mmap"
	"github.com/sparkoo/acctelemetry-go/pkg/types"
)

const PHYSICS_FILE_MMAP = "Local\\acpmf_physics"
const STATIC_FILE_MMAP = "Local\\acpmf_static"
const GRAPHIS_FILE_MMAP = "Local\\acpmf_graphics"

type accTelemetry struct {
	graphicsData *accDataHolder[types.AccGraphic]
	staticData   *accDataHolder[types.AccStatic]
	physicsData  *accDataHolder[types.AccPhysics]
}

type accDataHolder[T types.AccGraphic | types.AccPhysics | types.AccStatic] struct {
	mmap           *mmap.MMap
	data           *T
	callbackTicker *time.Ticker
}

func (d *accDataHolder[T]) Close() error {
	if d.callbackTicker != nil {
		d.callbackTicker.Stop()
	}
	if d.mmap != nil {
		d.mmap.Close()
	}
	d.data = nil
	return nil
}

func (d *accDataHolder[T]) subscribe(callback func(*T), pollRate time.Duration) {
	d.callbackTicker = time.NewTicker(pollRate)
	go func() {
		for range d.callbackTicker.C {
			callback(d.data)
		}
	}()
}

func AccTelemetry() (*accTelemetry, error) {
	var accGraphic types.AccGraphic
	graphicsMMap, err := mmap.MapFile(GRAPHIS_FILE_MMAP, unsafe.Sizeof(accGraphic))
	if err != nil {
		return nil, fmt.Errorf("Failed to create mapping to ACC physics file: %w", err)
	}

	var accStatic types.AccStatic
	staticMMap, err := mmap.MapFile(STATIC_FILE_MMAP, unsafe.Sizeof(accStatic))
	if err != nil {
		return nil, fmt.Errorf("Failed to create mapping to ACC static file: %w", err)
	}

	var AccPhysics types.AccPhysics
	physicsMMap, err := mmap.MapFile(PHYSICS_FILE_MMAP, unsafe.Sizeof(AccPhysics))
	if err != nil {
		return nil, fmt.Errorf("Failed to create mapping to ACC physics file: %w", err)
	}

	return &accTelemetry{
		graphicsData: &accDataHolder[types.AccGraphic]{
			mmap: graphicsMMap,
			data: (*types.AccGraphic)(graphicsMMap.Pointer()),
		},

		staticData: &accDataHolder[types.AccStatic]{
			mmap: staticMMap,
			data: (*types.AccStatic)(staticMMap.Pointer()),
		},

		physicsData: &accDataHolder[types.AccPhysics]{
			mmap: physicsMMap,
			data: (*types.AccPhysics)(physicsMMap.Pointer()),
		},
	}, nil
}

func (t *accTelemetry) SubscribeGraphic(pollRate time.Duration, callback func(*types.AccGraphic)) {
	t.graphicsData.subscribe(callback, pollRate)
}

func (t *accTelemetry) SubscribePhysics(pollRate time.Duration, callback func(*types.AccPhysics)) {
	t.physicsData.subscribe(callback, pollRate)
}

func (t *accTelemetry) SubscribeStatic(pollRate time.Duration, callback func(*types.AccStatic)) {
	t.staticData.subscribe(callback, pollRate)
}

func (t *accTelemetry) Close() error {
	t.graphicsData.Close()
	t.staticData.Close()
	t.physicsData.Close()
	return nil
}
