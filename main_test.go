package acctelemetry_test

import (
	"fmt"
	"testing"

	"github.com/sparkoo/acctelemetry-go"
)

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
