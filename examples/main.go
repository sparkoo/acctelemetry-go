package main

import (
	"fmt"
	"log"
	"time"

	"github.com/sparkoo/acctelemetry-go"
)

func main() {
	telemetry := acctelemetry.New(acctelemetry.DefaultUdpConfig())
	err := telemetry.Connect()
	if err != nil {
		log.Fatalf("unable to connect to ACC: %s", err)
	}

	// we can run this in loop
	ticker := time.NewTicker(1 * time.Second)
	for _ = range ticker.C {
		fmt.Printf("%+v\n\n", telemetry.StaticPointer())
		fmt.Printf("%+v\n\n", telemetry.PhysicsPointer())
		fmt.Printf("%+v\n\n", telemetry.GraphicsPointer())
		fmt.Printf("%+v\n\n", telemetry.RealtimeCarUpdate())
	}
}
