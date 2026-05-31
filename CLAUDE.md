# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build

```bash
go build ./...
go vet ./...
```

There are no tests in this module. `go build ./...` is the primary verification.

## Module Purpose

`github.com/go-gost/core` defines every interface and shared type used across the GOST project. It contains **no implementations** — all concrete types live in the `x/` module. Go stdlib interfaces (`net.Conn`, `net.Listener`, `net.Addr`) are the only external dependencies (plus `crypto/tls` for TLS config references).

The `go.mod` has zero third-party dependencies (only `go 1.22` / `toolchain go1.22.2`).

## Key Interfaces & Their Relationships

GOST tunneling follows a layered request path. Each layer is defined here as an interface:

```
Service → Listener → Handler → chain.Router → chain.Route → chain.Transporter → (dialer.Dialer → connector.Connector)
```

### Top-level: Service

[service/service.go](service/service.go) — `Serve() error / Addr() net.Addr / Close() error`. Wraps one `listener.Listener` + one `handler.Handler`.

### Accept side: Listener and Handler

- **[listener.Listener](listener/listener.go)** — `Init(metadata) / Accept() (net.Conn, error) / Addr() / Close()`. Accepts inbound connections. Its `Options` carry a `chain.Router` for upstream routing.
- **[handler.Handler](handler/handler.go)** — `Init(metadata) / Handle(ctx, conn) error`. Processes an inbound connection: authenticates, selects a route via `Forwarder.Forward(hop)`, and proxies traffic. Its `Options` carry a `chain.Router`, bypass rules, auth, rate/traffic limiters, TLS config, and recorders.

### Routing: Chain, Router, Hop, Node, Transporter

- **[chain.Router](chain/router.go)** — `Options() / Dial(ctx, network, address) (net.Conn, error) / Bind(ctx, network, address) (net.Listener, error)`. The top-level router that both `Listener` and `Handler` use. Combines a `Chainer`, `resolver.Resolver`, `hosts.HostMapper`, recorders, retries, and timeouts.
- **[chain.Chainer](chain/chain.go)** — `Route(ctx, network, address) Route`. Selects a `Route` for the target address.
- **[chain.Route](chain/route.go)** — `Dial(ctx, network, address) (net.Conn, error) / Bind(ctx, network, address) (net.Listener, error) / Nodes()`. Represents one path through the proxy chain. `Dial` connects outward; `Bind` is for reverse tunnels.
- **[hop.Hop](hop/hop.go)** — `Select(ctx, SelectOptions) *chain.Node`. Represents a group of nodes that serve the same role. Uses a `selector.Selector` to pick one node (load balancing: round-robin, random, etc.). The `SelectOptions` carry HTTP-level context (host, method, path, headers, client IP) for routing decisions.
- **[chain.Node](chain/node.go)** — A single proxy hop with name, address, transporter, bypass rules, resolver, host mapper, TLS/HTTP settings, and a `failMarker` for health tracking. Nodes are `selector.Markable` — failed nodes get marked so selectors can skip them.
- **[chain.Transporter](chain/transport.go)** — `Dial / Handshake / Connect / Bind / Multiplex`. Combines a `dialer.Dialer` + `connector.Connector` for one node in the chain. Dial reaches the next proxy hop, handshake negotiates the transport protocol, connect reaches the destination through the hop, bind sets up reverse listeners.

### Outbound: Dialer and Connector

- **[dialer.Dialer](dialer/dialer.go)** — `Init / Dial(ctx, addr) (net.Conn, error)`. Dials from GOST to the next-hop proxy server. Also defines `Handshaker` (transport-level handshake, e.g. TLS) and `Multiplexer` (can this connection be multiplexed?).
- **[connector.Connector](connector/connector.go)** — `Init / Connect(ctx, conn, network, address) (net.Conn, error)`. Connects from the final GOST node to the actual destination. Also defines `Handshaker` (application-level, e.g. SOCKS handshake) and `Binder` (for reverse tunnels / BIND).

### Cross-cutting: Registry, Metadata, Selector

- **[registry.Registry[T]](registry/registry.go)** — Generic `Register / Unregister / IsRegistered / Get / GetAll`. Every component type in `x/` has its own registry instance. Registration happens via `init()` side-effects and blank imports.
- **[metadata.Metadata](metadata/metadata.go)** — `IsExists / Set / Get` key-value map. Every component has an `Init(Metadata)` method. Config fields that don't map to explicit struct fields go into metadata.
- **[selector.Selector[T]](selector/selector.go)** — Generic `Select(ctx, ...T) T`, plus `Strategy[T]`, `Filter[T]`, `Markable`/`Marker` for health tracking. Used by `Hop` for node selection.

### Support interfaces

| Interface | File | Purpose |
|-----------|------|---------|
| `auth.Authenticator` | [auth/auth.go](auth/auth.go) | `Authenticate(ctx, user, password) (id, ok)` |
| `admission.Admission` | [admission/admission.go](admission/admission.go) | `Admit(ctx, addr) bool` — allow/deny connection |
| `bypass.Bypass` | [bypass/bypass.go](bypass/bypass.go) | `IsWhitelist() / Contains(ctx, network, addr) bool` — should this address bypass the proxy chain? |
| `resolver.Resolver` | [resolver/resolver.go](resolver/resolver.go) | `Resolve(ctx, network, host) ([]net.IP, error)` |
| `hosts.HostMapper` | [hosts/hosts.go](hosts/hosts.go) | `Lookup(ctx, network, host) ([]net.IP, bool)` — static host→IP mapping |
| `recorder.Recorder` | [recorder/recorder.go](recorder/recorder.go) | `Record(ctx, b []byte) error` — traffic recording (hex dump, HTTP body capture) |
| `observer.Observer` | [observer/observer.go](observer/observer.go) | `Observe(ctx, events) error` — observability events |
| `ingress.Ingress` | [ingress/ingress.go](ingress/ingress.go) | `SetRule / GetRule(host)` — hostname→tunnel endpoint mapping |
| `sd.SD` | [sd/sd.go](sd/sd.go) | `Register / Deregister / Renew / Get` — service discovery |
| `router.Router` | [router/router.go](router/router.go) | `GetRoute(ctx, dst)` — OS-level route table queries (TUN mode) |
| `routing.Matcher` | [routing/matcher.go](routing/matcher.go) | `Match(*Request) bool` — request-based routing rules |
| `logger.Logger` | [logger/logger.go](logger/logger.go) | Structured logging with levels (trace→fatal). `LoggerGroup` fans out to multiple loggers. Global default via `SetDefault`/`Default`. |
| `metrics.Metrics` | [metrics/metrics.go](metrics/metrics.go) | `Counter / Gauge / Observer` — Prometheus-style metrics abstraction |
| `conn.Limiter` | [limiter/conn/limiter.go](limiter/conn/limiter.go) | `Allow(n int) bool` — connection count limiting |
| `rate.Limiter` | [limiter/rate/limiter.go](limiter/rate/limiter.go) | `Allow(n int) bool` — rate limiting |
| `traffic.Limiter` | [limiter/traffic/limiter.go](limiter/traffic/limiter.go) | `Wait(ctx, n) int` — bandwidth/traffic limiting |
| `limiter` (scopes) | [limiter/limiter.go](limiter/limiter.go) | Shared `Options` struct and scope constants (`ScopeService`, `ScopeConn`, `ScopeClient`) |
| `observer/stats.Stats` | [observer/stats/stats.go](observer/stats/stats.go) | Connection/bytes/errors counters for observability |
| `bufpool` | [common/bufpool/pool.go](common/bufpool/pool.go) | Tiered `sync.Pool` for byte buffers (128B → 64KB). `Get(size)` / `Put(b)`. |
| `xnet.Dialer` | [common/net/dialer.go](common/net/dialer.go) | Minimal `Dial(ctx, network, addr) (net.Conn, error)` — used as a shared dialer interface in connector/dialer options |

## Patterns

### Functional Options

Every interface uses the functional options pattern. Each package typically has two option layers:

1. **Init-time options** (e.g. `handler.Option`, `listener.Option`) — set in `x/config/parsing/service/parse.go` when constructing components from config.
2. **Call-time options** (e.g. `handler.HandleOption`, `dialer.DialOption`) — passed at invocation time with per-request context.

Both follow the same pattern: an `Options` struct, and option functions `type FooOption func(*FooOptions)`.

### No Implementations Rule

This module must never contain implementations of any interface it defines. Implementations belong in `x/`. Exceptions are small utility types that are self-contained and used by multiple `x/` packages: `logger.LoggerGroup`, `selector.failMarker`, `listener.AcceptError`/`BindError`, and `bufpool`.

### Adding a New Interface

1. Define the interface in a new or existing package under `core/`
2. Add any init-time/call-time option types in the same package
3. Register implementations in `x/registry/` using the registration pattern
4. Keep `go.mod` dependency-free (third-party types belong in `x/`, not here)

## Relationship to Workspace

This is one of ~12 modules in the `go-gost` workspace. See [../../CLAUDE.md](../../CLAUDE.md) for the full project architecture, build commands, and the implementation module (`x/`). The workspace `go.work` ties all modules together — `go build ./...` from the workspace root resolves across all modules.
