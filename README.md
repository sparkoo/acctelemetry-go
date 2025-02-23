# acctelemetry-go
Simple Go library to read Assetto Corsa Competizione telemetry

Runs only on Windows!

## How to use

ACC must be running. Otherwise, `telemetry.Connect()` will fail with an error.

Methods `StaticPointer()`, `PhysicsPointer()` and `GraphicsPointer()` returns pointer to shared memory, so data will change over time. It's up to client code to create snapshot of the data if they wish so.

`RealtimeCarUpdate()` is read from UDP broadcast and it returns instance of latest `RealtimeCarUpdate` state.

### Example:

see [examples/main.go](examples/main.go)

### Configs:

see [telemetry_config.go](telemetry_config.go)
