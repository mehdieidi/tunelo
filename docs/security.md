# Security notes

This page describes what the checked-in implementation does today. Review these points before exposing a listener to an untrusted network or treating Tunelo as an authenticated secure tunnel.

## Transport protection

- `tcp` is plain TCP and provides no encryption or peer authentication.
- `ws` uses `ws://`, not `wss://`, in the client implementation. WebSocket framing alone does not encrypt traffic.
- `utls` wraps TCP in TLS and uses a Chrome 102 client hello, but the client sets `InsecureSkipVerify: true`. The connection therefore does not authenticate the server certificate or hostname. A TLS handshake is not sufficient protection against an active intermediary in this configuration.
- `pkg/xcrypto` contains ChaCha20-Poly1305 helpers, but the tunnel data path does not call them. Do not assume those helpers encrypt relayed VPN packets.

## Relay and exposure limits

- Both client UDP endpoints bind to loopback; the client-facing port is not exposed to other hosts by default.
- The server listener address is configurable. Binding it to a public interface exposes the selected transport listener to the network; the current code does not add client authentication or an access-control list.
- UDP and stream connections are bridged with `io.Copy` without an explicit datagram framing protocol. Stream read boundaries are not guaranteed to match original UDP packet boundaries, so packet behavior should be validated for the intended workload.
- The server shares one UDP connection to its configured VPN endpoint across accepted connections. The current implementation does not isolate peers into separate VPN-side UDP flows.

Use network firewall rules to restrict listener access. For production use, certificate validation, authenticated access, and explicit packet framing should be addressed in the implementation.
