# acctelemetry-go
Simple Go library to read Assetto Corsa Competizione telemetry

Runs only on Windows!

## How to use

First ACC must be running. Otherwise, it creating the `telemetry` will fail with error.
```go 
telemetry, err := acctelemetry.AccTelemetry()
```

Then it is possible to get actual data with:
```go
telemetry.ReadStatic()
telemetry.ReadPhysics()
telemetry.ReadGraphic()
```
These methods returns copy of the struct, so now it's yours.

Repetitive polling is up to the consumer.

## Full example:
```go
func main() {
  telemetry, err := acctelemetry.AccTelemetry()
  defer telemetry.Close()
  if err != nil {
    t.Error(fmt.Errorf("failed to create the telemetry: %w", err))
  }

  fmt.Printf("%+v\n\n", telemetry.ReadStatic())
  fmt.Printf("%+v\n\n", telemetry.ReadPhysics())
  fmt.Printf("%+v\n\n", telemetry.ReadGraphic())
}
```
