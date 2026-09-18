# Simulation (Go) - Loss Models

A Go implementation of teletraffic loss models simulations

## Run

```bash
go run main.go
```

## Test

```bash
go test ./...
```

## Layout

```
simulation/
├── go.mod              module definition
├── main.go             entry point
├── erlangb/            Erlang-B model components (arrivals, service, servers, blocking)
│   ├── arrivals.go
│   └── arrivals_test.go
└── README.md
```
