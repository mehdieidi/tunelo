# Development

Tunelo is a Go 1.20 module. The main dependency set includes uTLS, zerolog, `golang.org/x/crypto`, and `nhooyr.io/websocket`.

## Repository map

| Path | Responsibility |
| --- | --- |
| `client/` | Client executable, loopback UDP listener, and client-side transport implementations. |
| `server/` | Server executable, VPN-side UDP connection, and server-side transport listeners. |
| `pkg/xcrypto/` | ChaCha20-Poly1305 encrypt/decrypt helpers. |
| `pkg/logger/` | Logger interface and shared log metadata. |
| `pkg/logger/plain/` | Plain-text logger used by the current executables. |
| `pkg/logger/zerolog/` | JSON logger adapter. |
| `go.mod`, `go.sum` | Module and dependency lock data. |
| `Makefile` | `tidy`, `fmt`, and `lint` targets. |

## Useful commands

```sh
go build ./client ./server
go fmt ./...
go mod tidy
```

The Makefile's `lint` target expects `golangci-lint` and references `config/.golangci.yaml`; that config directory is not present in the current repository snapshot. Check the target configuration before relying on it.

There are no Go test files in the current repository. This documentation update does not add tests.
