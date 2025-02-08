# acctelemetry-go
Simple Go library to read Assetto Corsa Competizione telemetry

Runs only on Windows!

## How to use

First ACC must be running. Otherwise, it creating the `telemetry` will fail with error.

### Example:
```go
func main() {
  telemetry := acctelemetry.AccTelemetry()
	err := telemetry.Connect()
	if err != nil {
		t.Error(fmt.Errorf("unable to connect to ACC: %w", err))
	}

  fmt.Printf("%+v\n\n", telemetry.StaticPointer())
  fmt.Printf("%+v\n\n", telemetry.PhysicsPointer())
  fmt.Printf("%+v\n\n", telemetry.GraphicsPointer())
}
```
