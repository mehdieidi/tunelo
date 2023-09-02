# tunelo

__tunelo__ encrypts and tunnels UDP traffic (e.g. WireGuard) over a transport protocol like websocket or
TCP. Helping to use VPNs in restricted areas.

## WireGuard Client Config

You can use the WireGuard Android app to apply these settings. you need to replace the PrivateKey
and the PublicKey.

It needs to exclude the __tunelo__ server IP or exclude the __termux__ app in which the __tunelo__
client runs on.

The Endpoint in the peer part is pointing to the __tunelo__ client running on the phone.

``` WireGuard
[Interface]
Address = 10.8.0.2/32
DNS = 1.1.1.1
ExcludedApplications = com.termux
ListenPort = 23233
MTU = 1450
PrivateKey = 

[Peer]
AllowedIPs = 0.0.0.0/0
Endpoint = 127.0.0.1:23231
PersistentKeepalive = 25
PublicKey = D4PLMnAoDuXcgj7iTzyLs7NRptTND+z8vmxYA4Af218=
```

## WireGuard Server Config

This is the WireGuard config you should apply in the server that the __tunelo__ server will run.

Make sure to replace the private and public keys.

``` WireGuard
# Server
[Interface]
PrivateKey = 
Address = 10.8.0.1/24
MTU = 1450
ListenPort = 23233
SaveConfig = false
DNS = 1.1.1.1
DNS = 8.8.8.8

# Phone
[Peer]
PublicKey = r8KQuA7mtVVpHwDY6vTFmeMBcn+Y7omh6EPWroMWcD8=
AllowedIPs = 10.8.0.2/32
```

## Configure NAT rules on Linux

This iptables command is adding a NAT rule that masquerades (changes) the source IP address of
 packets originating from the specified source network (10.8.0.0/24) to match the public IP address
 of the router/firewall when those packets are leaving the system. This is commonly used in
 scenarios where you have a private network behind a NAT gateway or firewall, and you want the
internal devices to access the internet using the public IP address of the gateway.

``` bash
$ iptables -A POSTROUTING -t nat -s 10.8.0.0/24 -j MASQUERADE
...
```

The following command is used to enable IP forwarding on a Linux system. IP forwarding is a feature
 that allows a Linux system to route traffic between different network interfaces or subnets.

``` bash
$ echo 1 > /proc/sys/net/ipv4/ip_forward
...
```
