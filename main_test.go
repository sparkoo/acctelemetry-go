package acctelemetry_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/sparkoo/acctelemetry-go"
	"github.com/sparkoo/acctelemetry-go/pkg/types"
)

func TestTelemetry(t *testing.T) {
	telemetry, err := acctelemetry.AccTelemetry()
	defer telemetry.Close()
	if err != nil {
		t.Error(fmt.Errorf("failed to create the telemetry: %w", err))
	}
	telemetry.SubscribePhysics(1*time.Second, func(ag *types.AccPhysics) {
		fmt.Printf("%+v\n\n", ag)
	})

	time.Sleep(10 * time.Second)
	fmt.Println("finishing test")
}
