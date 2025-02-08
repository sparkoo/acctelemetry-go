package acctelemetry_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/sparkoo/acctelemetry-go"
)

func TestTelemetry(t *testing.T) {
	telemetry := acctelemetry.New()
	err := telemetry.Connect()
	if err != nil {
		t.Error(fmt.Errorf("unable to connect to ACC: %w", err))
	}

	fmt.Printf("Static: %+v\n\n", telemetry.StaticPointer())
	fmt.Printf("Graphics: %+v\n\n", telemetry.GraphicsPointer())
	fmt.Printf("Physics: %+v\n\n", telemetry.PhysicsPointer())

	time.Sleep(10 * time.Second)

	fmt.Printf("Static: %+v\n\n", telemetry.StaticPointer())
	fmt.Printf("Graphics: %+v\n\n", telemetry.GraphicsPointer())
	fmt.Printf("Physics: %+v\n\n", telemetry.PhysicsPointer())
}

func convertToString(chars []uint16) string {
	var str string
	for _, val := range chars {
		str += string(rune(val))
	}
	return str
}
