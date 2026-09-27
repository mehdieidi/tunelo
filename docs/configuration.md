# Configuration and deployment

Tunelo has two executable packages. Build them from the repository root:

```sh
go build -o tunelo-server ./server
go build -o tunelo-client ./client
```

The server must be able to reach its local VPN UDP endpoint. The client must run where it can reach the local VPN application and the server's network listener. The original WireGuard and Linux NAT/IP-forwarding walkthrough remains in the root [README](../README.md).

## Client flags

| Flag | Default | Purpose |
| --- | --- | --- |
| `-server_ip` | `127.0.0.1` | Remote Tunelo server address. |
| `-server_port` | `23230` | Remote Tunelo server port. |
| `-vpn_port` | `23233` | Local VPN UDP port on the client host. |
| `-client_port` | `23231` | Loopback UDP port exposed for the local VPN application. |
| `-p` | `ws` | Transport: `ws`, `tcp`, or `utls`. |
| `-server_domain` | empty | TLS server name supplied to the client in `utls` mode. Required in that mode. |

The client binds its UDP listener to `127.0.0.1`, so it is available only on the same host. Point the local VPN peer endpoint at `127.0.0.1:<client_port>` and exclude the Tunelo server from the VPN route to avoid routing the tunnel back through itself. The root README contains a sample WireGuard configuration.

Example WebSocket launch:

```sh
./tunelo-server -server_ip 0.0.0.0 -server_port 23230 -vpn_port 23233 -p ws
./tunelo-client -server_ip <server-address> -server_port 23230 -vpn_port 23233 -client_port 23231 -p ws
```

Replace `<server-address>` with the reachable address or DNS name of the host running the server.

## Server flags

| Flag | Default | Purpose |
| --- | --- | --- |
| `-server_ip` | `127.0.0.1` | Local address on which the selected transport listener binds. |
| `-server_port` | `23230` | Local listener port. |
| `-vpn_port` | `23233` | Local VPN UDP port on the server host. |
| `-p` | `ws` | Transport: `ws`, `tcp`, or `utls`. |

For a remotely reachable listener, bind to an appropriate interface (commonly `0.0.0.0`) and allow the chosen port through the host firewall. The WebSocket server registers `/ws`; the TCP and TLS listeners accept direct connections on the configured address and port.

## TLS/uTLS mode

For `-p utls`, place `cert.pem` and `key.pem` in the server's current working directory. The client must set `-server_domain` to the server name used in its TLS configuration. The client currently sets `InsecureSkipVerify: true`, so the certificate is not validated; see [security notes](security.md).

The checked-in code does not configure a certificate path flag, certificate renewal, or WSS for WebSocket mode.
