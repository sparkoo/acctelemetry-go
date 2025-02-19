# acctelemetry-go
Simple Go library to read Assetto Corsa Competizione telemetry

Runs only on Windows!

## How to use

ACC must be running. Otherwise, `telemetry.Connect()` will fail with an error.

Methods `StaticPointer()`, `PhysicsPointer()` and `GraphicsPointer()` returns pointer to shared memory, so data will change over time. It's up to client code to create snapshot of the data if they wish so.

### Example:
```go
func main() {
  telemetry := acctelemetry.New()
  err := telemetry.Connect()
  if err != nil {
    t.Error(fmt.Errorf("unable to connect to ACC: %w", err))
  }

  // we can run this in loop
  ticker := time.NewTicker(1 * time.Second)
  for _ = range ticker.C {
    fmt.Printf("%+v\n\n", telemetry.StaticPointer())
    fmt.Printf("%+v\n\n", telemetry.PhysicsPointer())
    fmt.Printf("%+v\n\n", telemetry.GraphicsPointer())
  }
}
```
