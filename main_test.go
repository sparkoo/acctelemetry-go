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
	// telemetry.SubscribeStatic(1*time.Second, func(ag *types.AccStatic) {
	// 	fmt.Printf("%+v\n\n", ag)
	// })

	telemetry.SubscribeGraphic(1*time.Second, func(ag *types.AccGraphic) {
		fmt.Printf("%+v\n\n", ag)
	})

	// telemetry.SubscribeStatic(1*time.Second, func(ag *types.AccStatic) {
	// 	fmt.Printf("%+v\n\n", ag)
	// })

	time.Sleep(100 * time.Second)
	fmt.Println("finishing test")
}

func TestOneHitRead(t *testing.T) {
	telemetry, err := acctelemetry.AccTelemetry()
	defer telemetry.Close()
	if err != nil {
		t.Error(fmt.Errorf("failed to create the telemetry: %w", err))
	}

	fmt.Printf("%+v\n\n", telemetry.ReadStatic())
	fmt.Printf("%+v\n\n", telemetry.ReadPhysics())
	fmt.Printf("%+v\n\n", telemetry.ReadGraphic())
}

func TestName(t *testing.T) {
	telemetry, err := acctelemetry.AccTelemetry()
	defer telemetry.Close()
	if err != nil {
		t.Error(fmt.Errorf("failed to create the telemetry: %w", err))
	}

	staticData := telemetry.ReadStatic()
	fmt.Println(convertToString(staticData.PlayerName[:]), convertToString(staticData.PlayerSurname[:]))
}

func convertToString(chars []uint16) string {
	var str string
	for _, val := range chars {
		str += string(rune(val))
	}
	return str
}
