package acctelemetry

import (
	"fmt"
	"net"
	"unsafe"
)

const STATIC_FILE_MMAP = "Local\\acpmf_static"
const PHYSICS_FILE_MMAP = "Local\\acpmf_physics"
const GRAPHIS_FILE_MMAP = "Local\\acpmf_graphics"

type AccTelemetryConfig struct {
	UdpIpPort                   string
	UdpConnectionPassword       string
	UdpCommandPassword          string
	UdpDisplayName              string
	UdpRealtimeUpdateIntervalMS int32
}

func DefaultConfig() *AccTelemetryConfig {
	return &AccTelemetryConfig{
		UdpIpPort:                   "127.0.0.1:9000",
		UdpConnectionPassword:       "asd",
		UdpCommandPassword:          "",
		UdpDisplayName:              "RaceMate",
		UdpRealtimeUpdateIntervalMS: 100,
	}
}

type accTelemetry struct {
	config *AccTelemetryConfig

	staticData   *accDataHolder[AccStatic]
	physicsData  *accDataHolder[AccPhysics]
	graphicsData *accDataHolder[AccGraphic]

	udpConnection *net.UDPConn
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
	udpErr := t.connectUdp()
	if udpErr != nil {
		return fmt.Errorf("failed UDP connection: %w", udpErr)
	}

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

func (telemetry *accTelemetry) connectUdp() error {
	udpAddress, err := net.ResolveUDPAddr("udp", telemetry.config.UdpIpPort)
	if err != nil {
		return fmt.Errorf("failed to resolve UDP address '%s': %w", telemetry.config.UdpIpPort, err)
	}

	udpConnection, err := net.DialUDP("udp", nil, udpAddress)
	if err != nil {
		return fmt.Errorf("failed to dial UDP: %w", err)
	}
	telemetry.udpConnection = udpConnection

	if err := telemetry.connect(); err != nil {
		return fmt.Errorf("failed to connect to UDP: %w", err)
	}
	return nil
}

func New(config *AccTelemetryConfig) *accTelemetry {
	return &accTelemetry{
		config: config,
	}
}

// this returns direct pointer to the memory so underlying struct will change over time
func (t *accTelemetry) GraphicsPointer() *AccGraphic {
	if t.graphicsData != nil {
		return t.graphicsData.data
	}
	return nil
}

// this returns direct pointer to the memory so underlying struct will change over time
func (t *accTelemetry) StaticPointer() *AccStatic {
	if t.staticData != nil {
		return t.staticData.data
	}
	return nil
}

// this returns direct pointer to the memory so underlying struct will change over time
func (t *accTelemetry) PhysicsPointer() *AccPhysics {
	if t.physicsData != nil {
		return t.physicsData.data
	}
	return nil
}

// reads from UDP
func (t *accTelemetry) RealtimeUpdate() *RealtimeUpdate {

}

func (t *accTelemetry) Close() error {
	if t.graphicsData != nil {
		t.graphicsData.Close()
	}

	if t.staticData != nil {
		t.staticData.Close()
	}

	if t.physicsData != nil {
		t.physicsData.Close()
	}

	if t.udpConnection != nil {
		t.udpConnection.Close()
	}

	return nil
}
