# acctelemetry-go
Simple Go library to read Assetto Corsa Competizione telemetry

## Example:
```
telemetry, err := acctelemetry.AccTelemetry()
defer telemetry.Close()
if err != nil {
  t.Error(fmt.Errorf("failed to create the telemetry: %w", err))
}

telemetry.SubscribeStatic(1*time.Second, func(ag *types.AccStatic) {
  fmt.Printf("Static data: %+v\n\n", ag)
})

telemetry.SubscribeGraphic(1*time.Second, func(ag *types.AccGraphic) {
  fmt.Printf("Graphic data: %+v\n\n", ag)
})

telemetry.SubscribePhysics(1*time.Second, func(ag *types.AccPhysics) {
  fmt.Printf("Physics data: %+v\n\n", ag)
})
```
