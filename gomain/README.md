# gomain

A main harness for Go applications that run for a long time. This harness supports:

* Run as Windows Service.
* Handle OS signals gracefully.
* Basic debug signal handling.

## Usage

```go
import (
  "github.com/jeremyje/gomain"
)

func main() {
  gomain.Run(runServer, gomain.Config{
    ServiceName:        "App",
    ServiceDescription: "App does stuff.",
    Command:            *flagValue,
  })
}

func runServer(wait func()) error {
  server := New()
  go func() {
    wait()
    // Terminates the server and causes Run to complete.
    server.Shutdown()
  }()
  return server.Run()
}

```
