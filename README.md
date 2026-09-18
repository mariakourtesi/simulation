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
├── main.go             entry point; runs one simulation, prints a report
├── erlangb/            Erlang-B (M/M/c/c) model components
│   ├── arrivals.go         exponential interarrival/service time draws
│   ├── events.go           min-heap of scheduled events (arrival/departure)
│   ├── formula.go          analytical Erlang-B recurrence B(c, A)
│   ├── simulation.go       the discrete-event simulation (RunSimulation)
│   └── *_test.go           tests for each of the above
└── README.md
```

## How the model works

The simulation keeps one state variable, the number of busy servers, and a
time-ordered event queue (a heap) of future arrivals and departures. Each
arrival is accepted if a server is free (and schedules its own departure) or
blocked if all servers are busy. Inter-event times are drawn from an
exponential distribution, which makes the arrival stream Poisson, the
assumption the Erlang-B formula is built on. Running long enough, the
simulated blocking probability converges to the analytical value, and
q(capacity) (fraction of time all servers are busy) converges to the same
number, demonstrating the PASTA property.
