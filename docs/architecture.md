# Architecture

Tunelo is a pair of Go programs that relay UDP traffic between a local VPN endpoint and a remote VPN endpoint. The client exposes a loopback UDP port for the local VPN application, then carries data over a selected network transport to the server. The server relays received data to its configured loopback UDP VPN port.

Open the [interactive, accessible HTML architecture diagram](architecture.html) to see the two VPN endpoints, both Tunelo processes, and the selectable client-to-server transport. The same diagram is available as a standalone [SVG vector](architecture.svg) for embedding and editing.

## Data path

1. The client binds a UDP socket at `127.0.0.1:<client_port>` for the local VPN application.
2. The client opens a separate UDP connection to `127.0.0.1:<vpn_port>` and selects the configured transport.
3. The client relays data between its local UDP socket, the selected network connection, and its VPN UDP connection.
4. The server accepts the network connection and copies data between it and a UDP connection to `127.0.0.1:<vpn_port>`.
5. Replies travel back over the same connections.

The default ports are `23231` for the client-facing UDP listener and `23233` for each host's local VPN endpoint. The network listener/client defaults to port `23230`.

## Transport modes

| Mode | Client connection | Server listener | Notes |
| --- | --- | --- | --- |
| `ws` (default) | WebSocket at `ws://<server>:<port>/ws` | HTTP WebSocket upgrade at `/ws` | Uses binary WebSocket messages and does not enable WSS in the current code. |
| `tcp` | Plain TCP | TCP listener | Carries the relay over an unencrypted TCP stream. |
| `utls` | TCP with a uTLS `HelloChrome_102` client hello | TLS listener | Server loads `cert.pem` and `key.pem`; the client currently skips certificate verification and requires `server_domain`. |

Unknown values of `-p` currently fall through to WebSocket mode. Prefer one of the three documented values.

## Implementation boundaries

- Transport implementations live in matching files under `client/` and `server/`.
- UDP and stream connections are bridged with `io.Copy`; the code does not add an explicit datagram framing layer.
- Each server process creates one UDP connection to its configured VPN endpoint and shares it across accepted transport handlers.
- `pkg/xcrypto` exposes ChaCha20-Poly1305 helper functions, but the current client/server relay path does not call them.
- The active programs use the plain logger. A zerolog adapter also exists under `pkg/logger/zerolog`.

See [configuration](configuration.md) for the flags and [security notes](security.md) for implications of the current transport and relay behavior.
