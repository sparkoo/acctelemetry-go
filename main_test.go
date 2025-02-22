package acctelemetry_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/sparkoo/acctelemetry-go"
)

func TestTelemetry(t *testing.T) {
	telemetry := acctelemetry.New(acctelemetry.DefaultConfig())
	defer telemetry.Close()
	err := telemetry.Connect()
	if err != nil {
		t.Error(fmt.Errorf("unable to connect to ACC: %w\n", err))
	}

	fmt.Printf("Static: %+v\n\n", telemetry.StaticPointer())
	fmt.Printf("Graphics: %+v\n\n", telemetry.GraphicsPointer())
	fmt.Printf("Physics: %+v\n\n", telemetry.PhysicsPointer())

	time.Sleep(10 * time.Second)

	fmt.Printf("Static: %+v\n\n", telemetry.StaticPointer())
	fmt.Printf("Graphics: %+v\n\n", telemetry.GraphicsPointer())
	fmt.Printf("Physics: %+v\n\n", telemetry.PhysicsPointer())
}

func TestUdp(t *testing.T) {
	telemetry := acctelemetry.New(acctelemetry.DefaultUdpConfig())
	defer telemetry.Close()
	if err := telemetry.Connect(); err != nil {
		fmt.Printf("%+v\n", err)
	}
	t.Fail()
}

func convertToString(chars []uint16) string {
	var str string
	for _, val := range chars {
		str += string(rune(val))
	}
	return str
}
