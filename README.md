# Donatello

Donatello is proxy tool.

## Client

### Client Flow

* Introduce the Wireguard port.
* Introduce the app port to the Wireguard.
* Introduce the remote server IP/port.
* Receive data from Wireguard.
* Encrypt the data.
* Send the data to the remote server.
* Receive back the response data.
* Decrypt the response data.
* Send the data to the Wireguard.

## Server

### Server Flow

* Introduce the Wireguard port.
* Introduce the app port to the Wireguard.
* Receive encrypted data.
* Decrypt the data.
* Send the data to the Wireguard.
* Get back the response data from the Wireguard.
* Encrypt the data.
* Send the data back to the client.
